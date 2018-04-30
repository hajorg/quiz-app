package database

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strconv"
)

var dbName string

func init() {
	if flag.Lookup("test.v") == nil {
		dbName = "quiz"
	} else {
		dbName = "quiztest"
	}
}

// Connect connects to the local database
func Connect(name string) *sql.DB {
	db, err := sql.Open("mysql", "root:guesswho@tcp(127.0.0.1:3306)/"+name)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}

// Insert inserts into a table
func Insert(table string, data map[string]interface{}) (int64, error) {
	var newData []map[string]interface{}

	for key, val := range data {
		newData = append(newData, map[string]interface{}{key: val})
	}

	db := Connect(dbName)
	defer db.Close()

	placeholder := "" // used as prepared statement binding i.e ?
	query := "INSERT INTO " + fmt.Sprint(table) + "("

	content := []interface{}{}
	for idx, val := range newData {
		for key := range val {
			query += ("" + fmt.Sprint(key))
			placeholder += "?"
			content = append(content, val[key])
			if idx != len(newData)-1 {
				query += ", "
				placeholder += ", "
			}
		}
	}

	query += ") VALUES(" + placeholder + ")"
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(content...)

	if err != nil {
		return 0, err
	}
	lastID, _ := result.LastInsertId()
	return lastID, nil
}

// GetAll gets all data for a paticular table
func GetAll(table string) []map[string]interface{} {
	db := Connect(dbName)
	defer db.Close()

	rows, err := db.Query("SELECT * FROM " + table)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// get table columns
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	// placeholder of result from the database
	data := make([]interface{}, len(columns))
	// an array of the results to be gotten from the db
	newData := make([]map[string]interface{}, 0)
	// holds the address of each interface{} value in the data slice
	scanArgs := make([]interface{}, len(columns))
	for i := range data {
		scanArgs[i] = &data[i]
	}

	for {
		if rows.Next() {
			err = rows.Scan(scanArgs...)
			if err != nil {
				panic(err)
			}
			entry := make(map[string]interface{})
			for i, col := range columns {
				val, ok := data[i].([]byte)
				if ok {
					entry[col] = string(val)
				} else {
					entry[col] = data[i]
				}
				// convert from interface to string, float, boolean etc
				switch value := entry[col].(type) {
				case string:
					if boolean, err := strconv.ParseBool(value); err == nil {
						entry[col] = boolean
					}

					if val, err := strconv.ParseFloat(value, 64); err == nil {
						entry[col] = val
					}
				}
			}
			newData = append(newData, entry)
		} else {
			break
		}
	}

	return newData
}

// GetWhere gets all data from a table where it meets the `where` condition
func GetWhere(table string, where []map[string]interface{}) []map[string]interface{} {
	db := Connect(dbName)
	defer db.Close()

	whereCount := len(where)
	whereValues := []interface{}{}
	whereCondition := "SELECT * FROM " + table + " WHERE"
	for idx, arr := range where {
		for key, value := range arr {
			whereCondition += " " + key + " = " + "?"
			whereValues = append(whereValues, value)
			if idx != whereCount-1 {
				whereCondition += " AND"
			}
		}
	}

	stmt, err := db.Prepare(whereCondition)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(whereValues...)
	if err != nil {
		panic(err)
	}

	// get table columns
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	// placeholder of result from the database
	data := make([]interface{}, len(columns))
	// an array of the results to be gotten from the db
	newData := make([]map[string]interface{}, 0)
	// holds the address of each interface{} value in the data slice
	scanArgs := make([]interface{}, len(columns))
	for idx := range data {
		scanArgs[idx] = &data[idx]
	}

	for {
		if rows.Next() {
			err = rows.Scan(scanArgs...)
			if err != nil {
				panic(err)
			}
			entry := make(map[string]interface{})
			for i, col := range columns {
				val, ok := data[i].([]byte)
				if ok {
					entry[col] = string(val)
				} else {
					entry[col] = data[i]
				}
				// convert from interface to string, float, boolean etc
				switch value := entry[col].(type) {
				case string:
					if boolean, err := strconv.ParseBool(value); err == nil {
						entry[col] = boolean
					}

					if val, err := strconv.ParseFloat(value, 64); err == nil {
						entry[col] = val
					}
				}
			}
			newData = append(newData, entry)
		} else {
			break
		}
	}

	return newData
}
