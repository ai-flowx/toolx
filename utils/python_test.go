package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPython(t *testing.T) {
	_, err := GetPythonVersion()
	assert.Equal(t, nil, err)
}
