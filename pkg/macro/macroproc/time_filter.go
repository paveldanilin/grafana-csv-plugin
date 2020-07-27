package macroproc

import (
	"fmt"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
)

func TimeFilter(args []string, scope *macro.Scope) (string, error) {
	return fmt.Sprintf(
		" %s BETWEEN '%s' AND '%s' ",
		// Table name
		args[0],
		scope.GetVar("timeFrom"),
		scope.GetVar("timeTo"),
	), nil

}
