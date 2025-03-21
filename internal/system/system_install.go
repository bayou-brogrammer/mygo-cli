package system

import (
	"fmt"
	"runtime"

	"github.com/bayou-brogrammer/mygo/internal/config"
	"github.com/bayou-brogrammer/mygo/internal/shell"
)

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
		// Install specific tool
		InstallWithPackageManagers(tool, pkgManager, options)
	} else {
		// Get the configuration to access preferred tools
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get configuration: %w", err)
		}

		// Run install over config tools
		for _, tool := range cfg.Tools {
			InstallWithPackageManagers(tool, pkgManager, options)
		}
	}

	return nil
}

func InstallWithPackageManagers(tool string, pkgManager string, options InstallOptions) {
	// Check if tool is already installed and should be skipped
	if shell.CommandExists(tool) && options.SkipExisting && !options.Force {
		fmt.Printf("%s is already installed, skipping\n", tool)
		return
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
		return
	}

	shell.PrintResult(result, options.Verbose)
}
