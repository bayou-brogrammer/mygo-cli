package dots

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bayou-brogrammer/mygo/internal/config"
	"github.com/bayou-brogrammer/mygo/internal/shell"
)

// Init initializes dotfiles from a repository or creates a new dotfiles repository
func Init(repoURL string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	// Create dotfiles directory if it doesn't exist
	if err := os.MkdirAll(cfg.DotfilesDir, 0755); err != nil {
		return fmt.Errorf("failed to create dotfiles directory: %w", err)
	}

	if repoURL != "" {
		// Clone the dotfiles repository
		result, err := shell.Execute("git", "clone", repoURL, cfg.DotfilesDir)
		if err != nil {
			return fmt.Errorf("failed to clone dotfiles repository: %w", err)
		}

		shell.PrintResult(result, true)

		// Update configuration
		cfg.DotfilesRepo = repoURL
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
	} else {
		// Initialize a new git repository
		result, err := shell.ExecuteInDir(cfg.DotfilesDir, "git", "init")
		if err != nil {
			return fmt.Errorf("failed to initialize git repository: %w", err)
		}

		shell.PrintResult(result, true)
	}

	return nil
}

// Apply applies dotfiles configuration to the system
func Apply() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	// Check if dotfiles directory exists
	if _, err := os.Stat(cfg.DotfilesDir); os.IsNotExist(err) {
		return fmt.Errorf("dotfiles directory not found: %s", cfg.DotfilesDir)
	}

	// Apply dotfiles (this is a simple implementation, can be expanded)
	// For now, we'll just create symlinks for all files in the dotfiles directory
	err = filepath.Walk(cfg.DotfilesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path from dotfiles directory
		relPath, err := filepath.Rel(cfg.DotfilesDir, path)
		if err != nil {
			return err
		}

		// Skip files in the root directory that don't start with a dot
		if !strings.HasPrefix(filepath.Base(relPath), ".") && filepath.Dir(relPath) == "." {
			return nil
		}

		// Create target path in home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		targetPath := filepath.Join(homeDir, relPath)

		// Create parent directories if they don't exist
		targetDir := filepath.Dir(targetPath)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return err
		}

		// Remove existing file or symlink
		if _, err := os.Lstat(targetPath); err == nil {
			if err := os.Remove(targetPath); err != nil {
				return err
			}
		}

		// Create symlink
		if err := os.Symlink(path, targetPath); err != nil {
			return err
		}

		fmt.Printf("Linked %s -> %s\n", targetPath, path)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to apply dotfiles: %w", err)
	}

	return nil
}

// Update updates dotfiles from the repository
func Update() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	// Check if dotfiles directory exists
	if _, err := os.Stat(cfg.DotfilesDir); os.IsNotExist(err) {
		return fmt.Errorf("dotfiles directory not found: %s", cfg.DotfilesDir)
	}

	// Check if it's a git repository
	gitDir := filepath.Join(cfg.DotfilesDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository: %s", cfg.DotfilesDir)
	}

	// Pull latest changes
	result, err := shell.ExecuteInDir(cfg.DotfilesDir, "git", "pull")
	if err != nil {
		return fmt.Errorf("failed to update dotfiles: %w", err)
	}

	shell.PrintResult(result, true)
	return nil
}

// Add adds a file to the dotfiles repository
func Add(filePath string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	// Check if dotfiles directory exists
	if _, err := os.Stat(cfg.DotfilesDir); os.IsNotExist(err) {
		return fmt.Errorf("dotfiles directory not found: %s", cfg.DotfilesDir)
	}

	// Check if it's a git repository
	gitDir := filepath.Join(cfg.DotfilesDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository: %s", cfg.DotfilesDir)
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filePath)
	}

	// Get absolute path of the file
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check if file is in home directory
	if !strings.HasPrefix(absPath, homeDir) {
		return fmt.Errorf("file must be in home directory: %s", absPath)
	}

	// Get relative path from home directory
	relPath, err := filepath.Rel(homeDir, absPath)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %w", err)
	}

	// Create target path in dotfiles directory
	targetPath := filepath.Join(cfg.DotfilesDir, relPath)

	// Create parent directories if they don't exist
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Copy file to dotfiles directory
	if err := copyFile(absPath, targetPath); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Add file to git
	result, err := shell.ExecuteInDir(cfg.DotfilesDir, "git", "add", relPath)
	if err != nil {
		return fmt.Errorf("failed to add file to git: %w", err)
	}

	shell.PrintResult(result, true)

	// Commit changes
	result, err = shell.ExecuteInDir(cfg.DotfilesDir, "git", "commit", "-m", fmt.Sprintf("Add %s", relPath))
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	shell.PrintResult(result, true)

	// Create symlink back to original location
	if err := os.Remove(absPath); err != nil {
		return fmt.Errorf("failed to remove original file: %w", err)
	}

	if err := os.Symlink(targetPath, absPath); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}

	fmt.Printf("Added %s to dotfiles\n", relPath)
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Write to destination file
	return os.WriteFile(dst, data, 0644)
}
