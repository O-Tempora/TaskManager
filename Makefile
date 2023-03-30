LINUX_NAME=server

.PHONY: run build

build: 
	go build -o $(LINUX_NAME) -v ./cmd/server

run: build
	./$(LINUX_NAME)
	

.DEFAULT_GOAL := build

