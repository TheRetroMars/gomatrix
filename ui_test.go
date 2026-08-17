package main

import (
	"testing"
	"github.com/gdamore/tcell/v2"
)

// TestRenderScreen verifies that rendering a frame
// successfully forces a redraw when the flag is set.
func TestRenderScreen(t *testing.T) {
	screen := tcell.NewSimulationScreen("")
	screen.Init()
	engine := NewEngine(10, 10, getTheme("Classic"), false)
	oldGrid := make([]Cell, 100)

	rendered := renderScreen(screen, engine, oldGrid, true, "")
	if !rendered {
		t.Errorf("expected renderScreen to return true on first call")
	}

	rendered = renderScreen(screen, engine, oldGrid, false, "")
	if rendered {
		t.Errorf("expected renderScreen to return false when nothing changed")
	}
}

// TestGradientStyle verifies the gradient logic.
func TestGradientStyle(t *testing.T) {
	e := NewEngine(10, 10, getTheme("Classic"), false)

	// Test Case A: Proper Config Assignment (T1 zone)
	e.T1Percent = DefaultT1Percent
	e.T2Percent = DefaultT2Percent
	e.GradientMode = ModeClassic // Ensure Classic mode

	// tailDist = 1, totalTail = 10 -> pct = 10.0 <= T1 (15)
	st1 := getStyle(e, 1, 10)
	fg1, _, _ := st1.Decompose()
	if fg1 != e.Theme.T1 {
		t.Errorf("Expected T1 color for 10%% tail, got %v", fg1)
	}

	// Test Case B: Boundary Check (T2 zone)
	// tailDist = 3, totalTail = 10 -> pct = 30.0 <= T2 (50)
	st2 := getStyle(e, 3, 10)
	fg2, _, _ := st2.Decompose()
	if fg2 != e.Theme.T2 {
		t.Errorf("Expected T2 color for 30%% tail, got %v", fg2)
	}

	// Test Case C: Bug Simulation (Zero Values)
	e.T1Percent = 0
	e.T2Percent = 0

	// tailDist = 1, totalTail = 10 -> pct = 10.0 > T2 (0) -> goes to T3 block
	stBug := getStyle(e, 1, 10)
	fgBug, _, _ := stBug.Decompose()

	// If the bug exists (t1=0, t2=0), it defaults to T3/Bg processing
	if fgBug == e.Theme.T1 {
		t.Errorf("Simulated bug failed, unexpectedly got T1 color")
	}
}
