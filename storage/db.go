package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// Database — глобальная переменная для хранения соединения с базой данных
var Database *sql.DB

// CheckOpenCloseDb открывает базу данных и создает необходимые таблицы, если их нет
func CheckOpenCloseDb() (*sql.DB, error) {
	// Получаем путь к исполняемому файлу приложения
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err) // Логируем и завершаем выполнение программы в случае ошибки
	}

	// Определяем путь к файлу базы данных, размещая его рядом с исполняемым файлом
	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")

	// Проверяем, существует ли файл базы данных
	_, err = os.Stat(dbFile)

	var install bool
	// Если файл базы данных не найден, устанавливаем флаг для создания новой базы
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return nil, fmt.Errorf("failed to check database file: %w", err)
		}
	}

	// Открываем соединение с базой данных SQLite
	Database, err = sql.Open("sqlite", dbFile)
	if err != nil {
		fmt.Println(err) // Выводим ошибку на экран
		return nil, err  // Возвращаем nil и ошибку в случае неудачи
	}

	// Если база данных новая (флаг install установлен), создаем таблицы и индексы
	if install {
		// SQL-запрос для создания таблицы scheduler
		createTableQuery := `CREATE TABLE IF NOT EXISTS scheduler (
			id      INTEGER PRIMARY KEY AUTOINCREMENT,
			date    CHAR(8) NOT NULL DEFAULT "",
			title   VARCHAR(128) NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat  VARCHAR(128) NOT NULL DEFAULT ""
		);`

		// Выполняем SQL-запрос для создания таблицы
		_, err = Database.Exec(createTableQuery)
		if err != nil {
			log.Println("Error creating table:", err) // Логируем ошибку создания таблицы
			return nil, err                           // Возвращаем nil и ошибку
		}

		// SQL-запрос для создания индекса на поле date
		createIndexQuery := `CREATE INDEX IF NOT EXISTS date_indx ON scheduler (date);`
		// Выполняем SQL-запрос для создания индекса
		_, err = Database.Exec(createIndexQuery)
		if err != nil {
			log.Println("Error creating index:", err) // Логируем ошибку создания индекса
			return nil, err                           // Возвращаем nil и ошибку
		}
	}

	// Возвращаем соединение с базой данных
	return Database, nil
}
