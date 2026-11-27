package ui

import (
	"fmt"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"chat-client/internal/chat"
)

type entry struct {
	kind kind
	text string
	nick string // set for kindChat; text still holds the full "[nick] ..." line
}

type model struct {
	conn     *chat.Conn
	input    textinput.Model
	viewport viewport.Model
	messages []entry
	width    int
	height   int
}

func New(c *chat.Conn) model {
	ti := textinput.New()
	ti.Placeholder = "/N <nickname>"
	ti.SetVirtualCursor(true)
	ti.Focus()

	vp := viewport.New()
	vp.SoftWrap = true // long lines wrap instead of being cut off

	m := model{input: ti, conn: c, viewport: vp}
	m.addMessage(kindSystem, fmt.Sprintf("connected to %s (Esc or Ctrl-C to quit)", c.RemoteAddr()))
	return m
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, readLine(m.conn))
}

func (m *model) push(e entry) {
	atBottom := m.viewport.AtBottom()
	m.messages = append(m.messages, e)
	m.viewport.SetContent(m.renderMessages())
	if atBottom {
		m.viewport.GotoBottom()
	}
}

func (m *model) addMessage(k kind, text string) {
	m.push(entry{kind: k, text: text})
}

// receive handles a line from the server: "[nick] text" is a chat
// message, anything else is a plain server line.
func (m *model) receive(line string) {
	if nick, ok := chatNick(line); ok {
		m.push(entry{kind: kindChat, nick: nick, text: line})
		return
	}
	m.addMessage(kindServer, line)
}

// chatNick extracts the nick from a "[nick] text" line.
func chatNick(line string) (string, bool) {
	rest, ok := strings.CutPrefix(line, "[")
	if !ok {
		return "", false
	}
	nick, _, ok := strings.Cut(rest, "] ")
	if !ok || nick == "" || strings.ContainsAny(nick, " \t\r\n]") {
		return "", false
	}
	return nick, true
}

func (m model) renderMessages() string {
	lines := make([]string, len(m.messages))
	for i, e := range m.messages {
		lines[i] = renderEntry(e)
	}
	return strings.Join(lines, "\n")
}

func renderEntry(e entry) string {
	base := styleFor(e.kind)
	if e.nick == "" {
		return base.Render(e.text)
	}
	labelLen := len(e.nick) + 2 // "[" + nick + "]"
	return nickStyle(e.nick).Render(e.text[:labelLen]) + base.Render(e.text[labelLen:])
}
