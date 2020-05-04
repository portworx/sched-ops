package externalsnapshotter

import (
	"github.com/kubernetes-csi/external-snapshotter/v2/pkg/apis/volumesnapshot/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SnapshotOps is an interface to perform k8s VolumeSnapshot operations
type SnapshotOps interface {
	// CreateSnapshot creates the given snapshot
	CreateSnapshot(snap *v1beta1.VolumeSnapshot) (*v1beta1.VolumeSnapshot, error)
	// GetSnapshot returns the snapshot for given name and namespace
	GetSnapshot(name string, namespace string) (*v1beta1.VolumeSnapshot, error)
	// ListSnapshots lists all snapshots in the given namespace
	ListSnapshots(namespace string) (*v1beta1.VolumeSnapshotList, error)
	// UpdateSnapshot updates the given snapshot
	UpdateSnapshot(snap *v1beta1.VolumeSnapshot) (*v1beta1.VolumeSnapshot, error)
	// DeleteSnapshot deletes the given snapshot
	DeleteSnapshot(name string, namespace string) error
}

// CreateSnapshot creates the given snapshot.
func (c *Client) CreateSnapshot(snap *v1beta1.VolumeSnapshot) (*v1beta1.VolumeSnapshot, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	return c.client.VolumeSnapshots(snap.Namespace).Create(snap)
}

// GetSnapshot returns the snapshot for given name and namespace
func (c *Client) GetSnapshot(name string, namespace string) (*v1beta1.VolumeSnapshot, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	return c.client.VolumeSnapshots(namespace).Get(name, metav1.GetOptions{})
}

// ListSnapshots lists all snapshots in the given namespace
func (c *Client) ListSnapshots(namespace string) (*v1beta1.VolumeSnapshotList, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	return c.client.VolumeSnapshots(namespace).List(metav1.ListOptions{})
}

// UpdateSnapshot updates the given snapshot
func (c *Client) UpdateSnapshot(snap *v1beta1.VolumeSnapshot) (*v1beta1.VolumeSnapshot, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	return c.client.VolumeSnapshots(snap.Namespace).Update(snap)
}

// DeleteSnapshot deletes the given snapshot
func (c *Client) DeleteSnapshot(name string, namespace string) error {
	if err := c.initClient(); err != nil {
		return err
	}
	return c.client.VolumeSnapshots(namespace).Delete(name, &metav1.DeleteOptions{})
}
