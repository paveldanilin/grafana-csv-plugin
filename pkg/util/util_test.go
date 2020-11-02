package util

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsNumber(t *testing.T) {
	if IsNumber("a") != false {
		t.Error("a: expected false")
	}
	if IsNumber("ThisIs0.45") != false {
		t.Error("ThisIs0.45: expected false")
	}
	if IsNumber("1") != true {
		t.Error("1: expected true")
	}
	if IsNumber("2.3") != true {
		t.Error("2.3: expected true")
	}
	if IsNumber(".3") != true {
		t.Error(".3: expected true")
	}
}

func TestIsInt(t *testing.T) {
	if IsInt(".45") != false {
		t.Error(".45: is not integer")
	}
	if IsInt("abc") != false {
		t.Error("abc: is not integer")
	}
	if IsInt("0.45") != false {
		t.Error("0.45: is not integer")
	}
	if IsInt("1589635810") != true {
		t.Error("1589635810: expected to be integer")
	}
}

func TestParseany(t *testing.T) {
	d, _ := dateparse.ParseAny("1595844722082")
	println(fmt.Sprintf("%v", d))
}

func TestStrToDur(t *testing.T) {
	dm := StrToDur("100000ms")
	fmt.Printf("%v", dm)
}

func TestInArray_False(t *testing.T) {
	r := InArray([]string{"a", "b"}, "c")

	assert.False(t, r)
}

func TestInArray_True(t *testing.T) {
	r := InArray([]string{"a", "b", "c"}, "b")

	assert.True(t, r)
}
