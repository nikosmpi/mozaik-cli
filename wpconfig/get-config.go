package wpconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func GetConfig() (WPConfig, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return WPConfig{}, err
	}
	configPath := filepath.Join(currentDir, "moz-config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return WPConfig{}, nil
	}
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return WPConfig{}, err
	}
	var config WPConfig
	if err := json.Unmarshal(configData, &config); err != nil {
		return WPConfig{}, err
	}
	return config, nil
}
