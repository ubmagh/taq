package parser

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/ubmagh/taq/host"
)

type ansibleGroup struct {
	Hosts    map[string]ansibleVars `yaml:"hosts"`
	Children map[string]ansibleGroup `yaml:"children"`
	Vars     ansibleVars             `yaml:"vars"`
}

type ansibleVars map[string]any

func (v ansibleVars) str(key string) string {
	if v == nil {
		return ""
	}
	val, ok := v[key]
	if !ok {
		return ""
	}
	switch s := val.(type) {
	case string:
		return s
	case int:
		return strconv.Itoa(s)
	default:
		return fmt.Sprintf("%v", s)
	}
}

func parseAnsibleDir(dir string) ([]host.Host, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var hosts []host.Host
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") && name != "hosts" {
			continue
		}
		fileHosts, err := parseAnsibleFile(filepath.Join(dir, name))
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		hosts = append(hosts, fileHosts...)
	}
	return hosts, nil
}

func parseAnsibleFile(path string) ([]host.Host, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var inv map[string]ansibleGroup
	if err := yaml.Unmarshal(data, &inv); err != nil {
		return nil, err
	}

	var hosts []host.Host
	for groupName, group := range inv {
		hosts = append(hosts, flattenAnsibleGroup(groupName, group, nil, nil)...)
	}
	return hosts, nil
}

func flattenAnsibleGroup(name string, group ansibleGroup, inheritedVars ansibleVars, parentGroups []string) []host.Host {
	// merge: inherited < group vars (group overrides inherited)
	groupVars := make(ansibleVars)
	maps.Copy(groupVars, inheritedVars)
	maps.Copy(groupVars, group.Vars)

	// copy parent groups and append current (exclude "all")
	groups := make([]string, len(parentGroups))
	copy(groups, parentGroups)
	if name != "all" {
		groups = append(groups, name)
	}

	var hosts []host.Host

	for hostname, hostVars := range group.Hosts {
		// merge: group vars < host vars (host overrides group)
		effective := make(ansibleVars)
		maps.Copy(effective, groupVars)
		maps.Copy(effective, hostVars)

		address := effective.str("ansible_host")
		if address == "" {
			address = hostname
		}

		var labels map[string]string
		if len(groups) > 0 {
			labels = map[string]string{"groups": strings.Join(groups, " ")}
		}

		hosts = append(hosts, host.Host{
			Name:    hostname,
			Address: address,
			User:    effective.str("ansible_user"),
			Port:    effective.str("ansible_port"),
			KeyPath: effective.str("ansible_ssh_private_key_file"),
			Labels:  labels,
		})
	}

	for childName, childGroup := range group.Children {
		hosts = append(hosts, flattenAnsibleGroup(childName, childGroup, groupVars, groups)...)
	}

	return hosts
}
