package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	URL        string   `mapstructure:"url"`
	APIKey     string   `mapstructure:"api_key"`
	Workspace  string   `mapstructure:"workspace"`
	Stream     bool     `mapstructure:"stream"`
	Mode       string   `mapstructure:"mode"`
	Priorities []string `mapstructure:"priorities"`
}

func Dir() (string, error) {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "anycli"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("finding home directory: %w", err)
	}
	return filepath.Join(home, ".config", "anycli"), nil
}

func Load() (*Config, error) {
	v := viper.New()

	dir, err := Dir()
	if err != nil {
		return nil, err
	}

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(dir)
	v.AddConfigPath(".")

	v.SetEnvPrefix("ANYCLI")
	v.AutomaticEnv()

	v.SetDefault("url", "http://localhost:3001")
	v.SetDefault("workspace", "vault")
	v.SetDefault("stream", true)
	v.SetDefault("mode", "chat")
	v.SetDefault("priorities", []string{
		"health (sleep >= 7h, exercise 4x/week)",
		"deep work blocks",
		"obligations (meetings, comms)",
	})

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("reading config: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	return &cfg, nil
}

func Write(url, apiKey, workspace string) (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("creating config directory: %w", err)
	}

	path := filepath.Join(dir, "config.yaml")
	content := fmt.Sprintf(`# anycli configuration
# Docs: https://github.com/stepinski/anycli

url: %q
api_key: %q
workspace: %q
stream: true
mode: chat

priorities:
  - "health (sleep >= 7h, exercise 4x/week)"
  - "deep work blocks"
  - "obligations (meetings, comms)"
`, url, apiKey, workspace)

	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		return "", fmt.Errorf("writing config: %w", err)
	}

	return path, nil
}
