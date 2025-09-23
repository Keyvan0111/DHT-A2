package routes

import 	(
	"github.com/gin-gonic/gin"

	"main/internal/controllers"
)

func SetupClusterRoutes(router *gin.Engine) {
	assistantGroup := router.Group("/cluster") 
	{
		assistantGroup.POST("/fetch_nodes", controllers.SendAllNodes())
	}
}
