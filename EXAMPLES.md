# API Examples - Banque Tunisienne

Ce fichier contient des exemples complets JSON pour tester tous les endpoints de l'API bancaire tunisienne avec Postman ou autres outils de test.

## Configuration Postman

**Base URL:** `http://localhost:8080`

## 1. Health Check

**Endpoint:** `GET /health`

**Réponse attendue:**

```json
{
  "status": "OK",
  "message": "Bank API is running"
}
```

## 2. Créer un Compte

**Endpoint:** `POST /api/accounts`

**Headers:**

```
Content-Type: application/json
```

**Corps de requête - Compte Courant:**

```json
{
  "first_name": "Mohamed",
  "last_name": "Ben Ahmed",
  "email": "mohamed.benahmed@email.tn",
  "phone": "+21625123456",
  "address": "15 Avenue Habib Bourguiba, Tunis",
  "date_of_birth": "1990-05-15",
  "password": "motdepasse123",
  "account_type": "COMPTE_COURANT",
  "currency": "TND",
  "initial_balance": 1500000
}
```

**Corps de requête - Compte Épargne:**

```json
{
  "first_name": "Fatma",
  "last_name": "Karray",
  "email": "fatma.karray@banque.tn",
  "phone": "+21698765432",
  "address": "22 Rue de la République, Sfax",
  "date_of_birth": "1985-03-20",
  "password": "motdepasse456",
  "account_type": "COMPTE_EPARGNE",
  "currency": "TND",
  "initial_balance": 5000000
}
```

**Corps de requête - Compte Entreprise:**

```json
{
  "first_name": "Ahmed",
  "last_name": "Trabelsi",
  "email": "ahmed.trabelsi@entreprise.tn",
  "phone": "+21655444333",
  "address": "Boulevard de l'Environnement, Sousse",
  "date_of_birth": "1978-11-10",
  "password": "motdepasse789",
  "account_type": "COMPTE_ENTREPRISE",
  "currency": "EUR",
  "initial_balance": 10000000
}
```

**Corps de requête - Compte Devises (USD):**

```json
{
  "first_name": "Leila",
  "last_name": "Mansouri",
  "email": "leila.mansouri@international.tn",
  "phone": "+21622111999",
  "address": "Zone Industrielle, Monastir",
  "date_of_birth": "1992-08-05",
  "password": "motdepasse321",
  "account_type": "COMPTE_DEVISES",
  "currency": "USD",
  "initial_balance": 2000000
}
```

**Réponse attendue:**

```json
{
  "id": 1,
  "account_number": "TN5901234567890123456789",
  "first_name": "Mohamed",
  "last_name": "Ben Ahmed",
  "email": "mohamed.benahmed@email.tn",
  "phone": "+21625123456",
  "address": "15 Avenue Habib Bourguiba, Tunis",
  "date_of_birth": "1990-05-15T00:00:00Z",
  "balance": 1500000,
  "currency": "TND",
  "account_type": "COMPTE_COURANT",
  "status": "ACTIF",
  "iban": "TN5901234567890123456789",
  "bic": "STBKTNTT",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

## 3. Connexion (Login)

**Endpoint:** `POST /api/login`

**Headers:**

```
Content-Type: application/json
```

**Corps de requête:**

```json
{
  "account_number": "TN5901234567890123456789",
  "password": "motdepasse123"
}
```

**Réponse attendue:**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "account_number": "TN5901234567890123456789",
    "first_name": "Mohamed",
    "last_name": "Ben Ahmed",
    "email": "mohamed.benahmed@email.tn",
    "balance": 1500000,
    "currency": "TND",
    "account_type": "COMPTE_COURANT"
  }
}
```

## 4. Virement (Transfer)

**Endpoint:** `POST /api/transfer`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {votre_token_jwt}
```

**Corps de requête - Virement TND:**

```json
{
  "from_account": "TN5901234567890123456789",
  "to_account": "TN5998765432109876543210",
  "amount": 250000,
  "currency": "TND",
  "description": "Paiement facture électricité STEG"
}
```

**Corps de requête - Virement EUR:**

```json
{
  "from_account": "TN5912345678901234567890",
  "to_account": "TN5987654321098765432109",
  "amount": 500000,
  "currency": "EUR",
  "description": "Transfert vers compte épargne"
}
```

**Corps de requête - Virement USD:**

```json
{
  "from_account": "TN5923456789012345678901",
  "to_account": "TN5976543210987654321098",
  "amount": 1000000,
  "currency": "USD",
  "description": "Paiement fournisseur international"
}
```

**Réponse attendue:**

```json
{
  "id": 1,
  "from_account_id": 1,
  "to_account_id": 2,
  "amount": 250000,
  "currency": "TND",
  "description": "Paiement facture électricité STEG",
  "status": "RÉUSSI",
  "transaction_type": "VIREMENT",
  "created_at": "2024-01-15T14:25:30Z"
}
```

## 5. Historique des Transactions

**Endpoint:** `GET /api/transactions/{account_number}`

**Headers:**

```
Authorization: Bearer {votre_token_jwt}
```

**Exemple d'URL:**

```
GET /api/transactions/TN5901234567890123456789
```

**Réponse attendue:**

```json
[
  {
    "id": 1,
    "from_account_id": 1,
    "to_account_id": 2,
    "amount": 250000,
    "currency": "TND",
    "description": "Paiement facture électricité STEG",
    "status": "RÉUSSI",
    "transaction_type": "VIREMENT",
    "created_at": "2024-01-15T14:25:30Z"
  },
  {
    "id": 2,
    "from_account_id": 2,
    "to_account_id": 1,
    "amount": 100000,
    "currency": "TND",
    "description": "Remboursement prêt personnel",
    "status": "RÉUSSI",
    "transaction_type": "VIREMENT",
    "created_at": "2024-01-14T09:15:45Z"
  }
]
```

## 6. Scénarios de Test Complets

### Scénario 1: Test complet avec TND

1. **Créer premier compte:**

```json
POST /api/accounts
{
    "first_name": "Sami",
    "last_name": "Bouazizi",
    "email": "sami.bouazizi@gmail.tn",
    "phone": "+21620123456",
    "address": "Cité El Khadra, Tunis",
    "date_of_birth": "1988-12-03",
    "password": "sami2024",
    "account_type": "COMPTE_COURANT",
    "currency": "TND",
    "initial_balance": 2000000
}
```

2. **Créer deuxième compte:**

```json
POST /api/accounts
{
    "first_name": "Nour",
    "last_name": "Hamdi",
    "email": "nour.hamdi@yahoo.tn",
    "phone": "+21655987654",
    "address": "Centre Ville, Nabeul",
    "date_of_birth": "1995-07-18",
    "password": "nour2024",
    "account_type": "COMPTE_EPARGNE",
    "currency": "TND",
    "initial_balance": 1000000
}
```

3. **Se connecter avec le premier compte:**

```json
POST /api/login
{
    "account_number": "{numéro_retourné_étape_1}",
    "password": "sami2024"
}
```

4. **Effectuer un virement:**

```json
POST /api/transfer
{
    "from_account": "{numéro_compte_sami}",
    "to_account": "{numéro_compte_nour}",
    "amount": 300000,
    "currency": "TND",
    "description": "Cadeau d'anniversaire"
}
```

5. **Consulter l'historique:**

```
GET /api/transactions/{numéro_compte_sami}
```

### Scénario 2: Test avec devises étrangères

1. **Créer compte EUR:**

```json
POST /api/accounts
{
    "first_name": "Karim",
    "last_name": "Zouari",
    "email": "karim.zouari@export.tn",
    "phone": "+21698123456",
    "address": "Zone Franche, Bizerte",
    "date_of_birth": "1982-04-25",
    "password": "karim2024",
    "account_type": "COMPTE_ENTREPRISE",
    "currency": "EUR",
    "initial_balance": 15000000
}
```

2. **Créer compte USD:**

```json
POST /api/accounts
{
    "first_name": "Olfa",
    "last_name": "Mejri",
    "email": "olfa.mejri@trading.tn",
    "phone": "+21644555666",
    "address": "La Marsa, Tunis",
    "date_of_birth": "1990-09-12",
    "password": "olfa2024",
    "account_type": "COMPTE_DEVISES",
    "currency": "USD",
    "initial_balance": 8000000
}
```

## 7. Codes d'Erreur Courants

**400 Bad Request:**

```json
{
  "error": "Données invalides"
}
```

**401 Unauthorized:**

```json
{
  "error": "Token manquant ou invalide"
}
```

**404 Not Found:**

```json
{
  "error": "Compte non trouvé"
}
```

**409 Conflict:**

```json
{
  "error": "Email déjà utilisé"
}
```

**500 Internal Server Error:**

```json
{
  "error": "Erreur serveur interne"
}
```

## 8. Notes Importantes

### Devises Supportées

- **TND**: Dinar Tunisien (devise principale)
- **EUR**: Euro
- **USD**: Dollar Américain

### Types de Comptes

- **COMPTE_COURANT**: Compte courant standard
- **COMPTE_EPARGNE**: Compte d'épargne avec intérêts
- **COMPTE_ENTREPRISE**: Compte professionnel pour entreprises
- **COMPTE_DEVISES**: Compte multi-devises pour transactions internationales

### Format Monétaire

Tous les montants sont en **millimes** (1 TND = 1000 millimes):

- 1 000 millimes = 1.000 TND
- 250 000 millimes = 250.000 TND
- 1 500 000 millimes = 1 500.000 TND

### IBAN Tunisien

Format: **TN** + 2 chiffres de contrôle + 20 chiffres
Exemple: `TN5901234567890123456789`

### BIC Codes Tunisiens

- **STBKTNTT**: Société Tunisienne de Banque
- **BIATTNTT**: Banque Internationale Arabe de Tunisie
- **BNTUTNTT**: Banque Nationale Agricole
- **ABTNTNTT**: Arab Tunisian Bank
- **UBCITNTT**: Union Bancaire pour le Commerce et l'Industrie
  -d '{
  "account_number": "YOUR_ACCOUNT_NUMBER",
  "password": "motdepasse123"
  }'

````

**Save the JWT token from the response for authenticated requests.**

### 4. Get Account Information

Replace `{TOKEN}` with the JWT token from login:

```bash
curl -X GET http://localhost:3000/api/v1/accounts/1 \
  -H "Authorization: Bearer {TOKEN}"
````

### 5. Register a Second Account (for transfers)

```bash
curl -X POST http://localhost:3000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Fatma",
    "last_name": "Karray",
    "email": "fatma.karray@example.tn",
    "phone": "+21698765432",
    "date_of_birth": "1985-05-20",
    "password": "motdepasse456",
    "account_type": "COMPTE_EPARGNE",
    "currency": "TND",
    "address": {
      "street": "Rue de la Liberté 456",
      "city": "Sfax",
      "postal_code": "3000",
      "country": "Tunisia",
      "state": "Sfax"
    }
  }'
```

### 6. Create a Transaction (Deposit in TND)

Replace `{TOKEN}` with a valid JWT token and use actual account numbers:

```bash
curl -X POST http://localhost:3000/api/v1/transactions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {TOKEN}" \
  -d '{
    "from_account_number": "YOUR_ACCOUNT_NUMBER",
    "to_account_number": "RECIPIENT_ACCOUNT_NUMBER",
    "amount": 50000,
    "currency": "TND",
    "description": "Virement pour services",
    "reference": "FACT-2025-001"
  }'
```

### 7. Get Account Transactions

```bash
curl -X GET "http://localhost:3000/api/v1/transactions/account/1?limit=10&offset=0" \
  -H "Authorization: Bearer {TOKEN}"
```

### 8. Update Account Information

```bash
curl -X PUT http://localhost:3000/api/v1/accounts/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer {TOKEN}" \
  -d '{
    "phone": "+1234567891",
    "address": {
      "street": "789 Pine St",
      "city": "Chicago",
      "postal_code": "60601",
      "country": "United States",
      "state": "IL"
    }
  }'
```

## Complete Workflow Example

Here's a complete example workflow using PowerShell on Windows:

```powershell
# 1. Start the server (in a separate terminal)
make dev

# 2. Health check
$response = Invoke-RestMethod -Uri "http://localhost:3000/api/v1/health" -Method GET
Write-Output $response

# 3. Register first account
$registerData = @{
    first_name = "John"
    last_name = "Doe"
    email = "john.doe@example.com"
    phone = "+1234567890"
    date_of_birth = "1990-01-15"
    password = "securePassword123"
    account_type = "CHECKING"
    currency = "USD"
    address = @{
        street = "123 Main St"
        city = "New York"
        postal_code = "10001"
        country = "United States"
        state = "NY"
    }
} | ConvertTo-Json -Depth 3

$account1 = Invoke-RestMethod -Uri "http://localhost:3000/api/v1/auth/register" -Method POST -Body $registerData -ContentType "application/json"
Write-Output $account1

# 4. Login
$loginData = @{
    email = "john.doe@example.com"
    password = "securePassword123"
} | ConvertTo-Json

$loginResponse = Invoke-RestMethod -Uri "http://localhost:3000/api/v1/auth/login" -Method POST -Body $loginData -ContentType "application/json"
$token = $loginResponse.data.token
Write-Output "Token: $token"

# 5. Get account info
$headers = @{ Authorization = "Bearer $token" }
$accountInfo = Invoke-RestMethod -Uri "http://localhost:3000/api/v1/accounts/1" -Method GET -Headers $headers
Write-Output $accountInfo

# 6. Register second account for transfers
$registerData2 = @{
    first_name = "Jane"
    last_name = "Smith"
    email = "jane.smith@example.com"
    phone = "+1987654321"
    date_of_birth = "1985-05-20"
    password = "anotherSecurePassword456"
    account_type = "SAVINGS"
    currency = "USD"
    address = @{
        street = "456 Oak Ave"
        city = "Boston"
        postal_code = "02101"
        country = "United States"
        state = "MA"
    }
} | ConvertTo-Json -Depth 3

$account2 = Invoke-RestMethod -Uri "http://localhost:3000/api/v1/auth/register" -Method POST -Body $registerData2 -ContentType "application/json"
Write-Output $account2

# 7. Create a transaction (use actual account numbers from registration responses)
$transactionData = @{
    from_account_number = $account1.data.account_number
    to_account_number = $account2.data.account_number
    amount = 5000
    currency = "USD"
    description = "Test transfer"
    reference = "TEST-001"
} | ConvertTo-Json

$transaction = Invoke-RestMethod -Uri "http://localhost:3000/api/v1/transactions" -Method POST -Body $transactionData -ContentType "application/json" -Headers $headers
Write-Output $transaction
```

## Notes

- All monetary amounts are in cents/minor currency units (e.g., 5000 = $50.00)
- Account numbers are automatically generated during registration
- IBAN and BIC codes are automatically generated based on account details
- JWT tokens expire after 24 hours by default
- The database is reset each time the server starts (for development)

## Error Handling

The API returns structured error responses:

```json
{
  "success": false,
  "error": "Error message describing what went wrong",
  "timestamp": "2025-05-31T06:15:30Z"
}
```

Common HTTP status codes:

- `200` - Success
- `201` - Created
- `400` - Bad Request (validation errors)
- `401` - Unauthorized (missing/invalid token)
- `404` - Not Found
- `500` - Internal Server Error
