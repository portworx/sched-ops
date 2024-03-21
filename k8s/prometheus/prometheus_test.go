package prometheus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestInstance tests the Instance function
func TestInstance(t *testing.T) {
	Instance()

	require.NotNil(t, instance, "instance should be initialized")
}
