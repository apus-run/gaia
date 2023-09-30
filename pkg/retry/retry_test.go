package retry

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRetry(t *testing.T) {

	r := Retry{
		Times:    3,
		Interval: 100 * time.Millisecond,
	}

	cnt := 0
	f := func() error {
		if cnt < r.Times {
			cnt++
			return fmt.Errorf("error, cnt=%d", cnt)
		}
		return nil
	}

	err := DoRetry("test", "dummy", "tag", r, f)
	assert.NotNil(t, err)
	assert.Equal(t, r.Times, cnt)

	cnt = 1
	err = DoRetry("test", "dummy", "tag", r, f)
	assert.Nil(t, err)
	assert.Equal(t, r.Times, cnt)

	f = func() error {
		cnt++
		return &ErrNoRetry{"No Retry Error"}
	}
	cnt = 0
	err = DoRetry("test", "dummy", "tag", r, f)
	assert.NotNil(t, err)
	assert.Equal(t, 1, cnt)
	assert.Equal(t, err.Error(), "No Retry Error")
}
