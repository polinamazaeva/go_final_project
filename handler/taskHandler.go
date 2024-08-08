package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"go_final_project/storage"
	"go_final_project/task"
)

// IDTask представляет собой структуру для хранения ID задачи, возвращаемого при создании задачи
type IDTask struct {
	Id int64 `json:"id"`
}

// emptyTask используется для отправки пустого JSON-ответа при обновлении задачи
var emptyTask = task.Task{
	Id:      "",
	Date:    "",
	Title:   "",
	Comment: "",
	Repeat:  "",
}

// TaskHandler обрабатывает HTTP-запросы для задач
// В зависимости от метода запроса, выполняется получение, создание, обновление или удаление задачи
func TaskHandler(db *storage.TaskStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// Извлекаем параметр ID из строки запроса
		queryParams := req.URL.Query().Get("id")

		var response []byte
		var err error
		var RespStatus int

		// Определяем действие в зависимости от метода HTTP-запроса
		switch req.Method {
		case http.MethodGet:
			// Обрабатываем запрос на получение задачи по ID
			if queryParams == "" {
				// Если параметр ID отсутствует, возвращаем ошибку 400 Bad Request
				http.Error(w, `{"error":"incorrect id"}`, http.StatusBadRequest)
				return
			}
			// Ищем задачу в базе данных
			task, err := db.TaskID(queryParams)
			if err != nil {
				if err == sql.ErrNoRows {
					// Если задача не найдена, возвращаем ошибку 404 Not Found
					http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
					return
				}
				// Если произошла ошибка при поиске, возвращаем ошибку 500 Internal Server Error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Преобразуем задачу в JSON
			response, err = json.Marshal(task)
			if err != nil {
				// Если произошла ошибка при преобразовании, возвращаем ошибку 500 Internal Server Error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			RespStatus = http.StatusOK

		case http.MethodPost:
			// Обрабатываем запрос на создание новой задачи
			var idtask IDTask
			// Проверяем и извлекаем данные задачи из запроса
			task, RespStatus, err := Check(req)
			if err != nil {
				// Если произошла ошибка при проверке, возвращаем соответствующую ошибку и статус
				http.Error(w, err.Error(), RespStatus)
				return
			}
			// Добавляем новую задачу в базу данных
			idtask.Id, err = db.AddTask(task)
			if err != nil {
				// Если произошла ошибка при добавлении задачи, возвращаем ошибку 500 Internal Server Error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Преобразуем ID новой задачи в JSON
			response, err = json.Marshal(idtask)
			if err != nil {
				// Если произошла ошибка при преобразовании, возвращаем ошибку 500 Internal Server Error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			RespStatus = http.StatusOK

		case http.MethodPut:
			// Обрабатываем запрос на обновление существующей задачи
			task, RespStatus, err := Check(req)
			if err != nil {
				// Если произошла ошибка при проверке, возвращаем соответствующую ошибку и статус
				http.Error(w, err.Error(), RespStatus)
				return
			}
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
			// Преобразуем пустую задачу в JSON (используется для подтверждения успешного обновления)
			response, err = json.Marshal(emptyTask)
			if err != nil {
				// Если произошла ошибка при преобразовании, возвращаем ошибку 500 Internal Server Error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			RespStatus = http.StatusOK

		case http.MethodDelete:
			// Обрабатываем запрос на удаление задачи
			err := db.DeleteTask(queryParams)
			if err != nil {
				if err.Error() == `{"error":"not found the task"}` {
					// Если задача не найдена при удалении, возвращаем ошибку 404 Not Found
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				// Если произошла ошибка при удалении, возвращаем ошибку 500 Internal Server Error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Создаем пустой JSON для ответа (используется для подтверждения успешного удаления)
			str := map[string]interface{}{}
			response, err = json.Marshal(str)
			if err != nil {
				// Если произошла ошибка при преобразовании, возвращаем ошибку 500 Internal Server Error
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			RespStatus = http.StatusOK

		default:
			// Если метод запроса не поддерживается, возвращаем ошибку 405 Method Not Allowed
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Устанавливаем заголовок ответа и отправляем ответ клиенту
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(RespStatus)
		_, err = w.Write(response)
		if err != nil {
			// Логируем ошибку записи ответа
			log.Printf("write response TaskHandler: %v", err)
		}
	}
}
