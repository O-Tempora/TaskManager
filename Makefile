LINUX_NAME=server

.PHONY: run build
build: 
	go build -o $(LINUX_NAME) -v ./cmd/server

run: build
	./$(LINUX_NAME)

.PHONY: compose build
compose: build
	docker compose up -d
	
.PHONY: git
git:
	git add .
	git commit -m "$m"
	git push

.DEFAULT_GOAL := build
