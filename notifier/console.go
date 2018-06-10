package notifier

import (
	"github.com/jouir/pgterminate/base"
	"log"
)

// Console notifier structure
type Console struct {
	sessions chan base.Session
}

// NewConsole creates a console notifier
func NewConsole(sessions chan base.Session) Notifier {
	return &Console{
		sessions: sessions,
	}
}

// Run starts console notifier
func (c *Console) Run() {
	for session := range c.sessions {
		log.Printf("%s", session)
	}
}

// Reload for handling SIGHUP signals
func (c *Console) Reload() {
}
