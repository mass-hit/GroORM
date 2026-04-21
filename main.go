package main

import (
	"GroORM/log"

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
	var tests []Test
	if err := s.Find(&tests); err != nil || len(tests) != 2 {
		log.Error("fail to find data")
	}
	num, err = s.Where("Name = ?", "zero").Update("Num", 2)
	if err != nil || num != 1 {
		log.Error("fail to update data")
	}
	var test Test
	if err := s.OrderBy("Num DESC").First(&test); err != nil || test.Num != 2 {
		log.Error("fail to find first")
	}
	num, err = s.Where("Num = ?", 1).Delete()
	if err != nil || num != 1 {
		log.Error("fail to delete data")
	}
	count, err := s.Count()
	if err != nil || count != 1 {
		log.Error("fail to count data")
	}
}
