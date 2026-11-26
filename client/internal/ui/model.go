package ui

import (
	"strings"

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
	ti.SetVirtualCursor(true)
	ti.Focus()

	return model{input: ti, conn: c, viewport: viewport.New()}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, readLine(m.conn))
}

func (m *model) addMessage(text string) {
	atBottom := m.viewport.AtBottom()
	m.messages = append(m.messages, text)
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	if atBottom {
		m.viewport.GotoBottom()
	}
}
