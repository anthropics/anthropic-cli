# Anthropic CLI

The official CLI for the [Anthropic REST API](https://docs.anthropic.com/claude/reference/).

## Installation

### Installing with Go

```sh
go install 'github.com/stainless-sdks/anthropic-cli'
```

## Usage

The CLI follows a resource-based command structure:

```sh
anthropic-cli [resource] [command] [flags]
```

```sh
anthropic-cli messages create \
  --max-tokens 1024 \
  --messages.content.id srvtoolu_SQfNkl1n_JR_ \
  --messages.content.name web_search \
  --messages.content.type server_tool_use \
  --messages.role user \
  --model claude-3-5-sonnet-latest
```

For details about specific commands, use the `--help` flag.

## Global Flags

- `--debug` - Enable debug logging (includes HTTP request/response details)
- `--version`, `-v` - Show the CLI version
