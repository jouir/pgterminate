package base

import (
	"database/sql"
	"github.com/lib/pq"
	"strconv"
)

const (
	maxQueryLength = 1000
)

// Db centralizes connection to the database
type Db struct {
	dsn  string
	conn *sql.DB
}

// NewDb creates a Db object
func NewDb(dsn string) *Db {
	return &Db{
		dsn: dsn,
	}
}

// Connect connects to the instance and ping it to ensure connection is working
func (db *Db) Connect() {
	conn, err := sql.Open("postgres", db.dsn)
	Panic(err)

	err = conn.Ping()
	Panic(err)

	db.conn = conn
}

// Disconnect ends connection cleanly
func (db *Db) Disconnect() {
	err := db.conn.Close()
	Panic(err)
}

// Sessions connects to the database and returns current sessions
func (db *Db) Sessions() (sessions []Session) {
	query := `select pid as pid, usename as user, datname as db, host(client_addr)::text || ':' || client_port::text as client, state as state, substring(query from 1 for ` + strconv.Itoa(maxQueryLength) + `) as query, coalesce(extract(epoch from now() - state_change), 0) as "stateDuration" from pg_catalog.pg_stat_activity where pid <> pg_backend_pid();`
	rows, err := db.conn.Query(query)
	Panic(err)
	defer rows.Close()

	for rows.Next() {
		var pid sql.NullInt64
		var user, db, client, state, query sql.NullString
		var stateDuration float64
		err := rows.Scan(&pid, &user, &db, &client, &state, &query, &stateDuration)
		Panic(err)

		if pid.Valid && user.Valid && db.Valid && client.Valid && state.Valid && query.Valid {
			sessions = append(sessions, NewSession(pid.Int64, user.String, db.String, client.String, state.String, query.String, stateDuration))
		}
	}

	return sessions
}

// TerminateSessions terminates a list of sessions
func (db *Db) TerminateSessions(sessions []Session) {
	var pids []int64
	for _, session := range sessions {
		pids = append(pids, session.Pid)
	}
	if len(pids) > 0 {
		query := `select pg_terminate_backend(pid) from pg_stat_activity where pid = any($1);`
		_, err := db.conn.Exec(query, pq.Array(pids))
		Panic(err)
	}
}
