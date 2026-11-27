package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:6666")
	if err != nil {
		log.Fatalln("Error listening:", err)
	}
	defer listener.Close()
	log.Println("Server running on :6666")

	man := NewManager()
	serve(listener, man)
}

func serve(ln net.Listener, man *Manager) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting:", err)
			continue
		}
		log.Println("Accepted connection from:", conn.RemoteAddr())

		w := NewWorker(conn, man)
		go w.Run()
	}
}
