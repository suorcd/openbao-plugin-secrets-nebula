# AGENTS.md — OpenBao Nebula Secrets Plugin

## Quick commands

```shell
nix develop                           # dev shell (Go 1.25, goreleaser, golangci-lint)
nix build                             # binary → result/bin/

go build -buildvcs=false ./...        # build (VCS stamping fails in detached worktrees)
go test -buildvcs=false ./...         # tests (one file: path_tidy_test.go)
go vet -buildvcs=false ./             # vet

make build                            # → bao/plugins/bao-plugin-secrets-nebula
make fmt                              # go fmt on all packages
```

Note: there is NO `make test` target. README is wrong.

## Go version

**1.25** (go.mod, CI). `.goversion` says 1.23.4 — ignore it, it's stale.

## Module path

`go.mod` declares `github.com/mkrauser/openbao-plugin-secrets-nebula`. All source imports use this path. Do not change without updating go.mod.

## Architecture

- Single Go package `nebula` at repo root.
- Entrypoint: `cmd/openbao-plugin-secrets-nebula/main.go`.
- Uses OpenBao v2 SDK `framework.Backend` pattern.
- Paths: generate/ca, config/ca, issue/:name, cert/:fingerprint, certs/, revoke, certs/revoked/, tidy, tidy-cancel, tidy-status, config/auto-tidy.

## Certificates

- **Nebula V2** only (`cert.Version2`).
- CA keys: Ed25519 (`ed25519.PrivateKey`), stored as JSON.
- Node keys: **native Curve25519 (X25519)**. Generated via `rand.Read` + `curve25519.X25519(priv, Basepoint)`. Do NOT use `ed25519.GenerateKey` for node keypairs — that was a crypto bug.
- `CertStorageEntry` (`util.go:12`) wraps certs as `{"pem":"..."}` JSON for storage (V2 certs are interfaces, can't directly JSON-serialize).
- Fingerprints: 64 hex chars → 79 with colons via `formatFingerprint`.

## Testing

- `path_tidy_test.go` only. Tests use `logical.InmemStorage` — no external services.
- Helper `createBackendWithStorage(t)` returns `(*backend, logical.Storage)`.
- Test CA bundle constant embedded in the test file.
- Tidy tests are async (background goroutine), poll tidy-status with sleeps.


## Known quirks

- `pathConfigCADelete` is a no-op (`path_ca.go:435`).
- `pathListRevokedCertsHandler` lists `certs/` instead of `revoked/` — likely a bug.
- `safety_buffer` field type is `framework.TypeDurationSecond` (integer seconds), not Go duration strings. README examples using `"168h"` are wrong for this parameter.
- No PR CI — only `release.yml` on `v*` tags via goreleaser. Test/lint must be run locally.
- `make all` starts an OpenBao dev server (terminal-blocking).
- `/bao` and `result` dirs are gitignored.

## Plugin registration (OpenBao CLI)

When registering/reloading the plugin on a live server:

```shell
# MUST be logged in first (root or policy with sys/plugins/catalog/secret/*)
bao login

SHA256=$(sha256sum /vault/data/plugins/openbao-plugin-secrets-nebula | cut -d' ' -f1)

# Flags BEFORE positional args. No -type flag.
bao plugin register \
    -sha256="${SHA256}" \
    -command="openbao-plugin-secrets-nebula" \
    secret openbao-plugin-secrets-nebula

bao plugin reload -plugin openbao-plugin-secrets-nebula
```

- Plugin checksum mismatch on reload can leave mounts in a permanent "cannot write during storage setup" state. If reload fails, check checksums match before retrying.
- Remount is required if the mount gets stuck in Setup deadlock.
