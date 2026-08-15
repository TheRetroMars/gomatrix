package main

import "testing"

// TestNewEngine verifies engine initialization
// logic including fallback conditions.
func TestNewEngine(t *testing.T) {
	theme := getTheme("Classic")

	t.Run("Dimensions", func(t *testing.T) {
		e := NewEngine(80, 24, theme, false)
		if e.Width != 80 || e.Height != 24 {
			t.Errorf("expected 80x24, got %dx%d", e.Width, e.Height)
		}
		if len(e.Grid) != 80*24 {
			t.Errorf("expected grid length %d, got %d", 80*24, len(e.Grid))
		}
	})

	t.Run("ASCIIFallback", func(t *testing.T) {
		t.Setenv("TERM", "dumb")
		e2 := NewEngine(80, 24, theme, false)
		if !e2.ASCIIOnly {
			t.Errorf("expected ASCII fallback for dumb terminal")
		}
	})
}
