FROM golang:1.25.3 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/main.go
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /go/bin/migrate .
COPY --from=builder /app/internal/migrations ./migrations
EXPOSE 8080
CMD ["sh", "-c", "./migrate -path ./migrations -database \"$DATABASE_URL\" up && ./main"]