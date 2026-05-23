package ssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ubmagh/taq/config"
	"github.com/ubmagh/taq/host"
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
		fmt.Println("⚠️  No address found for host.")
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
		fmt.Printf("❌ SSH connection failed: %v\n", err)
	}
}
