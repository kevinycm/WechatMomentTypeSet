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
	// dbDSN := "root:Look-good[123]@tcp(101.200.219.216:3306)/sas?charset=utf8mb4&parseTime=True&loc=Local"
	dbDSN := ""
	server := backend.NewServer(8888, basePath, dbDSN)
	if err := server.Start(); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
