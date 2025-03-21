package main

import (
	"fmt"

	"github.com/bayou-brogrammer/mygo/internal/system"
	"github.com/bayou-brogrammer/mygo/internal/ui"
	"github.com/spf13/cobra"
)

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Manage system configuration",
	Long:  `Commands for installing tools and configuring system preferences.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var (
	forceInstall   bool
	skipExisting   bool
	nonInteractive bool
)

var systemInstallCmd = &cobra.Command{
	Use:   "install [tool name]",
	Short: "Install development tools",
	Long:  `Install development tools from your configuration or a specific tool.`,
	Run: func(cmd *cobra.Command, args []string) {
		var toolName string
		if len(args) > 0 {
			toolName = args[0]
		}

		// Set options based on flags
		options := system.InstallOptions{
			Force:        forceInstall,
			SkipExisting: skipExisting,
			Verbose:      verbose,
		}

		err := system.InstallWithOptions(toolName, options)
		if err != nil {
			ui.PrintError("%v", err)
		}
	},
}

var systemConfigureCmd = &cobra.Command{
	Use:   "configure [component]",
	Short: "Configure system preferences",
	Long:  `Configure system preferences for development.`,
	Run: func(cmd *cobra.Command, args []string) {
		var component string
		if len(args) > 0 {
			component = args[0]
		}

		// Set options based on flags
		options := system.ConfigureOptions{
			NonInteractive: nonInteractive,
			Verbose:        verbose,
		}

		err := system.ConfigureWithOptions(component, options)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

var systemUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update system tools",
	Long:  `Update installed development tools.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Set options based on flags
		options := system.UpdateOptions{
			Verbose: verbose,
		}

		err := system.UpdateWithOptions(options)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}

func init() {
	// Add flags to install command
	systemInstallCmd.Flags().BoolVarP(&forceInstall, "force", "f", false, "Force installation even if the tool is already installed")
	systemInstallCmd.Flags().BoolVarP(&skipExisting, "skip-existing", "s", true, "Skip installation if the tool is already installed")

	// Add flags to configure command
	systemConfigureCmd.Flags().BoolVarP(&nonInteractive, "non-interactive", "n", false, "Run in non-interactive mode (requires environment variables for input)")

	systemCmd.AddCommand(systemInstallCmd)
	systemCmd.AddCommand(systemConfigureCmd)
	systemCmd.AddCommand(systemUpdateCmd)
	rootCmd.AddCommand(systemCmd)
}
