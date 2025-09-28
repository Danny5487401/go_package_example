<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [github.com/google/wire](#githubcomgooglewire)
  - [两个核心概念：提供者（providers）和注入器（injectors）](#%E4%B8%A4%E4%B8%AA%E6%A0%B8%E5%BF%83%E6%A6%82%E5%BF%B5%E6%8F%90%E4%BE%9B%E8%80%85providers%E5%92%8C%E6%B3%A8%E5%85%A5%E5%99%A8injectors)
  - [Wire 核心技术](#wire-%E6%A0%B8%E5%BF%83%E6%8A%80%E6%9C%AF)
  - [参考](#%E5%8F%82%E8%80%83)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# github.com/google/wire

Wire 是一个轻巧的 Golang 依赖注入工具。它由 Go Cloud 团队开发，通过自动生成代码的方式在编译期完成依赖注入。


Wire 分成两部分，一个是在项目中使用的依赖， 一个是命令行工具。

## 两个核心概念：提供者（providers）和注入器（injectors）


Provider: 生成组件的普通方法。这些方法接收所需依赖作为参数，创建组件并将其返回。


Injector: 由wire自动生成的函数。函数内部会按根据依赖顺序调用相关 provider

## Wire 核心技术


## 参考

