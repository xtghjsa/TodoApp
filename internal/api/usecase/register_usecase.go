package usecase

import (
	"main/internal/api/request"
	"main/internal/repository"
)

type RegisterUser struct {
	Repo repository.RegisterUserInterface
}

func (u *RegisterUser) Execute(regData request.RegistrationData) error {
	return u.Repo.RegisterUser(regData)
}
