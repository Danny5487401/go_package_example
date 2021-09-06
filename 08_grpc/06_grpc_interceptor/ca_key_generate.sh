#CA
#了保证证书的可靠性和有效性，在这里可引入 CA 颁发的根证书的概念。其遵守 X.509 标准

#生成RSA私钥(无加密)
openssl genrsa -out ca.key 2048

# 生成密钥
openssl req -new -x509 -days 7200 -key ca.key -out ca.pem

#填写信息
  #Country Name (2 letter code) []:
  #State or Province Name (full name) []:
  #Locality Name (eg, city) []:
  #Organization Name (eg, company) []:
  #Organizational Unit Name (eg, section) []:
  #Common Name (eg, fully qualified host name) []:go-grpc-example   (注意这名字)
  #Email Address []:


