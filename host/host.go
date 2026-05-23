package host

import (
	"fmt"
	"strings"
)

type Host struct {
	Name             string            `yaml:"name"`
	Address          string            `yaml:"address"`
	User             string            `yaml:"user,omitempty"`
	KeyPath          string            `yaml:"key_path,omitempty"`
	Labels           map[string]string `yaml:"labels,omitempty"`
	Port             string            `yaml:"port,omitempty"`
	searchable string
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
	h.searchable = strings.ToLower(sb.String())
}

func (h Host) Searchable() string { return h.searchable }

func (h Host) HostListDisplay() string {
	return fmt.Sprintf("%s (%s @ %s)", h.Name, h.User, h.Address)
}

type Group struct {
	Labels map[string]string `yaml:"labels,omitempty"`
	Hosts  []Host            `yaml:"hosts"`
}

type Inventory struct {
	Groups map[string]Group `yaml:"groups,omitempty"`
	Hosts  []Host           `yaml:"hosts,omitempty"`
}
