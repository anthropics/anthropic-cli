# CLI Overview

The Anthropic CLI (`ant`) provides a resource-based command structure for interacting with the Claude Developer Platform API.

## Command Structure

```
ant [resource] <command> [flags...]
```

Resources represent API entities (like `messages`, `models`), and commands represent actions (like `create`, `list`, `retrieve`).

## Available Resources

### Core Resources

| Resource | Description |
|----------|-------------|
| `completions` | Create completions with Claude models |
| `messages` | Create and count tokens for messages |
| `messages:batches` | Manage batch processing of messages |
| `models` | Retrieve and list available models |

### Beta Resources

| Resource | Description |
|----------|-------------|
| `beta:messages` | Beta message features |
| `beta:messages:batches` | Beta batch processing |
| `beta:models` | Beta model features |
| `beta:files` | File management |
| `beta:skills` | Skills management |
| `beta:skills:versions` | Skill version management |

## Getting Help

Every command supports the `--help` flag:

```bash
# Global help
ant --help

# Resource help
ant messages --help

# Command help
ant messages create --help
```

## Output Formats

Control how responses are displayed using the `--format` flag:

| Format | Description |
|--------|-------------|
| `auto` | Automatically select best format (default) |
| `json` | Raw JSON output |
| `yaml` | YAML formatted output |
| `pretty` | Human-readable formatted output |
| `raw` | Unformatted output |
| `gjson` | GJSON transformed output |

Examples:

```bash
# Pretty formatted output (default)
ant models list

# JSON output for scripting
ant models list --format json

# YAML output
ant models list --format yaml

# Raw output
ant models list --format raw
```

## Global Flags

These flags are available for all commands:

| Flag | Environment Variable | Description |
|------|---------------------|-------------|
| `--debug` | - | Enable debug logging with full request/response details |
| `--base-url` | - | Override the base URL for API requests |
| `--format` | - | Output format (auto, json, yaml, pretty, raw, gjson) |
| `--format-error` | - | Error format |
| `--transform` | - | GJSON transformation for data output |
| `--transform-error` | - | GJSON transformation for errors |
| `--api-key` | `ANTHROPIC_API_KEY` | Your Anthropic API key |
| `--auth-token` | `ANTHROPIC_AUTH_TOKEN` | Your Anthropic auth token |

## Examples

### Create a Message

```bash
ant messages create \
  --max-tokens 1024 \
  --message '{"role": "user", "content": [{"type": "text", "text": "Hello!"}]}' \
  --model claude-sonnet-4-5-20250929
```

### List Models

```bash
ant models list
```

### Create a Batch

```bash
ant messages:batches create \
  --requests 'requests.jsonl' \
  --endpoint /v1/messages
```

## Error Handling

The CLI provides detailed error messages with context:

```bash
$ ant messages create --max-tokens 100
Error: required flag "model" not set
```

With the `--debug` flag, you get full request/response details for troubleshooting.

## Shell Completion

Enable tab completion for your shell:

```bash
# Bash
ant @completion bash > /etc/bash_completion.d/ant

# Zsh
ant @completion zsh > /usr/share/zsh/site-functions/_ant

# Fish
ant @completion fish > ~/.config/fish/completions/ant.fish
```

## Next Steps

- Learn about [Installation](./installation.md)
- Set up [Authentication](./authentication.md)
- Explore [Global Flags](./global-flags.md) in detail
- Read the [Quick Start Guide](../guides/quickstart.md)
