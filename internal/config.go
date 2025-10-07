package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

type Source struct {
	Name   string `toml:"-"`
	Path   string `toml:"path"`
	Branch string `toml:"branch"`
	Remote string `toml:"remote"`
}

type Settings struct {
	BranchPrefix  string `toml:"branch_prefix"`
	DefaultSource string `toml:"default_source"`
}

type Config struct {
	Sources  map[string]Source `toml:"sources"`
	Settings Settings          `toml:"settings"`
}

func LoadConfig(configDir string) (*Config, error) {
	configPath := filepath.Join(configDir, "patch-review.config.toml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("no config at '%s'", configPath)
	}

	var config Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Set names in sources
	for name, source := range config.Sources {
		source.Name = name
		config.Sources[name] = source
	}

	return &config, nil
}

func (c *Config) GetSource(alias string) (*Source, error) {
	source, exists := c.Sources[alias]
	if !exists {
		return nil, fmt.Errorf("source not found: %s", alias)
	}
	return &source, nil
}

func (c *Config) GetDefaultSource() string {
	return c.Settings.DefaultSource
}

func (c *Config) GetBranchPrefix() string {
	if c.Settings.BranchPrefix == "" {
		return "patch-review"
	}
	return c.Settings.BranchPrefix
}

func ExpandPathString(pathStr string) string {
	if strings.HasPrefix(pathStr, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return pathStr
		}
		return filepath.Join(home, pathStr[2:])
	}
	return pathStr
}