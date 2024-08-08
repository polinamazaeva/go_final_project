package task

// Task представляет собой структуру для хранения информации о задаче
type Task struct {
	Id      string `json:"id"`      // Идентификатор задачи (строка). Поле сериализуется в JSON как "id"
	Date    string `json:"date"`    // Дата задачи в формате строки. Поле сериализуется в JSON как "date"
	Title   string `json:"title"`   // Название задачи. Поле сериализуется в JSON как "title"
	Comment string `json:"comment"` // Комментарий к задаче. Поле сериализуется в JSON как "comment"
	Repeat  string `json:"repeat"`  // Периодичность повторения задачи. Поле сериализуется в JSON как "repeat"
}
