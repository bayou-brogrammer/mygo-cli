package repo

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bayou-brogrammer/mygo/internal/config"
	"github.com/bayou-brogrammer/mygo/internal/shell"
)

// Clone clones a GitHub repository
func Clone(url string, destDir string) error {
	// Extract repository name from URL
	parts := strings.Split(url, "/")
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Clone the repository
	result, err := shell.Execute("git", "clone", url, filepath.Join(destDir, repoName))
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	shell.PrintResult(result, true)

	// Track the repository
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	cfg.TrackedRepos[repoName] = config.Repository{
		URL:         url,
		Path:        filepath.Join(destDir, repoName),
		Description: "",
		LastUpdated: time.Now().Format(time.RFC3339),
	}

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// Update updates a tracked repository
func Update(repoName string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	repo, exists := cfg.TrackedRepos[repoName]
	if !exists {
		return fmt.Errorf("repository not found: %s", repoName)
	}

	// Pull latest changes
	result, err := shell.ExecuteInDir(repo.Path, "git", "pull")
	if err != nil {
		return fmt.Errorf("failed to update repository: %w", err)
	}

	shell.PrintResult(result, true)

	// Update last updated timestamp
	repo.LastUpdated = time.Now().Format(time.RFC3339)
	cfg.TrackedRepos[repoName] = repo

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// UpdateAll updates all tracked repositories
func UpdateAll() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	for repoName := range cfg.TrackedRepos {
		fmt.Printf("Updating repository: %s\n", repoName)
		if err := Update(repoName); err != nil {
			fmt.Printf("Error updating %s: %v\n", repoName, err)
		}
	}

	return nil
}

// List lists all tracked repositories
func List() ([]config.Repository, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	repos := make([]config.Repository, 0, len(cfg.TrackedRepos))
	for _, repo := range cfg.TrackedRepos {
		repos = append(repos, repo)
	}

	return repos, nil
}

// Remove removes a tracked repository (without deleting files)
func Remove(repoName string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if _, exists := cfg.TrackedRepos[repoName]; !exists {
		return fmt.Errorf("repository not found: %s", repoName)
	}

	delete(cfg.TrackedRepos, repoName)

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}

// Delete removes a tracked repository and deletes the files
func Delete(repoName string) error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	repo, exists := cfg.TrackedRepos[repoName]
	if !exists {
		return fmt.Errorf("repository not found: %s", repoName)
	}

	// Delete the repository directory
	if err := os.RemoveAll(repo.Path); err != nil {
		return fmt.Errorf("failed to delete repository: %w", err)
	}

	// Remove from tracked repositories
	delete(cfg.TrackedRepos, repoName)

	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
