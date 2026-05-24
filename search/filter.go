package search

import (
	"strings"

	"github.com/sahilm/fuzzy"
	"github.com/ubmagh/taq/host"
)

// FilterHosts returns the subset of hosts matching query using fuzzy search
// plus an address substring fallback, preserving match order.
// If query is empty all hosts are returned unchanged.
// This is the same logic used by the TUI's search box.
func FilterHosts(hosts []host.Host, query string) []host.Host {
	query = strings.TrimSpace(query)
	if query == "" {
		return hosts
	}

	lq := strings.ToLower(query)

	searchables := make([]string, len(hosts))
	for i, h := range hosts {
		searchables[i] = h.Searchable()
	}

	matches := fuzzy.Find(lq, searchables)
	seen := make(map[int]bool, len(matches))
	filtered := make([]host.Host, 0, len(matches))
	for _, match := range matches {
		seen[match.Index] = true
		filtered = append(filtered, hosts[match.Index])
	}

	// fuzzy handles names/groups well but struggles with IPs — supplement
	// with a plain address substring match.
	for i, h := range hosts {
		if !seen[i] && strings.Contains(strings.ToLower(h.Address), lq) {
			filtered = append(filtered, h)
		}
	}

	return filtered
}
