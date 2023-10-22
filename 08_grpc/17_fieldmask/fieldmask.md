<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [field masks](#field-masks)
  - [参考资料](#%E5%8F%82%E8%80%83%E8%B5%84%E6%96%99)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# field masks


```protobuf
message UpdateBookRequest {
    // 操作人 
    string op = 1;
    // 要更新的书籍信息
    Book book = 2;
}
```

但是如果我们的Book中定义有很多很多字段时，我们不太可能每次请求都去全量更新Book的每个字段，因为通常每次操作只会更新1到2个字段。

答案是使用google/protobuf/field_mask.proto，它能够记录在一次更新请求中涉及到的具体字段路径。

## 参考资料
- [官方field masks](https://google.aip.dev/161)