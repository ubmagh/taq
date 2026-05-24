package config

import (
	"os"
	"strings"
)

const (
	defaultInventoryPath = "$HOME/.config/taq/inventory.yaml"

	envInventoryPath = "TAQ_INVENTORY_PATH"
	envDefaultUser   = "TAQ_DEFAULT_USER"
	envAnsibleInvs   = "TAQ_ANSIBLE_INVS"
	envSshKeyPath    = "TAQ_DEFAULT_SSH_KEY_PATH"
	envDisplayMode   = "TAQ_DISPLAY_MODE"
	envDebug         = "TAQ_DEBUG"
	envSSHTimeout    = "TAQ_SSH_TIMEOUT"
)

func IsCompactMode() bool {
	return os.Getenv(envDisplayMode) == "compact"
}

func GetDefaultInventoryPath() string {
	path := os.Getenv(envInventoryPath)
	if path == "" {
		path = defaultInventoryPath
	}
	return os.ExpandEnv(path)
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
