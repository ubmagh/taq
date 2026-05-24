package parser

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/ubmagh/taq/config"
	"github.com/ubmagh/taq/host"
	"github.com/ubmagh/taq/ui"
)

type ansibleGroup struct {
	Hosts    map[string]ansibleVars  `yaml:"hosts"`
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

// non-inventory Ansible dirs that never contain inventory files
var ansibleSkipDirs = map[string]bool{
	"group_vars": true, "host_vars": true, "roles": true,
	"tasks": true, "handlers": true, "templates": true,
	"files": true, "vars": true, "defaults": true, "meta": true,
}

func parseAnsibleDir(dir string) ([]host.Host, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var hosts []host.Host
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			if ansibleSkipDirs[name] {
				continue
			}
			subHosts, err := parseAnsibleDir(filepath.Join(dir, name))
			if err != nil {
				return nil, err
			}
			hosts = append(hosts, subHosts...)
			continue
		}
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
		// File is likely a playbook, vars file, or other non-inventory YAML — skip it.
		if config.IsDebugMode() {
			ui.Warn("skipping %s: not an inventory file (%v)", path, err)
		}
		return nil, nil
	}

	if len(inv) == 0 {
		if config.IsDebugMode() {
			ui.Warn("skipping %s: parsed as inventory but contains no groups", path)
		}
		return nil, nil
	}

	if config.IsDebugMode() {
		ui.Info("loading ansible inventory file: %s", path)
	}

	var hosts []host.Host
	for groupName, group := range inv {
		hosts = append(hosts, flattenAnsibleGroup(groupName, group, nil, nil)...)
	}
	return hosts, nil
}

func flattenAnsibleGroup(name string, group ansibleGroup, inheritedVars ansibleVars, parentGroups []string) []host.Host {
	// Lazy merge: only allocate a new map when this group actually has its own vars
	// to override; otherwise share the inherited reference (it is never mutated).
	var groupVars ansibleVars
	if len(group.Vars) > 0 {
		groupVars = make(ansibleVars, len(inheritedVars)+len(group.Vars))
		maps.Copy(groupVars, inheritedVars)
		maps.Copy(groupVars, group.Vars)
	} else {
		groupVars = inheritedVars
	}

	// copy parent groups and append current (exclude "all")
	groups := make([]string, len(parentGroups))
	copy(groups, parentGroups)
	if name != "all" {
		groups = append(groups, name)
	}

	var hosts []host.Host

	for hostname, hostVars := range group.Hosts {
		// Lazy merge: only allocate when the host has overrides.
		var effective ansibleVars
		if len(hostVars) > 0 {
			effective = make(ansibleVars, len(groupVars)+len(hostVars))
			maps.Copy(effective, groupVars)
			maps.Copy(effective, hostVars)
		} else {
			effective = groupVars
		}

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
