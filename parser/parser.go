package parser

import (
	"fmt"
	"maps"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/ubmagh/taq/config"
	"github.com/ubmagh/taq/host"
)

func ParseInventoryFile() ([]host.Host, error) {
	data, err := os.ReadFile(config.GetDefaultInventoryPath())
	if err != nil {
		return nil, fmt.Errorf("[Err] failed to read inventory file: %w", err)
	}

	var inv host.Inventory
	if err := yaml.Unmarshal(data, &inv); err != nil {
		return nil, err
	}

	var flattenedHosts []host.Host
	defaultUser := config.GetDefaultUser()

	for _, h := range inv.Hosts {
		flattenedHosts = append(flattenedHosts, h)
	}

	for gk, g := range inv.Groups {
		for _, h := range g.Hosts {
			mergedLabels := make(map[string]string)
			maps.Copy(mergedLabels, g.Labels)
			maps.Copy(mergedLabels, h.Labels)
			mergedLabels["groupName"] = gk
			h.Labels = mergedLabels
			flattenedHosts = append(flattenedHosts, h)
		}
	}

	for i := range flattenedHosts {
		if flattenedHosts[i].User == "" {
			flattenedHosts[i].User = defaultUser
		}
		flattenedHosts[i].BuildSearchable()
	}

	return flattenedHosts, nil
}
