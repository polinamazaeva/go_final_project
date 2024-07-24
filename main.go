package main

import (
	"net/http"

	"go_final_project/db"
	"go_final_project/handler"

	_ "modernc.org/sqlite"
)

func main() {
	db.CheckOpenCloseDb()
	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/api/nextdate", handler.NextDateHandler)
	http.HandleFunc("/api/task", handler.TaskHandler)
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}
}
