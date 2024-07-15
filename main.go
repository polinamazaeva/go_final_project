package main

import (
	"net/http"

	"go_final_project/db"

	_ "modernc.org/sqlite"
)

func main() {
	db.CheckOpenCloseDb()
	webDir := "./web"

	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err := http.ListenAndServe(":7540", nil)
	if err != nil {
		panic(err)
	}
}
