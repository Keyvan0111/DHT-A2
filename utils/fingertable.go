package utils

import (
	"main/models"
	"math"
)

// (start, end] modulo ring
func BetweenRightIncl(id, start, end int) bool {
	if start < end {
		return id > start && id <= end
	}
	return id > start || id <= end
}

// (start, end) modulo ring
func BetweenOpen(id, start, end int) bool {
	if start < end {
		return id > start && id < end
	}
	return id > start || id < end
}


// 15
func FindSuccessor(finger int, mynode *models.Node) (node *models.Node) {
	if IsResponsibleFor(finger, mynode) { // check if self responsible
		return mynode
	} else if finger > mynode.NodeId && finger <= mynode.Predecessor.NodeId { // successor responsible
		s := mynode.Successor
		return &models.Node{
			Host: s.Host,
			Port: s.Port,
			Addr: s.Addr,
			NodeId: s.NodeId,
			Hash: s.Hash,
		}
	} else {
		// response := http("GET", "/forward/1234")
		// return response
	}
	return nil
}


func FingerTableInit(node *models.Node) {
	for i := range HASHLEN {	
	 	key := (node.NodeId + (1 << i)) % int(math.Pow(2, HASHLEN))

		successor := FindSuccessor(key, node) // findsuccessor

		node.FingerTable[i] = models.FingerEntry{
			Key: key,
			Node: models.Peer{
                Host:   successor.Host,
                Port:   successor.Port,
                Addr:   successor.Addr,
                Hash:   successor.Hash,
                NodeId: successor.NodeId,
            },
		}
	}
}


// return best finger strictly in (n, id); else successor
func ClosestPrecedingFinger(n *models.Node, id int) models.Peer {
	for i := len(n.FingerTable) - 1; i >= 0; i-- {
		p := n.FingerTable[i].Node
		if BetweenOpen(p.NodeId, n.NodeId, id) {
			return p
		}
	}
	return n.Successor
}
