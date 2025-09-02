package main

import (
	"net/http"
	"time"

	"blog/internal/database"
	"blog/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	database.InitDB()
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

	handlers.ConfigPostsRoutes(router)
	handlers.ConfigTreesRoutes(router, database.DB)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			println("Server failed to start:", err.Error())
		}
	}()

	println("Server started on :8080")
	select {} // Keep main goroutine alive
}
