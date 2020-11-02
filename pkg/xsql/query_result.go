package xsql

import (
	"database/sql"
	"io"
)

type ColumnValueGenerator func (values []interface{}) interface{}

type GeneratedColumn struct {
	Name    string
	Func    ColumnValueGenerator
	Type    *sql.ColumnType
}

type QueryResult struct {
	rows             *sql.Rows
	columnNames      []string
	columnTypes      []*sql.ColumnType
	ptrs             []interface{}
	vals             []interface{}
	generatedColumns []GeneratedColumn
}

func NewQueryResult(rows *sql.Rows, generatedColumns []GeneratedColumn) (*QueryResult, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	generatedColumnNames := func() []string {
		names := make([]string, 0)
		for _, gc := range generatedColumns {
			names = append(names, gc.Name)
		}
		return names
	}

	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	generatedColumnTypes := func () []*sql.ColumnType {
		types := make([]*sql.ColumnType, 0)
		for _, gt := range generatedColumns {
			types = append(types, gt.Type)
		}
		return types
	}

	r := &QueryResult{
		rows:             rows,
		columnNames:      append(columns, generatedColumnNames()...),
		columnTypes:      append(types, generatedColumnTypes()...),
		generatedColumns: generatedColumns,
	}

	r.ptrs = make([]interface{}, len(columns))
	r.vals = make([]interface{}, len(r.columnNames))

	for i, _ := range r.ptrs {
		r.ptrs[i] = &r.vals[i]
	}
	return r, nil
}

func (r *QueryResult) Columns() []string {
	return r.columnNames
}

func (r *QueryResult) ColumnTypes() []*sql.ColumnType {
	return r.columnTypes
}

func (r *QueryResult) Next() ([]interface{}, error) {
	if ok := r.rows.Next(); !ok {
		return nil, io.EOF
	}
	err := r.rows.Scan(r.ptrs...)
	if err != nil {
		return nil, err
	}
	valGenIndex := 0
	for i := len(r.ptrs); i < len(r.vals); i++ {
		r.vals[i] = r.generatedColumns[valGenIndex].Func(r.vals)
		valGenIndex++
	}
	return r.vals, nil
}

func (r *QueryResult) Release() {
	_ = r.rows.Close()
	r.columnNames = nil
	r.columnTypes = nil
	r.generatedColumns = nil
	r.ptrs = nil
	r.vals = nil
}
