package main

import (
	"math/rand"
	"os"
	"strings"
)

var matrixSymbols = []rune("0123456789@#$%^&*+-=/\\|:<>?")

const (
	katakanaStart      = 0xFF66
	katakanaCount      = 56
	asciiStart         = 33
	asciiCount         = 94
	glitchChance       = 20
	spawnOffset        = 5
	matrixSymbolChance = 3
	matrixSymbolTotal  = 10
	tailDivisor        = 3
	envLang            = "LANG"
	envLcAll           = "LC_ALL"
	envTerm            = "TERM"
	termDumb           = "dumb"
	termLinux          = "linux"
	utf8Suffix         = "utf-8"
)

// Cell represents a single point on the screen.
type Cell struct {
	Char     rune
	Head     bool
	Tail     int
	Flash    bool
	T1Len    int
	T2Len    int
	CRTBleed bool
	CRTChar  rune
}

// Column tracks the state of a single falling matrix line.
type Column struct {
	YPos         int
	Tail         int
	PendingFlash bool
	T1Len        int
	T2Len        int
}

// Engine manages the state of the matrix effect.
type Engine struct {
	Width        int
	Height       int
	Grid         []Cell
	Columns      []Column
	Theme        Theme
	ASCIIOnly    bool
	StarMode     bool
	StarCount    int
	CRTMode      bool
	ShowHelp     bool
	GradientMode int
	UseTrueColor bool
	T1Percent    int
	T2Percent    int
}



// NewEngine creates and initializes a new matrix engine.
func NewEngine(w, h int, theme Theme, ascii bool) *Engine {
	if !ascii {
		lang := os.Getenv(envLang)
		lcAll := os.Getenv(envLcAll)
		term := os.Getenv(envTerm)
		if !strings.Contains(strings.ToLower(lang), utf8Suffix) &&
			!strings.Contains(strings.ToLower(lcAll), utf8Suffix) {
			ascii = true
		}
		if term == termDumb || term == termLinux {
			ascii = true
		}
	}
	e := &Engine{Theme: theme, ASCIIOnly: ascii, StarCount: DefaultStarCnt}
	e.resize(w, h)
	return e
}

// resize re-initializes the grid and columns for the given dimensions.
func (e *Engine) resize(w, h int) {
	if w <= 0 || h <= 0 {
		w, h = 1, 1
	}
	e.Width = w
	e.Height = h
	e.Grid = make([]Cell, w*h)
	e.Columns = make([]Column, w)
	for i := 0; i < w; i++ {
		e.initColumnLengths(i)
		e.Columns[i].YPos = rand.Intn(h) - h // start above screen
	}
}

// initColumnLengths assigns random total and gradient tail lengths.
func (e *Engine) initColumnLengths(x int) {
	minTail := e.Height / 2
	tailRange := e.Height - minTail
	e.Columns[x].Tail = rand.Intn(max(1, tailRange)) + minTail

	tMax := max(1, e.Columns[x].Tail/tailDivisor)
	e.Columns[x].T1Len = rand.Intn(tMax) + tailDivisor
	e.Columns[x].T2Len = rand.Intn(tMax) + tailDivisor
}

// randomRange returns a random integer in [min, max] inclusive.
func randomRange(min, max int) int {
	if min >= max {
		return min
	}
	// rand.Intn(n) is exclusive of n. We add an inclusive offset
	// to ensure the max value can actually be generated.
	inclusiveOffset := 1
	return rand.Intn(max-min+inclusiveOffset) + min
}

// randChar returns a random display character.
func (e *Engine) randChar() rune {
	if e.ASCIIOnly {
		return rune(randomRange(asciiStart, asciiStart + asciiCount - 1))
	}
	if rand.Intn(matrixSymbolTotal) < matrixSymbolChance {
		return matrixSymbols[rand.Intn(len(matrixSymbols))]
	}
	// Half-width katakana
	kEnd := katakanaStart + katakanaCount - 1
	return rune(randomRange(katakanaStart, kEnd))
}

// Step advances the matrix effect by one frame.
func (e *Engine) Step() {
	for x := 0; x < e.Width; x++ {
		e.Columns[x].YPos++
		y := e.Columns[x].YPos
		tail := e.Columns[x].Tail

		if y-tail > e.Height {
			e.Columns[x].YPos = randomRange(-spawnOffset, -1)
			e.initColumnLengths(x)
			e.Columns[x].PendingFlash = false
			continue
		}

		for cy := 0; cy < e.Height; cy++ {
			idx := x*e.Height + cy
			if cy == y {
				e.Grid[idx].Char = e.randChar()
				e.Grid[idx].Head = true
				e.Grid[idx].Tail = 0
				e.Grid[idx].Flash = e.Columns[x].PendingFlash
				e.Grid[idx].T1Len = e.Columns[x].T1Len
				e.Grid[idx].T2Len = e.Columns[x].T2Len
				e.Columns[x].PendingFlash = false
			} else if cy < y && cy >= y-tail {
				e.Grid[idx].Head = false
				e.Grid[idx].Tail = y - cy
				e.Grid[idx].Flash = false
				e.Grid[idx].T1Len = e.Columns[x].T1Len
				e.Grid[idx].T2Len = e.Columns[x].T2Len
				if rand.Intn(glitchChance) == 0 {
					e.Grid[idx].Char = e.randChar()
				}
			} else {
				e.Grid[idx].Char = ' '
				e.Grid[idx].Head = false
				e.Grid[idx].Tail = -1
				e.Grid[idx].Flash = false
				e.Grid[idx].T1Len = 0
				e.Grid[idx].T2Len = 0
			}
		}
	}

	if e.StarMode {
		e.generateStars()
	}

	for i := 0; i < len(e.Grid); i++ {
		e.Grid[i].CRTBleed = false
		e.Grid[i].CRTChar = ' '
	}

	if e.CRTMode {
		for x := 1; x < e.Width; x++ {
			for y := 0; y < e.Height; y++ {
				idx := x*e.Height + y
				leftIdx := (x-1)*e.Height + y
				if e.Grid[leftIdx].Flash {
					e.Grid[idx].CRTBleed = true
					e.Grid[idx].CRTChar = e.Grid[leftIdx].Char
				}
			}
		}
	}
}

// generateStars randomly selects cells to flash based on StarCount.
func (e *Engine) generateStars() {
	var candidates []int
	for i := 0; i < len(e.Grid); i++ {
		c := e.Grid[i]
		if c.Head || (c.Tail >= 0 && c.Tail < c.T1Len) {
			candidates = append(candidates, i)
		}
	}

	if len(candidates) > 0 {
		n := randomRange(0, e.StarCount)
		if n > len(candidates) {
			n = len(candidates)
		}
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
		for i := 0; i < n; i++ {
			e.Grid[candidates[i]].Flash = true
		}
	}
}

// TriggerFlash sets a random column to flash next tick.
func (e *Engine) TriggerFlash() {
	if e.Width == 0 {
		return
	}
	x := rand.Intn(e.Width)
	e.Columns[x].PendingFlash = true
}

// IncrementStarCount increases the star count up to maxCnt.
func (e *Engine) IncrementStarCount(maxCnt int) {
	if e.StarCount < maxCnt {
		e.StarCount++
	}
}

// DecrementStarCount decreases the star count down to minCnt.
func (e *Engine) DecrementStarCount(minCnt int) {
	if e.StarCount > minCnt {
		e.StarCount--
	}
}

// ToggleStarMode toggles the star background effect.
func (e *Engine) ToggleStarMode() {
	e.StarMode = !e.StarMode
}

// ToggleCRTMode toggles the CRT bleed effect.
func (e *Engine) ToggleCRTMode() {
	e.CRTMode = !e.CRTMode
}

// ToggleHelp toggles the help menu display.
func (e *Engine) ToggleHelp() {
	e.ShowHelp = !e.ShowHelp
}

// ToggleGradientMode cycles through the available gradient modes.
func (e *Engine) ToggleGradientMode(modeCount int) {
	e.GradientMode = cycleNext(e.GradientMode, modeCount)
}

// ToggleTrueColor toggles TrueColor vs Dithering rendering.
func (e *Engine) ToggleTrueColor() {
	e.UseTrueColor = !e.UseTrueColor
}
