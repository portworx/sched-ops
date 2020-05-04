package externalsnapshotter

import (
	"github.com/kubernetes-csi/external-snapshotter/v2/pkg/apis/volumesnapshot/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SnapshotClassOps is an interface to perform k8s VolumeSnapshotClass operations
type SnapshotClassOps interface {
	// CreateSnapshotClass creates the given snapshot class
	CreateSnapshotClass(snap *v1beta1.VolumeSnapshotClass) (*v1beta1.VolumeSnapshotClass, error)
	// GetSnapshotClass returns the snapshot class for given name
	GetSnapshotClass(name string) (*v1beta1.VolumeSnapshotClass, error)
	// ListSnapshotClasses lists all snapshot classes
	ListSnapshotClasses() (*v1beta1.VolumeSnapshotClassList, error)
	// UpdateSnapshotClass updates the given snapshot class
	UpdateSnapshotClass(snap *v1beta1.VolumeSnapshotClass) (*v1beta1.VolumeSnapshotClass, error)
	// DeleteSnapshotClass deletes the given snapshot class
	DeleteSnapshotClass(name string) error
}

// CreateSnapshotClass creates the given snapshot class.
func (c *Client) CreateSnapshotClass(snap *v1beta1.VolumeSnapshotClass) (*v1beta1.VolumeSnapshotClass, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	return c.client.VolumeSnapshotClasses().Create(snap)
}

// GetSnapshotClass returns the snapshot class for given name
func (c *Client) GetSnapshotClass(name string) (*v1beta1.VolumeSnapshotClass, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	return c.client.VolumeSnapshotClasses().Get(name, metav1.GetOptions{})
}

// ListSnapshotClasses lists all snapshot classes
func (c *Client) ListSnapshotClasses() (*v1beta1.VolumeSnapshotClassList, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	return c.client.VolumeSnapshotClasses().List(metav1.ListOptions{})
}

// UpdateSnapshotClass updates the given snapshot class
func (c *Client) UpdateSnapshotClass(snap *v1beta1.VolumeSnapshotClass) (*v1beta1.VolumeSnapshotClass, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	return c.client.VolumeSnapshotClasses().Update(snap)
}

// DeleteSnapshotClass deletes the given snapshot
func (c *Client) DeleteSnapshotClass(name string) error {
	if err := c.initClient(); err != nil {
		return err
	}
	return c.client.VolumeSnapshotClasses().Delete(name, &metav1.DeleteOptions{})
}
