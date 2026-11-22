# Этап 1. Сборка Go-приложения.
FROM golang:1.27-alpine AS builder

WORKDIR /src

# Сначала копируем файлы зависимостей — этот слой можно использовать повторно.
COPY go.mod go.sum ./
RUN go mod download

# Затем копируем исходный код и собираем бинарный файл.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/student-service .

# Этап 2. Минимальный образ для запуска.
FROM alpine:3.23

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S app -G app

WORKDIR /app
COPY --from=builder /out/student-service ./student-service

# Приложение слушает HTTP :8080 и gRPC :9090.
EXPOSE 8080 9090

# Запускаем контейнер не от root.
USER app
ENTRYPOINT ["/app/student-service"]
