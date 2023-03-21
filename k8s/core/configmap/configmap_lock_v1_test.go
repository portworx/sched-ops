package configmap

import (
	"fmt"
	"testing"
	"time"

	coreops "github.com/portworx/sched-ops/k8s/core"
	"github.com/stretchr/testify/require"
	fakek8sclient "k8s.io/client-go/kubernetes/fake"
)

func TestLock(t *testing.T) {
	fakeClient := fakek8sclient.NewSimpleClientset()
	coreops.SetInstance(coreops.New(fakeClient))
	cm, err := New("px-configmaps-test", nil, lockTimeout, 5, 0, 0)
	require.NoError(t, err, "Unexpected error on New")
	fmt.Println("testLock")

	id := "locktest"
	err = cm.Lock(id)
	require.NoError(t, err, "Unexpected error in lock")

	fmt.Println("\tunlock")
	err = cm.Unlock()
	require.NoError(t, err, "Unexpected error from Unlock")

	fmt.Println("\trelock")
	err = cm.Lock(id)
	require.NoError(t, err, "Failed to lock after unlock")

	fmt.Println("\treunlock")
	err = cm.Unlock()
	require.NoError(t, err, "Unexpected error from Unlock")

	fmt.Println("\trepeat lock once")
	err = cm.Lock(id)
	require.NoError(t, err, "Failed to lock unlock")

	done := 0
	go func() {
		time.Sleep(time.Second * 3)
		done = 1
		err := cm.Unlock()
		fmt.Println("\trepeat lock unlock once")
		require.NoError(t, err, "Unexpected error from Unlock")
	}()
	fmt.Println("\trepeat lock lock twice")
	err = cm.Lock(id)
	require.NoError(t, err, "Failed to lock")
	require.Equal(t, 1, done, "Locked before unlock")
	fmt.Println("\trepeat lock unlock twice")
	err = cm.Unlock()
	require.NoError(t, err, "Unexpected error from Unlock")

	for done == 0 {
		time.Sleep(time.Second)
	}

	id = "doubleLock"
	err = cm.Lock(id)
	require.NoError(t, err, "Unexpected error in lock")
	go func() {
		time.Sleep(3 * time.Second)
		err := cm.Unlock()
		require.NoError(t, err, "Unexpected error from Unlock")
	}()
	err = cm.Lock(id)
	require.NoError(t, err, "Double lock")
	err = cm.Unlock()
	require.NoError(t, err, "Unexpected error from Unlock")

	err = cm.Lock("id1")
	require.NoError(t, err, "Unexpected error in lock")
	go func() {
		time.Sleep(1 * time.Second)
		err := cm.Unlock()
		require.NoError(t, err, "Unexpected error from Unlock")
	}()
	err = cm.Lock("id2")
	require.NoError(t, err, "diff lock")
	err = cm.Unlock()
	require.NoError(t, err, "Unexpected error from Unlock")

	fmt.Println("\tlockExpiration")

	var lockTimedout bool
	fatalLockCb := func(format string, args ...interface{}) {
		fmt.Println("\tLock timeout called.")
		lockTimedout = true
		err := cm.Unlock()
		require.NoError(t, err, "Unexpected error from Unlock")
	}
	SetFatalCb(fatalLockCb)

	// Check lock expiration
	err = cm.Lock("id2")
	require.NoError(t, err, "Unexpected error in lock")
	time.Sleep(20 * time.Second)
	require.True(t, lockTimedout, "Lock hold timeout not triggered")

	err = cm.Lock("id3")
	require.NoError(t, err, "Lock should have expired")
	err = cm.Unlock()
	require.NoError(t, err, "Unexpected no error in unlock")
	err = cm.Delete()
	require.NoError(t, err, "Unexpected error on Delete")
}

func TestLockWithHoldTimeout(t *testing.T) {
	defaultHoldTimeout := 3 * time.Second
	customHoldTimeout := defaultHoldTimeout + v1DefaultK8sLockRefreshDuration + 10*time.Second
	fakeClient := fakek8sclient.NewSimpleClientset()
	coreops.SetInstance(coreops.New(fakeClient))
	cm, err := New("px-configmaps-test", nil, defaultHoldTimeout, 5, 0, 0)
	require.NoError(t, err, "Unexpected error on New")
	fmt.Println("TestLockWithHoldTimeout")

	var lockTimedout bool
	fatalLockCb := func(format string, args ...interface{}) {
		fmt.Println("\tLock timeout called.")
		lockTimedout = true
		err := cm.Unlock()
		require.NoError(t, err, "Unexpected error from Unlock")
	}
	SetFatalCb(fatalLockCb)

	// when custom lock hold timeout is more than the default lock hold timeout
	err = cm.LockWithHoldTimeout("id1", customHoldTimeout)
	require.NoError(t, err, "Unexpected error in lock")

	// lock hold timeout should not trigger after the default lock hold timeout period (plus refresh interval)
	time.Sleep(customHoldTimeout - 8*time.Second)
	require.False(t, lockTimedout, "Lock hold timeout should not have triggered")

	err = cm.Unlock()
	require.NoError(t, err, "Unexpected no error in unlock")

	err = cm.Delete()
	require.NoError(t, err, "Unexpected error on Delete")
}
