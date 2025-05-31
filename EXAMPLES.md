# 🏦 API Examples - Tunisian Bank API

This file contains complete JSON examples for testing all endpoints of the Tunisian Banking API with Postman or other testing tools.

## 🔧 Postman Configuration

**Base URL:** `http://localhost:8080`

## 1. 💓 Health Check

**Endpoint:** `GET /api/v1/health`

**Expected Response:**

```json
{
  "status": "OK",
  "message": "Bank API is running",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

## 2. 📝 Create Account

**Endpoint:** `POST /api/v1/accounts`

**Headers:**

```
Content-Type: application/json
```

**Request Body - Checking Account:**

```json
{
  "first_name": "Mohamed",
  "last_name": "Ben Ahmed",
  "email": "mohamed.benahmed@example.tn",
  "phone": "+21612345678",
  "date_of_birth": "1990-01-15T00:00:00Z",
  "password": "securepassword123",
  "account_type": "CHECKING",
  "currency": "TND",
  "address": {
    "street": "Avenue Habib Bourguiba 123",
    "city": "Tunis",
    "postal_code": "1001",
    "country": "Tunisia",
    "state": "Tunis"
  }
}
```

**Request Body - Savings Account:**

```json
{
  "first_name": "Fatma",
  "last_name": "Karray",
  "email": "fatma.karray@example.tn",
  "phone": "+21698765432",
  "date_of_birth": "1985-03-20T00:00:00Z",
  "password": "securepassword456",
  "account_type": "SAVINGS",
  "currency": "TND",
  "address": {
    "street": "Rue de la République 22",
    "city": "Sfax",
    "postal_code": "3000",
    "country": "Tunisia",
    "state": "Sfax"
  }
}
```

**Request Body - Business Account:**

```json
{
  "first_name": "Ahmed",
  "last_name": "Trabelsi",
  "email": "ahmed.trabelsi@business.tn",
  "phone": "+21655444333",
  "date_of_birth": "1978-11-10T00:00:00Z",
  "password": "securepassword789",
  "account_type": "BUSINESS",
  "currency": "EUR",
  "address": {
    "street": "Boulevard de l'Environnement 45",
    "city": "Sousse",
    "postal_code": "4000",
    "country": "Tunisia",
    "state": "Sousse"
  }
}
```

**Request Body - Foreign Currency Account (USD):**

```json
{
  "first_name": "Leila",
  "last_name": "Mansouri",
  "email": "leila.mansouri@international.tn",
  "phone": "+21622111999",
  "date_of_birth": "1992-08-05T00:00:00Z",
  "password": "securepassword321",
  "account_type": "FOREIGN_CURRENCY",
  "currency": "USD",
  "address": {
    "street": "Zone Industrielle 15",
    "city": "Monastir",
    "postal_code": "5000",
    "country": "Tunisia",
    "state": "Monastir"
  }
}
```

**Expected Response:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "customer_id": 1,
    "account_number": "TN5961705312451143542106",
    "iban": "TN5961705312451143542106",
    "bic": "STBKTNTT",
    "account_type": "CHECKING",
    "currency": "TND",
    "balance": 0,
    "status": "ACTIVE",
    "created_at": "2025-01-31T06:15:30Z",
    "updated_at": "2025-01-31T06:15:30Z"
  },
  "message": "Account created successfully",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

## 3. 🔐 Authentication (Login)

**Endpoint:** `POST /api/v1/auth/login`

**Headers:**

```
Content-Type: application/json
```

**Request Body:**

```json
{
  "account_number": "TN5961705312451143542106",
  "password": "securepassword123"
}
```

**Expected Response:**

````json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400,
    "account": {
      "id": 1,
      "customer_id": 1,
      "account_number": "TN5961705312451143542106",
      "account_type": "CHECKING",
      "currency": "TND",
      "status": "ACTIVE"
    }
  },
  "message": "Login successful",
  "timestamp": "2025-01-31T06:15:30Z"
}
```## 4. 🔄 Refresh Token

**Endpoint:** `POST /api/v1/auth/refresh`

**Headers:**

````

Authorization: Bearer {your_current_jwt_token}

````

**Expected Response:**

```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 86400
  },
  "message": "Token refreshed successfully",
  "timestamp": "2025-01-31T06:15:30Z"
}
````

## 5. 🚪 Logout

**Endpoint:** `POST /api/v1/auth/logout`

**Headers:**

```
Authorization: Bearer {your_jwt_token}
```

**Expected Response:**

```json
{
  "success": true,
  "message": "Logout successful",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

## 6. 💸 Transfer Money

**Endpoint:** `POST /api/v1/transactions/transfer`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {your_jwt_token}
```

**Request Body - TND Transfer:**

```json
{
  "from_account_number": "TN5961705312451143542106",
  "to_account_number": "TN5959238705041140193701",
  "amount": 25000,
  "currency": "TND",
  "description": "Electricity bill payment STEG",
  "reference": "TXN-2025-001"
}
```

**Request Body - EUR Transfer:**

```json
{
  "from_account_number": "TN5961705312451143542106",
  "to_account_number": "TN5959238705041140193701",
  "amount": 500,
  "currency": "EUR",
  "description": "Transfer to savings account",
  "reference": "TXN-2025-002"
}
```

**Request Body - USD Transfer:**

```json
{
  "from_account_number": "TN5961705312451143542106",
  "to_account_number": "TN5959238705041140193701",
  "amount": 1000,
  "currency": "USD",
  "description": "International supplier payment",
  "reference": "TXN-2025-003"
}
```

**Expected Response:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "from_account_id": 1,
    "to_account_id": 2,
    "from_account_number": "TN5961705312451143542106",
    "to_account_number": "TN5959238705041140193701",
    "amount": 25000,
    "currency": "TND",
    "description": "Electricity bill payment STEG",
    "reference": "TXN-2025-001",
    "transaction_type": "TRANSFER",
    "status": "COMPLETED",
    "fee": 500,
    "exchange_rate": 1.0,
    "created_at": "2025-01-31T06:15:30Z"
  },
  "message": "Transfer completed successfully",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

## 7. 📥 Deposit Money

**Endpoint:** `POST /api/v1/transactions/deposit`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {your_jwt_token}
```

**Request Body:**

```json
{
  "account_number": "TN5961705312451143542106",
  "amount": 100000,
  "currency": "TND",
  "description": "Cash deposit"
}
```

**Expected Response:**

```json
{
  "success": true,
  "data": {
    "id": 2,
    "account_id": 1,
    "account_number": "TN5961705312451143542106",
    "amount": 100000,
    "currency": "TND",
    "description": "Cash deposit",
    "transaction_type": "DEPOSIT",
    "status": "COMPLETED",
    "fee": 0,
    "created_at": "2025-01-31T06:15:30Z"
  },
  "message": "Deposit completed successfully",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

## 8. 📤 Withdraw Money

**Endpoint:** `POST /api/v1/transactions/withdraw`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {your_jwt_token}
```

**Request Body:**

```json
{
  "account_number": "TN5961705312451143542106",
  "amount": 10000,
  "currency": "TND",
  "description": "ATM withdrawal"
}
```

**Expected Response:**

```json
{
  "success": true,
  "data": {
    "id": 3,
    "account_id": 1,
    "account_number": "TN5961705312451143542106",
    "amount": 10000,
    "currency": "TND",
    "description": "ATM withdrawal",
    "transaction_type": "WITHDRAWAL",
    "status": "COMPLETED",
    "fee": 200,
    "created_at": "2025-01-31T06:15:30Z"
  },
  "message": "Withdrawal completed successfully",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

## 9. 💰 Get Account Balance

**Endpoint:** `GET /api/v1/accounts/{account_number}/balance`

**Headers:**

```
Authorization: Bearer {your_jwt_token}
```

**Example URL:**

```
GET /api/v1/accounts/TN5961705312451143542106/balance
```

**Expected Response:**

```json
{
  "success": true,
  "data": {
    "account_number": "TN5961705312451143542106",
    "balance": 114300,
    "currency": "TND",
    "available_balance": 114300,
    "pending_balance": 0
  },
  "message": "Balance retrieved successfully",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

## 10. 📋 Transaction History

**Endpoint:** `GET /api/v1/transactions/account/{account_id}`

**Headers:**

```
Authorization: Bearer {your_jwt_token}
```

**Example URL:**

```
GET /api/v1/transactions/account/1?limit=10&offset=0
```

**Expected Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "from_account_id": 1,
      "to_account_id": 2,
      "from_account_number": "TN5961705312451143542106",
      "to_account_number": "TN5959238705041140193701",
      "amount": 25000,
      "currency": "TND",
      "description": "Electricity bill payment STEG",
      "reference": "TXN-2025-001",
      "transaction_type": "TRANSFER",
      "status": "COMPLETED",
      "fee": 500,
      "created_at": "2025-01-31T06:15:30Z"
    },
    {
      "id": 2,
      "account_id": 1,
      "account_number": "TN5961705312451143542106",
      "amount": 100000,
      "currency": "TND",
      "description": "Cash deposit",
      "transaction_type": "DEPOSIT",
      "status": "COMPLETED",
      "fee": 0,
      "created_at": "2025-01-31T06:10:15Z"
    }
  ],
  "message": "Transaction history retrieved successfully",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

## 11. 👤 Get Account Information

**Endpoint:** `GET /api/v1/accounts/{id}`

**Headers:**

```
Authorization: Bearer {your_jwt_token}
```

**Example URL:**

```
GET /api/v1/accounts/1
```

**Expected Response:**

```json
{
  "success": true,
  "data": {
    "id": 1,
    "customer_id": 1,
    "account_number": "TN5961705312451143542106",
    "iban": "TN5961705312451143542106",
    "bic": "STBKTNTT",
    "account_type": "CHECKING",
    "currency": "TND",
    "balance": 114300,
    "status": "ACTIVE",
    "created_at": "2025-01-31T06:00:00Z",
    "updated_at": "2025-01-31T06:15:30Z",
    "customer": {
      "id": 1,
      "first_name": "Mohamed",
      "last_name": "Ben Ahmed",
      "email": "mohamed.benahmed@example.tn",
      "phone": "+21612345678",
      "date_of_birth": "1990-01-15T00:00:00Z",
      "address": {
        "street": "Avenue Habib Bourguiba 123",
        "city": "Tunis",
        "postal_code": "1001",
        "country": "Tunisia",
        "state": "Tunis"
      }
    }
  },
  "message": "Account information retrieved successfully",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

## 🧪 Complete Testing Scenarios

### Scenario 1: Complete TND Testing

1. **Create first account:**

```json
POST /api/v1/accounts
{
  "first_name": "Sami",
  "last_name": "Bouazizi",
  "email": "sami.bouazizi@example.tn",
  "phone": "+21620123456",
  "date_of_birth": "1988-12-03T00:00:00Z",
  "password": "sami2024",
  "account_type": "CHECKING",
  "currency": "TND",
  "address": {
    "street": "Cité El Khadra",
    "city": "Tunis",
    "postal_code": "1001",
    "country": "Tunisia",
    "state": "Tunis"
  }
}
```

2. **Create second account:**

```json
POST /api/v1/accounts
{
  "first_name": "Nour",
  "last_name": "Hamdi",
  "email": "nour.hamdi@example.tn",
  "phone": "+21655987654",
  "date_of_birth": "1995-07-18T00:00:00Z",
  "password": "nour2024",
  "account_type": "SAVINGS",
  "currency": "TND",
  "address": {
    "street": "Centre Ville",
    "city": "Nabeul",
    "postal_code": "8000",
    "country": "Tunisia",
    "state": "Nabeul"
  }
}
```

3. **Login with first account:**

```json
POST /api/v1/auth/login
{
  "account_number": "{account_number_from_step_1}",
  "password": "sami2024"
}
```

4. **Make a deposit:**

```json
POST /api/v1/transactions/deposit
{
  "account_number": "{sami_account_number}",
  "amount": 200000,
  "currency": "TND",
  "description": "Initial deposit"
}
```

5. **Transfer money:**

```json
POST /api/v1/transactions/transfer
{
  "from_account_number": "{sami_account_number}",
  "to_account_number": "{nour_account_number}",
  "amount": 30000,
  "currency": "TND",
  "description": "Birthday gift",
  "reference": "GIFT-001"
}
```

6. **Check transaction history:**

```
GET /api/v1/transactions/account/{sami_account_id}
```

### Scenario 2: Foreign Currency Testing

1. **Create EUR account:**

```json
POST /api/v1/accounts
{
  "first_name": "Karim",
  "last_name": "Zouari",
  "email": "karim.zouari@business.tn",
  "phone": "+21698123456",
  "date_of_birth": "1982-04-25T00:00:00Z",
  "password": "karim2024",
  "account_type": "BUSINESS",
  "currency": "EUR",
  "address": {
    "street": "Zone Franche",
    "city": "Bizerte",
    "postal_code": "7000",
    "country": "Tunisia",
    "state": "Bizerte"
  }
}
```

2. **Create USD account:**

```json
POST /api/v1/accounts
{
  "first_name": "Olfa",
  "last_name": "Mejri",
  "email": "olfa.mejri@trading.tn",
  "phone": "+21644555666",
  "date_of_birth": "1990-09-12T00:00:00Z",
  "password": "olfa2024",  "account_type": "FOREIGN_CURRENCY",
  "currency": "USD",
  "address": {
    "street": "La Marsa",
    "city": "Tunis",
    "postal_code": "2078",
    "country": "Tunisia",
    "state": "Tunis"
  }
}
```

3. **Login and make EUR transfer:**

```json
POST /api/v1/auth/login
{
  "account_number": "{karim_account_number}",
  "password": "karim2024"
}
```

4. **Make EUR deposit:**

```json
POST /api/v1/transactions/deposit
{
  "account_number": "{karim_account_number}",
  "amount": 5000,
  "currency": "EUR",
  "description": "Business revenue"
}
```

5. **Transfer EUR to USD account:**

```json
POST /api/v1/transactions/transfer
{
  "from_account_number": "{karim_account_number}",
  "to_account_number": "{olfa_account_number}",
  "amount": 1000,
  "currency": "EUR",
  "description": "International business payment",
  "reference": "BIZ-EUR-001"
}
```

## 🚨 Common Error Codes

**400 Bad Request:**

```json
{
  "success": false,
  "error": "Invalid data provided",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

**401 Unauthorized:**

```json
{
  "success": false,
  "error": "Missing or invalid token",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

**403 Forbidden:**

```json
{
  "success": false,
  "error": "Insufficient permissions",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

**404 Not Found:**

```json
{
  "success": false,
  "error": "Account not found",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

**409 Conflict:**

```json
{
  "success": false,
  "error": "Account with this email already exists",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

**422 Unprocessable Entity:**

```json
{
  "success": false,
  "error": "Insufficient balance for this transaction",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

**500 Internal Server Error:**

```json
{
  "success": false,
  "error": "Internal server error",
  "timestamp": "2025-01-31T06:15:30Z"
}
```

## 📊 Fee Structure

### Transaction Fees (in millimes):

- **🔄 Transfers (TND)**: 500 millimes (0.5 TND)
- **🔄 Transfers (EUR)**: 50 euro cents
- **🔄 Transfers (USD)**: 50 cents
- **📥 Deposits**: Free
- **📤 Withdrawals**: 200 millimes (0.2 TND)
- **💱 Currency Exchange**: 1% of transaction amount

### Account Maintenance Fees:

- **🏦 Checking Account**: Free
- **💰 Savings Account**: Free
- **🏢 Business Account**: 5 TND per month
- **💱 Foreign Currency Account**: 10 TND per month

## 🔐 Security Notes

1. **🔑 JWT Tokens** expire after 24 hours
2. **🔄 Refresh tokens** can be used to get new access tokens
3. **🛡️ Passwords** are hashed with bcrypt
4. **📱 Phone numbers** must be in international format (+216...)
5. **📧 Email addresses** must be unique across the system
6. **🏦 Account numbers** follow Tunisian IBAN format (TN59 + 20 digits)

## 🌐 Supported Currencies

- **TND**: Tunisian Dinar (primary currency, millimes precision)
- **EUR**: Euro (cents precision)
- **USD**: US Dollar (cents precision)

## 📝 Response Format

All API responses follow this standard format:

```json
{
  "success": boolean,
  "data": object|array|null,
  "message": string,
  "error": string|null,
  "timestamp": string (ISO 8601)
}
```

---

**📞 Support**: For API support, contact the development team or check the repository issues.

**🏦 Banking Standards**: This API follows Tunisian Central Bank (BCT) regulations and international banking standards.

````

**404 Not Found:**

```json
{
  "error": "Compte non trouvé"
}
````

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
