package system

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bayou-brogrammer/mygo/internal/shell"
)

// ConfigureOptions defines options for system configuration
type ConfigureOptions struct {
	// Run in non-interactive mode, using environment variables instead of prompts
	NonInteractive bool

	// Enable verbose output during configuration
	Verbose bool
}

// Configure configures system preferences for development
func Configure(component string) error {
	return ConfigureWithOptions(component, ConfigureOptions{
		NonInteractive: false,
		Verbose:        false,
	})
}

// ConfigureWithOptions configures system preferences with options
func ConfigureWithOptions(component string, options ConfigureOptions) error {
	if component != "" {
		// Configure specific component
		switch component {
		case "git":
			return configureGit(options)
		case "shell":
			return configureShell(options)
		default:
			return fmt.Errorf("unknown component: %s", component)
		}
	} else {
		// Configure all components
		if err := configureGit(options); err != nil {
			fmt.Printf("Failed to configure git: %v\n", err)
		}

		if err := configureShell(options); err != nil {
			fmt.Printf("Failed to configure shell: %v\n", err)
		}
	}

	return nil
}

// configureGit configures git settings
func configureGit(options ConfigureOptions) error {
	// Check if git is installed
	if !shell.CommandExists("git") {
		return fmt.Errorf("git is not installed, please install it first")
	}

	// Get git user name and email
	var username, email string

	if options.NonInteractive {
		// In non-interactive mode, try to get values from environment
		username = os.Getenv("GIT_USERNAME")
		email = os.Getenv("GIT_EMAIL")
		if username == "" || email == "" {
			return fmt.Errorf("in non-interactive mode, GIT_USERNAME and GIT_EMAIL environment variables must be set")
		}
	} else {
		// Prompt for git user name and email
		fmt.Print("Enter your Git username: ")
		fmt.Scanln(&username)

		fmt.Print("Enter your Git email: ")
		fmt.Scanln(&email)
	}

	// Configure git
	if username != "" {
		result, err := shell.Execute("git", "config", "--global", "user.name", username)
		if err != nil {
			return fmt.Errorf("failed to configure git username: %w", err)
		}
		shell.PrintResult(result, options.Verbose)
	}

	if email != "" {
		result, err := shell.Execute("git", "config", "--global", "user.email", email)
		if err != nil {
			return fmt.Errorf("failed to configure git email: %w", err)
		}
		shell.PrintResult(result, options.Verbose)
	}

	// Configure common git aliases
	aliases := map[string]string{
		"co":      "checkout",
		"br":      "branch",
		"ci":      "commit",
		"st":      "status",
		"unstage": "reset HEAD --",
		"last":    "log -1 HEAD",
	}

	for alias, command := range aliases {
		result, err := shell.Execute("git", "config", "--global", "alias."+alias, command)
		if err != nil {
			fmt.Printf("Failed to configure git alias %s: %v\n", alias, err)
			continue
		}
		shell.PrintResult(result, false)
	}

	fmt.Println("Git configured successfully")
	return nil
}

// configureShell configures shell settings
func configureShell(options ConfigureOptions) error {
	// Check if zsh is installed
	if !shell.CommandExists("zsh") {
		return fmt.Errorf("zsh is not installed, please install it first")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Check if Oh My Zsh is installed
	ohmyzshDir := filepath.Join(homeDir, ".oh-my-zsh")
	if _, err := os.Stat(ohmyzshDir); os.IsNotExist(err) {
		// Install Oh My Zsh
		fmt.Println("Installing Oh My Zsh...")
		result, err := shell.Execute("sh", "-c", "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)")
		if err != nil {
			return fmt.Errorf("failed to install Oh My Zsh: %w", err)
		}
		shell.PrintResult(result, options.Verbose)
	}

	fmt.Println("Shell configured successfully")
	return nil
}
