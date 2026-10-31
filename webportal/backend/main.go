package main

import (
	"fmt"
	"os"
	"time"

	"pemkab-tulungagung/config"
	"pemkab-tulungagung/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	r := gin.Default()

	r.Static("/uploads", "./public/uploads")

	// Setup CORS agar Nuxt (port 3000) bisa akses API (port 8080)
	r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}))

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "success", "message": "Backend Golang Berjalan!"})
	})

	// Setup Routes
	routes.SetupRouter(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("🚀 Server Golang berjalan di port %s\n", port)
	r.Run(":" + port)
}