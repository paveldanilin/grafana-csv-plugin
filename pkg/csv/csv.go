package csv

import (
	"encoding/csv"
	"errors"
	"os"
)

type FileDescriptor struct {
	Filename string
	Delimiter rune
	Comment rune
	TrimLeadingSpace bool
	FieldsPerRecord int
	// User defined or auto detected info about columns
	Columns []Column
}

// Deprecated
func Read(descriptor *FileDescriptor) (*Table, error) {
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

	csvFile := NewTable()
	if len(records) == 0 {
		return csvFile, nil
	}

	csvFile.AddColumns(records[0]...)

	for i := 1 ; i < len(records); i++ {
		csvFile.AddRow(records[i]...)
	}

	return csvFile, nil
}

func Read2(descriptor *FileDescriptor) (*Table2, error) {
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

	db, err := ToDb(csvReader, descriptor)
	if err != nil {
		return nil, err
	}

	return &Table2{
		db: db,
		columns: descriptor.Columns,
	}, nil
}
