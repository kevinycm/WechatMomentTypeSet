package main

import (
	"log"
	"os"

	"wechatmomenttypeset/backend"
)

func main() {
	// Get current working directory
	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting current directory:", err)
	}

	// Create and start server
	server := backend.NewServer(8888, basePath)
	if err := server.Start(); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
