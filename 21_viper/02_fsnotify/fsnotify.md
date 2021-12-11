# fsnotify源码分析

## 构建约束(Build Constraints)
fsnotify是一个跨平台的库, 源码中既包含了linux平台的实现逻辑, 也包含了mac平台和windows平台的实现逻辑
```go
// +build linux,386 darwin,!cgo
```
上面这条注释不是普通的注释, 而是构建约束, 把它写在代码文件的顶部(package声明的上面), 会被编译器在编译时按照目标平台来判断是否编译进可执行文件中. 上面这行构建约束的意思是(linux AND 386) OR (darwin AND (NOT cgo)).

## 初始化
```go
// NewWatcher establishes a new watcher with the underlying OS and begins waiting for events.
func NewWatcher() (*Watcher, error) {
	// 主要创建文件描述符fd
	kq, err := kqueue()
	if err != nil {
		return nil, err
	}

	w := &Watcher{
		kq:              kq,
		watches:         make(map[string]int),
		dirFlags:        make(map[string]uint32),
		paths:           make(map[int]pathInfo),
		fileExists:      make(map[string]bool),
		externalWatches: make(map[string]bool),
		Events:          make(chan Event),
		Errors:          make(chan error),
		done:            make(chan struct{}),
	}

	go w.readEvents()
	return w, nil
}
```
Note: 文件描述符需要我们自己去读取, 所以我们就需要有某种轮训机制
```go
func (w *Watcher) readEvents() {
	eventBuffer := make([]unix.Kevent_t, 10)

loop:
	for {
		// See if there is a message on the "done" channel
		select {
		case <-w.done:
			break loop
		default:
		}

		// Get new events
		kevents, err := read(w.kq, eventBuffer, &keventWaitTime)
		// EINTR is okay, the syscall was interrupted before timeout expired.
		if err != nil && err != unix.EINTR {
			select {
			case w.Errors <- err:
			case <-w.done:
				break loop
			}
			continue
		}

		// Flush the events we received to the Events channel
		// 直到读完为止
		for len(kevents) > 0 {
			kevent := &kevents[0]
			watchfd := int(kevent.Ident)
			mask := uint32(kevent.Fflags)
			w.mu.Lock()
			path := w.paths[watchfd]
			w.mu.Unlock()
			event := newEvent(path.name, mask)

			if path.isDir && !(event.Op&Remove == Remove) {
				// Double check to make sure the directory exists. This can happen when
				// we do a rm -fr on a recursively watched folders and we receive a
				// modification event first but the folder has been deleted and later
				// receive the delete event
				if _, err := os.Lstat(event.Name); os.IsNotExist(err) {
					// mark is as delete event
					event.Op |= Remove
				}
			}

			if event.Op&Rename == Rename || event.Op&Remove == Remove {
				w.Remove(event.Name)
				w.mu.Lock()
				delete(w.fileExists, event.Name)
				w.mu.Unlock()
			}

			if path.isDir && event.Op&Write == Write && !(event.Op&Remove == Remove) {
				w.sendDirectoryChangeEvents(event.Name)
			} else {
				// Send the event on the Events channel.
				select {
				case w.Events <- event:
				case <-w.done:
					break loop
				}
			}

			if event.Op&Remove == Remove {
				// Look for a file that may have overwritten this.
				// For example, mv f1 f2 will delete f2, then create f2.
				if path.isDir {
					fileDir := filepath.Clean(event.Name)
					w.mu.Lock()
					_, found := w.watches[fileDir]
					w.mu.Unlock()
					if found {
						// make sure the directory exists before we watch for changes. When we
						// do a recursive watch and perform rm -fr, the parent directory might
						// have gone missing, ignore the missing directory and let the
						// upcoming delete event remove the watch from the parent directory.
						if _, err := os.Lstat(fileDir); err == nil {
							w.sendDirectoryChangeEvents(fileDir)
						}
					}
				} else {
					filePath := filepath.Clean(event.Name)
					if fileInfo, err := os.Lstat(filePath); err == nil {
						w.sendFileCreatedEventIfNew(filePath, fileInfo)
					}
				}
			}

			// Move to next event
			kevents = kevents[1:]
		}
	}

	// cleanup
	err := unix.Close(w.kq)
	if err != nil {
		// only way the previous loop breaks is if w.done was closed so we need to async send to w.Errors.
		select {
		case w.Errors <- err:
		default:
		}
	}
	close(w.Events)
	close(w.Errors)
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