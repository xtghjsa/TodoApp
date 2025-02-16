package handler

import (
	"log"
	"main/internal/api/request"
	"main/internal/api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MarkTaskDoneHandler struct {
	Usecase *usecase.MarkTaskDone
}

// MarkTaskDone Change task status to done
func (h *MarkTaskDoneHandler) MarkTaskDone(c *gin.Context) {
	var taskDone request.MarkTaskDoneID
	if err := c.ShouldBind(&taskDone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "wrong user"})
		return
	}

	err := h.Usecase.Execute(taskDone, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't mark task done"})
		log.Printf("error marking task done: %v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": "marked done"})
}
