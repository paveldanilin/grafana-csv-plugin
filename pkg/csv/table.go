package csv

import (
	"errors"
	"fmt"
)

type Table struct {
	Columns []string
	Rows    [][]string
}

func (c *Table) RowsCount() int {
	return len(c.Rows)
}

func (c *Table) ColumnIndex(column string) int {
	for i, columnName := range c.Columns {
		if columnName == column {
			return i
		}
	}
	return -1
}

func (c *Table) HasColumn(column string) bool {
	return c.ColumnIndex(column) != -1
}

func (c *Table) GetRow(rowIndex int) ([]string, error) {
	if rowIndex > c.RowsCount() {
		return make([]string, 0), errors.New(fmt.Sprintf("the index [%d] is out of range", rowIndex))
	}
	return c.Rows[rowIndex], nil
}

// Returns nil when either there is no column with such name or there is no row with such index
func (c *Table) GetColumnValue(column string, rowIndex int) interface{} {
	columnIndex := c.ColumnIndex(column)
	if columnIndex == -1 {
		return nil
	}

	row, err := c.GetRow(rowIndex)
	if err != nil {
		return nil
	}

	return row[columnIndex]
}
