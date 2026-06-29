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
       -type=secret \
       -sha256="${SHA256}" \
       -command="openbao-plugin-secrets-nebula" \
       secret/openbao-plugin-secrets-nebula
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
   bao write nebula/sign/example.com \
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
       interval_duration="24h" \
       tidy_expired_certs=true \
       tidy_revoked_certs=true \
       safety_buffer="168h"  # 1 week safety buffer
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
       safety_buffer="48h"
   ```

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
    make test
    ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

## License

This project is licensed under the [MIT License](LICENSE).
