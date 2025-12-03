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
          ];

          shellHook = ''
            echo "Welcome to Ticker Rush Dev Environment!"
            echo "Run 'process-compose up' to start the stack."
          '';
        };
      }
    );
}
