package utils

import (
	"fmt"
	"main/models"

	"strconv"
	"strings"
)

func findNodeIndex(clusterNodes []models.ClusterNodes, target *models.Node) int {
	for i, node := range clusterNodes {
		if node.Host == target.Host && node.Port == target.Port {
			return i
		}
	}
	return -1 // not found
}

func SetPeers(myNode *models.Node, clusterNodes []models.ClusterNodes) {
	selfIndex := findNodeIndex(clusterNodes, myNode)
	successorIndex := (selfIndex + 1) % len(clusterNodes)
	predecessorIndex := (selfIndex - 1 + len(clusterNodes)) % len(clusterNodes)

	successor := clusterNodes[successorIndex]
	predecessor := clusterNodes[predecessorIndex]

	successorAddr := fmt.Sprintf("http://%s.ifi.uit.no:%s", successor.Host, successor.Port)
	predecessorAddr := fmt.Sprintf("http://%s.ifi.uit.no:%s", predecessor.Host, predecessor.Port)

	myNode.SuccessorId   = ConsistentHash(fmt.Sprintf("%s:%s", successor.Host, successor.Port))
	myNode.PredecessorId = ConsistentHash(fmt.Sprintf("%s:%s", predecessor.Host, predecessor.Port))

	myNode.SuccessorAddr = successorAddr
	myNode.PredecessorAddr = predecessorAddr
}

func SortNodes(clusterNodes []models.ClusterNodes) {
	getParts := func(name string) (int, int) {
		parts := strings.Split(strings.TrimPrefix(name, "c"), "-")
		if len(parts) != 2 {
			return 0, 0
		}
		a, _ := strconv.Atoi(parts[0])
		b, _ := strconv.Atoi(parts[1])
		return a, b
	}

	n := len(clusterNodes)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			a1, a2 := getParts(clusterNodes[j].Host)
			b1, b2 := getParts(clusterNodes[j+1].Host)

			// compare
			if a1 > b1 || (a1 == b1 && a2 > b2) {
				// swap
				clusterNodes[j], clusterNodes[j+1] = clusterNodes[j+1], clusterNodes[j]
			}
		}
	}
	}

