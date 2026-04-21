package session

import (
	"GroORM/clause"
	"errors"
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
	table := s.Model(reflect.New(destType).Interface()).RefTable()
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

// First queries the first matched record
func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("not found")
	}
	dest.Set(destSlice.Index(0))
	return nil
}

// Limit sets the LIMIT clause
func (s *Session) Limit(limit int) *Session {
	s.clause.Set(clause.LIMIT, limit)
	return s
}

// Where sets the WHERE clause
func (s *Session) Where(desc string, values ...interface{}) *Session {
	vars := make([]interface{}, 0, 1+len(values))
	vars = append(vars, desc)
	vars = append(vars, values...)
	s.clause.Set(clause.WHERE, vars...)
	return s
}

// OrderBy sets the ORDER BY clause
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

// Update updates records
func (s *Session) Update(values ...interface{}) (int64, error) {
	m, ok := values[0].(map[string]interface{})
	if !ok {
		valueLen := len(values)
		m = make(map[string]interface{}, valueLen/2)
		for i := 0; i < valueLen; i += 2 {
			m[values[i].(string)] = values[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete deletes records
func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Count returns the number of records
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
