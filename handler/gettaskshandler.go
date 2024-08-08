package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go_final_project/storage"
	"go_final_project/task"
)

// GetTasksHandler возвращает HTTP-обработчик для получения задач.
// Обработчик может возвращать либо одну задачу по ID, либо список задач.
func GetTasksHandler(ts *storage.TaskStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var response []byte    // Переменная для хранения ответа в виде JSON
		var responseStatus int // Переменная для хранения статуса HTTP-ответа
		var err error          // Переменная для хранения ошибок

		// Проверяем, что метод запроса GET. Если нет, возвращаем ошибку 405 Method Not Allowed.
		if req.Method != http.MethodGet {
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
			return
		}

		// Извлекаем параметры запроса: id задачи и лимит количества задач
		id := req.URL.Query().Get("id")
		limitParam := req.URL.Query().Get("limit")
		limit := storage.DefaultLimit // Используем значение по умолчанию для лимита

		// Если параметр лимита передан, пытаемся его преобразовать в число
		if limitParam != "" {
			if l, parseErr := strconv.Atoi(limitParam); parseErr == nil && l >= 10 && l <= 50 {
				limit = l // Если преобразование прошло успешно и значение в допустимом диапазоне, обновляем лимит
			}
		}

		// Если указан параметр id, пытаемся получить задачу по этому ID
		if id != "" {
			taskData, err := ts.GetTaskByID(id)
			if err != nil {
				// Если задача не найдена, возвращаем ошибку 404 Not Found
				if err.Error() == `{"error":"task not found"}` {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				// В случае другой ошибки возвращаем 500 Internal Server Error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Если задача найдена, сериализуем её в JSON
			response, err = json.Marshal(taskData)
			responseStatus = http.StatusOK
		} else {
			// Если параметр id не указан, возвращаем список задач с учетом лимита
			tasks, err := ts.GetTasks(limit)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Сериализуем список задач в JSON, обернув его в объект с ключом "tasks"
			response, err = json.Marshal(map[string][]task.Task{"tasks": tasks})
			responseStatus = http.StatusOK
		}

		// Если возникла ошибка при сериализации JSON, возвращаем её с соответствующим HTTP-статусом
		if err != nil {
			http.Error(w, err.Error(), responseStatus)
			return
		}

		// Устанавливаем заголовок Content-Type и статус ответа
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(responseStatus)
		// Пишем JSON-ответ в тело ответа
		w.Write(response)
	}
}
