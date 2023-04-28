package configmap

import (
	"time"

	"github.com/pborman/uuid"
	log "github.com/sirupsen/logrus"
)

func (c *configMap) Instance() ConfigMap {

	//iskvdbhealthy
	//kvdb.instance

	if c.pxNs == c.config.nameSpace {
		//fresh install
		//upgrade completed
		return c.config
	} else {
		existingConfig := c.config
		c.copylock.Lock(uuid.New())
		defer c.copylock.Unlock()
		lockMap, err := c.copylock.Get()
		if err != nil {
			log.Error("Error during fetching data from copy lock %s", err)
			return existingConfig
		}
		status := lockMap["UPGRADE_DONE"]
		if status == "true" {
			// upgrade is completed
			//create configmap in portworx namespace
			newConfig := &coreConfigMap{
				name:                   existingConfig.name,
				defaultLockHoldTimeout: existingConfig.defaultLockHoldTimeout,
				kLocksV2:               existingConfig.kLocksV2,
				lockAttempts:           existingConfig.lockAttempts,
				lockRefreshDuration:    existingConfig.lockRefreshDuration,
				lockK8sLockTTL:         existingConfig.lockK8sLockTTL,
				nameSpace:              "portworx",
			}

			configData, err := existingConfig.Get()
			if err != nil {
				log.Errorf("Error during fetching data from old config map %s", err)
				return existingConfig
			}
			//copy data from old configmap to new configmap
			if err = newConfig.Update(configData); err != nil {
				log.Errorf("Error during copying data from old config map %s", err)
				return existingConfig
			}

			//delete old configmap
			err = c.config.Delete()
			if err != nil {
				log.Errorf("Error during deleting configmap %s in namespace %s ", c.config.name, c.config.nameSpace)
			}
			c.config = newConfig
		} else {
			return existingConfig
		}
	}
	return c.config
}

func (c *configMap) Get() (map[string]string, error) {
	return c.Instance().Get()
}

func (c *configMap) Delete() error {
	return c.Instance().Delete()
}

func (c *configMap) Patch(data map[string]string) error {
	return c.Instance().Patch(data)
}

func (c *configMap) Update(data map[string]string) error {
	return c.Instance().Update(data)
}

func (c *configMap) Lock(id string) error {
	return c.Instance().Lock(id)
}

func (c *configMap) LockWithHoldTimeout(id string, holdTimeout time.Duration) error {
	return c.Instance().LockWithHoldTimeout(id, holdTimeout)
}

func (c *configMap) LockWithKey(owner, key string) error {
	return c.Instance().LockWithKey(owner, key)
}

func (c *configMap) Unlock() error {
	return c.Instance().Unlock()
}

func (c *configMap) UnlockWithKey(key string) error {
	return c.Instance().UnlockWithKey(key)
}
func (c *configMap) IsKeyLocked(key string) (bool, string, error) {
	return c.Instance().IsKeyLocked(key)
}
