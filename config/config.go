package config

import (
	"os"
	"strings"
)

const defaultInventoryPath = "~/.config/taq/inventory.yaml"

func GetDefaultInventoryPath() string {
	const envInventoryPath = "TAQ_INVENTORY_PATH"
	if path := os.Getenv(envInventoryPath); path != "" {
		return path
	}
	return defaultInventoryPath
}

func GetDefaultUser() string {
	const envDefaultUser = "TAQ_DEFAULT_USER"
	if user := os.Getenv(envDefaultUser); user != "" {
		return user
	}
	return os.Getenv("USER")
}

func GetAnsibleInventories() []string {
	const envAnsibleInvs = "TAQ_ANSIBLE_INVS"
	if inventoriesStr := os.Getenv(envAnsibleInvs); inventoriesStr != "" {
		return strings.Split(inventoriesStr, ";")
	}
	return []string{}
}

func GetDefaultSshKeyPath() string {
	const envSshKeyPath = "TAQ_DEFAULT_SSH_KEY_PATH"
	if path := os.Getenv(envSshKeyPath); path != "" {
		return path
	}
	return ""
}
