package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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
	connStr := fmt.Sprintf("user=postgres password=%s dbname=note sslmode=disable", os.Getenv("mypass"))
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db.Ping()
}

func syaHi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi")
}

func main() {
	loadEnv()
	if err := initDB(); err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", syaHi)
	if err := http.ListenAndServe(":8082", mux); err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Сервер запущен!")
}
