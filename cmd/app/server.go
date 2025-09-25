package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	"main/models"
	"main/routes"
	"main/utils"
)


func shortHost() string {
	h, err := os.Hostname()
	if err != nil || h == "" {
		return "unknown"
	}
	if i := strings.IndexByte(h, '.'); i > 0 {
		return h[:i]
	}
	return h
}

func main() {
	// Pick a free port on IPv4
	ln, err := net.Listen("tcp4", "0.0.0.0:0")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	host := shortHost()
	fullHost, _ := os.Hostname()

	NodeId := utils.ConsistentHash(fmt.Sprintf("%s:%d", host, port))

	myNode := &models.Node{
			Host: host,
			Port: fmt.Sprintf("%d", port),
			Addr: fmt.Sprintf("http://%s.ifi.uit.no:%d", host, port),
			NodeId: NodeId,
			SuccessorAddr: "",
			PredecessorAddr: "",
	}

	// If PORT_FILE is set, write the chosen port there so run.sh can read it
	if path := os.Getenv("PORT_FILE"); path != "" {
		f, err := os.Create(path)
		if err != nil {
			log.Fatalf("failed to create PORT_FILE %q: %v", path, err)
		}
		if _, err := fmt.Fprintf(f, "%d", port); err != nil {
			_ = f.Close()
			log.Fatalf("failed to write PORT_FILE %q: %v", path, err)
		}
		_ = f.Sync()
		_ = f.Close()
	}

	// Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	routes.SetupClusterRoutes(router, myNode)

	// Routes
	router.GET("/helloworld", func(c *gin.Context) {
		log.Println("Hello server guys!")
		c.String(200, "%s:%d", host, port)
	})

	fmt.Println("Hello guys im here!")
	log.Printf("listening on %s:%d\n", fullHost, port)

	// Use Gin's RunListener to serve on the already-open socket
	if err := router.RunListener(ln); err != nil {
		log.Fatalf("http serve error: %v", err)
	}
}
