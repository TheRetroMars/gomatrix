package main

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
)

// TestHandleKeyEvent verifies key event logic.
func TestHandleKeyEvent(t *testing.T) {
	cfg := &Config{Speed: 5, ColorTheme: "Classic"}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	engine := &Engine{}
	themeIdx := 0

	ev := tcell.NewEventKey(tcell.KeyRune, 'g', tcell.ModNone)
	forceRedraw, quit, _, _ := handleKeyEvent(ev, cfg, ticker, engine, &themeIdx)

	if !forceRedraw {
		t.Error("expected redraw on 'g'")
	}
	if quit {
		t.Error("expected no quit on 'g'")
	}
	if !engine.CRTMode {
		t.Error("expected CRT mode toggled")
	}

	evQuit := tcell.NewEventKey(tcell.KeyRune, 'q', tcell.ModNone)
	_, quit, _, _ = handleKeyEvent(evQuit, cfg, ticker, engine, &themeIdx)
	if !quit {
		t.Error("expected quit on 'q'")
	}
}

// TestAdjustSpeed verifies speed boundary constraints.
func TestAdjustSpeed(t *testing.T) {
	cfg := &Config{Speed: MinSpeed}
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	adjustSpeed(cfg, ticker, -1)
	if cfg.Speed != MinSpeed {
		t.Errorf("expected MinSpeed, got %d", cfg.Speed)
	}

	adjustSpeed(cfg, ticker, 1)
	if cfg.Speed != MinSpeed+1 {
		t.Errorf("expected %d, got %d", MinSpeed+1, cfg.Speed)
	}
}
