# Dockerfile pour l'API Bancaire Tunisienne
FROM golang:1.24-alpine AS builder

# Installer les dépendances système nécessaires
RUN apk add --no-cache git ca-certificates tzdata

# Définir le répertoire de travail
WORKDIR /app

# Copier les fichiers de modules Go
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod download

# Copier le code source
COPY . .

# Construire l'application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server

# Image finale - multi-stage build pour réduire la taille
FROM alpine:latest

# Installer ca-certificates pour les connexions HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Créer un utilisateur non-root pour la sécurité
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /root/

# Copier l'exécutable depuis l'image builder
COPY --from=builder /app/main .

# Changer le propriétaire du fichier
RUN chown appuser:appgroup main

# Changer vers l'utilisateur non-root
USER appuser

# Exposer le port 8080
EXPOSE 8080

# Variables d'environnement par défaut
ENV GIN_MODE=release
ENV PORT=8080

# Commande pour démarrer l'application
CMD ["./main"]
