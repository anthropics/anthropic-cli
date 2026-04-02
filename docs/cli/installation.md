# Installation

This guide covers all methods to install the Anthropic CLI (`ant`) on your system.

## Requirements

- **macOS, Linux, or Windows** (with WSL)
- **API Key** - Sign up at [console.anthropic.com](https://console.anthropic.com) to get your key

## Method 1: Homebrew (macOS/Linux)

The easiest way to install on macOS and Linux:

```bash
brew install anthropics/tap/ant
```

Upgrade to the latest version:

```bash
brew upgrade ant
```

## Method 2: Go Install

If you have Go 1.22 or later installed:

```bash
go install github.com/anthropics/anthropic-cli/cmd/ant@latest
```

The binary will be installed in your `$GOPATH/bin` directory (usually `~/go/bin`).

Add to your PATH if needed:

```bash
# Add to ~/.zshrc or ~/.bashrc
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Method 3: Download Pre-built Binary

Download the latest release for your platform:

```bash
# macOS (Apple Silicon)
curl -L "https://github.com/anthropics/anthropic-cli/releases/latest/download/ant-darwin-arm64" -o ant
chmod +x ant
sudo mv ant /usr/local/bin/

# macOS (Intel)
curl -L "https://github.com/anthropics/anthropic-cli/releases/latest/download/ant-darwin-amd64" -o ant
chmod +x ant
sudo mv ant /usr/local/bin/

# Linux (x86_64)
curl -L "https://github.com/anthropics/anthropic-cli/releases/latest/download/ant-linux-amd64" -o ant
chmod +x ant
sudo mv ant /usr/local/bin/
```

## Method 4: Build from Source

```bash
git clone https://github.com/anthropics/anthropic-cli.git
cd anthropic-cli
go build -o ant ./cmd/ant
sudo mv ant /usr/local/bin/
```

## Verify Installation

Check that the CLI is properly installed:

```bash
ant --version
```

You should see output like:

```
ant version 1.0.0
```

## Next Steps

1. [Set up authentication](./authentication.md)
2. Try the [Quick Start Guide](../guides/quickstart.md)

## Troubleshooting

### Command not found

If you get `command not found` after installation:

1. Check if the binary is in your PATH:
   ```bash
   which ant
   ```

2. If using `go install`, ensure `$(go env GOPATH)/bin` is in your PATH:
   ```bash
   echo $PATH | grep "$(go env GOPATH)/bin"
   ```

3. Reload your shell configuration:
   ```bash
   source ~/.zshrc  # or ~/.bashrc
   ```

### Permission denied

If you get `permission denied` errors:

```bash
# Make binary executable
chmod +x /path/to/ant

# Or if in /usr/local/bin
sudo chmod +x /usr/local/bin/ant
```

### macOS "unidentified developer" warning

On macOS, you may see a security warning. To resolve:

1. Go to **System Preferences** → **Security & Privacy** → **General**
2. Click **Allow Anyway** next to the message about `ant`
3. Run the command again

Or use the terminal:

```bash
xattr -d com.apple.quarantine /usr/local/bin/ant
```

## Uninstallation

### Homebrew

```bash
brew uninstall ant
```

### Manual

```bash
sudo rm /usr/local/bin/ant
```
