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

          openbao-plugin-secrets-nebula =
            let
              version = "2.10.1";
            in
            pkgs.buildGo127Module {
              pname = "openbao-plugin-secrets-nebula";
              inherit version; # Bumping to v2 for Nebula certs

              src = ./.;

              # We target the main command package specifically
              subPackages = [ "cmd/openbao-plugin-secrets-nebula" ];

              # Stamp the binary version so the nix-built artifact
              # self-identifies, matching goreleaser's -X main.version.
              ldflags = [
                "-s"
                "-w"
                "-X main.version=${version}"
              ];

              # Nix requires the vendor hash to guarantee reproducibility.
              # Leave this as fakeHash for the first build. It will fail and
              # give you the real hash, which you will paste here.
              vendorHash = "sha256-69vxgXHFPtkiWxUH8yjmWEs3U3+p+QKNKf4RAzRLSm0=";

              meta = with pkgs.lib; {
                description = "OpenBao Secrets Engine for Nebula PKI";
                homepage = "https://github.com/suorcd/openbao-plugin-secrets-nebula";
                license = licenses.mit;
                maintainers = [ ];
              };
            };
        };

        devShells.default = pkgs.mkShell {
          # This gives you a perfect local dev environment when you run `nix develop`
          buildInputs = with pkgs; [
            go_1_27
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
