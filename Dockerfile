FROM golang:1.24-alpine AS backend-builder

WORKDIR /app
COPY packages/backend/go.* ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY ./packages/backend/ ./

RUN swag init
RUN go build -o app .

FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend

COPY ./packages/frontend/package*.json ./
RUN npm install
COPY --from=backend-builder /app/docs ./docs
COPY ./packages/frontend/ ./
RUN npm run swag
RUN ls src
RUN npm run build

FROM golang:1.24.2-alpine

LABEL org.opencontainers.image.source="https://github.com/cableship/chart-sentinel"

WORKDIR /app
COPY packages/backend/go.* ./

COPY --from=frontend-builder /app/frontend/dist ./static
COPY --from=backend-builder /app/app ./
COPY --from=backend-builder /app/docs ./docs

ARG ARCH=
EXPOSE 8080

CMD ["./app"]

