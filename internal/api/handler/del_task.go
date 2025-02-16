package handler

import (
	"log"
	"main/internal/api/request"
	"main/internal/api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeleteTaskHandler struct {
	Usecase *usecase.DeleteTask
}

// DelTask Delete task from database
func (h *DeleteTaskHandler) DelTask(c *gin.Context) {
	var delTask request.DeleteTaskID
	if err := c.ShouldBind(&delTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user"})
		return
	}

	err := h.Usecase.Execute(delTask, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't delete task"})
		log.Printf("error deleting task: %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": "deleted"})
}
