package repository

import (
	"database/sql"
	"main/internal/api/request"
	"main/internal/api/response"
	"time"
)

type AddTaskInterface interface {
	AddTask(task request.AddTask, userID any) error
}

// AddTask Add task into database
func (db *PostgresRepository) AddTask(task request.AddTask, userID any) error {
	parsedDate, err := time.Parse("20060102", task.Date)
	parsedDateStr := parsedDate.Format("20060102")
	if err != nil {
		return err
	}
	_, err = db.DB.Exec("INSERT INTO tasks (title, description, date, user_id) VALUES ($1, $2, $3, $4)",
		task.Title, task.Description, parsedDateStr, userID)
	if err != nil {
		return err
	}
	return nil
}

type DeleteTaskInterface interface {
	DeleteTask(delete request.DeleteTaskID, userID any) error
}

// DeleteTask Delete task from database
func (db *PostgresRepository) DeleteTask(delete request.DeleteTaskID, userID any) error {
	_, err := db.DB.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", delete.ID, userID)
	if err != nil {
		return err
	}
	return nil
}

type MarkTaskDoneInterface interface {
	MarkTaskDone(done request.MarkTaskDoneID, userID any) error
}

// MarkTaskDone Mark task as done
func (db *PostgresRepository) MarkTaskDone(done request.MarkTaskDoneID, userID any) error {
	status := "done"
	_, err := db.DB.Exec("UPDATE tasks SET status = $1 WHERE id = $2 AND user_id = $3",
		status, done.ID, userID)
	if err != nil {
		return nil
	}
	return nil
}

type UpdateTaskInterface interface {
	UpdateTask(task request.UpdateTask, userID any) error
}

// UpdateTask Update task in database
func (db *PostgresRepository) UpdateTask(task request.UpdateTask, userID any) error {
	parsedDate, err := time.Parse("20060102", task.Date)
	parsedDateStr := parsedDate.Format("20060102")
	if err != nil {
		return err
	}
	_, err = db.DB.Exec("UPDATE tasks SET title = $1, description = $2, date = $3 WHERE id = $4 AND user_id = $5",
		task.Title, task.Description, parsedDateStr, task.ID, userID)
	if err != nil {
		return err
	}
	return nil

}

type ShowTasksInterface interface {
	ShowTasks(query, queryKey string, userID any) (response.ShowTasks, error)
}

// ShowTasks Get tasks from db
func (db *PostgresRepository) ShowTasks(query, queryKey string, userID any) (response.ShowTasks, error) {
	var taskResponse response.ShowTasks
	var rows *sql.Rows
	var err error
	switch queryKey {
	case "date":
		parsedDate, err := time.Parse("20060102", query)
		parsedDateStr := parsedDate.Format("20060102")
		if err != nil {
			return taskResponse, err
		}
		rows, err = db.DB.Query("SELECT id, title, description, date, status FROM tasks WHERE user_id = $1 AND date = $2", userID, parsedDateStr)
		if err != nil {
			return taskResponse, err
		}
	case "status":
		rows, err = db.DB.Query("SELECT id, title, description, date, status FROM tasks WHERE user_id = $1 AND status = $2", userID, query)
		if err != nil {
			return taskResponse, err
		}
	default:
		rows, err = db.DB.Query("SELECT id, title, description, date, status FROM tasks WHERE user_id = $1", userID)
		if err != nil {
			return taskResponse, err
		}
	}
	defer rows.Close()

	for rows.Next() {
		var task response.ShowTasksItem
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Date, &task.Status)
		if err != nil {
			return taskResponse, err
		}
		taskResponse.Tasks = append(taskResponse.Tasks, task)
	}
	return taskResponse, err
}
