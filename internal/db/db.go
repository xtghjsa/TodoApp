package db

import (
	"database/sql"
	"fmt"
	"log"
)

type DbConn struct {
	DB *sql.DB
}

func StartDb(PgsUser, PgsPass, PgsDbName, PgsHost, PgsPort string) (*DbConn, error) {
	// Check if env variables are not set
	if PgsUser == "" || PgsPass == "" || PgsDbName == "" || PgsHost == "" || PgsPort == "" {
		return nil, fmt.Errorf("lacking database environment variables, check .env file")

	}

	// Connecting to postgres server
	serverConnection := fmt.Sprintf("user=%s password=%s host=%s port=%s sslmode=disable",
		PgsUser, PgsPass, PgsHost, PgsPort)
	pgsServer, err := sql.Open("postgres", serverConnection)
	if err != nil {
		return nil, fmt.Errorf("error connecting to PostgreSQL server: %v", err)
	}
	defer pgsServer.Close()

	err = pgsServer.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging PostgreSQL server: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL server")

	//Check if database exists, if no - creating database
	var exists bool
	err = pgsServer.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname=$1)", PgsDbName).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("error checking if database exists: %v", err)
	}

	if !exists {
		_, err = pgsServer.Exec("CREATE DATABASE " + PgsDbName)
		if err != nil {
			return nil, fmt.Errorf("error creating database: %v", err)
		}
		log.Println("Successfully created " + PgsDbName + " database")
	}
	log.Println("Database " + PgsDbName + " exists")

	//Connecting to postgres database
	dbConnection := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		PgsUser, PgsPass, PgsHost, PgsPort, PgsDbName)
	pgsDb, err := sql.Open("postgres", dbConnection)
	if err != nil {
		return nil, fmt.Errorf("error connecting to PostgreSQL database: %v", err)
	}

	err = pgsDb.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging PostgreSQL server: %v", err)
	}

	// Creating `tasks` table if not exists
	createTasksTable := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		date VARCHAR(8),
		status VARCHAR(16) NOT NULL DEFAULT 'notDone',
		user_id INTEGER NOT NULL
	);`
	_, err = pgsDb.Exec(createTasksTable)
	if err != nil {
		return nil, fmt.Errorf("error creating `tasks` table: %v", err)
	}
	_, err = pgsDb.Exec("CREATE INDEX IF NOT EXISTS user_id_index ON tasks (user_id)")
	if err != nil {
		return nil, fmt.Errorf("error creating index on `tasks` table: %v", err)
	}

	// Creating `users` table if not exists
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(32) NOT NULL,
		password VARCHAR(128) NOT NULL
	);`
	_, err = pgsDb.Exec(createUsersTable)
	if err != nil {
		return nil, fmt.Errorf("error creating `users` table: %v", err)
	}
	_, err = pgsDb.Exec("CREATE UNIQUE INDEX IF NOT EXISTS username_index ON users (username)")
	if err != nil {
		return nil, fmt.Errorf("error creating index on `users` table: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return &DbConn{DB: pgsDb}, nil
}
