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
	serve(listener)
}

func serve(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting:", err)
			continue
		}
		log.Println("Accepted connection from:", conn.RemoteAddr())

		w := NewWorker(conn)
		go w.Run()
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error reading:", err)
			return
		}
		log.Printf("Recieved: %s", buffer[:n])
		conn.Write([]byte("Message received\n"))
	}
}
