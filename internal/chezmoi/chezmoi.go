package chezmoi

import (
	"fmt"
	"os"

	"github.com/bayou-brogrammer/mygo/internal/config"
	"github.com/bayou-brogrammer/mygo/internal/shell"
)

// Init initializes chezmoi with an optional dotfiles repository
func Init(repoURL string) error {
	// Check if chezmoi is installed
	if !shell.CommandExists("chezmoi") {
		return fmt.Errorf("chezmoi is not installed, please install it first")
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	// Initialize chezmoi
	var result *shell.Result
	if repoURL != "" {
		// Initialize with repository
		result, err = shell.Execute("chezmoi", "init", repoURL)
	} else {
		// Initialize without repository
		result, err = shell.Execute("chezmoi", "init")
	}

	if err != nil {
		return fmt.Errorf("failed to initialize chezmoi: %w", err)
	}

	shell.PrintResult(result, true)

	// Update configuration
	chezmoiDir, err := getChezmoiDir()
	if err != nil {
		return fmt.Errorf("failed to get chezmoi directory: %w", err)
	}

	cfg.ChezmoiDir = chezmoiDir
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// Apply applies chezmoi configuration to the system
func Apply() error {
	// Check if chezmoi is installed
	if !shell.CommandExists("chezmoi") {
		return fmt.Errorf("chezmoi is not installed, please install it first")
	}

	// Apply configuration
	result, err := shell.Execute("chezmoi", "apply")
	if err != nil {
		return fmt.Errorf("failed to apply chezmoi configuration: %w", err)
	}

	shell.PrintResult(result, true)
	return nil
}

// Update updates chezmoi-managed files from the source repository
func Update() error {
	// Check if chezmoi is installed
	if !shell.CommandExists("chezmoi") {
		return fmt.Errorf("chezmoi is not installed, please install it first")
	}

	// Update source repository
	result, err := shell.Execute("chezmoi", "update")
	if err != nil {
		return fmt.Errorf("failed to update chezmoi: %w", err)
	}

	shell.PrintResult(result, true)
	return nil
}

// Add adds a file to be managed by chezmoi
func Add(filePath string) error {
	// Check if chezmoi is installed
	if !shell.CommandExists("chezmoi") {
		return fmt.Errorf("chezmoi is not installed, please install it first")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Add file to chezmoi
	result, err := shell.Execute("chezmoi", "add", filePath)
	if err != nil {
		return fmt.Errorf("failed to add file to chezmoi: %w", err)
	}

	shell.PrintResult(result, true)
	return nil
}

// getChezmoiDir gets the chezmoi source directory
func getChezmoiDir() (string, error) {
	result, err := shell.Execute("chezmoi", "source-path")
	if err != nil {
		return "", fmt.Errorf("failed to get chezmoi source path: %w", err)
	}

	return result.Stdout, nil
}
