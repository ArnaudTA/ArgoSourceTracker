FROM golang:1.24-alpine

WORKDIR /app

# Install git and air for development
RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest

COPY go.* ./
RUN go mod download

COPY . .

EXPOSE 8080

# Use air for development (hot reload)
CMD ["air"]