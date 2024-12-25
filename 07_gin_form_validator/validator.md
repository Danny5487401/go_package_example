<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/go-playground/validator](#githubcomgo-playgroundvalidator)
  - [默认验证器](#%E9%BB%98%E8%AE%A4%E9%AA%8C%E8%AF%81%E5%99%A8)
  - [实现原理](#%E5%AE%9E%E7%8E%B0%E5%8E%9F%E7%90%86)
  - [初始化](#%E5%88%9D%E5%A7%8B%E5%8C%96)
  - [验证过程](#%E9%AA%8C%E8%AF%81%E8%BF%87%E7%A8%8B)
  - [第三方应用--gin](#%E7%AC%AC%E4%B8%89%E6%96%B9%E5%BA%94%E7%94%A8--gin)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/go-playground/validator

开发接口的时候需要对客户端传提交的参数进行参数校验，如果提交的参数只有一个两个，这样我们可以简单写个if判断，那么要是有很多的参数校验，那么满屏都是参数校验的if判断，效率不仅低还不美观.

## 默认验证器
```go
var bakedInValidators = map[string]Func{
		"required":                      hasValue,
		"required_if":                   requiredIf,
		"required_unless":               requiredUnless,
		"skip_unless":                   skipUnless,
		"required_with":                 requiredWith,
		"required_with_all":             requiredWithAll,
		"required_without":              requiredWithout,
		"required_without_all":          requiredWithoutAll,
		"excluded_if":                   excludedIf,
		"excluded_unless":               excludedUnless,
		"excluded_with":                 excludedWith,
		"excluded_with_all":             excludedWithAll,
		"excluded_without":              excludedWithout,
		"excluded_without_all":          excludedWithoutAll,
		"isdefault":                     isDefault,
		"len":                           hasLengthOf,
		"min":                           hasMinOf,
		"max":                           hasMaxOf,
		"eq":                            isEq,
		"eq_ignore_case":                isEqIgnoreCase,
		"ne":                            isNe,
		"ne_ignore_case":                isNeIgnoreCase,
		"lt":                            isLt,
		"lte":                           isLte,
		"gt":                            isGt,
		"gte":                           isGte,
		"eqfield":                       isEqField,
		"eqcsfield":                     isEqCrossStructField,
		"necsfield":                     isNeCrossStructField,
		"gtcsfield":                     isGtCrossStructField,
		"gtecsfield":                    isGteCrossStructField,
		"ltcsfield":                     isLtCrossStructField,
		"ltecsfield":                    isLteCrossStructField,
		"nefield":                       isNeField,
		"gtefield":                      isGteField,
		"gtfield":                       isGtField,
		"ltefield":                      isLteField,
		"ltfield":                       isLtField,
		"fieldcontains":                 fieldContains,
		"fieldexcludes":                 fieldExcludes,
		"alpha":                         isAlpha,
		"alphanum":                      isAlphanum,
		"alphaunicode":                  isAlphaUnicode,
		"alphanumunicode":               isAlphanumUnicode,
		"boolean":                       isBoolean,
		"numeric":                       isNumeric,
		"number":                        isNumber,
		"hexadecimal":                   isHexadecimal,
		"hexcolor":                      isHEXColor,
		"rgb":                           isRGB,
		"rgba":                          isRGBA,
		"hsl":                           isHSL,
		"hsla":                          isHSLA,
		"e164":                          isE164,
		"email":                         isEmail,
		"url":                           isURL,
		"http_url":                      isHttpURL,
		"uri":                           isURI,
		"urn_rfc2141":                   isUrnRFC2141, // RFC 2141
		"file":                          isFile,
		"filepath":                      isFilePath,
		"base64":                        isBase64,
		"base64url":                     isBase64URL,
        // ...
	}
```



## 实现原理
validator 应用了 Golang 的 Struct Tag 和 Reflect机制，基本思想是：在 Struct Tag 中为不同的字段定义各自类型的约束，然后通过 Reflect 获取这些约束的类型信息并在校验器中进行数据校验


## 初始化

结构体
```go
type Validate struct {
	tagName          string
	pool             *sync.Pool
	hasCustomFuncs   bool
	hasTagNameFunc   bool
	tagNameFunc      TagNameFunc
	structLevelFuncs map[reflect.Type]StructLevelFuncCtx
	customFuncs      map[reflect.Type]CustomTypeFunc
	aliases          map[string]string
	validations      map[string]internalValidationFuncWrapper
	transTagFunc     map[ut.Translator]map[string]TranslationFunc // map[<locale>]map[<tag>]TranslationFunc
	rules            map[reflect.Type]map[string]string
	tagCache         *tagCache
	structCache      *structCache
}
```

```go
func New() *Validate {

	tc := new(tagCache)
	tc.m.Store(make(map[string]*cTag))

	sc := new(structCache)
	sc.m.Store(make(map[reflect.Type]*cStruct))

	v := &Validate{
		tagName:     defaultTagName, // 默认tag使用 validate 的才会进行验证
		aliases:     make(map[string]string, len(bakedInAliases)),
		validations: make(map[string]internalValidationFuncWrapper, len(bakedInValidators)),
		tagCache:    tc,
		structCache: sc,
	}

	// must copy alias validators for separate validations to be used in each validator instance
	for k, val := range bakedInAliases {
		v.RegisterAlias(k, val)
	}

	// 注册默认的验证器，后续会会分发到Ctag上面
	for k, val := range bakedInValidators {

		switch k {
		// these require that even if the value is nil that the validation should run, omitempty still overrides this behaviour
		case requiredIfTag, requiredUnlessTag, requiredWithTag, requiredWithAllTag, requiredWithoutTag, requiredWithoutAllTag,
			excludedIfTag, excludedUnlessTag, excludedWithTag, excludedWithAllTag, excludedWithoutTag, excludedWithoutAllTag,
			skipUnlessTag:
			_ = v.registerValidation(k, wrapFunc(val), true, true)
		default:
			// no need to error check here, baked in will always be valid
			_ = v.registerValidation(k, wrapFunc(val), true, false)
		}
	}

	v.pool = &sync.Pool{
		New: func() interface{} {
			return &validate{
				v:        v,
				ns:       make([]byte, 0, 64),
				actualNs: make([]byte, 0, 64),
				misc:     make([]byte, 32),
			}
		},
	}

	return v
}

```
## 验证过程

这里拿结构体作为案例.

获取结构体元数据，也就是创建cStruct的过程


```go
// github.com/go-playground/validator/v10@v10.14.0/validator_instance.go
func (v *Validate) StructCtx(ctx context.Context, s interface{}) (err error) {

	val := reflect.ValueOf(s)
	top := val

	if val.Kind() == reflect.Ptr && !val.IsNil() {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct || val.Type().ConvertibleTo(timeType) {
		return &InvalidValidationError{Type: reflect.TypeOf(s)}
	}

	// good to validate
	vd := v.pool.Get().(*validate)
	vd.top = top
	vd.isPartial = false
	// vd.hasExcludes = false // only need to reset in StructPartial and StructExcept

	vd.validateStruct(ctx, top, val, val.Type(), vd.ns[0:0], vd.actualNs[0:0], nil)

	if len(vd.errs) > 0 {
		err = vd.errs
		vd.errs = nil
	}

	v.pool.Put(vd)

	return
}

```


```go
// parent and current will be the same the first run of validateStruct
func (v *validate) validateStruct(ctx context.Context, parent reflect.Value, current reflect.Value, typ reflect.Type, ns []byte, structNs []byte, ct *cTag) {
    // 获取 structCache 对象，该对象中保存着结构体的所有字段，每个字段里面包含tags对象链，tags对象中包含验证方法
	cs, ok := v.v.structCache.Get(typ)
	if !ok {
		cs = v.v.extractStructCache(current, typ.Name())
	}

	if len(ns) == 0 && len(cs.name) != 0 {

		ns = append(ns, cs.name...)
		ns = append(ns, '.')

		structNs = append(structNs, cs.name...)
		structNs = append(structNs, '.')
	}

	// ct is nil on top level struct, and structs as fields that have no tag info
	// so if nil or if not nil and the structonly tag isn't present
	if ct == nil || ct.typeof != typeStructOnly {

		var f *cField

		for i := 0; i < len(cs.fields); i++ { // 遍历每个字段

			f = cs.fields[i]

			if v.isPartial {

				if v.ffn != nil {
					// used with StructFiltered
					if v.ffn(append(structNs, f.name...)) {
						continue
					}

				} else {
					// used with StructPartial & StructExcept
					_, ok = v.includeExclude[string(append(structNs, f.name...))]

					if (ok && v.hasExcludes) || (!ok && !v.hasExcludes) {
						continue
					}
				}
			}
			// 针对每个字段的tags进行验证
			v.traverseField(ctx, current, current.Field(f.idx), ns, structNs, f, f.cTags)
		}
	}

	// check if any struct level validations, after all field validations already checked.
	// first iteration will have no info about nostructlevel tag, and is checked prior to
	// calling the next iteration of validateStruct called from traverseField.
	if cs.fn != nil {

		v.slflParent = parent
		v.slCurrent = current
		v.ns = ns
		v.actualNs = structNs

		cs.fn(ctx, v)
	}
}

```

## 第三方应用--gin

默认炎症期
```go
// github.com/gin-gonic/gin@v1.9.1/binding/binding.go
var Validator StructValidator = &defaultValidator{}
```

```go
// github.com/gin-gonic/gin@v1.9.1/binding/default_validator.go
type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}


// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *defaultValidator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}
```

这里拿 struct 作为案例
```go
func (v *defaultValidator) validateStruct(obj any) error {
	v.lazyinit()
	// 实际调用 validator.Validate 的 Struct 方法
	return v.validate.Struct(obj)
}

```


## 参考

