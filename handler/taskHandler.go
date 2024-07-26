package handler

import (
	"database/sql"
	"encoding/json"
	"go_final_project/nextdate"
	"go_final_project/task"
	"go_final_project/taskactions"
	"net/http"
	"strconv"
	"time"
)

func MakeTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var t task.Task
		var err error
		param := req.URL.Query().Get("id")
		var response []byte
		var responseStatus int

		switch req.Method {

		case http.MethodPost:
			decoder := json.NewDecoder(req.Body)
			err = decoder.Decode(&t)
			if err != nil {
				http.Error(w, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
				return
			}

			if t.Title == "" {
				http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
				return
			}

			layout := "20060102"
			now := time.Now().Truncate(24 * time.Hour)

			var taskDate time.Time
			if t.Date == "today" || t.Date == "" {
				taskDate = now
			} else {
				taskDate, err = time.Parse(layout, t.Date)
				if err != nil {
					http.Error(w, `{"error":"Дата представлена в неверном формате"}`, http.StatusBadRequest)
					return
				}
			}

			taskDate = taskDate.Truncate(24 * time.Hour)

			if taskDate.Format("20060102") < now.Format("20060102") {
				if t.Repeat == "" {
					taskDate = now
				} else {
					nextDate, err := nextdate.NextDate(now, taskDate.Format(layout), t.Repeat)
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
			}

			t.Date = taskDate.Format(layout)

			response, responseStatus, err = taskactions.AddTask(db, t)
			if err != nil {
				http.Error(w, err.Error(), responseStatus)
				return
			}

		case http.MethodGet:
			limitParam := req.URL.Query().Get("limit")
			limit, err := strconv.Atoi(limitParam)
			if err != nil || limit < 10 || limit > 50 {
				limit = 50 // Default limit
			}

			if param != "" {
				response, responseStatus, err = taskactions.GetTaskByID(db, param)
			} else {
				response, responseStatus, err = taskactions.GetTasks(db, limit)
			}

			if err != nil {
				http.Error(w, err.Error(), responseStatus)
				return
			}

		case http.MethodPut:
			decoder := json.NewDecoder(req.Body)
			err = decoder.Decode(&t)
			if err != nil {
				http.Error(w, `{"error":"Ошибка десериализации JSON"}`, http.StatusBadRequest)
				return
			}

			if t.Id == "" {
				http.Error(w, `{"error":"Не указан идентификатор задачи"}`, http.StatusBadRequest)
				return
			}

			if t.Title == "" {
				http.Error(w, `{"error":"Не указан заголовок задачи"}`, http.StatusBadRequest)
				return
			}

			layout := "20060102"
			now := time.Now().Truncate(24 * time.Hour)

			var taskDate time.Time
			if t.Date == "today" || t.Date == "" {
				taskDate = now
			} else {
				taskDate, err = time.Parse(layout, t.Date)
				if err != nil {
					http.Error(w, `{"error":"Дата представлена в неверном формате"}`, http.StatusBadRequest)
					return
				}
			}

			taskDate = taskDate.Truncate(24 * time.Hour)

			if taskDate.Format("20060102") < now.Format("20060102") {
				if t.Repeat == "" {
					taskDate = now
				} else {
					nextDate, err := nextdate.NextDate(now, taskDate.Format(layout), t.Repeat)
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
			}

			t.Date = taskDate.Format(layout)

			response, responseStatus, err = taskactions.UpdateTask(db, t)
			if err != nil {
				http.Error(w, err.Error(), responseStatus)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
