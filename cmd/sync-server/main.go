package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"github.com/techyshishy/nirn-revenue-service/internal/sync"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Empty handler for debugging
	http.Handle("/", http.NotFoundHandler())

	syncServer := sync.Register()
	syncServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	err := http.ListenAndServe(net.JoinHostPort("::", port), nil)
	if err != nil {
		log.Print(err)
		return
	}
}
