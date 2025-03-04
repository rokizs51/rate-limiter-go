package main

import (
	"log"
	"rateLimiter/internal/config"
	"rateLimiter/internal/database"
	"rateLimiter/internal/factory"
	"rateLimiter/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()
	database.Initialize(cfg)

	r := gin.Default()

	r.Use(middleware.RateLimiter(cfg, factory.TokenBucket))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	log.Println("Starting server on port 8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
