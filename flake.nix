{
  description = "Ticker Rush: Exchange & Fetcher";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # --- Environment Configuration ---
        env-config = {
          PGDATA = "$PWD/.data/postgres";
          GOPATH = "$PWD/.data/go";
          GOMODCACHE = "$PWD/.data/go/pkg/mod";
        };

        # --- Toolsets ---
        backend-tools = with pkgs; [
          go
          gopls
          delve
          golangci-lint
        ];
        frontend-tools = with pkgs; [
          nodejs_20
          nodePackages.pnpm
        ];
        proto-tools = with pkgs; [
          buf
          protobuf
          protoc-gen-go
          protoc-gen-go-grpc
        ];
        db-tools = with pkgs; [
          sqlc
          goose
          postgresql_16
          valkey
        ];
        infra-tools = with pkgs; [
          process-compose
          docker-compose
          go-task
          bun
        ];

      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = backend-tools ++ frontend-tools ++ proto-tools ++ db-tools ++ infra-tools;

          shellHook = ''
            # Export environment config for the setup script
            export PGDATA="${env-config.PGDATA}"
            export GOPATH="${env-config.GOPATH}"
            export GOMODCACHE="${env-config.GOMODCACHE}"

            # Source the setup script from the project folder
            if [ -f ./scripts/setup-env.sh ]; then
              source ./scripts/setup-env.sh
            fi

            echo ""
            echo "Welcome to Ticker Rush Dev Shell"
            echo "Go: $GO_VERSION | Node: $NODE_VERSION | pnpm: $PNPM_VERSION | Bun: $BUN_VERSION"
            echo "Run 'task dev' to start the stack"
            echo ""
          '';
        };
      }
    );
}
