package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadConfig_Defaults verifies that loadConfig returns the default values
// when no config file is found.
func TestLoadConfig_Defaults(t *testing.T) {
	cfg, err := loadConfig("")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Speed != 5 {
		t.Errorf("expected default speed 5, got %d", cfg.Speed)
	}
	if cfg.ColorTheme != "Classic" {
		t.Errorf("expected default Classic, got %s", cfg.ColorTheme)
	}
}

// TestLoadConfig_JSONFile verifies that loadConfig parses a valid JSON file
// successfully.
func TestLoadConfig_JSONFile(t *testing.T) {
	content := []byte(`{"speed":5, "color_theme":"Cyberpunk"}`)
	tmp := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(tmp, content, 0644); err != nil {
		t.Fatal(err)
	}
	cfg, err := loadConfig(tmp)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if cfg.Speed != 5 {
		t.Errorf("expected speed 5, got %d", cfg.Speed)
	}
	if cfg.ColorTheme != "Cyberpunk" {
		t.Errorf("expected Cyberpunk, got %s", cfg.ColorTheme)
	}
}
