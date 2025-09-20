package config

import (
	"os"
)

const defaultInventoryPath = "~/.config/taq.inventory.yaml"

func GetDefaultInventoryPath() string {
	const envInventoryPath = "TAQ_INVENTORY_PATH"
	if path := os.Getenv(envInventoryPath); path != "" {
		return path
	}
	return defaultInventoryPath
}
