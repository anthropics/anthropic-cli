# Global Flags

Global flags are available for all Anthropic CLI commands. They control authentication, output formatting, and debugging.

## Reference

| Flag | Shorthand | Environment Variable | Default | Description |
|------|-----------|---------------------|---------|-------------|
| `--api-key` | - | `ANTHROPIC_API_KEY` | - | Your Anthropic API key |
| `--auth-token` | - | `ANTHROPIC_AUTH_TOKEN` | - | Anthropic auth token |
| `--base-url` | - | - | - | Override the API base URL |
| `--debug` | - | - | `false` | Enable debug logging |
| `--format` | - | - | `auto` | Output format |
| `--format-error` | - | - | `auto` | Error output format |
| `--transform` | - | - | - | GJSON transformation for output |
| `--transform-error` | - | - | - | GJSON transformation for errors |

## Details

### `--api-key`

Your Anthropic API key for authentication.

```bash
ant messages create \
  --api-key "sk-ant-api03-..." \
  --max-tokens 100 \
  --message "Hello" \
  --model claude-sonnet-4-5-20250929
```

**Environment variable**: `ANTHROPIC_API_KEY` (recommended)

### `--auth-token`

Anthropic auth token for advanced features.

```bash
ant --auth-token "..." models list
```

**Environment variable**: `ANTHROPIC_AUTH_TOKEN`

### `--base-url`

Override the default API base URL.

```bash
ant --base-url "https://custom-api.example.com" messages create ...
```

Useful for:
- Testing with mock servers
- Using a proxy
- Enterprise deployments with custom endpoints

### `--debug`

Enable debug logging to see full request and response details.

```bash
ant --debug messages create \
  --max-tokens 100 \
  --message "Hello" \
  --model claude-sonnet-4-5-20250929
```

Debug output includes:
- HTTP request method and URL
- Request headers (with auth redacted)
- Request body
- Response status and headers
- Response body

⚠️ **Warning**: Debug output may contain sensitive information. Use with caution in shared environments.

### `--format`

Control the output format for successful responses.

Available formats:

| Format | Description | Use Case |
|--------|-------------|----------|
| `auto` | Automatically select best format | General use (default) |
| `json` | Raw JSON | Scripting, piping |
| `yaml` | YAML format | Human-readable structured data |
| `pretty` | Pretty-printed with colors | Interactive use |
| `raw` | Unmodified response | Raw API output |
| `gjson` | GJSON transformed | Data extraction |

Examples:

```bash
# JSON for scripting
ant models list --format json | jq '.data[0].id'

# YAML for readability
ant models list --format yaml

# Pretty for interactive use
ant messages create ... --format pretty
```

### `--format-error`

Control the output format for error responses. Uses the same options as `--format`.

```bash
ant messages create ... --format-error json
```

### `--transform`

Apply a [GJSON](https://github.com/tidwall/gjson) transformation to the output.

```bash
# Extract just the model ID from the first result
ant models list --format json --transform "data.0.id"

# Extract all model IDs
ant models list --format json --transform "data.#.id"
```

GJSON path syntax:
- `.` - Child separator
- `#` - Array iterator
- `*` - Wildcard
- `?()` - Filters

### `--transform-error`

Apply a GJSON transformation to error responses.

```bash
ant messages create ... --transform-error "error.message"
```

## Combining Flags

You can combine multiple global flags:

```bash
ant --debug \
    --format yaml \
    --api-key "$ANTHROPIC_API_KEY" \
    messages create \
    --max-tokens 1024 \
    --message "Hello" \
    --model claude-sonnet-4-5-20250929
```

**Order matters**: Global flags must come before the command name, while command-specific flags come after.

## Configuration File

You can set default values for global flags in a config file:

```yaml
# ~/.config/ant/config.yaml
format: yaml
debug: false
base-url: https://api.anthropic.com
```

Priority (highest to lowest):
1. Command line flags
2. Environment variables
3. Config file values
4. Built-in defaults

## Environment Variables

Set common flags via environment variables:

```bash
export ANTHROPIC_API_KEY="your-key"
export ANTHROPIC_AUTH_TOKEN="your-token"
```

All environment variables take precedence over config file but are overridden by command line flags.

## See Also

- [CLI Overview](./README.md)
- [Installation](./installation.md)
- [Authentication](./authentication.md)
- [Output Formats](./output-formats.md)
