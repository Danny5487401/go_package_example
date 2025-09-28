<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [DI（Dependency Injection Container  依赖注入容器）](#didependency-injection-container--%E4%BE%9D%E8%B5%96%E6%B3%A8%E5%85%A5%E5%AE%B9%E5%99%A8)
  - [工厂模式和DI容器的区别](#%E5%B7%A5%E5%8E%82%E6%A8%A1%E5%BC%8F%E5%92%8Cdi%E5%AE%B9%E5%99%A8%E7%9A%84%E5%8C%BA%E5%88%AB)
  - [DI 容器的核心功能](#di-%E5%AE%B9%E5%99%A8%E7%9A%84%E6%A0%B8%E5%BF%83%E5%8A%9F%E8%83%BD)
  - [第三方实现](#%E7%AC%AC%E4%B8%89%E6%96%B9%E5%AE%9E%E7%8E%B0)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# DI（Dependency Injection Container  依赖注入容器）

DI框架通常提供两种功能：
1. 一种“提供”新组件的机制。
这将告诉DI框架您需要构建自己的其他组件（您的依赖关系）以及在拥有这些组件后如何构建自己。

2. 一种“检索”构建组件的机制。
DI框架通常基于您所讲述的“提供者”构建一个图并确定如何构建您的对象

## 工厂模式和DI容器的区别

DI容器底层最基本的设计思路就是基于工厂模式的。DI容器相当于一个大的工厂类，负责在程序启动的时候，根据配置（要创建哪些类对象，每个类对象的创建需要依赖哪些其他类对象）事先创建好对象。
当应用程序需要使用某个类对象的时候，直接从容器中获取即可。正是因为它持有一堆对象，所以这个框架才被称为“容器”。


## DI 容器的核心功能
DI容器的核心功能一般有三个：配置解析、对象创建和对象生命周期管理

## 第三方实现

常用的主要是 google/wire, github.com/facebookarchive/inject( 2019 archive), uber/dig, uber/fx 等。

大体上看，分为两个派系：

- 代码生成 codegen
- 基于反射 reflect

dig 和 wire 对比

- dig 通过反射识别依赖关系，谷歌出的wire，这个是用抽象语法树在编译时实现的。
- uber出的dig，在运行时，用返射实现的，并基于dig库，写了一个依赖框架fx. dig 只能在代码运行时，才能知道哪个依赖不对，比如构造函数返回类型的是结构体指针，但是其他依赖的是interface，这样的错误只能在运行时发现，而wire可以在编译的时候就发现。