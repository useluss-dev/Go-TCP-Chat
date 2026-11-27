package ui

import (
	"errors"
	"io"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.layout()

	case tea.KeyPressMsg:
		return m.handleKey(msg)
	case incomingMsg:
		m.receive(string(msg))
		return m, readLine(m.conn) // wait for the next line

	case errMsg:
		if errors.Is(msg, io.EOF) {
			m.addMessage(kindSystem, "server closed the connection")
		} else {
			m.addMessage(kindSystem, "* "+msg.Error())
		}
		return m, nil
	}
	return m, nil
}

func (m model) handleKey(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit
	case "enter":
		line := strings.TrimSpace(m.input.Value())
		if line == "" {
			return m, nil
		}
		m.input.Reset()
		m.addMessage(kindSent, "> "+line)
		return m, send(m.conn, line)
	case "up", "down", "pgup", "pgdown":
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}

	m.input, cmd = m.input.Update(msg)

	return m, cmd
}

func (m *model) layout() {
	wasAtBottom := m.viewport.AtBottom()

	// input and message box heights
	ih := inputBox.GetVerticalFrameSize() + 1
	mh := max(0, m.height-ih)

	// change the size of the input and message boxes to the current width and heights
	inputBox = inputBox.Width(m.width)
	messageBox = messageBox.
		Width(m.width).
		Height(max(0, mh))

	input_inner := m.width - inputBox.GetHorizontalFrameSize() - lipgloss.Width(m.input.Prompt) - 1
	m.input.SetWidth(input_inner)
	m.input.CharLimit = input_inner

	m.viewport.SetWidth(m.width - messageBox.GetHorizontalFrameSize())
	m.viewport.SetHeight(max(0, mh-messageBox.GetVerticalFrameSize()))

	if wasAtBottom {
		m.viewport.GotoBottom()
	}
}
