package api

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка чтения запроса"})
		return
	}

	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	date, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "дата в неверном формате"})
		return
	}

	if task.Repeat != "" {
		_, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": "неверное правило повторения"})
			return
		}
	}

	if date.Before(now) && date.Format(DateFormat) != now.Format(DateFormat) {
		if task.Repeat == "" {
			task.Date = now.Format(DateFormat)
		} else {
			next, _ := NextDate(now, task.Date, task.Repeat)
			task.Date = next
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]interface{}{"id": id})
}

func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}
