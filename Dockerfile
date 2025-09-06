FROM golang:1.24-alpine AS base

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# # Install air for hot reload
RUN go install github.com/air-verse/air@latest

COPY . .

# Build the application
RUN go build -o main .

EXPOSE 8000

CMD ["air", "-c", ".air.toml"]