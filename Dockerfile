# Build stage
FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build from cmd/api directory (FIXED PATH)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/hrms-backend ./cmd/api

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata postgresql-client
RUN addgroup -g 1000 appgroup && adduser -u 1000 -G appgroup -D appuser
WORKDIR /app
COPY wait-for-postgres.sh /wait-for-postgres.sh
RUN chmod +x /wait-for-postgres.sh
COPY --from=builder /app/bin/hrms-backend .
RUN mkdir -p /app/storage && chown -R appuser:appgroup /app
USER appuser
EXPOSE 8081
CMD ["/wait-for-postgres.sh", "postgres", "5432", "./hrms-backend"]