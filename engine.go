package main

import (
	"GroORM/log"
	"GroORM/session"
	"database/sql"
)

// Engine wraps a database connection
type Engine struct {
	db *sql.DB
}

// NewEngine initializes a new Engine instance
func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	// Ensure the DB is closed
	defer func() {
		if err != nil {
			db.Close()
		}
	}()
	if err != nil {
		log.Error(err)
		return
	}
	// Verify the database is reachable
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	return &Engine{db: db}, nil
}

// Close releases the database connection
func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error(err)
	}
}

// NewSession creates a new Session
func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}

type TxFunc func(*session.Session) (interface{}, error)

func (e *Engine) Transaction(fn TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()
	return fn(s)
}
