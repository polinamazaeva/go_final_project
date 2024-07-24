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

func CheckOpenCloseDb() {
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

	Database, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer Database.Close()

	if install {
		query := `CREATE TABLE IF NOT EXISTS scheduler (
                id      INTEGER PRIMARY KEY AUTOINCREMENT,go run .
                date    CHAR(8) NOT NULL DEFAULT "",
                title   VARCHAR(128) NOT NULL DEFAULT "",
                comment TEXT NOT NULL DEFAULT "", 
                repeat VARCHAR(128) NOT NULL DEFAULT "" 
            );`

		// SQL-запрос для создания индекса
		query = `CREATE INDEX IF NOT EXISTS date_indx ON scheduler (date);`

		_, err = Database.Exec(query)
		if err != nil {
			log.Println("Error to create db", err)
		}
	}
}
