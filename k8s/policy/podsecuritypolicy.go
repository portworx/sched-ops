package policy

import (
	"fmt"
)

// PodSecurityPolicyOps is an interface to perform k8s Pod Security Policy operations
// NOTE: PodSecurityPolicy was removed in Kubernetes v1.25 and is no longer supported.
// These methods are kept for backward compatibility but will return errors.
// Use Pod Security Admission (PSA) or a policy engine like OPA/Kyverno instead.
type PodSecurityPolicyOps interface {
	// CreatePodSecurityPolicy creates the given pod security policy
	// DEPRECATED: PodSecurityPolicy was removed in Kubernetes v1.25
	CreatePodSecurityPolicy(policy interface{}) (interface{}, error)
	// GetPodSecurityPolicy gets the given pod security policy
	// DEPRECATED: PodSecurityPolicy was removed in Kubernetes v1.25
	GetPodSecurityPolicy(name string) (interface{}, error)
	// ListPodSecurityPolicies list pods security policies
	// DEPRECATED: PodSecurityPolicy was removed in Kubernetes v1.25
	ListPodSecurityPolicies() (interface{}, error)
	// UpdatePodSecurityPolicy updates the give pod security policy
	// DEPRECATED: PodSecurityPolicy was removed in Kubernetes v1.25
	UpdatePodSecurityPolicy(policy interface{}) (interface{}, error)
	// DeletePodSecurityPolicy deletes the given pod security policy
	// DEPRECATED: PodSecurityPolicy was removed in Kubernetes v1.25
	DeletePodSecurityPolicy(name string) error
}

var errPSPRemoved = fmt.Errorf("PodSecurityPolicy was removed in Kubernetes v1.25 and is no longer supported. Use Pod Security Admission (PSA) or a policy engine like OPA/Kyverno instead")

// CreatePodSecurityPolicy creates the given pod security policy
// DEPRECATED: PodSecurityPolicy was removed in Kubernetes v1.25
func (c *Client) CreatePodSecurityPolicy(policy interface{}) (interface{}, error) {
	return nil, errPSPRemoved
}

// GetPodSecurityPolicy gets the given pod security policy
// DEPRECATED: PodSecurityPolicy was removed in Kubernetes v1.25
func (c *Client) GetPodSecurityPolicy(name string) (interface{}, error) {
	return nil, errPSPRemoved
}

// ListPodSecurityPolicies gets the given pod security policy
// DEPRECATED: PodSecurityPolicy was removed in Kubernetes v1.25
func (c *Client) ListPodSecurityPolicies() (interface{}, error) {
	return nil, errPSPRemoved
}

// UpdatePodSecurityPolicy updates the give pod security policy
// DEPRECATED: PodSecurityPolicy was removed in Kubernetes v1.25
func (c *Client) UpdatePodSecurityPolicy(policy interface{}) (interface{}, error) {
	return nil, errPSPRemoved
}

// DeletePodSecurityPolicy deletes the given pod security policy
// DEPRECATED: PodSecurityPolicy was removed in Kubernetes v1.25
func (c *Client) DeletePodSecurityPolicy(name string) error {
	return errPSPRemoved
}
