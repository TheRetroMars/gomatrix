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

// TestClampSpeed verifies speed constraints against boundary values.
func TestClampSpeed(t *testing.T) {
	maxSpd := len(SpeedDelays)
	cases := []struct {
		input    int
		expected int
	}{
		{0, MinSpeed},
		{maxSpd + 1, maxSpd},
		{5, 5},
	}
	for _, tc := range cases {
		if got := clampSpeed(tc.input); got != tc.expected {
			t.Errorf("clampSpeed(%d) = %d; want %d", tc.input, got, tc.expected)
		}
	}
}

// TestStarCountBounds verifies star count clamping.
func TestStarCountBounds(t *testing.T) {
	e := &Engine{StarCount: 7}
	e.IncrementStarCount(20)
	if e.StarCount != 8 {
		t.Errorf("expected 8, got %d", e.StarCount)
	}
	e.StarCount = 20
	e.IncrementStarCount(20)
	if e.StarCount != 20 {
		t.Errorf("expected 20, got %d", e.StarCount)
	}

	e.StarCount = 2
	e.DecrementStarCount(1)
	if e.StarCount != 1 {
		t.Errorf("expected 1, got %d", e.StarCount)
	}
	e.DecrementStarCount(1)
	if e.StarCount != 1 {
		t.Errorf("expected 1, got %d", e.StarCount)
	}
}

// TestToggleModes verifies toggling logic.
func TestToggleModes(t *testing.T) {
	e := &Engine{}
	e.ToggleStarMode()
	if !e.StarMode {
		t.Error("expected StarMode true")
	}
	e.ToggleCRTMode()
	if !e.CRTMode {
		t.Error("expected CRTMode true")
	}
}
