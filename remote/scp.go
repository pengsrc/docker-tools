package remote

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Upload uploads a local path to remote server.
func Upload(
	host string, port int, username string, localPath string, remotePath string,
) (err error) {
	cmd := exec.Command(
		"scp", "-P", strconv.Itoa(port),
		localPath, fmt.Sprintf("%s@%s:%s", username, host, remotePath),
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Executing: %s %v\n", cmd.Path, cmd.Args)
	err = cmd.Run()
	if err != nil {
		return
	}

	return
}
