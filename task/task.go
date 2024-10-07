package task

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
)

//TODO: export the type: type Task func() (string, error)

// ErrTimedOut is returned when an operation times out
// Is this type used anywhere? If not we can get rid off it in favor context.DeadlineExceeded
type ErrTimedOut struct {
	// Reason is the reason for the timeout
	Reason string
}

// ExtraArgs defines a function type that takes a pointer to AdditionalInfo and modifies it.
// This is used for passing and applying configuration options.
type ExtraArgs func(*AdditionalInfo)

// AdditionalInfo holds configuration details that can be modified through functional options.
// Currently, it only holds a TestName, but more fields can be added as needed.
type AdditionalInfo struct {
	TestName string // TestName is a descriptor that can be used to identify or describe the test being performed.
}

// WithAdditionalInfo returns an ExtraArgs function that sets the TestName of an AdditionalInfo.
// This function is a functional option that allows callers to specify a test name for logging or identification purposes.
func WithAdditionalInfo(testName string) ExtraArgs {
	return func(cfg *AdditionalInfo) {
		cfg.TestName = testName // Sets the TestName field of the AdditionalInfo struct.
	}
}

func (e *ErrTimedOut) Error() string {
	errString := "timed out performing task."
	if len(e.Reason) > 0 {
		errString = fmt.Sprintf("%s, Error was: %s", errString, e.Reason)
	}

	return errString
}

// DoRetryWithTimeout performs given task with given timeout and timeBeforeRetry
// TODO(stgleb): In future I would like to add context as a first param to this function
// so calling code can cancel task.
func DoRetryWithTimeout(t func() (interface{}, bool, error), timeout, timeBeforeRetry time.Duration, opts ...ExtraArgs) (interface{}, error) {
	args := &AdditionalInfo{}
	for _, opt := range opts {
		opt(args)
	}
	if len(opts) > 0 && args.TestName != "" {
		log.Printf("In DoRetryWithTimeout method for test case: {%v}", args.TestName)
	}

	// Use context.Context as a standard go way of timeout and cancellation propagation amount goroutines.
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resultChan := make(chan interface{})
	errChan := make(chan error)
	errInRetires := make([]string, 0)

	go func() {
		for {
			select {
			case <-ctx.Done():

				if ctx.Err() != nil {
					errChan <- ctx.Err()
				}

				return
			default:
				out, retry, err := t()
				if err != nil {
					if retry {
						errInRetires = append(errInRetires, err.Error())
						if len(opts) > 0 && args.TestName != "" {
							log.Printf("DoRetryWithTimeout [%s] - Error: {%v}, Next try in [%v], timeout [%v]", args.TestName, err, timeBeforeRetry, timeout)
						} else {
							log.Printf("DoRetryWithTimeout - Error: {%v}, Next try in [%v], timeout [%v]", err, timeBeforeRetry, timeout)
						}
						time.Sleep(timeBeforeRetry)
					} else {
						errChan <- err
						return
					}
				} else {
					resultChan <- out
					return
				}
			}
		}
	}()

	select {
	case result := <-resultChan:
		return result, nil
	case err := <-errChan:
		if err == context.DeadlineExceeded {
			timeoutReason := fmt.Sprintf("DoRetryWithTimeout timed out. Errors generated in retries: {%s}", strings.Join(errInRetires, "}\n{"))
			if len(opts) > 0 && args.TestName != "" {
				timeoutReason = fmt.Sprintf("DoRetryWithTimeout [%s] timed out. Errors generated in retries: {%s}", args.TestName, strings.Join(errInRetires, "}\n{"))
			}
			return nil, &ErrTimedOut{
				Reason: timeoutReason,
			}
		}

		return nil, err
	}
}
