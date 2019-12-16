package task

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDoRetry(t *testing.T) {
	t1 := func() (interface{}, bool, error) {
		return "hello world", false, nil
	}

	output, err := DoRetryWithTimeout(t1, 1*time.Minute, 5*time.Second)
	require.NoError(t, err, "failed to run task")
	require.NotEmpty(t, output, "ask returned empty output")

	t2 := func() (interface{}, bool, error) {
		return "", true, fmt.Errorf("task is failing")
	}

	retryTill := time.Now().Add(10 * time.Second)
	output, err = DoRetryWithTimeout(t2, 10*time.Second, 2*time.Second)
	require.Error(t, err, "task was expected to fail")
	require.Empty(t, output, "task should have returned empty output")
	require.True(t, time.Now().After(retryTill) || time.Now().Equal(retryTill), "current time should be after expected timeout")

	// tighter retry loop
	t3 := func() (interface{}, bool, error) {
		return "", true, fmt.Errorf("task is failing")
	}

	retryTill = time.Now().Add(2 * time.Second)
	output, err = DoRetryWithTimeout(t3, 2*time.Second, 10*time.Millisecond)
	require.Error(t, err, "task was expected to fail")
	require.Empty(t, output, "task should have returned empty output")
	require.True(t, time.Now().After(retryTill) || time.Now().Equal(retryTill), "current time should be after expected timeout")

}

func TestDoRetryWithTimeoutSuccessAfter(t *testing.T) {
	counter := 0
	t4 := func() (interface{}, bool, error) {

		if counter > 3 {
			return "", false, nil
		}

		counter += 1
		return nil, true, fmt.Errorf("task is failing")
	}

	output, err := DoRetryWithTimeout(t4, 100*time.Millisecond, 10*time.Millisecond)
	require.NoError(t, err, "task must not fail")
	require.NotNil(t, output, "result must not  be nil")
}
