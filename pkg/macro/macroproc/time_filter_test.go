package macroproc

import (
	"fmt"
	"github.com/araddon/dateparse"
	"github.com/paveldanilin/grafana-csv-plugin/pkg/macro"
	"testing"
)

func TestTimeFilter(t *testing.T) {
	macro.Register("timeFilter", TimeFilter)

	datetimeFrom, _ := dateparse.ParseAny(fmt.Sprintf("%d", 1595844722082))
	datetimeTo, _ := dateparse.ParseAny(fmt.Sprintf("%d", 1595866322082))


	scope := macro.NewScope()
	scope.SetVar("timeFrom", datetimeFrom)
	scope.SetVar("timeTo", datetimeTo)

	text, err := macro.Interpolate("SELECT * FROM $__timeFilter(DATE_COLUMN)", scope)

	if err != nil {
		t.Error(err)
	}

	println(text)
}
