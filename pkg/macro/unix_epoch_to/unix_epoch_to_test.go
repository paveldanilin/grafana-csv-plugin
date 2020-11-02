package unix_epoch_to

import (
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnixEpochTo(t *testing.T) {
	macro.Register(MacroName, Processor)

	scope := macro.NewScope()
	scope.SetVar("timeToMs", 1595844722082)

	text, _, err := macro.Interpolate("SELECT * FROM my_table WHERE order_date >= $__unixEpochTo()", scope)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "SELECT * FROM my_table WHERE order_date >= 1595844722082", text)
}
