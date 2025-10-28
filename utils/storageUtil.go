package utils

import (
	"fmt"
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
	if node == nil {
		return false
	}
	node.Guard.RLock()
	start := node.Predecessor.NodeId
	end := node.NodeId
	node.Guard.RUnlock()

	return between(keyId, start, end)
}

func ForwardGet(node *models.Node, keyId int, key string, c *gin.Context, apiPath string) {
	// Linear implementation for now --> just hop to successor
	successor := FindPredecessorAddr(keyId, node)
	if successor == "" {
		node.Guard.RLock()
		successor = node.Successor.Addr
		node.Guard.RUnlock()
	}
	if successor == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no successor available"})
		return
	}
	fmt.Printf("Finding successor of key: %d", keyId)
	url := successor + apiPath + key

	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "forward GET failed", "detail": err.Error()})
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	c.Data(resp.StatusCode, "text/plain; charset=utf-8", body)
}

func ForwardPut(node *models.Node, keyId int, key string, value []byte, c *gin.Context, apiPath string) {
	// Forward same endpoint to successor
	successor := FindPredecessorAddr(keyId, node)
	if successor == "" {
		node.Guard.RLock()
		successor = node.Successor.Addr
		node.Guard.RUnlock()
	}
	if successor == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "no successor available"})
		return
	}
	fmt.Printf("Finding successor of key: %d", keyId)
	url := successor + apiPath + key

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
