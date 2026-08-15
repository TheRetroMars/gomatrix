package main

import "github.com/gdamore/tcell/v2"

// Theme defines the colors used for rendering the matrix effect.
type Theme struct {
	Head tcell.Color
	T1   tcell.Color
	T2   tcell.Color
	T3   tcell.Color
	Bg   tcell.Color
}

// getTheme returns the Theme struct corresponding to the given name.
// If the name is not recognized, it returns the Classic theme.
func getTheme(name string) Theme {
	switch name {
	case "AntiGravity":
		return Theme{
			Head: tcell.NewHexColor(0xF0F8FF),
			T1:   tcell.NewHexColor(0x00F0FF),
			T2:   tcell.NewHexColor(0xBD00FF),
			T3:   tcell.NewHexColor(0x120033),
			Bg:   tcell.NewHexColor(0x0A0B1E),
		}
	case "QuantumGold":
		return Theme{
			Head: tcell.NewHexColor(0xFFF8DC),
			T1:   tcell.NewHexColor(0xFFB800),
			T2:   tcell.NewHexColor(0xFF5500),
			T3:   tcell.NewHexColor(0x330A00),
			Bg:   tcell.NewHexColor(0x0C0A05),
		}
	case "Cyberpunk":
		return Theme{
			Head: tcell.NewHexColor(0xFFE6F2),
			T1:   tcell.NewHexColor(0xFF007F),
			T2:   tcell.NewHexColor(0x7FFF00),
			T3:   tcell.NewHexColor(0x200020),
			Bg:   tcell.NewHexColor(0x000000),
		}
	default: // Classic
		return Theme{
			Head: tcell.NewHexColor(0xFFFFFF),
			T1:   tcell.NewHexColor(0x00FF00),
			T2:   tcell.NewHexColor(0x003300),
			T3:   tcell.NewHexColor(0x003300),
			Bg:   tcell.NewHexColor(0x000000),
		}
	}
}
