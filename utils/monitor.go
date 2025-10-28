package utils

import (
	"fmt"
	"io"
	"log"
	"main/models"
	"net/http"
	"time"
)

const (
	monitorInterval    = 2 * time.Second
	monitorHTTPTimeout = time.Second
	requiredFailures   = 2
)

func StartHealthMonitor(node *models.Node) {
	go monitorSuccessor(node)
}

func monitorSuccessor(node *models.Node) {
	ticker := time.NewTicker(monitorInterval)
	defer ticker.Stop()

	client := &http.Client{Timeout: monitorHTTPTimeout}
	failures := make(map[string]int)

	for range ticker.C {
		node.Guard.RLock()
		if node.State == models.NodeStateCrashed || len(node.Nodes) <= 1 {
			node.Guard.RUnlock()
			continue
		}
		successor := node.Successor
		selfHost := node.Host
		selfPort := node.Port
		node.Guard.RUnlock()

		key := fmt.Sprintf("%s:%s", successor.Host, successor.Port)
		if key == ":" || successor.Addr == "" || (successor.Host == selfHost && successor.Port == selfPort) {
			continue
		}

		if probeNode(client, successor.Addr) {
			failures[key] = 0
			continue
		}

		failures[key]++
		if failures[key] < requiredFailures {
			continue
		}
		failures[key] = 0

		log.Printf("successor %s deemed unreachable, updating cluster", key)
		handleNodeFailure(node, successor.Host, successor.Port)
	}
}

func probeNode(client *http.Client, baseAddr string) bool {
	req, err := http.NewRequest(http.MethodGet, baseAddr+"/node-info", nil)
	if err != nil {
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	return resp.StatusCode == http.StatusOK
}

func handleNodeFailure(node *models.Node, host, port string) {
	node.Guard.RLock()
	current := CloneNodes(node.Nodes)
	node.Guard.RUnlock()

	updated := RemoveNode(current, host, port)
	if len(updated) == len(current) {
		return
	}

	UpdateClusterView(node, updated)
	if err := BroadcastCluster(updated); err != nil {
		log.Printf("failed to broadcast updated cluster after failure: %v", err)
	}
}
