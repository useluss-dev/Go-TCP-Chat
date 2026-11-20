package main

import "charm.land/lipgloss/v2"

var (
	inputBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(0, 1)

	messageBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder())
)
