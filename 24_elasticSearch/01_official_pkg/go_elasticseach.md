<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Go官方包源码分析](#go%E5%AE%98%E6%96%B9%E5%8C%85%E6%BA%90%E7%A0%81%E5%88%86%E6%9E%90)
  - [日志](#%E6%97%A5%E5%BF%97)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# Go官方包源码分析

## 日志
默认日志有以下四种：
- estransport.TextLogger: 将请求和响应的基本信息以明文的形式输出
- estransport.ColorLogger:在开发时能在终端将一些信息以不同颜色输出
- estransport.CurlLogger:将信息格式化为可运行的curl命令，当启用EnableResponseBody时会美化输出
- estransport.JSONLogger:将信息以 json 格式输出，适用于生产环境的日志