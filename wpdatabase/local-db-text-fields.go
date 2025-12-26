package wpdatabase

import (
	"database/sql"
	"fmt"
)

func LocalDBTextFields(db *sql.DB, tableName string) ([]string, error) {
	query := fmt.Sprintf("SHOW COLUMNS FROM `%s`", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var textFields []string
	var field, colType, null, key string
	var def, extra []byte // Το Default και Extra μπορεί να είναι null, οπότε τα βάζουμε ως bytes
	for rows.Next() {
		err := rows.Scan(&field, &colType, &null, &key, &def, &extra)
		if err != nil {
			return nil, err
		}
		typeStr := string(colType)

		if isTextType(typeStr) {
			textFields = append(textFields, field)
		}
	}
	return textFields, nil
}

func isTextType(sqlType string) bool {
	return (len(sqlType) >= 4 && sqlType[:4] == "char") ||
		(len(sqlType) >= 7 && sqlType[:7] == "varchar") ||
		(len(sqlType) >= 4 && sqlType[len(sqlType)-4:] == "text")
}
