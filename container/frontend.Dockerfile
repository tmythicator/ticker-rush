# Build stage
FROM node:20-alpine AS builder

WORKDIR /app

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN corepack enable pnpm && pnpm install --frozen-lockfile

COPY frontend/ .
RUN pnpm run build

# Runtime stage
FROM nginx:alpine

# Force fast shutdown, don't wait for SSE connections to close
STOPSIGNAL SIGTERM

COPY --from=builder /app/dist /usr/share/nginx/html
COPY container/nginx.conf.template /etc/nginx/templates/default.conf.template

CMD ["nginx", "-g", "daemon off;"]
