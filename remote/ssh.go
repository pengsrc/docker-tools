package remote

import (
	"fmt"
	"net"
	"os"
	"path"
	"time"

	"github.com/pengsrc/go-shared/utils"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
	"golang.org/x/crypto/ssh/terminal"
)

// NewSSHSession creates an SSH session.
func NewSSHSession(host string, port int, username string) (session *ssh.Session, err error) {
	// Connect to ssh-agent.
	connection, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return
	}

	// Get signers from ssh-agent.
	signers, err := agent.NewClient(connection).Signers()
	if err != nil {
		return
	}

	// Setup callback for known hosts.
	hostKeyCallback, err := knownhosts.New(path.Join(utils.GetHome(), ".ssh", "known_hosts"))
	if err != nil {
		return
	}

	// Connect.
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signers...)},
		HostKeyCallback: hostKeyCallback,
		Timeout:         15 * time.Second,
	}
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), config)
	if err != nil {
		return
	}

	// Create new SSH session.
	session, err = client.NewSession()
	if err != nil {
		return
	}

	// Setup input/output.
	session.Stdin = os.Stdin
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	// Don't allocate terminal when not in a terminal.
	if !terminal.IsTerminal(int(os.Stdin.Fd())) {
		return
	}

	// Set up terminal modes.
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	// Request pseudo terminal.
	err = session.RequestPty("xterm-256color", 24, 80, modes)
	if err != nil {
		return
	}

	// Set terminal window.
	go func() {
		var termWidth, termHeight int
		for {
			time.Sleep(100 * time.Microsecond)

			w, h, err := terminal.GetSize(int(os.Stdin.Fd()))
			if err != nil {
				fmt.Fprint(session.Stderr, "Failed to get current terminal size\n")
			}

			if termHeight != h || termWidth != w {
				termHeight, termWidth = h, w
				err = session.WindowChange(termHeight, termWidth)
				if err != nil {
					fmt.Fprint(session.Stderr, "Failed to change remote terminal size\n")
				}
			}
		}
	}()

	return
}
