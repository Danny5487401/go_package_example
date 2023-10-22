<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [OpenSSL是一个开源的项目，其由三个部分组成：](#openssl%E6%98%AF%E4%B8%80%E4%B8%AA%E5%BC%80%E6%BA%90%E7%9A%84%E9%A1%B9%E7%9B%AE%E5%85%B6%E7%94%B1%E4%B8%89%E4%B8%AA%E9%83%A8%E5%88%86%E7%BB%84%E6%88%90)
- [生成密钥对](#%E7%94%9F%E6%88%90%E5%AF%86%E9%92%A5%E5%AF%B9)
- [创建CA和申请证书](#%E5%88%9B%E5%BB%BAca%E5%92%8C%E7%94%B3%E8%AF%B7%E8%AF%81%E4%B9%A6)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# OpenSSL是一个开源的项目，其由三个部分组成：
1、openssl命令行工具；
2、libencrypt加密算法库；
3、libssl加密模块应用库；

# 生成密钥对
```shell
openssl genrsa [-out filename] [-passout arg] [-des] [-des3] [-idea] [-f4] [-3] [-rand file(s)] [-engine id] [numbits]
#常用选项：
#-out FILENAME：将生成的私钥保存至指定的文件中；
#[-des] [-des3] [-idea]：指定加密算法；
#numbits：指明生成的私钥大小，默认是512；

```
# 创建CA和申请证书
在使用OpenSSL命令创建证书前，可查看配置文件/etc/pki/tls/openss.conf文件，查看该文件定义了的证书存放位置及名称