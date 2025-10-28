package utils

import (
	"fmt"
	"main/models"
	"math"
	"sort"
)

// finds closest node for given key (nodes are sorted by Hash)
func FindClosestNodeByKey(key int, nodes []models.ClusterNodes) models.ClusterNodes {
	if len(nodes) == 0 {
		return models.ClusterNodes{}
	}

	sort.Slice(nodes, func(i, j int) bool { return nodes[i].Hash < nodes[j].Hash })

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
	if myNode == nil {
		return
	}

	myNode.FingerTable = myNode.FingerTable[:0] // reset in case it was already populated

	initialSplit := int(math.Pow(2, HASHLEN)) / 2
	for i := 0; i < HASHLEN; i++ {
		nextNode := (myNode.NodeId + initialSplit) % RINGSIZE
		initialSplit /= 2

		closestNode := FindClosestNodeByKey(nextNode, myNode.Nodes)
		peer := CreatePeer(closestNode.Host, closestNode.Port)

		entry := models.FingerEntry{
			Key:  nextNode,
			Node: peer,
		}

		myNode.FingerTable = append(myNode.FingerTable, entry)
	}

	sort.Slice(myNode.FingerTable, func(i, j int) bool { return myNode.FingerTable[i].Key < myNode.FingerTable[j].Key })
}

// Finger table must already be sorted by Key ascending on init.
func FindPredecessorAddr(key int, n *models.Node) string {
	if n == nil {
		return ""
	}

	n.Guard.RLock()
	table := make([]models.FingerEntry, len(n.FingerTable))
	copy(table, n.FingerTable)
	n.Guard.RUnlock()

	if len(table) == 0 {
		return ""
	}

	fmt.Println("Searching fingertable (predecessor)...")
	var (
		found bool
		addr  string
		k     int
	)

	for _, e := range table {
		if e.Key < key {
			found = true
			addr = e.Node.Addr
			k = e.Key
		} else {
			// since it's sorted, everything after is >= key; break early
			break
		}
	}

	if found {
		fmt.Printf("Found predecessor: %s (%d)\n", addr, k)
		return addr
	}

	// No entry with Key < key → wrap to the last entry
	last := table[len(table)-1]
	fmt.Printf("Wrapping to last entry: %s (%d)\n", last.Node.Host, last.Key)
	return last.Node.Addr
}
