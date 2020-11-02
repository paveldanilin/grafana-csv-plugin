package xsql

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"testing"
	"time"
)

func getTestDb() *sql.DB {
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE contacts (\ncontact_id INTEGER PRIMARY KEY,\nfirst_name TEXT NOT NULL,\nphone TEXT NOT NULL UNIQUE\n)")
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("INSERT INTO contacts VALUES(?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, _ = stmt.Exec(1, "Pasha", "999")
	_, _ = stmt.Exec(2, "Dima", "333")
	_, _ = stmt.Exec(3, "Lesha", "444")

	return db
}

func TestNewQueryResult(t *testing.T) {
	db := getTestDb()

	rows, err := db.Query("SELECT * FROM contacts")
	if err != nil {
		panic(err)
	}

	r, err := NewQueryResult(rows, []GeneratedColumn{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", r.Columns())

	for {
		row, err := r.Next()
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		fmt.Printf("%v\n", row)
	}
}

func TestQueryResult_GeneratedColumns(t *testing.T) {
	db := getTestDb()

	rows, err := db.Query("SELECT * FROM contacts")
	if err != nil {
		panic(err)
	}

	r, err := NewQueryResult(rows, []GeneratedColumn{
		{
			Name: "$autotime",
			Func: func(values []interface{}) interface{} {
				return time.Now()
			},
			Type: nil,
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v\n", r.Columns())

	for {
		row, err := r.Next()
		if err != nil && err != io.EOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}
		fmt.Printf("%v\n", row)
	}
}
