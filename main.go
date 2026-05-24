package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/ubmagh/taq/output"
	"github.com/ubmagh/taq/parser"
	"github.com/ubmagh/taq/search"
	"github.com/ubmagh/taq/ssh"
	"github.com/ubmagh/taq/ui"
)

var version = "dev"

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
	taq --list [query]           # list all hosts (or filter by query), then exit
	taq --list [query] -o fmt    # same, with output format: table (default), json, yaml, plain
	taq --validate               # parse inventories, report host count, then exit
	taq --debug,   -d            # enable verbose output (combine with any flag)
	taq --version, -v            # show version
	taq --help,    -h            # show this help message

	Environment Variables:
	TAQ_DEFAULT_USER         : Default SSH username [$USER]
	TAQ_DEFAULT_SSH_KEY_PATH : Default SSH key path []
	TAQ_INVENTORY_PATHS      : Semicolon-separated taq-inventory file paths [$HOME/.config/taq/inventory.yaml]
	TAQ_ANSIBLE_INVS         : Semicolon-separated list of Ansible inventory root directories []
	TAQ_DISPLAY_MODE         : List display mode: "detailed" (default) or "compact"
	TAQ_SSH_TIMEOUT          : SSH connect timeout in seconds (e.g. 5) []
	TAQ_DEBUG                : Set to any value to enable debug/verbose output []
`)
}

func main() {
	var debugMode, validateMode, listMode bool
	forwardMode := search.KindSSH
	var outputFormat, listQuery string

	// Index-based loop so we can consume the next arg as a flag value (e.g. -o json).
	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "--version" || arg == "-v":
			fmt.Printf("taq %s - %s\n", version, repoURL)
			return
		case arg == "--help" || arg == "-h":
			printHelp()
			return
		case arg == "--debug" || arg == "-d":
			debugMode = true
		case arg == "--validate":
			validateMode = true
		case arg == "--list":
			listMode = true
		case arg == "-l" || arg == "--local-forward":
			forwardMode = search.KindLocalForward
		case arg == "-r" || arg == "--remote-forward":
			forwardMode = search.KindRemoteForward
		case arg == "-o" || arg == "--output":
			if i+1 >= len(args) {
				fmt.Fprintln(os.Stderr, "flag -o/--output requires a value: table, json, yaml, plain")
				os.Exit(1)
			}
			i++
			outputFormat = args[i]
		case strings.HasPrefix(arg, "--output="):
			outputFormat = strings.TrimPrefix(arg, "--output=")
		case !strings.HasPrefix(arg, "-"):
			// Positional argument: treated as a search query for --list.
			listQuery = arg
		default:
			fmt.Fprintf(os.Stderr, "unknown flag: %s\n", arg)
			os.Exit(1)
		}
	}

	// Cross-flag validations.
	if outputFormat != "" && !listMode {
		fmt.Fprintln(os.Stderr, "flag -o/--output requires --list")
		os.Exit(1)
	}
	if listQuery != "" && !listMode {
		fmt.Fprintf(os.Stderr, "unexpected argument %q — did you mean: taq --list %q ?\n", listQuery, listQuery)
		os.Exit(1)
	}
	if listMode && forwardMode != search.KindSSH {
		fmt.Fprintln(os.Stderr, "--list cannot be combined with -l/--local-forward or -r/--remote-forward")
		os.Exit(1)
	}
	if listMode && validateMode {
		fmt.Fprintln(os.Stderr, "--list and --validate cannot be used together")
		os.Exit(1)
	}

	if debugMode {
		os.Setenv("TAQ_DEBUG", "1")
	}

	inventoryHosts, err := parser.Parse()
	if err != nil {
		ui.Error("%v", err)
		os.Exit(1)
	}

	// --validate
	if validateMode {
		if len(inventoryHosts) == 0 {
			ui.Warn("no hosts found — create %s or set %s", "$HOME/.config/taq/inventory.yaml", "TAQ_ANSIBLE_INVS")
		} else {
			ui.Info("inventory OK — %d host(s) loaded", len(inventoryHosts))
		}
		return
	}

	// --list [query] [-o format]
	if listMode {
		filtered := search.FilterHosts(inventoryHosts, listQuery)
		if len(filtered) == 0 {
			if listQuery != "" {
				ui.Warn("no hosts matched %q", listQuery)
				os.Exit(1)
			}
			ui.Warn("no hosts found — create %s or set %s", "$HOME/.config/taq/inventory.yaml", "TAQ_ANSIBLE_INVS")
			os.Exit(0)
		}
		if err := output.Print(filtered, outputFormat); err != nil {
			ui.Error("%v", err)
			os.Exit(1)
		}
		return
	}

	// Interactive mode.
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
