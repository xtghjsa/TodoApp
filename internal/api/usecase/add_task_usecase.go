package usecase

import (
	"main/internal/api/request"
	"main/internal/repository"
)

type AddTask struct {
	Repo repository.AddTaskInterface
}

func (u *AddTask) Execute(task request.AddTask, userID any) error {
	return u.Repo.AddTask(task, userID)
}
