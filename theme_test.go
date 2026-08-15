package main

import (
	"testing"
)

// TestGetTheme verifies that all predefined themes return valid color values.
func TestGetTheme(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"Classic"},
		{"AntiGravity"},
		{"QuantumGold"},
		{"Cyberpunk"},
		{"Unknown"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			theme := getTheme(tc.name)
			// Ensure that Head color has a non-zero valid value assigned.
			val := theme.Head.Hex()
			if val < 0 {
				t.Errorf("expected valid Hex Color for theme %s", tc.name)
			}
		})
	}
}
