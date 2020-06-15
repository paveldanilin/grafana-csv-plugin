package csv

import (
	"encoding/csv"
	"errors"
	"os"
	"time"
)

const TableName = "DataTable"

type FileDescriptor struct {
	Filename string
	fileSize int64
	fileModTime time.Time
	Delimiter rune
	Comment rune
	TrimLeadingSpace bool
	FieldsPerRecord int
	// User defined or auto detected info about columns
	Columns []Column
}

type reader struct {
	file *os.File
	csv  *csv.Reader
}

func newCsvReader(descriptor *FileDescriptor) (*reader, error) {
	if descriptor == nil {
		return nil, errors.New("file descriptor is missed")
	}

	fileStat, err := os.Stat(descriptor.Filename)
	if err != nil {
		return nil, err
	}
	descriptor.fileSize = fileStat.Size()
	descriptor.fileModTime = fileStat.ModTime()

	file, err := os.Open(descriptor.Filename)
	if err != nil {
		return nil, err
	}

	csvReader := csv.NewReader(file)
	csvReader.Comma = descriptor.Delimiter
	csvReader.Comment = descriptor.Comment
	csvReader.TrimLeadingSpace = descriptor.TrimLeadingSpace
	csvReader.FieldsPerRecord = descriptor.FieldsPerRecord

	return &reader{
		file: file,
		csv:  csvReader,
	}, nil
}

func (r *reader) close() {
	r.file.Close()
	r.file = nil
	r.csv = nil
}
