package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	engine, _ := NewEngine("mysql", "root:@tcp(127.0.0.1:3306)/gro?charset=utf8")
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name VARCHAR(255));").Exec()
	result, _ := s.Raw("INSERT INTO User(`Name`) values (?),(?)", "zero", "hhh").Exec()
	count, _ := result.RowsAffected()
	fmt.Println(count)
}
