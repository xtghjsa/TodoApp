package usecase

import (
	"main/internal/api/response"
	"main/internal/repository"
)

type ShowTasks struct {
	Repo repository.ShowTasksInterface
}

func (u *ShowTasks) Execute(query, queryKey string, userID any) (response.ShowTasks, error) {
	return u.Repo.ShowTasks(query, queryKey, userID)
}
