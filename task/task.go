package task

import (
	"fmt"
	"log"
	"time"
)

//TODO: export the type: type Task func() (string, error)

// ErrTimedOut is returned when an operation times out
type ErrTimedOut struct {
	// Reason is the reason for the timeout
	Reason string
}

func (e *ErrTimedOut) Error() string {
	return fmt.Sprintf("timed out performing task. Error was: %s", e.Reason)
}

// DoRetryWithTimeout performs given task with given timeout and timeBeforeRetry
func DoRetryWithTimeout(t func() (interface{}, bool, error), timeout, timeBeforeRetry time.Duration) (interface{}, error) {
	done := make(chan bool, 1)
	quit := make(chan bool, 1)
	var out interface{}
	var err error
	var retry bool

	go func() {
		count := 0
		for {
			select {
			case q := <-quit:
				if q {
					return
				}

			default:
				out, retry, err = t()
				if err == nil || !retry {
					done <- true
					return
				}

				log.Printf("%v Next retry in: %v", err, timeBeforeRetry)
				time.Sleep(timeBeforeRetry)
			}

			count++
		}
	}()

	select {
	case <-done:
		return out, err
	case <-time.After(timeout):
		quit <- true
		return out, &ErrTimedOut{
			Reason: err.Error(),
		}
	}
}
