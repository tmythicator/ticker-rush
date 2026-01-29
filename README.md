# Ticker Rush

![Ticker Rush Overview](assets/ticker-rush.png)

[![Backend Status](https://github.com/tmythicator/ticker-rush/actions/workflows/backend-test.yml/badge.svg)](https://github.com/tmythicator/ticker-rush/actions/workflows/backend-test.yml)
[![Frontend Status](https://github.com/tmythicator/ticker-rush/actions/workflows/frontend-test.yml/badge.svg)](https://github.com/tmythicator/ticker-rush/actions/workflows/frontend-test.yml)
[![Built With Nix](https://img.shields.io/badge/Built_With-Nix-5277C3.svg?logo=nixos&logoColor=white)](https://nixos.org)

## Overview

Ticker Rush is a high-performance, real-time stock trading simulator. The project's core mission is to make the experience of trading stocks enjoyable and risk-free. It provides a platform where users can have fun while learning market dynamics without financial exposure.

Participants start with equal virtual budgets and trade stocks and ETFs. The primary objective is to maximize portfolio value and compete for the top of the leaderboard.

## Project Scope & Roadmap

Ticker Rush is evolving from a core simulator to a fully featured trading ecosystem.

### Phase 1: Core Simulator (MVP)

- [x] High-performance trade execution engine (Go).
- [x] Real-time market data streaming via SSE.
- [x] Secure authentication (HttpOnly Cookie + JWT).
- [x] Versioned Protobuf API contracts (v1).
- [x] Basic Interactive Dashboard (React).

### Phase 2: Advanced Trading & Social [Current]

- [ ] Global Leaderboards & Social Profiles.
- [ ] Multi-asset support (Options, Crypto).
- [ ] Portfolio Performance Analytics.

### Phase 3: AI-Driven Ecosystem

- [ ] Competitive AI Trading Agents (LangChain + Bun).
- [ ] Natural Language Trading Interface.
- [ ] Automated Strategy Backtesting.

## Architecture

Ticker Rush uses a **Microservices Architecture within a Monorepo** setup:

- **Structure**: A single repository (Monorepo) holds all services, ensuring atomic changes across the stack and shared tooling (Nix, Taskfile).
- **Execution**: At runtime, services act as independent **Microservices** orchestrated via `docker-compose` (Production) or `process-compose` (Development).
- **Backend**: Go services built with the [Gin](https://gin-gonic.com/) web framework.
  - **API/Exchange**: High-performance trading logic and user management.
  - **Data Fetcher**: Real-time market data ingestion and processing.
  - **Persistence**: PostgreSQL using [sqlc](https://sqlc.dev/) and [goose](https://github.com/pressly/goose).
  - **Real-time**: Server-Sent Events (SSE) for sub-second updates.
- **Bot** (WIP): A lightweight trading agent built with [Bun](https://bun.sh/), TypeScript, and [LangChain](https://js.langchain.com/).
- **Frontend**: Modern React 19 + TypeScript + TanStack Query.
- **Infrastructure**: [Valkey](https://valkey.io/) for caching and Nix Flakes for reproducibility.

## Authentication

Ticker Rush implements a secure **HttpOnly Cookie + JWT** authentication strategy:

- **JWT**: Used for stateless session management.
- **HttpOnly Cookies**: Prevents XSS attacks by keeping tokens inaccessible to client-side scripts.
- **Stateless**: The backend remains lean and scalable.

## Protobuf & Type Safety

We use **Protocol Buffers (Protobuf)** as the single source of truth for all API contracts:

- **Source**: Definitions are located in versioned subdirectories within [`proto/`](./proto) (e.g., `proto/exchange/v1/`).
- **Generation**: We use [Buf](https://buf.build/) to generate type-safe code for all languages:
  - **Go**: Generated in `backend/internal/proto/`.
  - **TypeScript (Frontend)**: Generated in `frontend/src/lib/proto/`.
  - **TypeScript (Bot)**: Generated in `bot/src/lib/proto/`.
- **Consistency**: This ensures that a change in the schema is automatically propagated across the entire stack, preventing runtime type mismatches.

## Quick Start

The recommended way to start development is using [Nix](https://nixos.org):

1. **Bootstrap** (First time only):
   ```bash
   chmod +x scripts/bootstrap.sh
   ./scripts/bootstrap.sh
   ```
2. **Develop**:
   ```bash
   nix develop
   task dev
   ```

For detailed configuration and testing guidelines, refer to [CONTRIBUTING.md](./CONTRIBUTING.md).

## License

GNU Affero General Public License v3.0 (AGPL-3.0) — see [LICENSE](./LICENSE).

AGPLv3 ensures that Ticker Rush remains open-source. If you run a modified version as a network service, you must share the corresponding source code with your users.
