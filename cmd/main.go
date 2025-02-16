package main

import (
	"log"
	"main/internal/api"
	"main/pkg"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load environment variables from .env file: %v", err)
	}

	//Postgres .env variables
	pgsUser := os.Getenv("POSTGRES_USER")
	pgsPass := os.Getenv("POSTGRES_PASSWORD")
	pgsDbName := os.Getenv("POSTGRES_DBNAME")
	pgsHost := os.Getenv("POSTGRES_HOST")
	pgsPort := os.Getenv("POSTGRES_PORT")

	//Server .env variables
	serverHost := os.Getenv("SERVER_HOST")
	serverPort := os.Getenv("SERVER_PORT")

	//Redis .env variables
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")

	//Start Redis
	rds, ctx, err := pkg.InitializeRedis(redisHost, redisPort)
	if err != nil {
		log.Fatalf("Error connecting to redis: %v\n", err)
	}
	log.Println("Connected to Redis")

	//Connect to Postgres database
	db, err := pkg.InitializeDatabase(pgsUser, pgsPass, pgsDbName, pgsHost, pgsPort)
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}
	defer db.Close()
	log.Println("Connected to PostgreSQL database")
	//Starting server
	err = api.StartServer(db, rds, ctx, serverHost, serverPort)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}
}
