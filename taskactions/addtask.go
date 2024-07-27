package taskactions

import (
	"database/sql"
	"encoding/json"
	"go_final_project/check" // Убедитесь, что имя пакета и импорт правильны
	"net/http"
)

type Id struct {
	Id int64 `json:"id"`
}

func AddTask(db *sql.DB, req *http.Request) ([]byte, int, error) {
	var idresp Id

	task, ResponseStatus, err := check.Check(req) // Убедитесь, что check.Check доступен и правильно вызывается
	if err != nil {
		return []byte{}, ResponseStatus, err
	}

	result, err := db.Exec(`INSERT INTO scheduler (date, title, comment, repeat)
		VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return []byte{}, 500, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return []byte{}, 500, err
	}

	idresp.Id = id

	idResult, err := json.Marshal(idresp)
	if err != nil {
		return []byte{}, 500, err
	}
	return idResult, 200, nil
}
