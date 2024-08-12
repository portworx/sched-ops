package configmap

import (
	"errors"
	"regexp"
	"sync"
	"time"
)

const (
	// DefaultK8sLockAttempts is the number of times to try taking the lock before failing.
	// It defaults to 300, the same number for a kvdb lock.
	DefaultK8sLockAttempts = 300
	// DefaultK8sLockTimeout is time duration within which a lock should be released
	// else it assumes that the node is stuck and panics.
	DefaultK8sLockTimeout = 3 * time.Minute
	// v1DefaultK8sLockTTL is the time duration after which the lock will expire
	v1DefaultK8sLockTTL = 16 * time.Second
	// v1DefaultK8sLockRefreshDuration is the time duration after which a lock is refreshed
	v1DefaultK8sLockRefreshDuration = 8 * time.Second
	// v2DefaultK8sLockTTL is the time duration after which the lock will expire
	v2DefaultK8sLockTTL = 60 * time.Second
	// v2DefaultK8sLockRefreshDuration is the time duration after which a lock is refreshed
	v2DefaultK8sLockRefreshDuration = 20 * time.Second
	// k8sSystemNamespace is the namespace in which we create the ConfigMap
	k8sSystemNamespace = "kube-system"

	// ***********************
	//   ConfigMap Lock Keys
	// ***********************

	// pxOwnerKey is key which indicates the node holding the ConfigMap lock.
	// This is specifically for the deprecated Lock and Unlock methods.
	pxOwnerKey = "px-owner"
	// pxExpirationKey is the key which indicates the time at which the
	// current lock will expire.
	// This is specifically for the deprecated Lock and Unlock methods.
	pxExpirationKey = "px-expiration"

	// pxLockKey is the key which stores the lock data. The data in this key is stored in JSON as an array of lockData
	// objects.
	pxLockKey = "px-lock"

	// pxGenerationKey stores the generation of the configmap data. The value is incremented every time
	// the configmap data is updated via PatchKeyLocked or DeleteKeyLocked. This is used for diagnostics purposes only.
	pxGenerationKey = "px-generation"

	lockSleepDuration     = 1 * time.Second
	configMapUserLabelKey = "user"
	maxConflictRetries    = 3
)

var (
	// ErrConfigMapLocked is returned when the ConfigMap is locked
	ErrConfigMapLocked = errors.New("ConfigMap is locked")
	// ErrConfigMapLockLost is returned when the ConfigMap lock expired and was taken away
	ErrConfigMapLockLost = errors.New("ConfigMap lock was lost")
	fatalCb              FatalCb
	configMapNameRegex   = regexp.MustCompile("[^a-zA-Z0-9]+")
)

// FatalCb is a callback function which will be executed if the Lock
// routine encounters a panic situation
type FatalCb func(format string, args ...interface{})

// CheckLockCb is a callback function to check if the lock is still valid.
type CheckLockCb func(data map[string]string) bool

type configMap struct {
	name                   string
	kLockV1                k8sLock
	kLocksV2Mutex          sync.Mutex
	kLocksV2               map[string]*k8sLock
	lockHoldTimeoutV1      time.Duration
	defaultLockHoldTimeout time.Duration
	lockAttempts           uint
	lockRefreshDuration    time.Duration
	lockK8sLockTTL         time.Duration
}

type k8sLock struct {
	wg       sync.WaitGroup
	done     chan struct{}
	unlocked bool
	id       string
	sync.Mutex
}

// ConfigMap is an interface that provides a set of APIs over a single
// k8s configmap object. The data in the configMap is managed as a map of string
// to string.
//
// Rules:
//  1. Locks are non-reentrant.
//  2. For v1 locks, use Lock()/LockWithParams() and Unlock(). For v2 locks, use LockWithKey() and UnlockWithKey().
//  3. Use the same locking mechanism (v1 or v2) for a given key consistently. It is dangerous to change
//     the locking mechanism for a key (e.g. during upgrade).
//  4. Specify the correct locking mechanism and lock owner when using PatchKeyLocked() and DeleteKeyLocked().
//  5. Do not patch/delete any of the internal keys used for the locks and expirations.
type ConfigMap interface {
	// Lock locks a configMap where id is the identification of
	// the holder of the lock. Lock is non-reentrant.
	Lock(id string) error
	// LockWithParams similar to Lock but with custom params.
	// If lockAttempts is 0, the value passed to configmap.New() is used.
	LockWithParams(id string, holdTimeout time.Duration, lockAttempts uint) error
	// LockWithKey locks a configMap where owner is the identification
	// of the holder of the lock and key is the specific lock to take.  Lock is non-reentrant.
	LockWithKey(owner, key string) error
	// Unlock unlocks the configMap.
	Unlock() error
	// UnlockWithKey unlocks the given key in the configMap.
	UnlockWithKey(key string) error
	// IsKeyLocked returns if the given key is locked, and if so, by which owner.
	IsKeyLocked(key string) (bool, string, error)

	// PatchKeyLocked updates the specified key in the configMap. It verifies that
	// the lock is still held by the specified owner. Lock needs to be held by the lockOwner
	// throughout the patch operation for this function to succeed.
	PatchKeyLocked(isV1Lock bool, lockOwner, key, val string) error

	// DeleteKeyLocked deletes the specified key in the configMap.
	DeleteKeyLocked(isV1Lock bool, lockOwner, key string) error

	// Get returns the contents of the configMap
	Get() (map[string]string, error)
	// Delete deletes the configMap
	Delete() error
}

// lockData structs are serialized into JSON and stored as a list inside a ConfigMap.
// Each lockData struct contains the owner (usually which node took the lock), key
// (which specific lock it's taking), and an expiration time after which the lock is invalid.
type lockData struct {
	Owner      string    `json:"owner"`
	Key        string    `json:"key"`
	Expiration time.Time `json:"expiration"`
}
