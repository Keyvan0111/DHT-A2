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
		var clusterNodes []models.ClusterNodes
		if err := c.ShouldBindJSON(&clusterNodes); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"InputError" : "body did not match expected structure _ClusterNodes_"})
			return
		}
		utils.SortNodes(clusterNodes)
		myNode.Nodes = clusterNodes;
		utils.SetPeers(myNode, clusterNodes)
		utils.FingerTableInit(myNode)
		

		fmt.Println("My fingertable:")
		for _, entry := range myNode.FingerTable {
			fmt.Printf("%s      |      %d\n", entry.Node.Host, entry.Key)
		}

		c.JSON(http.StatusOK, gin.H{"message": "got all nodes "})
	}
}
