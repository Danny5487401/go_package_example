package main

/*
phpserialize 序列化对象后，可以很方便的将它传递给其他需要它的地方，且其类型和结构不会改变。

<?php
$sites = array('Google', 'Runoob', 'Facebook');
$serialized_data = serialize($sites);
echo  $serialized_data . PHP_EOL;
?>

结果：包含字节的长度,先后顺序
a:3:{i:0;s:6:"Google";i:1;s:6:"Runoob";i:2;s:8:"Facebook";}
*/

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/elliotchance/phpserialize"
)

const (
	cookieValidationKey = "123456"
)

var key = "age"
var value = "18"

func main() {
	// 加密
	encodeStr := Encode(key, value)
	fmt.Printf("加密后的数据是%v\n", encodeStr)
	// 解密
	k, v := Decode(encodeStr)
	fmt.Printf("原始数据是%v：%v", k, v)
}

func Encode(key, value string) (encodeStr string) {
	data, err := phpserialize.Marshal([]string{key, value}, nil)
	if err != nil {
		fmt.Println("加密错误", err.Error())
		return
	}
	mac := hmac.New(sha256.New, []byte(cookieValidationKey))
	mac.Write(data)
	hash := fmt.Sprintf("%x", mac.Sum(nil))
	value = hash + string(data)
	fmt.Println(value)
	return value
}

func Decode(encodeStr string) (key, value string) {
	sData := encodeStr
	mac := hmac.New(sha256.New, []byte(""))
	_, _ = mac.Write([]byte(""))
	byteInfo := mac.Sum(nil)            // 256 位 = 32字节
	test := fmt.Sprintf("%x", byteInfo) // 转为16 进制，2^4=16，即256/4= 64位
	hashLength := len(test)             // len() 函数的返回值的类型为 int ，在64位机器这里是int64, int64= 8 * uint8=64位，
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
	for k1, v := range data {
		k := k1.(int64)
		if k == 0 {
			key = v.(string)
		} else if k == 1 {
			value = v.(string)
		} else {
			return
		}
	}
	return
}
