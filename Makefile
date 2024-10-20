# Имя проекта
PROJECT_NAME = post-with-comments

# Команды для Docker
DC = docker-compose
DOCKER_BUILD = docker build
DOCKER_RUN = docker run

# Переменные
DB_CONTAINER = db
APP_CONTAINER = app
MIGRATION_PATH = internal/pkg/db/migrations/postgres
DB_URL = "postgres://postgres:dbpass@localhost:5432/posts_with_comments?sslmode=disable"

# Сборка всех контейнеров
build:
	$(DC) build

# Запуск всех контейнеров
up:
	$(DC) up -d

# Остановка всех контейнеров
down:
	$(DC) down

# Просмотр логов приложения
logs:
	$(DC) logs $(APP_CONTAINER)

# Пересобрать и запустить контейнеры с нуля
rebuild:
	$(DC) down
	$(DC) up --build -d

# Проверка статуса контейнеров
status:
	$(DC) ps

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

# Команды для миграций
# Создание новой миграции для users, posts и comments таблиц
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

