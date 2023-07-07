{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    gomod2nix_package.url = "github:nix-community/gomod2nix";
  };

  outputs = { nixpkgs, gomod2nix_package, ... }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
      pkgs = forAllSystems (system: (import nixpkgs {
        inherit system;
        overlays = [ gomod2nix_package.overlays.default ];
      }));
    in
    {
      packages = forAllSystems
        (system: {
          default =  pkgs.${system}.buildGoApplication
            {
              pname = "asciiquarium.live";
              version = "0.1";
              src = ./.;
              modules = ./gomod2nix.toml;
            };
        });

      devShells = forAllSystems (system: {
        default = pkgs.${system}.mkShell {
          buildInputs = with pkgs.${system}; [
            gomod2nix
            (mkGoEnv { pwd = ./.; })
          ];
        };
      });
    };
}
