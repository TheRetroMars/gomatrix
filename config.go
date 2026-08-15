package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// Config holds the runtime configuration loaded
// from a JSON file or built-in defaults.
type Config struct {
	Speed      int    `json:"speed"`
	ColorTheme string `json:"color_theme"`
	ASCIIOnly  bool   `json:"ascii_only"`
}

// loadConfig loads the configuration from the given file path or the default
// fallback locations.
func loadConfig(cliPath string) (*Config, error) {
	cfg := &Config{Speed: DefaultSpeed, ColorTheme: "Classic", ASCIIOnly: false}
	paths := []string{cliPath}
	home, err := os.UserHomeDir()
	if err == nil {
		p := filepath.Join(home, ".config", "gomatrix", "config.json")
		paths = append(paths, p)
	} else {
		log.Printf("Warning: failed to get user home directory: %v", err)
	}
	paths = append(paths, "./gomatrix.json")
	for _, p := range paths {
		if p == "" {
			continue
		}
		if data, err := os.ReadFile(p); err == nil {
			err = json.Unmarshal(data, cfg)
			return cfg, err
		}
	}
	return cfg, nil
}
