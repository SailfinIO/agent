# Sailfin Agent

## Installation

### Via Installer

```bash
curl -o- https://raw.githubusercontent.com/SailfinIO/agent/main/scripts/install.sh | bash
```

# Download the latest Linux binary (amd64)

curl -LO https://github.com/SailfinIO/agent/releases/latest/download/sailfin-linux-arm64.tar.gz

# Extract the tarball

tar -xzvf sailfin-linux-arm64.tar.gz

# Move the binary to a directory in your PATH (e.g., /usr/local/bin)

sudo mv sailfin-linux-arm64 /usr/local/bin/sailfin

# Verify installation

sailfin --help
