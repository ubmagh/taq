package types

type Host struct {
	Name    string            `yaml:"name"`
	Address string            `yaml:"address"`
	User    string            `yaml:"user,omitempty"`
	KeyPath string            `yaml:"key_path,omitempty"`
	Labels  map[string]string `yaml:"labels,omitempty"`
}

type Subgroup struct {
	Labels map[string]string `yaml:"labels,omitempty"`
	Hosts  []Host            `yaml:"hosts"`
	Groups map[string]string `yaml:"groups,omitempty"`
}

type Group struct {
	Labels map[string]string `yaml:"labels,omitempty"`
	Hosts  []Host            `yaml:"hosts"`
}

type Inventory struct {
	Groups map[string]Group `yaml:"groups,omitempty"`
	Hosts  []Host           `yaml:"hosts,omitempty"`
}
