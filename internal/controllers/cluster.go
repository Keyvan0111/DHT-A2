package controllers

import (
	"main/models"
	"main/utils"

	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SendAllNodes(myNode *models.Node) gin.HandlerFunc {
	return func(c *gin.Context) {
		if GuardCrash(myNode, c) {
			return
		}
		var clusterNodes []models.ClusterNodes
		if err := c.ShouldBindJSON(&clusterNodes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"InputError" : "body did not match expected structure _ClusterNodes_"})
			return
		}
		utils.UpdateClusterView(myNode, clusterNodes)

		fmt.Println("My fingertable:")
		for _, entry := range myNode.FingerTable {
			fmt.Printf("%s      |      %d\n", entry.Node.Host, entry.Key)
		}

		c.JSON(http.StatusOK, gin.H{"message": "got all nodes "})
	}
}

func FetchMembers(myNode *models.Node) gin.HandlerFunc {
	return func(c *gin.Context) {
		if GuardCrash(myNode, c) {
			return
		}
		myNode.Guard.RLock()
		defer myNode.Guard.RUnlock()

		c.JSON(http.StatusOK, myNode.Nodes)
	}
}
