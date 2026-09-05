package code

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGetPathSize calls with testdata path, checking
// for a valid return value.
func TestGetPathSize(t *testing.T) {
	result, err := GetPathSize("testdata", false, false, false)

	if err == nil {
		require.Equal(t, result, "5B")
	}
}

// TestGetPathSize calls with testdata path and recursive parameter, checking
// for a valid return value.
func TestGetPathSizeRecursive(t *testing.T) {
	result, err := GetPathSize("testdata", true, false, false)

	if err == nil {
		require.Equal(t, result, "10B")
	}
}

// TestGetPathSize calls with testdata path, recursive parameter and include hidden files, checking
// for a valid return value.
func TestGetPathSizeIncludeHidden(t *testing.T) {
	result, err := GetPathSize("testdata", true, false, true)

	if err == nil {
		require.Equal(t, result, "20B")
	}
}

// TestGetPathSize calls with empty string path, checking
// for an error to accure.
func TestGetPathSizeEmptyPathString(t *testing.T) {
	_, err := GetPathSize("", true, false, true)

	require.Error(t, err)
}
