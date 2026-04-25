package main

import (
	"GroORM/log"
	"GroORM/session"
	"database/sql"
	"errors"
	"fmt"
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

// TxFunc defines the function for transactional
type TxFunc func(*session.Session) (interface{}, error)

// Transaction executes the function within a transaction.
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

// difference returns elements present in `a` but not in `b`
func difference(a, b []string) (diff []string) {
	bMap := make(map[string]bool)
	for _, v := range b {
		bMap[v] = true
	}
	for _, v := range a {
		if _, ok := bMap[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate performs a simple schema migration
func (e *Engine) Migrate(value interface{}) error {
	s := e.NewSession()
	if !s.Model(value).TableIsExist() {
		return errors.New("table is not exist")
	}
	table := s.RefTable()
	rows, err := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
	if err != nil {
		return err
	}
	columns, err := rows.Columns()
	fieldNames := table.FieldNames
	addCols := difference(fieldNames, columns)
	delCols := difference(columns, fieldNames)
	for _, col := range addCols {
		field := table.Field(col)
		if _, err = s.Raw(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s %s", table.Name, field.Name, field.Type, field.Tag)).Exec(); err != nil {
			return err
		}
	}
	for _, col := range delCols {
		if _, err = s.Raw(fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s", table.Name, col)).Exec(); err != nil {
			return err
		}
	}
	return nil
}
