package main

import (
	"github.com/the-go-dragons/fake-airline-info-service/config"
	"github.com/the-go-dragons/fake-airline-info-service/server"
)

func main() {
	server.StartServer(config.Get(config.HostPort))
}
