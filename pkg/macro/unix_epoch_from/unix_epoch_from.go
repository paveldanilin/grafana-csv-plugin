package unix_epoch_from

import (
	"fmt"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
)

// $__unixEpochFrom()
const MacroName = "unixEpochFrom"

func Processor(args []string, scope *macro.Scope) (string, error) {
	return fmt.Sprintf(
		"%d",
		scope.GetVar("timeFromMs"),
	), nil
}
