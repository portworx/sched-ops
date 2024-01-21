package kubevirtdynamic

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetDataVolume(t *testing.T) {
	// Populate the variables below from a live cluster to test this manually.
	testKubeconfig := "TBD"
	testDVNamespace := "TBD"
	testDVName := "TBD"

	if testKubeconfig == "TBD" {
		t.Skip("Populate the test variables to run this test manually.")
	}
	os.Setenv("KUBECONFIG", testKubeconfig)

	Instance()

	require.NotNil(t, instance, "instance should be initialized")
	vmi, err := instance.GetDataVolume(context.TODO(), testDVNamespace, testDVName)
	if err != nil {
		t.Logf("Failed to get data volume: %v", err)
		t.FailNow()
	}
	t.Logf("Data volume: %v", vmi)
}

func TestListDataVolumes(t *testing.T) {
	// Populate the variables below from a live cluster to test this manually.
	testKubeconfig := "TBD"
	testDVNamespace := "TBD"

	if testKubeconfig == "TBD" {
		t.Skip("Populate the test variables to run this test manually.")
	}
	os.Setenv("KUBECONFIG", testKubeconfig)

	Instance()

	require.NotNil(t, instance, "instance should be initialized")
	dvs, err := instance.ListDataVolumes(
		context.TODO(), testDVNamespace, metav1.ListOptions{})
	if err != nil {
		t.Logf("Failed to list data volume: %v", err)
		t.FailNow()
	}
	for i, dv := range dvs {
		t.Logf("Data volume %d:", i)
		t.Logf("================")
		t.Logf("%v\n\n", dv)
	}
}
