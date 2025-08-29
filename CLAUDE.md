# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based blog backend API using the Gin web framework with PostgreSQL database integration via GORM. The application serves as a REST API for a blog system.

## Architecture

- **Main entry point**: `main.go` - Sets up Gin router and starts server on port 8080
- **Database layer**: `conn.go` - Handles PostgreSQL connection, migrations, and GORM setup
- **API routes**: `posts.go` - Defines blog post endpoints and route configuration
- **Container support**: `dockerfile` - Multi-stage Docker build for production deployment

## Key Components

### Database Configuration
- Uses GORM with PostgreSQL driver
- Environment variables for database connection (DB_HOST, DB_USER, DB_PASS, DB_NAME, DB_PORT, DB_SSLMODE)
- Auto-migration for Post model
- Loads configuration from `.env` file or system environment

### API Structure
- Base route group `/posts` for blog-related endpoints
- Health check endpoint at `/ping`
- Home endpoint fetches README from external SCM (uses SCM_ORIGIN and SCM_REPO env vars)

### Models
- `Post` struct with UUID primary key using `gen_random_uuid()`

## Development Commands

```bash
# Run the application
go run .

# Build the application
go build -o main

# Install dependencies
go mod download

# Update dependencies
go mod tidy

# Build Docker image
docker build -t blog-backend .

# Run Docker container
docker run -p 8080:8080 blog-backend
```

## Environment Variables Required

- `DB_HOST` - PostgreSQL host
- `DB_USER` - Database username  
- `DB_PASS` - Database password
- `DB_NAME` - Database name
- `DB_PORT` - Database port
- `DB_SSLMODE` - SSL mode for database connection
- `SCM_ORIGIN` - Source control origin URL (for home endpoint)
- `SCM_REPO` - Repository path (for home endpoint)

## Dependencies

Key dependencies from go.mod:
- `github.com/gin-gonic/gin` - Web framework
- `gorm.io/gorm` and `gorm.io/driver/postgres` - ORM and PostgreSQL driver
- `github.com/google/uuid` - UUID generation
- `github.com/joho/godotenv` - Environment variable loading