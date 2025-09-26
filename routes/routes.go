package routes

import (
	"github.com/gin-gonic/gin"

	"main/internal/controllers"
	"main/models"
)

func SetupClusterRoutes(router *gin.Engine, clusterNode *models.Node) {
	clusterGroup := router.Group("/cluster")
	{
		clusterGroup.POST("/fetch_nodes", controllers.SendAllNodes(clusterNode))
	}

	router.GET("/network", controllers.NetworkInfo(clusterNode))
}

func SetupStorageRoutes(router *gin.Engine, clusterNode *models.Node) {
	storageGroup := router.Group("/storage")
	{
		storageGroup.GET("/:key", controllers.GetValue(clusterNode))
		storageGroup.POST("/:key", controllers.PutValue(clusterNode))
	}
	
}

func SetupNetworkRoutes(router *gin.Engine, clusterNode *models.Node) {
	networkGroup := router.Group("/network")
	{
		networkGroup.GET("/", controllers.NetworkInfo(clusterNode))
	}
	
}
