package kubevirtdynamic

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetVM(t *testing.T) {
	// Populate the variables below from a live cluster to test this manually.
	testKubeconfig := "TBD"
	testVMNamespace := "TBD"
	testVMName := "TBD"

	if testKubeconfig == "TBD" {
		t.Skip("Populate the test variables to run this test manually.")
	}
	os.Setenv("KUBECONFIG", testKubeconfig)

	Instance()

	require.NotNil(t, instance, "instance should be initialized")
	vm, err := instance.GetVirtualMachine(context.TODO(), testVMNamespace, testVMName)
	if err != nil {
		t.Logf("Failed to get VM: %v", err)
		t.FailNow()
	}
	t.Logf("VM: %v", vm)
}

func TestListVMs(t *testing.T) {
	// Populate the variables below from a live cluster to test this manually.
	testKubeconfig := "TBD"
	testVMNamespace := "TBD"

	if testKubeconfig == "TBD" {
		t.Skip("Populate the test variables to run this test manually.")
	}
	os.Setenv("KUBECONFIG", testKubeconfig)

	Instance()

	require.NotNil(t, instance, "instance should be initialized")
	vms, err := instance.ListVirtualMachines(context.TODO(), testVMNamespace, metav1.ListOptions{})
	if err != nil {
		t.Logf("Failed to get VMs: %v", err)
		t.FailNow()
	}
	for i, vm := range vms {
		t.Logf("VM #%d: %v", i, vm)
	}
}
