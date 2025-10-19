package main

import (
	"fmt"
	"os"

	"github.com/ubmagh/taq/parser"
	"github.com/ubmagh/taq/search"
)

func main() {
	const Version = "v1.0.0"
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
