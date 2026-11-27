package main

import (
	"bufio"
	"log"
	"net"
	"regexp"
	"strings"
)

var (
	nickRegex  = regexp.MustCompile(`^[a-z][a-zA-Z0-9_]{1,19}$`) // 2-20 chars, starts lowercase
	groupRegex = regexp.MustCompile(`^[A-Z][a-zA-Z0-9_]{1,19}$`) // 2-20 chars, starts uppercase
)

const defaultGroup = "default"

type Worker struct {
	conn   net.Conn
	man    *Manager
	out    chan string
	nick   string
	groups map[string][]string
}

func NewWorker(conn net.Conn, man *Manager) *Worker {
	return &Worker{
		conn:   conn,
		man:    man,
		out:    make(chan string, 64),
		groups: map[string][]string{defaultGroup: {}},
	}
}

func (w *Worker) Run() {
	writerDone := make(chan struct{})
	go w.writeLoop(writerDone)

	defer func() {
		w.man.Remove(w)
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
			w.handleLine(line)
			log.Printf("%s: %s", w.conn.RemoteAddr(), line)
		}
	}
}

func (w *Worker) Deliver(text string) {
	select {
	case w.out <- text:
	default:
		log.Printf("dropping message for slow client %s", w.conn.RemoteAddr())
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

func (w *Worker) handleLine(line string) {
	log.Printf("%s", strings.TrimRight(line, "\r\n"))

	fields := strings.Fields(line)
	if len(fields) == 0 {
		return
	}
	cmd, args := normalize(fields[0]), fields[1:]

	switch cmd {
	case "/L": // allowed without nickname
		w.cmdList()
	case "/N":
		w.cmdNick(args)
	case "/M", "/B", "/D", "/G":
		if w.nick == "" {
			w.tell("set a nickname first with /N <nickname>")
			return
		}
		switch cmd {
		case "/M":
			w.cmdMessage(args)
		case "/B":
			w.cmdBroadcast(args)
		case "/D":
			w.cmdDefault(args)
		case "/G":
			w.cmdGroup(args)
		}
	default:
		w.tell(cmd + " is an invalid command")
	}
}

func (w *Worker) cmdList() {
	nicks := w.man.ListNicks()
	if len(nicks) == 0 {
		w.tell("no users online")
		return
	}
	w.tell("online:\n" + strings.Join(nicks, "\n"))
}

func (w *Worker) cmdNick(args []string) {
	if len(args) == 0 {
		w.tell("usage: /N <nickname>")
		return
	}

	nick := args[0]
	if !nickRegex.MatchString(nick) {
		w.tell("valid nicknames are 2-20 characters, start with a lowercase letter and only contain alphanumeric and underscore characters")
		return
	}

	switch w.man.SetNick(w, nick) {
	case NickNew:
		w.nick = nick
		w.tell("nickname successfully set to " + nick)
	case NickUpdated:
		w.nick = nick
		w.tell("nickname successfully updated to " + nick)
	case NickUnchanged:
		w.tell("you already have the nickname " + nick)
	case NickTaken:
		w.tell(nick + " is already taken please try another name")
	}
}

func (w *Worker) cmdMessage(args []string) {
	if len(args) < 2 {
		w.tell("usage: /M <name[,name...]> <message>")
		return
	}
	w.sendTo(w.expand(splitNames(args[0])), strings.Join(args[1:], " "))
}

func (w *Worker) cmdBroadcast(args []string) {
	if len(args) == 0 {
		w.tell("usage: /B <message>")
		return
	}
	w.sendTo(w.groups[defaultGroup], strings.Join(args, " "))
}

func (w *Worker) cmdDefault(args []string) {
	if len(args) == 0 {
		w.tell("usage: /D <name[,name...]>")
		return
	}
	nicks, bad, ok := w.resolve(splitNames(args[0]))
	if !ok {
		w.tell("invalid name: " + bad)
		return
	}
	w.groups[defaultGroup] = nicks
	w.tell("default group: " + strings.Join(nicks, ","))
}

func (w *Worker) cmdGroup(args []string) {
	if len(args) < 2 {
		w.tell("usage: /G <Group> <name[,name...]>")
		return
	}
	group := args[0]
	if !groupRegex.MatchString(group) {
		w.tell("group names must start with an uppercase letter and only contain alphanumeric and underscore characters")
		return
	}
	nicks, bad, ok := w.resolve(splitNames(args[1]))
	if !ok {
		w.tell("invalid name: " + bad)
		return
	}
	w.groups[group] = nicks
	w.tell("group " + group + ": " + strings.Join(nicks, ","))
}

func (w *Worker) sendTo(nicks []string, text string) {
	w.man.Send(w, nicks, text)
	w.tell("message sent")
}

func (w *Worker) expand(names []string) []string {
	var result []string
	seen := make(map[string]bool)
	add := func(n string) {
		if !seen[n] {
			seen[n] = true
			result = append(result, n)
		}
	}
	for _, name := range names {
		if groupRegex.MatchString(name) {
			for _, member := range w.groups[name] {
				add(member)
			}
		} else {
			add(name)
		}
	}
	return result
}

func (w *Worker) resolve(names []string) (nicks []string, invalid string, ok bool) {
	for _, n := range names {
		if !nickRegex.MatchString(n) && !groupRegex.MatchString(n) {
			return nil, n, false
		}
	}
	return w.expand(names), "", true
}

func normalize(command string) string {
	r := []rune(command)
	if len(r) > 2 {
		r = r[:2]
	}
	return strings.ToUpper(string(r))
}

func splitNames(names string) []string {
	var out []string
	for s := range strings.SplitSeq(names, ",") {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
