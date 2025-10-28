package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"main/models"
)

func GuardCrash(node *models.Node, c *gin.Context) bool {
	node.Guard.RLock()
	crashed := node.State == models.NodeStateCrashed
	node.Guard.RUnlock()

	if crashed {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "node unavailable"})
		return true
	}
	return false
}
