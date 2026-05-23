package types

import (
	"fmt"
	"strings"

	"github.com/ubmagh/taq/config"
)

type Host struct {
	Name             string            `yaml:"name"`
	Address          string            `yaml:"address"`
	User             string            `yaml:"user,omitempty"`
	KeyPath          string            `yaml:"key_path,omitempty"`
	Labels           map[string]string `yaml:"labels,omitempty"`
	Port             string            `yaml:"port,omitempty"`
	SearchableString string
}

func (h *Host) BuildSearchable() {
	var sb strings.Builder
	sb.WriteString(h.Name)
	sb.WriteByte(' ')
	sb.WriteString(h.Address)
	sb.WriteByte(' ')
	sb.WriteString(h.User)
	sb.WriteByte(' ')
	for _, v := range h.Labels {
		sb.WriteString(v)
		sb.WriteByte(' ')
	}
	h.SearchableString = strings.ToLower(sb.String())
}

func (h Host) HostListDisplay() string {
	return fmt.Sprintf("%s (%s @ %s)\n", h.Name, h.User, h.Address)
}

func (h Host) GetSshCommand() []string {
	args := []string{}
	if len(h.KeyPath) == 0 {
		args = append(args, fmt.Sprintf("-i \"%s\"", h.KeyPath))
	} else {
		defaultKey := config.GetDefaultSshKeyPath()
		if len(defaultKey) > 0 {
			args = append(args, fmt.Sprintf("-i \"%s\"", defaultKey))
		}
	}
	if len(h.Port) > 0 {
		args = append(args, fmt.Sprintf("-p %s", strings.TrimSpace(h.Port)))
	}
	args = append(args, fmt.Sprintf("%s@%s", h.User, h.Address))

	return args
}

type Group struct {
	Labels map[string]string `yaml:"labels,omitempty"`
	Hosts  []Host            `yaml:"hosts"`
}

type Inventory struct {
	Groups map[string]Group `yaml:"groups,omitempty"`
	Hosts  []Host           `yaml:"hosts,omitempty"`
}
