# OpenBao Nebula Secrets Plugin (V2 Certificate Support)

A secrets engine plugin for [OpenBao](https://github.com/openbao/openbao) (and HashiCorp Vault) that manages [Slack Nebula](https://github.com/slackhq/nebula) certificates.

> [!IMPORTANT]  
> **V2 Certificate Support:** This fork upgrades the underlying `slackhq/nebula` engine to `v1.10.3+` to fully support issuing, managing, and parsing modern **Nebula V2 Certificates**.

> [!NOTE]  
> **Vault Compatibility:** This plugin was developed and tested against OpenBao, but basic testing with HashiCorp Vault looks promising. If you test this with Vault, please open an issue if you run into any compatibility problems.

## Features

- **CA Management**:
  - Generate new Nebula V2 CA certificates
  - Import existing CA certificates seamlessly
  - Rotate CA certificates with backup preservation
  - View current and previous CA certificates
- **Node Certificate Management**:
  - Issue modern Nebula V2 node certificates dynamically
  - List all issued certificates
  - View individual certificate details and fingerprints
  - Revoke certificates
- **Automatic Maintenance**:
  - Configure automatic cleanup of expired certificates
  - Set safety buffer periods for certificate cleanup
  - Manual and scheduled cleanup operations

## Installation

1. Download the latest plugin release tarball from the [Releases page](https://github.com/suorcd/openbao-plugin-secrets-nebula/releases) (e.g., `openbao-plugin-secrets-nebula_Linux_x86_64.tar.gz`).
2. Extract the binary and register the plugin with OpenBao:

   ```shell
   # Download and extract the plugin
   wget [https://github.com/suorcd/openbao-plugin-secrets-nebula/releases/latest/download/openbao-plugin-secrets-nebula_Linux_x86_64.tar.gz](https://github.com/suorcd/openbao-plugin-secrets-nebula/releases/latest/download/openbao-plugin-secrets-nebula_Linux_x86_64.tar.gz)
   tar -xzf openbao-plugin-secrets-nebula_Linux_x86_64.tar.gz

   # Move the plugin to OpenBao's plugin directory
   mv openbao-plugin-secrets-nebula /etc/openbao/plugins/
   chmod +x /etc/openbao/plugins/openbao-plugin-secrets-nebula

   # Calculate the SHA256 sum of the plugin
   SHA256=$(sha256sum /etc/openbao/plugins/openbao-plugin-secrets-nebula | cut -d' ' -f1)

   # Register the plugin in the system catalog
   bao plugin register \
       -sha256="${SHA256}" \
       -command="openbao-plugin-secrets-nebula" \
       secret openbao-plugin-secrets-nebula
   ```

## Usage

### Enable the Plugin

    ```shell
    bao secrets enable -path=nebula -plugin-name=openbao-plugin-secrets-nebula plugin
    ```

### CA Certificate Management

1. Generate a new CA:

   ```shell
   # Generate a new CA with a 1-year validity period
   bao write nebula/generate/ca \
       name="my-nebula-ca" \
       duration="8760h" \
       ips="10.0.0.0/20" \
       groups="servers,clients"
   ```

2. Import an existing CA:

   ```shell
   # Import CA from a PEM bundle (private key + certificate)
   bao write nebula/config/ca pem_bundle=@bundle.pem
   ```

3. Read CA information:

   ```shell
   bao read nebula/config/ca
   ```

4. Rotate CA certificate:

   ```shell
   # Rotate with a new generated CA
   bao write nebula/generate/ca name="new-ca" rotate=true

   # Or rotate with an imported CA
   bao write nebula/config/ca pem_bundle=@new_bundle.pem rotate=true
   ```

### Node Certificate Management

1. Issue a node certificate:

   ```shell
    bao write nebula/issue/example.com \
        ip="10.0.0.1/32" \
        duration="720h" \
        groups="servers"
   ```

2. List all certificates:

   ```shell
   bao list nebula/certs
   ```

3. View certificate details:
   ```shell
   bao read nebula/cert/<fingerprint>
   ```

### Certificate Cleanup

1. Configure automatic cleanup:

   ```shell
    bao write nebula/config/auto-tidy \
        enabled=true \
        interval_duration=86400 \
        tidy_expired_certs=true \
        tidy_revoked_certs=true \
        safety_buffer=604800  # 1 week safety buffer in seconds
   ```

2. View cleanup configuration:

   ```shell
   bao read nebula/config/auto-tidy
   ```

3. Run manual cleanup:
   ```shell
    bao write nebula/tidy \
        tidy_expired_certs=true \
        tidy_revoked_certs=true \
        safety_buffer=172800  # 48 hours in seconds
   ```

## Upgrading the Plugin

> **WARNING**: Never overwrite the existing plugin catalog entry. Plugin checksum mismatches between the catalog and binary on disk will permanently lock mounts into a "cannot write to storage during setup" state that requires a full disable/re-enable.

OpenBao 2.5+ supports versioned plugin registration. Register each release as a distinct version, then tune mounts to the new version and reload globally. This avoids destructive overwrites and enables rollback.

### Step 1: Distribute the new binary

The binary must be on disk on **every** OpenBao node **before** updating the catalog. Never rely on init containers for upgrades — they only run at pod start.

**Kubernetes (Rancher / Helm chart):**

```shell
# Copy the binary to all pods while they are running
for pod in openbao-vault-0 openbao-vault-1 openbao-vault-2; do
  kubectl exec -n openbao $pod -- wget -qO- \
    https://github.com/suorcd/openbao-plugin-secrets-nebula/releases/download/v2.1.0/openbao-plugin-secrets-nebula_Linux_x86_64.tar.gz \
    | tar -xz -C /vault/data/plugins/
done
```

**Bare-metal / single-node:**

```shell
wget -qO- https://github.com/suorcd/openbao-plugin-secrets-nebula/releases/download/v2.1.0/openbao-plugin-secrets-nebula_Linux_x86_64.tar.gz \
  | tar -xz -C /etc/openbao/plugins/
```

Only after the binary is on **all** nodes should you proceed to step 2.

### Step 2: Register the new version

```shell
# Log in with a root token or a token with sys/plugins/catalog/secret/* permissions
bao login

SHA256=$(sha256sum /etc/openbao/plugins/openbao-plugin-secrets-nebula | cut -d' ' -f1)

# Flags MUST come before positional arguments. No -type flag.
bao plugin register \
    -sha256="${SHA256}" \
    -command="openbao-plugin-secrets-nebula" \
    -version=v2.1.0 \
    secret openbao-plugin-secrets-nebula
```

This adds `v2.1.0` to the catalog **without** removing the previous version.

### Step 3: Tune the mount to the new version

```shell
bao secrets tune -plugin-version=v2.1.0 nebula
```

The mount continues running the old version until reloaded.

### Step 4: Reload globally

```shell
bao plugin reload -plugin openbao-plugin-secrets-nebula -scope global
```

The `global` scope ensures all Raft replicas pick up the new binary. After reload, verify:

```shell
bao secrets list -detailed | grep nebula
```

The `Running Version` column should now match the `Version` column.

### Rollback

If the new version causes issues, tune back to the previous version and reload:

```shell
bao secrets tune -plugin-version=<previous-version> nebula
bao plugin reload -plugin openbao-plugin-secrets-nebula -scope global
```

### Troubleshooting

| Symptom | Cause | Fix |
|---|---|---|
| `cannot write to storage during setup` | Checksum mismatch — binary on disk doesn't match catalog SHA for the pinned version | Disable + re-enable the mount, then re-import CA |
| `checksums did not match` on reload | A Raft replica still has the old binary | Distribute the binary to the failing node and retry |
| Mount won't disable | Plugin process deadlocked | Restart the pod / OpenBao service |

## Development

### Prerequisites

- **Go 1.25** or higher (Required for Nebula V2 compatibility)
- OpenBao development environment

### Building with Nix (Recommended)

This repository includes a `flake.nix` for deterministic, reproducible builds and development environments.

    ```shell
    # Enter an ephemeral development shell containing Go 1.25, goreleaser, and linting tools
    nix develop

    # Or build the binary directly without polluting your host environment
    nix build
    # The compiled binary will be placed in ./result/bin/openbao-plugin-secrets-nebula
    ```

### Building with Standard Go Tools

    ```shell
    # Clone the repository
    git clone [https://github.com/suorcd/openbao-plugin-secrets-nebula](https://github.com/suorcd/openbao-plugin-secrets-nebula)
    cd openbao-plugin-secrets-nebula

    # Install dependencies and build
    go mod tidy
    go build ./cmd/openbao-plugin-secrets-nebula
    ```

### Testing

    ```shell
    # Run tests
    go test -buildvcs=false ./...
    ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the [MIT License](LICENSE).
