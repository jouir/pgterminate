package base

import (
	"fmt"
	"strings"
)

// Session represents a PostgreSQL backend
type Session struct {
	Pid           int64
	User          string
	Db            string
	Client        string
	State         string
	Query         string
	StateDuration float64
}

// NewSession instanciates a Session
func NewSession(pid int64, user string, db string, client string, state string, query string, stateDuration float64) Session {
	return Session{
		Pid:           pid,
		User:          user,
		Db:            db,
		Client:        client,
		State:         state,
		Query:         query,
		StateDuration: stateDuration,
	}
}

// String represents a Session as a string
func (s Session) String() string {
	var output []string
	if s.Pid != 0 {
		output = append(output, fmt.Sprintf("pid=%d", s.Pid))
	}
	if s.User != "" {
		output = append(output, fmt.Sprintf("user=%s", s.User))
	}
	if s.Db != "" {
		output = append(output, fmt.Sprintf("db=%s", s.Db))
	}
	if s.Client != "" {
		output = append(output, fmt.Sprintf("client=%s", s.Client))
	}
	if s.State != "" {
		output = append(output, fmt.Sprintf("state=%s", s.State))
	}
	if s.StateDuration != 0 {
		output = append(output, fmt.Sprintf("state_duration=%f", s.StateDuration))
	}
	if s.Query != "" && !s.IsIdle() {
		output = append(output, fmt.Sprintf("query=%s", s.Query))
	}
	return strings.Join(output, " ")
}

// IsIdle returns true when a session is doing nothing
func (s Session) IsIdle() bool {
	if s.State == "idle" || s.State == "idle in transaction" || s.State == "idle in transaction (aborted)" {
		return true
	}
	return false
}
