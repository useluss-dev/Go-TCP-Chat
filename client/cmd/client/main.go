package main

import (
	"log"
	"os"

	"chat-client/internal/chat"
	"chat-client/internal/ui"

	tea "charm.land/bubbletea/v2"
)

var DEFAULT_ADDR = "127.0.0.1:6666"

func main() {
	addr := DEFAULT_ADDR
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}

	conn := chat.Dial(addr)
	defer conn.Close()

	p := tea.NewProgram(ui.New(conn))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
