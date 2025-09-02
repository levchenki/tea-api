FROM golang:1.24.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./target/tea-api ./cmd/app

FROM alpine:3.20

COPY --from=builder /app/target/tea-api .
COPY --from=builder /app/migrations/postgres ./migrations/postgres

EXPOSE 8080

#fake commit
CMD ["./target/tea-api"]