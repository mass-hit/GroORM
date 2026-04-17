package main

import (
	"GroORM/log"

	_ "github.com/go-sql-driver/mysql"
)

type Test struct {
	Name string `gro:"PRIMARY KEY"`
	Num  int
}

func main() {
	engine, _ := NewEngine("mysql", "root:@tcp(127.0.0.1:3306)/gro?charset=utf8")
	defer engine.Close()
	s := engine.NewSession().Model(&Test{})
	_ = s.CreateTable()
	_ = s.DropTable()
	if s.TableIsExist() {
		log.Error("failed")
	}
}
