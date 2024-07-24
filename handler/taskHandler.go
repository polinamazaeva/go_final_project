package handler

import (
	"database/sql"
	"encoding/json"
	"go_final_project/nextdate"
	"go_final_project/task"
	"go_final_project/taskactions"
	"net/http"
	"time"
)

func MakeTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var task task.Task
		var err error

		if req.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		decoder := json.NewDecoder(req.Body)
		err = decoder.Decode(&task)
		if err != nil {
			http.Error(w, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
			return
		}

		if task.Title == "" {
			http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
			return
		}

		layout := "20060102"
		now := time.Now().Truncate(24 * time.Hour) // Заносим начало текущего дня

		// Обработка даты "today"
		if task.Date == "today" {
			task.Date = now.Format(layout)
		}

		var taskDate time.Time
		if task.Date == "" {
			taskDate = now
		} else {
			taskDate, err = time.Parse(layout, task.Date)
			if err != nil {
				http.Error(w, `{"error":"Дата представлена в неверном формате"}`, http.StatusBadRequest)
				return
			}
		}

		taskDate = taskDate.Truncate(24 * time.Hour)

		if taskDate.Before(now) {
			if task.Repeat == "" {
				// Если задача не повторяется, устанавливаем дату на сегодня
				taskDate = now
			} else {
				// Вычисляем следующую дату для повторения
				nextDate, err := nextdate.NextDate(now, taskDate.Format(layout), task.Repeat)
				if err != nil {
					http.Error(w, `{"error":"Неправильное правило повторения"}`, http.StatusBadRequest)
					return
				}
				taskDate, err = time.Parse(layout, nextDate)
				if err != nil {
					http.Error(w, `{"error":"Неправильный формат даты"}`, http.StatusBadRequest)
					return
				}
			}
		} else if task.Repeat != "" {
			// Если дата не прошла и задача повторяется
			nextDate, err := nextdate.NextDate(now, taskDate.Format(layout), task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Неправильное правило повторения"}`, http.StatusBadRequest)
				return
			}
			taskDate, err = time.Parse(layout, nextDate)
			if err != nil {
				http.Error(w, `{"error":"Неправильный формат даты"}`, http.StatusBadRequest)
				return
			}
		}

		task.Date = taskDate.Format(layout)

		Response, ResponseStatus, err := taskactions.AddTask(db, task)
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(Response)
	}
}
