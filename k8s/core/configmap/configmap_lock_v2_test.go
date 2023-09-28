package configmap

import (
	"testing"
	"time"

	coreops "github.com/portworx/sched-ops/k8s/core"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	fakek8sclient "k8s.io/client-go/kubernetes/fake"
)

const (
	lockTimeout = 15 * time.Second
)

func TestMultilock(t *testing.T) {
	t.Skip("skipped until PWX-31627 is fixed")
	fakeClient := fakek8sclient.NewSimpleClientset()
	coreops.SetInstance(coreops.New(fakeClient))
	cm, err := New("px-configmaps-test", nil, lockTimeout, 3, 0, 0)
	require.NoError(t, err, "Unexpected error on New")

	logrus.Infof("testMultilock")

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

	logrus.Infof("\tlocking key 2")
	err = cm.LockWithKey(id2, key2)
	require.NoError(t, err, "Unexpected error in LockWithKey(id2,key2)")

	locked, owner, err = cm.IsKeyLocked(key2)
	require.NoError(t, err)
	require.True(t, locked)
	require.Equal(t, id2, owner)

	logrus.Infof("\tunlocking key 1")
	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(id1,key1)")

	locked, owner, err = cm.IsKeyLocked(key1)
	require.NoError(t, err)
	require.False(t, locked)

	logrus.Infof("\tunlocking key 2")
	err = cm.UnlockWithKey(key2)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(id1,key1)")

	locked, owner, err = cm.IsKeyLocked(key2)
	require.NoError(t, err)
	require.False(t, locked)

	logrus.Infof("\tdouble lock with same id")
	logrus.Infof("\tlocking ID 1 key 1")
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")

	go func() {
		time.Sleep(1 * time.Second)
		// Second lock below will block until this unlock (because of same ID)
		logrus.Infof("\tunlocking key 1")
		err := cm.UnlockWithKey(key1)
		require.NoError(t, err, "Unexpected error from UnlockWithKey(key1)")
	}()
	logrus.Infof("\tlocking ID 1 key 2")
	err = cm.LockWithKey(id1, key2)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key2)")

	logrus.Infof("\tunlocking key 2")
	err = cm.UnlockWithKey(key2)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(key2)")

	logrus.Infof("\tdouble lock with same key")
	logrus.Infof("\tlocking ID 1 key 1")
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")

	go func() {
		time.Sleep(1 * time.Second)
		// Second lock below will block until this unlock (because of same key)
		logrus.Infof("\tunlocking key 1")
		err := cm.UnlockWithKey(key1)
		require.NoError(t, err, "Unexpected error from UnlockWithKey(key1)")
	}()
	logrus.Infof("\tlocking ID 2 key 1")
	err = cm.LockWithKey(id2, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id2,key1)")

	logrus.Infof("\tunlocking key 1")
	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(key1)")

	logrus.Infof("\tlockExpiration")

	var lockTimedout bool
	fatalLockCb := func(format string, args ...interface{}) {
		logrus.Infof("\tLock timeout called.")
		lockTimedout = true
	}
	SetFatalCb(fatalLockCb)

	// Check lock expiration
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in lock")

	// Locking again with same owner should not throw error
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in lock")
	time.Sleep((v2DefaultK8sLockRefreshDuration * 3) + (3 * time.Second))
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
