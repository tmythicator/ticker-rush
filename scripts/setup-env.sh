#!/usr/bin/env bash

# setup-env.sh: Source this script to configure the Ticker Rush dev environment.
# This is called automatically by flake.nix shellHook.

# 1. Docker Socket Detection
possible_sockets=(
    "$HOME/.colima/default/docker.sock" "$HOME/.colima/docker.sock"
    "$HOME/.orbstack/run/docker.sock" "$HOME/.docker/run/docker.sock"
    "/run/user/$UID/docker.sock" "/var/run/docker.sock"
)

for sock in "${possible_sockets[@]}"; do
    if [ -S "$sock" ]; then
        export DOCKER_HOST="unix://$sock"
        export TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE="/var/run/docker.sock"
        break
    fi
done

# 2. Database Initialization
if [ ! -d "$PGDATA" ]; then
    echo "Initializing Postgres data in $PGDATA..."
    initdb -U postgres --no-locale --encoding=UTF8 > /dev/null
fi

# 3. Gather versions for welcome message
export GO_VERSION=$(go version | awk '{print $3}')
export NODE_VERSION=$(node -v)
export PNPM_VERSION=$(pnpm -v)
export PYTHON_VERSION=$(python3 --version | awk '{print $2}')
