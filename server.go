package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type server struct {
	Name         string
	URL          string
	ReverseProxy *httputil.ReverseProxy
	Health       bool
}

func newServer(name, serverUrl string) *server {
	parsedUrl, _ := url.Parse(serverUrl)
	rp := httputil.NewSingleHostReverseProxy(parsedUrl)
	return &server{
		Name:         name,
		URL:          serverUrl,
		ReverseProxy: rp,
		Health:       true,
	}
}

func (s *server) checkHealth() bool {
	resp, err := http.Head(s.URL)
	if err != nil {
		s.Health = false
		return s.Health
	}
	if resp.StatusCode != http.StatusOK {
		s.Health = false
		return s.Health
	}
	s.Health = true
	return s.Health
}