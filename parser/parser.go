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

	// top level hosts
	for _, h := range inv.Hosts {
		flattened_hosts = append(flattened_hosts, h)
	}

	// // hosts in groups
	// for _, g := range inv.Groups {
	// 	// hosts directly under group
	// 	for _, h := range g.Hosts {
	// 		merged := make(map[string]string)
	// 		for k, v := range g.Labels {
	// 			merged[k] = v
	// 		}
	// 		for k, v := range h.Labels {
	// 			merged[k] = v
	// 		}
	// 		h.Labels = merged
	// 		flattened_hosts = append(flattened_hosts, h)
	// 	}

	// 	// hosts in subgroups
	// 	for _, sg := range g.Groups {
	// 		for _, h := range sg.Hosts {
	// 			merged := make(map[string]string)
	// 			for k, v := range g.Labels {
	// 				merged[k] = v
	// 			}
	// 			for k, v := range sg.Labels {
	// 				merged[k] = v
	// 			}
	// 			for k, v := range h.Labels {
	// 				merged[k] = v
	// 			}
	// 			h.Labels = merged
	// 			flattened_hosts = append(flattened_hosts, h)
	// 		}
	// 	}
	// }

	// if out != nil {

	// }
	// return flattened_hosts, nil
	// for _, g := range inv.Groups {
	// 	for _, h := range g.Hosts {
	// 		merged := make(map[string]string)
	// 		for k, v := range g.Labels {
	// 			merged[k] = v
	// 		}
	// 		for k, v := range h.Labels {
	// 			merged[k] = v
	// 		}
	// 		h.Labels = merged
	// 		hosts = append(hosts, h)
	// 	}
	// }

	// return hosts, nil
	return nil, nil
}
