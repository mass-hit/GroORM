package main

import (
	"GroORM/log"
	"GroORM/session"

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
	num, err := s.Insert(zero, one)
	if err != nil || num != 2 {
		log.Error("fail to insert data")
	}
	var test Test
	if err := s.OrderBy("Num").First(&test); err != nil || test.Num != 1 || test.Name != "***" {
		log.Error("fail to call hook method")
	}
}
