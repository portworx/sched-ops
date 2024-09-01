package configmap

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	coreops "github.com/portworx/sched-ops/k8s/core"
	"github.com/stretchr/testify/require"
	mpatch "github.com/undefinedlabs/go-mpatch"
)

const (
	testLockTimeout         = 5 * time.Second
	testLockAttempts        = 3
	testLockRefreshDuration = 1 * time.Second
	testLockTTL             = 3 * time.Second
)

func TestMultilock(t *testing.T) {
	setUpConfigMapTestCluster(t)
	cm, err := New("px-configmaps-multil-lock-test", nil, testLockTimeout, testLockAttempts, testLockRefreshDuration, testLockTTL)
	require.NoError(t, err, "Unexpected error on New")

	fmt.Println("testMultilock")

	id1 := "multi-lock-id1"
	id2 := "multi-lock-id2"
	key1 := "multi-lock-key1"
	key2 := "multi-lock-key2"

	locked, _, err := cm.IsKeyLocked(key1)
	require.NoError(t, err)
	require.False(t, locked)

	fmt.Println("\tlocking key 1")
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")
	fmt.Println("\tlocked ID 1 key 1")

	locked, owner, err := cm.IsKeyLocked(key1)
	require.NoError(t, err)
	require.True(t, locked)
	require.Equal(t, id1, owner)

	fmt.Println("\tlocking key 2")
	err = cm.LockWithKey(id2, key2)
	require.NoError(t, err, "Unexpected error in LockWithKey(id2,key2)")
	fmt.Println("\tlocked ID 2 key 2")

	locked, owner, err = cm.IsKeyLocked(key2)
	require.NoError(t, err)
	require.True(t, locked)
	require.Equal(t, id2, owner)

	fmt.Println("\tunlocking key 1")
	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(id1,key1)")
	fmt.Println("\tunlocked ID 1 key 1")

	locked, _, err = cm.IsKeyLocked(key1)
	require.NoError(t, err)
	require.False(t, locked)

	fmt.Println("\tunlocking key 2")
	err = cm.UnlockWithKey(key2)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(id1,key1)")
	fmt.Println("\tunlocked ID 2 key 2")

	locked, _, err = cm.IsKeyLocked(key2)
	require.NoError(t, err)
	require.False(t, locked)

	fmt.Println("\tdouble lock with same id")
	fmt.Println("\tlocking ID 1 key 1")
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")
	fmt.Println("\tlocked ID 1 key 1")

	go func() {
		time.Sleep(1 * time.Second)
		// Second lock below will block until this unlock (because of same ID)
		fmt.Println("\tunlocking ID 1 key 1")
		err := cm.UnlockWithKey(key1)
		require.NoError(t, err, "Unexpected error from UnlockWithKey(key1)")
		fmt.Println("\tunlocked ID 1 key 1")
	}()
	fmt.Println("\tlocking ID 1 key 2")
	err = cm.LockWithKey(id1, key2)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key2)")
	fmt.Println("\tlocked ID 1 key 2")

	fmt.Println("\tunlocking key 2")
	err = cm.UnlockWithKey(key2)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(key2)")
	fmt.Println("\tunlocked key 2")

	fmt.Println("\tdouble lock with same key")
	fmt.Println("\tlocking ID 1 key 1")
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")
	fmt.Println("\tlocked ID 1 key 1")

	go func() {
		time.Sleep(1 * time.Second)
		// Second lock below will block until this unlock (because of same key)
		fmt.Println("\tunlocking ID 1 key 1")
		err := cm.UnlockWithKey(key1)
		require.NoError(t, err, "Unexpected error from UnlockWithKey(key1)")
		fmt.Println("\tunlocked ID 1 key 1")
	}()
	fmt.Println("\tlocking ID 2 key 1")
	err = cm.LockWithKey(id2, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id2,key1)")
	fmt.Println("\tlocked ID 2 key 1")

	fmt.Println("\tunlocking key 1")
	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(key1)")
	fmt.Println("\tunlocked ID 2 key 1")

	// all keys should be unlocked now
	locked, _, err = cm.IsKeyLocked(key1)
	require.NoError(t, err)
	require.False(t, locked)

	locked, _, err = cm.IsKeyLocked(key2)
	require.NoError(t, err)
	require.False(t, locked)

	err = cm.Delete()
	require.NoError(t, err, "Unexpected error on Delete")
}

func TestLockHoldTimeout(t *testing.T) {
	setUpConfigMapTestCluster(t)
	cm, err := New("px-configmaps-lock-hold-v2-test", nil, testLockTimeout, testLockAttempts, testLockRefreshDuration, testLockTTL)
	require.NoError(t, err, "Unexpected error on New")

	fmt.Println("TestLockHoldTimeout")

	id1 := "hold-timeout-id1"
	key1 := "hold-timeout-key1"
	id2 := "hold-timeout-id2"

	var lockTimedout bool
	fatalLockCb := func(format string, args ...interface{}) {
		fmt.Println("\tLock timeout called.")
		lockTimedout = true
		// since our fatalCb does not panic, we need to unlock to stop the lock from refreshing
		go func() {
			fmt.Println("\tunlocking key 1 from fatalCb")
			cm.UnlockWithKey(key1)
			fmt.Println("\tunlocked key 1 from fatalCb")
		}()
	}
	SetFatalCb(fatalLockCb)

	// Check lock expiration
	fmt.Println("\tlocking ID 1 key 1")
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in lock")
	fmt.Println("\tlocked ID 1 key 1")

	// Wait for lock hold timeout which will stop the refresh
	time.Sleep(testLockTimeout + time.Second)
	require.True(t, lockTimedout, "Lock hold timeout not triggered")

	// Locking again should not throw error
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in lock")

	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected error in unlock")

	err = cm.LockWithKey(id2, key1)
	require.NoError(t, err, "Unexpected error in lock")
	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected error in unlock")

	err = cm.Delete()
	require.NoError(t, err, "Unexpected error on Delete")
}

func TestCMLockRefreshV2(t *testing.T) {
	setUpConfigMapTestCluster(t)

	expiry := 24 * time.Second
	cmIntf, err := New("px-configmaps-lock-refresh-v2-test", nil, 5*time.Minute, 1000, expiry/3, expiry)
	require.NoError(t, err, "Unexpected error on New")

	cm := cmIntf.(*configMap)

	id1 := "lock-refresh-id1"
	key1 := "lock-refresh-key1"

	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")

	err = cm.PatchKeyLocked(false, id1, key1, "val1")
	require.NoError(t, err, "Unexpected error in Patch")

	time.Sleep(time.Duration(rand.Intn(int(expiry - time.Second))))

	err = cm.PatchKeyLocked(false, id1, key1, "val2")
	require.NoError(t, err, "Unexpected error in Patch")

	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(key1)")

	resultMap, err := cm.Get()
	require.NoError(t, err, "Unexpected error in Get")
	t.Log(resultMap)
	require.Contains(t, resultMap, key1)
	require.Equal(t, "val2", resultMap[key1])

	var reentrantCheck atomic.Bool
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			err = cm.LockWithKey(id1, key1)
			require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")

			require.True(t, reentrantCheck.CompareAndSwap(false, true), "Reentrant lock detected")

			val := fmt.Sprintf("val%d", i)
			err = cm.PatchKeyLocked(false, id1, key1, val)
			require.NoError(t, err, "Unexpected error in Patch")

			// give some time for the refreshLock goroutine to start
			time.Sleep(time.Duration(rand.Intn(int(expiry - time.Second))))

			resultMap, err := cm.Get()
			require.NoError(t, err, "Unexpected error in Get")
			t.Log(resultMap)
			require.Contains(t, resultMap, key1)
			require.Equal(t, val, resultMap[key1])

			require.True(t, reentrantCheck.CompareAndSwap(true, false), "Reentrant lock detected")

			err = cm.UnlockWithKey(key1)
			require.NoError(t, err, "Unexpected error in UnlockWithKey(key1)")
		}(i)
	}
	wg.Wait()

	err = cm.Delete()
	require.NoError(t, err, "Unexpected error on Delete")
}

func TestCMLockLostV2(t *testing.T) {
	setUpConfigMapTestCluster(t)
	cmIntf, err := New("px-configmaps-lock-lost-v2-test", nil, testLockTimeout, testLockAttempts, testLockRefreshDuration, testLockTTL)
	require.NoError(t, err, "Unexpected error on New")

	cm := cmIntf.(*configMap)

	id1 := "lock-lost-id1"
	key1 := "lock-lost-key1"

	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")

	err = cm.PatchKeyLocked(false, id1, key1, "val1")
	require.NoError(t, err, "Unexpected error in Patch")

	// case: lock lost with NO new owner
	setV2LockOwnerForTesting(t, cm, key1, "", time.Time{})
	err = cm.PatchKeyLocked(false, id1, key1, "val2")
	require.Error(t, err, "Expected error in Patch")
	require.ErrorContains(t, err, "lock check failed")

	// case: re-take the lock and update
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")

	err = cm.PatchKeyLocked(false, id1, key1, "val1")
	require.NoError(t, err, "Unexpected error in Patch")

	// case: lock lost to a new owner
	setV2LockOwnerForTesting(t, cm, key1, "new-owner", time.Now().Add(5*time.Minute))
	err = cm.PatchKeyLocked(false, id1, key1, "val2")
	require.Error(t, err, "Expected error in Patch")
	require.ErrorContains(t, err, "lock check failed")

	// case: re-taking the lock should fail with "configmap is locked" error
	err = cm.LockWithKey(id1, key1)
	require.Error(t, err, "Expected error in Patch")
	require.ErrorContains(t, err, "ConfigMap is locked")

	// case: new owner releases the lock; then we should be able to take the lock
	setV2LockOwnerForTesting(t, cm, key1, "", time.Time{})
	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")

	err = cm.PatchKeyLocked(false, id1, key1, "val3")
	require.NoError(t, err, "Unexpected error in Patch")

	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(key1)")

	resultMap, err := cm.Get()
	require.NoError(t, err, "Unexpected error in Get")
	require.Contains(t, resultMap, key1)
	require.Equal(t, "val3", resultMap[key1])

	err = cm.Delete()
	require.NoError(t, err, "Unexpected error on Delete")
}

func TestDeleteKeyLockedV2(t *testing.T) {
	setUpConfigMapTestCluster(t)
	cmIntf, err := New("px-configmaps-delete-key-v2-test", nil, testLockTimeout, testLockAttempts, testLockRefreshDuration, testLockTTL)
	require.NoError(t, err, "Unexpected error on New")

	cm := cmIntf.(*configMap)

	id1 := "delete-lock-id1"
	key1 := "delete-lock-key1"

	err = cm.LockWithKey(id1, key1)
	require.NoError(t, err, "Unexpected error in LockWithKey(id1,key1)")

	err = cm.PatchKeyLocked(false, id1, key1, "val1")
	require.NoError(t, err, "Unexpected error in Patch")

	resultMap, err := cm.Get()
	require.NoError(t, err, "Unexpected error in Get")
	require.Contains(t, resultMap, key1)
	require.Equal(t, "val1", resultMap[key1])

	err = cm.DeleteKeyLocked(false, id1, key1)
	require.NoError(t, err, "Unexpected error in DeleteKeyLocked")

	err = cm.UnlockWithKey(key1)
	require.NoError(t, err, "Unexpected error in UnlockWithKey(key1)")

	resultMap, err = cm.Get()
	require.NoError(t, err, "Unexpected error in Get")
	require.NotContains(t, resultMap, key1)

	err = cm.Delete()
	require.NoError(t, err, "Unexpected error on Delete")
}

func setV2LockOwnerForTesting(t *testing.T, cm *configMap, key, owner string, expiration time.Time) {
	require.Eventually(t, func() bool {
		rawCM, err := coreops.Instance().GetConfigMap(cm.name, k8sSystemNamespace)
		if err != nil {
			return false
		}
		lockOwners, lockExpirations, err := cm.parseLocks(rawCM)
		if err != nil {
			return false
		}
		if owner == "" {
			delete(lockOwners, key)
			delete(lockExpirations, key)
		} else {
			lockOwners[key] = owner
			lockExpirations[key] = expiration
		}

		err = cm.generateConfigMapData(rawCM, lockOwners, lockExpirations)
		if err != nil {
			return false
		}
		_, err = coreops.Instance().UpdateConfigMap(rawCM)
		return err == nil
	}, 5*time.Second, 100*time.Millisecond)
}

func TestCMLockTTL(t *testing.T) {
	setUpConfigMapTestCluster(t)
	cmIntf, err := New("px-configmaps-lock-ttl-v2-test", nil, testLockTimeout, testLockAttempts, testLockRefreshDuration, testLockTTL)
	require.NoError(t, err, "Unexpected error on New")

	cm := cmIntf.(*configMap)

	currTime := time.Now()
	myTime := func() time.Time {
		return currTime
	}
	timePatch, err := mpatch.PatchMethod(time.Now, myTime)
	require.NoError(t, err, "Failed to patch time.Now()")
	defer timePatch.Unpatch()

	owner := "lock-ttl-id1"
	key := "lock-ttl-key1"
	lockOwners := map[string]string{}
	lockExpirations := map[string]time.Time{}

	gotOwner, err := cm.checkAndTakeLock(owner, key, false, lockOwners, lockExpirations)

	require.NoError(t, err, "Unexpected error in checkAndTakeLock")
	require.Equal(t, owner, gotOwner, "Unexpected owner in checkAndTakeLock")
	require.Equal(t, owner, lockOwners[key])
	require.Equal(t, currTime.Add(testLockTTL), lockExpirations[key])
}
