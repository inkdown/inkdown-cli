package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

type Config struct {
	Token string `json:"token,omitempty"`
	Email string `json:"email,omitempty"`
}

func ConfigPath() string {
	return filepath.Join(xdg.ConfigHome, "ink", "config.json")
}

func Load() (*Config, error) {
	configPath := ConfigPath()

	fmt.Printf("DEBUG: Config directory: %s \n", configPath)

	configDir := filepath.Dir(configPath)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		cfg := &Config{}
		if err := cfg.Save(); err != nil {
			return nil, err
		}
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return &Config{}, nil
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	configPath := ConfigPath()
	configDir := filepath.Dir(configPath)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

func (c *Config) IsAuthenticated() bool {
	return c.Token != "" && c.Email != ""
}

func (c *Config) SetToken(token string) error {
	c.Token = token
	return c.Save()
}

func (c *Config) ClearToken() error {
	c.Token = ""
	return c.Save()
}

func (c *Config) SetEmail(email string) error {
	c.Email = email
	return c.Save()
}
