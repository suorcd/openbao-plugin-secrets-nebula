{
  description = "OpenBao Secrets Engine for Nebula";

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
      in
      {
        packages = rec {
          default = openbao-plugin-secrets-nebula;

          openbao-plugin-secrets-nebula = pkgs.buildGoModule {
            pname = "openbao-plugin-secrets-nebula";
            version = "2.0.3"; # Bumping to v2 for Nebula certs

            src = ./.;

            # We target the main command package specifically
            subPackages = [ "cmd/openbao-plugin-secrets-nebula" ];

            # Nix requires the vendor hash to guarantee reproducibility.
            # Leave this as fakeHash for the first build. It will fail and
            # give you the real hash, which you will paste here.
            vendorHash = "sha256-QFiBVIaI+xXKyXpQqLwPgYc051TbUjcboQUPEkI9sqE=";

            meta = with pkgs.lib; {
              description = "OpenBao Secrets Engine for Nebula PKI";
              homepage = "https://github.com/suorcd/openbao-plugin-secrets-nebula";
              license = licenses.mpl20;
              maintainers = [ ];
            };
          };
        };

        devShells.default = pkgs.mkShell {
          # This gives you a perfect local dev environment when you run `nix develop`
          buildInputs = with pkgs; [
            go # Latest stable Go from nixpkgs
            goreleaser # For building your Github releases
            gnumake
            gotools
            golangci-lint
          ];

          shellHook = ''
            echo "OpenBao Nebula Plugin Dev Environment loaded."
            echo "Go version: $(go version)"
          '';
        };
      }
    );
}
