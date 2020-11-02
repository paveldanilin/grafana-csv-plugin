package csv

import (
	"encoding/csv"
	"errors"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
	"os"
)

type FileDescriptor struct {
	Filename string
	fileSize int64
	fileModTime int64
	Delimiter rune
	Comment rune
	TrimLeadingSpace bool
	FieldsPerRecord int
	// User defined or auto detected info about columns
	Columns []Column
}

func (d *FileDescriptor) GetFileSize() int64 {
	return d.fileSize
}

func (d *FileDescriptor) GetFileModTime() int64 {
	return d.fileModTime
}

type reader struct {
	file *os.File
	csv  *csv.Reader
}

func NewCsvReader(descriptor *FileDescriptor) (*reader, error) {
	if descriptor == nil {
		return nil, errors.New("file descriptor is missed")
	}

	descriptor.fileSize, descriptor.fileModTime = util.FileStat(descriptor.Filename)

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

func (r *reader) Read() ([]string, error){
	return r.csv.Read()
}

func (r *reader) Close() {
	_ = r.file.Close()
	r.file = nil
	r.csv = nil
}
