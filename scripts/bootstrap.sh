#!/usr/bin/env bash

# bootstrap.sh: Initialize the development environment for Ticker Rush.
# This script ensures Nix and Direnv are set up correctly.

set -e

echo "Bootstrapping Ticker Rush development environment..."

# 1. Check for Nix
if ! command -v nix &> /dev/null; then
    echo "Nix is not installed. Please install it from https://nixos.org/download.html"
    exit 1
fi
echo "Nix found."

# 2. Check for Direnv
if ! command -v direnv &> /dev/null; then
    echo "direnv is not installed. It is highly recommended for automatic environment loading."
    echo "Install it via your package manager or 'nix-env -iA nixpkgs.direnv'"
fi

# 3. Initialize Direnv if present
if [ -f .envrc ]; then
    echo "Allowing direnv..."
    direnv allow || true
fi

# 4. Check for Flakes support
if ! nix flake --help &> /dev/null; then
    echo "Nix Flakes are not enabled. Ensure 'experimental-features = nix-command flakes' is in your nix.conf"
fi

echo "Entering Nix development shell to verify tools..."
nix develop --command bash -c "echo 'All tools are ready!'; exit"

echo "Done! You can now run 'nix develop' or allow direnv to start coding."
