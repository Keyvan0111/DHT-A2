package models

import "sync"

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
	FingerTable []FingerEntry
}

type Peer struct {
	Host   string
	Port   string
	Addr   string
	Hash   string
	NodeId int
}
