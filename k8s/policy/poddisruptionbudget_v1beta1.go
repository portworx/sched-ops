package policy

import (
	"context"

	policyv1beta1 "k8s.io/api/policy/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodDisruptionBudgetOps is an interface to perform k8s Pod Disruption Budget operations
type PodDisruptionBudgetV1Beta1Ops interface {
	// CreatePodDisruptionBudget creates the given pod disruption budget
	CreatePodDisruptionBudgetV1beta1(policy *policyv1beta1.PodDisruptionBudget) (*policyv1beta1.PodDisruptionBudget, error)
	// GetPodDisruptionBudget gets the given pod disruption budget
	GetPodDisruptionBudgetV1beta1(name, namespace string) (*policyv1beta1.PodDisruptionBudget, error)
	// ListPodDisruptionBudget lists the pod disruption budgets
	ListPodDisruptionBudgetV1beta1(namespace string) (*policyv1beta1.PodDisruptionBudgetList, error)
	// UpdatePodDisruptionBudget updates the given pod disruption budget
	UpdatePodDisruptionBudgetV1beta1(policy *policyv1beta1.PodDisruptionBudget) (*policyv1beta1.PodDisruptionBudget, error)
	// DeletePodDisruptionBudget deletes the given pod disruption budget
	DeletePodDisruptionBudgetV1beta1(name, namespace string) error
}

// CreatePodDisruptionBudget creates the given pod disruption budget
func (c *Client) CreatePodDisruptionBudgetV1beta1(podDisruptionBudget *policyv1beta1.PodDisruptionBudget) (*policyv1beta1.PodDisruptionBudget, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.policy.PolicyV1beta1().PodDisruptionBudgets(podDisruptionBudget.Namespace).Create(context.TODO(), podDisruptionBudget, metav1.CreateOptions{})
}

// GetPodDisruptionBudget gets the given pod disruption budget
func (c *Client) GetPodDisruptionBudgetV1beta1(name, namespace string) (*policyv1beta1.PodDisruptionBudget, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.policy.PolicyV1beta1().PodDisruptionBudgets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// ListPodDisruptionBudget gets the given pod disruption budget
func (c *Client) ListPodDisruptionBudgetV1beta1(namespace string) (*policyv1beta1.PodDisruptionBudgetList, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.policy.PolicyV1beta1().PodDisruptionBudgets(namespace).List(context.TODO(), metav1.ListOptions{})
}

// UpdatePodDisruptionBudget updates the given pod disruption budget
func (c *Client) UpdatePodDisruptionBudgetV1beta1(podDisruptionBudget *policyv1beta1.PodDisruptionBudget) (*policyv1beta1.PodDisruptionBudget, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.policy.PolicyV1beta1().PodDisruptionBudgets(podDisruptionBudget.Namespace).Update(context.TODO(), podDisruptionBudget, metav1.UpdateOptions{})
}

// DeletePodDisruptionBudget deletes the given pod disruption budget
func (c *Client) DeletePodDisruptionBudgetV1beta1(name, namespace string) error {
	if err := c.initClient(); err != nil {
		return err
	}

	return c.policy.PolicyV1beta1().PodDisruptionBudgets(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}
