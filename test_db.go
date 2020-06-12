package main

import (
	"fmt"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/csv"
	"io"
	"runtime"
	"time"
)

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func main() {
	desc := &csv.FileDescriptor{
		Filename:         "./data/SacramentocrimeJanuary2006.csv",
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
	sql := "select * from DataTable"

	cols, rows, err := table.Query(sql)
	if err != nil {
		panic(err)
	}


	for _, row := range *rows {
		//println(fmt.Sprintf("%v", row))
		for _, value := range row {
			if _, ok := value.(time.Time); ok {
				switch typedValue := value.(type) {
				case time.Time:
					println(typedValue.Unix())
					break
				}
			} else {
				//println(fmt.Sprintf("%v=%v", value, reflect.TypeOf(value)))
			}
		}
	}

	println(fmt.Sprintf("%v", cols))

	PrintMemUsage()
	runtime.GC()
	PrintMemUsage()

	r, err := table.Query2(sql)
	if err != nil {
		panic(err)
	}

	for {
		_, err := r.Next()
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}

		//fmt.Printf("%v", row)
	}
	r.Release()

	PrintMemUsage()
}
