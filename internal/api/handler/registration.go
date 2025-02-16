package handler

import (
	"log"
	"main/internal/api/request"
	"main/internal/api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegUserHandler struct {
	Usecase *usecase.RegisterUser
}

// RegUser register new user
func (h *RegUserHandler) RegisterUser(c *gin.Context) {
	var regData request.RegistrationData

	if err := c.ShouldBindJSON(&regData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login data"})
		return
	}
	if regData.Username == "" || regData.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login data"})
		return
	}

	err := h.Usecase.Execute(regData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't register user"})
		log.Printf("error registering user: %v\n", err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}
