package taskactions

import (
	"database/sql"
	"encoding/json"
	"go_final_project/task"
	"net/http"
)

func GetTasks(db *sql.DB, limit int) ([]byte, int, error) {
	var tasks []task.Task

	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	for rows.Next() {
		var t task.Task
		err := rows.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		tasks = append(tasks, t)
	}

	err = rows.Err()
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	if tasks == nil {
		tasks = []task.Task{}
	}

	response, err := json.Marshal(map[string][]task.Task{"tasks": tasks})
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return response, http.StatusOK, nil
}

func GetTaskByID(db *sql.DB, id string) ([]byte, int, error) {
	var t task.Task
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := db.QueryRow(query, id)
	err := row.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, http.StatusNotFound, err
		}
		return nil, http.StatusInternalServerError, err
	}

	response, err := json.Marshal(t)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return response, http.StatusOK, nil
}
