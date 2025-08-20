package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReplaceExtension(t *testing.T) {
	result := ReplaceExtension("path/to/example.pdf", ".json")
	require.EqualValues(t, "path/to/example.json", result)
}
