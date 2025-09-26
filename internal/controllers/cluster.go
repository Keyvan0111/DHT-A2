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
		utils.SetPeers(myNode, clusterNodes)
		fmt.Println("printing node...")
		fmt.Println(myNode)

		fmt.Println("")
		c.JSON(http.StatusOK, gin.H{"message": "got all nodes "})
	}
}
