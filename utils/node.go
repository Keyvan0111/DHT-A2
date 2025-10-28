package utils

import (
	"fmt"
	"main/models"
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
	if len(clusterNodes) == 0 {
		return
	}

	selfIndex := findNodeIndex(clusterNodes, myNode)
	if selfIndex == -1 {
		return
	}

	successorIndex := (selfIndex + 1) % len(clusterNodes)
	predecessorIndex := (selfIndex - 1 + len(clusterNodes)) % len(clusterNodes)

	successor := clusterNodes[successorIndex]
	predecessor := clusterNodes[predecessorIndex]

	successorAddr := BuildHTTPAddr(successor.Host, successor.Port)
	predecessorAddr := BuildHTTPAddr(predecessor.Host, predecessor.Port)

	sucHash, sucID := ConsistentHash(successor.Host + ":" + successor.Port)
	predHash, predID := ConsistentHash(predecessor.Host + ":" + predecessor.Port)

	sucessorNode := &models.Peer{
		Host:   successor.Host,
		Port:   successor.Port,
		Addr:   successorAddr,
		NodeId: sucID,
		Hash:   sucHash,
	}

	predecessorNode := &models.Peer{
		Host:   predecessor.Host,
		Port:   predecessor.Port,
		Addr:   predecessorAddr,
		NodeId: predID,
		Hash:   predHash,
	}

	myNode.Successor = *sucessorNode
	myNode.Predecessor = *predecessorNode

	fmt.Printf("ID: %d, Im Node: %s:%s, Hash: %s\n", myNode.NodeId, myNode.Host, myNode.Port, myNode.Hash)
	fmt.Printf("ID: %d, My Successor is: %s:%s Hash: %s\n", sucessorNode.NodeId, sucessorNode.Host, sucessorNode.Port, sucessorNode.Hash)
	fmt.Printf("ID: %d, My Predecessor is: %s:%s Hash: %s\n\n", predecessorNode.NodeId, predecessor.Host, predecessor.Port, predecessorNode.Hash)

}

func SortNodes(clusterNodes []models.ClusterNodes) {
	// compute hash for each node
	for i := range clusterNodes {
		_, hashNum := ConsistentHash(
			clusterNodes[i].Host + ":" + clusterNodes[i].Port,
		)
		clusterNodes[i].Hash = hashNum
	}

	// bubble sort by Hash
	n := len(clusterNodes)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if clusterNodes[j].Hash > clusterNodes[j+1].Hash {
				clusterNodes[j], clusterNodes[j+1] = clusterNodes[j+1], clusterNodes[j]
			}
		}
	}
}

func CreatePeer(host string, port string) models.Peer {
	hash, id := ConsistentHash(host + ":" + port)
	return models.Peer{
		Host:   host,
		Port:   port,
		Addr:   BuildHTTPAddr(host, port),
		Hash:   hash,
		NodeId: id,
	}
}

func ResetToSingleNode(myNode *models.Node) {
	selfEntry := models.ClusterNodes{
		Host: myNode.Host,
		Port: myNode.Port,
	}
	selfPeer := CreatePeer(myNode.Host, myNode.Port)

	myNode.Guard.Lock()
	defer myNode.Guard.Unlock()

	myNode.Nodes = []models.ClusterNodes{selfEntry}
	myNode.Successor = selfPeer
	myNode.Predecessor = selfPeer
	myNode.FingerTable = myNode.FingerTable[:0]
	myNode.State = models.NodeStateSingle
	myNode.LastKnownPeer = ""
}

func UpdateClusterView(myNode *models.Node, clusterNodes []models.ClusterNodes) {
	if len(clusterNodes) == 0 {
		ResetToSingleNode(myNode)
		return
	}

	filtered := EnsureUniqueNodes(clusterNodes)
	SortNodes(filtered)

	myNode.Guard.Lock()
	myNode.Nodes = CloneNodes(filtered)
	SetPeers(myNode, myNode.Nodes)
	FingerTableInit(myNode)
	if len(myNode.Nodes) == 1 {
		myNode.State = models.NodeStateSingle
		myNode.LastKnownPeer = ""
	} else {
		myNode.State = models.NodeStateActive
		myNode.LastKnownPeer = fmt.Sprintf("%s:%s", myNode.Successor.Host, myNode.Successor.Port)
	}
	myNode.Guard.Unlock()
}
