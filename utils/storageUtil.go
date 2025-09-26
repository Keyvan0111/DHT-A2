package utils

import (
	"io"
	"main/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func between(id, start, end int) bool {
    if start < end {
        return id > start && id <= end
    }
    return id > start || id <= end
}

func IsResponsibleFor(keyId int, node *models.Node) bool {

    return between(keyId, node.Predecessor.NodeId, node.NodeId)
}

func ForwardGet(n *models.Node, key string, c *gin.Context) {
    // Linear implementation for now --> just hop to successor
    url := n.Successor.Addr + "/storage/" + key
    resp, err := http.Get(url)
    if err != nil {
        c.JSON(http.StatusBadGateway, gin.H{"error": "forward GET failed", "detail": err.Error()})
        return
    }
    defer resp.Body.Close()
    body, _ := io.ReadAll(resp.Body)
    c.Data(resp.StatusCode, "text/plain; charset=utf-8", body)
}

func ForwardPut(n *models.Node, key string, value []byte, c *gin.Context) {
    // Forward same endpoint to successor
    url := n.Successor.Addr + "/storage/" + key
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
