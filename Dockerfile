FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o bin/app ./src


FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/config.json ./
COPY --from=builder /app/bin/app .

EXPOSE 3000

CMD ["./app"]

# RUN apk add --no-cache bash

# CMD ["/bin/bash"]