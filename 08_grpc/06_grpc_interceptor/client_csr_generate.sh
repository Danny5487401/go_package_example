#生成 Key
openssl ecparam -genkey -name secp384r1 -out client.key

#生成CSR: Cerificate Signing Request 的英文缩写
#为证书请求文件。主要作用是 CA 会利用 CSR 文件进行签名使得攻击者无法伪装或篡改原有证书
openssl req -new -key client.key -out client.csr

#基于 CA 签发
openssl x509 -req -sha256 -CA ca.pem -CAkey ca.key -CAcreateserial -days 3650 -in client.csr -out client.pem
