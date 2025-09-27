package controllers

import (
	"main/models"
	// "main/utils"

	"github.com/gin-gonic/gin"
)

/*
Should find the closest node from the fingertable to the key
*/
func ForwardNode(mynode *models.Node) gin.HandlerFunc {
	return func(ctx *gin.Context) {
	// 	key := c.Param("key")
	// usatils.FindSuccessor(key, mynode)
	}
}
