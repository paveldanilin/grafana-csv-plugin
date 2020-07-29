package time_group

import (
	"errors"
	"fmt"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
	"strconv"
)

// $__timeGroup(dateColumn, interval)
const MacroName = "timeGroup"

func Processor(args []string, scope *macro.Scope) (string, error) {
	if len(args) != 2 {
		return "", errors.New(fmt.Sprintf("%s(dateColumn, interval): expected two arguments, but got %d", MacroName, len(args)))
	}

	dateColumn := args[0]
	// m = minutes
	// 1m, 5m, 10m
	// h = hour
	// 1h, 2h, 24h
	intervalExpr := args[1]
	var interval int64

	if intervalExpr[len(intervalExpr)-1] == 'm' {
		min,_ := strconv.ParseInt(intervalExpr[:len(intervalExpr) - 1], 10, 64)
		interval = min * 60
	} else if  intervalExpr[len(intervalExpr)-1] == 'h' {
		hour,_ := strconv.ParseInt(intervalExpr[:len(intervalExpr) - 1], 10, 64)
		interval = hour * 60 * 60
	} else {
		interval, _ = strconv.ParseInt(intervalExpr, 10, 64)
	}

	return fmt.Sprintf(
		"datetime((strftime('%%s', %s) / %d) * %d, 'unixepoch')",
		dateColumn,
		interval,
		interval,
	), nil
}
