package main

import (
	"GroORM/log"
	"GroORM/session"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

type Test struct {
	Name string `gro:"PRIMARY KEY"`
	Num  int
}

var (
	zero = &Test{Name: "zero", Num: 0}
	one  = &Test{Name: "one", Num: 1}
)

func (t *Test) BeforeInsert(s *session.Session) {
	log.Info("before insert")
	t.Num++
}

func (t *Test) AfterQuery(s *session.Session) {
	log.Info("after query")
	t.Name = "***"
}

func main() {
	engine, _ := NewEngine("mysql", "root:@tcp(127.0.0.1:3306)/gro?charset=utf8")
	defer engine.Close()
	s := engine.NewSession().Model(&Test{})
	_ = s.DropTable()
	_ = s.CreateTable()
	_ = s.Model(&Test{})
	var test Test
	if _, err := engine.Transaction(func(s *session.Session) (interface{}, error) {
		_, _ = s.Insert(zero, one)
		return nil, errors.New("error")
	}); err == nil || s.First(&test) == nil {
		log.Error("rollback fail")
	}
	if _, err := engine.Transaction(func(s *session.Session) (interface{}, error) {
		_, _ = s.Insert(zero, one)
		return nil, nil
	}); err != nil {
		log.Error("commit fail")
	}
	if num, err := s.Count(); err != nil || num != 2 {
		log.Error("count fail")
	}
}
