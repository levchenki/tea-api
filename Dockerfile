FROM golang:1.24.2-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY migrations/postgres /app/migrations/postgres

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ./target/tea-api ./cmd/app

EXPOSE 8080

CMD ["./target/tea-api"]