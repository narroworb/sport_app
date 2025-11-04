FROM golang:1.24-alpine AS builder

WORKDIR /app

# Скопируем только нужные файлы (ускоряет сборку)
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник (в папку /bin)
RUN go build -o /bin/auth-service ./cmd/main.go

# --- Runtime stage ---
FROM alpine:latest

WORKDIR /app
COPY --from=builder /bin/auth-service /app/auth-service

EXPOSE 8081

CMD ["./auth-service"]
