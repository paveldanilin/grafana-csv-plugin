package service

import (
	"github.com/paveldanilin/grafana-csv-plugin/pkg/csv"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDetectDatatypeString(t *testing.T) {
	strColType := detectDatatype("hello")
	assert.Equal(t, csv.ColumnTypeText, string(strColType))
}


func TestDetectDatatypeInteger(t *testing.T) {
	assert.Equal(t, csv.ColumnTypeInteger, string(detectDatatype("90")))
	assert.Equal(t, csv.ColumnTypeInteger, string(detectDatatype("0")))
	assert.Equal(t, csv.ColumnTypeInteger, string(detectDatatype("-90")))
	assert.Equal(t, csv.ColumnTypeInteger, string(detectDatatype("123456789")))
	assert.Equal(t, csv.ColumnTypeInteger, string(detectDatatype("1200")))
	assert.NotEqual(t, csv.ColumnTypeInteger, string(detectDatatype("0.34")))
}

func TestDetectDatatypeReal(t *testing.T) {
	assert.Equal(t, csv.ColumnTypeReal, string(detectDatatype("0.1")))
	assert.Equal(t, csv.ColumnTypeReal, string(detectDatatype("1.0")))
	assert.Equal(t, csv.ColumnTypeReal, string(detectDatatype("-0.01")))
}

func TestDetectDatatypeDate(t *testing.T) {
	assert.Equal(t, csv.ColumnTypeDate, string(detectDatatype("01.01.2020")))
	assert.Equal(t, csv.ColumnTypeDate, string(detectDatatype("2014-05-11 08:20:13")))
	assert.Equal(t, csv.ColumnTypeDate, string(detectDatatype("1499979795437000")))
}
