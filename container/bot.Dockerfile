FROM oven/bun:alpine

WORKDIR /app

COPY bot/package.json bot/bun.lock* ./
RUN bun install --frozen-lockfile

COPY bot/ .

EXPOSE 3000

CMD ["bun", "run", "index.ts"]
