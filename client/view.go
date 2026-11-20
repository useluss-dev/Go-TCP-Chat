package main

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func (m model) View() tea.View {
	input := inputBox.Render(m.input.View())

	messages := messageBox.Render(m.viewport.View())

	content := lipgloss.JoinVertical(lipgloss.Left, messages, input)

	v := tea.NewView(content)
	v.AltScreen = true

	return v
}
