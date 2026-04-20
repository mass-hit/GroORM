package schema

import (
	"GroORM/mysql"
	"reflect"
	"strings"
	"unicode"
)

// Field represents column
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema represents table
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
}

// RecordValues extracts the values of fields
func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}

// Parse constructs a Schema.
func Parse(dest interface{}) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model: dest,
		Name:  toSnakeCase(modelType.Name()),
	}
	for i := 0; i < modelType.NumField(); i++ {
		structField := modelType.Field(i)
		if !structField.Anonymous && structField.IsExported() {
			field := &Field{
				Name: structField.Name,
				Type: mysql.DataTypeOf(structField.Type),
			}
			if v, ok := structField.Tag.Lookup("gro"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, field.Name)
		}
	}
	return schema
}

// toSnakeCase converts a CamelCase string to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
