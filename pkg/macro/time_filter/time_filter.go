package time_filter

import (
	"errors"
	"fmt"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
)

// $__timeFilter(dateColumn)
const MacroName = "timeFilter"

func Processor(args []string, scope *macro.Scope) (string, error) {
	if len(args) != 1 {
		return "", errors.New("not found date column")
	}
	return fmt.Sprintf(
		" %s BETWEEN '%s' AND '%s' ",
		// Table name
		args[0],
		scope.GetVar("timeFrom"),
		scope.GetVar("timeTo"),
	), nil
}
