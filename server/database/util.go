package database

import (
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

func CheckRow(db *sqlx.DB, table, column string, value any) bool {
	row := db.QueryRow("SELECT EXISTS (SELECT 1 FROM " + table + " WHERE " + column + " = $1)", value)

	if row.Err() != nil {
		log.Println(row.Err().Error())
		return false
	}

	var exists bool 

	row.Scan(&exists)

	return exists
}