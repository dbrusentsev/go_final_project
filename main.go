package main

import (
	"fmt"
	"net/http"
	"os"

	"go_final_project/pkg/db"
)

func main() {

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	err := db.Init(dbFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.DB.Close()


	http.Handle("/", http.FileServer(http.Dir("web")))


	fmt.Println("Сервер запущен на порту 7540")
	err = http.ListenAndServe(":7540", nil)
	if err != nil {
		fmt.Println(err)
	}
}
