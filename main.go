package main

import (
	"github.com/gin-gonic/gin"
)

var router *gin.Engine

func main() {
	initDB()
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"res": "pong"})
	})
	configPostsRoute(router)
	router.Run(":8080")
}
