package build

import (
	"os/exec"
)

// CommandExists checks whether a command exists.
func CommandExists(command string) (isExists bool) {
	cmd := exec.Command("which", command)
	out, err := cmd.Output()
	if err != nil {
		return
	}

	if len(out) > 0 {
		isExists = true
	}
	return
}
