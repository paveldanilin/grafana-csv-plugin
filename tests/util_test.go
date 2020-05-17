package tests

import (
	"github.com/paveldanilin/grafana-csv-plugin/pkg/util"
	"testing"
)

func TestIsNumber(t *testing.T) {
	if util.IsNumber("a") != false {
		t.Error("a: expected false")
	}
	if util.IsNumber("ThisIs0.45") != false {
		t.Error("ThisIs0.45: expected false")
	}
	if util.IsNumber("1") != true {
		t.Error("1: expected true")
	}
	if util.IsNumber("2.3") != true {
		t.Error("2.3: expected true")
	}
	if util.IsNumber(".3") != true {
		t.Error(".3: expected true")
	}
}

func TestIsInt(t *testing.T) {
	if util.IsInt(".45") != false {
		t.Error(".45: is not integer")
	}
	if util.IsInt("abc") != false {
		t.Error("abc: is not integer")
	}
	if util.IsInt("0.45") != false {
		t.Error("0.45: is not integer")
	}
	if util.IsInt("1589635810") != true {
		t.Error("1589635810: expected to be integer")
	}
}
