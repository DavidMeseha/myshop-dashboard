package main

import (
	"log"
	"net/http"
	"os"
	"shop-dashboard/internal/api"
	"shop-dashboard/internal/database"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client, err := database.Connect()
	if err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	database.SetMongoClient(client)

	router := api.NewRouter()

	// Create a CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{os.Getenv("ORIGIN")}, // Add your frontend URL
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(router)

	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
