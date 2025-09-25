package models

type ClusterNodes struct {
    Host string `json:"host"`
    Port string `json:"port"`
}

type Node struct {
	Host          string
	Port          int
	Addr		  string
	NodeId        int // This is the sum of the hash bytes
	SucessorId    int
	PredecessorId int
	Nodes         []*Node // Store pointer to all node stucts
}
