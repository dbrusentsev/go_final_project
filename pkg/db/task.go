package db

import "fmt"

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func GetTask(id string) (*Task, error) {
	task := Task{}
	row := DB.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?", id)
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	return &task, err
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	_, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	return err
}

func DeleteTask(id string) error {
	_, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	return err
}

func Tasks(limit int) ([]*Task, error) {
	var tasks []*Task = []*Task{}

	rows, err := DB.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id      int64
			date    string
			title   string
			comment string
			repeat  string
		)

		err := rows.Scan(&id, &date, &title, &comment, &repeat)
		if err != nil {
			return nil, fmt.Errorf("error while scanning row: %w", err)
		}

		tasks = append(tasks, &Task{
			ID:      fmt.Sprintf("%d", id),
			Date:    date,
			Title:   title,
			Comment: comment,
			Repeat:  repeat,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
