<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Viper](#viper)
  - [支持](#%E6%94%AF%E6%8C%81)
  - [初始化](#%E5%88%9D%E5%A7%8B%E5%8C%96)
  - [设置默认值](#%E8%AE%BE%E7%BD%AE%E9%BB%98%E8%AE%A4%E5%80%BC)
  - [读取配置文件](#%E8%AF%BB%E5%8F%96%E9%85%8D%E7%BD%AE%E6%96%87%E4%BB%B6)
  - [访问配置](#%E8%AE%BF%E9%97%AE%E9%85%8D%E7%BD%AE)
    - [直接访问](#%E7%9B%B4%E6%8E%A5%E8%AE%BF%E9%97%AE)
    - [反序列化到struct或map之中](#%E5%8F%8D%E5%BA%8F%E5%88%97%E5%8C%96%E5%88%B0struct%E6%88%96map%E4%B9%8B%E4%B8%AD)
  - [监听文件变化](#%E7%9B%91%E5%90%AC%E6%96%87%E4%BB%B6%E5%8F%98%E5%8C%96)
  - [参考](#%E5%8F%82%E8%80%83)

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

Note: 目前Viper支持的Remote远程读取配置如 etcd, consul；目前还没有对 Nacos 进行支持，参考第三方 https://github.com/yoyofxteam/nacos-viper-remote



## 初始化

```go
// github.com/spf13/viper@v1.21.0/viper.go

// 初始化实例
func New() *Viper {
	v := new(Viper)
	v.keyDelim = "." // 切割 key 
	v.configName = "config"
	v.configPermissions = os.FileMode(0o644)
	v.fs = afero.NewOsFs()  // 文件系统
	v.config = make(map[string]any) // 配置文件
	v.parents = []string{}
	v.override = make(map[string]any)
	v.defaults = make(map[string]any)
	v.kvstore = make(map[string]any)
	v.pflags = make(map[string]FlagValue)
	v.env = make(map[string][]string)
	v.aliases = make(map[string]string)
	v.typeByDefValue = false
	v.logger = slog.New(&discardHandler{})

	codecRegistry := NewCodecRegistry()

	v.encoderRegistry = codecRegistry
	v.decoderRegistry = codecRegistry

	// 实验特征
	v.experimentalFinder = features.Finder 
	v.experimentalBindStruct = features.BindStruct

	return v
}

```

## 设置默认值
```go
func (v *Viper) SetDefault(key string, value any) {
	// 如果是别名,找别名
	key = v.realKey(strings.ToLower(key))
	value = toCaseInsensitiveValue(value)

	path := strings.Split(key, v.keyDelim)
	lastKey := strings.ToLower(path[len(path)-1])
	deepestMap := deepSearch(v.defaults, path[0:len(path)-1]) // 前缀搜索

	// set innermost value
	deepestMap[lastKey] = value
}
```
## 读取配置文件
```go
func (v *Viper) ReadInConfig() error {
	// 获取文件名
	filename, err := v.getConfigFile()
	if err != nil {
		return err
	}

	if !slices.Contains(SupportedExts, v.getConfigType()) {
		return UnsupportedConfigError(v.getConfigType())
	}

	// 读取文件
	v.logger.Debug("reading file", "file", filename)
	file, err := afero.ReadFile(v.fs, filename)
	if err != nil {
		return err
	}

	//  解析文件内容反序列化
	config := make(map[string]any)

	err = v.unmarshalReader(bytes.NewReader(file), config)
	if err != nil {
		return err
	}

	v.config = config
	return nil
}

```







## 访问配置

### 直接访问

```json
{
  "mysql":{
    "db":"test"
  },
  "host":{
	  "address":"localhost"
	  "ports":[
		  "8080",
		  "8081"
	  ]
  }
}


```

```go
// 多层级配置key，可以用逗号隔号
viper.Get("host.address")//输出：localhost

// 数组，可以用序列号访问
viper.Get("host.posts.1")//输出: 8081


//也可以使用sub函数解析某个key的下级配置,如：
hostViper := viper.Sub("host")
fmt.Println(hostViper.Get("address"))
fmt.Println(hostViper.Get("posts.1"))
```

获取 key,viper 配置键不区分大小写
```go
func (v *Viper) Get(key string) any {
	// 转成小写
	lcaseKey := strings.ToLower(key)
	val := v.find(lcaseKey, true)
	if val == nil {
		return nil
	}

	if v.typeByDefValue {
		// TODO(bep) this branch isn't covered by a single test.
		valType := val
		path := strings.Split(lcaseKey, v.keyDelim)
		defVal := v.searchMap(v.defaults, path)
		if defVal != nil {
			valType = defVal
		}

		switch valType.(type) {
		case bool:
			return cast.ToBool(val)
		case string:
			return cast.ToString(val)
        //  ... 其他类型
		}
	}

	return val
}
```

```go
// Note: this assumes a lower-cased key given.
func (v *Viper) find(lcaseKey string, flagDefault bool) any {
	var (
		val    any
		exists bool
		path   = strings.Split(lcaseKey, v.keyDelim)
		nested = len(path) > 1
	)

	// compute the path through the nested maps to the nested value
	if nested && v.isPathShadowedInDeepMap(path, castMapStringToMapInterface(v.aliases)) != "" {
		return nil
	}

	// if the requested key is an alias, then return the proper key
	lcaseKey = v.realKey(lcaseKey)
	path = strings.Split(lcaseKey, v.keyDelim)
	nested = len(path) > 1

	// Set() override first
	val = v.searchMap(v.override, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.override) != "" {
		return nil
	}

	// PFlag override next
	flag, exists := v.pflags[lcaseKey]
	if exists && flag.HasChanged() {
		switch flag.ValueType() {
		case "int", "int8", "int16", "int32", "int64":
			return cast.ToInt(flag.ValueString())
		case "bool":
			return cast.ToBool(flag.ValueString())
		case "stringSlice", "stringArray":
			s := strings.TrimPrefix(flag.ValueString(), "[")
			s = strings.TrimSuffix(s, "]")
			res, _ := readAsCSV(s)
			return res
        // 其他类型
		
		default:
			return flag.ValueString()
		}
	}
	if nested && v.isPathShadowedInFlatMap(path, v.pflags) != "" {
		return nil
	}

	// Env override next
	if v.automaticEnvApplied {
		envKey := strings.Join(append(v.parents, lcaseKey), ".")
		// even if it hasn't been registered, if automaticEnv is used,
		// check any Get request
		if val, ok := v.getEnv(v.mergeWithEnvPrefix(envKey)); ok {
			return val
		}
		if nested && v.isPathShadowedInAutoEnv(path) != "" {
			return nil
		}
	}
	envkeys, exists := v.env[lcaseKey]
	if exists {
		for _, envkey := range envkeys {
			if val, ok := v.getEnv(envkey); ok {
				return val
			}
		}
	}
	if nested && v.isPathShadowedInFlatMap(path, v.env) != "" {
		return nil
	}

	// Config file next
	val = v.searchIndexableWithPathPrefixes(v.config, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.config) != "" {
		return nil
	}

	// K/V store next
	val = v.searchMap(v.kvstore, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.kvstore) != "" {
		return nil
	}

	// Default next
	val = v.searchMap(v.defaults, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.defaults) != "" {
		return nil
	}

	if flagDefault {
		// last chance: if no value is found and a flag does exist for the key,
		// get the flag's default value even if the flag's value has not been set.
		if flag, exists := v.pflags[lcaseKey]; exists {
			switch flag.ValueType() {
			case "int", "int8", "int16", "int32", "int64":
				return cast.ToInt(flag.ValueString())
			case "bool":
				return cast.ToBool(flag.ValueString())
			case "stringSlice", "stringArray":
				s := strings.TrimPrefix(flag.ValueString(), "[")
				s = strings.TrimSuffix(s, "]")
				res, _ := readAsCSV(s)
				return res
			case "boolSlice":
				s := strings.TrimPrefix(flag.ValueString(), "[")
				s = strings.TrimSuffix(s, "]")
				res, _ := readAsCSV(s)
				return cast.ToBoolSlice(res)
			case "intSlice":
				s := strings.TrimPrefix(flag.ValueString(), "[")
				s = strings.TrimSuffix(s, "]")
				res, _ := readAsCSV(s)
				return cast.ToIntSlice(res)
			case "uintSlice":
				s := strings.TrimPrefix(flag.ValueString(), "[")
				s = strings.TrimSuffix(s, "]")
				res, _ := readAsCSV(s)
				return cast.ToUintSlice(res)
			case "float64Slice":
				s := strings.TrimPrefix(flag.ValueString(), "[")
				s = strings.TrimSuffix(s, "]")
				res, _ := readAsCSV(s)
				return cast.ToFloat64Slice(res)
			case "stringToString":
				return stringToStringConv(flag.ValueString())
			case "stringToInt":
				return stringToIntConv(flag.ValueString())
			case "durationSlice":
				s := strings.TrimPrefix(flag.ValueString(), "[")
				s = strings.TrimSuffix(s, "]")
				slice := strings.Split(s, ",")
				return cast.ToDurationSlice(slice)
			default:
				return flag.ValueString()
			}
		}
		// last item, no need to check shadowing
	}

	return nil
}

```
viper 优先级从高到低如下：

* 显式调用 Set()
* flag
* env
* config
* key/value store
* default

### 反序列化到struct或map之中

```go
func (v *Viper) Unmarshal(rawVal any, opts ...DecoderConfigOption) error {
	// 获取所有的 key 
	keys := v.AllKeys()

	if v.experimentalBindStruct { // 实验特征
		// TODO: make this optional?
		structKeys, err := v.decodeStructKeys(rawVal, opts...)
		if err != nil {
			return err
		}

		keys = append(keys, structKeys...)
	}

	// TODO: struct keys should be enough?
	return decode(v.getSettings(keys), v.defaultDecoderConfig(rawVal, opts...))
}


func (v *Viper) AllKeys() []string {
	m := map[string]bool{}
	// add all paths, by order of descending priority to ensure correct shadowing
	m = v.flattenAndMergeMap(m, castMapStringToMapInterface(v.aliases), "")
	m = v.flattenAndMergeMap(m, v.override, "")
	m = v.mergeFlatMap(m, castMapFlagToMapInterface(v.pflags))
	m = v.mergeFlatMap(m, castMapStringSliceToMapInterface(v.env))
	m = v.flattenAndMergeMap(m, v.config, "")
	m = v.flattenAndMergeMap(m, v.kvstore, "")
	m = v.flattenAndMergeMap(m, v.defaults, "")

	// convert set of paths to list
	a := make([]string, 0, len(m))
	for x := range m {
		a = append(a, x)
	}
	return a
}

```





## 监听文件变化

Viper支持在运行时让应用程序实时读取配置文件

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


## 参考

- [配置解析神器viper使用详解](https://juejin.cn/post/7096416508054044685)