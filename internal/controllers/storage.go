package controllers

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"main/models"
	"main/utils"
)


func between(id, start, end int) bool {
    if start < end {
        return id > start && id <= end
    }
    return id > start || id <= end
}

func responsibleFor(keyId int, n *models.Node) bool {

    return between(keyId, n.PredecessorId, n.NodeId)
}

func forwardGet(n *models.Node, key string, c *gin.Context) {
    // Linear implementation for now --> just hop to successor
    url := n.SuccessorAddr + "/storage/" + key
    resp, err := http.Get(url)
    if err != nil {
        c.JSON(http.StatusBadGateway, gin.H{"error": "forward GET failed", "detail": err.Error()})
        return
    }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    c.Data(resp.StatusCode, "text/plain; charset=utf-8", body)
}

func forwardPut(n *models.Node, key string, value []byte, c *gin.Context) {
    // Forward same endpoint to successor
    url := n.SuccessorAddr + "/storage/" + key
    req, _ := http.NewRequest(http.MethodPut, url, strings.NewReader(string(value)))
    req.Header.Set("Content-Type", "text/plain")
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        c.JSON(http.StatusBadGateway, gin.H{"error": "forward PUT failed", "detail": err.Error()})
        return
    }
    defer resp.Body.Close()
    io.Copy(io.Discard, resp.Body)
    c.Status(resp.StatusCode)
}


func GetValue(n *models.Node) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.Param("key")
        keyId := utils.ConsistentHash(key)

        if responsibleFor(keyId, n) {
            if v, ok := n.Store.Load(key); ok {
                c.String(http.StatusOK, "%s", v.(string))
                return
            }
            c.Status(http.StatusNotFound)
            return
        }
        forwardGet(n, key, c)
    }
}

func PutValue(n *models.Node) gin.HandlerFunc {
    return func(c *gin.Context) {
        key := c.Param("key")
        value, err := io.ReadAll(c.Request.Body)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "could not read body"})
            return
        }
        keyId := utils.ConsistentHash(key)

        if responsibleFor(keyId, n) {
            n.Store.Store(key, string(value))
            c.Status(http.StatusOK)
            return
        }
        forwardPut(n, key, value, c)
    }
}

func NetworkInfo(n *models.Node) gin.HandlerFunc {
    return func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "self": gin.H{
                "addr": n.Addr,
                "id":   n.NodeId,
            },
            "predecessor": gin.H{
                "addr": n.PredecessorAddr,
                "id":   n.PredecessorId,
            },
            "successor": gin.H{
                "addr": n.SuccessorAddr,
                "id":   n.SuccessorId,
            },
            "hashlen": utils.HASHLEN,
        })
    }
}
