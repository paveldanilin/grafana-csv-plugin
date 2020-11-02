package macro

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestInterpolate(t *testing.T) {
	Register("sum", func(args []string, scope *Scope) (string, error) {
		a, _ := strconv.ParseInt(args[0], 10, 64)
		b, _ := strconv.ParseInt(args[1], 10, 64)
		return fmt.Sprintf("%d", a + b), nil
	})

	interpolated, err := Interpolate("SELECT $__sum(1,2) FROM dual", nil)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "SELECT 3 FROM dual", interpolated.Text())
}

func TestInterpolate_Meta(t *testing.T) {
	interpolated, err := Interpolate("SELECT * FROM products WHERE product_name = 'TestProduct' $autoTime(1m)", nil)

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "SELECT * FROM products WHERE product_name = 'TestProduct'", interpolated.Text())
	assert.NotEmpty(t, interpolated.MetaList())
	assert.Equal(t, "$autoTime", interpolated.MetaList()[0].Name)
	assert.NotEmpty(t, interpolated.MetaList()[0].Options)
	assert.Equal(t, "1m", interpolated.MetaList()[0].Options[0])
}
