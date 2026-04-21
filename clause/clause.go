package clause

import "strings"

// Clause represents SQL fragment
type Clause struct {
	sql     map[int]string
	sqlVars map[int][]interface{}
}

const (
	INSERT = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

// Set registers a SQL fragment
func (c *Clause) Set(name int, vars ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[int]string)
		c.sqlVars = make(map[int][]interface{})
	}
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

// Build generate the final SQL statement
func (c *Clause) Build(orders ...int) (string, []interface{}) {
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql, ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
