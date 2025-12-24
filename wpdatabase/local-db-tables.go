package wpdatabase

import (
	"database/sql"
	"fmt"
)

func LocalDBTables(db *sql.DB) ([]string, error) {
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("δεν είναι δυνατή η σύνδεση στον server: %v", err)
	}
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}
	return tables, nil
}
