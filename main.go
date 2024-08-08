package main

import (
	"log"
	"net/http"

	"go_final_project/handler"
	"go_final_project/storage"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Открываем соединение с базой данных, используя функцию CheckOpenCloseDb из пакета storage
	database, err := storage.CheckOpenCloseDb()
	if err != nil {
		// Если возникает ошибка при инициализации базы данных, логируем её и завершаем программу
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// Закрываем соединение с базой данных, когда main функция завершится
	defer database.Close()

	// Создаем экземпляр структуры TaskStorage, передавая ему подключение к базе данных
	ts := &storage.TaskStorage{DB: database}

	// Указываем директорию, из которой будут раздаваться статические файлы для веб-интерфейса
	webDir := "./web"

	// Настраиваем обработчики
	http.Handle("/", http.FileServer(http.Dir(webDir)))            // Настраиваем обработчик для раздачи статических файлов из директории webDir
	http.HandleFunc("/api/nextdate", handler.NextDateHandler)      // Настраиваем обработчик для API-запросов на получение следующей даты (NextDateHandler)
	http.HandleFunc("/api/task", handler.TaskHandler(ts))          // Настраиваем обработчик для API-запросов на добавление или получение задачи (TaskHandler)
	http.HandleFunc("/api/tasks", handler.GetTasksHandler(ts))     // Настраиваем обработчик для API-запросов на получение всех задач (GetTasksHandler)
	http.HandleFunc("/api/task/done", handler.TaskDoneHandler(ts)) // Настраиваем обработчик для API-запросов на отметку задачи как выполненной (TaskDoneHandler)

	// Логируем сообщение о том, что сервер стартует на порту 7540
	log.Println("Server starting on :7540")

	// Запускаем HTTP-сервер на порту 7540 и обрабатываем запросы
	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		// Если возникает ошибка при запуске сервера, логируем её и завершаем программу
		log.Fatalf("Server failed: %v", err)
	}
}
