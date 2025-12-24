package wpconfig

import (
	"os"
	"path/filepath"
	"strings"
)

func AddGitignore() error {
	currentDir, err := os.Getwd()
	if err != nil {
		return err
	}
	gitignorePath := filepath.Join(currentDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		return nil
	}
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		return err
	}

	if !strings.Contains(string(content), "moz-config.json") {
		f, err := os.OpenFile(gitignorePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := f.WriteString("moz-config.json\n"); err != nil {
			return err
		}
	}
	return nil
}
