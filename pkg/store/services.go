// Package store содержит общие утилиты для работы с БД и HTTP.
package store

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Notes struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
type ViewData struct {
	Title   string
	Content string
}
type PageData struct {
	PageTitle string
	Notes     []ViewData
}

func GetAllNotes(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
		return
	}

	rows, err := db.Query("SELECT * FROM notes")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Database query error"})
		return
	}
	defer rows.Close()
	notes := []ViewData{}
	for rows.Next() {
		var note Notes
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Error iterating rows" + err.Error()})
			return
		}
		data := ViewData{
			Title:   note.Title,
			Content: note.Content,
		}
		notes = append(notes, data)
	}
	pageData := PageData{
		PageTitle: "Мои заметки",
		Notes:     notes,
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, pageData)
	w.WriteHeader(http.StatusOK)
}
func AddNote(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		_, err := db.Exec("INSERT INTO notes (title, content) VALUES ($1, $2)", title, content)
		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
