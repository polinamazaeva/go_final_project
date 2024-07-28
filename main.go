package main

import (
	"go_final_project/db"
	"go_final_project/handler"
	"log"
	"net/http"

	_ "modernc.org/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	database, err := db.CheckOpenCloseDb()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/task", handler.TaskHandler)
	http.HandleFunc("/api/task/done", handler.TaskDoneHandler)
	http.HandleFunc("/api/tasks", handler.GetTasksHandler(database))
	http.HandleFunc("/api/nextdate", handler.NextDateHandler)

	log.Println("Server starting on :7540")
	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
