<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/go-playground/validator](#githubcomgo-playgroundvalidator)
  - [默认验证器](#%E9%BB%98%E8%AE%A4%E9%AA%8C%E8%AF%81%E5%99%A8)
  - [实现原理](#%E5%AE%9E%E7%8E%B0%E5%8E%9F%E7%90%86)
    - [初始化](#%E5%88%9D%E5%A7%8B%E5%8C%96)
    - [验证过程](#%E9%AA%8C%E8%AF%81%E8%BF%87%E7%A8%8B)
    - [错误处理](#%E9%94%99%E8%AF%AF%E5%A4%84%E7%90%86)
  - [第三方应用--gin](#%E7%AC%AC%E4%B8%89%E6%96%B9%E5%BA%94%E7%94%A8--gin)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/go-playground/validator

开发接口的时候需要对客户端传提交的参数进行参数校验，如果提交的参数只有一个两个，这样我们可以简单写个if判断，那么要是有很多的参数校验，那么满屏都是参数校验的if判断，效率不仅低还不美观.

## 默认验证器
校验器主要
- Fields:对于结构体各个属性的校验，这里可以针对一个 field 与另一个 field 相互比较
- network: 网络相关的格式校验，可以用来校验 IP 格式，TCP, UDP, URL 等
- string: 字符串相关的校验，比如校验是否是数字，大小写，前后缀等
- format: 符合特定格式，如我们上面提到的 email，信用卡号，颜色，html，base64，json，经纬度，md5 等
- Comparisons: 比较大小
- other:杂项，各种通用能力


```go
type Func func(fl FieldLevel) bool
var bakedInValidators = map[string]Func{
		"required":                      hasValue,
		"required_if":                   requiredIf,
		"required_unless":               requiredUnless,
		"skip_unless":                   skipUnless,
		"required_with":                 requiredWith,
		"required_with_all":             requiredWithAll,
		"required_without":              requiredWithout,
		"required_without_all":          requiredWithoutAll,
        // ...
	}
```
这里拿 required 对应的 hasValue 作为案例
```go
func hasValue(fl FieldLevel) bool {
	field := fl.Field()
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		if fl.(*validate).fldIsPointer && field.Interface() != nil {
			return true
		}
		// 校验是否为空的经典解法，日常开发也用得上
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}
```
注册时会是wrap一层
```go
// FuncCtx accepts a context.Context and FieldLevel interface for all
// validation needs. The return value should be true when validation succeeds.
type FuncCtx func(ctx context.Context, fl FieldLevel) bool

// wrapFunc wraps normal Func makes it compatible with FuncCtx
func wrapFunc(fn Func) FuncCtx {
	if fn == nil {
		return nil // be sure not to wrap a bad function.
	}
	return func(ctx context.Context, fl FieldLevel) bool {
		return fn(fl)
	}
}

```



## 实现原理
validator 应用了 Golang 的 Struct Tag 和 Reflect机制，基本思想是：在 Struct Tag 中为不同的字段定义各自类型的约束，然后通过 Reflect 获取这些约束的类型信息并在校验器中进行数据校验



```go
// tag 类型
type tagType uint8

const (
	typeDefault tagType = iota
	typeOmitEmpty
	typeIsDefault
	typeNoStructLevel
	typeStructOnly
	typeDive // dive 的语义在于告诉 validator 不要停留在我这一级，而是继续往下校验
	typeOr // 符号 | 处理
	typeKeys
	typeEndKeys
)

```

### 初始化

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

### 验证过程

这里拿结构体作为案例.


重要结构体
```go
type cStruct struct {
	name   string // 结构体名字
	fields []*cField // 内部字段
	fn     StructLevelFuncCtx
}

type cField struct {
	idx        int
	name       string // 字段名字
	altName    string
	namesEqual bool
	cTags      *cTag
}

```




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
    // 获取 cStruct 对象，该对象中保存着结构体的所有字段，每个字段里面包含tags对象链，tags对象中包含验证方法
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
获取结构体元数据，也就是创建cStruct的过程
```go
// 不存在开始构建
func (v *Validate) extractStructCache(current reflect.Value, sName string) *cStruct {
	v.structCache.lock.Lock()
	defer v.structCache.lock.Unlock() // leave as defer! because if inner panics, it will never get unlocked otherwise!

	typ := current.Type()

	// could have been multiple trying to access, but once first is done this ensures struct
	// isn't parsed again.
	cs, ok := v.structCache.Get(typ)
	if ok {
		return cs
	}

	cs = &cStruct{name: sName, fields: make([]*cField, 0), fn: v.structLevelFuncs[typ]}

	numFields := current.NumField()
	rules := v.rules[typ]

	var ctag *cTag
	var fld reflect.StructField
	var tag string
	var customName string

	for i := 0; i < numFields; i++ {

		fld = typ.Field(i)

		if !fld.Anonymous && len(fld.PkgPath) > 0 {
			continue
		}

		if rtag, ok := rules[fld.Name]; ok {
			tag = rtag
		} else {
			tag = fld.Tag.Get(v.tagName)
		}
		// 如果是  "-" ,忽略
		if tag == skipValidationTag {
			continue
		}

		customName = fld.Name

		if v.hasTagNameFunc {
			name := v.tagNameFunc(fld)
			if len(name) > 0 {
				customName = name
			}
		}

		// NOTE: cannot use shared tag cache, because tags may be equal, but things like alias may be different
		// and so only struct level caching can be used instead of combined with Field tag caching

		if len(tag) > 0 {
			ctag, _ = v.parseFieldTagsRecursive(tag, fld.Name, "", false)
		} else {
			// even if field doesn't have validations need cTag for traversing to potential inner/nested
			// elements of the field.
			ctag = new(cTag)
		}

		cs.fields = append(cs.fields, &cField{
			idx:        i,
			name:       fld.Name,
			altName:    customName,
			cTags:      ctag,
			namesEqual: fld.Name == customName,
		})
	}
	v.structCache.Set(typ, cs)
	return cs
}

```

切分 tag 信息
```go
func (v *Validate) parseFieldTagsRecursive(tag string, fieldName string, alias string, hasAlias bool) (firstCtag *cTag, current *cTag) {
	var t string
	noAlias := len(alias) == 0
	// 通过 "," 进行切换
	tags := strings.Split(tag, tagSeparator)

	for i := 0; i < len(tags); i++ {
		t = tags[i]
		if noAlias {
			alias = t
		}

		// check map for alias and process new tags, otherwise process as usual
		if tagsVal, found := v.aliases[t]; found {
			if i == 0 {
				firstCtag, current = v.parseFieldTagsRecursive(tagsVal, fieldName, t, true)
			} else {
				next, curr := v.parseFieldTagsRecursive(tagsVal, fieldName, t, true)
				current.next, current = next, curr

			}
			continue
		}

		var prevTag tagType

		if i == 0 {
			current = &cTag{aliasTag: alias, hasAlias: hasAlias, hasTag: true, typeof: typeDefault}
			firstCtag = current
		} else {
			prevTag = current.typeof
			current.next = &cTag{aliasTag: alias, hasAlias: hasAlias, hasTag: true}
			current = current.next
		}

		switch t {
		case diveTag:
			current.typeof = typeDive
			continue

		case keysTag:
			current.typeof = typeKeys

			if i == 0 || prevTag != typeDive {
				panic(fmt.Sprintf("'%s' tag must be immediately preceded by the '%s' tag", keysTag, diveTag))
			}

			current.typeof = typeKeys

			// need to pass along only keys tag
			// need to increment i to skip over the keys tags
			b := make([]byte, 0, 64)

			i++

			for ; i < len(tags); i++ {

				b = append(b, tags[i]...)
				b = append(b, ',')

				if tags[i] == endKeysTag {
					break
				}
			}

			current.keys, _ = v.parseFieldTagsRecursive(string(b[:len(b)-1]), fieldName, "", false)
			continue

		case endKeysTag:
			current.typeof = typeEndKeys

			// if there are more in tags then there was no keysTag defined
			// and an error should be thrown
			if i != len(tags)-1 {
				panic(keysTagNotDefined)
			}
			return

		case omitempty:
			current.typeof = typeOmitEmpty
			continue

		case structOnlyTag:
			current.typeof = typeStructOnly
			continue

		case noStructLevelTag:
			current.typeof = typeNoStructLevel
			continue

		default:
			if t == isdefault {
				current.typeof = typeIsDefault
			}
			// if a pipe character is needed within the param you must use the utf8Pipe representation "0x7C"
			orVals := strings.Split(t, orSeparator)

			for j := 0; j < len(orVals); j++ {
				vals := strings.SplitN(orVals[j], tagKeySeparator, 2)
				if noAlias {
					alias = vals[0]
					current.aliasTag = alias
				} else {
					current.actualAliasTag = t
				}

				if j > 0 {
					current.next = &cTag{aliasTag: alias, actualAliasTag: current.actualAliasTag, hasAlias: hasAlias, hasTag: true}
					current = current.next
				}
				current.hasParam = len(vals) > 1

				current.tag = vals[0]
				if len(current.tag) == 0 {
					panic(strings.TrimSpace(fmt.Sprintf(invalidValidation, fieldName)))
				}

				if wrapper, ok := v.validations[current.tag]; ok {
					current.fn = wrapper.fn
					current.runValidationWhenNil = wrapper.runValidatinOnNil
				} else {
					panic(strings.TrimSpace(fmt.Sprintf(undefinedValidation, current.tag, fieldName)))
				}

				if len(orVals) > 1 {
					current.typeof = typeOr
				}

				if len(vals) > 1 {
					current.param = strings.Replace(strings.Replace(vals[1], utf8HexComma, ",", -1), utf8Pipe, "|", -1)
				}
			}
			current.isBlockEnd = true
		}
	}
	return
}

```


字段实际校验

```go
func (v *validate) traverseField(ctx context.Context, parent reflect.Value, current reflect.Value, ns []byte, structNs []byte, cf *cField, ct *cTag) {
	var typ reflect.Type
	var kind reflect.Kind

	current, kind, v.fldIsPointer = v.extractTypeInternal(current, false)

	switch kind {
	case reflect.Ptr, reflect.Interface, reflect.Invalid:

		if ct == nil {
			return
		}

		if ct.typeof == typeOmitEmpty || ct.typeof == typeIsDefault {
			return
		}

		if ct.hasTag {
			if kind == reflect.Invalid {
				v.str1 = string(append(ns, cf.altName...))
				if v.v.hasTagNameFunc {
					v.str2 = string(append(structNs, cf.name...))
				} else {
					v.str2 = v.str1
				}
				v.errs = append(v.errs,
					&fieldError{
						v:              v.v,
						tag:            ct.aliasTag,
						actualTag:      ct.tag,
						ns:             v.str1,
						structNs:       v.str2,
						fieldLen:       uint8(len(cf.altName)),
						structfieldLen: uint8(len(cf.name)),
						param:          ct.param,
						kind:           kind,
					},
				)
				return
			}

			v.str1 = string(append(ns, cf.altName...))
			if v.v.hasTagNameFunc {
				v.str2 = string(append(structNs, cf.name...))
			} else {
				v.str2 = v.str1
			}
			if !ct.runValidationWhenNil {
				v.errs = append(v.errs,
					&fieldError{
						v:              v.v,
						tag:            ct.aliasTag,
						actualTag:      ct.tag,
						ns:             v.str1,
						structNs:       v.str2,
						fieldLen:       uint8(len(cf.altName)),
						structfieldLen: uint8(len(cf.name)),
						value:          current.Interface(),
						param:          ct.param,
						kind:           kind,
						typ:            current.Type(),
					},
				)
				return
			}
		}

	case reflect.Struct:

		typ = current.Type()

		if !typ.ConvertibleTo(timeType) {

			if ct != nil {

				if ct.typeof == typeStructOnly {
					goto CONTINUE
				} else if ct.typeof == typeIsDefault {
					// set Field Level fields
					v.slflParent = parent
					v.flField = current
					v.cf = cf
					v.ct = ct

					if !ct.fn(ctx, v) {
						v.str1 = string(append(ns, cf.altName...))

						if v.v.hasTagNameFunc {
							v.str2 = string(append(structNs, cf.name...))
						} else {
							v.str2 = v.str1
						}

						v.errs = append(v.errs,
							&fieldError{
								v:              v.v,
								tag:            ct.aliasTag,
								actualTag:      ct.tag,
								ns:             v.str1,
								structNs:       v.str2,
								fieldLen:       uint8(len(cf.altName)),
								structfieldLen: uint8(len(cf.name)),
								value:          current.Interface(),
								param:          ct.param,
								kind:           kind,
								typ:            typ,
							},
						)
						return
					}
				}

				ct = ct.next
			}

			if ct != nil && ct.typeof == typeNoStructLevel {
				return
			}

		CONTINUE:
			// if len == 0 then validating using 'Var' or 'VarWithValue'
			// Var - doesn't make much sense to do it that way, should call 'Struct', but no harm...
			// VarWithField - this allows for validating against each field within the struct against a specific value
			//                pretty handy in certain situations
			if len(cf.name) > 0 {
				ns = append(append(ns, cf.altName...), '.')
				structNs = append(append(structNs, cf.name...), '.')
			}

			v.validateStruct(ctx, parent, current, typ, ns, structNs, ct)
			return
		}
	}

	if ct == nil || !ct.hasTag {
		return
	}

	typ = current.Type()

OUTER:
	for {
		if ct == nil {
			return
		}

		switch ct.typeof {

		case typeOmitEmpty:

			// set Field Level fields
			v.slflParent = parent
			v.flField = current
			v.cf = cf
			v.ct = ct

			if !hasValue(v) {
				return
			}

			ct = ct.next
			continue

		case typeEndKeys:
			return
		// ...	
		// 这里关心 default 即可
		default:

			// set Field Level fields
			v.slflParent = parent
			v.flField = current
			v.cf = cf
			v.ct = ct

			if !ct.fn(ctx, v) { // 校验不通过
				v.str1 = string(append(ns, cf.altName...))

				if v.v.hasTagNameFunc {
					v.str2 = string(append(structNs, cf.name...))
				} else {
					v.str2 = v.str1
				}

				v.errs = append(v.errs,
					&fieldError{
						v:              v.v,
						tag:            ct.aliasTag,
						actualTag:      ct.tag,
						ns:             v.str1,
						structNs:       v.str2,
						fieldLen:       uint8(len(cf.altName)),
						structfieldLen: uint8(len(cf.name)),
						value:          current.Interface(),
						param:          ct.param,
						kind:           kind,
						typ:            typ,
					},
				)

				return
			}
			ct = ct.next
		}
	}

}

```

### 错误处理
validator 返回的类型底层是 validator.ValidationErrors
```go
type ValidationErrors []FieldError

func (ve ValidationErrors) Error() string {

	buff := bytes.NewBufferString("")

	for i := 0; i < len(ve); i++ {

		buff.WriteString(ve[i].Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}
```

FieldError 实现的结构体
```go
type fieldError struct {
	v              *Validate
	tag            string
	actualTag      string
	ns             string
	structNs       string
	fieldLen       uint8
	structfieldLen uint8
	value          interface{}
	param          string
	kind           reflect.Kind
	typ            reflect.Type
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


- [解析 Golang 经典校验库 validator 设计和原理](https://juejin.cn/post/7136135907249225758)