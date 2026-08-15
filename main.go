package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

// MinSpeed is the lowest allowed speed level.
// MaxSpeed is the highest allowed speed level.
// DefaultSpeed is the initial speed on startup.
// ColorTailT1 is the tail age threshold for the first trail color tier.
// ColorTailT2 is the tail age threshold for the second trail color tier.
const (
	MinSpeed     = 1
	MaxSpeed     = 10
	DefaultSpeed = 5
	ColorTailT1  = 5
	ColorTailT2  = 10
)

// SpeedDelays maps speed levels (1-10) to their corresponding tick intervals.
var SpeedDelays = []time.Duration{
	1000 * time.Millisecond, 800 * time.Millisecond, 600 * time.Millisecond,
	450 * time.Millisecond, 300 * time.Millisecond, 200 * time.Millisecond,
	150 * time.Millisecond, 100 * time.Millisecond, 75 * time.Millisecond,
	50 * time.Millisecond,
}

var themes = []string{"Classic", "AntiGravity", "QuantumGold", "Cyberpunk"}

// clampSpeed constrains speed to
// [MinSpeed, MaxSpeed].
func clampSpeed(s int) int {
	if s < MinSpeed {
		return MinSpeed
	}
	if s > MaxSpeed {
		return MaxSpeed
	}
	return s
}

// adjustSpeed adjusts cfg.Speed by delta,
// clamping to [MinSpeed, MaxSpeed].
func adjustSpeed(
	cfg *Config,
	ticker *time.Ticker,
	delta int,
) {
	cfg.Speed = clampSpeed(cfg.Speed + delta)
	ticker.Reset(SpeedDelays[cfg.Speed-1])
}

// handleKeyEvent processes a key press and
// returns whether a full redraw is needed
// and whether to quit.
func handleKeyEvent(
	ev *tcell.EventKey,
	cfg *Config,
	ticker *time.Ticker,
	engine *Engine,
	themeIdx *int,
) (forceRedraw bool, quit bool) {
	switch ev.Key() {
	case tcell.KeyEscape, tcell.KeyCtrlC:
		return false, true
	}
	switch ev.Rune() {
	case 'q':
		return false, true
	case '+', '=':
		adjustSpeed(cfg, ticker, 1)
	case '-', '_':
		adjustSpeed(cfg, ticker, -1)
	case 'c':
		*themeIdx = (*themeIdx + 1) % len(themes)
		cfg.ColorTheme = themes[*themeIdx]
		engine.Theme = getTheme(cfg.ColorTheme)
		return true, false
	case 'r':
		cfg.Speed = DefaultSpeed
		ticker.Reset(SpeedDelays[cfg.Speed-1])
		*themeIdx = 0
		cfg.ColorTheme = "Classic"
		engine.Theme = getTheme(cfg.ColorTheme)
		return true, false
	}
	return false, false
}

// renderScreen draws changed cells to the
// screen and returns true if any cell was
// updated.
func renderScreen(
	screen tcell.Screen,
	engine *Engine,
	oldGrid []Cell,
	forceRedraw bool,
) bool {
	rendered := false
	for x := 0; x < engine.Width; x++ {
		for y := 0; y < engine.Height; y++ {
			idx := x*engine.Height + y
			if idx >= len(engine.Grid) || idx >= len(oldGrid) {
				continue // Guard against out-of-bounds access after resize.
			}
			c := engine.Grid[idx]
			if !forceRedraw && oldGrid[idx] == c {
				continue
			}
			oldGrid[idx] = c
			rendered = true

			st := tcell.StyleDefault.Background(engine.Theme.Bg)
			if c.Char == ' ' {
				screen.SetContent(x, y, ' ', nil, st)
				continue
			}
			if c.Head {
				st = st.Foreground(engine.Theme.Head)
			} else if c.Tail < ColorTailT1 {
				st = st.Foreground(engine.Theme.T1)
			} else if c.Tail < ColorTailT2 {
				st = st.Foreground(engine.Theme.T2)
			} else {
				st = st.Foreground(engine.Theme.T3)
			}
			screen.SetContent(x, y, c.Char, nil, st)
		}
	}
	return rendered
}

// parseFlags configures the CLI flags, sets a
// custom usage message, and parses arguments.
// It returns the config path. If the version
// flag is set, it prints the version and exits.
func parseFlags() string {
	var showVer bool
	flag.BoolVar(&showVer, "version", false, "Print version and exit")
	flag.BoolVar(&showVer, "v", false, "Print version (shorthand)")
	configPath := flag.String("config", "", "Path to config.json")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gomatrix %s\n\n", Version)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if showVer {
		fmt.Printf("gomatrix %s\n", Version)
		os.Exit(0)
	}

	return *configPath
}

func main() {
	configPath := parseFlags()

	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	cfg.Speed = clampSpeed(cfg.Speed)

	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}
	defer screen.Fini()

	w, h := screen.Size()
	engine := NewEngine(w, h, getTheme(cfg.ColorTheme), cfg.ASCIIOnly)
	themeIdx := 0
	for i, t := range themes {
		if t == cfg.ColorTheme {
			themeIdx = i
			break
		}
	}

	oldGrid := make([]Cell, w*h)
	forceRedraw := true

	ticker := time.NewTicker(SpeedDelays[cfg.Speed-1])
	defer ticker.Stop()
	eventCh := make(chan tcell.Event)

	go func() {
		for {
			eventCh <- screen.PollEvent()
		}
	}()

	for {
		select {
		case ev := <-eventCh:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				force, quit := handleKeyEvent(ev, cfg, ticker, engine, &themeIdx)
				if quit {
					return
				}
				if force {
					forceRedraw = true
				}
			case *tcell.EventResize:
				nw, nh := screen.Size()
				engine.resize(nw, nh)
				oldGrid = make([]Cell, nw*nh)
				forceRedraw = true
				screen.Sync()
			}
		case <-ticker.C:
			engine.Step()
			if renderScreen(screen, engine, oldGrid, forceRedraw) {
				screen.Show()
				forceRedraw = false
			}
		}
	}
}
