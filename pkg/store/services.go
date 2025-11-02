// Package store содержит общие утилиты для работы с БД и HTTP.
package store

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	ID      int
	Title   string
	Content string
	Editing bool
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

	rows, err := db.Query("SELECT * FROM notes ORDER BY id DESC")
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
			ID:      note.ID,
			Title:   note.Title,
			Content: note.Content,
			Editing: false,
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

func EditNote(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		idStr := r.URL.Path[len("/edit/"):]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT * FROM notes ORDER BY id DESC")
		if err != nil {
			http.Error(w, "Database query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		notes := []ViewData{}
		for rows.Next() {
			var note Notes
			err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
			if err != nil {
				http.Error(w, "Error iterating rows", http.StatusInternalServerError)
				return
			}

			// Помечаем редактируемую заметку
			editing := note.ID == id
			data := ViewData{
				ID:      note.ID,
				Title:   note.Title,
				Content: note.Content,
				Editing: editing,
			}
			notes = append(notes, data)
		}

		pageData := PageData{
			PageTitle: "Редактирование заметки",
			Notes:     notes,
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, pageData)
	}
}

func UpdateNote(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idStr := strings.TrimPrefix(r.URL.Path, "/delete/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")

		_, err = db.Exec("UPDATE notes SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3",
			title, content, id)
		if err != nil {
			log.Println("Error updating note:", err)
			http.Error(w, "Error updating note", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func DeleteNote(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		idStr := strings.TrimPrefix(r.URL.Path, "/delete/")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM notes WHERE id = $1", id)
		if err != nil {
			log.Println("Error deleting note:", err)
			http.Error(w, "Error deleting note", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
