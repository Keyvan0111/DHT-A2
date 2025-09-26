package models

type ClusterNodes struct {
    Host string `json:"host"`
    Port string `json:"port"`
}

type Node struct {
	Host          	string
	Port          	string
	Addr		  	string
	NodeId        	int // This is the sum of the hash bytes
	Successor   	Peer
	Predecessor 	Peer
}

type Peer struct {
	Host          	string
	Port          	string
	Addr		  	string
	NodeId        	int 
}
