package handlers

import (
	"log"
	"main/internal/db"
	"main/internal/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db.DbConn
}

// POST
func (h *Handler) AddTask(c *gin.Context) {
	var addTask types.TaskAdd
	if err := c.ShouldBind(&addTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task data"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		log.Println(userID)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user"})
		return
	}
	userIDStr := userID.(string)
	err := h.AddTaskDb(addTask.Title, addTask.Description, addTask.Date, userIDStr)
	if err != nil {
		log.Panicf("error adding task: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't add task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"task": "added"})
}

// DEL
func (h *Handler) DelTask(c *gin.Context) {
	var delTask types.TaskId
	if err := c.ShouldBind(&delTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user"})
		return
	}
	taskID, err := strconv.Atoi(delTask.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong id format"})
		return
	}
	userIDStr := userID.(string)
	err = h.DelTaskDb(taskID, userIDStr)
	if err != nil {
		log.Panicf("error deleting task: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't delete task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": "deleted"})
}

// PUT
func (h *Handler) UpdTask(c *gin.Context) {
	var updTask types.Task
	if err := c.ShouldBind(&updTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task udpate data"})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user"})
		return
	}
	taskId, err := strconv.Atoi(updTask.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong id format"})
		return
	}
	userIDStr := userID.(string)
	err = h.UpdTaskDb(updTask.Title, updTask.Description, updTask.Date, userIDStr, taskId)
	if err != nil {
		log.Panicf("error updating task: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't update task"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": "updated"})
}

// PUT
func (h *Handler) MarkTaskDone(c *gin.Context) {
	var taskDone types.TaskId
	if err := c.ShouldBind(&taskDone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user"})
		return
	}
	taskId, err := strconv.Atoi(taskDone.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wrong id format"})
		return
	}
	userIDStr := userID.(string)
	err = h.MarkTaskDoneDb(taskId, userIDStr)
	if err != nil {
		log.Panicf("error marking task done: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't mark task done"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": "marked done"})
}

// GET
func (h *Handler) ShowTasks(c *gin.Context) {
	//Search by date
	requiredDate := c.Query("date")
	//Search by done/notDone
	requiredStatus := c.Query("status")
	var query string
	var queryKey string

	if requiredDate != "" {
		queryKey = "date"
		query = requiredDate
	} else if requiredStatus != "" {
		queryKey = "status"
		query = requiredStatus
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user"})
		return
	}

	userIDStr := userID.(string)
	taskResponse, err := h.ShowTasksDb(query, queryKey, userIDStr)
	if err != nil {
		log.Panicf("error showing tasks: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't show tasks"})
		return
	}
	c.JSON(http.StatusOK, taskResponse)
}
