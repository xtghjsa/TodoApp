package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	t "main/internal/types"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

type RedisConn struct {
	Connect *redis.Client
	Ctx     *context.Context
}

type LoginStruct struct {
	Db    Handler
	Redis RedisConn
}

// RegUser register new user
func (h *Handler) RegUser(c *gin.Context) {
	var login t.LoginData

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login data"})
		return
	}
	if login.Username == "" || login.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login data"})
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't hash password"})
	}
	err = h.RegUserDb(login.Username, string(hashedPass))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't register user"})
	}
	c.JSON(http.StatusCreated, gin.H{"status": "created"})
}

func (ls *LoginStruct) Login(c *gin.Context) {
	var login t.LoginData
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login data"})
		return
	}

	userID, storedHashedPass, err := ls.Db.CheckUser(login.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid login data"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't not check data"})
		log.Panicf("error checking user data: %v", err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPass), []byte(login.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}

	sessionToken := fmt.Sprintf("%d", time.Now().UnixNano())
	err = ls.Redis.Connect.Set(*ls.Redis.Ctx, sessionToken, userID, 24*time.Hour).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't create authorized session"})
		log.Panicf("error creating authorized session: %v", err)
		return
	}
	log.Printf("Created session token for user: %s, token: %s\n", login.Username, sessionToken)
	c.JSON(http.StatusOK, gin.H{"token": sessionToken})

}

func AuthMw(rc RedisConn) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "need authorization header"})
			c.Abort()
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		if token == header {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "bearer token required"})
			c.Abort()
			return
		}
		userID, err := rc.Connect.Get(*rc.Ctx, token).Result()
		if errors.Is(err, redis.Nil) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token is invalid"})
			c.Abort()
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "can't check authorized session"})
			c.Abort()
			return
		}
		if userID == "" {
			fmt.Println("empty userID")
		}
		log.Printf("Authenticated user with ID: %v\n", userID)
		c.Set("userID", userID)
		c.Next()
	}
}
