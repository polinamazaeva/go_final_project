package handler

import (
	"database/sql"
	"encoding/json"
	"go_final_project/nextdate"
	"go_final_project/task"
	"go_final_project/taskactions"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

var db *sql.DB
var ResponseStatus int
var Response []byte

func TaskHandler(w http.ResponseWriter, req *http.Request) {

	var task task.Task
	var err error

	db, err = sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

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
	var taskDate time.Time
	if task.Date == "" {
		taskDate = time.Now()
	} else {
		taskDate, err = time.Parse(layout, task.Date)
		if err != nil {
			http.Error(w, `{"error":"Дата представлена в неверном формате"}`, http.StatusBadRequest)
			return
		}
	}

	today := time.Now()
	if taskDate.Before(today) {
		if task.Repeat == "" {
			taskDate = today
		} else {
			nextDate, err := nextdate.NextDate(today, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Неправильное правило повторения"}`, http.StatusBadRequest)
				return
			}
			taskDate, err = time.Parse("2006-01-02", nextDate)
		}
	} else {
		if task.Repeat != "" {
			nextDate, err := nextdate.NextDate(today, task.Date, task.Repeat)
			if err != nil {
				http.Error(w, `{"error":"Неправильное правило повторения"}`, http.StatusBadRequest)
				return
			}
			taskDate, err = time.Parse("2006-01-02", nextDate)

		}
	}
	task.Date = taskDate.Format(layout)

	switch req.Method {
	case http.MethodPost:
		Response, ResponseStatus, err = taskactions.AddTask(db, req)
		defer db.Close()
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(Response)
}
