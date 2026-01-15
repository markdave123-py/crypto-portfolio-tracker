# Crypto Portfolio Tracker

A production-style backend service for tracking crypto portfolios across chains, providing live pricing, transaction history, and portfolio valuation.

Built with Go using clean architecture principles and real blockchain data providers.
---

## Features

### Core

- Live token pricing (multi-chain) via CoinGecko

- Transaction history via Etherscan (Ethereum-compatible chains)

- Portfolio management with holdings CRUD


### Engineering 

- Clean modular architecture (handlers / services / repositories)

- Interface-based design with dependency injection

- Caching with TTL (Redis)

- Rate limiting and retry with exponential backoff

- Graceful shutdown

- Structured logging

- Swagger (OpenAPI) documentation

- Docker and docker-compose support

- Comprehensive unit tests
---

## Architecture Overview
This service follows a clean layered architecture:

```text
HTTP (Handlers)
   ↓
Services (Business Logic)
   ↓
Providers / Repositories
   ↓
External APIs (CoinGecko, Etherscan) / Cache
```

### Design Principles

- Separation of concerns (no business logic in handlers)

- Interfaces for testability and extensibility

- Explicit dependency wiring via AppContext

- Stateless API design

### Codebase Structure
```text
.
├── cmd/
│   ├── root.go
│   └── server.go
├── internal/
│   ├── app/
│   ├── cache/
│   ├── config/
│   ├── handlers/
│   ├── httpserver/
│   ├── logger/
│   ├── pricing/
│   │   └── coingecko/
│   ├── transactions/
│   │   └── etherscan/
│   ├── portfolio/
│   └── utils/
├── docs/
├── docker-compose.yml
├── Dockerfile
├── main.go
└── README.md
```

## High-Level Flows

### Pricing Flow

#### POST /prices
 → PricingHandler
 → PricingService
 → Redis Cache
 → CoinGecko Provider (batched, rate-limited)
 → Response


- Cache-aside pattern

- Batch requests per chain

- Retry with exponential backoff

### Transactions Flow
#### GET /wallets/{wallet}/transactions
 → TransactionsHandler
 → TransactionsService
 → Etherscan Provider
 → Classification (send / receive / swap / stake)
 → Filtering + pagination


- Direction detection (in / out)

- Type classification via method signatures

- Status detection (success / failed)

### Portfolio Flow
#### GET /wallets/{wallet}/portfolio
 → PortfolioHandler
 → PortfolioService
 → Holdings Repository
 → PricingService
 → Portfolio valuation


- Holdings CRUD

- Live price-based valuation

- Aggregated portfolio totals

## API Endpoints

### Prices

#### POST /prices

```json
Request:

{
  "assets": [
    {
      "chain": "ethereum",
      "contract_address": "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599"
    }
  ]
}
```

```json
Response:

{
  "success": true,
  "data": {
    "ethereum:0x2260fac5e5542a773aa44fbcfedf7c193bc2c599": 67187.33
  }
}
```

### Transactions

#### GET /wallets/{wallet}/transactions


#### Query parameters:

- chain (required)

- page (default: 1)

- limit (default: 20, max: 100)

- type (send | receive | swap | stake)

- status (success | failed)

- token

- start_date (RFC3339)

- end_date (RFC3339)


### Portfolio

#### GET    /wallets/{wallet}/portfolio
#### POST   /wallets/{wallet}/portfolio/holdings
#### PUT    /wallets/{wallet}/portfolio/holdings
#### DELETE /wallets/{wallet}/portfolio/holdings

## Swagger Documentation

Swagger UI is available at:

http://localhost:8080/swagger/index.html


Includes full request/response schemas, parameters, and examples.


## Configuration

### Required Environment Variables

| Variable          | Description              |
| ----------------- | ------------------------ |
| COINGECKO_API_KEY | CoinGecko API key        |
| ETHERSCAN_API_KEY | Etherscan API key        |
| REDIS_URL         | Redis connection URL     |
| SERVER_PORT       | API port (default: 8080) |

### Running with Docker
```bash
docker-compose up --build
```

#### Services:

- API on localhost:8080

- Redis (internal)

### Running Locally

Start Redis
docker run -p 6379:6379 redis

#### Export environment variables

export COINGECKO_API_KEY=your_key
export ETHERSCAN_API_KEY=your_key
export REDIS_URL=redis://localhost:6379

#### Run the server
```bash
go run . server
```

## Testing

Run all tests:

```bash
go test ./...
```

### Includes:

- Service unit tests

- Provider tests with mocks

- HTTP handler tests

- Security and Reliability

- Rate-limited outbound requests

- Retry with exponential backoff

- Graceful shutdown on SIGINT/SIGTERM

- Structured logging using Zap

## Notes

- PostgreSQL intentionally omitted (not required for assessment)

- Portfolio repository implemented in-memory for simplicity

- Architecture allows easy extension to persistent storage
