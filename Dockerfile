# Utilise une image Go officielle pour le build
FROM golang:1.24.2 AS builder

# Crée un dossier de travail
WORKDIR /app

# Copie les fichiers go
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copie le reste du code
COPY . .

# Compile l’app en binaire statique
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# Image finale minimaliste
FROM alpine:latest

WORKDIR /root/

# Copie le binaire depuis l’étape de build
COPY --from=builder /app/app .

# Expose le port
EXPOSE 8080

# Lancement
CMD ["./app"]
