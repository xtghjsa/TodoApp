package usecase

import (
	"context"
	"main/internal/api/request"
	"main/internal/repository"

	"github.com/redis/go-redis/v9"
)

type Login struct {
	Repo repository.CheckUserInterface
	Ctx  context.Context
	Rds  *redis.Client
}

func (u *Login) Execute(creds request.LoginData) (token string, err error) {
	return u.Repo.CheckUser(creds, u.Rds, u.Ctx)
}
