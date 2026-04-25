package main

import (
	"GroORM/log"
	"reflect"

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
	_, _ = s.Raw("DROP TABLE IF EXISTS `test`").Exec()
	_, _ = s.Raw("CREATE TABLE `test` (name VARCHAR(255) PRIMARY KEY,tmp INT)").Exec()
	err := engine.Migrate(&Test{})
	if err != nil {
		log.Error(err)
	}
	rows, err := s.Raw("SELECT * FROM `test`").QueryRows()
	if err != nil {
		log.Error(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil || !reflect.DeepEqual(columns, []string{"name", "num"}) {
		log.Error("migrate err")
	}
}
