package session

import (
	"GroORM/clause"
	"GroORM/log"
	"GroORM/schema"
	"database/sql"
	"strings"
)

// Session keep a database connection
type Session struct {
	db       *sql.DB
	refTable *schema.Schema
	clause   clause.Clause
	sql      strings.Builder
	sqlVars  []interface{}
}

// New creates a new Session
func New(db *sql.DB) *Session {
	return &Session{db: db}
}

// Clear resets the SQL builder and parameter
func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlVars = nil
	s.clause = clause.Clause{}
}

// Exec executes a non-query SQL
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if result, err = s.db.Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// QueryRow return at most one row
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	return s.db.QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows return multiple rows
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlVars)
	if rows, err = s.db.Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// Raw appends a SQL fragment and variables
func (s *Session) Raw(sql string, sqlVars ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, sqlVars...)
	return s
}
