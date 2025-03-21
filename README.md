# DevEnv CLI

A command-line tool for managing development environment setup, including GitHub repositories, dotfiles, and system configuration via chezmoi.

## Features

- GitHub repository management (clone, update, list)
- Dotfiles management and synchronization
- Chezmoi integration for configuration management
- System configuration and tool installation

## Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/devenv-cli.git

# Build the binary
cd devenv-cli
go build -o devenv

# Move to a directory in your PATH
mv devenv /usr/local/bin/
```

## Usage

```bash
# Initialize dotfiles
devenv dots init

# Clone a repository
devenv repo clone https://github.com/username/repo.git

# Apply chezmoi configuration
devenv chezmoi apply

# Install development tools
devenv system install
```

## Configuration

Configuration is stored in `~/.config/devenv/config.yaml`. You can edit this file directly or use the CLI to update settings.

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
