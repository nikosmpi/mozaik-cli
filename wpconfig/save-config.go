package wpconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveConfig() (bool, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return false, err
	}
	configPath := filepath.Join(currentDir, "moz-config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config, err := CreateConfig()
		if err != nil {
			return false, err
		}
		data, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			return false, err
		}
		err = os.WriteFile(configPath, data, 0644)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}
