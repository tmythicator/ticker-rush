#!/usr/bin/env bash

set -eo pipefail

readonly CONTAINER_NAME="nginx-test-suite"
readonly TEST_PORT=8080
readonly SITE_DOMAIN="localhost"
readonly SERVER_PORT=8081
readonly TEMP_CONF_FILE="container/nginx-test.run.conf"
readonly TEMP_HTML_DIR="container/test-html.run"

for cmd in docker hurl envsubst; do
  if ! command -v "$cmd" &>/dev/null; then
    echo "Error: Required command '$cmd' is not installed." >&2
    exit 1
  fi
done

cleanup() {
  echo "Cleaning up test container and temporary files..."
  docker rm -f "$CONTAINER_NAME" >/dev/null 2>&1 || true
  rm -f "$TEMP_CONF_FILE" || true
  rm -rf "$TEMP_HTML_DIR" || true
}
trap cleanup EXIT

# 1. Generate Nginx config from template
echo "Generating temporary Nginx configuration..."
export SITE_DOMAIN
export SERVER_PORT

# envsubst is scoped to only substitute site domain and server port.
# sed converts SSL server block to listen on port 8080 and comments out SSL directives.
envsubst '$SITE_DOMAIN,$SERVER_PORT' < container/nginx.conf.template \
  | sed "s/listen 443 ssl;/listen ${TEST_PORT};/g" \
  | sed 's/ssl_certificate/#ssl_certificate/g' \
  | sed 's/ssl_prefer_server_ciphers/#ssl_prefer_server_ciphers/g' \
  | sed 's/ssl_protocols/#ssl_protocols/g' \
  | sed 's/ssl_ciphers/#ssl_ciphers/g' \
  > "$TEMP_CONF_FILE"

# 2. Prepare flat HTML assets directory to avoid nested Docker volume mount placeholders on host
mkdir -p "$TEMP_HTML_DIR"
cp -r frontend/public/* "$TEMP_HTML_DIR/"
cp frontend/index.html "$TEMP_HTML_DIR/"

# 3. Run temporary Nginx container
echo "Starting test Nginx container..."
docker run -d \
  --name "$CONTAINER_NAME" \
  --add-host "exchange:127.0.0.1" \
  -p "${TEST_PORT}:${TEST_PORT}" \
  -v "$(pwd)/${TEMP_CONF_FILE}:/etc/nginx/conf.d/default.conf" \
  -v "$(pwd)/${TEMP_HTML_DIR}:/usr/share/nginx/html" \
  nginx:alpine

# Verify Nginx container started successfully
sleep 1.5
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
  echo "Error: Nginx container failed to start. Logs:" >&2
  docker logs "$CONTAINER_NAME" || true
  exit 1
fi

# 3. Run Hurl integration tests
echo "Running Hurl integration tests..."
hurl --test --variable host="http://localhost:${TEST_PORT}" \
  container/nginx_good_routes.hurl \
  container/nginx_bad_routes.hurl

echo "SUCCESS: All Nginx tests passed!"
