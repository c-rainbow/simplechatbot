package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Failure of this test means a problem with build.
func TestEmpty(t *testing.T) {
	assert.True(t, true)
}

func TestClientCreation(t *testing.T) {
	testClient := NewTwitchClient("test_username", "test_oauth")
	assert.NotNil(t, testClient)
}
