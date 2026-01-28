package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/parthsarthi-dutt/blog-api/internal/config"
	"github.com/parthsarthi-dutt/blog-api/internal/routes"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: no .env file found")
	}

	config.ConnectDB()

	r := gin.Default()

	routes.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port", port)
	log.Fatal(r.Run(":" + port))
}
