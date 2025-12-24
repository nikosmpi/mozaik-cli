package wpconfig

import (
	"os"
	"path/filepath"
)

func CreateConfig() (WPConfig, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return WPConfig{}, err
	}
	configPath := filepath.Join(currentDir, "wp-config.php")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return WPConfig{}, nil
	}
	content, err := os.ReadFile(configPath)
	if err != nil {
		return WPConfig{}, err
	}
	config, err := GetFromWPConfigPhp(string(content))
	if err != nil {
		return WPConfig{}, err
	}
	replace := Replace{
		Old: "https://old.mozaik.com",
		New: "http://new.mozaik.com",
	}
	config.ReplaceList = make([]Replace, 1)
	config.ReplaceList[0] = replace
	config.MySQLPath = "mysql"
	config.Staging = Remote{
		DBHost:     config.DBHost,
		DBUser:     config.DBUser,
		DBPass:     config.DBPass,
		DBName:     config.DBName,
		DBPrefix:   config.DBPrefix,
		SSHUser:    "staging",
		SSHKeyPath: "staging",
		SSHHost:    "staging",
	}

	return config, nil
}
