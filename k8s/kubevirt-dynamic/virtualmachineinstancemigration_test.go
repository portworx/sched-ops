package kubevirtdynamic

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		testMigrationNamespace, testMigrationName)
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
	migration, err := instance.CreateVirtualMachineInstanceMigration(testVMINamespace, testVMIName)
	if err != nil {
		t.Logf("Failed to create migration: %v", err)
	}
	t.Logf("Migration: %v", migration)
}

func TestCreateMigrationWithParams(t *testing.T) {
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

	migration, err := instance.CreateVirtualMachineInstanceMigrationWithParams(
		testVMINamespace, testVMIName, "migration1", "",
		map[string]string{"anno1": "val1"}, map[string]string{"label2": "val2"})
	if err != nil {
		t.Logf("Failed to create migration: %v", err)
	}
	t.Logf("Migration: %v", migration)
}

func TestListMigrations(t *testing.T) {
	// Populate the variables below from a live cluster to test this manually.
	testKubeconfig := "TBD"
	testMigrationNamespace := "TBD"

	if testKubeconfig == "TBD" {
		t.Skip("Populate the test variables to run this test manually.")
	}
	os.Setenv("KUBECONFIG", testKubeconfig)

	Instance()

	require.NotNil(t, instance, "instance should be initialized")
	migrations, err := instance.ListVirtualMachineInstanceMigrations(testMigrationNamespace, metav1.ListOptions{})
	if err != nil {
		t.Logf("Failed to list migrations: %v", err)
		t.FailNow()
	}
	for i, migration := range migrations {
		t.Logf("Migration %d:", i)
		t.Logf("================")
		t.Logf("%v\n\n", migration)
	}
}
