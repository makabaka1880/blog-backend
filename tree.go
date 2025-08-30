package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func configTreesRoutes(r *gin.Engine, d *gorm.DB) {
	trees := r.Group("tree")
	trees.POST("/full-update", handleFullUpdate(d))
}

func handleFullUpdate(d *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		// Extract token part
		token := strings.TrimPrefix(authHeader, "Bearer ")

		// Look up token in the database
		var key AuthKey
		if err := d.
			Where("key = ?", "auth").
			Where("val = ?", token).
			First(&key).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}
		commitSHA, err := getLatestCommitSHA()
		if err != nil {
			fmt.Printf("%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get latest commit"})
			return
		}

		tree, err := getTree(commitSHA)
		if err != nil || tree == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tree"})
			return
		}
		if err := processTreeAndSaveToDB(d, tree); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process tree"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "message": "Tree updated successfully"})
	}
}

func getLatestCommitSHA() (string, error) {
	reqURL := fmt.Sprintf("https://api.github.com/repos/%s/commits/main", os.Getenv("CMS_REPO"))
	res, err := http.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch commit: %w", err)
	}
	defer res.Body.Close()

	var commitResponse struct {
		Commit struct {
			Tree struct {
				SHA string `json:"sha"`
			} `json:"tree"`
		} `json:"commit"`
	}

	// bodyBytes, _ := io.ReadAll(res.Body)
	// res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	// fmt.Printf("Response body: %s\n", string(bodyBytes))
	if err := json.NewDecoder(res.Body).Decode(&commitResponse); err != nil {
		return "", fmt.Errorf("failed to parse commit response: %w", err)
	}
	// fmt.Printf("%v", commitResponse)

	return commitResponse.Commit.Tree.SHA, nil
}
