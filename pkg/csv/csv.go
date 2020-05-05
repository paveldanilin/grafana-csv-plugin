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
}

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

	csvFile := &Table{
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
