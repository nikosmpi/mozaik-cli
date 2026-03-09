package wpdatabase

import (
	"database/sql"
	"fmt"
	"strings"

	// Χρησιμοποιούμε το υπο-πακέτο php_serialize που παρέχει πλήρη έλεγχο στα Objects
	"github.com/nikosmpi/mozaik-cli/wpconfig"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
)

func LocalDBReplaceInColumn(db *sql.DB, tableName string, columnName string, data wpconfig.Replace) (int64, error) {
	pkName, err := getPrimaryKeyColumn(db, tableName)
	if err != nil {
		return 0, fmt.Errorf("primary key not found for table %s: %v", tableName, err)
	}

	querySelect := fmt.Sprintf(
		"SELECT `%s`, `%s` FROM `%s` WHERE `%s` LIKE ?",
		pkName, columnName, tableName, columnName,
	)

	likePattern := "%" + data.Old + "%"
	rows, err := db.Query(querySelect, likePattern)
	if err != nil {
		return 0, fmt.Errorf("error in select: %v", err)
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
			// Αν αποτύχει η αποκωδικοποίηση, το αγνοούμε και συνεχίζουμε
			// (Μπορείς να βγάλεις το σχόλιο από κάτω για debugging)
			// fmt.Printf("Warning: Serialization error for ID %v: %v\n", id, err)
			continue
		}

		if changed {
			updates[id] = newContent
		}
	}

	// Batch Update
	if len(updates) > 0 {
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
	}

	return count, nil
}

func smartReplace(content string, data wpconfig.Replace) (string, bool, error) {
	// Προσπάθεια Unserialize με τη βιβλιοθήκη php_session_decoder
	decoder := php_serialize.NewUnSerializer(content)
	obj, err := decoder.Decode()

	if err == nil {
		// Είναι έγκυρο serialized string (Object, Array, κλπ)
		changed, newObj := recursiveWalk(obj, data)

		if changed {
			encoder := php_serialize.NewSerializer()
			newData, err := encoder.Encode(newObj)
			if err != nil {
				return "", false, err
			}
			return newData, true, nil
		}
		return "", false, nil
	}

	// Fallback: Αν δεν είναι serialized, κάνουμε απλό string replace
	if strings.Contains(content, data.Old) {
		return strings.ReplaceAll(content, data.Old, data.New), true, nil
	}

	return "", false, nil
}

// recursiveWalk διατρέχει τη δομή και χρησιμοποιεί setters/getters για τα Objects
func recursiveWalk(data interface{}, replace wpconfig.Replace) (bool, interface{}) {
	changed := false

	switch v := data.(type) {
	case string:
		// Check for Double Serialization (nested serialized string)
		decoder := php_serialize.NewUnSerializer(v)
		if nestedObj, err := decoder.Decode(); err == nil {
			// Αν βρήκαμε ότι το string είναι serialized data, μπαίνουμε μέσα
			nestedChanged, newNestedObj := recursiveWalk(nestedObj, replace)
			if nestedChanged {
				encoder := php_serialize.NewSerializer()
				if newBytes, err := encoder.Encode(newNestedObj); err == nil {
					return true, newBytes
				}
			}
		}

		// Απλό replace στο string
		if strings.Contains(v, replace.Old) {
			return true, strings.ReplaceAll(v, replace.Old, replace.New)
		}
		return false, v

	case *php_serialize.PhpObject:
		// Χρησιμοποιούμε τους Getters/Setters της βιβλιοθήκης
		// v.GetMembers() επιστρέφει map[PhpValue]PhpValue
		members := v.GetMembers()
		newMembers := make(php_serialize.PhpArray) // PhpArray είναι alias του map

		objectChanged := false
		for key, val := range members {
			// Αναδρομή στην τιμή
			c, newVal := recursiveWalk(val, replace)
			if c {
				objectChanged = true
			}
			newMembers[key] = newVal
		}

		if objectChanged {
			v.SetMembers(newMembers)
			changed = true
		}
		return changed, v

	case php_serialize.PhpArray: // map[interface{}]interface{}
		newMap := make(php_serialize.PhpArray)
		for key, val := range v {
			c, newVal := recursiveWalk(val, replace)
			if c {
				changed = true
			}
			newMap[key] = newVal
		}
		return changed, newMap

	case map[interface{}]interface{}:
		// Fallback για απλά maps αν προκύψουν
		newMap := make(map[interface{}]interface{})
		for key, val := range v {
			c, newVal := recursiveWalk(val, replace)
			if c {
				changed = true
			}
			newMap[key] = newVal
		}
		return changed, newMap

	case []interface{}: // Slices
		newSlice := make([]interface{}, len(v))
		for i, val := range v {
			c, newVal := recursiveWalk(val, replace)
			if c {
				changed = true
			}
			newSlice[i] = newVal
		}
		return changed, newSlice

	default:
		return false, v
	}
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
			return "", fmt.Errorf("primary key not found for table %s", tableName)
		}
		return "", err
	}
	return columnName, nil
}
