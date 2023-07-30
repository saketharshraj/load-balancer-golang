package main

import (
	"fmt"
	"log"
	"net/http"
)


var (
	serverList = []*server{
		newServer("server-1", "http://127.0.0.1:5001"),
		newServer("server-2", "http://127.0.0.1:5002"),
		newServer("server-3", "http://127.0.0.1:5003"),
		newServer("server-4", "http://127.0.0.1:5004"),
		newServer("server-5", "http://127.0.0.1:5005"),
	}
	lastServedIndex = 0
)

func main() {
	http.HandleFunc("/", handleRequest)
	go startHealthCheck()
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handleRequest(res http.ResponseWriter, req *http.Request) {
	server, err := getHealthyServer()
	if err != nil {
		http.Error(res, "Couldn't process request: "+err.Error(), http.StatusServiceUnavailable)
		return
	}
	server.ReverseProxy.ServeHTTP(res, req)
	fmt.Printf("Served from %s server\n", server.Name)
}

func getHealthyServer() (*server, error) {
	for i := 0; i < len(serverList); i++ {
		server := getServer()
		if server.Health {
			return server, nil
		}
	}
	return nil, fmt.Errorf("all servers are down")
}

func getServer() *server {
	nextIndex := (lastServedIndex + 1) % len(serverList)
	server := serverList[nextIndex]
	lastServedIndex = nextIndex
	return server
}