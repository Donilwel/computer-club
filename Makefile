# Переменные
  BINARY_NAME = club-service
  DOCKER_COMPOSE = docker-compose
  GO = go

# Компиляция бинарника
build:
	$(GO) build -o $(BINARY_NAME) cmd/server/main.go

# Запуск сервера локально
run:
	$(GO) run cmd/server/main.go

# Запуск тестов с покрытием
test:
	$(GO) test ./... -cover

# Линтинг кода (если установлен golangci-lint)
lint:
	golangci-lint run || echo "golangci-lint not installed"

# Форматирование кода
fmt:
	$(GO) fmt ./...

# Установка зависимостей
deps:
	$(GO) mod tidy

# Запуск контейнеров
docker-up:
	$(DOCKER_COMPOSE) up -d

# Остановка всех контейнеров
docker-down:
	$(DOCKER_COMPOSE) down

# Пересборка контейнеров
docker-rebuild:
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) build
	$(DOCKER_COMPOSE) up -d

# Очистка временных файлов
clean:
	rm -f $(BINARY_NAME)
	$(GO) clean

# Полный цикл: форматирование, сборка, тесты, запуск
all: fmt build test run
