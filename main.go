package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)


type Config struct {
	Servers []*server `json:"servers"`
}

var (
	serverList       []*server
	lastServedIndex = 0
)

func loadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)

	if err != nil {
		return nil, err
	}

	// Initialize ReverseProxy for each server in the config using the newServer function
	for _, server := range config.Servers {
		server.ReverseProxy = newServer(server.Name, server.URL).ReverseProxy
	}
	return &config, nil
}

func main() {
	config, err :=loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	serverList = config.Servers
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