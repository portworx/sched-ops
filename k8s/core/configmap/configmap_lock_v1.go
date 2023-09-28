package configmap

import (
	"fmt"
	"time"

	"github.com/portworx/sched-ops/k8s/core"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
)

func (c *configMap) Lock(id string) error {
	logrus.Infof("trying to get lock = %v", id)
	return c.LockWithHoldTimeout(id, c.defaultLockHoldTimeout)
}

func (c *configMap) LockWithHoldTimeout(id string, holdTimeout time.Duration) error {
	logrus.Infof("** LockWithHoldTimeout id = %v, holdTimeOut = %v", id, holdTimeout)
	fn := "LockWithHoldTimeout"
	count := uint(0)
	// try acquiring a lock on the ConfigMap
	owner, err := c.tryLockV1(id, false)
	logrus.Infof("owner = %v, err = %v", owner, err)
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
	logrus.Infof("** Unlock")
	logrus.Infof("c.kLockV1.id = %v, c.kLockV1.unlocked = %v", c.kLockV1.id, c.kLockV1.unlocked)
	fn := "Unlock"
	// Get the existing ConfigMap
	c.kLockV1.Lock()
	logrus.Infof("** got lock ****")
	defer time.Sleep(10 * time.Second)
	defer logrus.Infof("defer c.kLockV1.id = %v, c.kLockV1.unlocked = %v", c.kLockV1.id, c.kLockV1.unlocked)
	defer c.kLockV1.Unlock()

	// defer func() {
	// 	if r := recover(); r != nil {
	// 		// Handle the panic here
	// 		logrus.Infof("Panic occurred in Unlock:", r)
	// 	}
	// 	// Always unlock the mutex, even in case of a panic
	// 	c.kLockV1.Unlock()
	// }()

	if c.kLockV1.unlocked {
		// The lock is already unlocked
		logrus.Infof("** return -2")
		return nil
	}
	logrus.Infof("--1--")
	c.kLockV1.unlocked = true
	c.kLockV1.done <- struct{}{}
	// defer func() {
	// 	c.kLockV1.unlocked = true
	// 	c.kLockV1.done <- struct{}{}
	// }()

	logrus.Infof("--2--")
	logrus.Infof("c.kLockV1.id = %v, c.kLockV1.unlocked = %v", c.kLockV1.id, c.kLockV1.unlocked)
	var (
		err error
		cm  *corev1.ConfigMap
	)
	logrus.Infof("3")
	logrus.Infof("c.kLockV1.id = %v, c.kLockV1.unlocked =%v", c.kLockV1.id, c.kLockV1.unlocked)
	for retries := 0; retries < maxConflictRetries; retries++ {
		logrus.Infof("retries = %v", retries)
		cm, err = core.Instance().GetConfigMap(
			c.name,
			k8sSystemNamespace,
		)
		if err != nil {
			// A ConfigMap should always be created.
			logrus.Info("** return -1")
			return err
		}

		currentOwner := cm.Data[pxOwnerKey]
		logrus.Infof("** currentOwner = %v", currentOwner)
		if currentOwner != c.kLockV1.id {
			// We are currently not holding the lock
			logrus.Infof("return 0")
			return nil
		}
		logrus.Info("*** delete")
		delete(cm.Data, pxOwnerKey)
		delete(cm.Data, pxExpirationKey)
		logrus.Info("** deleted")
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
		logrus.Info("return 1")
		return nil
	}
	logrus.Info("4")
	logrus.Info("return 2")
	return err
}

func (c *configMap) tryLockV1(id string, refresh bool) (string, error) {
	// Get the existing ConfigMap
	logrus.Info("*** inside tryLockV1")
	logrus.Infof("*** Unlocked status = %v", c.kLockV1.unlocked)
	logrus.Infof("*** inside tryLockV1 id = %v, refresh = %v", id, refresh)

	// res := c.kLockV1.TryLock()
	// for !res {
	// 	fmt.Printf("*** inside tryLockV1 %+v, id=%v \n", c.kLockV1, id)
	// 	res = c.kLockV1.TryLock()
	// 	logrus.Infof("*** res =", res)
	// }

	c.kLockV1.Lock()
	defer c.kLockV1.Unlock()
	cm, err := core.Instance().GetConfigMap(
		c.name,
		k8sSystemNamespace,
	)
	if err != nil {
		// A ConfigMap should always be created.
		logrus.Infof("return 1")
		return "", err
	}

	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}

	currentOwner := cm.Data[pxOwnerKey]
	logrus.Infof("currentOwner = %v, id = %v", currentOwner, id)
	if currentOwner != "" {
		logrus.Infof("currentOwner = %v", currentOwner)
		if currentOwner == id && refresh {
			// We already hold the lock just refresh
			// our expiry
			logrus.Infof("** increasing expiry = %v", id)
			goto increase_expiry
		} // refresh not requested
		// Someone might have a lock on the cm
		// Check expiration
		expiration := cm.Data[pxExpirationKey]
		if expiration != "" {
			logrus.Infof("expiration = %v", expiration)
			expiresAt, err := time.Parse(time.UnixDate, expiration)
			if err != nil {
				logrus.Infof("return 2")
				return currentOwner, err
			}
			logrus.Infof("expiresAt = %v", expiresAt)
			logrus.Infof("time.Now() = %v, time.Now().Before(expiresAt) = %v", time.Now(), time.Now().Before(expiresAt))
			if time.Now().Before(expiresAt) {
				// Lock is currently held by the owner
				// Retry after sometime
				logrus.Infof("returning here")
				logrus.Infof("return 3")
				return currentOwner, ErrConfigMapLocked
			} // else lock is expired. Try to lock it.
		}
	}

	// Take the lock or increase our expiration if we are already holding the lock
	cm.Data[pxOwnerKey] = id
increase_expiry:
	logrus.Infof("previous  = %v", cm.Data[pxExpirationKey])
	cm.Data[pxExpirationKey] = time.Now().Add(v1DefaultK8sLockTTL).Format(time.UnixDate)
	logrus.Infof("now  = %v", cm.Data[pxExpirationKey])
	if _, err = core.Instance().UpdateConfigMap(cm); err != nil {
		logrus.Infof("return 4")
		return "", err
	}
	logrus.Infof("return 5")
	return id, nil
}

func (c *configMap) refreshLockV1(id string) {
	logrus.Infof("** refreshLock = %v", id)
	fn := "refreshLock"
	refresh := time.NewTicker(v1DefaultK8sLockRefreshDuration)
	var (
		currentRefresh time.Time
		prevRefresh    time.Time
		startTime      time.Time
	)
	startTime = time.Now()
	logrus.Infof(" $$$$ startTime $$$$ = %v", startTime)
	defer refresh.Stop()
	for {
		select {
		case <-refresh.C:
			logrus.Infof("inside refresh = %v", c.kLockV1.id)
			// c.kLockV1.Lock()
			// defer c.kLockV1.Unlock()
			for !c.kLockV1.unlocked {
				logrus.Info("*** inside refrest switch")
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
						logrus.Infof("err = %v", err)
						continue
					}
				}
				logrus.Infof(" -- updated In Refresh -- ")
				prevRefresh = currentRefresh
				break
			}
			// c.kLockV1.Unlock()
		case <-c.kLockV1.done:
			logrus.Infof("** exiting refresh for = %v", id)
			return
		}
	}
}
