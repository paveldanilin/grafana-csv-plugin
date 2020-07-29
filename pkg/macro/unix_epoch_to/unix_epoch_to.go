package unix_epoch_to

import (
	"fmt"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
)

// $__unixEpochTo()
const MacroName = "unixEpochTo"

func Processor(args []string, scope *macro.Scope) (string, error) {
	return fmt.Sprintf(
		"%d",
		scope.GetVar("timeToMs"),
	), nil
}
