<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [yaml](#yaml)
  - [golang yaml v2和v3具体区别表现在v3拥有了许多新功能](#golang-yaml-v2%E5%92%8Cv3%E5%85%B7%E4%BD%93%E5%8C%BA%E5%88%AB%E8%A1%A8%E7%8E%B0%E5%9C%A8v3%E6%8B%A5%E6%9C%89%E4%BA%86%E8%AE%B8%E5%A4%9A%E6%96%B0%E5%8A%9F%E8%83%BD)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# yaml 

## golang yaml v2和v3具体区别表现在v3拥有了许多新功能
- 新的节点类型。虽然 YAML 文档通常被编码和解码为更高级别的类型，例如结构和映射，但 Node 是一种中间表示，允许对正在解码或编码的内容进行详细控制。
- 评论解码和编码。这可能是自该软件包被构思以来最需要的功能。
- 支持锚点和别名。新的 Node API 也支持锚点和别名操作，允许进行细粒度控制。
- Unmarshaler 接口现在使用 Node。Node API 是一种更强大的自定义解码方式，因此实现新的 yaml.Unmarshaler 接口的类型将采用 Node 作为参数。
- MapSlice 消失了。MapSlice 是一种在映射中自定义排序的简单方法，Node 提供了比这更多的功能，因此 MapSlice 完全从 API 中消失了。
- 解码值不再归零。在 v2 之前，yaml 包会将被解码的值清零，以便在迭代内部进行解码，避免之前迭代的剩余数据。 作为副作用，处理出于充分理由预初始化的类型变得更加困难。所以在 v3 中，这种行为发生了变化，这意味着预先存在的数据将被单独留在任何未正确解码的字段中。结果值是先前数据和成功解码的字段的总和