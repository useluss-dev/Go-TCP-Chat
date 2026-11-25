package main

import (
	"net"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	input    textinput.Model
	viewport viewport.Model
	messages []string
	width    int
	height   int
	conn     net.Conn
}

func initialModel(conn net.Conn) model {
	ti := textinput.New()
	ti.Placeholder = "/N <nickname>"
	ti.SetVirtualCursor(false)
	ti.Focus()

	return model{input: ti, conn: conn}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
