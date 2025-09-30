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
	selfIndex := findNodeIndex(clusterNodes, myNode)
	successorIndex := (selfIndex + 1) % len(clusterNodes)
	predecessorIndex := (selfIndex - 1 + len(clusterNodes)) % len(clusterNodes)

	successor := clusterNodes[successorIndex]
	predecessor := clusterNodes[predecessorIndex]

	successorAddr := fmt.Sprintf("http://%s.ifi.uit.no:%s", successor.Host, successor.Port)
	predecessorAddr := fmt.Sprintf("http://%s.ifi.uit.no:%s", predecessor.Host, predecessor.Port)

	sucHash, sucID := ConsistentHash(successor.Host+":"+successor.Port)
	predHash, predID := ConsistentHash(predecessor.Host+":"+predecessor.Port)
	
	sucessorNode := &models.Peer{
		Host: successor.Host,
		Port: successor.Port,
		Addr: successorAddr,
		NodeId: sucID,
		Hash: sucHash,
	}
	
	predecessorNode := &models.Peer{
		Host: predecessor.Host,
		Port: predecessor.Port,
		Addr: predecessorAddr,
		NodeId: predID,
		Hash: predHash,
	}
	
	myNode.Successor = *sucessorNode
	myNode.Predecessor = *predecessorNode

	fmt.Printf("ID: %d, Im Node: %s:%s, Hash: %s\n", myNode.NodeId ,myNode.Host, myNode.Port, myNode.Hash )
	fmt.Printf("ID: %d, My Successor is: %s:%s Hash: %s\n", sucessorNode.NodeId ,sucessorNode.Host, sucessorNode.Port, sucessorNode.Hash )
	fmt.Printf("ID: %d, My Predecessor is: %s:%s Hash: %s\n\n", predecessorNode.NodeId ,predecessor.Host, predecessor.Port, predecessorNode.Hash )

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
	hash, id := ConsistentHash(host+":"+port)
	return models.Peer{
		Host: host,
		Port: port,
		Addr: fmt.Sprintf("http://%s.ifi.uit.no:%s", host, port),
		Hash: hash,
		NodeId: id, 
	}
}
