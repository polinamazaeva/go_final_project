package main

import (
	"log"
	"net/http"

	"go_final_project/db"
	"go_final_project/handler"

	_ "modernc.org/sqlite"
)

func main() {
	database, err := db.CheckOpenCloseDb()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", handler.NextDateHandler)
	http.HandleFunc("/api/task", handler.MakeTaskHandler(database))
	http.HandleFunc("/api/tasks", handler.MakeTaskHandler(database))

	log.Println("Server starting on :7540")
	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
