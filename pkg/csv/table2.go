package csv

import (
	"database/sql"
	"io"
)

type ColumnType string

const (
	ColumnTypeText = "text"
	ColumnTypeInteger = "integer"
	ColumnTypeReal = "real"
	ColumnTypeTimestamp = "timestamp"
	ColumnTypeDate = "date"
)

type QueryResult struct {
	rows *sql.Rows
	columns []string
	ptrs []interface{}
	vals []interface{}
}

func newQueryResult(rows *sql.Rows) (*QueryResult, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	r := &QueryResult{
		rows: rows,
		columns: columns,
	}
	r.ptrs = make([]interface{}, len(columns))
	r.vals = make([]interface{}, len(columns))
	for i, _ := range r.ptrs {
		r.ptrs[i] = &r.vals[i]
	}
	return r, nil
}

func (r *QueryResult) Columns() ([]string, error) {
	return r.rows.Columns()
}

func (r *QueryResult) Next() ([]interface{}, error) {
	if ok := r.rows.Next(); !ok {
		return nil, io.EOF
	}
	err := r.rows.Scan(r.ptrs...)
	if err != nil {
		return nil, err
	}
	return r.vals, nil
}

func (r *QueryResult) Release() {
	r.rows.Close()
	r.columns = nil
	r.ptrs = nil
	r.vals = nil
}

type Column struct {
	Type ColumnType
	Name string
}

type Table2 struct {
	columns []Column
	db *sql.DB
}

func (t *Table2) Query2(sql string) (*QueryResult, error) {
	rows, err := t.db.Query(sql)
	if err != nil {
		return nil, err
	}
	return newQueryResult(rows)
}

// Deprecated
func (t *Table2) Query(sql string) ([]string, *[][]interface{}, error) {
	rows, err := t.db.Query(sql)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	resultColumns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	resultRows := make([][]interface{}, 0)

	for rows.Next() {
		c, err := rows.Columns()
		if err != nil {
			return nil, nil, err
		}

		pointers := make([]interface{}, len(c))
		curResultRow := len(resultRows)
		resultRows = append(resultRows, make([]interface{}, len(c)))

		for i, _ := range pointers {
			pointers[i] = &resultRows[curResultRow][i]
		}

		err = rows.Scan(pointers...)
		if err != nil {
			return nil, nil, err
		}
	}

	return resultColumns, &resultRows, nil
}

func (t *Table2) Destroy() {
	if t.db != nil {
		_ = t.db.Close()
		t.db = nil
	}
}

func ColumnTypeFromString(s string) ColumnType {
	switch s {
	case "text":
		return ColumnTypeText
	case "integer":
		return ColumnTypeInteger
	case "real":
		return ColumnTypeReal
	case "date":
		return ColumnTypeDate
	case "timestamp":
		return ColumnTypeTimestamp
	}
	return ""
}
