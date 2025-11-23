FROM golang:1.24.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o service ./cmd

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM alpine:latest

RUN apk --no-cache add postgresql-client

WORKDIR /app

COPY --from=builder /app/service .

COPY --from=builder /go/bin/migrate .

COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./service"]