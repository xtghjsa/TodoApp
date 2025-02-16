package handler

import (
	"log"
	"main/internal/api/request"
	"main/internal/api/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginHandler struct {
	Usecase *usecase.Login
}

// Login
func (h *LoginHandler) Login(c *gin.Context) {
	var creds request.LoginData
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login data"})
		return
	}

	token, err := h.Usecase.Execute(creds)
	if err != nil {
		log.Printf("Failed to create session token for user: %s\n", creds.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid login data"})
		return
	}

	log.Printf("Created session token for user: %s, token: %s\n", creds.Username, token)
	c.JSON(http.StatusOK, gin.H{"token": token})
}
