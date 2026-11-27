package ui

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

var (
	inputBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	messageBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())
)

type kind int

const (
	kindServer kind = iota // acks, errors, ...
	kindChat               // "[nick] text" from another user
	kindSent               // a line you sent
	kindSystem             // notice from the client itself
)

var (
	styleChat   = lipgloss.NewStyle()
	styleServer = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))              // Yellow
	styleSent   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))              // DarkGray
	styleSystem = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Italic(true) // Cyan
)

func styleFor(k kind) lipgloss.Style {
	switch k {
	case kindServer:
		return styleServer
	case kindSent:
		return styleSent
	case kindSystem:
		return styleSystem
	default:
		return styleChat
	}
}

// Yellow, cyan and gray are left out because server lines, system
// notices and sent lines already use them.
var nickPalette = []color.Color{
	lipgloss.Color("2"),  // Green
	lipgloss.Color("5"),  // Magenta
	lipgloss.Color("9"),  // LightRed
	lipgloss.Color("12"), // LightBlue
	lipgloss.Color("10"), // LightGreen
	lipgloss.Color("13"), // LightMagenta
}

// nickStyle gives each nick a stable color.
func nickStyle(nick string) lipgloss.Style {
	var h uint32
	for i := 0; i < len(nick); i++ {
		h = h*31 + uint32(nick[i])
	}
	return lipgloss.NewStyle().
		Foreground(nickPalette[int(h)%len(nickPalette)]).
		Bold(true)
}
