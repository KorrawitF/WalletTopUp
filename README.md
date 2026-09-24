# Wallet Top-Up System

A Go (Gin) REST service for topping up user wallets in two steps:

1. **Verify**: creates a top-up transaction (`verified`) that expires after 10 minutes.
2. **Confirm**: adds the amount to the wallet in a DB transaction and marks it `completed`, or `expired` if the time has run out.

**Stack:** Go 1.25, Gin, GORM + PostgreSQL, Redis (optional transaction cache), Docker Compose.

## API

| Method | Path              | Body                                                               |
| ------ | ----------------- | ------------------------------------------------------------------ |
| POST   | `/wallet/verify`  | `{"user_id": 1, "amount": 100, "payment_method": "credit_card"}`   |
| POST   | `/wallet/confirm` | `{"transaction_id": "<id>"}`                                       |

Payment methods: `credit_card`, `debit_card`, `bank_transfer`.

## Project Structure

```text
main.go                 # entrypoint, dependency wiring, graceful shutdown
internal/
  config/               # env config loading
  domain/               # entities, repository/cache interfaces, errors
  service/              # business logic (verify / confirm top-up)
  infra/
    database/           # PostgreSQL connection, models, repositories
    cache/              # Redis cache (with no-op fallback)
  interface/http/       # Gin router, handlers, DTOs, middleware
pkg/                    # shared logger, utils, constants
```

## Getting Started

**Docker (recommended)**

```bash
docker compose up --build
```

The API runs on `http://localhost:8000`, and the database is migrated automatically.

**Local** (needs PostgreSQL running, plus Redis if `CACHE_ENABLED=true`)

```bash
cp .env.example .env
go run . -db migrate
```

**Tests**

```bash
go test ./...
```
