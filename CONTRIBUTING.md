# Contributing to Ticker Rush

Ticker Rush is more than just a trading simulator; it's a project dedicated to making the complex world of markets accessible, enjoyable, and risk-free. Our mission is to build a high-performance environment where people can have fun while learning market dynamics, and eventually compete for supremacy in a dedicated human vs. AI division.

Thank you for your interest in helping us achieve this vision! This document provides guidelines for setting up your development environment and submitting contributions.

## Development Environment

We use **Nix** with Flakes to provide a hermetic, reproducible development environment. This ensures all contributors use the same versions of tools like Go, Node.js, and databases.

### 1. Using Nix (Recommended)

1.  **Enter the development shell**:

    ```bash
    nix develop
    ```

    This will automatically install all dependencies: Go, Node.js, Valkey, Postgres, Task, etc.

2.  **Configure environment**:

    ```bash
    cp .env-example .env
    # Add your FINNHUB_API_KEY
    ```

3.  **Start the full stack**:
    ```bash
    task dev
    ```

### 2. Using Docker

If you prefer Docker, you can use the provided `docker-compose.yml`:

```bash
cp .env-example .env
docker-compose up -d --build
```

### 3. Manual Setup (No Nix)

If you prefer manual installation:

- **Go**: 1.25+
- **Node.js**: 18+
- **Valkey/Redis**: port 6379
- **Postgres**: port 5432

Run services:

```bash
# Term A: Fetcher
cd backend && go run ./cmd/fetcher
# Term B: Exchange API
cd backend && go run ./cmd/exchange
# Term C: Frontend
cd frontend && pnpm install && pnpm run dev
```

## Contribution Guidelines

### Branching & PRs

- Create a feature branch from `main`.
- Ensure all tests pass (`task test`).
- Run linters (`task lint`).
- Submit a Pull Request with a clear description of your changes.

### Coding Style

- **Backend (Go)**: Follow standard `gofmt` and `golangci-lint`.
- **Frontend (React)**: Use TypeScript and follow established patterns in the codebase.
- **Licenses**: Ensure any new files include the AGPLv3 header.

## License

By contributing to Ticker Rush, you agree that your contributions will be licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**.
