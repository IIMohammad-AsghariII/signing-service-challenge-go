package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"regexp"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
)

func main() {

	// Regex to match the current working directory
	re := regexp.MustCompile(`^(.*` + "signing-service-challenge-go" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	// Load environment variables from .env file
	if err := godotenv.Load(string(rootPath) + `/.env`); err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get the DATA_STORE environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}
	listenAddress := ":" + port
	server := api.NewServer(listenAddress)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", listenAddress)
	}
}
