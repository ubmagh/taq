package config

import (
	"os"
	"strings"
)

const (
	defaultInventoryPath = "~/.config/taq/inventory.yaml"

	envInventoryPath = "TAQ_INVENTORY_PATH"
	envDefaultUser   = "TAQ_DEFAULT_USER"
	envAnsibleInvs   = "TAQ_ANSIBLE_INVS"
	envSshKeyPath    = "TAQ_DEFAULT_SSH_KEY_PATH"
)

func GetDefaultInventoryPath() string {
	if path := os.Getenv(envInventoryPath); path != "" {
		return path
	}
	return defaultInventoryPath
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
