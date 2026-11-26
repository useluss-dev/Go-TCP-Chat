package ui

import (
	tea "charm.land/bubbletea/v2"
	"chat-client/internal/chat"
)

type (
	incomingMsg string
	errMsg      error
)

func readLine(c *chat.Conn) tea.Cmd {
	return func() tea.Msg {
		line, err := c.ReadLine()
		if err != nil {
			return errMsg(err)
		}
		return incomingMsg(line)
	}
}

// message.go
func send(c *chat.Conn, line string) tea.Cmd {
	return func() tea.Msg {
		if err := c.Send(line); err != nil {
			return errMsg(err)
		}
		return nil
	}
}
