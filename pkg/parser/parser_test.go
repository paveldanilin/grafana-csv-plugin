package parser

import "testing"

func TestSqlParser_Parse(t *testing.T) {
	sqlp := NewSqlParser()

	sqlp.Parse("select * where price > 10")
}
