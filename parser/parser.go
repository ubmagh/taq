package parser

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/ubmagh/taq/config"
	"github.com/ubmagh/taq/types"
)

func ParseInventoryFile(out interface{}) ([]types.Host, error) {
	data, err := os.ReadFile(config.GetDefaultInventoryPath())
	if err != nil {
		return nil, fmt.Errorf("[Err] failed to read inventory file: %w", err)
	}

	var inv types.Inventory
	if err := yaml.Unmarshal(data, &inv); err != nil {
		return nil, err
	}

	var flattened_hosts []types.Host
	currentUser := os.Getenv("USER")

	// top level hosts
	for _, h := range inv.Hosts {
		flattened_hosts = append(flattened_hosts, h)
	}

	// hosts in groups
	for gk, g := range inv.Groups {
		for _, h := range g.Hosts {
			mergedLabels := make(map[string]string)
			for k, v := range g.Labels {
				mergedLabels[k] = v
			}
			for k, v := range h.Labels {
				mergedLabels[k] = v
			}
			mergedLabels["groupName"] = gk
			h.Labels = mergedLabels
			flattened_hosts = append(flattened_hosts, h)
		}
	}

	for i := range flattened_hosts {
		if flattened_hosts[i].User == "" {
			flattened_hosts[i].User = currentUser
		}
		flattened_hosts[i].BuildSearchable()
	}

	return flattened_hosts, nil
}
