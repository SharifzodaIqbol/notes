package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"notesApp/pkg/store"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func initDB() error {
	connStr := fmt.Sprintf("user=postgres password=%s dbname=notes sslmode=disable", os.Getenv("mypass"))
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}
	return db.Ping()
}

func main() {
	loadEnv()
	if err := initDB(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		store.FilterByName(db, w, r)
	})
	mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
		store.AddNote(db, w, r)
	})
	mux.HandleFunc("/edit/{id}", func(w http.ResponseWriter, r *http.Request) {
		store.EditNote(db, w, r)
	})
	mux.HandleFunc("/update/{id}", func(w http.ResponseWriter, r *http.Request) {
		store.UpdateNote(db, w, r)
	})
	mux.HandleFunc("/delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		store.DeleteNote(db, w, r)
	})
	fmt.Println("Сервер запущен!")
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Println(err)
		return
	}
}
