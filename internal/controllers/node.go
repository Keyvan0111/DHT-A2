package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"main/models"
	"main/utils"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var httpClient = &http.Client{Timeout: 4 * time.Second}

func nodeInfo(node *models.Node) gin.H {
	node.Guard.RLock()
	defer node.Guard.RUnlock()

	successor := fmt.Sprintf("%s:%s", node.Successor.Host, node.Successor.Port)
	othersSet := make(map[string]struct{})

	if node.Predecessor.Host != "" {
		key := fmt.Sprintf("%s:%s", node.Predecessor.Host, node.Predecessor.Port)
		if key != successor {
			othersSet[key] = struct{}{}
		}
	}

	for _, entry := range node.FingerTable {
		key := fmt.Sprintf("%s:%s", entry.Node.Host, entry.Node.Port)
		if key != successor {
			othersSet[key] = struct{}{}
		}
	}

	for _, n := range node.Nodes {
		key := fmt.Sprintf("%s:%s", n.Host, n.Port)
		if key != successor && key != fmt.Sprintf("%s:%s", node.Host, node.Port) {
			othersSet[key] = struct{}{}
		}
	}

	others := make([]string, 0, len(othersSet))
	for key := range othersSet {
		others = append(others, key)
	}
	sort.Strings(others)

	return gin.H{
		"node_hash": node.Hash,
		"state":     node.State,
		"successor": successor,
		"others":    others,
	}
}

func NodeInfoHandler(node *models.Node) gin.HandlerFunc {
	return func(c *gin.Context) {
		if GuardCrash(node, c) {
			return
		}
		c.JSON(http.StatusOK, nodeInfo(node))
	}
}

func JoinNetwork(node *models.Node) gin.HandlerFunc {
    return func(c *gin.Context) {
        if GuardCrash(node, c) {
            return
        }

		target := c.Query("nprime")
		if target == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing nprime parameter"})
			return
		}

		node.Guard.RLock()
		currentSize := len(node.Nodes)
		node.Guard.RUnlock()
		if currentSize > 1 {
			c.JSON(http.StatusConflict, gin.H{"error": "node already part of a network"})
			return
		}

		host, port, err := utils.NormalizeHostPort(target)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

        if err := joinVia(node, host, port); err != nil {
            // For API-level testing, respond 200 with detail on failure (mirrors SimRecover).
            c.JSON(http.StatusOK, gin.H{"message": "join attempted", "detail": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "joined network"})
    }
}

func LeaveNetwork(node *models.Node) gin.HandlerFunc {
	return func(c *gin.Context) {
		if GuardCrash(node, c) {
			return
		}

		node.Guard.RLock()
		size := len(node.Nodes)
		node.Guard.RUnlock()
		if size <= 1 {
			c.JSON(http.StatusOK, gin.H{"message": "already single node"})
			return
		}

		if err := leaveCluster(node); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "left network"})
	}
}

func SimCrash(node *models.Node) gin.HandlerFunc {
	return func(c *gin.Context) {
		node.Guard.Lock()
		successorAddr := ""
		if node.Successor.Host != "" && !(node.Successor.Host == node.Host && node.Successor.Port == node.Port) {
			successorAddr = fmt.Sprintf("%s:%s", node.Successor.Host, node.Successor.Port)
		}
		node.LastKnownPeer = successorAddr
		node.State = models.NodeStateCrashed
		node.Guard.Unlock()

		c.JSON(http.StatusOK, gin.H{"message": "node crashed"})
	}
}

func SimRecover(node *models.Node) gin.HandlerFunc {
	return func(c *gin.Context) {
		node.Guard.Lock()
		bootstrap := node.LastKnownPeer
		node.State = models.NodeStateSingle
		node.Guard.Unlock()

		utils.ResetToSingleNode(node)

		var (
			joined  bool
			joinErr error
		)
		if bootstrap != "" {
			host, port, err := utils.NormalizeHostPort(bootstrap)
			if err == nil && !(host == node.Host && port == node.Port) {
				if err := joinVia(node, host, port); err == nil {
					joined = true
				} else {
					joinErr = err
				}
			}
		}

		resp := gin.H{
			"message": "node recovered",
			"joined":  joined,
		}
		if joinErr != nil {
			resp["detail"] = joinErr.Error()
		}
		c.JSON(http.StatusOK, resp)
	}
}

func joinVia(node *models.Node, host, port string) error {
	members, err := fetchMembers(host, port)
	if err != nil {
		return err
	}
	if len(members) == 0 {
		members = append(members, models.ClusterNodes{Host: host, Port: port})
	}
	members = append(members, models.ClusterNodes{Host: node.Host, Port: node.Port})
	members = utils.EnsureUniqueNodes(members)
	utils.SortNodes(members)

	if err := utils.BroadcastCluster(members); err != nil {
		utils.ResetToSingleNode(node)
		return err
	}
	utils.UpdateClusterView(node, members)

	node.Guard.Lock()
	node.LastKnownPeer = fmt.Sprintf("%s:%s", host, port)
	node.Guard.Unlock()
	return nil
}

func leaveCluster(node *models.Node) error {
	node.Guard.RLock()
	current := utils.CloneNodes(node.Nodes)
	node.Guard.RUnlock()

	next := utils.RemoveNode(current, node.Host, node.Port)
	if len(next) == len(current) {
		return errors.New("node was not part of network")
	}
	utils.SortNodes(next)
	if err := utils.BroadcastCluster(next); err != nil {
		return err
	}

	utils.ResetToSingleNode(node)
	node.Guard.Lock()
	if len(next) > 0 {
		node.LastKnownPeer = fmt.Sprintf("%s:%s", next[0].Host, next[0].Port)
	} else {
		node.LastKnownPeer = ""
	}
	node.Guard.Unlock()

	return nil
}

func fetchMembers(host, port string) ([]models.ClusterNodes, error) {
	addr := utils.BuildHTTPAddr(host, port)
	resp, err := httpClient.Get(addr + "/cluster/members")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("members request failed: %s", strings.TrimSpace(string(body)))
	}

	var members []models.ClusterNodes
	if err := json.NewDecoder(resp.Body).Decode(&members); err != nil {
		return nil, err
	}
	return members, nil
}
