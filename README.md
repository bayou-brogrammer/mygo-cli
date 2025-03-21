# Milo CLI

A command-line tool for managing development environment setup, including GitHub repositories, dotfiles, and system configuration via chezmoi.

## Features

- GitHub repository management (clone, update, list)
- Dotfiles management and synchronization
- Chezmoi integration for configuration management
- System configuration and tool installation

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/milo-cli.git

# Build the binary
cd milo
go build -o milo

# Move to a directory in your PATH
mv milo /usr/local/bin/
```

## Usage

```bash
# Initialize dotfiles
milo dots init

# Clone a repository
milo repo clone https://github.com/username/repo.git

# Apply chezmoi configuration
milo chezmoi apply

# List repositories
milo repo list

# Install development tools
milo system install
```

## Configuration

Configuration is stored in `~/.config/milo/config.yaml`. You can edit this file directly or use the CLI to update settings.

## Development

This project uses Go modules for dependency management.

```bash
# Get dependencies
go mod tidy

# Run tests
go test ./...

# Build the project
go build
```

## License

MIT
