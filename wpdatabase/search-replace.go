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
	//tables, err := LocalDBTables(db)
	//if err != nil {
	//	return err
	//}
	for _, data := range config.ReplaceList {
		fmt.Printf("Old: %s, New: %s\n", data.Old, data.New)
	}

	return nil
}
