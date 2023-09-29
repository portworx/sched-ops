package kubevirtdynamic

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetMigration(t *testing.T) {
	// Populate the variables below from a live cluster to test this manually.
	testKubeconfig := "TBD"
	testMigrationNamespace := "TBD"
	testMigrationName := "TBD"

	if testKubeconfig == "TBD" {
		t.Skip("Populate the test variables to run this test manually.")
	}
	os.Setenv("KUBECONFIG", testKubeconfig)

	Instance()

	require.NotNil(t, instance, "instance should be initialized")
	migration, err := instance.GetVirtualMachineInstanceMigration(
		context.TODO(), testMigrationNamespace, testMigrationName)
	if err != nil {
		t.Logf("Failed to get migration: %v", err)
	}
	t.Logf("Migration: %v", migration)
}

func TestCreateMigration(t *testing.T) {
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
	migration, err := instance.CreateVirtualMachineInstanceMigration(context.TODO(), testVMINamespace, testVMIName)
	if err != nil {
		t.Logf("Failed to create migration: %v", err)
	}
	t.Logf("Migration: %v", migration)
}
