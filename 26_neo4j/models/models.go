package models

import (
	"fmt"
	"go_grpc_example/26_neo4j/database"
	"log"
	"strconv"
)

type Aomobj struct {
	NodeId    string `json:"id" 			form:"id"`
	ObjId     string `json:"objId" 		form:"objId"`
	ObjString string `json:"name" 		form:"name"`
}

type Aomedge struct {
	EdgeId   string `json:"id" form:"id"`
	EdgeName string `json:"relationship" form:"relationship"`
	StartId  string `json:"source" form:"source"`
	EndId    string `json:"target"  form:"target"`
}

type NodeData struct {
	Data Aomobj `json:"data" form:"data"`
}
type EdgeData struct {
	Data Aomedge `json:"data" form:"data"`
}

var (
	neo4jURL = "neo4j+s://28b56052.databases.neo4j.io"
)

func GetAomObjList(count int32) (nodes []NodeData) {

	nodes = make([]NodeData, 0)

	driver, err := database.CreateDriver(neo4jURL, "neo4j", "QHX9NbyAbhRihh266kuD7DZO8MobVFL-VjX9yWi1qt4")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil
	}

	data, err := database.NodeQuery(driver, fmt.Sprintf("MATCH (n:AOM) RETURN  n LIMIT %d", count), "")

	//log.Println("Lenght of data is :", len(data))
	for i := 0; i < len(data); i++ {
		var node NodeData
		node.Data.NodeId = strconv.FormatInt(data[i].Id, 10)
		node.Data.ObjId = data[i].Props["OID"].(string)      /// OID 是我自己创建 neo4j db entry 的时候，添加的私有属性
		node.Data.ObjString = data[i].Props["Desc"].(string) /// "Desc" is same with "OID"， 私有属性

		nodes = append(nodes, node)
	}

	database.CloseDriver(driver)
	return nodes
}

func GetAomObjRelationship(count int32) (Edges []EdgeData) {

	Edges = make([]EdgeData, 0)

	driver, err := database.CreateDriver(neo4jURL, "neo4j", "neo4j")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil
	}

	data, err := database.EdgeQuery(driver, fmt.Sprintf("MATCH ()-[r]->() RETURN r LIMIT %d", count), "")

	//log.Println("Lenght of data is :", len(data))
	for i := 0; i < len(data); i++ {
		var edge EdgeData
		edge.Data.EdgeId = "r" + strconv.FormatInt(data[i].Id, 10)
		edge.Data.EdgeName = data[i].Type
		edge.Data.StartId = strconv.FormatInt(data[i].StartId, 10)
		edge.Data.EndId = strconv.FormatInt(data[i].EndId, 10)
		Edges = append(Edges, edge)
	}

	database.CloseDriver(driver)
	return Edges
}
