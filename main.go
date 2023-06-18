package main

import (
	"os"
	"log"
	"github.com/joho/godotenv"
	"github.com/the-go-dragons/fake-airline-info-service/server"
)

const (
	PortNumber = "PORT_FAKE"
)

func main() {
	loadConfig()
	port := getPort(PortNumber)
	server.StartServer(port)
}

func loadConfig() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
}

func getPort(name string) string {
	return os.Getenv(name)
}
