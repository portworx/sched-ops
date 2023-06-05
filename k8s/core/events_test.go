package core

import (
	"testing"
	"time"

	"github.com/portworx/sched-ops/k8s/testutil"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRecordEvent(t *testing.T) {
	restCfg := testutil.SetUpTestCluster(t, "cluster1")

	testClient, err := NewForConfig(restCfg)
	require.NoError(t, err)

	SetInstance(testClient)

	testNamespace := "recordeventns1"
	recreateNamespace(t, testNamespace)

	testPod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod1",
			Namespace: testNamespace,
			UID:       "uid1",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "container1",
					Image: "image1",
				},
			},
		},
	}

	// Create a new Deployment object
	var replicas int32 = 1
	testDep := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dep1",
			Namespace: testNamespace,
			UID:       "uid2",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "app1",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "app1",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "container1",
							Image: "image1",
						},
					},
				},
			},
		},
	}
	// new interface
	Instance().RecordEventf("controller1", testPod, testDep, "Warning", "schedulingFailed", "scheduling",
		"%d/%d nodes not available for scheduling", 5, 5)

	// backwards-compatible interface that creates new-style event
	Instance().RecordEvent(corev1.EventSource{Component: "controller2", Host: "host2"}, testPod,
		"Normal", "scheduled", "pod scheduled on node n1")

	// backwards-compatible interface that creates an old-style event
	Instance().RecordEventLegacy(corev1.EventSource{Component: "controller3", Host: "host3"}, testDep,
		"Normal", "replicasReady", "all replicas are ready")

	// event creation is async in client-go. Need to wait a bit.
	var events *corev1.EventList
	waitUntil := time.Now().Add(10 * time.Second)
	for time.Now().Before(waitUntil) {
		events, err = Instance().ListEvents(testNamespace, metav1.ListOptions{})
		require.NoError(t, err)
		if len(events.Items) == 3 {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	require.Len(t, events.Items, 3)
	found := 0
	for _, event := range events.Items {
		switch event.ReportingController {
		case "controller1":
			require.Equal(t, event.Action, "scheduling")

			require.Equal(t, "Pod", event.InvolvedObject.Kind)
			require.Equal(t, testNamespace, event.InvolvedObject.Namespace)
			require.Equal(t, "pod1", event.InvolvedObject.Name)
			require.Equal(t, "uid1", string(event.InvolvedObject.UID))

			require.NotNil(t, event.Related)
			require.Equal(t, "Deployment", event.Related.Kind)
			require.Equal(t, testNamespace, event.Related.Namespace)
			require.Equal(t, "dep1", event.Related.Name)
			require.Equal(t, "uid2", string(event.Related.UID))

			require.Equal(t, "schedulingFailed", event.Reason)
			require.Equal(t, "5/5 nodes not available for scheduling", event.Message)
			require.Equal(t, "Warning", event.Type)
			require.Equal(t, "scheduling", event.Action)
			found++
		case "controller2":
			require.Equal(t, event.Action, "Unspecified")

			require.Equal(t, "Pod", event.InvolvedObject.Kind)
			require.Equal(t, testNamespace, event.InvolvedObject.Namespace)
			require.Equal(t, "pod1", event.InvolvedObject.Name)
			require.Equal(t, "uid1", string(event.InvolvedObject.UID))

			require.Nil(t, event.Related)
			require.Equal(t, "scheduled", event.Reason)
			require.Equal(t, "pod scheduled on node n1", event.Message)
			require.Equal(t, "Normal", event.Type)
			found++
		case "":
			require.Equal(t, "Deployment", event.InvolvedObject.Kind)
			require.Equal(t, testNamespace, event.InvolvedObject.Namespace)
			require.Equal(t, "dep1", event.InvolvedObject.Name)
			require.Equal(t, "uid2", string(event.InvolvedObject.UID))
			require.Equal(t, "replicasReady", event.Reason)
			require.Equal(t, "all replicas are ready", event.Message)
			require.Equal(t, "controller3", event.Source.Component)
			require.Equal(t, "host3", event.Source.Host)
			require.Equal(t, 1, int(event.Count))
			require.Equal(t, "Normal", event.Type)
			require.Equal(t, "", event.Action)
			found++
		}
	}
	require.Equal(t, 3, found)
}

// delete and recreate the specified namespace
func recreateNamespace(t *testing.T, name string) {
	err := Instance().DeleteNamespace(name)
	require.True(t, err == nil || errors.IsNotFound(err))
	for {
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: name}}
		_, err = Instance().CreateNamespace(ns)
		if errors.IsAlreadyExists(err) {
			t.Logf("waiting for the namespace %s to go away", name)
			// wait for the deleted namespace to get GC'ed
			time.Sleep(250 * time.Millisecond)
			continue
		}
		break
	}
	require.Nil(t, err)
}
