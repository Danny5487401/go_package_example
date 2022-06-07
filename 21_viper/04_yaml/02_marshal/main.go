package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type Users struct {
	Name    string   `yaml:"name"`
	Age     int8     `yaml:"age"`
	Address string   `yaml:"address"`
	Hobby   []string `yaml:"hobby"`
}

func main() {
	wanger := Users{
		Name:    "wanger",
		Age:     24,
		Address: "beijing",
		Hobby:   []string{"literature", "social"},
	}
	dongdong := Users{
		Name:    "冬哥",
		Age:     30,
		Address: "chengdu",
		Hobby:   []string{"basketball", "guitar"},
	}
	xialaoshi := Users{
		Name:    "夏老师",
		Age:     29,
		Address: "chengdu",
		Hobby:   []string{"吃吃喝喝"},
	}
	huazai := Users{
		Name:    "华子",
		Age:     28,
		Address: "shenzhen",
		Hobby:   []string{"王者荣耀"},
	}
	qiaoke := Users{
		Name:    "乔克",
		Age:     30,
		Address: "chongqing",
		Hobby:   []string{"阅读", "王者荣耀"},
	}
	jiangzong := Users{
		Name:    "姜总",
		Age:     25,
		Address: "shanghai",
		Hobby:   []string{"钓鱼", "音乐", "美食", "酒"},
	}
	zhengge := Users{
		Name:    "郑哥",
		Age:     30,
		Address: "beijing",
		Hobby:   []string{"阅读", "复读机"},
	}
	userlist := [7]Users{wanger, dongdong, huazai, qiaoke, xialaoshi, jiangzong, zhengge}

	yamlData, err := yaml.Marshal(&userlist)

	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}

	fmt.Println(string(yamlData))
	fileName := "21_viper/04_yaml/test.yaml"
	err = ioutil.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		panic("Unable to write data into the file")
	}
}
