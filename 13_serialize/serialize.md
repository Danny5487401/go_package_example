# 序列化Marshal和反序列化UnMarshal

## Go 原生 encoding/json
使用 json.Unmarshal 和 json.Marshal 函数，可以轻松将 JSON 格式的二进制数据反序列化到指定的 Go 结构体中，以及将 Go 结构体序列化为二进制流。
而对于未知结构或不确定结构的数据，则支持将二进制反序列化到 map[string]interface{} 类型中，使用 KV 的模式进行数据的存取
json特性
json 包解析的是一个 JSON 数据，而 JSON 数据既可以是对象（object），也可以是数组（array），同时也可以是字符串（string）、数值（number）、布尔值（boolean）以及空值（null）。
```go
var s string
err := json.Unmarshal([]byte(`"Hello, world!"`), &s)
// 注意字符串中的双引号不能缺，如果仅仅是 `Hello, world`，则这不是一个合法的 JSON 序列，会返回错误。
```

## 第三方包 jsoniter
从性能上，jsoniter 能够比众多大神联合开发的官方库性能还快的主要原因，一个是尽量减少不必要的内存复制，另一个是减少 reflect 的使用——同一类型的对象，jsoniter 只调用 reflect 解析一次之后即缓存下来。
不过随着 go 版本的迭代，原生 json 库的性能也越来越高，jsonter 的性能优势也越来越窄
