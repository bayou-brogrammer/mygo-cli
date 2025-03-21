package main

import (
	"fmt"
	"os"

	"github.com/bayou-brogrammer/mygo/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool
)

var rootCmd = &cobra.Command{
	Use:   "devenv",
	Short: "A CLI tool for managing development environment",
	Long: `DevEnv CLI is a comprehensive tool for managing your development environment,
including GitHub repositories, dotfiles, and system configuration via chezmoi.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize configuration
		_, err := config.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommands are provided, print help
		cmd.Help()
	},
}

func init() {
	// Initialize cobra
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/devenv/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")

	// Register commands - these are defined in their respective files
	// and will be automatically registered when those files are imported

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("DevEnv CLI v0.1.0")
		},
	})
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
			os.Exit(1)
		}

		// Search config in home directory with name ".devenv" (without extension)
		viper.AddConfigPath(fmt.Sprintf("%s/.config/devenv", home))
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// Read in environment variables that match
	viper.AutomaticEnv()

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil && verbose {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
