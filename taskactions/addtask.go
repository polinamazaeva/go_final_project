package taskactions

import (
	"database/sql"
	"encoding/json"
	"go_final_project/task"
)

type Id struct {
	Id int64 `json:"id"`
}

func AddTask(db *sql.DB, task task.Task) ([]byte, int, error) {
	var id Id

	// Используем параметризированный запрос для вставки
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

	// Получение последнего вставленного ID
	getId, err := res.LastInsertId()
	if err != nil {
		return []byte{}, 500, err
	}

	id.Id = getId

	// Преобразование ID в JSON
	idResult, err := json.Marshal(id)
	if err != nil {
		return []byte{}, 500, err
	}
	return idResult, 200, nil
}
