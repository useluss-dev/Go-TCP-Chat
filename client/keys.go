package main

import "charm.land/bubbles/v2/key"

type keyMap struct {
	Quit key.Binding
}

var keys = keyMap{
	Quit: key.NewBinding(key.WithKeys("q", "ctrl+c")),
}
