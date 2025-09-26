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
func GetValue(n *models.Node) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.Param("key")
        _, keyId := utils.ConsistentHash(key)

        if utils.IsResponsibleFor(keyId, n) {
            if v, ok := n.Store.Load(key); ok {
                c.String(http.StatusOK, "%s", v.(string))
                return
            }
            c.Status(http.StatusNotFound)
            return
        }
        utils.ForwardGet(n, key, c)
    }
}

// PUT
func PutValue(n *models.Node) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.Param("key")
        value, err := io.ReadAll(c.Request.Body)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
            return
        }

		
        _, keyId := utils.ConsistentHash(key)
		fmt.Println("value: ", string(value), "key: ", keyId)

        if utils.x(keyId, n) {
			fmt.Printf("Keyid: %d, pred: %d", keyId, n.Predecessor.NodeId)
            n.Store.Store(key, string(value))
            c.Status(http.StatusOK)
            return
        }
        utils.ForwardPut(n, key, value, c)
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
