package configmap

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"

	"github.com/libopenstorage/openstorage/pkg/dbg"
	"github.com/portworx/sched-ops/k8s/core"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
)

func (c *configMap) LockWithKey(owner, key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	fn := "LockWithKey"
	configMapLog(fn, c.name, owner, key, nil).Debugf("Taking the lock")

	count := uint(0)
	// try acquiring a lock on the ConfigMap
	newOwner, err := c.tryLock(owner, key, false)
	// if it fails, keep trying for the provided number of retries until it succeeds
	for maxCount := c.lockAttempts; err != nil && count < maxCount; count++ {
		time.Sleep(lockSleepDuration)
		newOwner, err = c.tryLock(owner, key, false)
		if count > 0 && count%15 == 0 && err != nil {
			configMapLog(fn, c.name, newOwner, key, err).Warnf("Locked for"+
				" %v seconds", float64(count)*lockSleepDuration.Seconds())
		}
	}
	if err != nil {
		// We failed to acquire the lock
		return err
	}
	if count >= 30 {
		configMapLog(fn, c.name, newOwner, key, err).Warnf("Spent %v iteration"+
			" locking.", count)
	}

	// We have acquired the lock on the configmap. If the previous owner was the same node, old refreshLock
	// goroutine may still be running since we don't write to the Done chan. Take care of the old lock
	// so that the old refreshLock goroutine does not interfere with the new lock.
	c.kLocksV2Mutex.Lock()
	oldLock := c.kLocksV2[key]
	c.kLocksV2Mutex.Unlock()
	if oldLock != nil {
		oldLock.Lock()
		if !oldLock.unlocked {
			configMapLog(fn, c.name, newOwner, key, err).Warn("Found old lock still locked. Unlocking...")
			close(oldLock.done)
			oldLock.unlocked = true
		}
		oldLock.Unlock()
	}

	// Create a new lock and store it
	lock := &k8sLock{done: make(chan struct{}), id: owner}
	c.kLocksV2Mutex.Lock()
	c.kLocksV2[key] = lock
	c.kLocksV2Mutex.Unlock()

	configMapLog(fn, c.name, owner, key, nil).Debugf("Starting lock refresh")
	go c.refreshLock(owner, key)
	return nil
}

func (c *configMap) UnlockWithKey(key string) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	fn := "UnlockWithKey"
	configMapLog(fn, c.name, "", key, nil).Debugf("Releasing the lock for %s", key)

	// Get the lock reference now so we don't have to keep locking and unlocking
	c.kLocksV2Mutex.Lock()
	lock, ok := c.kLocksV2[key]
	c.kLocksV2Mutex.Unlock()
	if !ok {
		return nil
	}

	lock.Lock()
	defer lock.Unlock()

	if lock.unlocked {
		// The lock is already unlocked. Chan cannot be closed again. Return immediately.
		return nil
	}
	lock.unlocked = true
	// Don't write to the chan since the refresh goroutine might have exited already if the lock was lost.
	// If we write to the chan, we may block indefinitely. Note that LockWithKey will also close the old chan if
	// old lock is not yet unlocked.
	close(lock.done)

	var (
		err error
		cm  *v1.ConfigMap
	)

	// Get the existing ConfigMap
	for retries := 0; retries < maxConflictRetries; retries++ {
		cm, err = core.Instance().GetConfigMap(
			c.name,
			k8sSystemNamespace,
		)
		if err != nil {
			// A ConfigMap should always be created.
			return err
		}

		lockOwners, lockExpirations, err := c.parseLocks(cm)
		if err != nil {
			return fmt.Errorf("failed to get locks from configmap: %v", err)
		}

		currentOwner := lockOwners[key]
		if currentOwner != lock.id {
			return nil
		}

		// We are holding the lock, let's remove it
		delete(lockOwners, key)
		delete(lockExpirations, key)

		err = c.generateConfigMapData(cm, lockOwners, lockExpirations)
		if err != nil {
			return err
		}

		if k8sConflict, err := c.updateConfigMap(cm); err != nil {
			configMapLog(fn, c.name, "", "", err).Errorf("Failed to update" +
				" config map during unlock")
			if k8sConflict {
				// try unlocking again
				continue
			}
			// else unknown error - return immediately
			return err
		}

		// Clean up the lock
		c.kLocksV2Mutex.Lock()
		delete(c.kLocksV2, key)
		c.kLocksV2Mutex.Unlock()
		return nil
	}

	return err
}

// If the lock hasn't expired but the owner doesn't have a refresh goroutine,
// we return unlocked for the owner but locked for others
func (c *configMap) IsKeyLocked(key, requester string) (bool, string, error) {
	// Get the existing ConfigMap
	cm, err := core.Instance().GetConfigMap(
		c.name,
		k8sSystemNamespace,
	)
	if err != nil {
		return false, "", err
	}

	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	lockIDs, lockExpirations, err := c.parseLocks(cm)
	if err != nil {
		return false, "", fmt.Errorf("failed to get locks from configmap: %v", err)
	}

	if owner, ok := lockIDs[key]; ok {
		// Existing key is unlocked if
		//   1. Lock has expired; or
		//   2. Lock owned by itself but refresh goroutine is not running
		expiration := lockExpirations[key]
		if time.Now().After(expiration) {
			return false, "", nil
		}
		if c.ifRequesterIsLockOwnerWithoutGoroutine(requester, owner, key) {
			return false, owner, nil
		}
		return true, owner, nil
	}

	// Nobody owns the lock
	return false, "", nil
}

func (c *configMap) tryLock(owner string, key string, refresh bool) (string, error) {
	// Get the existing ConfigMap
	cm, err := core.Instance().GetConfigMap(
		c.name,
		k8sSystemNamespace,
	)
	if err != nil {
		// A ConfigMap should always be created.
		return "", err
	}

	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	lockIDs, lockExpirations, err := c.parseLocks(cm)
	if err != nil {
		return "", fmt.Errorf("failed to get locks from configmap: %v", err)
	}

	finalOwner, err := c.checkAndTakeLock(owner, key, refresh, lockIDs, lockExpirations)
	if err != nil {
		return finalOwner, err
	}

	err = c.generateConfigMapData(cm, lockIDs, lockExpirations)
	if err != nil {
		return finalOwner, err
	}

	if _, err := c.updateConfigMap(cm); err != nil {
		return "", err
	}
	return owner, nil
}

// Returns the lock owner for the v2 lock regardless of whether the lock has expired.
func (c *configMap) getV2LockOwnerIncludeExpired(cm *v1.ConfigMap, key string) (string, error) {
	lockOwners, _, err := c.parseLocks(cm)
	if err != nil {
		return "", err
	}
	return lockOwners[key], nil
}

// parseLocks reads the lock data from the given ConfigMap and then converts it to:
// * a map of keys to lock owners
// * a map of keys to lock expiration times
func (c *configMap) parseLocks(cm *v1.ConfigMap) (map[string]string, map[string]time.Time, error) {
	// Check all the locks: will be an empty string if key is not present indicating no lock
	parsedLocks := []lockData{}
	if lock, ok := cm.Data[pxLockKey]; ok && len(lock) > 0 {
		err := json.Unmarshal([]byte(cm.Data[pxLockKey]), &parsedLocks)
		if err != nil {
			return nil, nil, err
		}
	}

	// Check all the locks first and store them, makes the looping a little easier
	lockOwners := map[string]string{}
	lockExpirations := map[string]time.Time{}

	for _, lock := range parsedLocks {
		lockOwners[lock.Key] = lock.Owner
		lockExpirations[lock.Key] = lock.Expiration
	}

	return lockOwners, lockExpirations, nil
}

// checkAndTakeLock checks if we can take the desired lock (refresh=false) or extend the expiration of the lock
// we have taken already (refresh=true). If either condition is true, it updates the in-memory state in
// lockOwners and lockExpirations. "refresh" argument indicates if this is the refreshLock goroutine refreshing
// the lock or an initial Lock call taking the lock.
func (c *configMap) checkAndTakeLock(
	owner, key string,
	refresh bool,
	lockOwners map[string]string,
	lockExpirations map[string]time.Time,
) (string, error) {
	fn := "checkAndTakeLock"
	currentOwner, ownerOK := lockOwners[key]
	_, expOK := lockExpirations[key]

	// Just check to make sure these are consistent and that we either have both or don't
	if ownerOK != expOK {
		return "", fmt.Errorf("inconsistent lock ID and expiration")
	}
	k8sTTL := c.lockK8sLockTTL

	if refresh {
		if currentOwner != owner {
			// We lost our lock probably because we could not refresh it in time. Note that even if the
			// lock is available now, it is not safe to just re-acquire the lock here in refresh.
			// We will fail any subsequent Patch/Delete calls being made under this lock. Caller must start over and
			// call Lock() again.
			configMapLog(fn, c.name, "", "", nil).Warnf(
				"Lost our lock on key %s in the configMap %s to a new owner %q", key, c.name, currentOwner)
			return currentOwner, ErrConfigMapLockLost
		}
		// We hold the lock; just refresh the expiry
		lockExpirations[key] = time.Now().Add(k8sTTL)
		return owner, nil
	}

	if currentOwner == "" {
		lockOwners[key] = owner
		lockExpirations[key] = time.Now().Add(k8sTTL)
		return owner, nil
	}

	if time.Now().Before(lockExpirations[key]) {
		if c.ifRequesterIsLockOwnerWithoutGoroutine(owner, currentOwner, key) {
			lockExpirations[key] = time.Now().Add(k8sTTL)
			return owner, nil
		}
		return lockOwners[key], ErrConfigMapLocked
	}

	configMapLog(fn, c.name, owner, key, nil).Infof(
		"Lock from owner '%s' is expired, now claiming for new owner '%s'", currentOwner, owner)

	lockOwners[key] = owner
	lockExpirations[key] = time.Now().Add(k8sTTL)
	return owner, nil
}

// generateConfigMapData converts the given lock data (lockOwners, lockExpirations) to JSON and
// stores it in the given ConfigMap.
func (c *configMap) generateConfigMapData(cm *v1.ConfigMap, lockOwners map[string]string, lockExpirations map[string]time.Time) error {
	var locks []lockData
	for key, lockOwner := range lockOwners {
		locks = append(locks, lockData{
			Owner:      lockOwner,
			Key:        key,
			Expiration: lockExpirations[key],
		})
	}

	cmData, err := json.Marshal(locks)
	if err != nil {
		return err
	}
	cm.Data[pxLockKey] = string(cmData)
	return nil
}

func (c *configMap) updateConfigMap(cm *v1.ConfigMap) (bool, error) {
	if _, err := core.Instance().UpdateConfigMap(cm); err != nil {
		return k8s_errors.IsConflict(err), err
	}
	return false, nil
}

// refreshLock is the goroutine running in the background after calling LockWithKey.
// It keeps the lock refreshed in k8s until we call Unlock. This is so that if the
// node dies, the lock can have a short timeout and expire quickly but we can still
// take longer-term locks.
func (c *configMap) refreshLock(id, key string) {
	fn := "refreshLock"
	refresh := time.NewTicker(c.lockRefreshDuration)
	defer refresh.Stop()
	var (
		currentRefresh time.Time
		prevRefresh    time.Time
		startTime      time.Time
	)

	// get a reference to the lock object so we don't have to hold open a
	// map reference - this makes it easier for concurrency purposes (can't
	// lock in a select condition)
	c.kLocksV2Mutex.Lock()
	lock := c.kLocksV2[key]
	c.kLocksV2Mutex.Unlock()

	if lock == nil {
		// could happen if the lock was unlocked before the goroutine started
		configMapLog(fn, c.name, "", key, nil).Warnf("Lock not found for key %s; refresh goroutine exiting", key)
		return
	}
	defer func() {
		lock.Lock()
		lock.refreshing = false
		lock.Unlock()
	}()
	lock.Lock()
	lock.refreshing = true
	lock.Unlock()

	startTime = time.Now()
	for {
		select {
		case <-refresh.C:
			lock.Lock()

			for !lock.unlocked {
				c.checkLockTimeout(c.defaultLockHoldTimeout, startTime, id)
				currentRefresh = time.Now()
				if _, err := c.tryLock(id, key, true); err != nil {
					if k8s_errors.IsConflict(err) {
						// try refreshing again
						continue
					}
					configMapLog(fn, c.name, "", key, err).Errorf(
						"Error refreshing lock. [ID %v] [Key %v] [Err: %v]"+
							" [Current Refresh: %v] [Previous Refresh: %v]",
						id, key, err, currentRefresh, prevRefresh,
					)
					if errors.Is(err, ErrConfigMapLockLost) {
						// there is no coming back from this
						lock.unlocked = true
						lock.Unlock()
						return
					}
				}
				thresh := c.lockRefreshDuration * 3 / 2
				if !prevRefresh.IsZero() && prevRefresh.Add(thresh).Before(currentRefresh) {
					configMapLog(fn, c.name, "", "", nil).Warnf(
						"V2 lock refresh is taking too long. [Owner %v] [Current Refresh: %v] [Previous Refresh: %v]",
						id, currentRefresh, prevRefresh)
				}
				prevRefresh = currentRefresh
				break
			}
			lock.Unlock()
		case <-lock.done:
			return
		}
	}

}

func (c *configMap) checkLockTimeout(holdTimeout time.Duration, startTime time.Time, id string) {
	if holdTimeout > 0 && time.Since(startTime) > holdTimeout {
		panicMsg := fmt.Sprintf("Lock hold timeout (%v) triggered for K8s configmap lock key %s", holdTimeout, id)
		if fatalCb != nil {
			fatalCb(panicMsg)
		} else {
			dbg.Panicf(panicMsg)
		}
	}
}

// We want to avoid locks being re-entrant AND want the owner to re-aquire the lock if there has been a restart.
// For that reason we check if the refresh routine is running.
// On a restart, the routine would have been cancelled.
func (c *configMap) ifRequesterIsLockOwnerWithoutGoroutine(requester, owner, key string) bool {
	c.kLocksV2Mutex.Lock()
	lock := c.kLocksV2[key]
	c.kLocksV2Mutex.Unlock()
	if lock == nil {
		return requester == owner
	}
	lock.Lock()
	defer lock.Unlock()
	return requester == owner && !lock.refreshing
}
