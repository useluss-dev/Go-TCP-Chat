package ui

import (
	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"chat-client/internal/chat"
)

type model struct {
	conn     *chat.Conn
	input    textinput.Model
	viewport viewport.Model
	messages []string
	width    int
	height   int
}

func New(c *chat.Conn) model {
	ti := textinput.New()
	ti.Placeholder = "/N <nickname>"
	ti.SetVirtualCursor(false)
	ti.Focus()

	return model{input: ti, conn: c}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}
