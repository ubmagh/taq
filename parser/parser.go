package parser

import (
	"errors"
	"fmt"
	"maps"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/ubmagh/taq/config"
	"github.com/ubmagh/taq/host"
	"github.com/ubmagh/taq/ui"
)

func Parse() ([]host.Host, error) {
	var all []host.Host

	for _, path := range config.GetInventoryPaths() {
		taqHosts, err := parseTaqFile(path)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("inventory %q: %w", path, err)
		}
		if config.IsDebugMode() && len(taqHosts) > 0 {
			ui.Info("loaded %d host(s) from %s", len(taqHosts), path)
		}
		all = append(all, taqHosts...)
	}

	for _, dir := range config.GetAnsibleInventories() {
		ansibleHosts, err := parseAnsibleDir(dir)
		if err != nil {
			return nil, fmt.Errorf("ansible inventory %q: %w", dir, err)
		}
		all = append(all, ansibleHosts...)
	}

	return applyDefaults(all), nil
}

func parseTaqFile(path string) ([]host.Host, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var inv host.Inventory
	if err := yaml.Unmarshal(data, &inv); err != nil {
		return nil, err
	}

	var hosts []host.Host
	hosts = append(hosts, inv.Hosts...)

	for gk, g := range inv.Groups {
		for _, h := range g.Hosts {
			mergedLabels := make(map[string]string)
			maps.Copy(mergedLabels, g.Labels)
			maps.Copy(mergedLabels, h.Labels)
			mergedLabels["groupName"] = gk
			h.Labels = mergedLabels
			hosts = append(hosts, h)
		}
	}

	return hosts, nil
}

func applyDefaults(hosts []host.Host) []host.Host {
	defaultUser := config.GetDefaultUser()
	seen := make(map[string]bool)
	result := make([]host.Host, 0, len(hosts))

	for _, h := range hosts {
		key := h.Name + "|" + h.Address
		if seen[key] {
			if config.IsDebugMode() {
				ui.Warn("duplicate host skipped: %q (%s)", h.Name, h.Address)
			}
			continue
		}
		seen[key] = true
		if h.User == "" {
			h.User = defaultUser
		}
		h.BuildSearchable()
		result = append(result, h)
	}

	return result
}
