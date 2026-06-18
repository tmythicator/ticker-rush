# Test Blocked Vulnerability Scans & Config Probes
# All these must return HTTP 403 Forbidden under the strict allowlist rules

# 1. Admin/DB scan targets
GET {{host}}/wp-admin
HTTP 403

GET {{host}}/wp-login.php
HTTP 403

GET {{host}}/phpmyadmin
HTTP 403

# 2. Config & Lock files
GET {{host}}/firebase.json
HTTP 403

GET {{host}}/package.json
HTTP 403

GET {{host}}/pnpm-lock.yaml
HTTP 403

# 3. Environment & Git folders
GET {{host}}/.env
HTTP 403

GET {{host}}/.git
HTTP 403

# 4. Scanner targets (like webui)
GET {{host}}/webui
HTTP 403

# 5. General random path (any path not explicitly allowed)
GET {{host}}/random-vulnerability-probe-path
HTTP 403
