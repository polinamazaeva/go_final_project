package nextdate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Функция для преобразования строки в формат времени
func parseTime(s string) (time.Time, error) {
	t, err := time.Parse("20060102", s)
	return t, err
}

// Функция для вычисления следующей даты в соответствии с правилами повторения
func NextDate(now time.Time, date string, repeat string) (string, error) {

	// Преобразование входных параметров во время
	nowTime := now //Truncate(24 * time.Hour)
	dateTime, err := parseTime(date)
	if err != nil {
		return "", fmt.Errorf("некорректная дата date: %w", err)
	}

	symbols := strings.Split(repeat, " ")
	if len(symbols) == 0 {
		return "", fmt.Errorf("invalid repeat format")
	}
	firstsymbol := symbols[0]

	if !strings.ContainsAny(firstsymbol, "dywm") {
		return "", fmt.Errorf(`{"error":"incorrect symbol %s"}`, firstsymbol)
	}

	switch firstsymbol {
	case "y": // Ежегодно
		if len(symbols) != 1 {
			return "", errors.New("invalid repeat format for 'y'")
		}
		return nextYearlyDate(dateTime, nowTime)

	case "d": // Через указанное число дней
		if len(symbols) != 2 {
			return "", fmt.Errorf("некорректный интервал дней в правиле повторения: %w", err)
		}
		secondsymbol := symbols[1]
		return nextDayRepeat(dateTime, nowTime, secondsymbol)

	default: // Неподдерживаемые форматы
		return "", errors.New("unsupported repeat format")
	}
}

// nextYearlyDate переносит дату на один год вперед
func nextYearlyDate(dateTime time.Time, nowTime time.Time) (string, error) {
	next := dateTime

	for {
		if next.After(nowTime) && next.After(dateTime) {
			break
		}
		next = next.AddDate(1, 0, 0)
	}

	return next.Format("20060102"), nil
}

// nextDayRepeat вычисляет следующую дату на основании повторения через указанное количество дней
func nextDayRepeat(dateTime time.Time, nowTime time.Time, daysInterval string) (string, error) {
	daysIntervalInt, err := strconv.Atoi(daysInterval)
	if err != nil {
		return "", fmt.Errorf("некорректный интервал дней в правиле повторения: %w", err)
	}

	if daysIntervalInt > 400 {
		return "", fmt.Errorf("максимальный интервал дней равен 400")
	}

	nextDate := dateTime.AddDate(0, 0, daysIntervalInt)

	for {
		if nextDate.After(nowTime) {
			break
		}
		nextDate = nextDate.AddDate(0, 0, daysIntervalInt)
	}
	return nextDate.Format("20060102"), nil
}
