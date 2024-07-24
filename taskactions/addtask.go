package taskactions

import (
	"database/sql"
	"encoding/json"
	"go_final_project/task"
	"net/http"
)

var Response []byte

type Id struct {
	Id int64 `json:"id"`
}

func AddTask(db *sql.DB, req *http.Request) ([]byte, int, error) {

	var id Id
	var task task.Task

	res, err := db.Exec(`INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return []byte{}, 500, err
	}

	getId, err := res.LastInsertId()
	if err != nil {
		return []byte{}, 500, err
	}

	id.Id = getId

	idResult, err := json.Marshal(id)
	if err != nil {
		return []byte{}, 500, err
	}
	return idResult, 200, nil
}
