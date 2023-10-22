<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [zerolog](#zerolog)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# zerolog

zerolog 包提供了一个专门用于 JSON 输出的简单快速的 Logger。

zerolog 的 API 旨在为开发者提供出色的体验和令人惊叹的性能[1]。其独特的链式 API 允许通过避免内存分配和反射来写入 JSON ( 或 CBOR ) 日志。

uber 的 zap 库开创了这种方法，zerolog 通过更简单的应用编程接口和更好的性能，将这一概念提升到了更高的层次