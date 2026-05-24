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
	taq                          # launch interactive SSH search
	taq -l, --local-forward      # launch in local port-forward mode (-L)
	taq -r, --remote-forward     # launch in remote/reverse port-forward mode (-R)
	taq --validate               # parse inventories, report host count, then exit
	taq --debug,   -d            # enable verbose output (combine with any flag)
	taq --version, -v            # show version
	taq --help,    -h            # show this help message

	Environment Variables:
	TAQ_DEFAULT_USER         : Default SSH username [$USER]
	TAQ_DEFAULT_SSH_KEY_PATH : Default SSH key path []
	TAQ_ANSIBLE_INVS         : Semicolon-separated list of Ansible inventory root directories []
	TAQ_INVENTORY_PATH       : Inventory file path [$HOME/.config/taq/inventory.yaml]
	TAQ_DISPLAY_MODE         : List display mode: "detailed" (default) or "compact"
	TAQ_SSH_TIMEOUT          : SSH connect timeout in seconds (e.g. 5) []
	TAQ_DEBUG                : Set to any value to enable debug/verbose output []
`)
}

func main() {
	var debugMode, validateMode bool
	forwardMode := search.KindSSH

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
		case "-l", "--local-forward":
			forwardMode = search.KindLocalForward
		case "-r", "--remote-forward":
			forwardMode = search.KindRemoteForward
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

	if result, ok := search.RunSearcher(inventoryHosts, forwardMode); ok {
		switch result.Kind {
		case search.KindSSH:
			ssh.OpenSSHSession(result.Host)
		case search.KindLocalForward:
			ssh.OpenPortForwardSession(result.Host, "-L", result.Rules)
		case search.KindRemoteForward:
			ssh.OpenPortForwardSession(result.Host, "-R", result.Rules)
		}
	}
}
