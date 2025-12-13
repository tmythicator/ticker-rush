{
  description = "Ticker Rush: Exchange & Fetcher";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        
        # Use writeText to create a sourcable script (no bin wrapper)
        setup-env = pkgs.writeText "setup-env.sh" ''
          # Define potential Docker socket paths. Needed for `task test:backend` to work
          possible_sockets=(
            "$HOME/.colima/default/docker.sock" # Colima Default
            "$HOME/.colima/docker.sock"         # Colima Legacy
            "$HOME/.orbstack/run/docker.sock"   # OrbStack
            "$HOME/.docker/run/docker.sock"     # Docker Desktop (Mac/Linux)
            "/run/user/$UID/docker.sock"        # Rootless Linux
            "/var/run/docker.sock"              # Standard Linux / System / WSL2
          )

          # Find the first existing socket
          for sock in "''${possible_sockets[@]}"; do
            if [ -S "$sock" ]; then
              export DOCKER_HOST="unix://$sock"
              echo "Found Docker socket: $sock"
              break
            fi
          done
          
          # Fix for Testcontainers + Colima (Ryuk socket mount)
          export TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE="/var/run/docker.sock"
          
          # Setup Postgres
          export PGDATA="$PWD/.data/postgres"
          if [ ! -d "$PGDATA" ]; then
            echo "Initializing Postgres data..."
            initdb -U postgres --no-locale --encoding=UTF8 > /dev/null
          fi
        '';
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Backend
            go
            gopls
            delve
            golangci-lint

            # Infrastructure
            valkey
            process-compose

            # Frontend
            nodejs_20
            nodePackages.pnpm

            # Protobuf
            protobuf
            protoc-gen-go
            protoc-gen-go-grpc

            # Database
            sqlc
            goose
            postgresql_16

            # Task Runner
            go-task

            # Trading Bot
            python3
            python3Packages.grpcio
            python3Packages.grpcio-tools
            python3Packages.protobuf
          ];

          shellHook = ''
            source ${setup-env}
            echo "Welcome to Ticker Rush Dev Environment!"
            echo "Run 'task dev' to start the stack."
          '';
        };
      }
    );
}
