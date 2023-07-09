{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    gomod2nix_package = {
      url = "github:nix-community/gomod2nix";
      inputs = {
        nixpkgs.follows = "nixpkgs";
      };
    };
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
          default = pkgs.${system}.buildGoApplication
            {
              pname = "asciiquarium.live";
              version = "0.1";
              src = ./.;
              modules = ./gomod2nix.toml;

              nativeBuildInputs = with pkgs.${system}; [ makeWrapper ];

              installPhase = ''
                runHook preInstall

                mkdir -p $out
                dir="$GOPATH/bin"

                [ -e "$dir" ] && cp -r $dir $out

                wrapProgram $out/bin/asciiquarium.live \
                  --prefix PATH : ${pkgs.${system}.lib.makeBinPath [ pkgs.${system}.asciiquarium ]}

                runHook postInstall
              '';
            };
        });

      devShells = forAllSystems (system: {
        default = pkgs.${system}.mkShell {
          buildInputs = with pkgs.${system}; [
            gomod2nix
            asciiquarium
            (mkGoEnv { pwd = ./.; })
          ];
        };
      });
    };
}
