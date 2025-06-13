package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertError(t *testing.T, expectedError, actualError error) {
	t.Helper()
	if expectedError != nil {
		require.Error(t, actualError)
		require.ErrorIs(t, actualError, expectedError)
	} else {
		assert.NoError(t, actualError)
	}
}
