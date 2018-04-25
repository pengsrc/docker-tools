package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseImageTag(t *testing.T) {
	var image, tag string

	image, tag = ParseImageTag("private/test-service:v0.0.1")
	assert.Equal(t, "private/test-service", image)
	assert.Equal(t, "v0.0.1", tag)

	image, tag = ParseImageTag("private/test-service:")
	assert.Equal(t, "", image)
	assert.Equal(t, "", tag)

	image, tag = ParseImageTag("private/test-service")
	assert.Equal(t, "private/test-service", image)
	assert.Equal(t, "", tag)
}

func TestCheckCommandExists(t *testing.T) {
	assert.True(t, CommandExists("ls"))
}
