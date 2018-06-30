package terminator

import (
	"github.com/jouir/pgterminate/base"
	"github.com/jouir/pgterminate/log"
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
	log.Info("Starting terminator")
	t.db = base.NewDb(t.config.Dsn())
	log.Info("Connecting to instance")
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
				t.terminateAndNotify(t.filter(actives))
			}

			if t.config.IdleTimeout != 0 {
				idles := idleSessions(sessions, t.config.IdleTimeout)
				t.terminateAndNotify(t.filter(idles))
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

// filter removes sessions according to include and exclude users settings
// when include users slice and regex are not set, append all sessions except excluded users
// otherwise, append included users
func (t *Terminator) filter(sessions []base.Session) (filtered []base.Session) {
	includeUsers, includeRegex := t.config.IncludeUsers, t.config.IncludeUsersRegexCompiled
	excludeUsers, excludeRegex := t.config.ExcludeUsers, t.config.ExcludeUsersRegexCompiled

	for _, session := range sessions {
		if t.config.IncludeUsers == nil && includeRegex == nil {
			// append all sessions except excluded users
			if !base.InSlice(session.User, excludeUsers) || (excludeRegex != nil && !excludeRegex.MatchString(session.User)) {
				filtered = append(filtered, session)
			}
		} else {
			// append included users only
			if base.InSlice(session.User, includeUsers) || (includeRegex != nil && includeRegex.MatchString(session.User)) {
				filtered = append(filtered, session)
			}
		}
	}

	return filtered
}

// terminate terminates gracefully
func (t *Terminator) terminate() {
	log.Info("Disconnecting from instance")
	t.db.Disconnect()
}

// activeSessions returns a list of active sessions
// A session is active when state is "active" and state has changed before elapsed seconds
// seconds
func activeSessions(sessions []base.Session, elapsed float64) (result []base.Session) {
	for _, session := range sessions {
		if session.State == "active" && session.StateDuration > elapsed {
			result = append(result, session)
		}
	}
	return result
}

// idleSessions returns a list of idle sessions
// A sessions is idle when state is "idle",  "idle in transaction" or "idle in transaction
// (aborted)"and state has changed before elapsed seconds
func idleSessions(sessions []base.Session, elapsed float64) (result []base.Session) {
	for _, session := range sessions {
		if session.IsIdle() && session.StateDuration > elapsed {
			result = append(result, session)
		}
	}
	return result
}
