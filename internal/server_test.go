package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContinousIntegration(t *testing.T) {
	result := MockTest()
	assert.Equal(t, "ok", result)
}
