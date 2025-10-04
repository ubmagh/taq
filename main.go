package main

import (
	"fmt"
	"os"

	"github.com/ubmagh/taq/parser"
)

func main() {
	fmt.Print("Hey there.\n")

	inventory_hosts, err := parser.ParseInventoryFile(nil)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed inventory:\n%+v\n\n", inventory_hosts)

	os.Exit(0)
}
