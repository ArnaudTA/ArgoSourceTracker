FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend
COPY ./packages/frontend/package*.json ./
RUN npm install
COPY ./packages/frontend/ ./
RUN npm run build

FROM golang:1.24-alpine

WORKDIR /app
COPY packages/backend/go.* ./
RUN go mod download

COPY ./packages/backend/ ./
COPY --from=frontend-builder /app/frontend/dist ./static

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]