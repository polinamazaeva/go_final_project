package taskactions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/task"
)

func UpdateTask(db *sql.DB, t task.Task) ([]byte, int, error) {

	res, err := db.Exec(`UPDATE scheduler SET
	date = :date, title = :title, comment = :comment, repeat = :repeat
	WHERE id = :id`,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.Id))
	if err != nil {
		return []byte{}, 500, fmt.Errorf(`{"error":"task is not found" %s}`, err)
	}

	result, err := res.RowsAffected()
	if err != nil {
		return []byte{}, 500, fmt.Errorf(`{"error":"task is not found" %s}`, err)
	}
	if result == 0 {
		return []byte{}, 400, fmt.Errorf(`{"error":"task is not found"}`)
	}
	var str task.Task
	response, err := json.Marshal(str)
	if err != nil {
		return []byte{}, 500, err
	}

	return response, 200, nil
}
