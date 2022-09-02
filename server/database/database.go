package database

import (
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func exec(db *sqlx.DB, query string) {

	_, err := db.Exec(query)

	if err != nil {
		log.Fatalln("Query '", query, "' failed with error", err.Error())
	}
}

func initTables(db *sqlx.DB) {

	const create_str = "CREATE TABLE IF NOT EXISTS "

	exec(db, create_str+`users
		(
			id TEXT PRIMARY KEY,	
			username TEXT UNIQUE, 
			email TEXT UNIQUE,
			bio TEXT, 
			salt TEXT,
			bookmarks TEXT [],
			password BYTEA
		)`)

	exec(db, create_str+`albums
		(
			id TEXT PRIMARY KEY, 
			name TEXT, 
			authorId TEXT, 
			description TEXT,
			createdAt  TIMESTAMP,
			likes INTEGER DEFAULT 0,
			tags TEXT [], 
			files TEXT []
		)`)

	exec(db, create_str+`sessions
		(
			id TEXT PRIMARY KEY,
			userId TEXT
		)
	`)

	exec(db, create_str+`comments
		(
			albumId TEXT PRIMARY KEY,
			authorId TEXT,
			contents TEXT,
			createdAt TIMESTAMP
		)
	`)
}

func readDbConfig() string {
	bytes, err := os.ReadFile("./psql_config")

	if err != nil {
		log.Fatal("could not read psql_config file")
	}

	return string(bytes)
}

func New() *sqlx.DB {

	conn_str := readDbConfig()

	db, err := sqlx.Open("postgres", conn_str)

	if err != nil {
		log.Fatal(err.Error())
	}

	initTables(db)

	return db
}
