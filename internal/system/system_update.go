package system

import (
	"fmt"
	"runtime"

	"github.com/bayou-brogrammer/mygo/internal/shell"
)

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
