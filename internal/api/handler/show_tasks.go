package handler

import (
	"log"
	"net/http"

	"main/internal/api/usecase"

	"github.com/gin-gonic/gin"
)

type ShowTasksHandler struct {
	Usecase *usecase.ShowTasks
}

// ShowTasks Show tasks depending on query
func (h *ShowTasksHandler) ShowTasks(c *gin.Context) {
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
	taskResponse, err := h.Usecase.Execute(query, queryKey, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't show tasks"})
		log.Printf("error showing tasks: %v\n", err)
		return
	}
	c.JSON(http.StatusOK, taskResponse)
}
