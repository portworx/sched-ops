package storage

import (
	"context"

	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CsiDriversOps interface {
	// ListCSIDrivers returns all CSIDrivers
	ListCSIDrivers() (*storagev1.CSIDriverList, error)
}

func (c *Client) ListCSIDrivers() (*storagev1.CSIDriverList, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.storage.CSIDrivers().List(context.Background(), metav1.ListOptions{})
}
