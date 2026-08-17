package main

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// clampFloat restricts a float64 to the range [min, max].
func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// half returns half of the given float64 value.
func half(v float64) float64 {
	return v / 2.0
}

// lerpColor interpolates between two colors based on p.
func lerpColor(c1, c2 tcell.Color, p float64) tcell.Color {
	p = clampFloat(p, MinRatio, MaxRatio)
	r1, g1, b1 := c1.RGB()
	r2, g2, b2 := c2.RGB()
	nr := int32(float64(r1)*(1-p) + float64(r2)*p)
	ng := int32(float64(g1)*(1-p) + float64(g2)*p)
	nb := int32(float64(b1)*(1-p) + float64(b2)*p)
	return tcell.NewRGBColor(nr, ng, nb)
}

// getStyle returns the styling for a given tail distance.
func getStyle(e *Engine, tailDist, totalTail int) tcell.Style {
	st := tcell.StyleDefault.Background(e.Theme.Bg)

	if tailDist < 0 || totalTail <= 0 {
		return st
	}

	pct := float64(tailDist) / float64(totalTail) * MaxPercent
	t1 := float64(e.T1Percent)
	t2 := float64(e.T2Percent)

	var c1, c2 tcell.Color
	var p float64

	// Determine segment and anchor colors
	if pct <= t1 {
		c1, c2 = e.Theme.Head, e.Theme.T1
		if t1 > 0 {
			p = pct / t1
		}
	} else if pct <= t2 {
		c1, c2 = e.Theme.T1, e.Theme.T2
		if t2 > t1 {
			p = (pct - t1) / (t2 - t1)
		}
	} else {
		c1, c2 = e.Theme.T2, e.Theme.Bg
		if MaxPercent > t2 {
			p = (pct - t2) / (MaxPercent - t2)
		}
	}

	// Apply gradient mode logic
	if e.GradientMode == ModeClassic {
		if pct < t1 {
			return st.Foreground(e.Theme.T1)
		} else if pct < t2 {
			return st.Foreground(e.Theme.T2)
		}
		// Fade to Bg in classic mode as well per instruction
		c1, c2 = e.Theme.T3, e.Theme.Bg
		if MaxPercent > t2 {
			p = (pct - t2) / (MaxPercent - t2)
		}
		if e.UseTrueColor {
			return st.Foreground(lerpColor(c1, c2, p))
		}
		if rand.Float64() > p {
			return st.Foreground(c1)
		}
		return st.Foreground(c2)
	}

	if e.GradientMode == ModeSmooth {
		buffer := BufferSmooth
		halfBuf := half(buffer)

		if pct >= t1-halfBuf && pct <= t1+halfBuf {
			c1, c2 = e.Theme.T1, e.Theme.T2
			p = (pct - (t1 - halfBuf)) / buffer
		} else if pct >= t2-halfBuf && pct <= t2+halfBuf {
			c1, c2 = e.Theme.T2, e.Theme.T3
			p = (pct - (t2 - halfBuf)) / buffer
		} else if pct < t1-halfBuf {
			c1, c2 = e.Theme.T1, e.Theme.T1
			p = 0
		} else if pct > t1+halfBuf && pct < t2-halfBuf {
			c1, c2 = e.Theme.T2, e.Theme.T2
			p = 0
		} else if pct > t2+halfBuf {
			c1, c2 = e.Theme.T3, e.Theme.Bg
			if MaxPercent > (t2 + halfBuf) {
				p = (pct - (t2 + halfBuf)) / (MaxPercent - (t2 + halfBuf))
			}
		}
	} else if e.GradientMode == ModeVerySmooth {
		if pct > t2 {
			c1, c2 = e.Theme.T2, e.Theme.Bg
			if MaxPercent > t2 {
				p = (pct - t2) / (MaxPercent - t2)
			}
		}
	}

	p = clampFloat(p, MinRatio, MaxRatio)

	// Apply TrueColor or Dithering
	if e.UseTrueColor {
		return st.Foreground(lerpColor(c1, c2, p))
	}

	if rand.Float64() > p {
		return st.Foreground(c1)
	}
	return st.Foreground(c2)
}

func drawHelp(screen tcell.Screen, e *Engine) {
	menuH := len(helpText) + 2

	if e.Width < HelpMenuWidth || e.Height < menuH {
		st := tcell.StyleDefault.Background(e.Theme.Head).Foreground(e.Theme.Bg)
		for i, c := range msgHelpWarning {
			if i < e.Width {
				screen.SetContent(i, 0, c, nil, st)
			}
		}
		return
	}

	startX := (e.Width - HelpMenuWidth) / 2
	startY := (e.Height - menuH) / 2

	borderSt := tcell.StyleDefault.Foreground(e.Theme.T1).Background(e.Theme.Bg)
	textSt := tcell.StyleDefault.Foreground(e.Theme.Head).Background(e.Theme.Bg)

	for y := 0; y < menuH; y++ {
		for x := 0; x < HelpMenuWidth; x++ {
			char := ' '
			st := borderSt
			if y > 0 && y <= len(helpText) && x > 0 && x <= len(helpText[y-1]) {
				char = rune(helpText[y-1][x-1])
				st = textSt
			}
			screen.SetContent(startX+x, startY+y, char, nil, st)
		}
	}
}

// renderScreen draws changed cells to the screen.
func renderScreen(
	screen tcell.Screen,
	engine *Engine,
	oldGrid []Cell,
	forceRedraw bool,
	overlayMsg string,
) bool {
	rendered := false

	if engine.ShowHelp {
		drawHelp(screen, engine)
		return true
	}

	for x := 0; x < engine.Width; x++ {
		colTotalTail := engine.Columns[x].Tail
		for y := 0; y < engine.Height; y++ {
			idx := x*engine.Height + y
			if idx >= len(engine.Grid) || idx >= len(oldGrid) {
				continue
			}
			c := engine.Grid[idx]
			if !forceRedraw && oldGrid[idx] == c {
				continue
			}
			oldGrid[idx] = c
			rendered = true

			if overlayMsg != "" && y == 0 && x < len(overlayMsg) {
				st := tcell.StyleDefault.Background(engine.Theme.Head).Foreground(engine.Theme.Bg)
				screen.SetContent(x, y, rune(overlayMsg[x]), nil, st)
				continue
			}

			st := tcell.StyleDefault.Background(engine.Theme.Bg)
			if c.Char == ' ' && !c.CRTBleed {
				screen.SetContent(x, y, ' ', nil, st)
				continue
			}

			if c.Flash {
				st = st.Background(engine.Theme.Head).Foreground(engine.Theme.Bg)
				screen.SetContent(x, y, c.Char, nil, st)
			} else if c.CRTBleed {
				st = st.Foreground(tcell.ColorRed)
				charToDraw := c.Char
				if charToDraw == ' ' {
					charToDraw = c.CRTChar
				}
				screen.SetContent(x, y, charToDraw, nil, st)
			} else {
				if c.Head {
					st = st.Foreground(engine.Theme.Head)
				} else {
					st = getStyle(engine, c.Tail, colTotalTail)
				}
				screen.SetContent(x, y, c.Char, nil, st)
			}
		}
	}
	return rendered
}
