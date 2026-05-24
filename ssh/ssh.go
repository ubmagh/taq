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

// buildBaseArgs returns the SSH args shared by all connection types:
// identity file, port, and connect timeout. The caller appends
// mode-specific flags and the destination.
func buildBaseArgs(h host.Host) []string {
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
	if timeout := config.GetSSHTimeout(); timeout != "" {
		args = append(args, "-o", "ConnectTimeout="+timeout)
	}
	return args
}

// clearScreen runs the clear command before handing the terminal to SSH.
func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}

// handleSSHError silently ignores a Ctrl+C exit (code 130) and prints
// an error message for everything else.
func handleSSHError(err error, action string) {
	if err == nil {
		return
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 130 {
		return
	}
	ui.Error("%s: %v", action, err)
}

// parseForwardRule converts the friendly shorthand "8080->3000" or "8080→3000"
// into the SSH spec "8080:localhost:3000", assuming localhost on the remote side.
// If the input is already in full SSH format it is passed through unchanged.
func parseForwardRule(rule string) string {
	normalized := strings.ReplaceAll(rule, "→", "->")
	parts := strings.SplitN(normalized, "->", 2)
	if len(parts) == 2 {
		local := strings.TrimSpace(parts[0])
		remote := strings.TrimSpace(parts[1])
		return fmt.Sprintf("%s:localhost:%s", local, remote)
	}
	return rule
}

// OpenPortForwardSession opens a port-forwarding-only SSH session (-N).
// flag must be "-L" (local) or "-R" (remote/reverse).
// rules use the shorthand format "localPort->remotePort".
func OpenPortForwardSession(h host.Host, flag string, rules []string) {
	args := buildBaseArgs(h)
	for _, rule := range rules {
		args = append(args, flag, parseForwardRule(rule))
	}
	args = append(args, "-N", fmt.Sprintf("%s@%s", h.User, h.Address))

	kindStr := "Local"
	if flag == "-R" {
		kindStr = "Remote"
	}

	clearScreen()
	fmt.Printf("Port Forwarding [%s] — %s (%s@%s)\n", kindStr, h.Name, h.User, h.Address)
	for _, r := range rules {
		fmt.Printf("  → %s\n", r)
	}
	fmt.Println("\nCtrl+C to stop.")

	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	handleSSHError(cmd.Run(), "port forwarding failed")
}

func OpenSSHSession(h host.Host) {
	if h.Address == "" {
		ui.Warn("no address configured for host %q", h.Name)
		return
	}
	args := buildBaseArgs(h)
	args = append(args, fmt.Sprintf("%s@%s", h.User, h.Address))

	clearScreen()
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	handleSSHError(cmd.Run(), "SSH connection failed")
}
