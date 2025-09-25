package controllers

import (
	"main/models"
	"net/http"
	"fmt"

	"github.com/gin-gonic/gin"
)

func SendAllNodes(myNode *models.Node) gin.HandlerFunc {
	return func(c *gin.Context) {
		var clusterNodes []models.ClusterNodes
		if err := c.ShouldBindJSON(&clusterNodes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"InputError" : "body did not match expected structure _ClusterNodes_"})
			return
		}

		fmt.Printf("Im node: %s", myNode.Addr)

		for _, node:= range clusterNodes {
			fmt.Println(node)
		}

		fmt.Println("")
		c.JSON(http.StatusOK, gin.H{"message": "got all nodes "})
	}

}
