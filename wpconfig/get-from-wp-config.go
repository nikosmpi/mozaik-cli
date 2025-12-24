package wpconfig

import (
	"regexp"
)

func GetFromWPConfigPhp(content string) (WPConfig, error) {
	config := WPConfig{}
	reDefine := regexp.MustCompile(`define\s*\(\s*['"](DB_NAME|DB_USER|DB_PASSWORD|DB_HOST)['"]\s*,\s*['"](.*?)['"]\s*\);`)
	matches := reDefine.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		key := match[1]
		value := match[2]
		switch key {
		case "DB_NAME":
			config.DBName = value
		case "DB_USER":
			config.DBUser = value
		case "DB_PASSWORD":
			config.DBPass = value
		case "DB_HOST":
			config.DBHost = value
		}
	}
	rePrefix := regexp.MustCompile(`\$table_prefix\s*=\s*['"](.*?)['"];`)
	prefixMatch := rePrefix.FindStringSubmatch(content)
	if len(prefixMatch) > 1 {
		config.DBPrefix = prefixMatch[1]
	}
	return config, nil
}
