# README.md/index.md Detection and TOC Generation System - Implementation Guide

## Overview

This implementation adds intelligent README.md/index.md detection and Table of Contents (TOC) generation to the blog backend. The system automatically processes GitHub repository structures and provides appropriate content for directories.

## Key Features

1. **Automatic README/index Detection**: Scans directories for README.md or index.md files
2. **Smart Content Selection**: Uses README/index content when available, generates TOC otherwise
3. **Extensible Architecture**: Modular design for future enhancements like AI summaries
4. **GitHub Integration**: Works with GitHub API to fetch repository structures

## Architecture

### Database Schema

#### Tree Table
- `ID` (uuid) - Primary key
- `Name` (varchar) - File/directory name  
- `ParentID` (uuid) - Reference to parent directory (nullable)
- `SHA` (varchar) - GitHub blob SHA (for files)
- `URL` (varchar) - GitHub API URL

#### Content Table
- `ID` (uuid) - Primary key
- `Content` (text) - Markdown content
- `NodeID` (uuid) - Reference to Tree node

### Core Components

#### 1. Tree Processing (`treemanip.go`)
- **Function**: `processTreeAndSaveToDB()`
- **Purpose**: Processes GitHub tree response and creates database entries
- **Key Logic**:
  - Handles both files (`type: "blob"`) and directories (`type: "tree"`)
  - Builds hierarchical structure with proper parent-child relationships
  - Tracks directories for content processing

#### 2. TOC Generation (`toc.go`)
- **Interface**: `TOCGenerator`
- **Default Implementation**: `DefaultTOCGenerator`
- **Features**:
  - Groups items by type (directories vs files)
  - Creates markdown links for navigation
  - Handles empty directories gracefully

#### 3. Content Management (`toc.go`)
- **Function**: `ProcessDirectoryContent()`
- **Logic Flow**:
  1. Check if directory has README.md or index.md via `HasReadmeOrIndex()`
  2. If README exists: fetch content using `fetchContent()` and save to database
  3. If no README: generate TOC using `TOCManager` and save to database

## API Endpoints

### `GET /tree/full-update`
- **Purpose**: Processes entire repository structure
- **Flow**:
  1. Fetches latest commit SHA from GitHub
  2. Retrieves complete repository tree
  3. Processes all files and directories
  4. Generates appropriate content for each directory

## Implementation Details

### GitHub Integration

The system uses GitHub's Git Data API:
- Tree endpoint: `/repos/{owner}/{repo}/git/trees/{sha}?recursive=1`
- Blob endpoint: `/repos/{owner}/{repo}/git/blobs/{sha}`

### Content Processing Logic

```go
// For each directory in directoryNodes map
if directory has README.md or index.md {
    fetch README content from GitHub
    save content to database with directory's NodeID
} else {
    generate TOC using DefaultTOCGenerator
    save TOC content to database
}
```

### TOC Generation Algorithm

1. **Fetch Children**: Get all items in the directory
2. **Categorize**: Separate into directories and files
3. **Generate Markdown**:
   - Directories section with 📁 emoji and trailing slashes
   - Files section with 📄 emoji
   - Markdown files get clean display names (without .md extension)
   - Non-markdown files shown as plain text

## Extensibility Patterns

### Custom TOC Generators

Implement the `TOCGenerator` interface:

```go
type CustomTOCGenerator struct{}

func (g *CustomTOCGenerator) GenerateTOC(db *gorm.DB, nodeID uuid.UUID, path string) (string, error) {
    // Custom generation logic
    return customContent, nil
}
```

### AI Integration Example

```go
type AITOCGenerator struct {
    aiClient *AIClient
}

func (g *AITOCGenerator) GenerateTOC(db *gorm.DB, nodeID uuid.UUID, path string) (string, error) {
    children := getChildrenFromDB(db, nodeID)
    summary := g.aiClient.SummarizeDirectory(children, path)
    return fmt.Sprintf("# AI Summary\n\n%s\n\n## Contents\n\n...", summary), nil
}
```

## Environment Variables

Required for GitHub integration:
- `SCM_REPO`: GitHub repository (format: "owner/repo")
- `SCM_ORIGIN`: GitHub API base URL (optional, defaults to GitHub)

Database configuration:
- `DB_HOST`, `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_PORT`, `DB_SSLMODE`

## Usage Example

```bash
# Trigger full repository processing
curl -X GET "http://localhost:8080/tree/full-update"
```

## Expected Output

### Directory with README.md
Uses the README content directly.

### Directory without README/index
Generates TOC like:
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

### Empty Directory
```markdown
# Empty Directory

This directory doesn't contain any files or subdirectories.
```

## Error Handling

- GitHub API failures are properly handled with error messages
- Database operations include proper error checking
- Content fetching includes fallback mechanisms

## Performance Considerations

- Batch processing of directories after tree creation
- Efficient database queries with proper indexing
- Memory-efficient content handling

## Future Enhancements

1. **Incremental Updates**: Process only changed directories
2. **Content Caching**: Cache GitHub API responses
3. **Custom Templates**: Allow custom TOC templates
4. **AI Integration**: Add AI-powered summaries and content generation
5. **Web Interface**: Admin panel for manual content management

## Testing

The system includes test examples in `test_toc.go` demonstrating:
- TOC generation with various directory structures
- README detection logic
- Empty directory handling
- Database integration

This implementation provides a robust foundation for intelligent content management in GitHub-based blog systems.