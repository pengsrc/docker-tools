package build

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ImageTagForGitRepo returns valid git tag or commit hash for given git repo.
func ImageTagForGitRepo(repoDir string, tag string) (tagOrCommitHash string, err error) {
	git := "git"
	if !CommandExists(git) {
		err = fmt.Errorf("command not exists: %s", git)
		return
	}

	if !GitRepoExists(repoDir) {
		err = fmt.Errorf("git repo not exists: %s", repoDir)
		return
	}

	if tag != "" && GitTagExists(repoDir, tag) {
		tagOrCommitHash = tag
		return
	}

	tagOrCommitHash = LatestGitCommitHash(repoDir, true)
	if tagOrCommitHash == "" {
		err = fmt.Errorf("failed to get latest commit hash")
		return
	}
	return
}

// GitRepoExists checks whether a git repo exists.
func GitRepoExists(repoDir string) (isExists bool) {
	cmd := exec.Command("git", "status")
	cmd.Dir = repoDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return
	}

	if cmd.ProcessState.Success() {
		isExists = true
	}
	return
}

// GitTagExists checks whether a git tag exists.
func GitTagExists(repoDir, tag string) (isExists bool) {
	cmd := exec.Command("git", "rev-parse", tag)
	cmd.Dir = repoDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return
	}

	if cmd.ProcessState.Success() {
		isExists = true
	}
	return
}

// LatestGitCommitHash gets the latest git commit hash.
func LatestGitCommitHash(repoDir string, isShort bool) (commitHash string) {
	var cmd *exec.Cmd
	if isShort {
		cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
	} else {
		cmd = exec.Command("git", "rev-parse", "HEAD")
	}
	cmd.Dir = repoDir
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return
	}
	if !cmd.ProcessState.Success() {
		return
	}

	hash := strings.Trim(string(out), "\n")
	if !strings.Contains(hash, " ") {
		if isShort {
			if len(hash) == 7 {
				commitHash = hash
			}
		} else {
			if len(hash) == 40 {
				commitHash = hash
			}
		}
	}
	return
}
