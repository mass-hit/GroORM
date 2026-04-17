package session

import (
	"GroORM/log"
	"GroORM/schema"
	"fmt"
	"reflect"
	"strings"
)

func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil || reflect.TypeOf(value) != reflect.TypeOf(s.refTable) {
		s.refTable = schema.Parse(value)
	}
	return s
}

func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("refTable is nil")
	}
	return s.refTable
}

// CreateTable executes CREATE TABLE
func (s *Session) CreateTable() error {
	var columns []string
	for _, field := range s.refTable.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw(fmt.Sprintf("CREATE TABLE %s (%s)", s.RefTable().Name, desc)).Exec()
	return err
}

// DropTable executes DROP TABLE
func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE %s", s.RefTable().Name)).Exec()
	return err
}

// TableIsExist checks whether the table exists
func (s *Session) TableIsExist() bool {
	tableName := s.RefTable().Name
	sql := "SELECT table_name FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	row := s.Raw(sql, tableName).QueryRow()
	var tmp string
	_ = row.Scan(&tmp)
	return tmp == tableName
}
