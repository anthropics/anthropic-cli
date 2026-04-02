# Authentication

The Anthropic CLI requires authentication for all API requests. This guide covers how to configure your credentials securely.

## API Keys

Anthropic uses API keys to authenticate requests. You can create and manage your API keys in the [Anthropic Console](https://console.anthropic.com).

### Security Best Practices

⚠️ **Important**: Never commit your API key to version control or share it publicly.

1. Store your API key in environment variables or secure secret managers
2. Rotate your keys regularly
3. Use separate keys for development and production
4. Monitor your API usage for unexpected activity

## Setting Your API Key

### Method 1: Environment Variable (Recommended)

Set the `ANTHROPIC_API_KEY` environment variable:

```bash
# macOS/Linux
export ANTHROPIC_API_KEY="your-api-key-here"

# Windows (PowerShell)
$env:ANTHROPIC_API_KEY="your-api-key-here"

# Windows (CMD)
set ANTHROPIC_API_KEY=your-api-key-here
```

Add to your shell profile for persistence:

```bash
# ~/.zshrc or ~/.bashrc
echo 'export ANTHROPIC_API_KEY="your-api-key-here"' >> ~/.zshrc
source ~/.zshrc
```

### Method 2: Command Line Flag

Pass the API key directly to any command:

```bash
ant messages create \
  --api-key "your-api-key-here" \
  --max-tokens 1024 \
  --message "Hello!" \
  --model claude-sonnet-4-5-20250929
```

⚠️ Note: This method may expose your key in shell history.

### Method 3: Config File

Create a configuration file at `~/.config/ant/config.yaml`:

```yaml
api-key: your-api-key-here
```

Set appropriate permissions:

```bash
chmod 600 ~/.config/ant/config.yaml
```

## Authentication Priority

The CLI checks for credentials in this order (first match wins):

1. Command line flag (`--api-key`)
2. Environment variable (`ANTHROPIC_API_KEY`)
3. Config file (`~/.config/ant/config.yaml`)

## Verifying Authentication

Test your authentication:

```bash
ant models list
```

If authentication is successful, you'll see a list of available models. If not, you'll see an error:

```
Error: authentication failed: invalid API key
```

## Auth Tokens (Optional)

For certain advanced features, you may need an auth token:

```bash
export ANTHROPIC_AUTH_TOKEN="your-auth-token"
```

Or use the flag:

```bash
ant --auth-token "your-auth-token" <command>
```

## Troubleshooting

### "No API key provided" Error

```
Error: required flag "api-key" not set
```

**Solution**: Set the `ANTHROPIC_API_KEY` environment variable or pass `--api-key` flag.

### "Invalid API key" Error

```
Error: authentication failed: invalid API key
```

**Solutions**:
- Verify your API key is correct (no extra spaces or characters)
- Check if your key has been revoked in the console
- Generate a new key if needed

### Permission Denied on Config File

```
Error: cannot read config file: permission denied
```

**Solution**: Fix file permissions:

```bash
chmod 600 ~/.config/ant/config.yaml
```

## Key Rotation

To rotate your API key:

1. Generate a new key in the [Anthropic Console](https://console.anthropic.com)
2. Update your environment variable or config file
3. Test the new key
4. Revoke the old key

## Next Steps

- [Quick Start Guide](../guides/quickstart.md) — Make your first API call
- [CLI Commands Reference](./commands/) — Explore all available commands
