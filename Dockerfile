FROM golang:1.23.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main .

FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

ENTRYPOINT ["/app/main"]