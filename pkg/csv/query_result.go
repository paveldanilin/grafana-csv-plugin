package csv

import (
	"database/sql"
	"io"
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

func (r *QueryResult) ColumnTypes() ([]*sql.ColumnType, error) {
	return r.rows.ColumnTypes()
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
