package usecase

import (
	"main/internal/api/request"
	"main/internal/repository"
)

type MarkTaskDone struct {
	Repo repository.MarkTaskDoneInterface
}

func (u *MarkTaskDone) Execute(id request.MarkTaskDoneID, userID any) error {
	return u.Repo.MarkTaskDone(id, userID)
}
