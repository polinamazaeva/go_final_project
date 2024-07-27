package check

import (
	"bytes"
	"encoding/json"
	"errors"
	"go_final_project/nextdate"
	"go_final_project/task"

	"net/http"
	"time"
)

// CheckTask проверяет и обрабатывает задачу из HTTP-запроса.
func Check(req *http.Request) (task.Task, int, error) {
	var task task.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return task, http.StatusInternalServerError, err
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		return task, http.StatusInternalServerError, err
	}

	if task.Title == "" {
		return task, http.StatusBadRequest, errors.New(`{"error":"task title is not specified"}`)
	}

	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	if task.Date == "" {
		task.Date = now.Format("20060102")
	}

	dateParse, err := time.Parse("20060102", task.Date)
	if err != nil {
		return task, http.StatusBadRequest, errors.New(`{"error":"incorrect date"}`)
	}

	var dateNew string
	if task.Repeat != "" {
		dateNew, err = nextdate.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return task, http.StatusBadRequest, err
		}
	}

	if task.Date == now.Format("20060102") {
		task.Date = now.Format("20060102")
	}

	if dateParse.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			task.Date = dateNew
		}
	}

	return task, http.StatusOK, nil
}
