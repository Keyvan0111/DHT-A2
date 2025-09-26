package models

import "sync"

type ClusterNodes struct {
    Host string `json:"host"`
    Port string `json:"port"`
}

type Node struct {
	Host          string
	Port          string
	Addr		  string
	NodeId        int // This is the sum of the hash bytes
	SuccessorAddr    string
	PredecessorAddr string
	SuccessorId     int
  PredecessorId   int
	Nodes         []*Node // Store pointer to all node stucts
	Store sync.Map
}
