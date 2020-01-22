package storage

import (
	storagebeta "k8s.io/api/storage/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VolumeAttachmentOps is an interface to perform k8s VolumeAttachmentOps operations
type VolumeAttachmentOps interface {
	// ListVolumeAttachments lists all volume attachments
	ListVolumeAttachments() (*storagebeta.VolumeAttachmentList, error)
	// DeleteVolumeAttachment deletes a given Volume Attachment by name
	DeleteVolumeAttachment(name string) error
	// CreateVolumeAttachment creates a volume attachment
	CreateVolumeAttachment(*storagebeta.VolumeAttachment) (*storagebeta.VolumeAttachment, error)
	// UpdateVolumeAttachment updates a volume attachment
	UpdateVolumeAttachment(*storagebeta.VolumeAttachment) (*storagebeta.VolumeAttachment, error)
	// UpdateVolumeAttachmentStatus updates a volume attachment status
	UpdateVolumeAttachmentStatus(*storagebeta.VolumeAttachment) (*storagebeta.VolumeAttachment, error)
}

// ListVolumeAttachments lists all volume attachments
func (c *Client) ListVolumeAttachments() (*storagebeta.VolumeAttachmentList, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.storagebeta.VolumeAttachments().List(metav1.ListOptions{})
}

// DeleteVolumeAttachment deletes a given Volume Attachment by name
func (c *Client) DeleteVolumeAttachment(name string) error {
	if err := c.initClient(); err != nil {
		return err
	}

	return c.storagebeta.VolumeAttachments().Delete(name, &metav1.DeleteOptions{})
}

// CreateVolumeAttachment creates a volume attachment
func (c *Client) CreateVolumeAttachment(volumeAttachment *storagebeta.VolumeAttachment) (*storagebeta.VolumeAttachment, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.storagebeta.VolumeAttachments().Create(volumeAttachment)
}

// UpdateVolumeAttachment updates a volume attachment
func (c *Client) UpdateVolumeAttachment(volumeAttachment *storagebeta.VolumeAttachment) (*storagebeta.VolumeAttachment, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.storagebeta.VolumeAttachments().Update(volumeAttachment)
}

// UpdateVolumeAttachmentStatus updates a volume attachment status
func (c *Client) UpdateVolumeAttachmentStatus(volumeAttachment *storagebeta.VolumeAttachment) (*storagebeta.VolumeAttachment, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.storagebeta.VolumeAttachments().UpdateStatus(volumeAttachment)
}
