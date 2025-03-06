# Используем Go 1.23
FROM golang:1.23 AS builder

WORKDIR /app

# Копируем файлы
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Компилируем
RUN go build -o club-service cmd/server/main.go

# Второй этап (минимальный контейнер)
FROM debian:latest
WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/club-service .

# Копируем .env в контейнер
COPY .env .env

# Загружаем переменные окружения
CMD ["./club-service"]
