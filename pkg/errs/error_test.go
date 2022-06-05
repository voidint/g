package errs

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLUnreachableError(t *testing.T) {
	t.Run("URL不可达错误", func(t *testing.T) {
		url := "https://github.com/voidint"
		core := errors.New("hello error")

		err := NewURLUnreachableError(url, core)
		assert.NotNil(t, err)

		e, ok := err.(*URLUnreachableError)
		assert.True(t, ok)
		assert.NotNil(t, e)
		assert.Equal(t, url, e.URL())
		assert.Equal(t, core, e.Err())
		assert.Equal(t, fmt.Sprintf("URL %q is unreachable ==> %s", url, core.Error()), e.Error())
	})
}
