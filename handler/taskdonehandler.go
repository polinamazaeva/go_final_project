package handler

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"go_final_project/nextdate"
	"go_final_project/storage"
)

// TaskDoneHandler обрабатывает запрос на завершение задачи
// Возвращает обработчик HTTP-запроса для завершения задачи
func TaskDoneHandler(db *storage.TaskStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Получаем ID задачи из параметров запроса
		id := req.URL.Query().Get("id")
		if id == "" {
			// Если ID не предоставлен, возвращаем ошибку 400 Bad Request
			http.Error(w, `{"error":"missing id parameter"}`, http.StatusBadRequest)
			return
		}

		// Ищем задачу по ID в базе данных
		task, err := db.TaskDone(id)
		if err != nil {
			if err == sql.ErrNoRows {
				// Если задача не найдена, возвращаем ошибку 404 Not Found
				http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
				return
			}
			// Если произошла ошибка при сканировании, возвращаем ошибку 500 Internal Server Error
			http.Error(w, `{"error":"error scanning task: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}

		// Проверяем, указано ли поле repeat
		if task.Repeat == "" {
			// Если поле repeat пустое, удаляем задачу
			err := db.DeleteTask(id)
			if err != nil {
				if err.Error() == "task not found" {
					// Если задача не найдена при удалении, возвращаем ошибку 404 Not Found
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				// Если произошла ошибка при удалении, возвращаем ошибку 500 Internal Server Error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// Если поле repeat указано, обновляем дату задачи
			now := time.Now()
			// Вычисляем следующую дату повторения задачи
			dateNew, err := nextdate.NextDate(now, task.Date, task.Repeat)
			if err != nil {
				// Если произошла ошибка при вычислении следующей даты, возвращаем ошибку 400 Bad Request
				http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
				return
			}

			// Преобразуем строку с новой датой в формат времени
			nextDate, err := time.Parse("20060102", dateNew)
			if err != nil {
				// Если произошла ошибка при парсинге новой даты, возвращаем ошибку 500 Internal Server Error
				http.Error(w, `{"error":"error parsing next date: `+err.Error()+`"}`, http.StatusInternalServerError)
				return
			}

			// Обновляем дату задачи
			task.Date = nextDate.Format("20060102")
			// Обновляем задачу в базе данных
			err = db.UpdateTask(task)
			if err != nil {
				if err.Error() == `{"error":"not found the task"}` {
					// Если задача не найдена при обновлении, возвращаем ошибку 404 Not Found
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				// Если произошла ошибка при обновлении, возвращаем ошибку 500 Internal Server Error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Отправляем пустой JSON-ответ с кодом 200 OK
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(`{}`))
		if err != nil {
			// Логируем ошибку записи ответа
			log.Printf("write response TaskDoneHandler: %v", err)
		}
	}
}
