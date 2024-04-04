package anthos

import (
	"context"

	v1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
)

// ClusterOps is an interface to perform k8s cluster operations
type ClusterOps interface {
	// ListCluster lists all kubernetes clusters
	ListCluster(ctx context.Context, options metav1.ListOptions) (*v1beta1.ClusterList, error)
	// GetCluster returns a cluster for the given name
	GetCluster(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.Cluster, error)
	// DescribeCluster gets the cluster status
	DescribeCluster(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.ClusterStatus, error)
}

// ListCluster lists all kubernetes clusters
func (c *Client) ListCluster(ctx context.Context, options metav1.ListOptions) (result *v1beta1.ClusterList, err error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	result = &v1beta1.ClusterList{}
	err = c.RESTClient().Get().
		Resource("clusters").
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)

	return
}

// GetCluster returns a cluster for the given name
func (c *Client) GetCluster(ctx context.Context, name string, options metav1.GetOptions) (result *v1beta1.Cluster, err error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	result = &v1beta1.Cluster{}
	err = c.RESTClient().Get().
		Resource("clusters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).Do(ctx).Into(result)
	return
}

// DescribeCluster describe the cluster status
func (c *Client) DescribeCluster(ctx context.Context, name string, options metav1.GetOptions) (*v1beta1.ClusterStatus, error) {
	cluster := &v1beta1.Cluster{}
	err := c.RESTClient().Get().
		Resource("clusters").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).Do(ctx).Into(cluster)

	return &cluster.Status, err
}
