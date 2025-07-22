# Build Stage
FROM golang:1.24.5-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8083

CMD [ "/app/main" ]