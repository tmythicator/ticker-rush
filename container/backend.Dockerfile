# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

# --- Build Exchange Stage ---
FROM builder AS build-exchange
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/exchange ./cmd/exchange

# --- Build Fetcher Stage ---
FROM builder AS build-fetcher
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/fetcher ./cmd/fetcher

# --- Build Migrate Stage ---
FROM builder AS build-migrate
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /bin/migrate ./cmd/migrate

# --- Final Exchange Image ---
FROM gcr.io/distroless/static-debian12:latest AS exchange-image
WORKDIR /app
COPY --from=build-exchange /bin/exchange .
USER nonroot:nonroot
CMD ["./exchange"]

# --- Final Fetcher Image ---
FROM gcr.io/distroless/static-debian12:latest AS fetcher-image
WORKDIR /app
COPY --from=build-fetcher /bin/fetcher .
USER nonroot:nonroot
CMD ["./fetcher"]

# --- Final Migrate Image ---
FROM gcr.io/distroless/static-debian12:latest AS migrate-image
WORKDIR /app
COPY --from=build-migrate /bin/migrate .
USER nonroot:nonroot
CMD ["./migrate"]