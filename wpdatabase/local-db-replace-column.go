package wpdatabase

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/elliotchance/phpserialize"
	"github.com/nikosmpi/mozaik-cli/wpconfig"
)

func LocalDBReplaceInColumn(db *sql.DB, tableName string, columnName string, data wpconfig.Replace) (int64, error) {
	pkName, err := getPrimaryKeyColumn(db, tableName)
	if err != nil {
		return 0, fmt.Errorf("δεν βρέθηκε primary key για τον πίνακα %s: %v", tableName, err)
	}
	querySelect := fmt.Sprintf(
		"SELECT `%s`, `%s` FROM `%s` WHERE `%s` LIKE ?",
		pkName, columnName, tableName, columnName,
	)
	likePattern := "%" + data.Old + "%"
	rows, err := db.Query(querySelect, likePattern)
	if err != nil {
		return 0, fmt.Errorf("σφάλμα στο select: %v", err)
	}
	defer rows.Close()
	var count int64 = 0
	updates := make(map[interface{}]string)
	for rows.Next() {
		var id interface{}
		var content string
		if err := rows.Scan(&id, &content); err != nil {
			return count, err
		}
		newContent, changed, err := smartReplace(content, data)
		if err != nil {
			fmt.Printf("Warning: Failed to parse serialization for ID %v: %v\n", id, err)
			continue
		}
		if changed {
			updates[id] = newContent
		}
	}
	queryUpdate := fmt.Sprintf("UPDATE `%s` SET `%s` = ? WHERE `%s` = ?", tableName, columnName, pkName)
	stmt, err := db.Prepare(queryUpdate)
	if err != nil {
		return count, err
	}
	defer stmt.Close()
	for id, newContent := range updates {
		res, err := stmt.Exec(newContent, id)
		if err != nil {
			fmt.Printf("Error updating row %v: %v\n", id, err)
			continue
		}
		affected, _ := res.RowsAffected()
		count += affected
	}
	return count, nil
}

func getPrimaryKeyColumn(db *sql.DB, tableName string) (string, error) {
	query := `
		SELECT COLUMN_NAME
		FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE()
		AND TABLE_NAME = ?
		AND COLUMN_KEY = 'PRI'
		LIMIT 1
	`
	var columnName string
	err := db.QueryRow(query, tableName).Scan(&columnName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("δεν βρέθηκε primary key για τον πίνακα %s", tableName)
		}
		return "", err
	}
	return columnName, nil
}

func smartReplace(content string, data wpconfig.Replace) (string, bool, error) {
	var obj interface{}
	err := phpserialize.Unmarshal([]byte(content), &obj)
	if err == nil {
		changed, newObj := recursiveWalk(obj, data)
		if changed {
			newData, err := phpserialize.Marshal(newObj, nil)
			if err != nil {
				return "", false, err
			}
			return string(newData), true, nil
		}
		return "", false, nil
	}
	if strings.Contains(content, data.Old) {
		return strings.ReplaceAll(content, data.Old, data.New), true, nil
	}
	return "", false, nil
}

func recursiveWalk(data interface{}, replace wpconfig.Replace) (bool, interface{}) {
	changed := false
	switch v := data.(type) {
	case string:
		if strings.Contains(v, replace.Old) {
			return true, strings.ReplaceAll(v, replace.Old, replace.New)
		}
		return false, v
	case map[interface{}]interface{}:
		newMap := make(map[interface{}]interface{})
		for key, val := range v {
			c, newVal := recursiveWalk(val, replace)
			if c {
				changed = true
			}
			newMap[key] = newVal
		}
		return changed, newMap

	default:
		return false, v
	}
}
