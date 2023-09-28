package configmap

import (
	"fmt"
	"time"

	"github.com/portworx/sched-ops/k8s/core"
	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
)

func (c *configMap) Lock(id string) error {
	fmt.Println("trying to get lock", id)
	return c.LockWithHoldTimeout(id, c.defaultLockHoldTimeout)
}

func (c *configMap) LockWithHoldTimeout(id string, holdTimeout time.Duration) error {
	fmt.Println("** LockWithHoldTimeout id, holdTimeOut", id, holdTimeout)
	fn := "LockWithHoldTimeout"
	count := uint(0)
	// try acquiring a lock on the ConfigMap
	owner, err := c.tryLockV1(id, false)
	fmt.Println("owner, err", owner, err)
	// This is the same no. of times (300) we try while acquiring a kvdb lock
	for maxCount := c.lockAttempts; err != nil && count < maxCount; count++ {
		time.Sleep(lockSleepDuration)
		fmt.Printf("count=%v \n", count)
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
	c.lockHoldTimeoutV1 = holdTimeout
	// c.kLockV1 = k8sLock{done: make(chan struct{}), id: id}
	c.kLockV1.id = id

	c.kLockV1.unlocked = false
	go c.refreshLockV1(id)
	return nil
}

func (c *configMap) Unlock() error {
	fmt.Println("** Unlock")
	fmt.Println("c.kLockV1.id, c.kLockV1.unlocked", c.kLockV1.id, c.kLockV1.unlocked)
	fn := "Unlock"
	// Get the existing ConfigMap
	c.kLockV1.Lock()
	fmt.Println("** got lock ****")
	defer c.kLockV1.Unlock()
	if c.kLockV1.unlocked {
		// The lock is already unlocked
		fmt.Println("** return -2")
		return nil
	}
	fmt.Println("1")
	c.kLockV1.unlocked = true
	c.kLockV1.done <- struct{}{}
	fmt.Println("2")
	fmt.Println("c.kLockV1.id, c.kLockV1.unlocked", c.kLockV1.id, c.kLockV1.unlocked)
	var (
		err error
		cm  *corev1.ConfigMap
	)
	fmt.Println("3")
	fmt.Println("c.kLockV1.id, c.kLockV1.unlocked", c.kLockV1.id, c.kLockV1.unlocked)
	for retries := 0; retries < maxConflictRetries; retries++ {
		fmt.Println("retries=", retries)
		cm, err = core.Instance().GetConfigMap(
			c.name,
			k8sSystemNamespace,
		)
		if err != nil {
			// A ConfigMap should always be created.
			fmt.Println("** return -1")
			return err
		}

		currentOwner := cm.Data[pxOwnerKey]
		fmt.Println("** currentOwner", currentOwner)
		if currentOwner != c.kLockV1.id {
			// We are currently not holding the lock
			fmt.Println("return 0")
			return nil
		}
		fmt.Println("*** delete")
		delete(cm.Data, pxOwnerKey)
		delete(cm.Data, pxExpirationKey)
		fmt.Println("** deleted")
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
		fmt.Println("return 1")
		return nil
	}
	fmt.Println("4")
	fmt.Println("return 2")
	return err
}

func (c *configMap) tryLockV1(id string, refresh bool) (string, error) {
	// Get the existing ConfigMap
	fmt.Println("*** inside tryLockV1 id, refresh", id, refresh)
	cm, err := core.Instance().GetConfigMap(
		c.name,
		k8sSystemNamespace,
	)
	if err != nil {
		// A ConfigMap should always be created.
		fmt.Println("return 1")
		return "", err
	}

	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	currentOwner := cm.Data[pxOwnerKey]
	fmt.Println("currentOwner, id", currentOwner, id)
	if currentOwner != "" {
		fmt.Println("currentOwner", currentOwner)
		if currentOwner == id && refresh {
			// We already hold the lock just refresh
			// our expiry
			fmt.Println("** increasing expiry", id)
			goto increase_expiry
		} // refresh not requested
		// Someone might have a lock on the cm
		// Check expiration
		expiration := cm.Data[pxExpirationKey]
		if expiration != "" {
			fmt.Println("expiration-", expiration)
			expiresAt, err := time.Parse(time.UnixDate, expiration)
			if err != nil {
				fmt.Println("return 2")
				return currentOwner, err
			}
			fmt.Println("expiresAt", expiresAt)
			fmt.Println("time.Now().Before(expiresAt)", time.Now().Before(expiresAt))
			if time.Now().Before(expiresAt) {
				// Lock is currently held by the owner
				// Retry after sometime
				fmt.Println("returning here")
				fmt.Println("return 3")
				return currentOwner, ErrConfigMapLocked
			} // else lock is expired. Try to lock it.
		}
	}

	// Take the lock or increase our expiration if we are already holding the lock
	cm.Data[pxOwnerKey] = id
increase_expiry:
	fmt.Println("previous = ", cm.Data[pxExpirationKey])
	cm.Data[pxExpirationKey] = time.Now().Add(v1DefaultK8sLockTTL).Format(time.UnixDate)
	fmt.Println("now = ", cm.Data[pxExpirationKey])
	if _, err = core.Instance().UpdateConfigMap(cm); err != nil {
		fmt.Println("return 4")
		return "", err
	}
	fmt.Println("return 5")
	return id, nil
}

func (c *configMap) refreshLockV1(id string) {
	fmt.Println("** refreshLock", id)
	fn := "refreshLock"
	refresh := time.NewTicker(v1DefaultK8sLockRefreshDuration)
	var (
		currentRefresh time.Time
		prevRefresh    time.Time
		startTime      time.Time
	)
	startTime = time.Now()
	fmt.Println(" $$$$ startTime $$$$", startTime)
	defer refresh.Stop()
	for {
		select {
		case <-refresh.C:
			fmt.Println("inside refresh", c.kLockV1.id)
			c.kLockV1.Lock()
			// defer c.kLockV1.Unlock()
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
