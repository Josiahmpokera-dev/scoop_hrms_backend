# Build stage
FROM golang:1.25-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o hrms-backend ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o hrms-migrator ./cmd/migrator

# Production stage
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata postgresql-client
RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -D appuser
WORKDIR /app
COPY --from=builder /app/hrms-backend .
COPY --from=builder /app/hrms-migrator .
RUN mkdir -p /app/storage && chown -R appuser:appgroup /app
USER appuser
EXPOSE 8080
# Default runtime is app only (no automatic migration).
# Run migrations explicitly when needed:
#   ./hrms-migrator
CMD ["./hrms-backend"]
