#Makefile for TodoApp

#Variables
APP_NAME = TodoApp
SRC_DIR = ./cmd
MOCK_DIR = ./mocks

all:build run generate-mocks

build:
	go build -o ./${APP_NAME} ${SRC_DIR}/main.go

run: build
	./${APP_NAME}

generate-mocks:
	mockery --all --output=${MOCK_DIR}