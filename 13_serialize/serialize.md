# 序列化Marshal和反序列化UnMarshal

- 序列化（编码）是将对象序列化为二进制形式（字节数组），主要用于网络传输、数据持久化等；
- 而反序列化（解码）则是将从网络、磁盘等读取的字节数组还原成原始对象，主要用于网络传输对象的解码，以便完成远程调用。

影响序列化性能的关键因素：序列化后的码流大小（网络带宽的占用）、序列化的性能（CPU资源占用）；是否支持跨语言（异构系统的对接和开发语言切换）。

## 常见序列化协议

### xml（Extensible Markup Language）
- 优点：人机可读性好，可指定元素或特性的名称。
- 缺点：序列化数据只包含数据本身以及类的结构，不包括类型标识和程序集信息；只能序列化公共属性和字段；不能序列化方法；文件庞大，文件格式复杂，传输占带宽。

适用场景：当做配置文件存储数据，实时数据转换。

### JSON(JavaScript Object Notation, JS 对象标记) 
是一种轻量级的数据交换格式，
- 优点：兼容性高、数据格式比较简单，易于读写、序列化后数据较小，可扩展性好，兼容性好、与XML相比，其协议比较简单，解析速度比较快。
- 缺点：数据的描述性比XML差、不适合性能要求为ms级别的情况、额外空间开销比较大。

适用场景（可替代ＸＭＬ）：跨防火墙访问、可调式性要求高、基于Web browser的Ajax请求、传输数据量相对小，实时性要求相对低（例如秒级别）的服务。

### Thrift
不仅是序列化协议，还是一个RPC框架
- 优点：序列化后的体积小, 速度快、支持多种语言和丰富的数据类型、对于数据字段的增删具有较强的兼容性、支持二进制压缩编码。
- 缺点：使用者较少、跨防火墙访问时，不安全、不具有可读性，调试代码时相对困难、不能与其他传输层协议共同使用（例如HTTP）、无法支持向持久层直接读写数据，即不适合做数据持久化序列化协议。

适用场景：分布式系统的RPC解决方案

### Avro
Hadoop的一个子项目，解决了JSON的冗长和没有IDL的问题。

- 优点：支持丰富的数据类型、简单的动态语言结合功能、具有自我描述属性、提高了数据解析速度、快速可压缩的二进制数据形式、可以实现远程过程调用RPC、支持跨编程语言实现。
- 缺点：对于习惯于静态类型语言的用户不直观。

适用场景：在Hadoop中做Hive、Pig和MapReduce的持久化数据格式。

## Protobuf
将数据结构以.proto文件进行描述，通过代码生成工具可以生成对应数据结构的POJO对象和Protobuf相关的方法和属性
- 优点：序列化后码流小，性能高、结构化数据存储格式（XML JSON等）、通过标识字段的顺序，可以实现协议的前向兼容、结构化的文档更容易管理和维护。
- 缺点：需要依赖于工具生成代码、支持的语言相对较少，官方只支持Java 、C++ 、python,但是可以扩展。

适用场景：对性能要求高的RPC调用、具有良好的跨防火墙的访问属性、适合应用层对象的持久化




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
