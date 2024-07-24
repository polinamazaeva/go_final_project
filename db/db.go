package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var Database *sql.DB

func CheckOpenCloseDb() (*sql.DB, error) {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}

	Database, err = sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	if install {
		// SQL-запрос для создания таблицы
		createTableQuery := `CREATE TABLE IF NOT EXISTS scheduler (
			id      INTEGER PRIMARY KEY AUTOINCREMENT,
			date    CHAR(8) NOT NULL DEFAULT "",
			title   VARCHAR(128) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "", 
			repeat  VARCHAR(128) NOT NULL DEFAULT "" 
		);`

		_, err = Database.Exec(createTableQuery)
		if err != nil {
			log.Println("Error creating table", err)
			return nil, err
		}

		// SQL-запрос для создания индекса
		createIndexQuery := `CREATE INDEX IF NOT EXISTS date_indx ON scheduler (date);`

		_, err = Database.Exec(createIndexQuery)
		if err != nil {
			log.Println("Error creating index", err)
			return nil, err
		}
	}
	return Database, nil
}
