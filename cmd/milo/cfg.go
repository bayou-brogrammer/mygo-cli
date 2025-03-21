package main

import (
	"fmt"

	"slices"

	"github.com/bayou-brogrammer/mygo/internal/config"
	"github.com/spf13/cobra"
)

var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "Manage configuration",
	Long:  `Commands for managing milo configuration files.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var cfgRegenerateCmd = &cobra.Command{
	Use:   "regenerate",
	Short: "Regenerate configuration",
	Long:  `Regenerate configuration from the repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Regenerating configuration")
		// TODO: Implement configuration regeneration logic
	},
}

var cfgToolsCmd = &cobra.Command{
	Use:   "tools",
	Short: "Manage preferred tools",
	Long:  `View and manage the list of preferred tools in your configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the configuration
		cfg, err := config.GetConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Display the list of preferred tools
		fmt.Println("Tools in your configuration:")
		for i, tool := range cfg.Tools {
			fmt.Printf("%d. %s\n", i+1, tool)
		}
	},
}

var cfgToolsAddCmd = &cobra.Command{
	Use:   "add [tool name]",
	Short: "Add a tool to preferred tools",
	Long:  `Add a tool to the list of preferred tools in your configuration.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		toolName := args[0]

		// Get the configuration
		cfg, err := config.GetConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Check if the tool is already in the list
		if slices.Contains(cfg.Tools, toolName) {
			fmt.Printf("%s is already in your tools list\n", toolName)
			return
		}

		// Add the tool to the list
		cfg.Tools = append(cfg.Tools, toolName)

		// Save the configuration
		if err := cfg.Save(); err != nil {
			fmt.Printf("Error saving configuration: %v\n", err)
			return
		}

		fmt.Printf("%s added to your preferred tools list\n", toolName)
	},
}

var cfgToolsRemoveCmd = &cobra.Command{
	Use:   "remove [tool name]",
	Short: "Remove a tool from preferred tools",
	Long:  `Remove a tool from the list of preferred tools in your configuration.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		toolName := args[0]

		// Get the configuration
		cfg, err := config.GetConfig()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		// Check if the tool is in the list
		found := false
		newTools := []string{}
		for _, tool := range cfg.Tools {
			if tool == toolName {
				found = true
			} else {
				newTools = append(newTools, tool)
			}
		}

		if !found {
			fmt.Printf("%s is not in your tools list\n", toolName)
			return
		}

		// Update the list
		cfg.Tools = newTools

		// Save the configuration
		if err := cfg.Save(); err != nil {
			fmt.Printf("Error saving configuration: %v\n", err)
			return
		}

		fmt.Printf("%s removed from your tools list\n", toolName)
	},
}

func init() {
	// Add tools subcommands
	cfgToolsCmd.AddCommand(cfgToolsAddCmd)
	cfgToolsCmd.AddCommand(cfgToolsRemoveCmd)

	// Add all subcommands to cfg command
	cfgCmd.AddCommand(cfgRegenerateCmd)
	cfgCmd.AddCommand(cfgToolsCmd)

	// Add cfg command to root
	rootCmd.AddCommand(cfgCmd)
}
