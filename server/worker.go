package main

import (
	"bufio"
	"log"
	"net"
)

const defaultGroup = "default"

type Worker struct {
	conn   net.Conn
	out    chan string
	nick   string
	groups map[string][]string
}

func NewWorker(conn net.Conn) *Worker {
	return &Worker{
		conn:   conn,
		out:    make(chan string, 64),
		groups: map[string][]string{defaultGroup: {}},
	}
}

func (w *Worker) Run() {
	writerDone := make(chan struct{})
	go w.writeLoop(writerDone)

	defer func() {
		close(w.out) // flush writer and exit
		<-writerDone
		w.conn.Close()
		log.Printf("client disconnected: %s", w.conn.RemoteAddr())
	}()

	reader := bufio.NewReader(w.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		if len(line) > 0 {
			// w.handleLine(line)
			w.tell("message received") // acknowledge every line
			log.Printf("%s: %s", w.conn.RemoteAddr(), line)
		}
	}
}

func (w *Worker) writeLoop(done chan<- struct{}) {
	defer close(done)
	failed := false
	for line := range w.out {
		if failed {
			continue // keep draining so senders never block
		}
		if _, err := w.conn.Write([]byte(line + "\n")); err != nil {
			log.Printf("send failed: %v", err)
			failed = true
			w.conn.Close() // makes the read loop return
		}
	}
}

func (w *Worker) tell(text string) {
	w.out <- text
}
