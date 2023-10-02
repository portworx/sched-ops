package kubevirt

import (
	"context"

	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubevirtv1 "kubevirt.io/api/core/v1"
)

// VirtualMachineInstanceOps is an interface to perform kubevirt VirtualMachineInstance operations
type VirtualMachineInstanceOps interface {
	// CreateVirtualMachineInstance calls VirtualMachineInstance create client method
	CreateVirtualMachineInstance(ctx context.Context, vmi *kubevirtv1.VirtualMachineInstance) (*kubevirtv1.VirtualMachineInstance, error)
	// GetVirtualMachineInstance Get updated VirtualMachineInstance from client matching name and namespace
	GetVirtualMachineInstance(ctx context.Context, name string, namespace string) (*kubevirtv1.VirtualMachineInstance, error)
}

// CreateVirtualMachineInstance calls VirtualMachineInstance create client method
func (c *Client) CreateVirtualMachineInstance(ctx context.Context, vmi *kubevirtv1.VirtualMachineInstance) (*kubevirtv1.VirtualMachineInstance, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.kubevirt.VirtualMachineInstance(vmi.GetNamespace()).Create(ctx, vmi)
}

// GetVirtualMachineInstance Get updated VirtualMachineInstance from client matching name and namespace
func (c *Client) GetVirtualMachineInstance(ctx context.Context, name string, namespace string) (*kubevirtv1.VirtualMachineInstance, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.kubevirt.VirtualMachineInstance(namespace).Get(ctx, name, &k8smetav1.GetOptions{})
}
