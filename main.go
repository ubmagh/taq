package main

import (
	"fmt"
	"os"

	"github.com/ubmagh/taq/parser"
	"github.com/ubmagh/taq/search"
)

func printHelp() {
	fmt.Println(`
	taq - fast SSH search and connect CLI
	
	Usage:
	taq            # launch interactive search
	taq --help,-h  # show this help message
	taq --version,-v  # show version

	Config:
	TAQ_INVENTORY_PATH : environment variable to specify inventory file path, default : "~/.config/taq.inventory.yaml"

	Features:
	• Search hosts by name, address, user, and labels
	• Interactive fuzzy search with up/down arrows
	• Launch SSH session directly from the list
	• Uses inventory from YAML file`)
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

	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("taq", Version)
		return
	}

	inventory_hosts, err := parser.ParseInventoryFile(nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	search.RunSearcher(inventory_hosts)

	os.Exit(0)
}
