package main

import (
	"context"
	"fmt"
	"log"
	"main/internal/db"

	h "main/internal/handlers"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load environment variables from .env file: %v", err)
	}

	//Postgres env variables
	pgsUser := os.Getenv("POSTGRES_USER")
	pgsPass := os.Getenv("POSTGRES_PASSWORD")
	pgsDbName := os.Getenv("POSTGRES_DBNAME")
	pgsHost := os.Getenv("POSTGRES_HOST")
	pgsPort := os.Getenv("POSTGRES_PORT")

	//Server env variables
	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")

	//Redis env variables
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	//Start Redis
	ctx := context.Background()
	rds := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", redisHost, redisPort),
		DB:   0,
	})
	_, err = rds.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("failed to connect to Redis: %v\n", err)
	}
	//Redis connection struct
	redisConn := h.RedisConn{
		Connect: rds,
		Ctx:     &ctx,
	}
	log.Println("Connected to Redis")
	//Connect to Postgres database
	Conn, err := db.StartDb(pgsUser, pgsPass, pgsDbName, pgsHost, pgsPort)
	if err != nil {
		log.Fatalf("failed to start database: %v\n", err)
	}
	defer Conn.DB.Close()
	//Postgres connection struct
	dbHandler := h.Handler{DbConn: *Conn}

	//Login Struct
	loginHandler := h.LoginStruct{
		Db:    dbHandler,
		Redis: redisConn,
	}

	r := gin.Default()
	//Server routs
	r.POST("/reg", dbHandler.RegUser)
	r.POST("/login", loginHandler.Login)

	needAuth := r.Group("/todoapp")
	needAuth.Use(h.AuthMw(redisConn))
	{
		needAuth.POST("/add", dbHandler.AddTask)
		needAuth.PUT("/update", dbHandler.UpdTask)
		needAuth.PUT("/done", dbHandler.MarkTaskDone)
		needAuth.DELETE("/delete", dbHandler.DelTask)
		needAuth.GET("/show", dbHandler.ShowTasks)
	}
	//Start server
	if err := r.Run(serverHost + ":" + serverPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

}
