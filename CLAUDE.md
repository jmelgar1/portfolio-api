# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go portfolio API project using Go 1.24.5. The project follows a clean architecture pattern with handlers, services, and a modular structure.

## Development Guidelines

### Code Standards
- Use Go 1.24.5 features and syntax
- Follow Go naming conventions (PascalCase for exported, camelCase for unexported)
- Use structured logging with the `log` package
- Organize code in packages by domain/feature (e.g., `user`, `images`)
- Keep handlers thin - business logic should be in services

### Project Structure
```
cmd/
├── main.go           # Application entry point
├── api/              # API server setup and routing
└── service/          # Business logic organized by domain
    └── user/         # User-related handlers and logic
```

### Development Commands

#### Running the Application
```bash
go run cmd/main.go
```

#### Building the Application
```bash
go build -o portfolio-api cmd/main.go
```

#### Running Tests
```bash
go test ./...
```

#### Module Management
```bash
go mod tidy    # Clean up dependencies
go mod vendor  # Create vendor directory
```

## Architecture Guidelines

### HTTP Routing
- Use `http.NewServeMux()` for routing
- Organize routes by API version (`/api/v1/`)
- Use method-specific patterns where supported (e.g., `"GET /images"`)
- Register routes in dedicated handler files

### Error Handling
- Always handle errors appropriately
- Use structured logging for debugging
- Return appropriate HTTP status codes

### Code Organization
- Keep handlers in `cmd/service/{domain}/routes.go`
- Use dependency injection for database connections
- Create separate packages for different business domains

## Change Tracking

When making modifications to this codebase, maintain a record of changes in `claude/memories/CHANGELOG.md`. Include:
- What was changed and why
- Files modified
- Key decisions made
- Any breaking changes or considerations for future development
- For changelogs created during sessions in the @claude/memories/ folder, ensure they are named properly in the following format (mm-dd-yyyy s1, 2, 3, etc.), the "s" stands for session and we will increment as necessary.