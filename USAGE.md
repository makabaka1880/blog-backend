# README.md/index.md Detection and TOC Generation - Usage Guide

## Overview

This system automatically processes GitHub repository structures to provide intelligent content for directories. It detects README.md/index.md files and generates Table of Contents when needed.

## Quick Start

### 1. Environment Setup

Create a `.env` file with required variables:

```bash
# GitHub Repository Configuration
SCM_REPO=your-username/your-repository
SCM_ORIGIN=https://github.com

# Database Configuration
DB_HOST=localhost
DB_USER=your_username
DB_PASS=your_password
DB_NAME=blog_db
DB_PORT=5432
DB_SSLMODE=disable
```

### 2. Start the Application

```bash
# Build and run
go build -o main .
./main

# Or run directly
go run .
```

### 3. Trigger Repository Processing

```bash
# Process entire repository structure
curl -X GET "http://localhost:8080/tree/full-update"
```

## API Endpoints

### Process Repository Tree

**Endpoint**: `GET /tree/full-update`

**Description**: Processes the entire GitHub repository structure, creates database entries, and generates appropriate content for all directories.

**Response**:
```json
{
  "status": "ok",
  "message": "Tree updated successfully"
}
```

### Health Check

**Endpoint**: `GET /ping`

**Description**: Simple health check endpoint.

**Response**:
```json
{
  "res": "pong"
}
```

### Home Page Content

**Endpoint**: `GET /posts/home`

**Description**: Fetches the repository's README.md file.

## How It Works

### Automatic Content Detection

The system automatically:

1. **Scans directories** for README.md or index.md files
2. **Uses existing content** when README/index files are found
3. **Generates TOC** when no README/index files are available
4. **Handles empty directories** gracefully

### Example Directory Structures

#### Directory with README.md
```
project/
├── README.md          # Content used directly
├── src/
│   └── main.go
└── docs/
    └── guide.md
```

#### Directory without README/index
```
utils/
├── helpers.go
├── validators.go
└── formatters.go
# Generates TOC automatically
```

### Generated TOC Example

For a directory without README.md:

```markdown
# Table of Contents

## 📁 Directories

- [src/](src/)
- [docs/](docs/)

## 📄 Files

- [configuration](configuration.md)
- [getting-started](getting-started.md)
- LICENSE
```

## Database Schema

### Tree Table
Stores the repository structure:
- `ID` - UUID primary key
- `Name` - File/directory name
- `ParentID` - Reference to parent directory
- `SHA` - GitHub blob SHA (files only)
- `URL` - GitHub API URL

### Content Table
Stores generated content:
- `ID` - UUID primary key  
- `Content` - Markdown content
- `NodeID` - Reference to Tree node

## Customization

### Custom TOC Generators

Implement the `TOCGenerator` interface:

```go
type CustomTOCGenerator struct{}

func (g *CustomTOCGenerator) GenerateTOC(db *gorm.DB, nodeID uuid.UUID, path string) (string, error) {
    // Your custom generation logic
    return customContent, nil
}
```

### Using Custom Generator

```go
// Replace the default generator
tocManager := NewTOCManager(db, &CustomTOCGenerator{})
```

## Error Handling

The system includes comprehensive error handling for:
- GitHub API failures
- Database connection issues
- Invalid repository configurations
- Network timeouts

## Monitoring

Check application logs for:
- Tree processing status
- Content generation results
- Error messages and warnings

## Troubleshooting

### Common Issues

1. **GitHub API Rate Limits**: Ensure proper authentication if needed
2. **Database Connection**: Verify database credentials and connectivity
3. **Repository Access**: Confirm the repository exists and is accessible

### Debug Mode

Add debug logging by modifying the application to output more detailed information during processing.

## Performance Considerations

- Processing large repositories may take time
- Consider incremental updates for frequent changes
- Database indexing is optimized for tree traversal

## Future Enhancements

Planned features:
- Incremental directory updates
- Content caching
- AI-powered summaries
- Web-based administration interface

## Support

For issues or questions:
1. Check application logs
2. Verify environment variables
3. Ensure repository accessibility
4. Review database connectivity

This system provides a robust foundation for intelligent content management in GitHub-based applications.