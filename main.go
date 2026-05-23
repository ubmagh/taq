package main

import (
	"fmt"
	"os"

	"github.com/ubmagh/taq/parser"
	"github.com/ubmagh/taq/search"
	"github.com/ubmagh/taq/ssh"
)

var version = "dev"

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
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--version", "-v":
			fmt.Println("taq", version)
			return
		case "--help", "-h":
			printHelp()
			return
		}
	}

	inventoryHosts, err := parser.ParseInventoryFile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if h, ok := search.RunSearcher(inventoryHosts); ok {
		ssh.OpenSSHSession(h)
	}
}
