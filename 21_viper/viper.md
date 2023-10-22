<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Viper](#viper)
  - [支持](#%E6%94%AF%E6%8C%81)
  - [监听文件变化](#%E7%9B%91%E5%90%AC%E6%96%87%E4%BB%B6%E5%8F%98%E5%8C%96)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Viper
Viper是Go应用程序的完整配置解决方案，包括12-Factor应用程序。
## 支持
- 设置默认值
- 从JSON，TOML，YAML，HCL和Java属性配置文件中读取
- 实时观看和重新读取配置文件（可选）
- 从环境变量中读取
- 从远程配置系统（etcd或Consul）读取，并观察变化
- 从命令行标志读取
- 从缓冲区读取
- 设置显式值

Note: 目前Viper支持的Remote远程读取配置如 etcd, consul；目前还没有对Nacos进行支持，参考第三方https://github.com/yoyofxteam/nacos-viper-remote
## 监听文件变化
```go
func (v *Viper) WatchConfig() {
	// 保证监听器初始化
	initWG := sync.WaitGroup{}
	initWG.Add(1)
	go func() {
		// 调用NewWatcher创建一个监听器
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()
		//获取配置文件路径，抽出文件名、目录，配置文件如果是一个符号链接，获得链接指向的路径
		// we have to watch the entire directory to pick up renames/atomic saves in a cross-platform way
		filename, err := v.getConfigFile()
		if err != nil {
			log.Printf("error: %v\n", err)
			initWG.Done()
			return
		}

		configFile := filepath.Clean(filename)
		configDir, _ := filepath.Split(configFile)
		realConfigFile, _ := filepath.EvalSymlinks(filename)

		// eventsWG是在事件通道关闭，或配置被删除了，或遇到错误时退出事件处理循环
		eventsWG := sync.WaitGroup{}
		eventsWG.Add(1)
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok { // 'Events' channel is closed
						eventsWG.Done()
						return
					}
					currentConfigFile, _ := filepath.EvalSymlinks(filename)
					// we only care about the config file with the following cases:
					// 1 - if the config file was modified or created
					// 2 - if the real path to the config file changed (eg: k8s ConfigMap replacement)
					const writeOrCreateMask = fsnotify.Write | fsnotify.Create
					if (filepath.Clean(event.Name) == configFile &&
						event.Op&writeOrCreateMask != 0) ||
						(currentConfigFile != "" && currentConfigFile != realConfigFile) {
						realConfigFile = currentConfigFile
						err := v.ReadInConfig()
						if err != nil {
							log.Printf("error reading config file: %v\n", err)
						}
						if v.onConfigChange != nil {
							v.onConfigChange(event)
						}
					} else if filepath.Clean(event.Name) == configFile &&
						event.Op&fsnotify.Remove&fsnotify.Remove != 0 {
						eventsWG.Done()
						return
					}

				case err, ok := <-watcher.Errors:
					if ok { // 'Errors' channel is not closed
						log.Printf("watcher error: %v\n", err)
					}
					eventsWG.Done()
					return
				}
			}
		}()
		// 监听配置文件所在目录
		watcher.Add(configDir)
		initWG.Done()   // done initializing the watch in this go routine, so the parent routine can move on...
		eventsWG.Wait() // now, wait for event loop to end in this go-routine...
	}()
	initWG.Wait() // make sure that the go routine above fully ended before returning
}
```