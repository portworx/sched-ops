package configmap

import (
	"errors"
	"time"

	"github.com/portworx/sched-ops/k8s/core"
	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
)

func (c *configMap) Lock(id string) error {
	return c.LockWithParams(id, c.defaultLockHoldTimeout, c.lockAttempts)
}

func (c *configMap) LockWithParams(id string, holdTimeout time.Duration, numAttempts uint) error {
	fn := "LockWithParams"

	// wait for any previous refreshLockV1 goroutine to exit
	c.kLockV1.wg.Wait()

	if numAttempts == 0 {
		// This is the same no. of times (300) we try while acquiring a kvdb lock
		numAttempts = c.lockAttempts
	}
	count := uint(0)
	// try acquiring a lock on the ConfigMap
	owner, err := c.tryLockV1(id, false)
	for maxCount := numAttempts; err != nil && count < maxCount; count++ {
		time.Sleep(lockSleepDuration)
		owner, err = c.tryLockV1(id, false)
		if count > 0 && count%15 == 0 && err != nil {
			configMapLog(fn, c.name, owner, "", err).Warnf("Locked for"+
				" %v seconds", float64(count)*lockSleepDuration.Seconds())
		}
	}
	if err != nil {
		// We failed to acquire the lock
		return err
	}
	if count >= 30 {
		configMapLog(fn, c.name, owner, "", err).Warnf("Spent %v iteration"+
			" locking.", count)
	}
	c.kLockV1.Lock()
	defer c.kLockV1.Unlock()
	c.lockHoldTimeoutV1 = holdTimeout
	c.kLockV1.id = id
	c.kLockV1.unlocked = false
	c.kLockV1.done = make(chan struct{})

	c.kLockV1.wg.Add(1)
	go func() {
		c.refreshLockV1(id)
		c.kLockV1.wg.Done()
	}()
	return nil
}

func (c *configMap) Unlock() error {
	fn := "Unlock"
	// Get the existing ConfigMap
	c.kLockV1.Lock()
	defer c.kLockV1.Unlock()

	if c.kLockV1.unlocked {
		// The lock is already unlocked
		return nil
	}
	c.kLockV1.unlocked = true
	// Don't write to the chan since the refresh goroutine might have exited already if the lock was lost.
	// If we write to the chan, we may block indefinitely.
	close(c.kLockV1.done)

	var (
		err error
		cm  *corev1.ConfigMap
	)
	for retries := 0; retries < maxConflictRetries; retries++ {
		cm, err = core.Instance().GetConfigMap(
			c.name,
			k8sSystemNamespace,
		)
		if err != nil {
			// A ConfigMap should always be created.
			return err
		}

		currentOwner := cm.Data[pxOwnerKey]
		if currentOwner != c.kLockV1.id {
			// We are currently not holding the lock
			return nil
		}
		delete(cm.Data, pxOwnerKey)
		delete(cm.Data, pxExpirationKey)
		if _, err = core.Instance().UpdateConfigMap(cm); err != nil {
			configMapLog(fn, c.name, "", "", err).Errorf("Failed to update" +
				" config map during unlock")
			if k8s_errors.IsConflict(err) {
				// try unlocking again
				continue
			} // else unknown error - return immediately
			return err
		}
		c.kLockV1.id = ""
		return nil
	}
	return err
}

func (c *configMap) tryLockV1(id string, refresh bool) (string, error) {
	fn := "tryLockV1"
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

	currentOwner := cm.Data[pxOwnerKey]
	if refresh {
		if currentOwner != id {
			// We lost our lock probably because we could not refresh it in time. Note that even if the
			// lock is available now, it is not safe to just re-acquire the lock here in refresh.
			// We will fail any subsequent Patch/Delete calls being made under this lock. Caller must start over and
			// call Lock() again.
			configMapLog(fn, c.name, "", "", nil).Warnf(
				"Lost our lock on the configMap %s to a new owner %q", c.name, currentOwner)
			return "", ErrConfigMapLockLost
		}
		// We already hold the lock just refresh our expiry (fall through)
	} else if currentOwner != "" {
		// Someone else might have a lock on the cm. Check expiration.
		expiration := cm.Data[pxExpirationKey]
		if expiration == "" {
			configMapLog(fn, c.name, "", "", nil).Warnf(
				"ConfigMap %v has an owner %s but no expiry; taking the lock away", c.name, currentOwner)
		} else {
			expiresAt, err := time.Parse(time.UnixDate, expiration)
			if err != nil {
				return currentOwner, err
			}
			if time.Now().Before(expiresAt) {
				// Lock is currently held by the owner
				// Retry after sometime
				return currentOwner, ErrConfigMapLocked
			} // else lock is expired. Try to lock it.
		}
	}
	// Take the lock or increase our expiration if we are already holding the lock
	cm.Data[pxOwnerKey] = id
	cm.Data[pxExpirationKey] = time.Now().Add(v1DefaultK8sLockTTL).Format(time.UnixDate)
	if _, err = core.Instance().UpdateConfigMap(cm); err != nil {
		return "", err
	}
	return id, nil
}

func (c *configMap) refreshLockV1(id string) {
	fn := "refreshLockV1"
	refresh := time.NewTicker(v1DefaultK8sLockRefreshDuration)
	var (
		currentRefresh time.Time
		prevRefresh    time.Time
		startTime      time.Time
	)
	startTime = time.Now()
	defer refresh.Stop()
	for {
		select {
		case <-refresh.C:
			c.kLockV1.Lock()
			for !c.kLockV1.unlocked {
				c.checkLockTimeout(c.lockHoldTimeoutV1, startTime, id)
				currentRefresh = time.Now()

				if _, err := c.tryLockV1(id, true); err != nil {
					configMapLog(fn, c.name, "", "", err).Errorf(
						"Error refreshing lock. [Owner %v] [Err: %v]"+
							" [Current Refresh: %v] [Previous Refresh: %v]",
						id, err, currentRefresh, prevRefresh,
					)
					if k8s_errors.IsConflict(err) {
						// try refreshing again
						continue
					}
					if errors.Is(err, ErrConfigMapLockLost) {
						// there is no coming back from this
						c.kLockV1.unlocked = true
						c.kLockV1.done = nil
						c.kLockV1.Unlock()
						return
					}
				}
				thresh := v1DefaultK8sLockRefreshDuration * 3 / 2
				if !prevRefresh.IsZero() && prevRefresh.Add(thresh).Before(currentRefresh) {
					configMapLog(fn, c.name, "", "", nil).Warnf(
						"V1 lock refresh is taking too long. [Owner %v] [Current Refresh: %v] [Previous Refresh: %v]",
						id, currentRefresh, prevRefresh)
				}
				prevRefresh = currentRefresh
				break
			}
			c.kLockV1.Unlock()
		case <-c.kLockV1.done:
			return
		}
	}
}
