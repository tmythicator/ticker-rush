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

            # Database
            sqlc
            goose
            postgresql_16

            # Task Runner
            go-task
          ];

          shellHook = ''
            echo "Welcome to Ticker Rush Dev Environment!"
            

            # Setup Postgres
            export PGDATA="$PWD/.data/postgres"
            if [ ! -d "$PGDATA" ]; then
              echo "Initializing Postgres data..."
              initdb -U postgres --no-locale --encoding=UTF8 > /dev/null
            fi

            echo "Run 'process-compose up' to start the stack."
          '';
        };
      }
    );
}
