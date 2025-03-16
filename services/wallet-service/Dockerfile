# Используем базовый образ Golang
FROM golang:1.20

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY . .

# Загружаем зависимости
RUN go mod tidy

# Сборка приложения
RUN go build -o main .

# Запуск сервиса
CMD ["/app/main"]
