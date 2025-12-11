# Claude Developer Platform CLI

The official CLI for the [Anthropic REST API](https://docs.anthropic.com/claude/reference/).

## Installation

### Installing with Go

```sh
go install 'github.com/stainless-sdks/anthropic-cli/cmd/cdp@latest'
```

### Running Locally

```sh
go run cmd/cdp/main.go
```

## Usage

The CLI follows a resource-based command structure:

```sh
cdp [resource] [command] [flags]
```

```sh
cdp messages create \
  --max-tokens 1024 \
  --message '{content: [{text: x, type: text}], role: user}' \
  --model claude-sonnet-4-5-20250929
```

For details about specific commands, use the `--help` flag.

## Global Flags

- `--debug` - Enable debug logging (includes HTTP request/response details)
- `--version`, `-v` - Show the CLI version
