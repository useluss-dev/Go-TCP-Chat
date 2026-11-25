package main

import (
	"fmt"
	"log"
	"net"
	"os"

	tea "charm.land/bubbletea/v2"
)

var DEFAULT_PORT = "127.0.0.1:6666"

func main() {
	port := DEFAULT_PORT

	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	addr := port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Fprintf(conn, "/N alice\n")

	p := tea.NewProgram(initialModel(conn))
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
