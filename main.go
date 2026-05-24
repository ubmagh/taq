package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/ubmagh/taq/parser"
	"github.com/ubmagh/taq/search"
	"github.com/ubmagh/taq/ssh"
	"github.com/ubmagh/taq/ui"
)

var version = "main"

func init() {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
	}
}

func init() {
	if version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
	}
}

const repoURL = "https://github.com/ubmagh/taq"

func printHelp() {
	fmt.Print(`
	taq - fast SSH search and connect CLI

	Usage:
	taq               # launch interactive search
	taq --help,-h     # show this help message
	taq --version,-v  # show version

	Environment Variables:
	TAQ_DEFAULT_USER         : Specifies default SSH username [$USER]
	TAQ_DEFAULT_SSH_KEY_PATH : Default SSH key path. []
	TAQ_ANSIBLE_INVS         : Semicolon-separated list of Ansible inventory root directories (e.g. inventories/).  []
	TAQ_INVENTORY_PATH       : Inventory file path ["~/.config/taq/inventory.yaml"]
	TAQ_DISPLAY_MODE         : List display mode: "detailed" (default) or "compact"
`)
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Printf("taq %s - %s\n", version, repoURL)
			return
		case "--help", "-h":
			printHelp()
			return
		}
	}

	inventoryHosts, err := parser.Parse()
	if err != nil {
		ui.Error("%v", err)
		os.Exit(1)
	}

	if len(inventoryHosts) == 0 {
		ui.Warn("no hosts found — create %s or set %s", "$HOME/.config/taq/inventory.yaml", "TAQ_ANSIBLE_INVS")
		os.Exit(0)
	}

	if h, ok := search.RunSearcher(inventoryHosts); ok {
		ssh.OpenSSHSession(h)
	}
}
