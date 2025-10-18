package main

import (
	"fmt"
	"os"

	"github.com/ubmagh/taq/parser"
	"github.com/ubmagh/taq/search"
)

func main() {

	inventory_hosts, err := parser.ParseInventoryFile(nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed inventory:\n%+v\n\n", inventory_hosts)

	search.RunSearcher(inventory_hosts)

	os.Exit(0)
}
