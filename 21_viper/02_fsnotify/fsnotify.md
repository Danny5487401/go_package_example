# fsnotify

fsnotify ，就是封装了系统调用，用来监控文件事件的。当指定目录或者文件，发生了创建，删除，修改，重命名的事件，里面就能得到通知。


fsnotify是一个跨平台的库, 源码中既包含了linux平台的实现逻辑, 也包含了mac平台和windows平台的实现逻辑

构建约束(Build Constraints)
```go
// +build linux,386 darwin,!cgo
```
上面这条注释不是普通的注释, 而是构建约束, 把它写在代码文件的顶部(package声明的上面), 会被编译器在编译时按照目标平台来判断是否编译进可执行文件中. 上面这行构建约束的意思是(linux AND 386) OR (darwin AND (NOT cgo)).

fsnotify 本质上就是对系统能力的一个浅层封装，主要封装了操作系统提供的两个机制：

- inotify 机制；
- epoll 机制；

## 需求
- 监听一个目录中所有文件，文件大小到一定阀值，则处理；
- 监控某个目录，当有文件新增，立马处理；
- 监控某个目录或文件，当有文件被修改或者删除，立马能感知，进行处理；

## 做法
1. 第一种：当事人主动通知你，这是侵入式的，需要当事人修改这部分代码来支持，依赖于当事人的自觉；
2. 第二种：轮询观察，这个是无侵入式的，你可以自己写个轮询程序，每隔一段时间唤醒一次，对文件和目录做各种判断，从而得到这个目录的变化；
3. 第三种：操作系统支持，以事件的方式通知到订阅这个事件的用户，达到及时处理的目的；

第三种最好：

- 纯旁路的逻辑，对线上程序无侵入；
- 操作系统直接支持，以事件的形式通知，性能也最好，100% 准确率（比较自己轮询判断要好）；




## inotify机制

这是一个内核用于通知用户空间程序文件系统变化的机制。

其实 inotify 机制的诞生源于一个通用的需求，由于IO/硬件管理都在内核，但用户态是有获悉内核事件的强烈需求，比如磁盘的热插拔，文件的增删改。这里就诞生了三个异曲同工的机制：hotplug 机制、udev 管理机制、inotify 机制.

### inotify 的三个接口
```C
// fs/notify/inotify/inotify_user.c

// 创建 notify fd
inotify_init1

// 添加监控路径
inotify_add_watch

// 删除一个监控
inotify_rm_watch
```

### inotify 怎么实现监控的？
inotify 支持监听的事件非常多，除了增删改，还有访问，移动，打开，关闭，设备卸载等等事件。

内核要上报这些文件 api 事件必然要采集这些事件。在哪一个内核层次采集的呢？
```css
系统调用 -> vfs -> 具体文件系统（ ext4 ）-> 块层 -> scsi 层
```

答案是：vfs （virtual File System）层。因为这是所有“文件”操作的入口。

以 vfs 的 read/write 为例
```C
ssize_t vfs_read(struct file *file, char __user *buf, size_t count, loff_t *pos)
{
    // ...
    ret = __vfs_read(file, buf, count, pos);
    if (ret > 0) {
        // 事件采集点：访问事件
        fsnotify_access(file);
    }

}

ssize_t vfs_write(struct file *file, const char __user *buf, size_t count, loff_t *pos)
{
    // ...
    ret = __vfs_write(file, buf, count, pos);
    if (ret > 0) {
        // 事件采集点：修改事件
        fsnotify_modify(file);
    }
}
```
fsnotify_access 和 fsnotify_modify 就是 inotify 机制的一员。有一系列 fsnotify_xxx 的函数，定义在 include/linux/fsnotify.h ，这函数里面全都调用到 fsnotify 这个函数。
```C
static inline void fsnotify_modify(struct file *file)
{
    // 获取到 inode

    if (!(file->f_mode & FMODE_NONOTIFY)) {
        fsnotify_parent(path, NULL, mask);
        // 采集事件，通知到指定结构
        fsnotify(inode, mask, path, FSNOTIFY_EVENT_PATH, NULL, 0);
    }
}
```


## watcher结构体
```go
// Watcher watches a set of files, delivering events to a channel.
type Watcher struct {
	Events chan Event
	Errors chan error
	done   chan struct{} // Channel for sending a "quit message" to the reader goroutine

	kq int // 文件描述符 (as returned by the kqueue() syscall).

	mu              sync.Mutex        // Protects access to watcher data
	watches         map[string]int    // Map of watched file descriptors (key: path).
	externalWatches map[string]bool   // Map of watches added by user of the library.
	dirFlags        map[string]uint32 // Map of watched directories to fflags used in kqueue.
	paths           map[int]pathInfo  // Map file descriptors to path names for processing kqueue events.
	fileExists      map[string]bool   // Keep track of if we know this file exists (to stop duplicate create events).
	isClosed        bool              // Set to true when Close() is first called
}
```

## 初始化
- linux: inoyify.go
- darwin: kqueue.go

```go
func NewWatcher() (*Watcher, error) {
	// 调用InotifyInit创建notify fd
	fd, errno := unix.InotifyInit1(unix.IN_CLOEXEC)
	if fd == -1 {
		return nil, errno
	}
	// Create epoll
	poller, err := newFdPoller(fd)
	if err != nil {
		unix.Close(fd)
		return nil, err
	}
	w := &Watcher{
		fd:       fd,
		poller:   poller,
		watches:  make(map[string]*watch),
		paths:    make(map[int]string),
		Events:   make(chan Event),
		Errors:   make(chan error),
		done:     make(chan struct{}),
		doneResp: make(chan struct{}),
	}

	go w.readEvents()
	return w, nil
}
```
Note: 文件描述符需要我们自己去读取, 所以我们就需要有某种轮训机制
```go
func (w *Watcher) readEvents() {
	var (
		buf   [unix.SizeofInotifyEvent * 4096]byte // Buffer for a maximum of 4096 raw events
		n     int                                  // Number of bytes read with read()
		errno error                                // Syscall errno
		ok    bool                                 // For poller.wait
	)

	defer close(w.doneResp)
	defer close(w.Errors)
	defer close(w.Events)
	defer unix.Close(w.fd)
	defer w.poller.close()

	for {
		// See if we have been closed.
		if w.isClosed() {
			return
		}

		ok, errno = w.poller.wait()
		if errno != nil {
			select {
			case w.Errors <- errno:
			case <-w.done:
				return
			}
			continue
		}

		if !ok {
			continue
		}

		n, errno = unix.Read(w.fd, buf[:])
		// If a signal interrupted execution, see if we've been asked to close, and try again.
		// http://man7.org/linux/man-pages/man7/signal.7.html :
		// "Before Linux 3.8, reads from an inotify(7) file descriptor were not restartable"
		if errno == unix.EINTR {
			continue
		}

		// unix.Read might have been woken up by Close. If so, we're done.
		if w.isClosed() {
			return
		}

		if n < unix.SizeofInotifyEvent {
			var err error
			if n == 0 {
				// If EOF is received. This should really never happen.
				err = io.EOF
			} else if n < 0 {
				// If an error occurred while reading.
				err = errno
			} else {
				// Read was too short.
				err = errors.New("notify: short read in readEvents()")
			}
			select {
			case w.Errors <- err:
			case <-w.done:
				return
			}
			continue
		}

		var offset uint32
		// We don't know how many events we just read into the buffer
		// While the offset points to at least one whole event...
		for offset <= uint32(n-unix.SizeofInotifyEvent) {
			// Point "raw" to the event in the buffer
			raw := (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))

			mask := uint32(raw.Mask)
			nameLen := uint32(raw.Len)

			if mask&unix.IN_Q_OVERFLOW != 0 {
				select {
				case w.Errors <- ErrEventOverflow:
				case <-w.done:
					return
				}
			}

			// If the event happened to the watched directory or the watched file, the kernel
			// doesn't append the filename to the event, but we would like to always fill the
			// the "Name" field with a valid filename. We retrieve the path of the watch from
			// the "paths" map.
			w.mu.Lock()
			name, ok := w.paths[int(raw.Wd)]
			// IN_DELETE_SELF occurs when the file/directory being watched is removed.
			// This is a sign to clean up the maps, otherwise we are no longer in sync
			// with the inotify kernel state which has already deleted the watch
			// automatically.
			if ok && mask&unix.IN_DELETE_SELF == unix.IN_DELETE_SELF {
				delete(w.paths, int(raw.Wd))
				delete(w.watches, name)
			}
			w.mu.Unlock()

			if nameLen > 0 {
				// Point "bytes" at the first byte of the filename
				bytes := (*[unix.PathMax]byte)(unsafe.Pointer(&buf[offset+unix.SizeofInotifyEvent]))[:nameLen:nameLen]
				// The filename is padded with NULL bytes. TrimRight() gets rid of those.
				name += "/" + strings.TrimRight(string(bytes[0:nameLen]), "\000")
			}

			event := newEvent(name, mask)

			// Send the events that are not ignored on the events channel
			if !event.ignoreLinux(mask) {
				select {
				case w.Events <- event:
				case <-w.done:
					return
				}
			}

			// Move to the next event in the buffer
			offset += unix.SizeofInotifyEvent + nameLen
		}
	}
}
```

## 事件
```go
type Event struct {
	Name string // Name表示发生变化的文件或目录名
	Op   Op     // 具体的变化
}
```
具体的变化
```go
// 事件中的Op是按照位来存储的，可以存储多个，可以通过&操作判断对应事件是不是发生了
// Op describes a set of file operations.
type Op uint32

// These are the generalized file operations that can trigger a notification.
const (
	Create Op = 1 << iota
	Write
	Remove
	Rename
	Chmod
)

```
查看行为
```go
func (op Op) String() string {
    // Use a buffer for efficient string concatenation
    var buffer bytes.Buffer
    
    if op&Create == Create {
        buffer.WriteString("|CREATE")
    }
    if op&Remove == Remove {
        buffer.WriteString("|REMOVE")
    }
    if op&Write == Write {
       buffer.WriteString("|WRITE")
    }
    if op&Rename == Rename {
       buffer.WriteString("|RENAME")
    }
    if op&Chmod == Chmod {
      buffer.WriteString("|CHMOD")
    }
    if buffer.Len() == 0 {
      return ""
    }
    return buffer.String()[1:] // Strip leading pipe
}
```


添加watch路径
```go
func (w *Watcher) Add(name string) error {
	w.mu.Lock()
	w.externalWatches[name] = true
	w.mu.Unlock()
	_, err := w.addWatch(name, noteAllEvents)
	return err
}

func (w *Watcher) addWatch(name string, flags uint32) (string, error) {
	var isDir bool
	// 获取标准路径, 如/tmp//////too经过Clean后就成了/tmp/too
	// Make ./name and name equivalent
	name = filepath.Clean(name)

	w.mu.Lock()
	if w.isClosed {
		w.mu.Unlock()
		return "", errors.New("kevent instance already closed")
	}
    // 取出上下文里的watch路径(如果存在的话)
	watchfd, alreadyWatching := w.watches[name]
	// We already have a watch, but we can still override flags.
	if alreadyWatching {
		isDir = w.paths[watchfd].isDir
	}
	w.mu.Unlock()

	if !alreadyWatching {
		fi, err := os.Lstat(name)
		if err != nil {
			return "", err
		}

		// Don't watch sockets.
		if fi.Mode()&os.ModeSocket == os.ModeSocket {
			return "", nil
		}

		// Don't watch named pipes.
		if fi.Mode()&os.ModeNamedPipe == os.ModeNamedPipe {
			return "", nil
		}

		// Follow Symlinks
		// Unfortunately, Linux can add bogus symlinks to watch list without
		// issue, and Windows can't do symlinks period (AFAIK). To  maintain
		// consistency, we will act like everything is fine. There will simply
		// be no file events for broken symlinks.
		// Hence the returns of nil on errors.
		if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
			name, err = filepath.EvalSymlinks(name)
			if err != nil {
				return "", nil
			}

			w.mu.Lock()
			_, alreadyWatching = w.watches[name]
			w.mu.Unlock()

			if alreadyWatching {
				return name, nil
			}

			fi, err = os.Lstat(name)
			if err != nil {
				return "", nil
			}
		}

		watchfd, err = unix.Open(name, openMode, 0700)
		if watchfd == -1 {
			return "", err
		}

		isDir = fi.IsDir()
	}

	const registerAdd = unix.EV_ADD | unix.EV_CLEAR | unix.EV_ENABLE
	if err := register(w.kq, []int{watchfd}, registerAdd, flags); err != nil {
		unix.Close(watchfd)
		return "", err
	}

	if !alreadyWatching {
        // 如果上下文里不存在此路径, 表明这是一个新的watch, 添加到上下文
		w.mu.Lock()
		w.watches[name] = watchfd
		w.paths[watchfd] = pathInfo{name: name, isDir: isDir}
		w.mu.Unlock()
	}

	if isDir {
		// Watch the directory if it has not been watched before,
		// or if it was watched before, but perhaps only a NOTE_DELETE (watchDirectoryFiles)
		w.mu.Lock()

		watchDir := (flags&unix.NOTE_WRITE) == unix.NOTE_WRITE &&
			(!alreadyWatching || (w.dirFlags[name]&unix.NOTE_WRITE) != unix.NOTE_WRITE)
		// Store flags so this watch can be updated later
		w.dirFlags[name] = flags
		w.mu.Unlock()

		if watchDir {
			if err := w.watchDirectoryFiles(name); err != nil {
				return "", err
			}
		}
	}
	return name, nil
}
```

## 添加文件或则目录监听

```go
func (w *Watcher) Add(name string) error {
	name = filepath.Clean(name)
	if w.isClosed() {
		return errors.New("inotify instance already closed")
	}

	const agnosticEvents = unix.IN_MOVED_TO | unix.IN_MOVED_FROM |
		unix.IN_CREATE | unix.IN_ATTRIB | unix.IN_MODIFY |
		unix.IN_MOVE_SELF | unix.IN_DELETE | unix.IN_DELETE_SELF

	var flags uint32 = agnosticEvents

	w.mu.Lock()
	defer w.mu.Unlock()
	watchEntry := w.watches[name]
	if watchEntry != nil {
		flags |= watchEntry.flags | unix.IN_MASK_ADD
	}
	wd, errno := unix.InotifyAddWatch(w.fd, name, flags)
	if wd == -1 {
		return errno
	}

	if watchEntry == nil {
		w.watches[name] = &watch{wd: uint32(wd), flags: flags}
		w.paths[wd] = name
	} else {
		watchEntry.wd = uint32(wd)
		watchEntry.flags = flags
	}

	return nil
}

```
```go
func InotifyAddWatch(fd int, pathname string, mask uint32) (watchdesc int, err error) {
	var _p0 *byte
	_p0, err = BytePtrFromString(pathname)
	if err != nil {
		return
	}
	// 调用inotify的添加监控路径
	r0, _, e1 := Syscall(SYS_INOTIFY_ADD_WATCH, uintptr(fd), uintptr(unsafe.Pointer(_p0)), uintptr(mask))
	watchdesc = int(r0)
	if e1 != 0 {
		err = errnoErr(e1)
	}
	return
}
```