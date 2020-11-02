package time_group

import (
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimeGroup(t *testing.T) {
	macro.Register(MacroName, Processor)

	// min
	minSql, _, err := macro.Interpolate("SELECT $__timeGroup(order_date, 1m) FROM my_table GROUP BY interval", nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "SELECT datetime((strftime('%s', order_date) / 60) * 60, 'unixepoch') FROM my_table GROUP BY interval", minSql)

	// hour
	hourSql, _, err := macro.Interpolate("SELECT $__timeGroup(order_date, 24h) FROM my_table GROUP BY interval", nil)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, "SELECT datetime((strftime('%s', order_date) / 86400) * 86400, 'unixepoch') FROM my_table GROUP BY interval", hourSql)
}
