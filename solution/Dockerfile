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
COPY --from=builder /app/backend .

CMD ["/app/backend"]