package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

type FileDescriptor struct {
	Filename string
	Delimiter rune
	Comment rune
	TrimLeadingSpace bool
	FieldsPerRecord int
}

type File struct {
	Columns []string
	Rows    [][]string
}

func (c *File) RowsCount() int {
	return len(c.Rows)
}

func (c *File) ColumnIndex(column string) int {
	for i, columnName := range c.Columns {
		if columnName == column {
			return i
		}
	}
	return -1
}

func (c *File) HasColumn(column string) bool {
	return c.ColumnIndex(column) != -1
}

func (c *File) GetRow(rowIndex int) ([]string, error) {
	if rowIndex > c.RowsCount() {
		return make([]string, 0), errors.New(fmt.Sprintf("the index [%d] is out of range", rowIndex))
	}
	return c.Rows[rowIndex], nil
}

// Returns nil when either there is no column with such name or there is no row with such index
func (c *File) GetColumnValue(column string, rowIndex int) interface{} {
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

func Read(descriptor *FileDescriptor) (*File, error) {
	if descriptor == nil {
		return nil, errors.New("file descriptor is missed")
	}

	file, err := os.Open(descriptor.Filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = descriptor.Delimiter
	csvReader.Comment = descriptor.Comment
	csvReader.TrimLeadingSpace = descriptor.TrimLeadingSpace
	csvReader.FieldsPerRecord = descriptor.FieldsPerRecord

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	csvFile := &File{
		Columns: make([]string, 0),
		Rows:    make([][]string, 0),
	}

	if len(records) == 0 {
		return csvFile, nil
	}

	// Build Columns
	headers := records[0]
	for _, headerName := range headers {
		csvFile.Columns = append(csvFile.Columns, headerName)
	}

	csvFile.Rows = records[1:]

	return csvFile, nil
}
