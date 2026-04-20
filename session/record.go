package session

import (
	"GroORM/clause"
	"reflect"
)

// Insert inserts one or more records into the database
func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0, len(values))
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	// use vars... to expand the slice
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Find queries records from the database
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var fields []interface{}
		for _, name := range table.FieldNames {
			fields = append(fields, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(fields...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return nil
}
