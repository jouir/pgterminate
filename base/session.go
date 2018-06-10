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
	BackendDuration float64
	XactDuration    float64
	QueryDuration   float64
}

// NewSession instanciates a Session
func NewSession(pid int64, user string, db string, client string, state string, query string, backendDuration float64, xactDuration float64, queryDuration float64) Session {
	return Session{
		Pid:             pid,
		User:            user,
		Db:              db,
		Client:          client,
		State:           state,
		Query:           query,
		BackendDuration: backendDuration,
		XactDuration:    xactDuration,
		QueryDuration:   queryDuration,
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
	if s.BackendDuration != 0 {
		output = append(output, fmt.Sprintf("backend_duration=%f", s.BackendDuration))
	}
	if s.XactDuration != 0 {
		output = append(output, fmt.Sprintf("xact_duration=%f", s.XactDuration))
	}
	if s.QueryDuration != 0 {
		output = append(output, fmt.Sprintf("query_duration=%f", s.QueryDuration))
	}
	if s.Query != "" {
		output = append(output, fmt.Sprintf("query=%s", s.Query))
	}
	return strings.Join(output, " ")
}
