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
	cm, err := New("px-configmaps-test", nil, testLockTimeout, 5, 0, 0)
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
	}
	SetFatalCb(fatalLockCb)

	err = cm.Lock("id2")
	require.NoError(t, err, "Unexpected error in lock")
	time.Sleep(20 * time.Second)
	require.True(t, lockTimedout, "Lock hold timeout not triggered")

	err = cm.Unlock()
	require.NoError(t, err, "Unexpected no error in unlock")

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
	}
	SetFatalCb(fatalLockCb)

	// when custom lock hold timeout is more than the default lock hold timeout
	err = cm.LockWithParams("id1", customHoldTimeout, 0)
	require.NoError(t, err, "Unexpected error in lock")
	time.Sleep(20 * time.Second)

	err = cm.Unlock()
	require.NoError(t, err, "Unexpected no error in unlock")

	// lock hold timeout should not trigger after the default lock hold timeout period (plus refresh interval)
	time.Sleep(customHoldTimeout - 8*time.Second)
	require.False(t, lockTimedout, "Lock hold timeout should not have triggered")

	err = cm.Unlock()
	require.NoError(t, err, "Unexpected no error in unlock")

	err = cm.Delete()
	require.NoError(t, err, "Unexpected error on Delete")
}

func TestPatchKeyLockedV1(t *testing.T) {
	fakeClient := fakek8sclient.NewSimpleClientset()
	coreops.SetInstance(coreops.New(fakeClient))

	configData := map[string]string{
		"key1": "val1",
	}

	cm, err := New("px-configmaps-test", configData, testLockTimeout, 5, 0, 0)
	require.NoError(t, err, "Unexpected error in creating configmap")

	// case: empty lock owner while CM is not locked
	err = cm.PatchKeyLocked(true, "", "key2", "val2")
	require.Error(t, err, "Expected error in Patch")

	// case: non-empty lock owner while CM is not locked
	err = cm.PatchKeyLocked(true, "no-such-owner", "key2", "val2")
	require.Error(t, err, "Expected error in Patch")

	err = cm.Lock("lock-owner")
	require.NoError(t, err, "Unexpected error in Lock")

	// case: empty lock owner
	err = cm.PatchKeyLocked(true, "", "key2", "val2")
	require.Error(t, err, "Expected error in Patch")

	// case: wrong lock owner
	err = cm.PatchKeyLocked(true, "no-such-owner", "key2", "val2")
	require.Error(t, err, "Expected error in Patch")

	// case: correct lock owner
	err = cm.PatchKeyLocked(true, "lock-owner", "key2", "val2")
	require.NoError(t, err, "Unexpected error in Patch")

	err = cm.Unlock()
	require.NoError(t, err, "Unexpected error in Unlock")

	resultMap, err := cm.Get()
	require.NoError(t, err, "Unexpected error in Get")
	require.Contains(t, resultMap, "key1")
	require.Contains(t, resultMap, "key2")
	require.Equal(t, "val1", resultMap["key1"])
	require.Equal(t, "val2", resultMap["key2"])
	require.Contains(t, resultMap, pxGenerationKey)
	require.Equal(t, "1", resultMap[pxGenerationKey])
	fmt.Println(resultMap)

	// case: check generation increments after 2 more updates
	err = cm.Lock("lock-owner")
	require.NoError(t, err, "Unexpected error in Lock")

	err = cm.PatchKeyLocked(true, "lock-owner", "key1", "val2")
	require.NoError(t, err, "Unexpected error in Patch")

	err = cm.PatchKeyLocked(true, "lock-owner", "key2", "val2")
	require.NoError(t, err, "Unexpected error in Patch")

	err = cm.Unlock()
	require.NoError(t, err, "Unexpected error in Unlock")

	resultMap, err = cm.Get()
	require.NoError(t, err, "Unexpected error in Get")
	require.Contains(t, resultMap, pxGenerationKey)
	require.Equal(t, "3", resultMap[pxGenerationKey])
}

func TestDeleteKeyLockedV1(t *testing.T) {
	fakeClient := fakek8sclient.NewSimpleClientset()
	coreops.SetInstance(coreops.New(fakeClient))

	configData := map[string]string{
		"key1": "val1",
	}

	cm, err := New("px-configmaps-test", configData, testLockTimeout, 5, 0, 0)
	require.NoError(t, err, "Unexpected error in creating configmap")

	err = cm.Lock("lock-owner")
	require.NoError(t, err, "Unexpected error in Lock")

	err = cm.PatchKeyLocked(true, "lock-owner", "key2", "val2")
	require.NoError(t, err, "Unexpected error in Patch")

	resultMap, err := cm.Get()
	require.NoError(t, err, "Unexpected error in Get")
	require.Contains(t, resultMap, "key1")
	require.Contains(t, resultMap, "key2")

	err = cm.DeleteKeyLocked(true, "lock-owner", "key1")
	require.NoError(t, err, "Unexpected error in Patch")

	err = cm.Unlock()
	require.NoError(t, err, "Unexpected error in Unlock")

	resultMap, err = cm.Get()
	require.NoError(t, err, "Unexpected error in Get")
	require.Contains(t, resultMap, "key2")
	require.NotContains(t, resultMap, "key1")
	require.Equal(t, "2", resultMap[pxGenerationKey])
}

func TestCMLockLostV1(t *testing.T) {
	fakeClient := fakek8sclient.NewSimpleClientset()
	coreops.SetInstance(coreops.New(fakeClient))

	configData := map[string]string{
		"key1": "val1",
	}

	cm, err := New("px-configmaps-test", configData, 0, 0, 0, 0)
	require.NoError(t, err, "Unexpected error in creating configmap")

	err = cm.Lock("lock-owner")
	require.NoError(t, err, "Unexpected error in Lock")

	err = cm.PatchKeyLocked(true, "lock-owner", "key1", "val2")
	require.NoError(t, err, "Unexpected error in Patch")

	// case: lock lost with NO new owner
	setV1LockOwnerForTesting(t, "px-configmaps-test", "", time.Time{})
	err = cm.PatchKeyLocked(true, "lock-owner", "key1", "val3")
	require.Error(t, err, "Expected error in Patch")
	require.ErrorContains(t, err, "lock check failed")

	// case: re-take the lock and update
	err = cm.Lock("lock-owner")
	require.NoError(t, err, "Unexpected error in Lock")

	err = cm.PatchKeyLocked(true, "lock-owner", "key1", "val2")
	require.NoError(t, err, "Unexpected error in Patch")

	// case: lock lost to a new owner
	setV1LockOwnerForTesting(t, "px-configmaps-test", "new-owner", time.Now().Add(5*time.Minute))
	err = cm.PatchKeyLocked(true, "lock-owner", "key2", "val2")
	require.Error(t, err, "Expected error in Patch")
	require.ErrorContains(t, err, "lock check failed")

	// case: re-taking the lock should fail with "configmap is locked" error
	err = cm.Lock("lock-owner")
	require.Error(t, err, "Expected error in Lock")
	require.ErrorContains(t, err, "ConfigMap is locked")

	// case: new owner releases the lock; then we should be able to take the lock
	setV1LockOwnerForTesting(t, "px-configmaps-test", "", time.Time{})
	err = cm.Lock("lock-owner")
	require.NoError(t, err, "Unexpected error in Lock")

	err = cm.PatchKeyLocked(true, "lock-owner", "key1", "val2")
	require.NoError(t, err, "Unexpected error in Patch")

	err = cm.Unlock()
	require.NoError(t, err, "Unexpected error in Unlock")
}

func setV1LockOwnerForTesting(t *testing.T, cmName, owner string, expiration time.Time) {
	require.Eventually(t, func() bool {
		rawCM, err := coreops.Instance().GetConfigMap(cmName, k8sSystemNamespace)
		if err != nil {
			return false
		}
		if owner == "" {
			delete(rawCM.Data, pxOwnerKey)
			delete(rawCM.Data, pxExpirationKey)
		} else {
			rawCM.Data[pxOwnerKey] = owner
			rawCM.Data[pxExpirationKey] = expiration.Format(time.UnixDate)
		}
		_, err = coreops.Instance().UpdateConfigMap(rawCM)
		return err == nil
	}, 5*time.Second, 100*time.Millisecond)
}
