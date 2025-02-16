package handler

import (
	"log"
	"main/internal/api/request"
	"main/internal/api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddTaskHandler struct {
	Usecase *usecase.AddTask
}

// AddTask Add task to database
func (h *AddTaskHandler) AddTask(c *gin.Context) {
	var addTask request.AddTask
	if err := c.ShouldBind(&addTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task data"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user"})
		return
	}

	err := h.Usecase.Execute(addTask, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can't add task"})
		log.Printf("error adding task: %v\n", err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"task": "added"})
}
