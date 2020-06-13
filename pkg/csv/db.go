package csv

import (
	"database/sql"
)

type ColumnType string

const (
	ColumnTypeText = "text"
	ColumnTypeInteger = "integer"
	ColumnTypeReal = "real"
	ColumnTypeTimestamp = "timestamp"
	ColumnTypeDate = "date"
)

type Column struct {
	Type ColumnType
	Name string
}

type DB struct {
	columns []Column
	db *sql.DB
}

func (t *DB) Query(sql string) (*QueryResult, error) {
	rows, err := t.db.Query(sql)
	if err != nil {
		return nil, err
	}
	return newQueryResult(rows)
}

func (t *DB) Release() {
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
