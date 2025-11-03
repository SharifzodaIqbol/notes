// Package store содержит общие утилиты для работы с БД и HTTP.
package store

import (
	"database/sql"
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
	PageTitle   string
	Notes       []ViewData
	FilterValue string
}

func AddNote(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		title := r.FormValue("title")
		content := r.FormValue("content")
		_, err := db.Exec("INSERT INTO notes (title, content, updated_at) VALUES ($1, $2, CURRENT_TIMESTAMP)", title, content)
		if err != nil {
			log.Println("Error adding note:", err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditNote(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		idStr := r.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid note ID", http.StatusBadRequest)
			return
		}
		notes, err := fetchNotes(db, "")
		if err != nil {
			http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		for i := range notes {
			notes[i].Editing = notes[i].ID == id
		}

		pageData := PageData{
			PageTitle: "Редактирование заметки",
			Notes:     notes,
		}

		tmpl, err := template.ParseFiles("templates/index.html")
		if err != nil {
			http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, pageData)
	}
}

func UpdateNote(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		idStr := r.PathValue("id")
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
		idStr := r.PathValue("id")
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

func fetchNotes(db *sql.DB, filter string) ([]ViewData, error) {
	query := "SELECT id, title, content, created_at, updated_at FROM notes"
	arg := []any{}

	if filter != "" {
		query += " WHERE lower(title) LIKE $1"
		arg = append(arg, "%"+strings.ToLower(filter)+"%")
	}

	query += " ORDER BY ID DESC"

	rows, err := db.Query(query, arg...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := []ViewData{}
	for rows.Next() {
		var note Notes
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, ViewData{
			ID:      note.ID,
			Title:   note.Title,
			Content: note.Content,
			Editing: false,
		})
	}
	return notes, nil
}
func FilterByName(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	filter := strings.TrimSpace(r.URL.Query().Get("filter"))

	notes, err := fetchNotes(db, filter)
	if err != nil {
		http.Error(w, "Database query error: "+err.Error(), http.StatusInternalServerError)
		log.Println(err)
		return
	}

	pageData := PageData{
		PageTitle:   "Мои заметки",
		Notes:       notes,
		FilterValue: filter,
	}
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, pageData)
}
