package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/BondMachineHQ/BMBuildkit/pkg/bondmachined"
)

func main() {
	socketPath := "/var/run/bm.sock"

	// Create the Unix socket
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	// Create a new HTTPServeMux
	mux := http.NewServeMux()

	// Handle GET requests to list available board
	mux.HandleFunc("/discovery", bondmachined.DiscoveryHandler)
	// Handle POST requests to pull artifact
	mux.HandleFunc("/pull", bondmachined.PullHandler)
	// Handle POST requests to pull artifact
	mux.HandleFunc("/load", bondmachined.LoadHandler)

	// Serve HTTP requests on the Unix socket
	http.Handle("/", mux)
	err = http.Serve(listener, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
