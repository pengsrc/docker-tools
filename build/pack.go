package build

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/pengsrc/docker-tools/constants"
)

// CreatePackage creates source package.
func CreatePackage(
	sourceDir string, includes []string, excludes []string,
	useGitArchive bool, tagOrCommitHash string,
) (packagePath string, err error) {
	f, err := ioutil.TempFile(os.TempDir(), fmt.Sprintf("%d-", time.Now().Unix()))
	if err != nil {
		return
	}
	defer f.Close()

	// Check tar or gtar.
	tar := "tar"
	if CommandExists(constants.GNUPrefix + tar) {
		tar = constants.GNUPrefix + tar
	} else {
		if !CommandExists(tar) {
			err = fmt.Errorf(`command "%s" not found`, tar)
			return
		}
	}

	// Check gzip.
	gzip := "gzip"
	if !CommandExists(gzip) {
		err = fmt.Errorf(`command "%s" not found`, gzip)
		return
	}

	// Commands to execute later.
	cmds := []*exec.Cmd{}

	// Create tarball.
	if useGitArchive {
		// Check git.
		git := "git"
		if !CommandExists(git) {
			err = fmt.Errorf(`command "%s" not found`, git)
			return
		}

		command := exec.Command(git, "archive", "--format", "tar", tagOrCommitHash)
		command.Stdout = f
		cmds = append(cmds, exec.Command("pwd"))
		cmds = append(cmds, command)
	} else {
		cmds = append(cmds, exec.Command(tar, "--transform", "s,^./,,g", "-cf", f.Name(), "."))
	}

	// Include directories.
	if len(includes) > 0 {
		args := []string{"--transform", "s,^./,,g", "-rf", f.Name()}
		args = append(args, includes...)
		cmds = append(cmds, exec.Command(tar, args...))
	}

	// Exclude directories.
	if len(excludes) > 0 {
		args := []string{"--delete", "-f", f.Name()}
		args = append(args, excludes...)
		cmds = append(cmds, exec.Command(tar, args...))
	}

	// Compress tarball.
	mv := "mv"
	if !CommandExists(mv) {
		err = fmt.Errorf(`command "%s" not found`, mv)
		return
	}

	cmds = append(cmds, exec.Command(gzip, "-9f", f.Name()))
	cmds = append(cmds, exec.Command(mv, fmt.Sprintf("%s.gz", f.Name()), f.Name()))

	// List tarball.
	cmds = append(cmds, exec.Command(tar, "-tf", f.Name()))

	// Execute commands.
	for _, cmd := range cmds {
		fmt.Printf("Executing: %s %v\n", cmd.Path, cmd.Args)
		cmd.Dir = sourceDir
		if cmd.Stdin == nil {
			cmd.Stdin = os.Stdin
		}
		if cmd.Stdout == nil {
			cmd.Stdout = os.Stdout
		}
		if cmd.Stderr == nil {
			cmd.Stderr = os.Stderr
		}
		err = cmd.Run()
		if err != nil {
			return
		}
	}

	packagePath = f.Name()
	return
}
