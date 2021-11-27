package apis

import (
	"github.com/gin-gonic/gin"
	"go_grpc_example/26_neo4j/models"
	"log"
	"net/http"
)

type elements struct {
	Nodes []models.NodeData `json:"nodes" form:"nodes"`
	Edges []models.EdgeData `json:"edges" form:"edges"`
}

func GetAomObj(c *gin.Context) {

	log.Println("call GetStudTeacherRelationObj")
	// 获取节点
	nodes := models.GetStuObjList(200)
	// 获取关系
	edges := models.GetTeachStudentRelationship(200)

	e := elements{Nodes: nodes, Edges: edges}

	c.JSON(http.StatusOK, gin.H{"elements": e})
}

func CreateObj(c *gin.Context) {

	models.CreateStudent(c.Query("name"), c.Query("class"))
}
