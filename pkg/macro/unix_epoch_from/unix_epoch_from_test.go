package unix_epoch_from

import (
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnixEpochFrom(t *testing.T) {
	macro.Register(MacroName, Processor)

	scope := macro.NewScope()
	scope.SetVar("timeFromMs", 1595844722082)

	text, _, err := macro.Interpolate("SELECT * FROM my_table WHERE order_date >= $__unixEpochFrom()", scope)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "SELECT * FROM my_table WHERE order_date >= 1595844722082", text)
}
