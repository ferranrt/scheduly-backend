# DBTools CLI

A command-line interface for managing the Scheduly backend database operations.

## Installation

The CLI is part of the Scheduly backend project. To use it, navigate to the project root and run:

```bash
go run cmd/dbtools/main.go [command]
```

## Available Commands

### `migrate`

Run database migrations to create or update the database schema.

```bash
go run cmd/dbtools/main.go migrate
```

This command will:

- Connect to the database using the configuration from environment variables
- Execute all pending migrations
- Create or update tables based on the defined models

### `cleanup`

Drop all tables and recreate them from scratch. **Use with caution!**

```bash
go run cmd/dbtools/main.go cleanup
```

This command will:

- Connect to the database
- Drop all existing tables
- Recreate all tables from scratch
- Useful for development/testing environments

### `--version`

Display the CLI version.

```bash
go run cmd/dbtools/main.go --version
```

### `--help`

Display help information for any command.

```bash
go run cmd/dbtools/main.go --help
go run cmd/dbtools/main.go migrate --help
go run cmd/dbtools/main.go cleanup --help
```

## Configuration

The CLI uses the same configuration as the main application. Make sure your environment variables are properly set:

- Database connection parameters
- Any other required configuration

## Examples

```bash
# Run migrations
go run cmd/dbtools/main.go migrate

# Clean up database (development only)
go run cmd/dbtools/main.go cleanup

# Check version
go run cmd/dbtools/main.go --version

# Get help
go run cmd/dbtools/main.go --help
```
