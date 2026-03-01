package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		updateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJSON(w, map[string]string{"error": "метод не поддерживается"}, http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка чтения запроса"}, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON"}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	date, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "дата в неверном формате"}, http.StatusBadRequest)
		return
	}

	if task.Repeat != "" {
		_, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": "неверное правило повторения"}, http.StatusBadRequest)
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
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{"id": id}, http.StatusOK)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}

	writeJSON(w, task, http.StatusOK)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка чтения запроса"}, http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка десериализации JSON"}, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	_, err = db.GetTask(task.ID)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	date, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeJSON(w, map[string]string{"error": "дата в неверном формате"}, http.StatusBadRequest)
		return
	}

	if task.Repeat != "" {
		_, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": "неверное правило повторения"}, http.StatusBadRequest)
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

	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("ошибка записи JSON ответа: %v", err)
	}
}

func doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
	} else {
		now := time.Now()
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
		task.Date = next
		err = db.UpdateTask(task)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	_, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"}, http.StatusNotFound)
		return
	}

	err = db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]interface{}{}, http.StatusOK)
}
