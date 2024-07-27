package handler

import (
	"database/sql"
	"encoding/json"
	"go_final_project/nextdate"
	"go_final_project/task"
	"go_final_project/taskactions"
	"log"
	"net/http"
	"time"
)

type ResponseForPostTask struct {
	Id int64 `json:"id"`
}

var ResponseStatus int

// TaskHandler возвращает обработчик для создания и обновления задач
func TaskHandler(w http.ResponseWriter, req *http.Request) {
	param := req.URL.Query().Get("id")

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		http.Error(w, "error opening database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var response []byte

	switch req.Method {
	case http.MethodGet:
		if param == "" {
			http.Error(w, `{"error":"incorrect id"}`, http.StatusBadRequest)
			return
		}
		response, ResponseStatus, err = taskactions.TaskID(db, param)
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}

	case http.MethodPost:
		response, ResponseStatus, err = taskactions.AddTask(db, req)
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}
	case http.MethodPut:
		response, ResponseStatus, err = taskactions.UptadeTaskID(db, req)
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}
	case http.MethodDelete:
		ResponseStatus, err = taskactions.DeleteTask(db, param)
		if err != nil {
			http.Error(w, err.Error(), ResponseStatus)
			return
		}
		response = []byte(`{}`)
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Обработчик для завершения задачи
func TaskDoneHandler(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"missing id parameter"}`, http.StatusBadRequest)
		return
	}

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		logError(err)
		http.Error(w, `{"error":"error opening database: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var taskID task.Task
	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err = row.Scan(&taskID.Id, &taskID.Date, &taskID.Title, &taskID.Comment, &taskID.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"task not found"}`, http.StatusNotFound)
		} else {
			http.Error(w, `{"error":"error scanning task: `+err.Error()+`"}`, http.StatusInternalServerError)
		}
		return
	}

	if taskID.Repeat == "" {
		_, err := db.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			logError(err)
			http.Error(w, `{"error":"error deleting task: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		nextDate, err := nextdate.NextDate(now, taskID.Date, taskID.Repeat)
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		_, err = db.Exec("UPDATE scheduler SET date = ? WHERE id = ?", nextDate, taskID.Id)
		if err != nil {
			http.Error(w, `{"error":"error updating task: `+err.Error()+`"}`, http.StatusInternalServerError)
			return
		}
	}

	response := map[string]string{
		"status": "success",
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, `{"error":"error marshalling response: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write(responseData)
}

func logError(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
}
