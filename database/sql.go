package database

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// SelectBuilder builds select statement.
type SelectBuilder struct {
	From    string
	Columns []string
	Where   []string
	OrderBy []string
	Limit   interface{}
	Offset  interface{}
	Join    []string
	GroupBy []string
	Having  []string
}

// String generates the sql statement.
func (br SelectBuilder) String() string {
	sql := &bytes.Buffer{}

	sql.WriteString("SELECT ")
	if len(br.Columns) > 0 {
		sql.WriteString(strings.Join(br.Columns, ", "))
	} else {
		sql.WriteString("*")
	}

	if len(br.From) > 0 {
		sql.WriteString(" FROM ")
		sql.WriteString(br.From)
	}

	if len(br.Join) > 0 {
		sql.WriteString(" ")
		sql.WriteString(strings.Join(br.Join, " "))
	}

	if len(br.Where) > 0 {
		sql.WriteString(" WHERE ")
		sql.WriteString(strings.Join(br.Where, " AND "))
	}

	if len(br.GroupBy) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(br.GroupBy, ", "))
	}

	if len(br.Having) > 0 {
		sql.WriteString(" HAVING ")
		sql.WriteString(strings.Join(br.Having, " AND "))
	}

	if len(br.OrderBy) > 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(br.OrderBy, ", "))
	}

	if br.Limit != nil {
		sql.WriteString(fmt.Sprintf(" LIMIT %v", br.Limit))
	}
	if br.Offset != nil {
		sql.WriteString(fmt.Sprintf(" OFFSET %v", br.Offset))
	}

	return sql.String()
}

// InsertBuilder builds insert statement.
type InsertBuilder struct {
	Into   string
	Fields map[string]interface{}
}

// String generates the sql statement.
func (br InsertBuilder) String() string {
	sql := &bytes.Buffer{}

	sql.WriteString("INSERT INTO ")
	sql.WriteString(br.Into)

	cols, binds := builderFields(br.Fields)
	sql.WriteString(" (")
	sql.WriteString(strings.Join(cols, ", "))
	sql.WriteString(") VALUES (")
	sql.WriteString(strings.Join(binds, ", "))
	sql.WriteString(")")

	return sql.String()
}

// UpdateBuilder builds update statement.
type UpdateBuilder struct {
	From   string
	Fields map[string]interface{}
	Where  []string
}

// String generates the sql statement.
func (br UpdateBuilder) String() string {
	sql := &bytes.Buffer{}

	sql.WriteString("UPDATE ")
	sql.WriteString(br.From)

	sql.WriteString(" SET ")
	cols, binds := builderFields(br.Fields)
	for i := 0; i < len(br.Fields); i++ {
		if i > 0 {
			sql.WriteString(", ")
		}
		c, b := cols[i], binds[i]
		sql.WriteString(c + " = " + b)
	}

	if len(br.Where) > 0 {
		sql.WriteString(" WHERE ")
		sql.WriteString(strings.Join(br.Where, " AND "))
	}

	return sql.String()
}

// DeleteBuilder builds delete statement.
type DeleteBuilder struct {
	From    string
	Where   []string
	OrderBy []string
	Limit   interface{}
	Offset  interface{}
}

// String generates the sql statement.
func (br DeleteBuilder) String() string {
	sql := &bytes.Buffer{}

	sql.WriteString("DELETE FROM ")
	sql.WriteString(br.From)

	if len(br.Where) > 0 {
		sql.WriteString(" WHERE ")
		sql.WriteString(strings.Join(br.Where, " AND "))
	}

	if len(br.OrderBy) > 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(br.OrderBy, ", "))
	}

	if br.Limit != nil {
		sql.WriteString(fmt.Sprintf(" LIMIT %v", br.Limit))
	}
	if br.Offset != nil {
		sql.WriteString(fmt.Sprintf(" OFFSET %v", br.Offset))
	}

	return sql.String()
}

// builderFields splits fields map into columns and binding names
// and ensures the order consistency for testing.
func builderFields(m map[string]interface{}) ([]string, []string) {
	l := len(m)
	ks := make([]string, l)
	i := 0
	for k := range m {
		ks[i] = k
		i++
	}
	sort.Strings(ks)

	cols := make([]string, l)
	binds := make([]string, l)
	for i, k := range ks {
		cols[i] = k
		b := m[k]
		if b == nil || b.(string) == "" {
			binds[i] = fmt.Sprintf(":%s", k)
		} else {
			binds[i] = fmt.Sprintf("%v", b)
		}
	}
	return cols, binds
}

func NamespacedColumn(columnNames []string, tableName string) []string {
	if len(columnNames) == 0 {
		return []string{tableName + ".*"}
	}
	nsCols := make([]string, len(columnNames))
	for i, col := range columnNames {
		if strings.ContainsAny(col, ".") {
			nsCols[i] = col
		} else {
			nsCols[i] = tableName + "." + col
		}
	}
	return nsCols
}
