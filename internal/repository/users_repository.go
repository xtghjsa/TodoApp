package repository

import (
	"context"
	"fmt"
	"main/internal/api/request"
	"time"

	"github.com/redis/go-redis/v9"

	"golang.org/x/crypto/bcrypt"
)

type RegisterUserInterface interface {
	RegisterUser(regData request.RegistrationData) error
}

// RegisterUser adds new user data to database `users` table
func (db *PostgresRepository) RegisterUser(regData request.RegistrationData) error {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(regData.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = db.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", regData.Username, hashedPass)
	if err != nil {
		return err
	}
	return nil
}

type CheckUserInterface interface {
	CheckUser(creds request.LoginData, rds *redis.Client, ctx context.Context) (token string, err error)
}

// CheckUser Check if user registered in database
func (db *PostgresRepository) CheckUser(creds request.LoginData, rds *redis.Client, ctx context.Context) (token string, err error) {
	var storedUserID int
	var storedHashedPass string
	err = db.DB.QueryRow("SELECT id, password FROM users WHERE username=$1", creds.Username).Scan(&storedUserID, &storedHashedPass)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPass), []byte(creds.Password))
	if err != nil {
		return "", err
	}

	sessionToken := fmt.Sprintf("%d", time.Now().UnixNano())
	err = rds.Set(ctx, sessionToken, storedUserID, 24*time.Hour).Err()
	if err != nil {
		return "", err
	}
	return sessionToken, nil
}
