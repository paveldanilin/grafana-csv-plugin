package csv

import (
	"github.com/paveldanilin/grafana-csv-plugin/pkg/xsql"
)

const (
	ColumnTypeText = "text"
	ColumnTypeInteger = "integer"
	ColumnTypeReal = "real"
	ColumnTypeTimestamp = "timestamp"
	ColumnTypeDate = "date"
)

type ColumnType string

type Column struct {
	Type ColumnType
	Name string
}

type DB interface {
	Init() error
	Query(sql string, generatedColumns []xsql.GeneratedColumn) (*xsql.QueryResult, error)
	LoadCSV(tableName string, descriptor *FileDescriptor) error
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
