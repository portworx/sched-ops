package kubevirtdynamic

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstance(t *testing.T) {
	Instance()

	require.NotNil(t, instance, "instance should be initialized")
}

func TestGetVMI(t *testing.T) {
	// Populate the variables below from a live cluster to test this manually.
	testKubeconfig := "TBD"
	testVMINamespace := "TBD"
	testVMIName := "TBD"

	if testKubeconfig == "TBD" {
		t.Skip("Populate the test variables to run this test manually.")
	}
	os.Setenv("KUBECONFIG", testKubeconfig)

	Instance()

	require.NotNil(t, instance, "instance should be initialized")
	vmi, err := instance.GetVirtualMachineInstance(testVMINamespace, testVMIName)
	if err != nil {
		t.Logf("Failed to get VMI: %v", err)
	}
	t.Logf("VMI: %v", vmi)
}
