{
  description = "yamusic-tui — terminal client for Yandex Music (alpineQ fork with caching)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        packages.default = pkgs.buildGoModule {
          pname = "yamusic-tui";
          version = "0.1.0-fork";

          src = ./.;

          vendorHash = "sha256-ZacMhp67tFv7f9DbF1c+6oxQeZuCWCEamZ7CqgXMca8=";

          subPackages = [ "." ];

          nativeBuildInputs = [ pkgs.pkg-config ];
          buildInputs = [ pkgs.alsa-lib ];

          ldflags = [ "-s" "-w" ];

          meta = with pkgs.lib; {
            description = "Terminal client for Yandex Music (fork with caching)";
            homepage = "https://github.com/alpineQ/yamusic-tui";
            license = licenses.gpl3Only;
            mainProgram = "yamusic-tui";
            platforms = platforms.linux;
          };
        };

        devShells.default = pkgs.mkShell {
          packages = [ pkgs.go pkgs.gopls pkgs.alsa-lib pkgs.pkg-config ];
        };
      });
}
