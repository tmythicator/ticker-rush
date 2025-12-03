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
        
        startupScript = ''
          echo "Welcome to Ticker Rush Dev Environment!"

          if [ -f .env ]; then
            set -a
            source .env
            set +a
            echo "✅ Loaded environment variables from .env"
          else
            echo "⚠️ .env not found. Please create it if you need environment variables."
          fi

          echo "Run 'process-compose up' to start the stack."
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

            # Frontend
            nodejs_20
            nodePackages.typescript-language-server

            # Infrastructure
            valkey
            process-compose
          ];

          shellHook = startupScript;
        };
      }
    );
}
