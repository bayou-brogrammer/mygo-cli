package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// DefaultToolsFile defines the default list of tools
const DefaultToolsFile = "data/default_tools.yml"

// Config holds the application configuration
type Config struct {
	// General configuration
	ConfigDir string

	// Repository configuration
	ReposDir     string
	TrackedRepos map[string]Repository

	// Dotfiles configuration
	DotfilesRepo string
	DotfilesDir  string

	// Chezmoi configuration
	ChezmoiDir string

	// System configuration
	Tools []string
}

// Repository represents a tracked GitHub repository
type Repository struct {
	URL         string
	Path        string
	Description string
	LastUpdated string
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}

	// Read default tools from file
	DefaultTools, err := readDefaultTools(DefaultToolsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading default tools: %v\n", err)
		os.Exit(1)
	}

	return &Config{
		ConfigDir: filepath.Join(homeDir, ".config", "devenv"),
		// ReposDir:     filepath.Join(homeDir, "Projects"),
		// TrackedRepos: make(map[string]Repository),
		// DotfilesRepo: "",
		// DotfilesDir:  filepath.Join(homeDir, ".dotfiles"),
		// ChezmoiDir:   filepath.Join(homeDir, ".local", "share", "chezmoi"),
		Tools: DefaultTools,
	}
}

var cfg *Config

// Init initializes the configuration
func Init() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = DefaultConfig()

	// Ensure config directory exists
	if err := os.MkdirAll(cfg.ConfigDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set up Viper for configuration
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(cfg.ConfigDir)

	// Set defaults
	// viper.SetDefault("repos_dir", cfg.ReposDir)
	// viper.SetDefault("dotfiles_dir", cfg.DotfilesDir)
	// viper.SetDefault("chezmoi_dir", cfg.ChezmoiDir)
	viper.SetDefault("tools", cfg.Tools)

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create it
			if err := viper.SafeWriteConfig(); err != nil {
				return nil, fmt.Errorf("failed to create default config file: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Load configuration into struct
	// cfg.ReposDir = viper.GetString("repos_dir")
	// cfg.DotfilesRepo = viper.GetString("dotfiles_repo")
	// cfg.DotfilesDir = viper.GetString("dotfiles_dir")
	// cfg.ChezmoiDir = viper.GetString("chezmoi_dir")
	cfg.Tools = viper.GetStringSlice("tools")

	// Load tracked repositories
	reposFile := filepath.Join(cfg.ConfigDir, "repos.yaml")
	if _, err := os.Stat(reposFile); err == nil {
		reposViper := viper.New()
		reposViper.SetConfigFile(reposFile)
		if err := reposViper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read repos file: %w", err)
		}

		if err := reposViper.UnmarshalKey("repos", &cfg.TrackedRepos); err != nil {
			return nil, fmt.Errorf("failed to unmarshal repos: %w", err)
		}
	}

	return cfg, nil
}

// Save saves the current configuration
func (c *Config) Save() error {
	// Save main configuration
	// viper.Set("repos_dir", c.ReposDir)
	// viper.Set("dotfiles_repo", c.DotfilesRepo)
	// viper.Set("dotfiles_dir", c.DotfilesDir)
	// viper.Set("chezmoi_dir", c.ChezmoiDir)
	viper.Set("tools", c.Tools)

	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Save tracked repositories
	reposViper := viper.New()
	reposViper.SetConfigFile(filepath.Join(c.ConfigDir, "repos.yaml"))
	reposViper.Set("repos", c.TrackedRepos)

	if err := reposViper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write repos file: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func GetConfig() (*Config, error) {
	if cfg == nil {
		return Init()
	}
	return cfg, nil
}

// readDefaultTools reads the default tools from a YAML file
func readDefaultTools(filename string) ([]string, error) {
	var tools []string

	viper.SetConfigFile(filename)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read default tools file: %w", err)
	}

	if err := viper.UnmarshalKey("tools", &tools); err != nil {
		return nil, fmt.Errorf("failed to unmarshal default tools: %w", err)
	}

	return tools, nil
}
