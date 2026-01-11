package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/julianfbeck/bring-cli/internal/api"
	"gopkg.in/yaml.v3"
)

const (
	configDir  = "bring-cli"
	configFile = "config.yaml"
)

// Config holds the CLI configuration.
type Config struct {
	Credentials *api.Credentials `yaml:"credentials,omitempty"`
}

// GetConfigPath returns the path to the config file.
func GetConfigPath() (string, error) {
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("getting home directory: %w", err)
		}
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, configDir, configFile), nil
}

// Load loads the configuration from disk.
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to disk.
func Save(cfg *Config) error {
	path, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}

	return nil
}

// SaveCredentials saves credentials to the config.
func SaveCredentials(creds *api.Credentials) error {
	cfg, err := Load()
	if err != nil {
		cfg = &Config{}
	}
	cfg.Credentials = creds
	return Save(cfg)
}

// GetCredentials returns stored credentials.
func GetCredentials() (*api.Credentials, error) {
	cfg, err := Load()
	if err != nil {
		return nil, err
	}
	return cfg.Credentials, nil
}

// ClearCredentials removes stored credentials.
func ClearCredentials() error {
	cfg, err := Load()
	if err != nil {
		return err
	}
	cfg.Credentials = nil
	return Save(cfg)
}

// SetDefaultList sets the default list UUID.
func SetDefaultList(listUUID string) error {
	cfg, err := Load()
	if err != nil {
		cfg = &Config{}
	}
	if cfg.Credentials == nil {
		cfg.Credentials = &api.Credentials{}
	}
	cfg.Credentials.DefaultList = listUUID
	return Save(cfg)
}

// GetDefaultList returns the stored default list UUID.
func GetDefaultList() string {
	cfg, err := Load()
	if err != nil || cfg.Credentials == nil {
		return ""
	}
	return cfg.Credentials.DefaultList
}
