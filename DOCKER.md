# Guide Docker - API Bancaire Tunisienne

Ce guide explique comment utiliser Docker pour développer et déployer l'API bancaire tunisienne.

## 🐳 Prérequis

- Docker Desktop ou Docker Engine
- Docker Compose
- Git

## 🚀 Démarrage Rapide

### 1. Cloner et Démarrer

```bash
git clone <repository-url>
cd bank-api
make docker-run
```

L'API sera disponible sur `http://localhost:8080`

### 2. Vérifier le Fonctionnement

```bash
# Vérifier la santé de l'API
curl http://localhost:8080/health

# Voir les logs
make docker-logs
```

## 📋 Commandes Disponibles

| Commande            | Description                      |
| ------------------- | -------------------------------- |
| `make docker-build` | Construire l'image Docker        |
| `make docker-run`   | Démarrer en arrière-plan         |
| `make docker-dev`   | Démarrer en mode développement   |
| `make docker-stop`  | Arrêter les conteneurs           |
| `make docker-clean` | Nettoyer tout (données incluses) |
| `make docker-logs`  | Voir les logs de l'application   |

## 🏗️ Architecture Docker

### Services

#### 1. **bank-api** (Application principale)

- **Port**: 8080
- **Image**: Construite depuis le Dockerfile local
- **Santé**: Endpoint `/health` vérifié toutes les 30s
- **Dépendances**: PostgreSQL doit être en santé

#### 2. **postgres** (Base de données)

- **Port**: 5433 (mappé depuis 5432)
- **Image**: `postgres:15-alpine`
- **Volume**: Données persistantes
- **Initialisation**: Script `init.sql` exécuté au premier démarrage

#### 3. **pgadmin** (Interface d'administration)

- **Port**: 5050
- **Image**: `dpage/pgadmin4:latest`
- **Accès**: `http://localhost:5050`
- **Identifiants**:
  - Email: `admin@banque-tunisia.tn`
  - Mot de passe: `admin123`

### Réseau

Tous les services communiquent via le réseau `bank-tunisia-network`.

### Volumes

- `postgres-bank-tunisia-data`: Données PostgreSQL persistantes
- `pgadmin-bank-tunisia-data`: Configuration PgAdmin persistante

## 🔧 Configuration

### Variables d'Environnement

L'application utilise ces variables d'environnement dans Docker:

```yaml
# Configuration base de données
DB_HOST: postgres
DB_PORT: 5432
DB_USER: bankgo
DB_PASSWORD: testbank
DB_NAME: bankdb_tunisia
DB_SSLMODE: disable

# Configuration application
PORT: 8080
JWT_SECRET: 'votre_cle_jwt_secrete_pour_banque_tunisienne_2024'
GIN_MODE: release

# Configuration bancaire tunisienne
BANK_COUNTRY: TN
DEFAULT_CURRENCY: TND
SUPPORTED_CURRENCIES: 'TND,EUR,USD'
```

### Personnalisation

Pour personnaliser la configuration:

1. Copiez `docker-compose.yml` vers `docker-compose.override.yml`
2. Modifiez les variables d'environnement selon vos besoins
3. Redémarrez avec `make docker-run`

## 📊 Surveillance et Debug

### Vérification des Services

```bash
# Statut des conteneurs
docker ps

# Logs détaillés
docker-compose logs -f

# Logs d'un service spécifique
docker-compose logs -f bank-api
docker-compose logs -f postgres
```

### Accès aux Conteneurs

```bash
# Shell dans l'application
docker-compose exec bank-api sh

# Shell dans PostgreSQL
docker-compose exec postgres psql -U bankgo -d bankdb_tunisia
```

### Vérification Base de Données

```bash
# Se connecter à PostgreSQL
docker-compose exec postgres psql -U bankgo -d bankdb_tunisia

# Vérifier les tables
\dt

# Vérifier les comptes de test
SELECT * FROM accounts;
```

## 🔍 Résolution de Problèmes

### Problème: L'API ne démarre pas

**Solution 1**: Vérifier que PostgreSQL est prêt

```bash
docker-compose logs postgres
```

**Solution 2**: Reconstruire l'image

```bash
make docker-clean
make docker-build
make docker-run
```

### Problème: Erreur de connexion base de données

**Vérification**:

```bash
# Tester la connexion depuis l'application
docker-compose exec bank-api ping postgres

# Vérifier PostgreSQL
docker-compose exec postgres pg_isready -U bankgo
```

### Problème: Port déjà utilisé

**Solution**: Modifier les ports dans `docker-compose.yml`

```yaml
services:
  bank-api:
    ports:
      - '8081:8080' # Utiliser le port 8081 au lieu de 8080
```

### Problème: Données corrompues

**Solution**: Nettoyer et redémarrer

```bash
make docker-clean
make docker-run
```

## 🔐 Sécurité en Production

### 1. Variables Sensibles

Utilisez Docker secrets ou variables d'environnement externes:

```bash
# Exemple avec .env
echo "JWT_SECRET=votre_vraie_cle_secrete_production" > .env
docker-compose --env-file .env up
```

### 2. Réseau

En production, utilisez un réseau privé:

```yaml
networks:
  bank-network:
    driver: bridge
    internal: true # Pas d'accès Internet direct
```

### 3. Volumes

Sauvegardez régulièrement le volume PostgreSQL:

```bash
# Sauvegarde
docker run --rm -v postgres-bank-tunisia-data:/data -v $(pwd):/backup alpine tar czf /backup/backup.tar.gz /data

# Restauration
docker run --rm -v postgres-bank-tunisia-data:/data -v $(pwd):/backup alpine tar xzf /backup/backup.tar.gz -C /
```

## 🚀 Déploiement en Production

### 1. Optimisations Image

Le Dockerfile utilise déjà:

- Multi-stage build pour réduire la taille
- Image Alpine légère
- Utilisateur non-root
- Certificats CA inclus

### 2. Configuration Production

```yaml
# docker-compose.prod.yml
version: '3.8'
services:
  bank-api:
    environment:
      GIN_MODE: release
      JWT_SECRET: ${JWT_SECRET_PROD}
    restart: always
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
```

### 3. Monitoring

Ajoutez des services de monitoring:

```yaml
services:
  prometheus:
    image: prom/prometheus
    ports:
      - '9090:9090'

  grafana:
    image: grafana/grafana
    ports:
      - '3000:3000'
```

## 📝 Tests avec Docker

### Tests Automatisés

```bash
# Construire et tester
docker build -t bank-api-test --target builder .
docker run --rm bank-api-test go test ./tests/...
```

### Tests d'Intégration

```bash
# Démarrer l'environnement de test
docker-compose -f docker-compose.test.yml up -d

# Exécuter les tests
curl http://localhost:8080/health

# Nettoyer
docker-compose -f docker-compose.test.yml down
```

---

**🏦 API Bancaire Tunisienne - Prête pour la production avec Docker! 🇹🇳**
