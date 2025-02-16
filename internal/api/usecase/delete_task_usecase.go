package usecase

import (
	"main/internal/api/request"
	"main/internal/repository"
)

type DeleteTask struct {
	Repo repository.DeleteTaskInterface
}

func (u *DeleteTask) Execute(id request.DeleteTaskID, userID any) error {
	return u.Repo.DeleteTask(id, userID)
}
