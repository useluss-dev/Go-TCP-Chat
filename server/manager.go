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

func (m *Manager) Send(from *Worker, nicks []string, text string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	fromNick, ok := m.byWorker[from]
	if !ok {
		return false
	}
	msg := fmt.Sprintf("[%s] %s", fromNick, text)

	seen := make(map[string]bool, len(nicks))
	for _, nick := range nicks {
		if seen[nick] {
			continue
		}
		seen[nick] = true
		if target, ok := m.byNick[nick]; ok {
			target.Deliver(msg)
		}
	}
	return true
}

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

func (m *Manager) Remove(w *Worker) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if nick, ok := m.byWorker[w]; ok {
		delete(m.byNick, nick)
		delete(m.byWorker, w)
		log.Println("nickname released:", nick)
	}
}
