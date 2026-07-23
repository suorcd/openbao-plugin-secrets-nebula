# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v2.0.4] - 2026-07-23

### Added
- Bumped `slackhq/nebula` dependency to v1.11.0 (cert API unchanged, no source changes)

### Changed
- Bumped Go version to 1.26.0, updated flake inputs

## [Unreleased]

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

[Unreleased]: https://github.com/mkrauser/openbao-plugin-secrets-nebula/compare/v2.0.4...HEAD
[v2.0.4]: https://github.com/mkrauser/openbao-plugin-secrets-nebula/compare/v2.0.3...v2.0.4
[v2.0.3]: https://github.com/mkrauser/openbao-plugin-secrets-nebula/compare/v2.0.2...v2.0.3
[v2.0.2]: https://github.com/mkrauser/openbao-plugin-secrets-nebula/compare/v2.0.1...v2.0.2
[v2.0.1]: https://github.com/mkrauser/openbao-plugin-secrets-nebula/compare/v2.0.0...v2.0.1
[v2.0.0]: https://github.com/mkrauser/openbao-plugin-secrets-nebula/compare/v1.0.1...v2.0.0
[1.0.1]: https://github.com/mkrauser/openbao-plugin-secrets-nebula/releases/tag/v1.0.1
[1.0.0]: https://github.com/mkrauser/openbao-plugin-secrets-nebula/releases/tag/v1.0.0