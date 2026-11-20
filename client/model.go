package main

import (
	tea "charm.land/bubbletea/v2"
)

type model struct {
	width  int
	height int
}

func initialModel() model {
	return model{}
}

func (m model) Init() tea.Cmd {
	return nil
}
