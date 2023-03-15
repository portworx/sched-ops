package configmap

import (
	"fmt"
	"testing"
	"time"

	coreops "github.com/portworx/sched-ops/k8s/core"
	"github.com/stretchr/testify/require"
	fakek8sclient "k8s.io/client-go/kubernetes/fake"
)

const (
	lockTimeout = 15 * time.Second
)

func TestMultilock(t *testing.T) {
	fakeClient := fakek8sclient.NewSimpleClientset()
	coreops.SetInstance(coreops.New(fakeClient))
	cm, err := New("px-configmaps-test", nil, lockTimeout, 3, 0, 0, "")
	require.NoError(t, err, "Unexpected error on New")

	fmt.Println("testMultilock")

	id1 := "id1"
	id2 := "id2"
	key1 := "key1"
	key2 := "key2"

	locked, _, err := cm.IsKeyLocked(key1)
	require.NoError(t, err)
	require.False(t, locked)

	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")

	locked, owner, err := cm.IsKeyLocked(key1)
	require.NoError(t, err)
	require.True(t, locked)
	require.Equal(t, id1, owner)

	fmt.Println("\tlocking key 2")
	err = cm.LockWithKey(id2, key2)
	require.NoError(t, err, "Unexpected error in LockWithKey(id2,key2)")

	locked, owner, err = cm.IsKeyLocked(key2)
	require.NoError(t, err)
	require.True(t, locked)
	require.Equal(t, id2, owner)

	fmt.Println("\tunlocking key 1")
	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(id1,key1)")

	locked, owner, err = cm.IsKeyLocked(key1)
	require.NoError(t, err)
	require.False(t, locked)

	fmt.Println("\tunlocking key 2")
	err = cm.UnlockWithKey(key2)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(id1,key1)")

	locked, owner, err = cm.IsKeyLocked(key2)
	require.NoError(t, err)
	require.False(t, locked)

	fmt.Println("\tdouble lock with same id")
	fmt.Println("\tlocking ID 1 key 1")
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")

	go func() {
		time.Sleep(1 * time.Second)
		// Second lock below will block until this unlock (because of same ID)
		fmt.Println("\tunlocking key 1")
		err := cm.UnlockWithKey(key1)
		require.NoError(t, err, "Unexpected error from UnlockWithKey(key1)")
	}()
	fmt.Println("\tlocking ID 1 key 2")
	err = cm.LockWithKey(id1, key2)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key2)")

	fmt.Println("\tunlocking key 2")
	err = cm.UnlockWithKey(key2)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(key2)")

	fmt.Println("\tdouble lock with same key")
	fmt.Println("\tlocking ID 1 key 1")
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")

	go func() {
		time.Sleep(1 * time.Second)
		// Second lock below will block until this unlock (because of same key)
		fmt.Println("\tunlocking key 1")
		err := cm.UnlockWithKey(key1)
		require.NoError(t, err, "Unexpected error from UnlockWithKey(key1)")
	}()
	fmt.Println("\tlocking ID 2 key 1")
	err = cm.LockWithKey(id2, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id2,key1)")

	fmt.Println("\tunlocking key 1")
	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(key1)")

	fmt.Println("\tlockExpiration")

	var lockTimedout bool
	fatalLockCb := func(format string, args ...interface{}) {
		fmt.Println("\tLock timeout called.")
		lockTimedout = true
	}
	SetFatalCb(fatalLockCb)

	// Check lock expiration
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in lock")

	// Locking again with same owner should not throw error
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in lock")
	fmt.Println("time : ", time.Now())
	time.Sleep((v2DefaultK8sLockRefreshDuration * 3) + (3 * time.Second))
	fmt.Println("time : ", time.Now())
	require.True(t, lockTimedout, "Lock hold timeout not triggered")

	// Locking again with expired lock should not throw error
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in lock")

	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected no error in unlock")

	err = cm.LockWithKey(id2, key1)
	require.NoError(t, err, "Lock should have expired")
	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected no error in unlock")
	err = cm.Delete()
	require.NoError(t, err, "Unexpected error on Delete")
}
