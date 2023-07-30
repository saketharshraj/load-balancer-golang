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

func handleRequest(w http.ResponseWriter, r *http.Request) {
	server := serverList[lastServedIndex]
	server.ReverseProxy.ServeHTTP(w, r)
	fmt.Printf("Used server %d for handling request\n", lastServedIndex+1)
	lastServedIndex = (lastServedIndex + 1) % len(serverList);
	
}

func main() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8000", nil))
}