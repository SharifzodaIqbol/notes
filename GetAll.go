package main

import (
	"encoding/json"
	"net/http"
)

type Notes struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func GetAllNotes(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Database not init"})
		return
	}

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
	
	var notes []Notes
	for rows.Next() {
		var note Notes
		err := rows.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
		if err != nil {
			json.NewEncoder(w).Encode(map[string]string{"error": "Error iterating rows" + err.Error()})
			return
		}
		notes = append(notes, note)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(notes); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}