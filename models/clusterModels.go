package models

import (
	"sync"
)

type NodeState string

const (
	NodeStateSingle  NodeState = "single"
	NodeStateActive  NodeState = "active"
	NodeStateLeaving NodeState = "leaving"
	NodeStateCrashed NodeState = "crashed"
)

type ClusterNodes struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Hash int    `json:"-"`
}

type FingerEntry struct {
	Key  int
	Node Peer
}

type Node struct {
	Host        string
	Port        string
	Addr        string
	NodeId      int // This is the sum of the hash bytes
	Hash        string
	Successor   Peer
	Predecessor Peer
	Store       sync.Map
	Nodes 		[]ClusterNodes

	FingerTable []FingerEntry
	Guard         sync.RWMutex
	State         NodeState
	LastKnownPeer string
}

type Peer struct {
	Host   string
	Port   string
	Addr   string
	Hash   string
	NodeId int
}
