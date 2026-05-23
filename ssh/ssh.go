package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ubmagh/taq/config"
	"github.com/ubmagh/taq/host"
	"github.com/ubmagh/taq/ui"
)

func sshArgs(h host.Host) []string {
	var args []string
	keyPath := h.KeyPath
	if keyPath == "" {
		keyPath = config.GetDefaultSshKeyPath()
	}
	if keyPath != "" {
		args = append(args, "-i", keyPath)
	}
	if h.Port != "" {
		args = append(args, "-p", strings.TrimSpace(h.Port))
	}
	args = append(args, fmt.Sprintf("%s@%s", h.User, h.Address))
	return args
}

func OpenSSHSession(h host.Host) {
	if h.Address == "" {
		ui.Warn("no address configured for host %q", h.Name)
		return
	}
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
	cmd = exec.Command("ssh", sshArgs(h)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
			return
		}
		ui.Error("SSH connection failed: %v", err)
	}
}
