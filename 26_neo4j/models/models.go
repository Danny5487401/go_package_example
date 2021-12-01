package models

import (
	"fmt"
	"go_grpc_example/26_neo4j/database"
	"log"
	"strconv"
)

type StuObj struct {
	NodeId string `json:"id" form:"id"`
	Class  string `json:"class" form:"class"`
	Name   string `json:"name" form:"name"`
}

type StudentTeacherdge struct {
	EdgeId  string `json:"id" form:"id"`
	Teach   string `json:"teach" form:"teach"`
	StartId string `json:"source" form:"source"`
	EndId   string `json:"target"  form:"target"`
}

type NodeData struct {
	Data StuObj `json:"data" form:"data"`
}
type EdgeData struct {
	Data StudentTeacherdge `json:"data" form:"data"`
}

var (
	neo4jURL = "bolt://tencent.danny.games:7687"
	db       = "neo4j"
)

func CreateStudent(name string, class string) {
	driver, err := database.CreateDriver(neo4jURL, "neo4j", "chuanzhi")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return
	}
	sql := fmt.Sprintf(`create n:Student{name:%v, class: %v } `, name, class)
	_ = database.NodeCreate(driver, sql, db)
}

func GetStuObjList(count int32) (nodes []NodeData) {

	nodes = make([]NodeData, 0)

	driver, err := database.CreateDriver(neo4jURL, "neo4j", "chuanzhi")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil
	}

	data, err := database.NodeQuery(driver, fmt.Sprintf("match(n:Student) return (n) LIMIT %d", count), db)
	if err != nil {
		return nil
	}
	//log.Println("Lenght of data is :", len(data))
	for i := 0; i < len(data); i++ {
		var node NodeData
		node.Data.NodeId = strconv.FormatInt(data[i].Id, 10)
		name, ok := data[i].Props["name"].(string)
		if ok {
			node.Data.Name = name
		}
		class, ok := data[i].Props["class"].(string)
		if ok {
			node.Data.Class = class
		}

		nodes = append(nodes, node)
	}

	database.CloseDriver(driver)
	return nodes
}

func GetTeachStudentRelationship(count int32) (Edges []EdgeData) {

	Edges = make([]EdgeData, 0)

	driver, err := database.CreateDriver(neo4jURL, "neo4j", "chuanzhi")
	if err != nil {
		log.Println("error connecting to neo4j:", err)
		return nil
	}

	data, err := database.EdgeQuery(driver, fmt.Sprintf("MATCH ()-[r]->() RETURN r LIMIT %d", count), "")

	//log.Println("Lenght of data is :", len(data))
	for i := 0; i < len(data); i++ {
		var edge EdgeData
		edge.Data.EdgeId = "r" + strconv.FormatInt(data[i].Id, 10)
		edge.Data.Teach = data[i].Type
		edge.Data.StartId = strconv.FormatInt(data[i].StartId, 10)
		edge.Data.EndId = strconv.FormatInt(data[i].EndId, 10)
		Edges = append(Edges, edge)
	}

	database.CloseDriver(driver)
	return Edges
}
