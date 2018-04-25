package build

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitRepoExists(t *testing.T) {
	assert.False(t, GitRepoExists(os.TempDir()))
}
