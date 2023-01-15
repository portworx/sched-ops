package rancher

import (
	"context"

	rancherv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ProjectOps is an interface to perfrom rancher project operations
type ProjectOps interface {
	// ListProjects lists all Projects in a given namespace
	ListProjects(namespace string) (*rancherv3.ProjectList, error)
}

// ListProjects lists all Projects in the given namespace
func (c *Client) ListProjects(namespace string) (*rancherv3.ProjectList, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	projectList := &rancherv3.ProjectList{}
	if err := c.crClient.List(
		context.TODO(),
		projectList,
		&client.ListOptions{
			Namespace: namespace,
		},
	); err != nil {
		return projectList, err
	}
	return projectList, nil
}
