package wpdatabase

import (
	"fmt"

	"github.com/nikosmpi/mozaik-cli/wpconfig"
)

func SearchReplace(config wpconfig.WPConfig) error {
	db, err := LocalDB(config)
	if err != nil {
		return err
	}
	defer db.Close()
	tables, err := LocalDBTables(db)
	if err != nil {
		return err
	}
	for _, data := range config.ReplaceList {
		mainCount := 0
		fmt.Printf("Old: %s, New: %s\n", data.Old, data.New)
		for _, table := range tables {
			fields, err := LocalDBTextFields(db, table)
			if err != nil {
				return err
			}
			for _, field := range fields {
				count, err := LocalDBReplaceInColumn(db, table, field, data)
				if err != nil {
					return err
				}
				if count > 0 {
					fmt.Printf("%s.%s (%d)\n", table, field, count)
					mainCount += int(count)
				}
			}
		}
		fmt.Printf("Total: %d\n", mainCount)
	}

	return nil
}
