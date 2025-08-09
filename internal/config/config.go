package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	Jira JiraConfig `mapstructure:"jira"`
}

// JiraConfig represents Jira-specific configuration
type JiraConfig struct {
	Project string `mapstructure:"project"`
}

// Load loads configuration from .jitt.yaml file
func Load() (*Config, error) {
	viper.SetConfigName(".jitt")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("jira.project", "")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found")
		}
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// Exists checks if the config file exists
func Exists() bool {
	_, err := os.Stat(".jitt.yaml")
	return err == nil
}

// Create creates a new config file with the given project
func Create(project string) error {
	viper.Set("jira.project", project)

	return viper.WriteConfigAs(".jitt.yaml")
}
