package taskactions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/task"
	"net/http"
)

func TaskID(db *sql.DB, id string) ([]byte, int, error) {
	var task task.Task

	row := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return []byte{}, http.StatusNotFound, fmt.Errorf(`{"error":"task not found"}`)
		}
		return []byte{}, http.StatusInternalServerError, fmt.Errorf(`{"error":"error reading data: %v"}`, err)
	}

	result, err := json.Marshal(task)
	if err != nil {
		return []byte{}, http.StatusInternalServerError, err
	}

	return result, http.StatusOK, nil
}
