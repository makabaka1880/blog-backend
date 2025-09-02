package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"blog/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TreeResponse struct {
	SHA       string `json:"sha"`
	URL       string `json:"url"`
	Truncated bool   `json:"truncated"`
	Tree      []struct {
		Path string `json:"path"`
		Mode string `json:"mode"`
		Type string `json:"type"`
		SHA  string `json:"sha"`
		Size int    `json:"size"`
		URL  string `json:"url"`
	} `json:"tree"`
}

func GetTree(tree string) (*TreeResponse, error) {
	reqURL := fmt.Sprintf("https://api.github.com/repos/%s/git/trees/%s?recursive=1", os.Getenv("CMS_REPO"), tree)
	res, err := http.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer res.Body.Close()

	var parsedTree TreeResponse
	if err := json.NewDecoder(res.Body).Decode(&parsedTree); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &parsedTree, nil
}

func ProcessTreeAndSaveToDB(d *gorm.DB, tree *TreeResponse) error {
	// Delete all existing tree nodes first (THIS SHOULD ONLY BE USED IN THE UPDATE ALL ENDPOINT)
	if err := d.Where("1 = 1").Delete(&models.Tree{}).Error; err != nil {
		return fmt.Errorf("failed to delete existing tree nodes: %w", err)
	}
	parentMap := make(map[string]uuid.UUID)
	directoryNodes := make(map[string]uuid.UUID)

	for _, item := range tree.Tree {
		parts := strings.Split(item.Path, "/")
		var currentParentID *uuid.UUID
		currentPath := ""

		// Create directory structure first (only for files, directories are handled separately)
		if item.Type == "blob" { // Only files need directory structure creation
			for i, part := range parts[:len(parts)-1] {
				if i > 0 {
					currentPath += "/"
				}
				currentPath += part

				// Check if this path already exists in our map
				if existingID, exists := parentMap[currentPath]; exists {
					currentParentID = &existingID
					continue
				}

				treeNode := models.Tree{
					Name:     part,
					ParentID: currentParentID,
					SHA:      "",
					URL:      "",
				}

				// Create the tree node in database
				if err := d.Create(&treeNode).Error; err != nil {
					return fmt.Errorf("failed to create tree node: %w", err)
				}

				// Store the ID in our map for future reference
				parentMap[currentPath] = treeNode.ID
				directoryNodes[currentPath] = treeNode.ID
				currentParentID = &treeNode.ID
			}
		}

		// Handle the actual item (file or directory)
		if item.Type == "blob" { // File
			finalPart := parts[len(parts)-1]
			if len(parts) > 1 {
				currentPath += "/"
			}
			currentPath += finalPart

			// Check if this final path already exists
			if _, exists := parentMap[currentPath]; exists {
				continue
			}

			treeNode := models.Tree{
				Name:     finalPart,
				ParentID: currentParentID,
				SHA:      item.SHA,
				URL:      item.URL,
			}

			// Create the tree node in database
			if err := d.Create(&treeNode).Error; err != nil {
				return fmt.Errorf("failed to create tree node: %w", err)
			}

			// Store the ID in our map for future reference
			parentMap[currentPath] = treeNode.ID

			// If this is a markdown file, fetch and store its content
			if strings.HasSuffix(finalPart, ".md") {
				content, err := FetchContent(item.SHA)
				if err != nil {
					return fmt.Errorf("failed to fetch content for SHA %s: %w", item.SHA, err)
				}

				contentRecord := models.Content{
					Content: content,
					NodeID:  treeNode.ID,
				}
				if err := d.Create(&contentRecord).Error; err != nil {
					return fmt.Errorf("failed to create content record: %w", err)
				}
			}
		} else if item.Type == "tree" { // Directory
			// Handle directory items directly from GitHub response
			currentPath := item.Path

			// Check if this directory already exists
			if _, exists := parentMap[currentPath]; exists {
				continue
			}

			// Determine parent ID for this directory
			var parentID *uuid.UUID
			if len(parts) > 1 {
				parentPath := strings.Join(parts[:len(parts)-1], "/")
				if parentIDVal, exists := parentMap[parentPath]; exists {
					parentID = &parentIDVal
				}
			}

			treeNode := models.Tree{
				Name:     parts[len(parts)-1],
				ParentID: parentID,
				SHA:      "",
				URL:      "",
			}

			// Create the tree node in database
			if err := d.Create(&treeNode).Error; err != nil {
				return fmt.Errorf("failed to create tree node: %w", err)
			}

			// Store the ID in our map for future reference
			parentMap[currentPath] = treeNode.ID
			directoryNodes[currentPath] = treeNode.ID
		}
	}

	// Process all directories to generate TOC or use README/index files
	for path, nodeID := range directoryNodes {
		if err := ProcessDirectoryContent(d, nodeID, path); err != nil {
			return fmt.Errorf("failed to process directory content for %s: %w", path, err)
		}
	}

	return nil
}

func FetchContent(sha string) (string, error) {
	reqURL := fmt.Sprintf("https://api.github.com/repos/%s/git/blobs/%s", os.Getenv("CMS_REPO"), sha)
	res, err := http.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch content: %w", err)
	}
	defer res.Body.Close()

	var blobResponse struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(res.Body).Decode(&blobResponse); err != nil {
		return "", fmt.Errorf("failed to decode content: %w", err)
	}

	decodedContent, err := base64.StdEncoding.DecodeString(blobResponse.Content)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 content: %w", err)
	}

	return string(decodedContent), nil
}
