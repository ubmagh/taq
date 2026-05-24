package config

import (
	"os"
	"strings"
)

const (
	defaultInventoryPath = "$HOME/.config/taq/inventory.yaml"

	envInventoryPaths = "TAQ_INVENTORY_PATHS"
	envDefaultUser    = "TAQ_DEFAULT_USER"
	envAnsibleInvs    = "TAQ_ANSIBLE_INVS"
	envSshKeyPath     = "TAQ_DEFAULT_SSH_KEY_PATH"
	envDisplayMode    = "TAQ_DISPLAY_MODE"
	envDebug          = "TAQ_DEBUG"
	envSSHTimeout     = "TAQ_SSH_TIMEOUT"
)

func IsCompactMode() bool {
	return os.Getenv(envDisplayMode) == "compact"
}

// GetInventoryPaths returns all taq-inventory file paths to load.
// Set TAQ_INVENTORY_PATHS to a semicolon-separated list of paths (one or many).
// Falls back to the default path if the variable is not set.
func GetInventoryPaths() []string {
	if v := os.Getenv(envInventoryPaths); v != "" {
		raw := strings.Split(v, ";")
		paths := make([]string, 0, len(raw))
		for _, p := range raw {
			if p = strings.TrimSpace(os.ExpandEnv(p)); p != "" {
				paths = append(paths, p)
			}
		}
		if len(paths) > 0 {
			return paths
		}
	}
	return []string{os.ExpandEnv(defaultInventoryPath)}
}

func GetDefaultUser() string {
	if user := os.Getenv(envDefaultUser); user != "" {
		return user
	}
	return os.Getenv("USER")
}

func GetAnsibleInventories() []string {
	if v := os.Getenv(envAnsibleInvs); v != "" {
		return strings.Split(v, ";")
	}
	return nil
}

func GetDefaultSshKeyPath() string {
	return os.Getenv(envSshKeyPath)
}

// IsDebugMode returns true when TAQ_DEBUG is set to any non-empty value.
func IsDebugMode() bool {
	return os.Getenv(envDebug) != ""
}

// GetSSHTimeout returns the value of TAQ_SSH_TIMEOUT (seconds).
// An empty string means no explicit timeout is set.
func GetSSHTimeout() string {
	return os.Getenv(envSSHTimeout)
}
