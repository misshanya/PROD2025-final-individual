FROM golang:1.23.6-alpine as builder

WORKDIR /app

COPY . .

# Download requirements
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o backend ./cmd

# Alpine as runner
FROM alpine as runner

WORKDIR /app
COPY --from=builder /app/internal/infrastructure/db/migrations/ ./migrations
COPY --from=builder /app/backend .

# Install goose for db migrations
RUN wget -O /usr/local/bin/goose https://github.com/pressly/goose/releases/download/v3.24.1/goose_linux_x86_64
RUN chmod +x /usr/local/bin/goose

# Migrate db and run server
CMD ["sh", "-c", "goose -dir ./migrations postgres $DATABASE_URL up && ./backend"]