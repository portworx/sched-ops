package configmap

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/portworx/sched-ops/k8s/core"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// New returns the ConfigMap interface. It also creates a new
// configmap in k8s for the given name if not present and puts the data in it.
func New(
	name string,
	data map[string]string,
	lockTimeout time.Duration,
	lockAttempts uint,
	v2LockRefreshDuration time.Duration,
	v2LockK8sLockTTL time.Duration,
) (ConfigMap, error) {
	if data == nil {
		data = make(map[string]string)
	}

	labels := map[string]string{
		configMapUserLabelKey: TruncateLabel(name),
	}
	data[pxOwnerKey] = ""

	cm := &corev1.ConfigMap{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      name,
			Namespace: k8sSystemNamespace,
			Labels:    labels,
		},
		Data: data,
	}

	if _, err := core.Instance().CreateConfigMap(cm); err != nil &&
		!k8s_errors.IsAlreadyExists(err) {
		return nil, fmt.Errorf("failed to create configmap %v: %v",
			name, err)
	}

	if v2LockK8sLockTTL == 0 {
		v2LockK8sLockTTL = v2DefaultK8sLockTTL
	}

	if v2LockRefreshDuration == 0 {
		v2LockRefreshDuration = v2DefaultK8sLockRefreshDuration
	}

	return &configMap{
		name:                   name,
		defaultLockHoldTimeout: lockTimeout,
		kLocksV2:               map[string]*k8sLock{},
		lockAttempts:           lockAttempts,
		lockRefreshDuration:    v2LockRefreshDuration,
		lockK8sLockTTL:         v2LockK8sLockTTL,
	}, nil
}

func (c *configMap) Get() (map[string]string, error) {
	cm, err := core.Instance().GetConfigMap(
		c.name,
		k8sSystemNamespace,
	)
	if err != nil {
		return nil, err
	}

	return cm.Data, nil
}

func (c *configMap) Delete() error {
	return core.Instance().DeleteConfigMap(
		c.name,
		k8sSystemNamespace,
	)
}

func (c *configMap) PatchKeyLocked(isV1Lock bool, lockOwner, key, val string) error {
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
			return err
		}
		if err := c.checkLockOwner(isV1Lock, key, lockOwner, cm); err != nil {
			return fmt.Errorf("lock check failed: %w", err)
		}
		if cm.Data == nil {
			cm.Data = make(map[string]string, 0)
		}

		cm.Data[key] = val
		newGen := c.incrementGeneration(cm)
		_, err = core.Instance().UpdateConfigMap(cm)
		if k8s_errors.IsConflict(err) {
			// try again
			continue
		}
		if err == nil {
			logrus.Infof("Updated key %s in configmap %s/%s with generation %d and lockOwner %s",
				key, k8sSystemNamespace, c.name, newGen, lockOwner)
		}
		return err
	}
	return err
}

func (c *configMap) incrementGeneration(cm *corev1.ConfigMap) uint64 {
	val := cm.Data[pxGenerationKey]
	if val == "" {
		val = "0"
	}
	gen, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		logrus.Errorf("Failed to parse generation %s; resetting to 1: %v", val, err)
		val = "0"
		gen = 0
	}
	newGen := gen + 1
	if newGen == 0 {
		logrus.Warnf("Resetting generation to 1 from %s", val)
		newGen = 1
	}
	cm.Data[pxGenerationKey] = strconv.FormatUint(newGen, 10)
	return newGen
}

func (c *configMap) checkLockOwner(isV1Lock bool, v2Key string, expectedLockOwner string, cm *corev1.ConfigMap) error {
	if expectedLockOwner == "" {
		return errors.New("expected lock owner cannot be empty")
	}
	if isV1Lock {
		if cm.Data[pxOwnerKey] != expectedLockOwner {
			return fmt.Errorf("v1 lock owner is %q instead of expected %q", cm.Data[pxOwnerKey], expectedLockOwner)
		}
		return nil
	}
	if v2Key == "" {
		return errors.New("v2 key cannot be empty when checking v2 lock owner")
	}
	currentOwner, err := c.getV2LockOwnerIncludeExpired(cm, v2Key)
	if err != nil {
		return fmt.Errorf("failed to parse v2 locks: %w", err)
	}
	if currentOwner != expectedLockOwner {
		return fmt.Errorf("v2 lock owner is %q instead of expected %q", currentOwner, expectedLockOwner)
	}
	return nil
}

func (c *configMap) DeleteKeyLocked(isV1Lock bool, lockOwner string, key string) error {
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
			return err
		}
		if err := c.checkLockOwner(isV1Lock, key, lockOwner, cm); err != nil {
			return fmt.Errorf("lock check failed: %w", err)
		}
		delete(cm.Data, key)

		newGen := c.incrementGeneration(cm)
		_, err = core.Instance().UpdateConfigMap(cm)
		if k8s_errors.IsConflict(err) {
			// try again
			continue
		}
		if err == nil {
			logrus.Infof("Deleted key %s in configmap %s/%s with generation %d and lockOwner %s",
				key, k8sSystemNamespace, c.name, newGen, lockOwner)
		}
		return err
	}
	return err
}

// SetFatalCb sets the fatal callback for the package which will get invoked in panic situations
func SetFatalCb(fb FatalCb) {
	fatalCb = fb
}

func configMapLog(fn, name, owner, key string, err error) *logrus.Entry {
	if len(owner) > 0 && len(key) > 0 {
		return logrus.WithFields(logrus.Fields{
			"Module":   "ConfigMap",
			"Name":     name,
			"Owner":    owner,
			"Key":      key,
			"Function": fn,
			"Error":    err,
		})
	}
	if len(owner) > 0 {
		return logrus.WithFields(logrus.Fields{
			"Module":   "ConfigMap",
			"Name":     name,
			"Owner":    owner,
			"Function": fn,
			"Error":    err,
		})
	}
	return logrus.WithFields(logrus.Fields{
		"Module":   "ConfigMap",
		"Name":     name,
		"Function": fn,
		"Error":    err,
	})
}

// GetName is a helper function that returns a valid k8s
// configmap name given a prefix identifying the component using
// the configmap and a clusterID
func GetName(prefix, clusterID string) string {
	return prefix + strings.ToLower(configMapNameRegex.ReplaceAllString(clusterID, ""))
}

// TruncateLabel is a helper function that returns a valid k8s
// label stripped down to 63 characters. It removes the trailing characters
func TruncateLabel(label string) string {
	if len(label) > 63 {
		return label[:63]
	}
	return label
}
