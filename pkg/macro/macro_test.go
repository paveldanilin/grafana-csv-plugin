package macro

import (
	"fmt"
	"strconv"
	"testing"
)

func TestInterpolate(t *testing.T) {
	Register("sum", func(args []string, scope *Scope) (string, error) {
		a, _ := strconv.ParseInt(args[0], 10, 64)
		b, _ := strconv.ParseInt(args[1], 10, 64)
		return fmt.Sprintf("%d", a + b), nil
	})

	text, err := Interpolate("SELECT $__sum(1,2) FROM dual", nil)

	if err != nil {
		t.Error(err)
	}

	println(text)
}
