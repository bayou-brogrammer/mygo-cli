package main

import (
	"fmt"
	"os"

	"github.com/bayou-brogrammer/mygo/internal/config"
	"github.com/bayou-brogrammer/mygo/internal/logger"
	"github.com/bayou-brogrammer/mygo/internal/ui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	verbose  bool
	logLevel string
)

var rootCmd = &cobra.Command{
	Use:   "milo",
	Short: "A CLI tool for managing development environment",
	Long: `Milo CLI is a comprehensive tool for managing your development environment,
including GitHub repositories, dotfiles, and system configuration via chezmoi.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize logger with appropriate level
		loggerLevel := logger.ParseLevel(logLevel)
		if verbose {
			loggerLevel = logger.LevelDebug
		}

		logger.Init(loggerLevel)
		logger.Debug("Logger initialized with level: %s", loggerLevel)

		// Initialize configuration
		_, err := config.Init()
		if err != nil {
			logger.Fatal("Error initializing config: %v", err)
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/milo/config.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose output")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "set log level (debug, info, warn, error, fatal)")

	// Register commands - these are defined in their respective files
	// and will be automatically registered when those files are imported

	// Add version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			ui.PrintInfo("Milo CLI v0.1.0")
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

		// Search config in home directory with name ".milo" (without extension)
		viper.AddConfigPath(fmt.Sprintf("%s/.config/milo", home))
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
	// Ensure logger is closed when program exits
	defer logger.Close()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
