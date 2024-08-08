<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Go 原生 encoding/json](#go-%E5%8E%9F%E7%94%9F-encodingjson)
  - [接口](#%E6%8E%A5%E5%8F%A3)
  - [标签](#%E6%A0%87%E7%AD%BE)
  - [序列化](#%E5%BA%8F%E5%88%97%E5%8C%96)
    - [源码](#%E6%BA%90%E7%A0%81)
      - [问题](#%E9%97%AE%E9%A2%98)
  - [反序列化](#%E5%8F%8D%E5%BA%8F%E5%88%97%E5%8C%96)
    - [反序列化源码](#%E5%8F%8D%E5%BA%8F%E5%88%97%E5%8C%96%E6%BA%90%E7%A0%81)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Go 原生 encoding/json
使用 json.Unmarshal 和 json.Marshal 函数，可以轻松将 JSON 格式的二进制数据反序列化到指定的 Go 结构体中，以及将 Go 结构体序列化为二进制流。
而对于未知结构或不确定结构的数据，则支持将二进制反序列化到 map[string]interface{} 类型中，使用 KV 的模式进行数据的存取



序列化和反序列化的开销完全不同，JSON 反序列化的开销是序列化开销的好几倍。


## 接口

```go
// go1.22.2/src/encoding/json/encode.go
type Marshaler interface {
	MarshalJSON() ([]byte, error)
}

```


```go
// go1.22.2/src/encoding/json/decode.go
type Unmarshaler interface {
	UnmarshalJSON([]byte) error
}

```

在 JSON 序列化和反序列化的过程中，它会使用反射判断结构体类型是否实现了上述接口，如果实现了上述接口就会优先使用对应的方法进行编码和解码操作，
除了这两个方法之外，Go 语言其实还提供了另外两个用于控制编解码结果的方法，即 encoding.TextMarshaler 和 encoding.TextUnmarshaler

```go
// go1.22.2/src/encoding/encoding.go
type TextMarshaler interface {
	MarshalText() (text []byte, err error)
}

type TextUnmarshaler interface {
	UnmarshalText(text []byte) error
}


```

## 标签

- string 表示当前的整数或者浮点数是由 JSON 中的字符串表示的. 只能 string, floating point, integer, or boolean types
- omitempty 会在字段为零值时，直接在生成的 JSON 中忽略对应的键值对.
- inline 标记会告诉编码/解码器将嵌套结构体字段展开到父结构体中。

## 序列化


Marshal 递归地遍历参数v 。如果遇到一个实现了Marshaler接口的值并且它是非nil的，就会调用这个值的MarshalJSON方法；没有的话就检查encoding.TextMarshaler接口的MarshalText方法。

如果上述两个接口都没有，则按照以下规则：

翻译如下：
- Boolean类型转化为 JSON 的 boolean类型.
- 浮点，整形，数字类型都转化为JSON 的 number类型.
- String 则强制使用UTF-8编码，将无效字节转化为Unicode形式的rune。并且默认使用 HTMLEscape （将 < > & U+2028 U+2029 替换为 \u003c \u003e \u0026 \u2028 \u2029），这个行为可以通过 Encoder 来自定义。
- Array, slice 相应地转化为 JSON array，但是 []byte会编为 base64字符串。空切片则转换为 JSON null 。
- Struct 转为为 JSON object 。每个公开成员（大写开头）都会作为object的一个成员，使用字段名作为键，除非：
  - 可以通过字段tag中 json 来指定名称，作为 object的键；
  - 名称后面，可以用逗号分割来附带一些额外的配置；名称可以留空，以保留默认的键，同时附带配置。
  - 配置omitempty时，则当字段为空值（零值）的时候忽略这个字段。
  - 如果名称指定为-，则总是忽略这个字段。（注意如果想让键名就是-，则要写-,）
  - 除了omitempty之外还有个string选项，它会把相应的值以string的形式转化，这只对数字和布尔类型生效
- 键名必须符合如下规则：只包含Unicode中的字母、数字 和 除了引号、反斜杠、逗号之外的ASCII标点符号
- 对于匿名字段
  - 如果没有给它指定tag，那么它其中的字段会平铺在父对象中
  - 如果指定了tag，则视为一个子对象
  - 如果是interface类型，一律视为子对象
- 键有冲突时，优先使用有tag指定的字段
- 指针会被转换为其所指向的值。空指针转化为null
- channel, complex, 函数 不能被序列化，会返回错误
- 循环引用会返回错误
### 源码
```go
func Marshal(v any) ([]byte, error) {
	// 创建并复用对象
	e := newEncodeState()
	defer encodeStatePool.Put(e)

	err := e.marshal(v, encOpts{escapeHTML: true})
	if err != nil {
		return nil, err
	}
	buf := append([]byte(nil), e.Bytes()...)

	return buf, nil
}
```

初始化获取全局 encodeState
```go
var encodeStatePool sync.Pool

func newEncodeState() *encodeState {
	if v := encodeStatePool.Get(); v != nil {
		e := v.(*encodeState)
		e.Reset()
		if len(e.ptrSeen) > 0 {
			panic("ptrEncoder.encode should have emptied ptrSeen via defers")
		}
		e.ptrLevel = 0
		return e
	}
	return &encodeState{ptrSeen: make(map[any]struct{})}
}


type encodeState struct {
	bytes.Buffer // accumulated output

	// 下面两个字段用来防止循环引用.
	ptrLevel uint
	ptrSeen  map[any]struct{}
}
```

根据类型来选择一个处理的函数(encoderFunc)，然后调用这个函数去处理
```go
func (e *encodeState) reflectValue(v reflect.Value, opts encOpts) {
	valueEncoder(v)(e, v, opts)
}

func newTypeEncoder(t reflect.Type, allowAddr bool) encoderFunc {
	// 如果当前不是指针类型、可以寻址的并且它的指针实现了Marshaler,那么就可以得到一个条件编码，然后再根据这个类型是否可以寻址去得到addrMarshalerEncoder编码方法或使用内置的编码方法
	if t.Kind() != reflect.Pointer && allowAddr && reflect.PointerTo(t).Implements(marshalerType) {
		return newCondAddrEncoder(addrMarshalerEncoder, newTypeEncoder(t, false))
	}
	
	// 当前的数据类型直接实现了Marshaler,返回marshalerEncoder编码方法
	if t.Implements(marshalerType) {
		return marshalerEncoder
	}
	
	// 如果当前不是指针类型、可以寻址的并且它的指针实现了TextMarshaler,那么就可以得到一个条件编码，然后再根据这个类型是否可以寻址去得到addrTextMarshalerEncoder编码方法或使用内置的编码方法
	if t.Kind() != reflect.Pointer && allowAddr && reflect.PointerTo(t).Implements(textMarshalerType) {
		return newCondAddrEncoder(addrTextMarshalerEncoder, newTypeEncoder(t, false))
	}
	
	// 当前数据类型直接实现了TextMarshaler,返回textMarshalerEncoder编码方法
	if t.Implements(textMarshalerType) {
		return textMarshalerEncoder
	}

	switch t.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.Float32:
		return float32Encoder
	case reflect.Float64:
		return float64Encoder
	case reflect.String:
		return stringEncoder
	case reflect.Interface:
		return interfaceEncoder
	case reflect.Struct:
		return newStructEncoder(t)
	case reflect.Map:
		return newMapEncoder(t)
	case reflect.Slice:
		return newSliceEncoder(t)
	case reflect.Array:
		return newArrayEncoder(t)
	case reflect.Pointer:
		return newPtrEncoder(t)
	default:
		return unsupportedTypeEncoder
	}
}
```


案例一：mapEncoder 
```go
type mapEncoder struct {
   elemEnc encoderFunc
}

func (me mapEncoder) encode(e *encodeState, v reflect.Value, opts encOpts) {
   if v.IsNil() {//为nil时，返回null
      e.WriteString("null")
      return
   }
   e.WriteByte('{')

   // Extract and sort the keys.
   keys := v.MapKeys()//获取map中的所有keys
   sv := make([]reflectWithString, len(keys))
   for i, v := range keys {
      sv[i].v = v
      if err := sv[i].resolve(); err != nil {//处理key，尤其是非string（int/uint）类型的key转string
         e.error(&MarshalerError{v.Type(), err})
      }
   }
   //排序，升序，直接比较字符串
   sort.Slice(sv, func(i, j int) bool { return sv[i].s < sv[j].s })

   for i, kv := range sv {
      if i > 0 {
         e.WriteByte(',')
      }
      e.string(kv.s, opts.escapeHTML)
      e.WriteByte(':')
      me.elemEnc(e, v.MapIndex(kv.v), opts)
   }
   e.WriteByte('}')
}

func newMapEncoder(t reflect.Type) encoderFunc {
    // 检查 键 的类型
    switch t.Key().Kind() {
    case reflect.String,
      reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
      reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
    default:
      if !t.Key().Implements(textMarshalerType) {
         return unsupportedTypeEncoder
      }
    }
    me := mapEncoder{typeEncoder(t.Elem())}
    return me.encode
}

```


案例二：structEncoder

```go
type structEncoder struct {
	fields structFields
}

type structFields struct {
	list         []field
	byExactName  map[string]*field
	byFoldedName map[string]*field
}

func (se structEncoder) encode(e *encodeState, v reflect.Value, opts encOpts) {
	next := byte('{')
FieldLoop:
	for i := range se.fields.list {
		f := &se.fields.list[i]

		// Find the nested struct field by following f.index.
		fv := v
		for _, i := range f.index {
			if fv.Kind() == reflect.Pointer {
				if fv.IsNil() {
					continue FieldLoop
				}
				fv = fv.Elem()
			}
			fv = fv.Field(i)
		}

		if f.omitEmpty && isEmptyValue(fv) {
			continue
		}
		e.WriteByte(next)
		next = ','
		if opts.escapeHTML {
			e.WriteString(f.nameEscHTML)
		} else {
			e.WriteString(f.nameNonEsc)
		}
		opts.quoted = f.quoted
		f.encoder(e, fv, opts)
	}
	if next == '{' {
		e.WriteString("{}")
	} else {
		e.WriteByte('}')
	}
}

func newStructEncoder(t reflect.Type) encoderFunc {
	se := structEncoder{fields: cachedTypeFields(t)}
	return se.encode
}
```

缓存子段
```go
func cachedTypeFields(t reflect.Type) structFields {
	if f, ok := fieldCache.Load(t); ok {
		return f.(structFields)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.(structFields)
}


func typeFields(t reflect.Type) structFields {
	// Anonymous fields to explore at the current level and the next.
	current := []field{}
	next := []field{{typ: t}}

	// Count of queued names for current level and the next.
	var count, nextCount map[reflect.Type]int

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	// Fields found.
	var fields []field

	// Buffer to run appendHTMLEscape on field names.
	var nameEscBuf []byte

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true

			// 遍历结构体的每个字段
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				if sf.Anonymous {
                    //表示是否是嵌套类型
					t := sf.Type
					if t.Kind() == reflect.Pointer {
						t = t.Elem()
					}
					if !sf.IsExported() && t.Kind() != reflect.Struct {
						// Ignore embedded fields of unexported non-struct types.
						continue
					}
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
				} else if !sf.IsExported() {
					// Ignore unexported non-embedded fields.
					continue
				}
				tag := sf.Tag.Get("json")
				if tag == "-" {
					continue
				}
				name, opts := parseTag(tag)
				if !isValidTag(name) {
					name = ""
				}
				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Pointer {
					// Follow pointer.
					ft = ft.Elem()
				}

				// Only strings, floats, integers, and booleans can be quoted.
				quoted := false
				if opts.Contains("string") {
					switch ft.Kind() {
					case reflect.Bool,
						reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
						reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
						reflect.Float32, reflect.Float64,
						reflect.String:
						quoted = true
					}
				}

				// Record found field and index sequence.
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}
					field := field{
						name:      name,
						tag:       tagged,
						index:     index,
						typ:       ft,
						omitEmpty: opts.Contains("omitempty"),
						quoted:    quoted,
					}
					field.nameBytes = []byte(field.name)

					// Build nameEscHTML and nameNonEsc ahead of time.
					nameEscBuf = appendHTMLEscape(nameEscBuf[:0], field.nameBytes)
					field.nameEscHTML = `"` + string(nameEscBuf) + `":`
					field.nameNonEsc = `"` + field.name + `":`

					fields = append(fields, field)
					if count[f.typ] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 and 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}
					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					next = append(next, field{name: ft.Name(), index: index, typ: ft})
				}
			}
		}
	}

	sort.Slice(fields, func(i, j int) bool {
		x := fields
		// sort field by name, breaking ties with depth, then
		// breaking ties with "name came from json tag", then
		// breaking ties with index sequence.
		if x[i].name != x[j].name {
			return x[i].name < x[j].name
		}
		if len(x[i].index) != len(x[j].index) {
			return len(x[i].index) < len(x[j].index)
		}
		if x[i].tag != x[j].tag {
			return x[i].tag
		}
		return byIndex(x).Less(i, j)
	})

	// Delete all fields that are hidden by the Go rules for embedded fields,
	// except that fields with JSON tags are promoted.

	// The fields are sorted in primary order of name, secondary order
	// of field index length. Loop over names; for each name, delete
	// hidden fields by choosing the one dominant field that survives.
	out := fields[:0]
	for advance, i := 0, 0; i < len(fields); i += advance {
		// One iteration per name.
		// Find the sequence of fields with the name of this first field.
		fi := fields[i]
		name := fi.name
		for advance = 1; i+advance < len(fields); advance++ {
			fj := fields[i+advance]
			if fj.name != name {
				break
			}
		}
		if advance == 1 { // Only one field with this name
			out = append(out, fi)
			continue
		}
		dominant, ok := dominantField(fields[i : i+advance])
		if ok {
			out = append(out, dominant)
		}
	}

	fields = out
	sort.Sort(byIndex(fields))

	for i := range fields {
		f := &fields[i]
		f.encoder = typeEncoder(typeByIndex(t, f.index))
	}
	exactNameIndex := make(map[string]*field, len(fields))
	foldedNameIndex := make(map[string]*field, len(fields))
	for i, field := range fields {
		exactNameIndex[field.name] = &fields[i]
		// For historical reasons, first folded match takes precedence.
		if _, ok := foldedNameIndex[string(foldName(field.nameBytes))]; !ok {
			foldedNameIndex[string(foldName(field.nameBytes))] = &fields[i]
		}
	}
	return structFields{fields, exactNameIndex, foldedNameIndex}
}
```




#### 问题

map是无序的，每次取出key/value的顺序都可能不一致，但map转json的顺序是不是也是无序的吗？

结论：map转json是有序的，按照ASCII码升序排列key。

## 反序列化

翻译：
- 传入空指针或者非指针的话，会返回错误
- 它的过程与Marshal相反；它会分配map, slice, 指针 等，按照以下规律：
  - 首先检查 JSON null，如果是则把指针设为nil；如果JSON有数据，则将其填入指针所指向的数据内存中；如果指针为空，则会new一个。
  - 先检查Unmarshaler（包括JSON null），然后如果JSON字段是quoted字符串，则检查TextUnmarshaler。
- 结构体
  - 先查找与key一样的json标签，找到则赋值给该标签对应的变量(如Name)。
  - 没有json标签的，就从上往下依次查找变量名与key一样的变量，如Age。或者变量名忽略大小写后与key一样的变量。如HIgh，Class。第一个匹配的就赋值，后面就算有匹配的也忽略。
    (前提是该变量必需是可导出的，即首字母大写)。
  - 反序列化过程中，先检查JSON的键。如果struct中没有对应的字段，则默认情况下会忽略。可以选择设置 Decoder.DisallowUnknownFields
- interface
  - JSON boolean -> bool
  - JSON numbers -> float64
  - JSON string -> string
  - JSON array -> []interface{}
  - JSON objects -> map[string]interface{}
  - JSON null -> nil
- array转化为切片：先将切片长度设置为0，然后逐个append进去
- array转化为数组：多出的会被抛弃，不足的会被设为零值；
- object转为map：如果map是nil则会创建一个，如果有旧的map则用旧的map；键的类型必须是string, integer或者实现了json.Unmarshaler或encoding.TextUnmarshaler；
- 如果一个JSON的值与目标类型不匹配，或者number超出了范围，则会跳过这个值，并继续尽可能地完成剩下的部分。如果后续没有更严重的错误，则会返回UnmarshalTypeError来描述遇到的第一个不匹配的类型。注意，如果出现类型不匹配的情况，那么不保证后续字段都会正常工作
- 当解析字符串的值的时候，无效的 utf-8 或者 utf-16 字符不会被视为一种错误；这些无效字符会被替换为 替换字符 U+FFFD



### 反序列化源码

使用
```go
var s string
err := json.Unmarshal([]byte(`"Hello, world!"`), &s)
// 注意字符串中的双引号不能缺，如果仅仅是 `Hello, world`，则这不是一个合法的 JSON 序列，会返回错误。
```


```go
// encoding/json/decode.go
func Unmarshal(data []byte, v interface{}) error {
	var d decodeState
	err := checkValid(data, &d.scan) // 使用一个状态机来判断是否是合法的json字符串
	if err != nil {
		return err
	}

	d.init(data) //初始化数据
	return d.unmarshal(v) //反序列化
}


func (d *decodeState) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	// 保证v是指针,所以这个对象必须是可写入的
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

    // 。。。
	// We decode rv not rv.Elem because the Unmarshaler interface
	// test must be applied at the top level of the value.
	err := d.value(rv)
    // ...
}
```


在反序列化过程中,会先创建一个scanner类,用于扫描数据、更新解析的状态、分析下一步的步骤等。
解析特殊符号主要集中在两个方法中stateBeginValue和stateEndValue,分为开始符和结束符号,开始符号主要有:{、[、"、-、t(true)、f(false)、n(null)、数字等,结束符则为:、,、}、]等,反序列化正是不断在起始符和结束符中扭转，
```go
func (s *scanner) reset() {
	s.step = stateBeginValue
	s.parseState = s.parseState[0:0]
	s.err = nil
	s.endTop = false
}

// stateBeginValue is the state at the beginning of the input.
func stateBeginValue(s *scanner, c byte) int {
	if isSpace(c) {
		return scanSkipSpace
	}
	switch c {
	case '{':
		s.step = stateBeginStringOrEmpty
		// 说明当前扫描得是对象,当前要要先解析对象的key值
		return s.pushParseState(c, parseObjectKey, scanBeginObject)
	case '[':
		s.step = stateBeginValueOrEmpty
		// 说明当前扫描的是数组,当前要先解析对象的元素值
		return s.pushParseState(c, parseArrayValue, scanBeginArray)
	case '"':
		s.step = stateInString
		// 解析的是字面量信息,key或value
		return scanBeginLiteral
	case '-':
		s.step = stateNeg
		return scanBeginLiteral
	case '0': // beginning of 0.123
		s.step = state0
		return scanBeginLiteral
	case 't': // beginning of true
		s.step = stateT
		return scanBeginLiteral
	case 'f': // beginning of false
		s.step = stateF
		return scanBeginLiteral
	case 'n': // beginning of null
		s.step = stateN
		return scanBeginLiteral
	}
	if '1' <= c && c <= '9' { // beginning of 1234.5
		s.step = state1
		return scanBeginLiteral
	}
	return s.error(c, "looking for beginning of value")
}
```


```go
func (d *decodeState) value(v reflect.Value) error {
    switch d.opcode {
    default:
        panic(phasePanicMsg)
    // 数组 
    case scanBeginArray:
        //...
    // 对象：结构体或map
    case scanBeginObject:
		if v.IsValid() {
            if err := d.object(v); err != nil {
                return err
            }
        } else {
            d.skip()
        }
        d.scanNext()
        // ...
    // 字面量，包括 int、string、float 等
    case scanBeginLiteral:
        // ...
    }
    return nil
}
```


以解析对象为例

检查是map还是struct：

```go
func (d *decodeState) object(v reflect.Value) error {
    // ......
    switch v.Kind() {
    case reflect.Map:
        // ...
    case reflect.Struct:
        fields = cachedTypeFields(t)
    default:
        d.saveError(&UnmarshalTypeError{Value: "object", Type: t, Offset: int64(d.off)})
        d.skip()
        return nil
    }
    // ......
	
}
```
巨大的循环，一个一个地将JSON字符串中的键和值取出来
```go
func (d *decodeState) object(v reflect.Value) error {
    for {
        // 1. 扫描一个键
        d.scanWhile(scanSkipSpace)
        item := d.data[start:d.readIndex()]

        if v.Kind() == reflect.Map {
            // ...
        } else {
            // 2. 找到这个键对应的结构体字段
            var f *field = &fields.list[i]
            d.errorContext.FieldStack = append(d.errorContext.FieldStack, f.name)
        }

        // 3. 把值给取出来
        if err := d.value(subv); err != nil {
            return err
        }
	}	
}
```

## 参考


- [你需要知道的那些go语言json技巧](https://www.liwenzhou.com/posts/Go/json_tricks_in_go/)