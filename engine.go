package main

import (
	"math/rand"
	"os"
	"strings"
)

var matrixSymbols = []rune("0123456789@#$%^&*+-=/\\|:<>?")

const (
	katakanaStart = 0xFF66
	katakanaCount = 56
)

// Cell represents a single point on the screen.
type Cell struct {
	Char rune
	Head bool
	Tail int
}

// Column tracks the state of a single falling matrix line.
type Column struct {
	YPos int
	Tail int
}

// Engine manages the state of the matrix effect.
type Engine struct {
	Width     int
	Height    int
	Grid      []Cell
	Columns   []Column
	Theme     Theme
	ASCIIOnly bool
}

// NewEngine initializes a new engine with the given dimensions.
func NewEngine(w, h int, theme Theme, ascii bool) *Engine {
	if !ascii {
		lang := os.Getenv("LANG")
		lcAll := os.Getenv("LC_ALL")
		term := os.Getenv("TERM")
		if !strings.Contains(strings.ToLower(lang), "utf-8") &&
			!strings.Contains(strings.ToLower(lcAll), "utf-8") {
			ascii = true
		}
		if term == "dumb" || term == "linux" {
			ascii = true
		}
	}
	e := &Engine{Theme: theme, ASCIIOnly: ascii}
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
		e.Columns[i].Tail = rand.Intn(max(1, h/2)) + 5
		e.Columns[i].YPos = rand.Intn(h) - h // start above screen
	}
}

// randChar returns a random display character.
func (e *Engine) randChar() rune {
	if e.ASCIIOnly {
		return rune(rand.Intn(94) + 33) // printable ASCII
	}
	if rand.Intn(10) < 3 {
		return matrixSymbols[rand.Intn(len(matrixSymbols))]
	}
	return rune(rand.Intn(katakanaCount) + katakanaStart) // Half-width katakana
}

// Step advances the matrix effect by one frame.
func (e *Engine) Step() {
	for x := 0; x < e.Width; x++ {
		e.Columns[x].YPos++
		y := e.Columns[x].YPos
		tail := e.Columns[x].Tail

		if y-tail > e.Height {
			e.Columns[x].YPos = rand.Intn(5) - 5
			e.Columns[x].Tail = rand.Intn(max(1, e.Height/2)) + 5
			continue
		}

		for cy := 0; cy < e.Height; cy++ {
			idx := x*e.Height + cy
			if cy == y {
				e.Grid[idx].Char = e.randChar()
				e.Grid[idx].Head = true
				e.Grid[idx].Tail = 0
			} else if cy < y && cy >= y-tail {
				e.Grid[idx].Head = false
				e.Grid[idx].Tail = y - cy
				if rand.Intn(20) == 0 {
					e.Grid[idx].Char = e.randChar()
				}
			} else {
				e.Grid[idx].Char = ' '
				e.Grid[idx].Head = false
				e.Grid[idx].Tail = -1
			}
		}
	}
}
