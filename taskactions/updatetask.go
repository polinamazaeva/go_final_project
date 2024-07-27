package taskactions

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go_final_project/check" // Убедитесь, что имя пакета и импорт правильны
	"go_final_project/task"  // Проверьте, нужен ли этот импорт
	"net/http"
)

func UptadeTaskID(db *sql.DB, req *http.Request) ([]byte, int, error) {

	taskID, ResponseStatus, err := check.Check(req) // Убедитесь, что check.Check доступен и правильно вызывается
	if err != nil {
		return []byte{}, ResponseStatus, err
	}

	res, err := db.Exec(`UPDATE scheduler SET
	date = :date, title = :title, comment = :comment, repeat = :repeat
	WHERE id = :id`,
		sql.Named("date", taskID.Date),
		sql.Named("title", taskID.Title),
		sql.Named("comment", taskID.Comment),
		sql.Named("repeat", taskID.Repeat),
		sql.Named("id", taskID.Id))
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
