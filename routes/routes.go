package routes

import (
	"github.com/gin-gonic/gin"

	"main/internal/controllers"
	"main/models"
)

func SetupClusterRoutes(router *gin.Engine, clusterNode *models.Node) {
	assistantGroup := router.Group("/cluster")
	{
		assistantGroup.POST("/fetch_nodes", controllers.SendAllNodes(clusterNode))
	}

	 	router.GET("/network", controllers.NetworkInfo(clusterNode))
  	router.GET("/storage/:key", controllers.GetValue(clusterNode))
    router.PUT("/storage/:key", controllers.PutValue(clusterNode))
}
