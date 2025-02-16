package api

import (
	"context"
	"database/sql"
	"log"
	"main/internal/api/handler"
	"main/internal/api/middleware"
	"main/internal/api/usecase"
	"main/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func StartServer(db *sql.DB, rds *redis.Client, ctx context.Context, serverHost, serverPort string) error {

	//Initialize repository
	repo := &repository.PostgresRepository{DB: db}
	//Initialize usecases
	addTaskUsecase := &usecase.AddTask{Repo: repo}
	deleteTaskUsecase := &usecase.DeleteTask{Repo: repo}
	updateTaskUsecase := &usecase.UpdateTask{Repo: repo}
	markTaskDoneUsecase := &usecase.MarkTaskDone{Repo: repo}
	showTasksUsecase := &usecase.ShowTasks{Repo: repo}
	regUserUsecase := &usecase.RegisterUser{Repo: repo}
	loginUsecase := &usecase.Login{Repo: repo, Ctx: ctx, Rds: rds}

	//Initialize handlers
	addTaskHandler := &handler.AddTaskHandler{Usecase: addTaskUsecase}
	deleteTaskHandler := &handler.DeleteTaskHandler{Usecase: deleteTaskUsecase}
	updateTaskHandler := &handler.UpdateTaskHandler{Usecase: updateTaskUsecase}
	markTaskDoneHandler := &handler.MarkTaskDoneHandler{Usecase: markTaskDoneUsecase}
	showTasksHandler := &handler.ShowTasksHandler{Usecase: showTasksUsecase}
	regUserHandler := &handler.RegUserHandler{Usecase: regUserUsecase}
	loginHandler := &handler.LoginHandler{Usecase: loginUsecase}

	r := gin.Default()

	r.POST("/reg", regUserHandler.RegisterUser)
	r.POST("/login", loginHandler.Login)

	needAuth := r.Group("/todoapp")
	needAuth.Use(middleware.AuthenticationMiddleware(rds, ctx))
	{
		needAuth.POST("/add", addTaskHandler.AddTask)
		needAuth.PUT("/update", updateTaskHandler.UpdateTask)
		needAuth.PUT("/done", markTaskDoneHandler.MarkTaskDone)
		needAuth.DELETE("/delete", deleteTaskHandler.DelTask)
		needAuth.GET("/show", showTasksHandler.ShowTasks)
	}
	// Run server
	if err := r.Run(serverHost + ":" + serverPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
	return nil
}
