package main

import (
	"log"
	"portfolio-api/cmd/api"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	server := api.NewAPIServer(":8080", nil)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
