package db

import (
	"database/sql"
	"fmt"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"

	"log"
)

var db *sql.DB

func Init(dbName string) error {
	var err error
	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %s", err)
	}
	log.Println("Connected Successfully to the Database")
	return createMigrations()
}

func createMigrations() error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT, 
			username TEXT NOT NULL,
			email TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS articles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			portal TEXT NOT NULL,
			link TEXT NOT NULL,
			published_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
	}

	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute migration statement: %s", err)
		}
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}
