package system

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bayou-brogrammer/mygo/internal/config"
	"github.com/bayou-brogrammer/mygo/internal/shell"
)

// Common development tools by platform
var commonTools = map[string][]string{
	"darwin": {
		"git", "vim", "tmux", "zsh", "curl", "wget", "jq", "ripgrep", "fd", "bat", "exa",
		"go", "node", "python", "ruby", "rust",
	},
	"linux": {
		"git", "vim", "tmux", "zsh", "curl", "wget", "jq", "ripgrep", "fd", "bat", "exa",
		"go", "nodejs", "python", "ruby", "rust",
	},
}

// Package managers by platform
var packageManagers = map[string]string{
	"darwin": "brew",
	"linux":  "apt", // Default to apt, can be changed based on distro
}

// InstallOptions defines options for tool installation
type InstallOptions struct {
	// Force reinstallation even if the tool is already installed
	Force bool

	// Skip installation if the tool is already installed
	SkipExisting bool

	// Enable verbose output during installation
	Verbose bool
}

// Install installs a specific tool or common development tools
func Install(tool string) error {
	return InstallWithOptions(tool, InstallOptions{
		SkipExisting: true,
		Verbose:      false,
	})
}

// InstallWithOptions installs a specific tool or common development tools with options
func InstallWithOptions(tool string, options InstallOptions) error {
	// Get package manager for current platform
	pkgManager, ok := packageManagers[runtime.GOOS]
	if !ok {
		return fmt.Errorf("unsupported platform: %s - only darwin (macOS) and linux are supported", runtime.GOOS)
	}

	// Check if package manager is installed
	if !shell.CommandExists(pkgManager) {
		switch pkgManager {
		case "brew":
			return fmt.Errorf("%s is not installed, please install Homebrew first: https://brew.sh", pkgManager)
		case "apt":
			return fmt.Errorf("%s is not installed, please install apt-get first", pkgManager)
		default:
			return fmt.Errorf("%s is not installed, please install it first", pkgManager)
		}
	}

	if tool != "" {
		// Check if tool is already installed and should be skipped
		if shell.CommandExists(tool) && options.SkipExisting && !options.Force {
			fmt.Printf("%s is already installed, skipping\n", tool)
			return nil
		}

		// Install specific tool
		fmt.Printf("Installing %s...\n", tool)
		var result *shell.Result
		var err error

		switch pkgManager {
		case "brew":
			if options.Force {
				result, err = shell.Execute("brew", "reinstall", tool)
			} else {
				result, err = shell.Execute("brew", "install", tool)
			}
		case "apt":
			if options.Force {
				result, err = shell.Execute("sudo", "apt", "install", "--reinstall", "-y", tool)
			} else {
				result, err = shell.Execute("sudo", "apt", "install", "-y", tool)
			}
		}

		if err != nil {
			return fmt.Errorf("failed to install %s: %w", tool, err)
		}

		shell.PrintResult(result, options.Verbose)
	} else {
		// Get the configuration to access preferred tools
		cfg, err := config.GetConfig()
		fmt.Printf("%v\n", cfg.PreferredTools)

		if err != nil {
			return fmt.Errorf("failed to get configuration: %w", err)
		}

		// Run install over config tools
		for _, tool := range cfg.PreferredTools {
			// Check if tool is already installed and should be skipped
			if shell.CommandExists(tool) && options.SkipExisting && !options.Force {
				fmt.Printf("%s is already installed, skipping\n", tool)
				continue
			}

			fmt.Printf("Installing %s...\n", tool)
			var result *shell.Result
			var err error

			switch pkgManager {
			case "brew":
				if options.Force {
					result, err = shell.Execute("brew", "reinstall", tool)
				} else {
					result, err = shell.Execute("brew", "install", tool)
				}
			case "apt":
				if options.Force {
					result, err = shell.Execute("sudo", "apt", "install", "--reinstall", "-y", tool)
				} else {
					result, err = shell.Execute("sudo", "apt", "install", "-y", tool)
				}
			}

			if err != nil {
				fmt.Printf("Failed to install %s: %v\n", tool, err)
				continue
			}

			shell.PrintResult(result, options.Verbose)
		}
	}

	return nil
}

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

// UpdateOptions defines options for system update
type UpdateOptions struct {
	// Enable verbose output during update
	Verbose bool
}

// Update updates installed development tools
func Update() error {
	return UpdateWithOptions(UpdateOptions{
		Verbose: false,
	})
}

// UpdateWithOptions updates installed development tools with options
func UpdateWithOptions(options UpdateOptions) error {
	// Get package manager for current platform
	pkgManager, ok := packageManagers[runtime.GOOS]
	if !ok {
		return fmt.Errorf("unsupported platform: %s - only darwin (macOS) and linux are supported", runtime.GOOS)
	}

	// Check if package manager is installed
	if !shell.CommandExists(pkgManager) {
		switch pkgManager {
		case "brew":
			return fmt.Errorf("%s is not installed, please install Homebrew first: https://brew.sh", pkgManager)
		case "apt":
			return fmt.Errorf("%s is not installed, please install apt-get first", pkgManager)
		default:
			return fmt.Errorf("%s is not installed, please install it first", pkgManager)
		}
	}

	// Update package manager
	var result *shell.Result
	var err error

	switch pkgManager {
	case "brew":
		fmt.Println("Updating Homebrew...")
		result, err = shell.Execute("brew", "update")
		if err != nil {
			return fmt.Errorf("failed to update Homebrew: %w", err)
		}
		shell.PrintResult(result, options.Verbose)

		fmt.Println("Upgrading packages...")
		result, err = shell.Execute("brew", "upgrade")
	case "apt":
		fmt.Println("Updating package lists...")
		result, err = shell.Execute("sudo", "apt", "update")
		if err != nil {
			return fmt.Errorf("failed to update package lists: %w", err)
		}
		shell.PrintResult(result, options.Verbose)

		fmt.Println("Upgrading packages...")
		result, err = shell.Execute("sudo", "apt", "upgrade", "-y")
	}

	if err != nil {
		return fmt.Errorf("failed to upgrade packages: %w", err)
	}

	shell.PrintResult(result, options.Verbose)
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
