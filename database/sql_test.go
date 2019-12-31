package database

import "testing"

func TestSelectBuilder(t *testing.T) {
	tests := []struct {
		name  string
		br    SelectBuilder
		wants string
	}{
		{
			name: "select",
			br: SelectBuilder{
				From:    "table",
				Columns: []string{"a", "b", "c"},
			},
			wants: "SELECT a, b, c FROM table",
		},
		{
			name: "order by",
			br: SelectBuilder{
				From:    "table",
				Columns: []string{"d", "e", "f"},
				OrderBy: []string{"d ASC", "e DESC"},
			},
			wants: "SELECT d, e, f FROM table ORDER BY d ASC, e DESC",
		},
		{
			name: "limit (int)",
			br: SelectBuilder{
				From:    "table",
				Columns: []string{"d", "e", "f"},
				Limit:   1,
			},
			wants: "SELECT d, e, f FROM table LIMIT 1",
		},
		{
			name: "limit (string)",
			br: SelectBuilder{
				From:    "table",
				Columns: []string{"d", "e", "f"},
				Limit:   "1",
			},
			wants: "SELECT d, e, f FROM table LIMIT 1",
		},
		{
			name: "limit offset",
			br: SelectBuilder{
				From:    "table",
				Columns: []string{"g"},
				Limit:   1,
				Offset:  2,
			},
			wants: "SELECT g FROM table LIMIT 1 OFFSET 2",
		},
		{
			name: "join",
			br: SelectBuilder{
				From:    "table",
				Columns: []string{"h"},
				Join:    []string{"JOIN sth ON table.id = sth.fid"},
			},
			wants: "SELECT h FROM table JOIN sth ON table.id = sth.fid",
		},
		{
			name: "where",
			br: SelectBuilder{
				From:    "table",
				Columns: []string{"a", "b", "c", "d"},
				Where:   []string{"a = $1", "b IS NULL", "c = :c", "d = ?"},
			},
			wants: "SELECT a, b, c, d FROM table WHERE a = $1 AND b IS NULL AND c = :c AND d = ?",
		},
		{
			name: "escape",
			br: SelectBuilder{
				From:    "table",
				Columns: []string{"`a`", `"b"`},
			},
			wants: "SELECT `a`, \"b\" FROM table",
		},
		{
			name: "group by",
			br: SelectBuilder{
				From:    "table",
				Columns: []string{"`a`", `"b"`},
				GroupBy: []string{"a"},
			},
			wants: "SELECT `a`, \"b\" FROM table GROUP BY a",
		},
	}

	for _, test := range tests {
		got := test.br.String()
		if got != test.wants {
			t.Errorf("case %q fails, got %q", test.name, got)
		}
	}
}

func TestInsertBuilder(t *testing.T) {
	tests := []struct {
		name  string
		br    InsertBuilder
		wants string
	}{
		{
			name: "insert",
			br: InsertBuilder{
				Into:   "table",
				Fields: map[string]interface{}{"a": nil},
			},
			wants: "INSERT INTO table (a) VALUES (:a)",
		},
		{
			name: "bind",
			br: InsertBuilder{
				Into: "table",
				Fields: map[string]interface{}{
					"mysql": "?",
					"pq":    "$1",
				},
			},
			wants: "INSERT INTO table (mysql, pq) VALUES (?, $1)",
		},
		{
			name: "escape",
			br: InsertBuilder{
				Into:   "table",
				Fields: map[string]interface{}{`"escape"`: ":escape"},
			},
			wants: "INSERT INTO table (\"escape\") VALUES (:escape)",
		},
	}

	for _, test := range tests {
		got := test.br.String()
		if got != test.wants {
			t.Errorf("case %q fails, got %q", test.name, got)
		}
	}
}

func TestUpdateBuilder(t *testing.T) {
	tests := []struct {
		name  string
		br    UpdateBuilder
		wants string
	}{
		{
			name: "update",
			br: UpdateBuilder{
				From:   "table",
				Fields: map[string]interface{}{"a": nil},
			},
			wants: "UPDATE table SET a = :a",
		},
		{
			name: "bind",
			br: UpdateBuilder{
				From: "table",
				Fields: map[string]interface{}{
					"mysql": "?",
					"pq":    "$1",
				},
			},
			wants: "UPDATE table SET mysql = ?, pq = $1",
		},
		{
			name: "where",
			br: UpdateBuilder{
				From:   "table",
				Fields: map[string]interface{}{"a": nil},
				Where:  []string{"b = :b"},
			},
			wants: "UPDATE table SET a = :a WHERE b = :b",
		},
		{
			name: "escape",
			br: UpdateBuilder{
				From:   "table",
				Fields: map[string]interface{}{`"a"`: ":escape"},
			},
			wants: "UPDATE table SET \"a\" = :escape",
		},
	}

	for _, test := range tests {
		got := test.br.String()
		if got != test.wants {
			t.Errorf("case %q fails, got %q", test.name, got)
		}
	}
}

func TestDeleteBuilder(t *testing.T) {
	tests := []struct {
		name  string
		br    DeleteBuilder
		wants string
	}{
		{
			name: "delete",
			br: DeleteBuilder{
				From: "table",
			},
			wants: "DELETE FROM table",
		},
		{
			name: "where",
			br: DeleteBuilder{
				From:  "table",
				Where: []string{"a = :a"},
			},
			wants: "DELETE FROM table WHERE a = :a",
		},
		{
			name: "limit",
			br: DeleteBuilder{
				From:  "table",
				Limit: 1,
			},
			wants: "DELETE FROM table LIMIT 1",
		},
		{
			name: "limit offset",
			br: DeleteBuilder{
				From:   "table",
				Limit:  1,
				Offset: 2,
			},
			wants: "DELETE FROM table LIMIT 1 OFFSET 2",
		},
		{
			name: "offset",
			br: DeleteBuilder{
				From:   "table",
				Where:  []string{"1 = 1"},
				Offset: 2,
			},
			wants: "DELETE FROM table WHERE 1 = 1 OFFSET 2",
		},
	}

	for _, test := range tests {
		got := test.br.String()
		if got != test.wants {
			t.Errorf("case %q fails, got %q", test.name, got)
		}
	}
}
