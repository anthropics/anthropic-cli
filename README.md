# Anthropic CLI

The official CLI for the [Anthropic REST API](https://docs.anthropic.com/claude/reference/).

## Installation

### Installing with Go

```sh
go install 'github.com/stainless-sdks/anthropic-cli/cmd/anthropic-cli@latest'
```

## Usage

The CLI follows a resource-based command structure:

```sh
anthropic-cli [resource] [command] [flags]
```

```sh
anthropic-cli messages create \
  --max-tokens 1024 \
  --messages.content.text x \
  --messages.content.type text \
  --messages.role user \
  --model claude-sonnet-4-20250514
```

For details about specific commands, use the `--help` flag.

## Global Flags

- `--debug` - Enable debug logging (includes HTTP request/response details)
- `--version`, `-v` - Show the CLI version
