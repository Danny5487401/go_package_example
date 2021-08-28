package main

/*

加密技术包括两个元素：算法和密钥
算法：
	将普通的信息或者可以理解的信息与一串数字（密钥）结合，产生不可理解的密文的步骤。

密钥：
	用来对数据进行编码和解密的一种算法。
加密解密方法分类：
	基本加密方法
	对称加密方法
	非对称加密方法
对称加密和非对称加密的区别:
	对称加密中只有一个钥匙也就是KEY,加解密都依靠这组密钥
	非对称加密中有公私钥之分,私钥可以生产公钥(比特币的钱包地址就是公钥),一般加密通过公钥加密，私钥解密(也有私钥加密公钥解密)
RSA使用场景:
	我们最熟悉的就是HTTPS中就是使用的RSA加密,CA机构给你颁发的就是私钥给到我们进行配置,在请求过程中端用CA内置到系统的公钥加密,
	请求道服务器由服务器进行解密验证,保障了传输过程中的请求加密

	高安全场景(比如金融设备银联交易等)下的双向认证(一机一密钥),每台机器本地都会生成一组公私钥对,并且吧公钥发送给服务器,
	这个使用发起的请求模型如下
*/

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/elliotchance/phpserialize"
	"time"
)

const (
	cookieValidationKey = "123456"
)

var name = "UAF"
var value = time.Now().Format("20060102")

func main() {
	// 加密
	Decode()

}

func Encode() (encodeStr string) {
	data, err := phpserialize.Marshal([]string{name, value}, nil)
	if err != nil {
		fmt.Println("加密错误", err.Error())
	}
	mac := hmac.New(sha256.New, []byte(cookieValidationKey))
	mac.Write(data)
	hash := fmt.Sprintf("%x", mac.Sum(nil))
	value = hash + string(data)
	fmt.Println(value)
	return value
}

func Decode() {
	sData := Encode()
	mac := hmac.New(sha256.New, []byte(""))
	_, _ = mac.Write([]byte(""))
	test := fmt.Sprintf("%x", mac.Sum(nil))
	hashLength := len(test)
	if len(sData) < hashLength {
		return
	}

	hash := sData[0:hashLength]
	pureData := sData[hashLength:]
	mac2 := hmac.New(sha256.New, []byte(cookieValidationKey))
	_, _ = mac2.Write([]byte(pureData))
	if hash != fmt.Sprintf("%x", mac2.Sum(nil)) {
		return
	}
	// 解密后的数据
	var data map[interface{}]interface{}
	err := phpserialize.Unmarshal([]byte(pureData), &data)
	if err != nil {
		return
	}
	for k, v := range data {
		if k.(int64) == 1 {
			fmt.Println(v.(string))
		}
	}
}
