package main

import (
	"fmt"
	"github.com/araddon/dateparse"
)

func main() {
	t, err := dateparse.ParseAny("3600")
	if err != nil {
		println(err)
		return
	}

	println(fmt.Sprintf("%v", t))
}
