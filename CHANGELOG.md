# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- CI workflow (`test.yml`) running on push to `feature/main-v2-dev` and on pull requests: `go mod tidy -diff`, build, vet, test, `goreleaser check`, plus a `nix build` job that guards the flake `vendorHash`

### Changed
- Release notes no longer exclude `chore:` commits; that filter had left the Changelog section of published releases empty

### Fixed
- Backfilled the v2.10.1 release notes with the change list

## [v2.10.1] - 2026-10-09

### Changed
- Upgraded OpenBao deps to `api/v2 v2.7.1` and `sdk/v2 v2.7.1` (includes the SDK audit-log plaintext leak fix GHSA-8xxq-mq9m-xmhw, present in sdk v2.6.3+), plus transitive module updates from `go get -u`
- Bumped `slackhq/nebula` dependency to v1.11.2 (cert API unchanged, no source changes)
- Bumped `golang.org/x/crypto` to v0.58.0
- Bumped `.goversion` to 1.27.2; release workflow moved to `actions/setup-go@v7` and now reads the pin via `go-version-file`
- Updated flake inputs and Nixpkgs pin

### Fixed
- README download URLs pointed at the nonexistent `v2.1.0` release (404); they now reference the real tag
- CHANGELOG reconciled with the actual release tags (`v2.10.0`, previously recorded as `v2.1.0`; added missing `v2.0.5`)
- Release-notes footer compared against the upstream `mkrauser` repo where the v2.x tags don't exist; it now compares on `suorcd`. The install snippet also now matches the README (`bao plugin register ... secret openbao-plugin-secrets-nebula`)
- The `-X main.version` ldflag had no target; the plugin now declares a `version` variable and logs it at startup
- LICENSE placeholders filled (Matthias Krauser 2025, suorcd 2026) and the flake metadata corrected to MIT to match the LICENSE file and README
- The Nix-built binary now stamps its version via `-X main.version` (flake `ldflags`); it previously reported `dev`
- README clone snippet and indented code fences fixed; stale AGENTS.md notes corrected and the commit-trailer convention documented

## [v2.10.0] - 2026-09-03

### Added
- Bumped `slackhq/nebula` dependency to v1.11.1 (cert API unchanged, no source changes)

### Changed
- Bumped Go toolchain to Go 1.27 (`go.mod`), with `.goversion` and release workflow on 1.27.1
- Updated flake inputs and Nixpkgs pin
- Upgraded OpenBao deps to `api/v2 v2.6.0` and `sdk/v2 v2.6.2`, plus transitive module updates from `go get -u`

## [v2.0.5] - 2026-07-23

### Fixed
- Updated GoReleaser config to the v2 schema and bumped CI Go to 1.26

## [v2.0.4] - 2026-07-23

### Added
- Bumped `slackhq/nebula` dependency to v1.11.0 (cert API unchanged, no source changes)

### Changed
- Bumped Go version to 1.26.0, updated flake inputs

## [v2.0.3] - 2026-06-30

### Fixed
- **Critical**: replaced `ed25519.GenerateKey` with native `curve25519.X25519` key
  generation for node certificates (`/issue` endpoint). Previously the cert embedded
  an Ed25519 public key but claimed `Curve_CURVE25519`, causing Nebula V2 to reject
  the keypair due to a mathematical mismatch during scalar multiplication.

## [v2.0.2] - 2026-06-26

### Added
- OpenAPI specification (`openapi.yaml`) for all plugin endpoints

### Changed
- Renamed `/sign/{name}` endpoint to `/issue/{name}` to reflect that the plugin
  generates keypairs internally
- Updated README with comprehensive usage documentation

### Fixed
- Fixed X25519 private key PEM encoding to use the raw 32-byte seed

## [v2.0.1] - 2026-06-26

### Changed
- Updated CI workflow (GitHub Actions) to latest versions
- Bumped Go version to 1.25

## [v2.0.0] - 2026-06-25

### Added
- **Nebula V2 certificate support**: upgraded `slackhq/nebula` to v1.10.3+
- Nix flake for deterministic builds and development environment
- Build and dev server targets via Makefile

### Changed
- All certificate operations now issue and parse Nebula V2 certificates exclusively
- Imported CAs must be Nebula V2 format

## [1.0.1] - 2025-08-05

### Added
- Changelog

### Changed
- Fixed typo in goreleaser config

## [1.0.0] - 2025-08-05

### Added
- Basic CA certificate management
  - Generate new CA certificates
  - Import existing CA certificates
  - Read CA certificate information
- Node certificate management
  - Issue node certificates
  - List issued certificates
  - View certificate details
  - Revoke certificates
- CA certificate rotation functionality
  - Support for rotating CA with backup preservation
  - Ability to read both current and old CA certificates
  - Automatic backup of old CA during rotation
- Automated certificate cleanup (tidy) functionality
  - Configurable cleanup of expired certificates
  - Configurable cleanup of revoked certificates
  - Safety buffer period configuration
  - Manual and scheduled cleanup operations

### Changed
- Improved CA certificate management
  - Added validation for CA rotation operations
  - Enhanced error messages for CA operations
- Updated README with comprehensive documentation
- Restructured project layout for better maintainability

### Security
- Added validation to prevent unintended CA overwrites
- Added safety checks for CA rotation operations

[Unreleased]: https://github.com/suorcd/openbao-plugin-secrets-nebula/compare/v2.10.1...HEAD
[v2.10.1]: https://github.com/suorcd/openbao-plugin-secrets-nebula/compare/v2.10.0...v2.10.1
[v2.10.0]: https://github.com/suorcd/openbao-plugin-secrets-nebula/compare/v2.0.5...v2.10.0
[v2.0.5]: https://github.com/suorcd/openbao-plugin-secrets-nebula/compare/v2.0.4...v2.0.5
[v2.0.4]: https://github.com/suorcd/openbao-plugin-secrets-nebula/compare/v2.0.3...v2.0.4
[v2.0.3]: https://github.com/suorcd/openbao-plugin-secrets-nebula/compare/v2.0.2...v2.0.3
[v2.0.2]: https://github.com/suorcd/openbao-plugin-secrets-nebula/compare/v2.0.1...v2.0.2
[v2.0.1]: https://github.com/suorcd/openbao-plugin-secrets-nebula/compare/v2.0.0...v2.0.1
[v2.0.0]: https://github.com/suorcd/openbao-plugin-secrets-nebula/compare/v1.0.1...v2.0.0
[1.0.1]: https://github.com/suorcd/openbao-plugin-secrets-nebula/releases/tag/v1.0.1
[1.0.0]: https://github.com/suorcd/openbao-plugin-secrets-nebula/releases/tag/v1.0.0
