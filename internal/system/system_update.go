package system

import (
	"fmt"
	"runtime"

	"github.com/bayou-brogrammer/mygo/internal/logger"
	"github.com/bayou-brogrammer/mygo/internal/shell"
	"github.com/bayou-brogrammer/mygo/internal/ui"
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
		errMsg := fmt.Sprintf("Unsupported platform: %s - only darwin (macOS) and linux are supported", runtime.GOOS)
		logger.Error(errMsg)
		ui.PrintError(errMsg)
		return fmt.Errorf(errMsg)
	}

	// Check if package manager is installed
	if !shell.CommandExists(pkgManager) {
		var errMsg string
		switch pkgManager {
		case "brew":
			errMsg = fmt.Sprintf("%s is not installed, please install Homebrew first: https://brew.sh", pkgManager)
		case "apt":
			errMsg = fmt.Sprintf("%s is not installed, please install apt-get first", pkgManager)
		default:
			errMsg = fmt.Sprintf("%s is not installed, please install it first", pkgManager)
		}
		logger.Error(errMsg)
		ui.PrintError(errMsg)
		return fmt.Errorf(errMsg)
	}

	// Update package manager
	var result *shell.Result
	var err error

	ui.PrintTitle("System Update")

	switch pkgManager {
	case "brew":
		ui.PrintInfo("Updating Homebrew...")
		result, err = shell.Execute("brew", "update")
		if err != nil {
			errMsg := fmt.Sprintf("Failed to update Homebrew: %v", err)
			logger.Error(errMsg)
			ui.PrintError(errMsg)
			return fmt.Errorf("failed to update Homebrew: %w", err)
		}
		if options.Verbose {
			if result.Stdout != "" {
				ui.PrintBox(fmt.Sprintf("Output:\n%s", result.Stdout))
			}
		}

		ui.PrintInfo("Upgrading Homebrew packages...")
		result, err = shell.Execute("brew", "upgrade")
	case "apt":
		ui.PrintInfo("Updating apt repositories...")
		result, err = shell.Execute("sudo", "apt", "update")
		if err != nil {
			errMsg := fmt.Sprintf("Failed to update apt repositories: %v", err)
			logger.Error(errMsg)
			ui.PrintError(errMsg)
			return fmt.Errorf("failed to update apt repositories: %w", err)
		}
		if options.Verbose {
			if result.Stdout != "" {
				ui.PrintBox(fmt.Sprintf("Output:\n%s", result.Stdout))
			}
		}

		ui.PrintInfo("Upgrading apt packages...")
		result, err = shell.Execute("sudo", "apt", "upgrade", "-y")
	}

	if err != nil {
		errMsg := fmt.Sprintf("Failed to upgrade packages: %v", err)
		logger.Error(errMsg)
		ui.PrintError(errMsg)
		return fmt.Errorf("failed to upgrade packages: %w", err)
	}

	if options.Verbose {
		if result.Stdout != "" {
			ui.PrintBox(fmt.Sprintf("Output:\n%s", result.Stdout))
		}
	}

	ui.PrintSuccess("System update completed successfully")
	return nil
}
