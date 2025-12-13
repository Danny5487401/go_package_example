<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Viper 远程配置](#viper-%E8%BF%9C%E7%A8%8B%E9%85%8D%E7%BD%AE)
  - [默认支持的插件](#%E9%BB%98%E8%AE%A4%E6%94%AF%E6%8C%81%E7%9A%84%E6%8F%92%E4%BB%B6)
  - [RemoteProvider 接口](#remoteprovider-%E6%8E%A5%E5%8F%A3)
    - [nacos-viper 插件](#nacos-viper-%E6%8F%92%E4%BB%B6)
  - [加载远程配置](#%E5%8A%A0%E8%BD%BD%E8%BF%9C%E7%A8%8B%E9%85%8D%E7%BD%AE)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Viper 远程配置

## 默认支持的插件
```go
// github.com/spf13/viper@v1.8.1/viper.go

// 默认支持的远程插件
var SupportedRemoteProviders = []string{"etcd", "consul", "firestore"}
```



## RemoteProvider 接口
```go
type RemoteProvider interface {
	Provider() string
	Endpoint() string
	Path() string
	SecretKeyring() string
}
```

### nacos-viper 插件

```go
// github.com/yoyofxteam/nacos-viper-remote@v0.4.0/nacosprovider.go
type nacosRemoteProvider struct {
	provider      string
	endpoint      string
	path          string
	secretKeyring string
}

func DefaultRemoteProvider() *nacosRemoteProvider {
	return &nacosRemoteProvider{provider: "nacos", endpoint: "localhost", path: "", secretKeyring: ""}
}

func (rp nacosRemoteProvider) Provider() string {
	return rp.provider
}

func (rp nacosRemoteProvider) Endpoint() string {
	return rp.endpoint
}

func (rp nacosRemoteProvider) Path() string {
	return rp.path
}

func (rp nacosRemoteProvider) SecretKeyring() string {
	return rp.secretKeyring
}

```

```go
func SetOptions(option *Option) {
	manager, _ := NewNacosConfigManager(option)
	// 覆盖默认支持的远程插件
	viper.SupportedRemoteProviders = []string{"nacos"}
	viper.RemoteConfig = &remoteConfigProvider{ConfigManager: manager}
}

```




## 加载远程配置
```go
func (v *Viper) ReadRemoteConfig() error {
	return v.getKeyValueConfig()
}

func (v *Viper) getKeyValueConfig() error {
    // 校验

	for _, rp := range v.remoteProviders {
		// 拿到一个远程配置即可 
		val, err := v.getRemoteConfig(rp)
		if err != nil {
			jww.ERROR.Printf("get remote config: %s", err)

			continue
		}

		v.kvstore = val

		return nil
	}
	return RemoteConfigError("No Files Found")
}


// 从注册的RemoteProvider获取配置
func (v *Viper) getRemoteConfig(provider RemoteProvider) (map[string]interface{}, error) {
	reader, err := RemoteConfig.Get(provider)
	if err != nil {
		return nil, err
	}
	// 把数据写入v.kvstore
	err = v.unmarshalReader(reader, v.kvstore)
	return v.kvstore, err
}
```
从注册工厂中获取nacos服务

注册工厂需要实现的方法
```go
// github.com/spf13/viper@v1.8.1/viper.go
type remoteConfigFactory interface {
	Get(rp RemoteProvider) (io.Reader, error)
	Watch(rp RemoteProvider) (io.Reader, error)
	WatchChannel(rp RemoteProvider) (<-chan *RemoteResponse, chan bool)
}
// RemoteConfig is optional, see the remote package
var RemoteConfig remoteConfigFactory
```


默认工厂
```go
// github.com/spf13/viper@v1.8.1/remote/remote.go
func init() {
    viper.RemoteConfig = &remoteConfigProvider{}
}

func (rc remoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	cm, err := getConfigManager(rp)
	if err != nil {
		return nil, err
	}
	b, err := cm.Get(rp.Path())
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
```

代码中我们注册nacos工厂
```go
// 注册nacos配置
remote.SetOptions(&remote.Option{
    Url:         "tencent.danny.games",
    Port:        8848,
    NamespaceId: "public",
    GroupName:   "DEFAULT_GROUP",
    Config:      remote.Config{DataId: "config_dev.yaml"},
    Auth: &remote.Auth{User: "nacos",
        Password: "nacos"},
})
```

```go
// github.com/yoyofxteam/nacos-viper-remote@v0.4.0/viper_remote.go
func SetOptions(option *Option) {
	manager, _ := NewNacosConfigManager(option)
	viper.SupportedRemoteProviders = []string{"nacos"}
	// 注册nacos manager
	viper.RemoteConfig = &remoteConfigProvider{ConfigManager: manager}
}

type remoteConfigProvider struct {
	ConfigManager *nacosConfigManager
}

// 读取配置
func (rc *remoteConfigProvider) Get(rp viper.RemoteProvider) (io.Reader, error) {
	cmt, err := rc.getConfigManager(rp)
	if err != nil {
		return nil, err
	}
	var b []byte
	switch cm := cmt.(type) {
	case viperConfigManager:
		b, err = cm.Get(rp.Path())
	}
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func (rc *remoteConfigProvider) Watch(rp viper.RemoteProvider) (io.Reader, error) {
	return rc.Get(rp)
}

func (rc *remoteConfigProvider) WatchChannel(rp viper.RemoteProvider) (<-chan *viper.RemoteResponse, chan bool) {
	cmt, err := rc.getConfigManager(rp)
	if err != nil {
		return nil, nil
	}

	switch cm := cmt.(type) {
	case viperConfigManager:
		quit := make(chan bool)
		viperResponseCh := cm.Watch("dataId", quit)
		return viperResponseCh, quit
	}

	return nil, nil
}

func (rc *remoteConfigProvider) getConfigManager(rp viper.RemoteProvider) (interface{}, error) {
	if rp.Provider() == "nacos" {
		return rc.ConfigManager, nil
	} else {
		return nil, errors.New("The Nacos configuration manager is not supported!")
	}
}
```
实际拉取配置
```go
func (cm *nacosConfigManager) Get(dataId string) ([]byte, error) {
	//get config
	content, err := cm.client.GetConfig(vo.ConfigParam{
		DataId: cm.option.Config.DataId,
		Group:  cm.option.GroupName,
	})
	return []byte(content), err
}
```

寻找数值
```go
//
// Viper will check to see if an alias exists first.
// Viper will then check in the following order:
// flag, env, config file, key/value store.
// Lastly, if no value was found and flagDefault is true, and if the key
// corresponds to a flag, the flag's default value is returned.
//
// Note: this assumes a lower-cased key given.
func (v *Viper) find(lcaseKey string, flagDefault bool) interface{} {
	var (
		val    interface{}
		exists bool
		path   = strings.Split(lcaseKey, v.keyDelim)
		nested = len(path) > 1
	)

	// compute the path through the nested maps to the nested value
	if nested && v.isPathShadowedInDeepMap(path, castMapStringToMapInterface(v.aliases)) != "" {
		return nil
	}

	// 如果是别名，先找原来的数值
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
		case "intSlice":
			s := strings.TrimPrefix(flag.ValueString(), "[")
			s = strings.TrimSuffix(s, "]")
			res, _ := readAsCSV(s)
			return cast.ToIntSlice(res)
		case "stringToString":
			return stringToStringConv(flag.ValueString())
		default:
			return flag.ValueString()
		}
	}
	if nested && v.isPathShadowedInFlatMap(path, v.pflags) != "" {
		return nil
	}

	// 环境变量
	if v.automaticEnvApplied {
		// even if it hasn't been registered, if automaticEnv is used,
		// check any Get request
		if val, ok := v.getEnv(v.mergeWithEnvPrefix(lcaseKey)); ok {
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

	// 配置文件当中
	val = v.searchIndexableWithPathPrefixes(v.config, path)
	if val != nil {
		return val
	}
	if nested && v.isPathShadowedInDeepMap(path, v.config) != "" {
		return nil
	}

	// K/V store next  --刚刚nacos远程配置信息就是在这
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
			case "intSlice":
				s := strings.TrimPrefix(flag.ValueString(), "[")
				s = strings.TrimSuffix(s, "]")
				res, _ := readAsCSV(s)
				return cast.ToIntSlice(res)
			case "stringToString":
				return stringToStringConv(flag.ValueString())
			default:
				return flag.ValueString()
			}
		}
		// last item, no need to check shadowing
	}

	return nil
}
```

## 参考



