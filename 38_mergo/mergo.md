<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/imdario/mergo](#githubcomimdariomergo)
  - [源码分析](#%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
    - [Merge](#merge)
    - [Map](#map)
  - [应用](#%E5%BA%94%E7%94%A8)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/imdario/mergo

mergo 可以在相同的结构体或map之间赋值，可以将结构体的字段赋值到map中，可以将map的值赋值给结构体的字段


## 源码分析

配置

```go
// github.com/imdario/mergo@v0.3.8/merge.go
type Config struct {
	Overwrite                    bool
	AppendSlice                  bool
	TypeCheck                    bool
	Transformers                 Transformers
	overwriteWithEmptyValue      bool
	overwriteSliceWithEmptyValue bool
}

```

### Merge

```go
func Merge(dst, src interface{}, opts ...func(*Config)) error {
	return merge(dst, src, opts...)
}


func merge(dst, src interface{}, opts ...func(*Config)) error {
	var (
		vDst, vSrc reflect.Value
		err        error
	)
    
	
	config := &Config{}
    // 配置选项
	for _, opt := range opts {
		opt(config)
	}

	// 校验要求 dst 是 struct 或则 map
	if vDst, vSrc, err = resolveValues(dst, src); err != nil {
		return err
	}
	// 校验类型相同
	if vDst.Type() != vSrc.Type() {
		return ErrDifferentArgumentsTypes
	}
	return deepMerge(vDst, vSrc, make(map[uintptr]*visit), 0, config)
}


func deepMerge(dst, src reflect.Value, visited map[uintptr]*visit, depth int, config *Config) (err error) {
	overwrite := config.Overwrite
	typeCheck := config.TypeCheck
	overwriteWithEmptySrc := config.overwriteWithEmptyValue
	overwriteSliceWithEmptySrc := config.overwriteSliceWithEmptyValue
	config.overwriteWithEmptyValue = false

	if !src.IsValid() {
		return
	}
	if dst.CanAddr() {
		addr := dst.UnsafeAddr()
		h := 17 * addr
		seen := visited[h]
		typ := dst.Type()
		for p := seen; p != nil; p = p.next {
			if p.ptr == addr && p.typ == typ {
				return nil
			}
		}
		// Remember, remember...
		visited[h] = &visit{addr, typ, seen}
	}

	if config.Transformers != nil && !isEmptyValue(dst) {
		if fn := config.Transformers.Transformer(dst.Type()); fn != nil {
			err = fn(dst, src)
			return
		}
	}

	switch dst.Kind() {
	case reflect.Struct:
		if hasExportedField(dst) {
			for i, n := 0, dst.NumField(); i < n; i++ {
				// 逐个字段递归
				if err = deepMerge(dst.Field(i), src.Field(i), visited, depth+1, config); err != nil {
					return
				}
			}
		} else {
			if dst.CanSet() && (!isEmptyValue(src) || overwriteWithEmptySrc) && (overwrite || isEmptyValue(dst)) {
				dst.Set(src)
			}
		}
	case reflect.Map:
        // ...
	case reflect.Slice:
		// ....
	case reflect.Ptr:
		fallthrough
	case reflect.Interface:
        // ...
	default:
		if dst.CanSet() && (!isEmptyValue(src) || overwriteWithEmptySrc) && (overwrite || isEmptyValue(dst)) {
			// 设置值
			dst.Set(src)
		}
	}

	return
}

```


### Map 


## 应用
在 kubernetes 中应用



```go
// https://github.com/kubernetes/kubernetes/blob/88dfa51b6003c90e8f0a0508939a1d79950a40df/staging/src/k8s.io/client-go/tools/clientcmd/loader.go

func (rules *ClientConfigLoadingRules) Load() (*clientcmdapi.Config, error) {
	if err := rules.Migrate(); err != nil {
		return nil, err
	}

	errlist := []error{}
	missingList := []string{}

	kubeConfigFiles := []string{}

    // 。。。。。

	// first merge all of our maps
	mapConfig := clientcmdapi.NewConfig()

	for _, kubeconfig := range kubeconfigs {
		mergo.Merge(mapConfig, kubeconfig, mergo.WithOverride)
	}

	// merge all of the struct values in the reverse order so that priority is given correctly
	// errors are not added to the list the second time
	nonMapConfig := clientcmdapi.NewConfig()
	for i := len(kubeconfigs) - 1; i >= 0; i-- {
		kubeconfig := kubeconfigs[i]
		mergo.Merge(nonMapConfig, kubeconfig, mergo.WithOverride)
	}

	// since values are overwritten, but maps values are not, we can merge the non-map config on top of the map config and
	// get the values we expect.
	config := clientcmdapi.NewConfig()
	mergo.Merge(config, mapConfig, mergo.WithOverride)
	mergo.Merge(config, nonMapConfig, mergo.WithOverride)

	if rules.ResolvePaths() {
		if err := ResolveLocalPaths(config); err != nil {
			errlist = append(errlist, err)
		}
	}
	return config, utilerrors.NewAggregate(errlist)
}
```

## 参考

- [Go 每日一库之 mergo](https://darjun.github.io/2020/03/11/godailylib/mergo/)