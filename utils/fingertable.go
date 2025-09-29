package utils

import (
	"main/models"
	"math"
)

// finds closest node for given key
func FindClosestNodeByKey(key int, nodes []models.ClusterNodes) models.ClusterNodes {
	prevNode := nodes[0]

	for _, node := range nodes {
		if key <= node.Hash {
			return node
		}
		prevNode = node
	}

	return prevNode
}
 
func FingerTableInit(myNode *models.Node) {
	initialSplit := int(math.Pow(2, HASHLEN)) / 2;
	for i := 0; i < HASHLEN; i++ {
		nextNode := (myNode.NodeId + initialSplit) % RINGSIZE
		initialSplit /= 2

		closestNode := FindClosestNodeByKey(nextNode, myNode.Nodes)
		peer := CreatePeer(closestNode.Host, closestNode.Port)

		entry := models.FingerEntry{
			Key: nextNode,
			Node: peer,
		}

		myNode.FingerTable = append(myNode.FingerTable, entry)
	}
}

func FindSuccessorAddr(key int, myNode *models.Node) string {
	fingerTable := myNode.FingerTable
	if len(fingerTable) == 0 {
		return ""
	}

	for _, entry := range fingerTable {
		if key <= entry.Key {
			return entry.Node.Addr
		}
	}

	// wrap-around: successor is the first entry
	return fingerTable[0].Node.Addr
}
