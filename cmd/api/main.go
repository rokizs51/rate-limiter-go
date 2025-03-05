package main

import (
	"fmt"
	"log"
	"rateLimiter/internal/config"
	"rateLimiter/internal/database"
	"rateLimiter/internal/factory"
	handlers "rateLimiter/internal/handler"
	"rateLimiter/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()
	database.Initialize(cfg)
	database.InitializeRedis(cfg)
	fmt.Println(cfg)
	gin.SetMode(gin.DebugMode)
	r := gin.Default()

	// r.Use(middleware.Timeout(time.Second * 5))
	r.Use(middleware.RateLimiter(cfg, factory.TokenBucket))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/rate-limits", handlers.GetRateLimiterData())

	log.Println("Starting server on port 8080")

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
