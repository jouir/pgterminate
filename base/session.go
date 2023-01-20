package base

import (
	"fmt"
	"strings"
)

// Session represents a PostgreSQL backend
type Session struct {
	Pid             int64
	User            string
	Db              string
	Client          string
	State           string
	Query           string
	StateDuration   float64
	ApplicationName string
}

// NewSession instanciates a Session
func NewSession(pid int64, user string, db string, client string, state string, query string, stateDuration float64, applicationName string) *Session {
	return &Session{
		Pid:             pid,
		User:            user,
		Db:              db,
		Client:          client,
		State:           state,
		Query:           query,
		StateDuration:   stateDuration,
		ApplicationName: applicationName,
	}
}

// Format returns a Session as a string by replacing placeholders with their respective value
func (s *Session) Format(format string) string {
	definitions := map[string]string{
		"%p": fmt.Sprintf("%d", s.Pid),
		"%u": s.User,
		"%d": s.Db,
		"%r": s.Client,
		"%s": s.State,
		"%m": fmt.Sprintf("%f", s.StateDuration),
		"%q": s.Query,
		"%a": s.ApplicationName,
	}

	output := format

	for placeholder, value := range definitions {
		output = strings.Replace(output, placeholder, value, -1)
	}

	return output
}

// IsIdle returns true when a session is doing nothing
func (s *Session) IsIdle() bool {
	if s.State == "idle" || s.State == "idle in transaction" || s.State == "idle in transaction (aborted)" {
		return true
	}
	return false
}

// Equal returns true when two sessions share the same process id
func (s *Session) Equal(session *Session) bool {
	if s.Pid == 0 {
		return s.User == session.User && s.Db == session.Db && s.Client == session.Client
	}
	return s.Pid == session.Pid
}

// InSlice returns true when this sessions in present in the slice
func (s *Session) InSlice(sessions []*Session) bool {
	for _, session := range sessions {
		if s.Equal(session) {
			return true
		}
	}
	return false
}
