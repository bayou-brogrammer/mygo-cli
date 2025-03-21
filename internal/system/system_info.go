package system

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/bayou-brogrammer/mygo/internal/logger"
	"github.com/bayou-brogrammer/mygo/internal/shell"
	"github.com/bayou-brogrammer/mygo/internal/ui"
)

// InfoOptions defines options for displaying system information
type InfoOptions struct {
	// Enable verbose output
	Verbose bool
}

// Info displays system information
func Info() error {
	return InfoWithOptions(InfoOptions{
		Verbose: false,
	})
}

// InfoWithOptions displays system information with options
func InfoWithOptions(options InfoOptions) error {
	ui.PrintTitle("System Information")

	// Display OS information
	osInfo := fmt.Sprintf("OS: %s (%s)", runtime.GOOS, runtime.GOARCH)
	ui.PrintInfo(osInfo)

	// Display Go version
	goVersion := fmt.Sprintf("Go Version: %s", runtime.Version())
	ui.PrintInfo(goVersion)

	// Get package manager information
	pkgManager, ok := packageManagers[runtime.GOOS]
	if !ok {
		errMsg := fmt.Sprintf("Unsupported platform: %s", runtime.GOOS)
		logger.Error(errMsg)
		ui.PrintError(errMsg)
		return fmt.Errorf(errMsg)
	}

	// Check if package manager is installed
	if !shell.CommandExists(pkgManager) {
		errMsg := fmt.Sprintf("Package manager %s is not installed", pkgManager)
		ui.PrintWarning(errMsg)
	} else {
		// Get package manager version
		var version string
		var err error

		switch pkgManager {
		case "brew":
			cmd := exec.Command("brew", "--version")
			output, err := cmd.Output()
			if err == nil {
				// Extract just the first line
				lines := strings.Split(string(output), "\n")
				if len(lines) > 0 {
					version = lines[0]
				}
			}
		case "apt":
			cmd := exec.Command("apt", "--version")
			output, err := cmd.Output()
			if err == nil {
				// Extract just the first line
				lines := strings.Split(string(output), "\n")
				if len(lines) > 0 {
					version = lines[0]
				}
			}
		}

		if err != nil {
		} else if version != "" {
			pkgInfo := fmt.Sprintf("Package Manager: %s", version)
			ui.PrintInfo(pkgInfo)
		}
	}

	// Display installed tools if verbose
	if options.Verbose {
		ui.PrintSubtitle("Installed Tools")

		// Common developer tools to check
		tools := []string{
			"git", "node", "npm", "python", "pip", "docker",
			"kubectl", "terraform", "ansible", "make", "gcc",
		}

		for _, tool := range tools {
			if shell.CommandExists(tool) {
				// Try to get version
				var version string
				result, err := shell.Execute(tool, "--version")
				if err == nil && result.Stdout != "" {
					// Extract just the first line
					lines := strings.Split(result.Stdout, "\n")
					if len(lines) > 0 {
						version = strings.TrimSpace(lines[0])
					}
				}

				if version != "" {
					toolInfo := fmt.Sprintf("%s: %s", tool, version)
					ui.PrintInfo(toolInfo)
				} else {
					ui.PrintInfo("%s: installed (version unknown)", tool)
				}
			}
		}
	}

	ui.PrintSuccess("System information gathered successfully")
	return nil
}
