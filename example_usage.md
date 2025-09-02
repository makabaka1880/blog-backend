# README.md/index.md Detection and TOC Generation System

This document explains how the new README.md/index.md detection and TOC generation system works.

## Overview

The system automatically processes GitHub repository structures and:
1. Detects README.md or index.md files in directories
2. Uses their content when available
3. Generates Table of Contents (TOC) when no README/index is found
4. Provides an extensible architecture for future enhancements like AI summaries

## How It Works

### 1. Tree Processing
When the `/tree/full-update` endpoint is called:
- Fetches the complete GitHub repository tree
- Processes each item (files and directories)
- Creates database entries for all nodes

### 2. Directory Content Detection
For each directory:
```go
hasReadme, readmeFile, err := HasReadmeOrIndex(db, nodeID)
if hasReadme {
    // Use README/index content
    content := fetchContent(readmeFile.SHA)
    saveContentToDB(content, nodeID)
} else {
    // Generate TOC
    tocManager.GenerateTOCForDirectory(nodeID, path)
}
```

### 3. TOC Generation
The default TOC generator creates markdown content with:
- 📁 Directories section with links
- 📄 Files section with proper formatting
- Clean organization and navigation

## Example Output

### Directory with README.md
If a directory contains `README.md`, that content is used directly.

### Directory without README/index
Generates a TOC like:
```markdown
# Table of Contents

## 📁 Directories

- [docs/](docs/)
- [images/](images/)

## 📄 Files

- [getting-started](getting-started.md)
- [configuration](configuration.md)
- LICENSE
```

## Extensibility

### Custom TOC Generators
Implement the `TOCGenerator` interface:
```go
type CustomTOCGenerator struct{}

func (g *CustomTOCGenerator) GenerateTOC(db *gorm.DB, nodeID uuid.UUID, path string) (string, error) {
    // Custom TOC generation logic
    return customTOCContent, nil
}
```

### AI Integration
Future enhancement example:
```go
type AITOCGenerator struct {
    aiClient *AIClient
}

func (g *AITOCGenerator) GenerateTOC(db *gorm.DB, nodeID uuid.UUID, path string) (string, error) {
    children := getChildrenFromDB(db, nodeID)
    summary := g.aiClient.SummarizeDirectory(children, path)
    return fmt.Sprintf("# AI-Generated Summary\n\n%s\n\n## Contents\n\n...", summary), nil
}
```

## API Endpoints

- `GET /tree/full-update` - Processes entire repository tree
- (Future) `GET /tree/update/:path` - Updates specific directory
- (Future) `GET /content/:path` - Retrieves content for path

## Database Schema

### Tree Table
- `ID` - UUID primary key
- `Name` - File/directory name
- `ParentID` - Reference to parent directory
- `SHA` - GitHub blob SHA (for files)
- `URL` - GitHub API URL

### Content Table
- `ID` - UUID primary key
- `Content` - Markdown content
- `NodeID` - Reference to Tree node

## Environment Variables

Ensure these are set:
- `SCM_REPO` - GitHub repository (e.g., "username/repo")
- Database connection variables

## Usage Example

```bash
curl -X GET "http://localhost:8080/tree/full-update"
```

This will process the entire repository structure and generate appropriate content for all directories.
