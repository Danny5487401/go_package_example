<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Go 原生 encoding/json](#go-%E5%8E%9F%E7%94%9F-encodingjson)
  - [序列化](#%E5%BA%8F%E5%88%97%E5%8C%96)
  - [反序列化](#%E5%8F%8D%E5%BA%8F%E5%88%97%E5%8C%96)
- [第三方包 jsoniter](#%E7%AC%AC%E4%B8%89%E6%96%B9%E5%8C%85-jsoniter)
- [反序列化](#%E5%8F%8D%E5%BA%8F%E5%88%97%E5%8C%96-1)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Go 原生 encoding/json
使用 json.Unmarshal 和 json.Marshal 函数，可以轻松将 JSON 格式的二进制数据反序列化到指定的 Go 结构体中，以及将 Go 结构体序列化为二进制流。
而对于未知结构或不确定结构的数据，则支持将二进制反序列化到 map[string]interface{} 类型中，使用 KV 的模式进行数据的存取

特性：
json 包解析的是一个 JSON 数据，而 JSON 数据既可以是对象（object），也可以是数组（array），同时也可以是字符串（string）、数值（number）、布尔值（boolean）以及空值（null）。

### 序列化

问题：map是无序的，每次取出key/value的顺序都可能不一致，但map转json的顺序是不是也是无序的吗？

结论：map转json是有序的，按照ASCII码升序排列key。

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

### 反序列化
```go
var s string
err := json.Unmarshal([]byte(`"Hello, world!"`), &s)
// 注意字符串中的双引号不能缺，如果仅仅是 `Hello, world`，则这不是一个合法的 JSON 序列，会返回错误。
```

```go
// encoding/json/decode.go
func Unmarshal(data []byte, v interface{}) error {
	// Check for well-formedness.
	// Avoids filling out half a data structure
	// before discovering a JSON syntax error.
	var d decodeState
	err := checkValid(data, &d.scan)
	if err != nil {
		return err
	}

	d.init(data)
	return d.unmarshal(v)
}


func (d *decodeState) unmarshal(v interface{}) error {
	rv := reflect.ValueOf(v)
	// 保证v是指针,即结果
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}

    // 。。。
	// We decode rv not rv.Elem because the Unmarshaler interface
	// test must be applied at the top level of the value.
	err := d.value(rv)
    // ...
}

func (d *decodeState) value(v reflect.Value) error {
    switch d.opcode {
    default:
        panic(phasePanicMsg)
    // 数组 
    case scanBeginArray:
        ...
    // 结构体或map
    case scanBeginObject:
		if v.IsValid() {
            if err := d.object(v); err != nil {
                return err
            }
        } else {
            d.skip()
        }
        d.scanNext()
        ...
    // 字面量，包括 int、string、float 等
    case scanBeginLiteral:
        ...
    }
    return nil
}
```

以解析对象为例：
```go
func (d *decodeState) object(v reflect.Value) error {
    u, ut, pv := indirect(v, false)
    // ...
    v = pv
    t := v.Type()
    ...  
    var fields structFields
    // 检验这个对象的类型是 map 还是 结构体
    switch v.Kind() {
    case reflect.Map: 
        ...
    case reflect.Struct:
        // 缓存结构体的字段到 fields 对象中
        fields = cachedTypeFields(t)
        // ok
    default:
        d.saveError(&UnmarshalTypeError{Value: "object", Type: t, Offset: int64(d.off)})
        d.skip()
        return nil
    }

    var mapElem reflect.Value
    origErrorContext := d.errorContext
    // 循环一个个解析JSON字符串中的 key value 值
    for {  
        start := d.readIndex()
        d.rescanLiteral()
        item := d.data[start:d.readIndex()]
        // 获取 key 值
        key, ok := unquoteBytes(item)
        if !ok {
            panic(phasePanicMsg)
        } 
        var subv reflect.Value
        destring := false   
        ... 
        // 根据 value 的类型反射设置 value 值 
        if destring {
            // value 值是字面量会进入到这里
            switch qv := d.valueQuoted().(type) {
            case nil:
                if err := d.literalStore(nullLiteral, subv, false); err != nil {
                    return err
                }
            case string:
                if err := d.literalStore([]byte(qv), subv, true); err != nil {
                    return err
                }
            default:
                d.saveError(fmt.Errorf("json: invalid use of ,string struct tag, trying to unmarshal unquoted value into %v", subv.Type()))
            }
        } else {
            // 数组或对象会递归调用 value 方法
            if err := d.value(subv); err != nil {
                return err
            }
        }
        ...
        // 直到遇到 } 最后退出循环
        if d.opcode == scanEndObject {
            break
        }
        if d.opcode != scanObjectValue {
            panic(phasePanicMsg)
        }
    }
    return nil
}
```
流程：
1. 首先会缓存结构体对象；
2. 循环遍历结构体对象；
3. 找到结构体中的 key 值之后再找到结构体中同名字段类型；
4. 递归调用 value 方法反射设置结构体对应的值；
5. 直到遍历到 JSON 中结尾 }结束循环。

Note：通过看 Unmarshal 源码中可以看到其中使用了大量的反射来获取字段值，如果是多层嵌套的 JSON 的话，那么还需要递归进行反射获取值，可想而知性能是非常差的了。


## 第三方包 jsoniter
从性能上，jsoniter 能够比众多大神联合开发的官方库性能还快的主要原因，一个是尽量减少不必要的内存复制，另一个是减少 reflect 的使用——同一类型的对象，jsoniter 只调用 reflect 解析一次之后即缓存下来。
不过随着 go 版本的迭代，原生 json 库的性能也越来越高，jsonter 的性能优势也越来越窄

## 反序列化

1. json字符串解析时，需要一个“接收体”接受解析后的数据，且Unmarshal时接收体必须传递指针。否则解析虽不报错，但数据无法赋值到接受体中。如这里用的是StuRead{}接收。

2. 解析时，接收体可自行定义。json串中的key自动在接收体中寻找匹配的项进行赋值。匹配规则是：
    - 先查找与key一样的json标签，找到则赋值给该标签对应的变量(如Name)。
    - 没有json标签的，就从上往下依次查找变量名与key一样的变量，如Age。或者变量名忽略大小写后与key一样的变量。如HIgh，Class。第一个匹配的就赋值，后面就算有匹配的也忽略。
    (前提是该变量必需是可导出的，即首字母大写)。
3. 不可导出的变量无法被解析（如sex变量，虽然json串中有key为sex的k-v，解析后其值仍为nil,即空值）

4. 当接收体中存在json串中匹配不了的项时，解析会自动忽略该项，该项仍保留原值。如变量Test，保留空值nil。

5. 你一定会发现，变量Class貌似没有解析为我们期待样子。因为此时的Class是个interface{}类型的变量，而json串中key为CLASS的value是个复合结构，
不是可以直接解析的简单类型数据（如“张三”，18，true等）。
所以解析时，由于没有指定变量Class的具体类型，json自动将value为复合结构的数据解析为map[string]interface{}类型的项。
也就是说，此时的struct Class对象与StuRead中的Class变量没有半毛钱关系，故与这次的json解析没有半毛钱关系。


interface{}类型变量在json解析前，打印出的类型都为nil，就是没有具体类型，这是空接口（interface{}类型）的特点。

json解析后，json串中value，只要是”简单数据”，都会按照默认的类型赋值，如”张三”被赋值成string类型到Name变量中，数字18对应float64，true对应bool类型。
“简单数据”：是指不能再进行二次json解析的数据，如”name”:”张三”只能进行一次json解析。
“复合数据”：类似”CLASS\”:{\”naME\”:\”1班\”,\”GradE\”:3}这样的数据，是可进行二次甚至多次json解析的，因为它的value也是个可被解析的独立json。即第一次解析key为CLASS的value，第二次解析value中的key为naME和GradE的value

对于”复合数据”，如果接收体中配的项被声明为interface{}类型，go都会默认解析成map[string]interface{}类型。如果我们想直接解析到struct Class对象中，可以将接受体对应的项定义为该struct类型
```go
type StuRead struct {
	//普通struct类型
	Class Class `json:"class"`
	//指针类型
	Class *Class `json:"class"`
}
```
打印结果
```shell
Class类型：{张三 18 true <nil> {1班 3} <nil>}
*Class类型：{张三 18 true <nil> 0xc42008a0c0 <nil>}
```
可以看出，传递Class类型的指针时，stu中的Class变量存的是指针，我们可通过该指针直接访问所属的数据，如stu.Class.Name/stu.Class.Grade
