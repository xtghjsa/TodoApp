package usecase

import (
	"main/internal/api/request"
	"main/internal/repository"
)

type UpdateTask struct {
	Repo repository.UpdateTaskInterface
}

func (u *UpdateTask) Execute(task request.UpdateTask, userID any) error {
	return u.Repo.UpdateTask(task, userID)
}
