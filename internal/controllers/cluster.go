package controllers

import (
	"main/models"
	"main/utils"

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

		utils.SortNodes(clusterNodes)
		utils.SetPeers(myNode, clusterNodes)
		fmt.Println(myNode)

		fmt.Println("")
		c.JSON(http.StatusOK, gin.H{"message": "got all nodes "})
	}

}

func GetDataWithKey() gin.HandlerFunc {
	return func(c *gin.Context) {

		/*
		if key {
			c.JSON(http.statusOK, {data: DHT[key]})
			} else {
				
			c.JSON(http.statusNotFound, {message: "Key not found")
		}
			*/
	}
}

func GetKnownNodes() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if key {
		// 	c.JSON(http.statusOK, {data: Nodes})
		// 	}
	}
}

func SendData() gin.HandlerFunc {
	return func(c *gin.Context) {
		/*
		var addr string;
		c.Shouldbindjason(&addr); err != nil{
			http badrequest
		}

		business logic....

		c.JSON(http.statusOK, Message: "persistetd")
		*/

	}	
}
