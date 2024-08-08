package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"go_final_project/nextdate"
	"go_final_project/task"
)

// Константа для формата даты, используемого в приложении
const DateFormat = "20060102"

// Функция Check выполняет проверку и валидацию задачи, полученной из HTTP-запроса.
// Возвращает валидированную задачу, HTTP-статус и возможную ошибку.
func Check(req *http.Request) (task.Task, int, error) {
	var t task.Task      // Переменная для хранения задачи
	var buf bytes.Buffer // Буфер для хранения данных из тела запроса

	// Чтение данных из тела запроса в буфер
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		return t, http.StatusInternalServerError, err // Возвращаем ошибку, если чтение не удалось
	}

	// Декодирование JSON из буфера в структуру задачи
	if err = json.Unmarshal(buf.Bytes(), &t); err != nil {
		return t, http.StatusInternalServerError, err // Возвращаем ошибку, если декодирование не удалось
	}

	// Проверка наличия обязательного поля title (название задачи)
	if t.Title == "" {
		return t, http.StatusBadRequest, errors.New(`{"error":"task title is not specified"}`) // Возвращаем ошибку, если поле пустое
	}

	now := time.Now()
	// Приводим время к полуночи, чтобы сравнение дат было корректным
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	// Если дата не указана, устанавливаем текущую дату
	if t.Date == "" {
		t.Date = now.Format(DateFormat)
	}

	// Парсим дату задачи в формате DateFormat
	dateParse, err := time.Parse(DateFormat, t.Date)
	if err != nil {
		return t, http.StatusBadRequest, errors.New(`{"error":"incorrect date"}`) // Возвращаем ошибку, если дата некорректна
	}

	var dateNew string
	// Если указано поле repeat, проверяем его корректность и вычисляем следующую дату
	if t.Repeat != "" {
		dateNew, err = nextdate.NextDate(now, t.Date, t.Repeat)
		if err != nil {
			return t, http.StatusBadRequest, err // Возвращаем ошибку, если вычисление следующей даты не удалось
		}
	}

	// Если дата задачи совпадает с текущей датой, оставляем её неизменной
	if t.Date == now.Format(DateFormat) {
		t.Date = now.Format(DateFormat)
	}

	// Если дата задачи раньше текущей
	if dateParse.Before(now) {
		// Если поле repeat пустое, устанавливаем текущую дату
		if t.Repeat == "" {
			t.Date = now.Format(DateFormat)
		} else {
			// Иначе устанавливаем следующую дату повторения
			t.Date = dateNew
		}
	}

	// Возвращаем валидированную задачу, HTTP-статус 200 OK и nil, что означает отсутствие ошибок
	return t, http.StatusOK, nil
}
