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
