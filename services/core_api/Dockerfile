FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bin/core_api ./cmd/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /bin/core_api /app/core_api

EXPOSE 8082

CMD ["./core_api"]
