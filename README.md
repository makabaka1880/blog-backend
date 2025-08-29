# Blog Backend API

A lightweight Go REST API for a blog system using Gin web framework and PostgreSQL.

## Features

- RESTful API endpoints for blog posts
- PostgreSQL database integration with GORM
- Docker support for containerized deployment
- Environment-based configuration
- Health check endpoint

## Quick Start

### Prerequisites

- Go 1.25+
- PostgreSQL database
- Docker (optional)

### Environment Variables

Create a `.env` file:

```env
DB_HOST=localhost
DB_USER=your_username
DB_PASS=your_password
DB_NAME=blog_db
DB_PORT=5432
DB_SSLMODE=disable
SCM_ORIGIN=https://api.github.com/repos
SCM_REPO=username/repo
```

### Running Locally

```bash
# Install dependencies
go mod download

# Run the application
go run .
```

The server starts on `http://localhost:8080`

### Using Docker

```bash
# Build image
docker build -t blog-backend .

# Run container
docker run -p 8080:8080 --env-file .env blog-backend
```

## API Endpoints

- `GET /ping` - Health check
- `GET /posts/` - Get posts
- `GET /posts/home` - Get README content from configured repository

## Database Schema

The application automatically creates a `posts` table with:
- `id` - UUID primary key

## Development

```bash
# Build binary
go build -o main

# Update dependencies
go mod tidy
```

## Tech Stack

- **Framework**: Gin
- **Database**: PostgreSQL with GORM
- **Language**: Go 1.25
- **Container**: Docker with Alpine Linux