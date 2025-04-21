FROM golang:1.24-alpine

WORKDIR /app

# Install git and air for development
RUN apk add --no-cache git && \
    go install github.com/air-verse/air@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

COPY --chown=1000:1000 go.* ./
RUN go mod download

# ENV GOCACHE=/app/.cache
RUN mkdir /.cache && chown -R 1000:1000 /.cache
ENV USER=1000
USER 1000

EXPOSE 8080

# Use air for development (hot reload)
CMD ["air"]