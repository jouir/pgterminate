package notifier

import (
	"github.com/jouir/pgterminate/base"
	"github.com/jouir/pgterminate/log"
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
	log.Info("Starting console notifier")
	for session := range c.sessions {
		log.Infof("%s\n", session)
	}
}

// Reload for handling SIGHUP signals
func (c *Console) Reload() {
	log.Info("Reloading console notifier")
}
