package output

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/goccy/go-yaml"
	"github.com/ubmagh/taq/host"
)

const (
	FormatTable = "table"
	FormatJSON  = "json"
	FormatYAML  = "yaml"
	FormatPlain = "plain"
)

// Print writes hosts to stdout in the requested format.
// An empty format defaults to table.
func Print(hosts []host.Host, format string) error {
	switch format {
	case FormatTable, "":
		return printTable(hosts)
	case FormatJSON:
		return printJSON(hosts)
	case FormatYAML:
		return printYAML(hosts)
	case FormatPlain:
		return printPlain(hosts)
	default:
		return fmt.Errorf("unknown format %q — valid: table, json, yaml, plain", format)
	}
}

// groups returns a human-readable group string for a host, checking both
// the Ansible "groups" label and the taq-inventory "groupName" label.
func groups(h host.Host) string {
	if g := h.Labels["groups"]; g != "" {
		return g
	}
	return h.Labels["groupName"]
}

func printTable(hosts []host.Host) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tADDRESS\tUSER\tPORT\tGROUPS")
	fmt.Fprintln(w, "----\t-------\t----\t----\t------")
	for _, h := range hosts {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", h.Name, h.Address, h.User, h.Port, groups(h))
	}
	return w.Flush()
}

func printJSON(hosts []host.Host) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(hosts)
}

func printYAML(hosts []host.Host) error {
	data, err := yaml.Marshal(hosts)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	return err
}

func printPlain(hosts []host.Host) error {
	for _, h := range hosts {
		fmt.Printf("%s %s\n", h.Name, h.Address)
	}
	return nil
}
