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

	log.Println("call GetAomObj")
	nodes := models.GetAomObjList(200)
	edges := models.GetAomObjRelationship(200)

	e := elements{Nodes: nodes, Edges: edges}

	c.JSON(http.StatusOK, gin.H{"elements": e})
}
