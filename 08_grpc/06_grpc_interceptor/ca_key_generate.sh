#CA
#了保证证书的可靠性和有效性，在这里可引入 CA 颁发的根证书的概念。其遵守 X.509 标准

#生成RSA私钥(无加密)
openssl genrsa -out ca.key 2048

# 生成密钥
openssl req -new -x509 -days 7200 -key ca.key -out ca.pem

#-new：表示生成一个新的证书签署请求；
#-x509：专用于生成CA自签证书；
#-key：指定生成证书用到的私钥文件；
#-out FILNAME：指定生成的证书的保存路径；
#-days：指定证书的有效期限，单位为day，默认是365天；


#填写信息
  #Country Name (2 letter code) []:
  #State or Province Name (full name) []:
  #Locality Name (eg, city) []:
  #Organization Name (eg, company) []:
  #Organizational Unit Name (eg, section) []:
  #Common Name (eg, fully qualified host name) []:go-grpc-example   (注意这名字)
  #Email Address []:





