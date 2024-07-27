package taskactions

import (
	"database/sql"
	"errors"
	"fmt"
)

func DeleteTask(db *sql.DB, id string) (int, error) {
	task, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return 500, fmt.Errorf(`{"error":"%s"}`, err)
	}

	rowsAffected, err := task.RowsAffected()
	if err != nil {
		return 500, err
	}

	if rowsAffected == 0 {
		return 400, errors.New(`{"error":"not found the task"}`)
	}
	return 200, nil
}
