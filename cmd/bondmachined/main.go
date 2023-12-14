package main

import (
	"fmt"
	"net"
	"net/http"
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

	// Handle GET requests to /message
	mux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from the Go daemon!")
	})

	// Serve HTTP requests on the Unix socket
	http.Handle("/", mux)
	err = http.Serve(listener, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
