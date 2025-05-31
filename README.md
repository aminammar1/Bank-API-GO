# 🏦 Bank API

A modern, secure banking API built with Go that follows Tunisian banking standards and regulations.

## ✨ Features

### 🇹🇳 Tunisian Banking Compliance

- **🏧 Tunisian IBAN** format support (TN59 + 20 digits)
- **🏛️ Tunisian BIC** codes (STB, BIAT, BNA, ATB, UBCI)
- **💰 Tunisian Dinar (TND)** as primary currency with millimes precision
- **💱 Multi-currency** support for EUR/USD foreign currency accounts
- **📋 Account types** following Central Bank of Tunisia (BCT) regulations

### 🔐 Security & Authentication

- **🔑 JWT-based authentication** with secure token management
- **🛡️ Password hashing** with bcrypt encryption
- **🚪 Protected endpoints** with middleware authorization
- **🌐 CORS support** for web applications
- **⏰ Token refresh** functionality

### 💳 Transaction Management

- **💸 Multi-currency transactions** (TND, EUR, USD)
- **📊 Automatic fee calculation** based on transaction type
- **⚡ Real-time balance updates** in millimes precision
- **📈 Transaction status tracking** (PENDING, COMPLETED, FAILED)
- **📚 Comprehensive transaction history** with filtering options

### 🔧 API Design

- **🎯 RESTful API design** with clean endpoints
- **✅ Proper HTTP status codes** for all responses
- **📝 Structured JSON responses** with consistent format
- **🔍 Input validation** with detailed error messages
- **⚠️ Error handling** with timestamps and tracking

## 🚀 Quick Start

### 📋 Prerequisites

- **🔧 Go 1.23.0** or higher
- **🐘 PostgreSQL** database
- **📦 Git** version control
- **🐳 Docker & Docker Compose** (for containerized deployment)

### 🛠️ Installation

## Option 1: 🏠 Local Development

1. **📥 Clone the repository**

   ```bash
   git clone <repository-url>
   cd bank-api
   ```

2. **📦 Install dependencies**

   ```bash
   go mod tidy
   ```

3. **⚙️ Set up environment variables**

   **CRITICAL: Copy and configure environment file:**

   ```bash
   # Copy the example file
   copy .env.example .env

   # Edit .env with your SECURE values
   # ⚠️  NEVER use the example values in production!
   ```

   **Required secure configuration in `.env`:**

   ```env
   # IMPORTANT: Generate secure values for production!

   # Database - Use strong credentials
   DB_USER=your_secure_db_user
   DB_PASSWORD=your_very_secure_password_min_16_chars
   DB_NAME=your_database_name

   # JWT - Generate with: openssl rand -hex 32
   JWT_SECRET=your_32_char_minimum_secure_jwt_secret

   # PgAdmin (dev only)
   PGADMIN_DEFAULT_EMAIL=admin@your-company.local
   PGADMIN_DEFAULT_PASSWORD=your_secure_pgadmin_password
   ```

4. **🗄️ Start PostgreSQL database**

   You can use Docker to run PostgreSQL:

   ```bash
   # Uses environment variables from .env file
   docker run --name postgres-bank `
     -e POSTGRES_USER=$env:DB_USER `
     -e POSTGRES_PASSWORD=$env:DB_PASSWORD `
     -e POSTGRES_DB=$env:DB_NAME `
     -p 5434:5432 `
     -d postgres:15-alpine
   ```

   Or make sure PostgreSQL is running with the configured database.

5. **🚀 Run the application**

   ```bash
   make dev
   ```

   The API will be available at `http://localhost:8080`

## Option 2: 🐳 Docker (Recommended)

1. **📥 Clone the repository**

   ```bash
   git clone <repository-url>
   cd bank-api
   ```

2. **🐳 Start with Docker Compose** ```bash

   # Start in background

   make docker-run

   # Or start in development mode (foreground)

   make docker-dev

   ```

   This will:
   - 🐘 Create and start PostgreSQL with initialization data
   - 🏗️ Build and start the banking API
   - 🔧 Optionally start PgAdmin for database administration

   ```

3. **✅ Check status**

   ```bash
   # View logs
   make docker-logs

   # Check API health
   curl http://localhost:8080/api/v1/health
   ```

4. **🛑 Stop services**

   ```bash
   # Stop containers
   make docker-stop

   # Clean everything (removes volumes and data)
   make docker-clean
   ```

### 🐳 Docker Services

The Docker Compose setup starts these services:

- **🏦 bank-api**: Main API on configured port (default: 8080)
- **🐘 postgres**: PostgreSQL database on configured port (default: 5434)
- **🔧 pgadmin**: PostgreSQL admin interface (development profile only)
  - Access via your configured credentials in `.env` file
  - URL: `http://localhost:5050`

### 🔨 Building for Production

```bash
go build -o bin/bank-api cmd/server/main.go
./bin/bank-api
```

## 📚 API Documentation

### 🌐 Base URL

```
http://localhost:8080/api/v1
```

### 🔐 Authentication

Most endpoints require JWT authentication. Include the token in the `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```

### 🛣️ Endpoints

#### 💓 Health Check

```http
GET /api/v1/health
```

Returns the API health status.

#### 🔑 Authentication

##### 📝 Register Account

```http
POST /api/v1/accounts
Content-Type: application/json

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

##### 🔐 Login

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "account_number": "TN59...",
  "password": "securepassword123"
}
```

##### 🔄 Refresh Token

```http
POST /api/v1/auth/refresh
Authorization: Bearer <token>
```

##### 🚪 Logout

```http
POST /api/v1/auth/logout
Authorization: Bearer <token>
```

#### 🏦 Account Management

##### 👤 Get Account by ID

```http
GET /api/v1/accounts/{id}
Authorization: Bearer <token>
```

##### 🔢 Get Account by Account Number

```http
GET /api/v1/accounts/number/{account_number}
Authorization: Bearer <token>
```

##### 👥 Get Accounts by Customer ID

```http
GET /api/v1/accounts/customer/{customer_id}
Authorization: Bearer <token>
```

##### 📋 Get All Accounts (Admin)

```http
GET /api/v1/accounts?limit=10&offset=0
Authorization: Bearer <token>
```

##### ✏️ Update Account

```http
PUT /api/v1/accounts/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "phone": "+1234567891",
  "address": {
    "street": "456 Oak St",
    "city": "Boston",
    "postal_code": "02101",
    "country": "United States",
    "state": "MA"
  }
}
```

##### 🗑️ Delete Account

```http
DELETE /api/v1/accounts/{id}
Authorization: Bearer <token>
```

##### 💰 Get Account Balance

```http
GET /api/v1/accounts/{account_number}/balance
Authorization: Bearer <token>
```

##### 🔧 Update Account Status

```http
PATCH /api/v1/accounts/{id}/status
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "ACTIVE"
}
```

#### 💳 Transaction Management

##### 💸 Transfer Money

```http
POST /api/v1/transactions/transfer
Authorization: Bearer <token>
Content-Type: application/json

{
  "from_account_number": "TN5961705312451143542106",
  "to_account_number": "TN5959238705041140193701",
  "amount": 25000,
  "currency": "TND",
  "description": "Transfer to friend",
  "reference": "TXN-2025-001"
}
```

##### 📥 Deposit Money

```http
POST /api/v1/transactions/deposit
Authorization: Bearer <token>
Content-Type: application/json

{
  "account_number": "TN5961705312451143542106",
  "amount": 100000,
  "currency": "TND",
  "description": "Initial deposit"
}
```

##### 📤 Withdraw Money

```http
POST /api/v1/transactions/withdraw
Authorization: Bearer <token>
Content-Type: application/json

{
  "account_number": "TN5961705312451143542106",
  "amount": 10000,
  "currency": "TND",
  "description": "ATM withdrawal"
}
```

##### 📄 Get Transaction by ID

```http
GET /api/v1/transactions/{id}
Authorization: Bearer <token>
```

##### 📋 Get Account Transactions

```http
GET /api/v1/transactions/account/{account_id}?limit=10&offset=0
Authorization: Bearer <token>
```

##### 📊 Get All Transactions (Admin)

```http
GET /api/v1/transactions?limit=10&offset=0
Authorization: Bearer <token>
```

##### 🔄 Update Transaction Status

```http
PUT /api/v1/transactions/{id}/status
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "COMPLETED"
}
```

### Response Format

#### Success Response

```json
{
  "success": true,
  "data": {
    // Response data
  },
  "message": "Success message",
  "timestamp": "2025-05-31T06:15:30Z"
}
```

#### Error Response

```json
{
  "success": false,
  "error": "Error message",
  "timestamp": "2025-05-31T06:15:30Z"
}
```

### 🏦 Account Types (BCT Compliant)

- `CHECKING` - Standard checking account
- `SAVINGS` - High-yield savings account
- `BUSINESS` - Business banking account
- `FOREIGN_CURRENCY` - Multi-currency account (EUR/USD)

### 💳 Transaction Types

- `TRANSFER` - Transfer between accounts
- `DEPOSIT` - Account deposit
- `WITHDRAWAL` - Account withdrawal
- `PAYMENT` - Payment transaction

### 📊 Transaction Status

- `PENDING` - Awaiting processing
- `COMPLETED` - Successfully completed transaction
- `FAILED` - Failed transaction

### 💰 Supported Currencies

- `TND` - Tunisian Dinar (primary currency)
- `EUR` - Euro (foreign currency accounts)
- `USD` - US Dollar (foreign currency accounts)

## 🧪 Testing

Run the comprehensive test suite:

```bash
go test ./tests/...
```

The test suite includes:

- Account creation and management
- Authentication flow
- Transaction processing
- Error handling
- Edge cases

## 🏗️ Project Structure

```
bank-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── api/
│   │   ├── handlers/            # HTTP request handlers
│   │   ├── middleware/          # Authentication, logging, CORS
│   │   └── routes/              # Route definitions
│   ├── config/                  # Configuration management
│   ├── models/                  # Data models and DTOs
│   ├── repository/              # Data access layer
│   ├── services/                # Business logic layer
│   └── utils/                   # Utility functions
├── tests/                       # Comprehensive test suite          
└── README.md                    # This file
```

## 🔧 Configuration

The application uses environment variables for configuration:

### Server Settings

- `SERVER_PORT` - Server port (default: 3000)
- `SERVER_HOST` - Server host (default: localhost)
- `SERVER_READ_TIMEOUT` - Read timeout (default: 30s)
- `SERVER_WRITE_TIMEOUT` - Write timeout (default: 30s)

### Database Settings

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5433)
- `DB_USER` - Database user (default: bankgo)
- `DB_PASSWORD` - Database password (default: testbank)
- `DB_NAME` - Database name (default: bankdb)
- `DB_SSLMODE` - SSL mode (default: disable)

### JWT Settings

- `JWT_SECRET` - JWT signing secret (required in production)
- `JWT_EXPIRES_IN` - Token expiration time (default: 24h)
- `JWT_ISSUER` - JWT issuer (default: bank-api)

## 🛡️ Security Best Practices

1. **Change the JWT secret** in production
2. **Use HTTPS** in production environments
3. **Implement rate limiting** for API endpoints
4. **Regular security audits** of dependencies
5. **Input validation** on all endpoints
6. **Secure database connections** with SSL in production

## 📈 Performance Considerations

- Database indexes on frequently queried fields
- Connection pooling for database connections
- Pagination for large result sets
- Efficient JSON marshaling/unmarshaling
- Proper HTTP caching headers

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License .

## 🆘 Support

For support and questions, please open an issue in the repository.

---

**Built by Mohamed Amine Ammar with ❤️ using Go and following international banking standards**
