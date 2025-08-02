# Scheduly Backend

A Go-based backend service for the Scheduly application.

## Environment Configuration

### JWT Configuration

The application supports configurable JWT token durations through environment variables:

- `JWT_ACCESS_TOKEN_DURATION`: Duration for access tokens (e.g., "15m", "1h")
- `JWT_REFRESH_TOKEN_DURATION`: Duration for refresh tokens (e.g., "24h", "30d")
- `JWT_DURATION`: Fallback duration for both token types (for backward compatibility)
- `JWT_SECRET_KEY`: Secret key for signing JWT tokens

**Priority order:**
1. `JWT_ACCESS_TOKEN_DURATION` / `JWT_REFRESH_TOKEN_DURATION` (specific)
2. `JWT_DURATION` (fallback)
3. Default values (15 minutes for access, 30 days for refresh)

**Examples:**
```bash
# Set specific durations
JWT_ACCESS_TOKEN_DURATION=30m
JWT_REFRESH_TOKEN_DURATION=7d

# Or use fallback for both
JWT_DURATION=1h
```

## Database Tools

This project includes a CLI tool for database management operations. See [cmd/dbtools/README.md](cmd/dbtools/README.md) for detailed usage instructions.

### Quick Start

```bash
# Run database migrations
go run cmd/dbtools/main.go migrate

# Clean up database (development only)
go run cmd/dbtools/main.go cleanup
```
