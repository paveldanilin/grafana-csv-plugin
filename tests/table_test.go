package tests

import (
	"github.com/paveldanilin/grafana-csv-plugin/pkg/csv"
	"testing"
)

func TestTableAddColumns(t *testing.T) {
	table := csv.NewTable()

	table.AddColumns("ID", "Name", "Sename")

	if len(table.Columns) != 3 {
		t.Errorf("Expected 3 columns, but got %d", len(table.Columns))
	}
}

func TestTableAddRow(t *testing.T) {
	table := csv.NewTable()

	table.AddColumns("ID", "Name", "Date")

	table.AddRow("1", "Pavel", "12.05.2020")
	table.AddRow("2", "Ivan", "12.05.2020")

	val, err := table.GetValue("Name", 1)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if val != "Ivan" {
		t.Errorf("Expected 'Ivan' string, but got '%v'", val)
	}

	date, err := table.GetValue("Date", 0)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if date != "12.05.2020" {
		t.Errorf("Expected '12.05.2020' string, but got '%v'", val)
	}

	id, err := table.GetInt64("ID", 1)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if id != 2 {
		t.Errorf("Expected 1 int, but got '%v'", id)
	}
}

func TestTableNumericFilter(t *testing.T) {
	table := csv.NewTable()

	table.AddColumn("A")
	table.AddColumn("B")
	table.AddColumn("C")

	table.AddRow("1", "2", "3")
	table.AddRow("4", "5", "6")
	table.AddRow("7", "8", "9")

	filteredTable, err := table.Filter("A > 2 and B == 5")
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if filteredTable.RowsCount() != 1 {
		t.Error("Expected exactly 1 row in filtered table")
	}
}

func TestTableAddRowNil(t *testing.T) {
	table := csv.NewTable()

	table.AddColumn("ID")
	table.AddColumn("Name")

	table.AddRow("1")
	table.AddRow("2")

	val, err := table.GetValue("Name", 0)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if val != nil {
		t.Error("Expected nil value")
	}
}
