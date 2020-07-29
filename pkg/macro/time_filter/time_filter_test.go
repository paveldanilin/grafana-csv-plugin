package time_filter

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimeFilter(t *testing.T) {
	macro.Register(MacroName, Processor)

	datetimeFrom, _ := dateparse.ParseAny(fmt.Sprintf("%d", 1595844722082))
	datetimeTo, _ := dateparse.ParseAny(fmt.Sprintf("%d", 1595866322082))

	scope := macro.NewScope()
	scope.SetVar("timeFrom", datetimeFrom)
	scope.SetVar("timeTo", datetimeTo)

	text, err := macro.Interpolate("SELECT * FROM my_table WHERE $__timeFilter(DATE_COLUMN) AND status=0", scope)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "SELECT * FROM my_table WHERE  DATE_COLUMN BETWEEN '2020-07-27 20:12:02.082 +1000 +10' AND '2020-07-28 02:12:02.082 +1000 +10'  AND status=0", text)
}
