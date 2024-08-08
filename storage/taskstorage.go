package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"go_final_project/task"
)

// DefaultLimit определяет значение лимита по умолчанию для запросов задач
const DefaultLimit = 25

// TaskStorage представляет собой структуру, которая содержит соединение с базой данных
type TaskStorage struct {
	DB *sql.DB // Указатель на объект базы данных
}

// TaskID получает задачу по её ID из базы данных
func (ts *TaskStorage) TaskID(id string) (task.Task, error) {
	var task task.Task // Переменная для хранения задачи

	// Выполняем SQL-запрос для получения задачи по ID
	row := ts.DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	// Заполняем поля задачи из результата запроса
	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		// Если задача не найдена, возвращаем ошибку с сообщением "task not found"
		if err == sql.ErrNoRows {
			return task, errors.New(`{"error":"not find the task"}`)
		}
		// Возвращаем ошибку в случае других проблем
		return task, err
	}

	// Если возникла ошибка при работе с результатом запроса, возвращаем её
	if err := row.Err(); err != nil {
		return task, err
	}

	// Возвращаем задачу
	return task, nil
}

// AddTask добавляет новую задачу в базу данных
func (ts *TaskStorage) AddTask(t task.Task) (int64, error) {
	var id int64 // Переменная для хранения ID новой задачи

	// Выполняем SQL-запрос для добавления новой задачи в базу данных
	result, err := ts.DB.Exec(`INSERT INTO scheduler (date, title, comment, repeat)
        VALUES (:date, :title, :comment, :repeat)`,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	)
	if err != nil {
		return id, err // Возвращаем ошибку в случае неудачи
	}

	// Получаем ID добавленной задачи
	id, err = result.LastInsertId()
	if err != nil {
		return id, err // Возвращаем ошибку в случае неудачи
	}

	// Возвращаем ID новой задачи
	return id, nil
}

// UpdateTask обновляет существующую задачу по её ID
func (ts *TaskStorage) UpdateTask(t task.Task) error {
	// Выполняем SQL-запрос для обновления задачи
	res, err := ts.DB.Exec(`UPDATE scheduler SET
	date = :date, title = :title, comment = :comment, repeat = :repeat
	WHERE id = :id`,
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.Id))
	if err != nil {
		return err // Возвращаем ошибку в случае неудачи
	}

	// Проверяем, было ли обновлено хотя бы одно поле
	result, err := res.RowsAffected()
	if err != nil {
		return err // Возвращаем ошибку в случае неудачи
	}
	if result == 0 {
		return errors.New(`{"error":"not found the task"}`) // Возвращаем ошибку, если задача не найдена
	}

	// Возвращаем nil, если обновление прошло успешно
	return nil
}

// DeleteTask удаляет задачу по её ID
func (ts *TaskStorage) DeleteTask(id string) error {
	// Выполняем SQL-запрос для удаления задачи по ID
	result, err := ts.DB.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("internal server error: %w", err) // Возвращаем ошибку при неудаче
	}

	// Проверяем, была ли удалена хотя бы одна задача
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("internal server error: %w", err) // Возвращаем ошибку при неудаче
	}

	// Если ни одна задача не была удалена, возвращаем ошибку
	if rowsAffected == 0 {
		return errors.New("task not found")
	}

	// Возвращаем nil, если удаление прошло успешно
	return nil
}

// TaskDone возвращает задачу по её ID
func (ts *TaskStorage) TaskDone(id string) (task.Task, error) {
	var task task.Task // Переменная для хранения задачи

	// Выполняем SQL-запрос для получения задачи по ID
	row := ts.DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id))

	// Заполняем поля задачи из результата запроса
	err := row.Scan(&task.Id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		// Если задача не найдена, возвращаем ошибку с сообщением "task not found"
		if err == sql.ErrNoRows {
			return task, errors.New(`{"error":"not find the task"}`)
		}
		// Возвращаем ошибку в случае других проблем
		return task, err
	}

	// Если возникла ошибка при работе с результатом запроса, возвращаем её
	if err := row.Err(); err != nil {
		return task, err
	}

	// Возвращаем задачу
	return task, nil
}

// GetTasks получает список задач с ограничением на количество
func (ts *TaskStorage) GetTasks(limit int) ([]task.Task, error) {
	var tasks []task.Task // Срез для хранения списка задач

	// Выполняем SQL-запрос для получения списка задач с ограничением на количество
	query := "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?"
	rows, err := ts.DB.Query(query, limit)
	if err != nil {
		return nil, err // Возвращаем ошибку в случае неудачи
	}
	defer rows.Close()

	// Итерируем по результатам запроса и добавляем задачи в срез
	for rows.Next() {
		var t task.Task
		if err := rows.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err // Возвращаем ошибку в случае проблем с получением данных
		}
		tasks = append(tasks, t) // Добавляем задачу в список
	}

	// Проверяем наличие ошибок при итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, err // Возвращаем ошибку в случае проблем с итерацией
	}

	// Если задач нет, возвращаем пустой список
	if len(tasks) == 0 {
		tasks = []task.Task{}
	}

	// Возвращаем список задач
	return tasks, nil
}

// GetTaskByID получает задачу по её ID
func (ts *TaskStorage) GetTaskByID(id string) (task.Task, error) {
	var t task.Task // Переменная для хранения задачи

	// Выполняем SQL-запрос для получения задачи по ID
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	row := ts.DB.QueryRow(query, id)

	// Заполняем поля задачи из результата запроса
	err := row.Scan(&t.Id, &t.Date, &t.Title, &t.Comment, &t.Repeat)
	if err != nil {
		// Если задача не найдена, возвращаем ошибку с сообщением "task not found"
		if err == sql.ErrNoRows {
			return t, errors.New(`{"error":"task not found"}`)
		}
		// Возвращаем ошибку в случае других проблем
		return t, err
	}

	// Возвращаем задачу
	return t, nil
}
