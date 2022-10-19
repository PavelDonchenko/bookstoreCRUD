APP_NAME=bookstorecrud

build: ## Build the release and develoment container. The development
	docker-compose build

delete: ## Delete container and images
	docker system prune -a --volumes

dev: ## Run container in development mode
	docker-compose build --no-cache $(APP_NAME) && docker-compose run $(APP_NAME)

# Build and run the container
up: ## Spin up the project
	docker-compose up

stop: ## Stop running containers
	docker stop

rm: stop ## Stop and remove running containers
	docker rm $(APP_NAME)

lint: ## Run golangci-lint
	golangci-lint run
	go vet ./...

run: ##run application
	go run main.go

docker-test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit