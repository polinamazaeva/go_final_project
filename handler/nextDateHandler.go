package handler

import (
	"log"
	"net/http"
	"time"

	"go_final_project/nextdate"
)

// NextDateHandler обрабатывает HTTP-запросы для получения следующей даты на основе текущей даты, даты начала и повторяющегося интервала.
func NextDateHandler(w http.ResponseWriter, req *http.Request) {
	// Извлекаем параметры запроса
	queryParams := req.URL.Query()
	now := queryParams.Get("now")       // Текущая дата и время
	day := queryParams.Get("date")      // Дата начала
	repeat := queryParams.Get("repeat") // Интервал повторения

	// Парсим текущую дату из строки
	timeNow, err := time.Parse(nextdate.DateFormat, now)
	if err != nil {
		// Если произошла ошибка при парсинге, возвращаем ошибку 500 Internal Server Error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Вызываем функцию NextDate для вычисления следующей даты
	result, err := nextdate.NextDate(timeNow, day, repeat)
	if err != nil {
		// Если произошла ошибка при вычислении даты, возвращаем ошибку 400 Bad Request
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отправляем результат клиенту
	_, err = w.Write([]byte(result))
	if err != nil {
		// Логируем ошибку записи ответа
		log.Printf("error writing next date: %v", err)
	}
	// Устанавливаем заголовок ответа и отправляем результат
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
