package session

import (
	"GroORM/log"
	"errors"
)

// Begin starts a new transaction
func (s *Session) Begin() (err error) {
	if s.tx != nil {
		return errors.New("transaction already started")
	}
	log.Info("begin transaction")
	if s.tx, err = s.db.Begin(); err != nil {
		log.Error(err)
	}
	return
}

// Commit commits the transaction
func (s *Session) Commit() (err error) {
	if s.tx == nil {
		return errors.New("no transaction")
	}
	log.Info("commit transaction")
	defer func() {
		s.tx = nil
	}()
	if err = s.tx.Commit(); err != nil {
		log.Error(err)
	}
	// reset transaction
	s.tx = nil
	return
}

// Rollback aborts the transaction
func (s *Session) Rollback() (err error) {
	if s.tx == nil {
		return errors.New("no transaction")
	}
	log.Info("rollback transaction")
	defer func() {
		s.tx = nil
	}()
	if err = s.tx.Rollback(); err != nil {
		log.Error(err)
	}
	// reset transaction
	s.tx = nil
	return
}
