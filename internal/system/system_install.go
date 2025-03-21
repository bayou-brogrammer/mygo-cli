package system

import (
	"fmt"
	"runtime"

	"github.com/bayou-brogrammer/mygo/internal/config"
	"github.com/bayou-brogrammer/mygo/internal/logger"
	"github.com/bayou-brogrammer/mygo/internal/shell"
	"github.com/bayou-brogrammer/mygo/internal/ui"
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
		logger.Error("Unsupported platform: %s - only darwin (macOS) and linux are supported", runtime.GOOS)
		return fmt.Errorf("unsupported platform: %s - only darwin (macOS) and linux are supported", runtime.GOOS)
	}

	// Check if package manager is installed
	if !shell.CommandExists(pkgManager) {
		var errMsg error
		switch pkgManager {
		case "brew":
			errMsg = fmt.Errorf("%s is not installed, please install Homebrew first: https://brew.sh", pkgManager)
		case "apt":
			errMsg = fmt.Errorf("%s is not installed, please install apt-get first", pkgManager)
		default:
			errMsg = fmt.Errorf("%s is not installed, please install it first", pkgManager)
		}

		return errMsg
	}

	if tool != "" {
		// Install specific tool
		ui.PrintInfo("Installing tool: %s", tool)
		InstallWithPackageManagers(tool, pkgManager, options)
	} else {
		// Get the configuration to access preferred tools
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get configuration: %w", err)
		}

		// Run install over config tools
		ui.PrintTitle("Installing Tools from Configuration")

		for _, tool := range cfg.Tools {
			err := InstallWithPackageManagers(tool, pkgManager, options)
			if err != nil {
				return err
			}
		}

		ui.PrintSuccess("Finished installing all tools from configuration")
	}

	return nil
}

func InstallWithPackageManagers(tool string, pkgManager string, options InstallOptions) error {
	// Format the tool name with accent color for better visibility
	highlightedTool := ui.FormatTextWithColor(tool, &ui.StyleCommand, ui.ColorInfo)

	// Check if tool is already installed and should be skipped
	if shell.CommandExists(tool) && options.SkipExisting && !options.Force {
		ui.PrintInfo("%s is already installed, skipping", highlightedTool)
		return nil
	}

	ui.PrintCommand("Installing %s", highlightedTool)

	var result *shell.Result
	var err error

	// Display the command being executed with proper formatting
	switch pkgManager {
	case "brew":
		if options.Force {
			// Show the command with nice formatting
			ui.PrintInfo("Running: %s %s %s",
				ui.FormatCommand("brew"),
				ui.FormatValue("reinstall"),
				highlightedTool)
			result, err = shell.Execute("brew", "reinstall", tool)
		} else {
			// Show the command with nice formatting
			ui.PrintInfo("Running: %s %s %s",
				ui.FormatCommand("brew"),
				ui.FormatValue("install"),
				highlightedTool)
			result, err = shell.Execute("brew", "install", tool)
		}
	case "apt":
		if options.Force {
			// Show the command with nice formatting
			ui.PrintInfo("Running: %s %s %s %s %s %s",
				ui.FormatCommand("sudo"),
				ui.FormatCommand("apt"),
				ui.FormatValue("install"),
				ui.FormatValue("--reinstall"),
				ui.FormatValue("-y"),
				highlightedTool)
			result, err = shell.Execute("sudo", "apt", "install", "--reinstall", "-y", tool)
		} else {
			// Show the command with nice formatting
			ui.PrintInfo("Running: %s %s %s %s %s",
				ui.FormatCommand("sudo"),
				ui.FormatCommand("apt"),
				ui.FormatValue("install"),
				ui.FormatValue("-y"),
				highlightedTool)
			result, err = shell.Execute("sudo", "apt", "install", "-y", tool)
		}
	}

	if err != nil {
		return fmt.Errorf("failed to install %s: %v", tool, err)
	}

	// Print success message
	ui.PrintSuccess("Successfully installed %s", highlightedTool)

	// Print detailed output if verbose
	if options.Verbose {
		if result.Stdout != "" {
			ui.PrintBox(fmt.Sprintf("Output:\n%s", result.Stdout))
		}
		if result.Stderr != "" {
			ui.PrintErrorBox(fmt.Sprintf("Errors:\n%s", result.Stderr))
		}
	}

	return nil
}
