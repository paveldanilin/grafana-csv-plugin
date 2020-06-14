package tests

import (
	"github.com/paveldanilin/grafana-csv-plugin/pkg/csv"
	"testing"
)

func TestCsvReadSimple(t *testing.T) {
	f, err := csv.NewRead(&csv.FileDescriptor{
		Filename:         "../data/test.csv",
		Delimiter:        ',',
		Comment:          '#',
		TrimLeadingSpace: true,
		FieldsPerRecord:  0,
	})

	if err != nil {
		t.Error(err)
	}

	if len(f.Columns) != 3 {
		t.Errorf("csv.NewRead(test.csv) FAILED, expected %d columns, but got %d", 3, len(f.Columns))
	}

	if f.RowsCount() != 3 {
		t.Errorf("csv.NewRead(test.csv) FAILED, expected %d columns, but got %d", 3, f.RowsCount())
	}
}

func TestCsvReadCrime2006(t *testing.T) {
	f, err := csv.NewRead(&csv.FileDescriptor{
		Filename:         "../data/SacramentocrimeJanuary2006.csv",
		Delimiter:        ',',
		Comment:          '#',
		TrimLeadingSpace: true,
		FieldsPerRecord:  0,
	})

	if err != nil {
		t.Error(err)
	}

	if len(f.Columns) != 9 {
		t.Errorf("csv.NewRead(SacramentocrimeJanuary2006.csv) FAILED, expected %d columns, but got %d", 9, len(f.Columns))
	}
}
