PROJECT_NAME = post-with-comments

DC = docker-compose
DOCKER_BUILD = docker build
DOCKER_RUN = docker run

# Переменные
DB_CONTAINER = db
APP_CONTAINER = app
MIGRATION_PATH = internal/pkg/db/migrations/postgres
DB_URL = "postgres://postgres:dbpass@localhost:5432/posts_with_comments?sslmode=disable"

.PHONY: test

up_build:
	$(DC) build
	$(DC) up -d

build:
	$(DC) build

up:
	$(DC) up -d

down:
	$(DC) down

logs:
	$(DC) logs $(APP_CONTAINER)

rebuild:
	$(DC) down
	$(DC) up --build -d


# Запуск контейнера базы данных отдельно
run-db:
	$(DOCKER_RUN) -p 5432:5432 --name $(DB_CONTAINER) \
	    -e "POSTGRES_PASSWORD=dbpass" \
	    -e "POSTGRES_DB=posts_with_comments" \
	    -d postgres:latest

# Удаление всех контейнеров и образов
clean:
	$(DC) down --rmi all
	docker system prune -f

# Линтер кода 
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run ./...

# Форматирование кода
fmt:
	gofmt -w .

# Запуск всех тестов
test:
	go test ./... -v

# Команды для миграций
migration-create-users:
	cd internal/pkg/db/migrations/ && migrate create -ext sql -dir postgres -seq create_users_table

migration-create-posts:
	cd internal/pkg/db/migrations/ && migrate create -ext sql -dir postgres -seq create_posts_table

migration-create-comments:
	cd internal/pkg/db/migrations/ && migrate create -ext sql -dir postgres -seq create_comments_table

# Применение миграций
migrate-up:
	migrate -verbose -database $(DB_URL) -path $(MIGRATION_PATH) up

# Откат миграций
migrate-down:
	migrate -verbose -database $(DB_URL) -path $(MIGRATION_PATH) down

