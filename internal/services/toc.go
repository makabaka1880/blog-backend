package services

import (
	"fmt"
	"strings"

	"blog/pkg/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TOCGenerator defines the interface for generating table of contents
type TOCGenerator interface {
	GenerateTOC(db *gorm.DB, nodeID uuid.UUID, path string) (string, error)
}

// DefaultTOCGenerator implements basic TOC generation
type DefaultTOCGenerator struct{}

// GenerateTOC creates a table of contents for a directory
func (g *DefaultTOCGenerator) GenerateTOC(db *gorm.DB, nodeID uuid.UUID, path string) (string, error) {
	var children []models.Tree
	if err := db.Where("parent_id = ?", nodeID).Find(&children).Error; err != nil {
		return "", fmt.Errorf("failed to fetch children: %w", err)
	}

	if len(children) == 0 {
		return "# Empty Directory\n\nThis directory doesn't contain any files or subdirectories.", nil
	}

	var tocBuilder strings.Builder
	tocBuilder.WriteString("# Table of Contents\n\n")

	// Group by type (file vs directory)
	var files []models.Tree
	var directories []models.Tree

	for _, child := range children {
		// Check if this child has children (indicating it's a directory)
		var childCount int64
		db.Model(&models.Tree{}).Where("parent_id = ?", child.ID).Count(&childCount)

		if childCount > 0 || child.SHA == "" {
			directories = append(directories, child)
		} else {
			files = append(files, child)
		}
	}

	// Add directories section
	if len(directories) > 0 {
		tocBuilder.WriteString("## 📁 Directories\n\n")
		for _, dir := range directories {
			tocBuilder.WriteString(fmt.Sprintf("- [%s/](%s/)\n", dir.Name, dir.Name))
		}
		tocBuilder.WriteString("\n")
	}

	// Add files section
	if len(files) > 0 {
		tocBuilder.WriteString("## 📄 Files\n\n")
		for _, file := range files {
			if strings.HasSuffix(file.Name, ".md") {
				// For markdown files, remove .md extension in display
				displayName := strings.TrimSuffix(file.Name, ".md")
				tocBuilder.WriteString(fmt.Sprintf("- [%s](%s)\n", displayName, file.Name))
			} else {
				tocBuilder.WriteString(fmt.Sprintf("- %s\n", file.Name))
			}
		}
	}

	return tocBuilder.String(), nil
}

// TOCManager handles TOC generation and management
type TOCManager struct {
	generator TOCGenerator
	db        *gorm.DB
}

// NewTOCManager creates a new TOC manager
func NewTOCManager(db *gorm.DB, generator TOCGenerator) *TOCManager {
	return &TOCManager{
		generator: generator,
		db:        db,
	}
}

// GenerateTOCForDirectory generates TOC for a specific directory
func (m *TOCManager) GenerateTOCForDirectory(nodeID uuid.UUID, path string) error {
	content, err := m.generator.GenerateTOC(m.db, nodeID, path)
	if err != nil {
		return fmt.Errorf("failed to generate TOC: %w", err)
	}

	// Check if content already exists for this node
	var existingContent models.Content
	if err := m.db.Where("node_id = ?", nodeID).First(&existingContent).Error; err == nil {
		// Update existing content
		existingContent.Content = content
		if err := m.db.Save(&existingContent).Error; err != nil {
			return fmt.Errorf("failed to update TOC content: %w", err)
		}
	} else {
		// Create new content
		contentRecord := models.Content{
			Content: content,
			NodeID:  nodeID,
		}
		if err := m.db.Create(&contentRecord).Error; err != nil {
			return fmt.Errorf("failed to create TOC content: %w", err)
		}
	}

	return nil
}

// HasReadmeOrIndex checks if a directory has README.md or index.md files
func HasReadmeOrIndex(db *gorm.DB, nodeID uuid.UUID) (bool, *models.Tree, error) {
	var readmeFiles []models.Tree
	if err := db.Where("parent_id = ? AND (name = 'README.md' OR name = 'index.md')", nodeID).Find(&readmeFiles).Error; err != nil {
		return false, nil, fmt.Errorf("failed to check for README/index files: %w", err)
	}

	if len(readmeFiles) > 0 {
		return true, &readmeFiles[0], nil
	}

	return false, nil, nil
}

// ProcessDirectoryContent handles content generation for directories
func ProcessDirectoryContent(db *gorm.DB, nodeID uuid.UUID, path string) error {
	hasReadme, readmeFile, err := HasReadmeOrIndex(db, nodeID)
	if err != nil {
		return fmt.Errorf("failed to check directory content: %w", err)
	}

	if hasReadme && readmeFile != nil {
		// Directory has README/index file, fetch its content
		content, err := FetchContent(readmeFile.SHA)
		if err != nil {
			return fmt.Errorf("failed to fetch README content: %w", err)
		}

		// Create or update content for the directory
		var existingContent models.Content
		if err := db.Where("node_id = ?", nodeID).First(&existingContent).Error; err == nil {
			existingContent.Content = content
			if err := db.Save(&existingContent).Error; err != nil {
				return fmt.Errorf("failed to update directory content: %w", err)
			}
		} else {
			contentRecord := models.Content{
				Content: content,
				NodeID:  nodeID,
			}
			if err := db.Create(&contentRecord).Error; err != nil {
				return fmt.Errorf("failed to create directory content: %w", err)
			}
		}
	} else {
		// Directory doesn't have README/index, generate TOC
		tocManager := NewTOCManager(db, &DefaultTOCGenerator{})
		if err := tocManager.GenerateTOCForDirectory(nodeID, path); err != nil {
			return fmt.Errorf("failed to generate TOC: %w", err)
		}
	}

	return nil
}
