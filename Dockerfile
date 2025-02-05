FROM golang:1.23-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/main ./cmd/judge_integrator

FROM golang:1.23-alpine AS migrator
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM alpine:3.18
COPY --from=migrator /go/bin/migrate /usr/local/bin/migrate

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/main /app/main
COPY --from=migrator /usr/local/bin/migrate /usr/local/bin/migrate

COPY migrations /app/migrations
COPY .env /app/

EXPOSE 8080

CMD ["sh", "-c", "migrate -path ./migrations -database ${DB_CONN} up && /app/main"]
