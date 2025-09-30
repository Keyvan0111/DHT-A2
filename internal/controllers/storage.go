package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"main/models"
	"main/utils"
)

// GET
func GetValue(node *models.Node) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		_, keyId := utils.ConsistentHash(key)
		fmt.Printf("Getting Value: %s (id: %d) into table...\n\n", key, keyId)

		fmt.Printf("Im node: %d Checking responsibility...\n", node.NodeId)
		if utils.IsResponsibleFor(keyId, node) {
			fmt.Printf("Im responsible for %d\n", keyId)
			if v, ok := node.Store.Load(key); ok {
				c.String(http.StatusOK, "%s", v.(string))
				return
			}
			c.Status(http.StatusNotFound)
			return
		}
		fmt.Printf("Im not responsible for %d\n", keyId)
		fmt.Printf("Forwarding...\n\n")
		utils.ForwardGet(node, keyId, key, c, "/storage/")
	}
}

// PUT
func PutValue(node *models.Node) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")
		value, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
			return
		}
		fmt.Printf("Putting Value into table...\n\n")

		_, keyId := utils.ConsistentHash(key)
		fmt.Printf("Hashing the value: %s -> %d\n", key, keyId)

		fmt.Printf("Im node: %d Checking responsibility...\n", node.NodeId)
		if utils.IsResponsibleFor(keyId, node) {
			fmt.Printf("Im responsible for %d\n", keyId)
			node.Store.Store(key, string(value))
			c.Status(http.StatusOK)
			return
		}
		fmt.Printf("Im not responsible for %d\n", keyId)
		fmt.Printf("Forwarding...\n\n")
		utils.ForwardPut(node, keyId, key, value, c, "/storage/")
	}
}

// GET
func NetworkInfo(n *models.Node) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"self": gin.H{
				"addr": n.Addr,
				"id":   n.NodeId,
			},
			"predecessor": gin.H{
				"addr": n.Predecessor.Addr,
				"id":   n.Predecessor.NodeId,
			},
			"successor": gin.H{
				"addr": n.Successor.Addr,
				"id":   n.Successor.NodeId,
			},
			"hashlen": utils.HASHLEN,
		})
	}
}
