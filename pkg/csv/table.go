package csv

import (
	"errors"
	"fmt"
	"github.com/antonmedv/expr"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
	"strconv"
	"strings"
)

type Table struct {
	Columns []string
	Rows    [][]interface{}
}

func NewTable() *Table {
	return &Table{
		Columns: make([]string, 0),
		Rows:    make([][]interface{}, 0),
	}
}

func (t *Table) String() string {
	return fmt.Sprintf("Columns=%d, Rows=%d", len(t.Columns), len(t.Rows))
}

func (t *Table) AddColumn(column string) {
	t.Columns = append(t.Columns, column)
}

func (t *Table) AddColumns(column ...string) {
	for _, name := range column {
		t.Columns = append(t.Columns, name)
	}
}

func (t *Table) RowsCount() int {
	return len(t.Rows)
}

func (t *Table) ColumnIndex(column string) int {
	for i, columnName := range t.Columns {
		if columnName == column {
			return i
		}
	}
	return -1
}

func (t *Table) HasColumn(column string) bool {
	return t.ColumnIndex(column) != -1
}

func (t *Table) GetRow(rowIndex int) ([]interface{}, error) {
	if rowIndex > t.RowsCount() {
		return make([]interface{}, 0), errors.New(fmt.Sprintf("the index [%d] is out of range", rowIndex))
	}
	return t.Rows[rowIndex], nil
}

func (t *Table) AddRow(values ...string) {
	cols := len(t.Columns)
	vlen := len(values)
	newRow := make([]interface{}, cols)

	for i := 0; i < cols; i++ {
		if i > vlen - 1 {
			newRow[i] = nil
			continue
		}
		v := strings.TrimSpace(values[i])
		if util.IsNumber(v) {
			if util.IsInt(v) {
				newRow[i], _ = strconv.ParseInt(v, 10, 64)
			} else {
				newRow[i], _ = strconv.ParseFloat(v, 64)
			}
		} else {
			newRow[i] = v
		}
	}

	t.Rows = append(t.Rows, newRow)
}

// Returns nil when either there is no column with such name or there is no row with such index
func (t *Table) GetValue(column string, rowIndex int) (interface{}, error) {
	columnIndex := t.ColumnIndex(column)
	if columnIndex == -1 {
		return nil, errors.New(fmt.Sprintf("unknown column `%s`", column))
	}

	row, err := t.GetRow(rowIndex)
	if err != nil {
		return nil, err
	}

	return row[columnIndex], nil
}

func (t *Table) GetInt64(column string, rowIndex int) (int64, error) {
	val, err := t.GetValue(column, rowIndex)
	if err != nil {
		return 0, err
	}
	return val.(int64), nil
}

func (t *Table) GetFloat64(column string, rowIndex int) (float64, error) {
	val, err := t.GetValue(column, rowIndex)
	if err != nil {
		return 0, err
	}
	return val.(float64), nil
}

func (t *Table) GetString(column string, rowIndex int) (string, error) {
	val, err := t.GetValue(column, rowIndex)
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

// https://github.com/antonmedv/expr/blob/master/docs/Language-Definition.md
func (t *Table) Filter(filterExpr string) (*Table, error) {
	filtered := NewTable()
	for _, col := range t.Columns {
		filtered.AddColumn(col)
	}

	for i, row := range t.Rows {
		namedRow, err := t.toNamedRow(i)
		if err != nil {
			return nil, err
		}

		result, err := expr.Eval(filterExpr, namedRow)
		if err != nil {
			return nil, err
		}

		if result.(bool) == true {
			filtered.Rows = append(filtered.Rows, row)
		}
	}

	return filtered, nil
}

func (t *Table) toNamedRow(rowIndex int) (map[string]interface{}, error) {
	namedRow := make(map[string]interface{})
	row, err := t.GetRow(rowIndex)
	if err != nil {
		return nil, err
	}
	for i, colName := range t.Columns {
		namedRow[colName] = row[i]
	}
	return namedRow, nil
}
