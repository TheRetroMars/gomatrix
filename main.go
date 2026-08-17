package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// MinSpeed is the lowest allowed speed level.
// DefaultSpeed is the initial speed on startup.
const (
	MinSpeed       = 1
	MaxSpeed       = 10
	DefaultSpeed   = 8
	MinStarCount   = 1
	MaxStarCount   = 20
	DefaultStarCnt = 7

	DefaultT1Percent = 30
	DefaultT2Percent = 50

	ModeClassic    = 0
	ModeSmooth     = 1
	ModeVerySmooth = 2

	BufferSmooth  = 45.0
	MaxPercent    = 100.0
	MinRatio      = 0.0
	MaxRatio      = 1.0
	HelpMenuWidth = 32

	msgThemePrefix = "Theme: "
	msgCRTOff      = "CRT Effect: OFF"
	msgCRTOn       = "CRT Effect: ON"
	msgStarOff     = "Star: OFF"
	msgStarOn      = "Star: ON"
	msgStarMax     = "Star Max: %d"
	themeClassic   = "Classic"

	msgGradientPrefix = "Gradient: "
	msgTypeTrueColor  = "Render: TrueColor"
	msgTypeDithering  = "Render: Dithering"
	msgHelpWarning    = "Screen too small for help"
)

var helpText = []string{
	"  [h, ?]     Toggle Help      ",
	"  [Space]    Trigger Flash    ",
	"  [+, =]     Speed Up         ",
	"  [-, _]     Speed Down       ",
	"  [c]        Change Theme     ",
	"  [r]        Reset State      ",
	"  [g]        Toggle CRT Mode  ",
	"  [s]        Toggle Star Mode ",
	"  [[, ]]     Adjust Star Max  ",
	"  [m]        Gradient Mode    ",
	"  [t]        Render Engine    ",
	"  [q, Esc]   Quit             ",
}

// SpeedDelays maps speed levels to their corresponding tick intervals.
var SpeedDelays = []time.Duration{
	60000 * time.Millisecond, 45000 * time.Millisecond,
	30000 * time.Millisecond, 15000 * time.Millisecond,
	10000 * time.Millisecond, 5000 * time.Millisecond,
	2000 * time.Millisecond, 1000 * time.Millisecond,
	800 * time.Millisecond, 600 * time.Millisecond,
	450 * time.Millisecond, 300 * time.Millisecond,
	200 * time.Millisecond, 150 * time.Millisecond,
	100 * time.Millisecond, 75 * time.Millisecond,
	50 * time.Millisecond,
}

var themes = []string{"Classic", "AntiGravity", "QuantumGold", "Cyberpunk"}
var gradientModes = []string{"Classic", "Smooth", "VerySmooth"}

// applyTheme applies the selected theme to the configuration and engine.
func applyTheme(cfg *Config, engine *Engine, themeName string) {
	cfg.ColorTheme = themeName
	engine.Theme = getTheme(themeName)
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
) (forceRedraw bool, quit bool, spdChg bool, msg string) {
	if engine.ShowHelp {
		engine.ToggleHelp()
		return true, false, false, ""
	}
	switch ev.Key() {
	case tcell.KeyEscape, tcell.KeyCtrlC:
		return false, true, false, ""
	}
	switch ev.Rune() {
	case 'h', '?':
		engine.ToggleHelp()
		return true, false, false, ""
	case 'm':
		engine.ToggleGradientMode(len(gradientModes))
		return true, false, false, msgGradientPrefix + gradientModes[engine.GradientMode]
	case 't':
		engine.ToggleTrueColor()
		if engine.UseTrueColor {
			return true, false, false, msgTypeTrueColor
		}
		return true, false, false, msgTypeDithering
	case ' ':
		engine.TriggerFlash()
		return false, false, false, ""
	case 'q':
		return false, true, false, ""
	case '+', '=':
		adjustSpeed(cfg, ticker, 1)
		return true, false, true, ""
	case '-', '_':
		adjustSpeed(cfg, ticker, -1)
		return true, false, true, ""
	case 'c':
		*themeIdx = cycleNext(*themeIdx, len(themes))
		applyTheme(cfg, engine, themes[*themeIdx])
		return true, false, false, msgThemePrefix + themes[*themeIdx]
	case 'r':
		cfg.Speed = DefaultSpeed
		ticker.Reset(SpeedDelays[cfg.Speed-1])
		*themeIdx = 0
		applyTheme(cfg, engine, themeClassic)
		return true, false, true, ""
	case 'g':
		engine.ToggleCRTMode()
		if engine.CRTMode {
			return true, false, false, msgCRTOn
		}
		return true, false, false, msgCRTOff
	case 's':
		engine.ToggleStarMode()
		if engine.StarMode {
			return true, false, false, msgStarOn
		}
		return true, false, false, msgStarOff
	case ']':
		if !engine.StarMode {
			return false, false, false, ""
		}
		engine.IncrementStarCount(MaxStarCount)
		return true, false, false, fmt.Sprintf(msgStarMax, engine.StarCount)
	case '[':
		if !engine.StarMode {
			return false, false, false, ""
		}
		engine.DecrementStarCount(MinStarCount)
		return true, false, false, fmt.Sprintf(msgStarMax, engine.StarCount)
	}
	return false, false, false, ""
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

// formatFloat formats a float to max 3 decimal places.
func formatFloat(v float64) string {
	s := fmt.Sprintf("%.3f", v)
	s = strings.TrimRight(s, "0")
	return strings.TrimRight(s, ".")
}

// setOverlayMsg updates the overlay message and sets a 1-second timeout.
func setOverlayMsg(msg string, target *string, ch *<-chan time.Time) {
	*target = msg
	*ch = time.After(1 * time.Second)
}

// cycleNext returns the next index in a circular array of size count.
func cycleNext(current, count int) int {
	if count <= 0 {
		return current
	}
	stepForward := 1
	return (current + stepForward) % count
}



// main initializes the configuration, sets up the terminal screen,
// and starts the main event loop for rendering the matrix effect.
func main() {
	configPath := parseFlags()

	cfg, err := loadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

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
	engine.StarMode = cfg.StarMode
	if cfg.StarCount < MinStarCount {
		cfg.StarCount = MinStarCount
	} else if cfg.StarCount > MaxStarCount {
		cfg.StarCount = MaxStarCount
	}
	engine.StarCount = cfg.StarCount
	engine.CRTMode = cfg.CRTMode
	engine.T1Percent = cfg.T1Percent
	engine.T2Percent = cfg.T2Percent
	engine.UseTrueColor = cfg.TrueColor
	engine.GradientMode = cfg.GradientMode
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

	var overlayMsg string
	var clearMsgCh <-chan time.Time

	go func() {
		for {
			eventCh <- screen.PollEvent()
		}
	}()

	for {
		draw := false
		select {
		case ev := <-eventCh:
			switch ev := ev.(type) {
			case *tcell.EventKey:
				force, quit, spd, msg := handleKeyEvent(
					ev, cfg, ticker, engine, &themeIdx)
				if quit {
					return
				}
				if spd {
					delay := float64(SpeedDelays[cfg.Speed-1].Milliseconds())
					msg = fmt.Sprintf("%s Sec/Frame", formatFloat(delay/1000.0))
				}
				if msg != "" {
					setOverlayMsg(msg, &overlayMsg, &clearMsgCh)
				}
				if force || msg != "" {
					forceRedraw = true
					draw = true
				}
			case *tcell.EventResize:
				nw, nh := screen.Size()
				engine.resize(nw, nh)
				oldGrid = make([]Cell, nw*nh)
				forceRedraw = true
				draw = true
				screen.Sync()
			}
		case <-clearMsgCh:
			overlayMsg = ""
			forceRedraw = true
			clearMsgCh = nil
			draw = true
		case <-ticker.C:
			engine.Step()
			draw = true
		}

		if draw {
			if renderScreen(screen, engine, oldGrid, forceRedraw, overlayMsg) {
				screen.Show()
				forceRedraw = false
			}
		}
	}
}
