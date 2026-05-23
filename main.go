package main

import (
	"fmt"
	"os"

	"github.com/ubmagh/taq/parser"
	"github.com/ubmagh/taq/search"
	"github.com/ubmagh/taq/ssh"
)

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
	TAQ_ANSIBLE_INVS         : List of ansible projects inventories, (;) separated.  []
	TAQ_INVENTORY_PATH       : Inventory file path ["~/.config/taq/inventory.yaml"]
`)
}

func main() {
	const Version = "v1.0.0"

	if len(os.Args) > 1 {
		arg := os.Args[1]

		switch arg {
		case "--version", "-v":
			fmt.Println("taq", Version)
			return
		case "--help", "-h":
			printHelp()
			return
		}
	}

	inventory_hosts, err := parser.ParseInventoryFile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if host, ok := search.RunSearcher(inventory_hosts); ok {
		ssh.OpenSSHSession(host)
	}

	os.Exit(0)
}
