package main

import (
	"fmt"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/csv"
)

func main() {
	desc := &csv.FileDescriptor{
		Filename:         "./data/dtp_opendata.csv",
		Delimiter:        ',',
		Comment:          '#',
		TrimLeadingSpace: false,
		FieldsPerRecord:  0,
	}

	table, err := csv.Read2(desc)

	if err != nil {
		panic(err)
	}

	//sql := "Select * FROM DataTable where Price = 3600 and Transaction_date BETWEEN date('2009-01-01') AND date('2009-01-15')  order by Transaction_date"
	//sql := "select * from DataTable where cdatetime BETWEEN date('2006-01-01') AND date('2006-01-02')  order by cdatetime"
	sql := "select * from DataTable where message_id = 2467"

	cols, rows, err := table.Query(sql)
	if err != nil {
		panic(err)
	}


	for _, row := range rows {
		println(fmt.Sprintf("%v", row))
	}

	println(fmt.Sprintf("%v", cols))
}
