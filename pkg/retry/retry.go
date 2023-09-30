package retry

import (
	"fmt"
	"time"

	"github.com/apus-run/sea-kit/log"
)

// Retry is the settings of retry
type Retry struct {
	Times    int           `yaml:"times" json:"times,omitempty" jsonschema:"title=Retry Times,description=how many times need to retry,minimum=1"`
	Interval time.Duration `yaml:"interval" json:"interval,omitempty" jsonschema:"type=string,format=duration,title=Retry Interval,description=the interval between each retry"`
}

// ErrNoRetry is the error need not retry
type ErrNoRetry struct {
	Message string
}

func (e *ErrNoRetry) Error() string {
	return e.Message
}

// DoRetry is a help function to retry the function if it returns error
func DoRetry(kind, name, tag string, r Retry, fn func() error) error {
	var err error
	for i := 0; i < r.Times; i++ {
		err = fn()
		_, ok := err.(*ErrNoRetry)
		if err == nil || ok {
			return err
		}
		log.Warnf("[%s / %s / %s] Retried to send %d/%d - %v", kind, name, tag, i+1, r.Times, err)

		// last time no need to sleep
		if i < r.Times-1 {
			time.Sleep(r.Interval)
		}
	}
	return fmt.Errorf("[%s / %s / %s] failed after %d retries - %v", kind, name, tag, r.Times, err)
}
