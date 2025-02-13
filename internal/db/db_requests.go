package db

import (
	"database/sql"
	"main/internal/types"
	"time"
)

// Add task into database
func (db *DbConn) AddTaskDb(title, description, date, userID string) error {
	parsedDate, err := time.Parse("20060102", date)
	parsedDateStr := parsedDate.Format("20060102")
	if err != nil {
		return err
	}
	_, err = db.DB.Exec("INSERT INTO tasks (title, description, date, user_id) VALUES ($1, $2, $3, $4)",
		title, description, parsedDateStr, userID)
	if err != nil {
		return err
	}
	return nil
}

// Delete task from database
func (db *DbConn) DelTaskDb(id int, userID string) error {
	_, err := db.DB.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		return err
	}
	return nil
}

// Update task in database
func (db *DbConn) UpdTaskDb(title, description, date, userID string, id int) error {
	parsedDate, err := time.Parse("20060102", date)
	parsedDateStr := parsedDate.Format("20060102")
	if err != nil {
		return err
	}
	_, err = db.DB.Exec("UPDATE tasks SET title = $1, description = $2, date = $3 WHERE id = $4 AND user_id = $5",
		title, description, parsedDateStr, id, userID)
	if err != nil {
		return err
	}
	return nil

}

// Mark task as done
func (db *DbConn) MarkTaskDoneDb(id int, userID string) error {
	status := "done"
	_, err := db.DB.Exec("UPDATE tasks SET status = $1 WHERE id = $2 AND user_id = $3",
		status, id, userID)
	if err != nil {
		return nil
	}
	return nil
}

// Get tasks from db
func (db *DbConn) ShowTasksDb(query, queryKey, userID string) (types.ShowTasks, error) {
	var taskResponse types.ShowTasks
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
		var task types.Task
		err := rows.Scan(&task.Id, &task.Title, &task.Description, &task.Date, &task.Status)
		if err != nil {
			return taskResponse, err
		}
		taskResponse.Tasks = append(taskResponse.Tasks, task)
	}
	return taskResponse, err
}

// Authorization
func (db *DbConn) RegUserDb(username, hashedPassword string) error {
	_, err := db.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (db *DbConn) CheckUser(username string) (userID int, storedHashedPass string, err error) {
	var userIdFound int
	var storedPassFound string
	err = db.DB.QueryRow("SELECT id, password FROM users WHERE username=$1", username).Scan(&userIdFound, &storedPassFound)
	if err != nil {
		return 0, "", err
	}
	return userIdFound, storedPassFound, nil
}
