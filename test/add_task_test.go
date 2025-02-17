package test

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"main/internal/api/request"
	"main/internal/api/usecase"
	"main/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddTask_ExecuteWrongDate(t *testing.T) {
	mockRepo := &mocks.AddTaskInterface{}

	addTaskUsecase := &usecase.AddTask{
		Repo: mockRepo,
	}

	validTask := request.AddTask{
		Title:       "test",
		Description: "test",
		Date:        "07072025",
	}

	invalidTask := request.AddTask{
		Title:       "test",
		Description: "test",
		Date:        "12345", // Invalid date format
	}
	userID := 1

	expectedErr := errors.New("invalid date format, expected YYYYMMDD")
	//No error if length of date is 8
	mockRepo.On("AddTask", mock.MatchedBy(func(t request.AddTask) bool {
		return len(t.Date) == 8
	}), userID).Return(nil)
	//Error if length of date is not 8
	mockRepo.On("AddTask", mock.MatchedBy(func(task request.AddTask) bool {
		return len(task.Date) != 8
	}), userID).Return(expectedErr)

	//Test valid task
	err := addTaskUsecase.Execute(validTask, userID)
	assert.NoError(t, err)
	//Test invalid task
	err = addTaskUsecase.Execute(invalidTask, userID)
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)

	mockRepo.AssertExpectations(t)
}
