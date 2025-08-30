package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func configPostsRoutes(r *gin.Engine) {
	posts := r.Group("/posts")
	posts.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})
	posts.GET("/home", func(c *gin.Context) {
		uri := fmt.Sprintf("%s/%s/refs/heads/%s/README.md", os.Getenv("CMS_ORIGIN"), os.Getenv("CMS_REPO"), os.Getenv("CMS_BRANCH"))
		resp, err := http.Get(uri)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch README"})
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Data(http.StatusOK, "text/plain; charset=utf-8", body)
	})
}
