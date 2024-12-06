package batch

import (
	"context"
	"fmt"
	"time"

	"github.com/portworx/sched-ops/task"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JobOps is an interface to perform job operations
type JobOps interface {
	// CreateJob creates the given job
	CreateJob(job *batchv1.Job) (*batchv1.Job, error)
	// GetJob returns the job from given namespace and name
	GetJob(name, namespace string) (*batchv1.Job, error)
	// DeleteJob deletes the job with given namespace and name
	DeleteJob(name, namespace string) error
	// UpdateJob updates the given job
	UpdateJob(job *batchv1.Job) (*batchv1.Job, error)
	// DeleteJobWithForce deletes the job with given namespace and name
	DeleteJobWithForce(name, namespace string) error
	// ValidateJob validates if the job with given namespace and name succeeds.
	// It waits for timeout duration for job to succeed
	ValidateJob(name, namespace string, timeout time.Duration) error
	// ListAllJobs returns the jobs from given namespace
	ListAllJobs(namespace string, options metav1.ListOptions) (*batchv1.JobList, error)
}

// CreateJob creates the given job
func (c *Client) CreateJob(job *batchv1.Job) (*batchv1.Job, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.batch.Jobs(job.Namespace).Create(context.TODO(), job, metav1.CreateOptions{})
}

// GetJob returns the job from given namespace and name
func (c *Client) GetJob(name, namespace string) (*batchv1.Job, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.batch.Jobs(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

// DeleteJob deletes the job with given namespace and name
func (c *Client) DeleteJob(name, namespace string) error {
	if err := c.initClient(); err != nil {
		return err
	}

	return c.batch.Jobs(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deleteForegroundPolicy,
	})
}

// UpdateJob updates the given job.
func (c *Client) UpdateJob(job *batchv1.Job) (*batchv1.Job, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}
	return c.batch.Jobs(job.Namespace).Update(context.TODO(), job, metav1.UpdateOptions{})
}

// DeleteJobWithForce deletes the job with given namespace and name with force option
func (c *Client) DeleteJobWithForce(name, namespace string) error {
	if err := c.initClient(); err != nil {
		return err
	}
	gracePeriodSec := int64(0)
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy:  &deleteForegroundPolicy,
		GracePeriodSeconds: &gracePeriodSec,
	}
	return c.batch.Jobs(namespace).Delete(context.TODO(), name, deleteOptions)
}

// ValidateJob validates if the job with given namespace and name succeeds.
// It waits for timeout duration for job to succeed
func (c *Client) ValidateJob(name, namespace string, timeout time.Duration) error {
	t := func() (interface{}, bool, error) {
		job, err := c.GetJob(name, namespace)
		if err != nil {
			return nil, true, err
		}

		if job.Status.Failed > 0 {
			return nil, false, fmt.Errorf("job: [%s] %s has %d failed pod(s)", namespace, name, job.Status.Failed)
		}

		if job.Status.Active > 0 {
			return nil, true, fmt.Errorf("job: [%s] %s still has %d active pod(s)", namespace, name, job.Status.Active)
		}

		if job.Status.Succeeded == 0 {
			return nil, true, fmt.Errorf("job: [%s] %s no pod(s) that have succeeded", namespace, name)
		}

		return nil, false, nil
	}

	if _, err := task.DoRetryWithTimeout(t, timeout, 10*time.Second); err != nil {
		return err
	}

	return nil
}

// ListAllJobs returns the jobs from given namespace and options
func (c *Client) ListAllJobs(namespace string, filterOptions metav1.ListOptions) (*batchv1.JobList, error) {
	if err := c.initClient(); err != nil {
		return nil, err
	}

	return c.batch.Jobs(namespace).List(context.TODO(), filterOptions)
}
