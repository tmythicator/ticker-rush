# Contributing to Ticker Rush

Ticker Rush is a high-performance, real-time stock trading simulator. Our mission is to build a hermetic, scalable, and enjoyable platform for learning market dynamics. Eventually, we aim to establish a Human vs. AI Division, where trading bots compete on an equal footing with human players.

---

## Development Environment

We use Nix to provide a fully reproducible, "zero-config" development environment.

### 1. Startup Flow (Nix)

1. **Enter the shell**:

   ```bash
   nix develop
   ```

   This initializes everything in the project-local .data/ directory, including an isolated Go workspace.

2. **Initialize the project** (First time or after clean):

   ```bash
   task init
   ```

   This handles dependency installation (Go, Node, Bun) and code generation.

3. **Configure Environment**:

   ```bash
   cp .env-example .env
   # Add your FINNHUB_API_KEY to .env
   ```

4. **Start the stack**:
   ```bash
   task dev
   ```
   This command starts databases, runs migrations, and launches all services.

### 2. Manual Setup

We strongly recommend Nix, but if you prefer manual setup:

- Go 1.25+: Backend API & Exchange logic.
- Node.js 20+ / pnpm: Frontend development.
- Bun: Bot execution and TS scripting.
- Postgres & Valkey: Core persistence and caching.

---

## Technical Standards

### Data Layer & State

- **Server State**: Use TanStack Query (React Query) for all frontend data.
- **Real-time**: We use a hybrid architecture (Initial fetch via Query + background updates via SSE).
- **Protobuf**: Protocol Buffers are the Single Source of Truth. All API contracts must be defined in proto/v1/ before implementation.

### Architecture

- **Backend**: Go services built with the Gin framework.
- **Authentication**: HttpOnly Cookies + JWT (stateless, secure against XSS).
- **AI Division**: Trading bots are written in TypeScript using Bun and LangChain.

---

## Contribution Workflow

### 1. Code Quality

- **Formatting**: Run task format before committing.
- **Linting**: Ensure task lint passes (Go & TypeScript).
- **Type Safety**: Never use any. Use generated Protobuf types for all cross-service communication.

### 2. Submitting Changes

- Create a branch for your feature/fix.
- Ensure all tests pass with task test.
- Submissions must be under the AGPL-3.0 License.

Thank you for building the future of trading with us!
