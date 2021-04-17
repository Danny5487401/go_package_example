package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
)

type Student struct {
	Name string `json:"name"`
	Age int `json:"age"`
	IsMarried bool `json:"isMarried"`
	
}
// 注意elasticSearch 版本
func main()  {
	// 初始化连接
	client,err := elastic.NewClient(elastic.SetSniff(false),elastic.SetURL("http://81.68.197.3:9200"))
	if err != nil{
		panic(err)
	}
	fmt.Println("connected to elasticSearch")
	//s1 := &Student{
	//	Name: "Durant",
	//	Age: 50,
	//	IsMarried:true,
	//}
	// 1. 添加数据
	// 链式操作
 	//indexRsp, err := client.Index().Index("student").Type("go").BodyJson(s1).Do(context.Background())
 	//if err != nil{
 	//	panic(err)
	//}
	//fmt.Printf("StudentId: %s ,index :%s ,type :%s\n",indexRsp.Id,indexRsp.Index,indexRsp.Type)


 	// 2。 查询数据
 	var stud []Student  // 用于收集数据

 	// 构建搜索资源
 	searchRes := elastic.NewSearchSource()
 	searchRes.Query(elastic.NewMatchQuery("isMarried","true"))
 	searchRsp,err := client.Search().Index("student").SearchSource(searchRes).Do(context.Background())
	if err != nil{
		panic(err)
	}
 	for _,hit := range searchRsp.Hits.Hits{
		var student Student
		err := json.Unmarshal(*hit.Source, &student)
		if err != nil {
			fmt.Println("[Getting Students][Unmarshal] Err=", err)
		}

		stud = append(stud, student)
	}
	for _, s := range stud {
		fmt.Printf("Student found Name: %s, Age: %d, Score: %v \n", s.Name, s.Age, s.IsMarried)
	}



}
