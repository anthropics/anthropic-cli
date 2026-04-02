# Anthropic CLI Documentation

Welcome to the official documentation for the **Anthropic CLI** (`ant`) — a powerful command-line interface for the [Claude Developer Platform](https://platform.anthropic.com).

## What is the Anthropic CLI?

The Anthropic CLI provides a resource-based command structure for interacting with Anthropic's APIs directly from your terminal. It supports:

- **Message creation** with Claude models
- **Batch processing** for high-throughput workloads
- **Model management** and metadata retrieval
- **Beta features** including files and skills management
- **Advanced formatting** with JSON, YAML, and GJSON transformations

## Quick Start

### Installation

```bash
# macOS with Homebrew
brew install anthropics/tap/ant

# Or install with Go
go install github.com/anthropics/anthropic-cli/cmd/ant@latest
```

### Authentication

Set your API key:

```bash
export ANTHROPIC_API_KEY="your-api-key-here"
```

Or pass it directly:

```bash
ant messages create \
  --api-key $ANTHROPIC_API_KEY \
  --max-tokens 1024 \
  --message '{"role": "user", "content": [{"type": "text", "text": "Hello!"}]}' \
  --model claude-sonnet-4-5-20250929
```

### First Command

Create your first message:

```bash
ant messages create \
  --max-tokens 1024 \
  --message '{"role": "user", "content": [{"type": "text", "text": "Say hello!"}]}' \
  --model claude-sonnet-4-5-20250929
```

## Documentation Structure

- **[CLI Reference](cli/)** — Complete command documentation
- **[Guides](guides/)** — Tutorials and how-to articles
- **[API Reference](api/)** — API overview and concepts

## Features

### Resource-Based Commands

The CLI organizes commands by resource type:

| Resource | Description |
|----------|-------------|
| `messages` | Create and manage messages with Claude |
| `messages:batches` | Batch processing for multiple messages |
| `models` | List and retrieve model information |
| `beta:files` | File management (beta) |
| `beta:skills` | Skills management (beta) |

### Output Formats

Control how responses are displayed:

- `auto` — Automatically select best format (default)
- `json` — Raw JSON output
- `yaml` — YAML formatted output
- `pretty` — Human-readable formatted output
- `raw` — Unformatted output
- `gjson` — GJSON transformed output

### Debug Mode

Enable debug logging to see full request/response details:

```bash
ant --debug messages create --max-tokens 100 --message "Hello"
```

## Getting Help

- Use `ant --help` for global options
- Use `ant <command> --help` for command-specific help
- Visit [Anthropic Support](https://support.anthropic.com) for assistance

## Contributing

This documentation is open source. See [CONTRIBUTING.md](../CONTRIBUTING.md) for details on how to contribute.

---

**Version:** v1.0.0  
**Last Updated:** 2026-04-01
