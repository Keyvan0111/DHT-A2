package routes

import (
	"github.com/gin-gonic/gin"

	"main/internal/controllers"
	"main/models"
)

func SetupClusterRoutes(router *gin.Engine, clusterNode *models.Node) {
	clusterGroup := router.Group("/cluster")
	{
		clusterGroup.GET("/forward/:key", controllers.ForwardNode(clusterNode))
		clusterGroup.POST("/fetch_nodes", controllers.SendAllNodes(clusterNode))
		clusterGroup.GET("/members", controllers.FetchMembers(clusterNode))
	}
}

func SetupStorageRoutes(router *gin.Engine, clusterNode *models.Node) {
	storageGroup := router.Group("/storage")
	{
		storageGroup.GET("/:key", controllers.GetValue(clusterNode))
		storageGroup.PUT("/:key", controllers.PutValue(clusterNode))
	}

}

func SetupNetworkRoutes(router *gin.Engine, clusterNode *models.Node) {
	networkGroup := router.Group("/network")
	{
		ng.GET("", controllers.NetworkPeers(n))     // serves /network
		ng.GET("/info", controllers.NetworkInfo(n)) // serves /network/info
	}

}

func SetupNodeRoutes(r *gin.Engine, n *models.Node) {
	r.GET("/node-info", controllers.NodeInfoHandler(n))
	r.POST("/join", controllers.JoinNetwork(n))
	r.POST("/leave", controllers.LeaveNetwork(n))
	r.POST("/sim-crash", controllers.SimCrash(n))
	r.POST("/sim-recover", controllers.SimRecover(n))
}
