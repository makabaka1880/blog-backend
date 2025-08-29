package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	initDB()
	router = gin.Default() // fixed

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",         // dev
			"http://localhost:8080",         // dev (backend UI)
			"http://blog.makabaka1880.xyz",  // prod http
			"https://blog.makabaka1880.xyz", // prod https
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"res": "pong"})
	})

	configPostsRoute(router)
	router.Run(":8080")
}
