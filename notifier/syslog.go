package notifier

import (
	"fmt"
	"github.com/jouir/pgterminate/base"
	"github.com/jouir/pgterminate/log"
	"log/syslog"
)

// Syslog notifier
type Syslog struct {
	sessions chan base.Session
	ident    string
	priority syslog.Priority
	writer   *syslog.Writer
}

// NewSyslog creates a syslog notifier
func NewSyslog(facility string, ident string, sessions chan base.Session) Notifier {
	var priority syslog.Priority
	switch facility {
	case "LOCAL0":
		priority = syslog.LOG_INFO | syslog.LOG_LOCAL0
	case "LOCAL1":
		priority = syslog.LOG_INFO | syslog.LOG_LOCAL1
	case "LOCAL2":
		priority = syslog.LOG_INFO | syslog.LOG_LOCAL2
	case "LOCAL3":
		priority = syslog.LOG_INFO | syslog.LOG_LOCAL3
	case "LOCAL4":
		priority = syslog.LOG_INFO | syslog.LOG_LOCAL4
	case "LOCAL5":
		priority = syslog.LOG_INFO | syslog.LOG_LOCAL5
	case "LOCAL6":
		priority = syslog.LOG_INFO | syslog.LOG_LOCAL6
	default: // LOCAL7
		priority = syslog.LOG_INFO | syslog.LOG_LOCAL7
	}
	return &Syslog{
		sessions: sessions,
		ident:    ident,
		priority: priority,
	}
}

// Run starts syslog notifier
func (s *Syslog) Run() {
	log.Info("Starting syslog notifier")
	var err error
	if s.writer, err = syslog.New(s.priority, s.ident); err != nil {
		base.Panic(err)
	}
	for session := range s.sessions {
		s.writer.Info(fmt.Sprintf("%s", session))
	}
}

// Reload disconnect from syslog daemon and re-connect
// Executed when receiving SIGHUP signal
func (s *Syslog) Reload() {
	log.Info("Reloading syslog notifier")
	if s.writer != nil {
		log.Debug("Re-connecting to syslog daemon")
		s.disconnect()
		s.connect()
	}
}

// connect to syslog daemon
func (s *Syslog) connect() {
	var err error
	if s.writer, err = syslog.New(s.priority, s.ident); err != nil {
		base.Panic(err)
	}
}

// disconnect from syslog daemon
func (s *Syslog) disconnect() {
	var err error
	if err = s.writer.Close(); err != nil {
		base.Panic(err)
	}
}
