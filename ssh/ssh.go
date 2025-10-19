package ssh

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/ubmagh/taq/types"
)

func OpenSSHSession(h types.Host) {
	if h.Address == "" {
		fmt.Println("⚠️  No address found for host.")
		return
	}

	cmd := exec.Command("ssh", h.GetSshCommand())
	// Attach to current terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run SSH directly (blocking until exit)
	if err := cmd.Run(); err != nil {
		fmt.Printf("❌ SSH connection failed: %v\n", err)
	}
}
