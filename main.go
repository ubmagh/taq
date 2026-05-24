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

const repoURL = "https://github.com/ubmagh/taq"

func printHelp() {
	fmt.Print(`
	taq - fast SSH search and connect CLI

	Usage:
	taq               # launch interactive search
	taq --help,    -h # show this help message
	taq --version, -v # show version
	taq --validate    # parse inventories and report host count, then exit
	taq --debug,   -d # enable verbose output (combine with --validate or normal run)

	Environment Variables:
	TAQ_DEFAULT_USER         : Specifies default SSH username [$USER]
	TAQ_DEFAULT_SSH_KEY_PATH : Default SSH key path []
	TAQ_ANSIBLE_INVS         : Semicolon-separated list of Ansible inventory root directories (e.g. inventories/) []
	TAQ_INVENTORY_PATH       : Inventory file path ["~/.config/taq/inventory.yaml"]
	TAQ_DISPLAY_MODE         : List display mode: "detailed" (default) or "compact"
	TAQ_SSH_TIMEOUT          : SSH connect timeout in seconds (e.g. 5) []
	TAQ_DEBUG                : Set to any value to enable debug/verbose output []
`)
}

func main() {
	var debugMode, validateMode bool

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--version", "-v":
			fmt.Printf("taq %s - %s\n", version, repoURL)
			return
		case "--help", "-h":
			printHelp()
			return
		case "--debug", "-d":
			debugMode = true
		case "--validate":
			validateMode = true
		default:
			fmt.Fprintf(os.Stderr, "unknown flag: %s\n", arg)
			os.Exit(1)
		}
	}

	if debugMode {
		os.Setenv("TAQ_DEBUG", "1")
	}

	inventoryHosts, err := parser.Parse()
	if err != nil {
		ui.Error("%v", err)
		os.Exit(1)
	}

	if validateMode {
		if len(inventoryHosts) == 0 {
			ui.Warn("no hosts found — create %s or set %s", "$HOME/.config/taq/inventory.yaml", "TAQ_ANSIBLE_INVS")
		} else {
			ui.Info("inventory OK — %d host(s) loaded", len(inventoryHosts))
		}
		return
	}

	if len(inventoryHosts) == 0 {
		ui.Warn("no hosts found — create %s or set %s", "$HOME/.config/taq/inventory.yaml", "TAQ_ANSIBLE_INVS")
		os.Exit(0)
	}

	if h, ok := search.RunSearcher(inventoryHosts); ok {
		ssh.OpenSSHSession(h)
	}
}
