FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./src .

RUN go build -o bin


FROM alpine:latest


WORKDIR /app

COPY --from=builder /app/bin .

EXPOSE 8080

CMD ["./go-api"]