package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
)

// SetResult describes the outcome of a nickname request.
type SetResult int

const (
	NickNew       SetResult = iota // client had no nickname, now has one
	NickUpdated                    // client changed an existing nickname
	NickUnchanged                  // client asked for the nickname it already has
	NickTaken                      // someone else owns the nickname
)

// SendResult describes the outcome of a send
type SendResult int

const (
	SuccessfulSend SendResult = iota // Message successfully sent to all named clients
	PartialSend                      // Message was only able to be sent to some clients
	FailedSend                       // Message failed to send to any client
)

// Manager keeps track of clients and manages their shared data (nicknames)
type Manager struct {
	mu       sync.Mutex
	byNick   map[string]*Worker
	byWorker map[*Worker]string
}

func NewManager() *Manager {
	return &Manager{
		byNick:   make(map[string]*Worker),
		byWorker: make(map[*Worker]string),
	}
}

// Tries to set the clients nickname returns SetResult based on what happened
func (m *Manager) SetNick(w *Worker, nick string) SetResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	if owner, ok := m.byNick[nick]; ok {
		if owner == w {
			return NickUnchanged
		}
		return NickTaken
	}

	old, had := m.byWorker[w]
	if had {
		delete(m.byNick, old) // Release previous nickname
	}
	m.byNick[nick] = w
	m.byWorker[w] = nick

	if had {
		log.Printf("nickname change: %s -> %s", old, nick)
		return NickUpdated
	}
	log.Println("New nickname registered:", nick)
	return NickNew
}

// Sends the workers message to the clients named if they exist. Returns SendResult
// based on how many messages were successfully sent.
func (m *Manager) Send(from *Worker, nicks []string, text string) SendResult {
	m.mu.Lock()
	defer m.mu.Unlock()

	fromNick, ok := m.byWorker[from]
	if !ok {
		return FailedSend
	}
	msg := fmt.Sprintf("[%s] %s", fromNick, text)

	// The number of messages successfully sent
	num_sent := 0

	// Only send message to each client only once, even if they were named more
	seen := make(map[string]bool, len(nicks))
	for _, nick := range nicks {
		if seen[nick] {
			continue
		}
		seen[nick] = true
		if target, ok := m.byNick[nick]; ok {
			num_sent++
			target.Deliver(msg)
		}
	}

	// Message wasn't sent to any clients
	if num_sent == 0 {
		return FailedSend
	}

	if num_sent < len(nicks) {
		return PartialSend
	}

	return SuccessfulSend
}

// Returns a sorted list of all the nicknames
func (m *Manager) ListNicks() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	nicks := make([]string, 0, len(m.byNick))
	for nick := range m.byNick {
		nicks = append(nicks, nick)
	}
	sort.Strings(nicks)
	return nicks
}

// Removes a worker and releases their nickname
func (m *Manager) Remove(w *Worker) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if nick, ok := m.byWorker[w]; ok {
		delete(m.byNick, nick)
		delete(m.byWorker, w)
		log.Println("nickname released:", nick)
	}
}
