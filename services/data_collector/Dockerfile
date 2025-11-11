FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /bin/data_collector ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /bin/data_collector /app/data_collector

ENV DATABASE_ADDR=${DATABASE_ADDR}
ENV DB_NAME=${DB_NAME}
ENV DB_USER=${DB_USER}
ENV API_URL=${API_URL}
ENV UPDATE_INTERVAL_MINUTES=${UPDATE_INTERVAL_MINUTES}

EXPOSE 8090

CMD ["./data_collector"]
