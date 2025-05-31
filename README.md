# Tunisian Bank API

A modern banking API built with Go that follows Tunisian banking standards

## Features

### Tunisian Banking Compliance

- **Tunisian IBAN** format support (TN59 + 20 digits)
- **Tunisian BIC** codes (STB, BIAT, BNA, ATB, UBCI)
- **Tunisian Dinar (TND)** as primary currency with millimes precision
- Support for EUR/USD foreign currency accounts
- Tunisian account types following BCT regulations

### Security & Authentication

- JWT-based authentication
- Password hashing with bcrypt
- Protected endpoints with middleware
- CORS support

### Transaction Management

- Multi-currency support (TND, EUR, USD)
- Transaction fees calculation
- Real-time balance updates in millimes
- Transaction status tracking
- Comprehensive transaction history

### API Design

- RESTful API design
- Proper HTTP status codes
- Structured JSON responses in French/Arabic context
- Input validation
- Error handling with timestamps

## 🚀 Quick Start

### Prerequisites

- Go 1.23.0 or higher
- PostgreSQL database
- Git

### Installation

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd bank-api
   ```

2. **Install dependencies**

   ```bash
   go mod tidy
   ```

3. **Set up environment variables**

   Create a `.env` file or set environment variables:

   ```env
   # Server Configuration
   SERVER_PORT=3000
   SERVER_HOST=localhost   # Database Configuration
   DB_HOST=localhost
   DB_PORT=5433
   DB_USER=bankgo
   DB_PASSWORD=testbank
   DB_NAME=bankdb_tunisia
   DB_SSLMODE=disable

   # JWT Configuration
   JWT_SECRET=your-secret-key-change-in-production
   JWT_EXPIRES_IN=24h
   JWT_ISSUER=tunisia-bank-api
   ```

4. **Start PostgreSQL database**

   You can use Docker to run PostgreSQL:

   ```bash
   docker run --name postgres-bank \
     -e POSTGRES_USER=bankgo \
     -e POSTGRES_PASSWORD=testbank \
     -e POSTGRES_DB=bankdb_tunisia \
     -p 5433:5432 \
     -d postgres:15-alpine
   ```

   Or make sure PostgreSQL is running with the configured database.

5. **Run the application**

   ```bash
   go run cmd/server/main.go
   ```

   The API will be available at `http://localhost:3000`

### Building for Production

```bash
go build -o bin/bank-api cmd/server/main.go
./bin/bank-api
```

## 📚 API Documentation

### Base URL

```
http://localhost:3000/api/v1
```

### Authentication

Most endpoints require JWT authentication. Include the token in the `Authorization` header:

```
Authorization: Bearer <your-jwt-token>
```

### Endpoints

#### Health Check

```http
GET /health
```

Returns the API health status.

#### Authentication

##### Register Account

```http
POST /auth/register
Content-Type: application/json

{
  "first_name": "Mohamed",
  "last_name": "Ben Ahmed",
  "email": "mohamed.benahmed@example.tn",
  "phone": "+21612345678",
  "date_of_birth": "1990-01-15",
  "password": "motdepasse123",
  "account_type": "COMPTE_COURANT",
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

##### Login

```http
POST /auth/login
Content-Type: application/json

{
  "account_number": "TN59...",
  "password": "motdepasse123"
}
```

#### Account Management

##### Get Account by ID

```http
GET /accounts/{id}
Authorization: Bearer <token>
```

##### Get Account by Account Number

```http
GET /accounts/number/{account_number}
Authorization: Bearer <token>
```

##### Get Accounts by Customer ID

```http
GET /accounts/customer/{customer_id}
Authorization: Bearer <token>
```

##### Get All Accounts (Admin)

```http
GET /accounts?limit=10&offset=0
Authorization: Bearer <token>
```

##### Update Account

```http
PUT /accounts/{id}
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

##### Delete Account

```http
DELETE /accounts/{id}
Authorization: Bearer <token>
```

#### Transaction Management

##### Create Transaction

```http
POST /transactions
Authorization: Bearer <token>
Content-Type: application/json

{
  "from_account_number": "ACC001234567890",
  "to_account_number": "ACC987654321098",
  "amount": 10000,
  "currency": "USD",
  "description": "Payment for services",
  "reference": "INV-2023-001"
}
```

##### Get Transaction by ID

```http
GET /transactions/{id}
Authorization: Bearer <token>
```

##### Get Account Transactions

```http
GET /transactions/account/{account_id}?limit=10&offset=0
Authorization: Bearer <token>
```

##### Get All Transactions (Admin)

```http
GET /transactions?limit=10&offset=0
Authorization: Bearer <token>
```

##### Update Transaction Status

```http
PUT /transactions/{id}/status
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

### Account Types (BCT Compliant)

- `COMPTE_COURANT` - Compte courant (checking account)
- `COMPTE_EPARGNE` - Compte d'épargne (savings account)
- `COMPTE_ENTREPRISE` - Compte entreprise (business account)
- `COMPTE_DEVISES` - Compte en devises (foreign currency account)

### Transaction Types

- `TRANSFER` - Virement entre comptes
- `DEPOSIT` - Dépôt sur compte
- `WITHDRAWAL` - Retrait du compte
- `PAYMENT` - Paiement

### Transaction Status

- `PENDING` - En attente de traitement
- `COMPLETED` - Transaction terminée avec succès
- `FAILED` - Échec de la transaction

### Supported Currencies

- `TND` - Dinar Tunisien (currency principale)
- `EUR` - Euro (comptes en devises)
- `USD` - Dollar Américain (comptes en devises)

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
├── old_files/                   # Legacy files (for reference)
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

This project is licensed under the MIT License - see the LICENSE file for details.

## 🆘 Support

For support and questions, please open an issue in the repository.

---

**Built with ❤️ using Go and following international banking standards**
