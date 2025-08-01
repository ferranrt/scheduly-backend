# Scheduly Backend

A Go-based backend service for the Scheduly application.

## Database Tools

This project includes a CLI tool for database management operations. See [cmd/dbtools/README.md](cmd/dbtools/README.md) for detailed usage instructions.

### Quick Start

```bash
# Run database migrations
go run cmd/dbtools/main.go migrate

# Clean up database (development only)
go run cmd/dbtools/main.go cleanup
```
