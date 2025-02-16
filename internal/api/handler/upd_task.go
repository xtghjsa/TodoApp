package handler

import (
	"log"
	"main/internal/api/request"
	"main/internal/api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateTaskHandler struct {
	Usecase *usecase.UpdateTask
}

// UpdTask Update existing task in database
func (h *UpdateTaskHandler) UpdateTask(c *gin.Context) {
	var updTask request.UpdateTask
	if err := c.ShouldBind(&updTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task udpate data"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user"})
		return
	}

	err := h.Usecase.Execute(updTask, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't update task"})
		log.Printf("error updating task: %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": "updated"})
}
