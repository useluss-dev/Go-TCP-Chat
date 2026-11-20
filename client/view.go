package main

import tea "charm.land/bubbletea/v2"

func (m model) View() tea.View {
	v := tea.NewView("")
	v.AltScreen = true

	return v
}
