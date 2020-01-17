package core

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInstance(t *testing.T) {
	Instance()

	require.NotNil(t, instance, "instance should be initialized")
}
