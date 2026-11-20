package main

import (
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
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "/N <nickname>"
	ti.SetVirtualCursor(false)
	ti.Focus()

	return model{input: ti}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
