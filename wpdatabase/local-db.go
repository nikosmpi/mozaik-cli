package wpdatabase

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/nikosmpi/mozaik-cli/wpconfig"
)

func LocalDB(config wpconfig.WPConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		config.DBUser,
		config.DBPass,
		config.DBHost,
		config.DBName,
	)
	db, err := sql.Open("sql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
