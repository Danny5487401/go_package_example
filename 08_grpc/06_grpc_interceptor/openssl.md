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