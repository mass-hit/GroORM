package mysql

import (
	"fmt"
	"reflect"
	"time"
)

func DataTypeOf(typ reflect.Type) string {
	switch typ.Kind() {
	case reflect.Bool:
		return `BOOL`
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return `INT`
	case reflect.Int64, reflect.Uint64:
		return `BIGINT`
	case reflect.Float32, reflect.Float64:
		return `DOUBLE`
	case reflect.String:
		return `VARCHAR(255)`
	case reflect.Slice, reflect.Array:
		if typ.Elem().Kind() == reflect.Uint8 {
			return `BLOB`
		}
		return `TEXT`
	case reflect.Struct:
		if typ == reflect.TypeOf(time.Time{}) {
			return `TIMESTAMP`
		}
		return `TEXT`
	}
	panic(fmt.Sprintf("unsupported type: %s (%s)", typ.Name(), typ.Kind()))
}
