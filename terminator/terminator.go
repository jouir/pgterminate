package terminator

import (
	"github.com/jouir/pgterminate/base"
	"log"
	"time"
)

// Terminator looks for sessions, filters actives and idles, terminate them and notify sessions channel
// It ends itself gracefully when done channel is triggered
type Terminator struct {
	config   *base.Config
	db       *base.Db
	sessions chan base.Session
	done     chan bool
}

// NewTerminator instanciates a Terminator
func NewTerminator(ctx *base.Context) *Terminator {
	return &Terminator{
		config:   ctx.Config,
		sessions: ctx.Sessions,
		done:     ctx.Done,
	}
}

// Run starts the Terminator
func (t *Terminator) Run() {
	t.db = base.NewDb(t.config.Dsn())
	log.Println("Connecting to instance")
	t.db.Connect()
	defer t.terminate()

	for {
		select {
		case <-t.done:
			return
		default:
			sessions := t.db.Sessions()
			if t.config.ActiveTimeout != 0 {
				actives := activeSessions(sessions, t.config.ActiveTimeout)
				go t.terminateAndNotify(actives)
			}

			if t.config.IdleTimeout != 0 {
				idles := idleSessions(sessions, t.config.IdleTimeout)
				go t.terminateAndNotify(idles)
			}
			time.Sleep(time.Duration(t.config.Interval*1000) * time.Millisecond)
		}

	}
}

// terminateAndNotify terminates a list of sessions and notifies channel
func (t *Terminator) terminateAndNotify(sessions []base.Session) {
	t.db.TerminateSessions(sessions)
	for _, session := range sessions {
		t.sessions <- session
	}
}

// terminate terminates gracefully
func (t *Terminator) terminate() {
	log.Println("Disconnecting from instance")
	t.db.Disconnect()
}

// activeSessions returns a list of active sessions
// A session is active when state is "active" and backend has started before elapsed
// seconds
func activeSessions(sessions []base.Session, elapsed float64) (result []base.Session) {
	for _, session := range sessions {
		if session.State == "active" && session.QueryDuration > elapsed {
			result = append(result, session)
		}
	}
	return result
}

// idleSessions returns a list of idle sessions
// A sessions is idle when state is "idle" and backend has started before elapsed seconds
// and when state is "idle in transaction" or "idle in transaction (aborted)" and
// transaction has started before elapsed seconds
func idleSessions(sessions []base.Session, elapsed float64) (result []base.Session) {
	for _, session := range sessions {
		if session.State == "idle" && session.BackendDuration > elapsed {
			result = append(result, session)
		} else if (session.State == "idle in transaction" || session.State == "idle in transaction (aborted)") && session.XactDuration > elapsed {
			result = append(result, session)
		}
	}
	return result
}
