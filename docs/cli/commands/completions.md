# completions create

Create a completion using Claude models.

## Synopsis

```
ant completions create [flags]
```

## Description

Creates a completion request to the Claude API. This is a simplified interface for generating completions.

## Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--max-tokens` | Yes | - | Maximum number of tokens to generate |
| `--message` | Yes | - | Message content (JSON format) |
| `--model` | Yes | - | Model ID to use |
| `--metadata` | No | - | Additional metadata (JSON format) |
| `--stop-sequences` | No | - | Sequences that stop generation |
| `--stream` | No | false | Stream the response |
| `--system` | No | - | System prompt |
| `--temperature` | No | 1.0 | Sampling temperature |
| `--tool-choice` | No | - | Tool choice configuration |
| `--tools` | No | - | Tools available to the model |
| `--top-k` | No | - | Top-k sampling |
| `--top-p` | No | - | Nucleus sampling |

## Examples

### Basic completion

```bash
ant completions create \
  --max-tokens 100 \
  --message '{"role": "user", "content": [{"type": "text", "text": "Hello!"}]}' \
  --model claude-sonnet-4-5-20250929
```

### With system prompt

```bash
ant completions create \
  --max-tokens 200 \
  --system "You are a helpful assistant." \
  --message '{"role": "user", "content": [{"type": "text", "text": "Explain quantum computing"}]}' \
  --model claude-sonnet-4-5-20250929
```

### Streaming response

```bash
ant completions create \
  --max-tokens 500 \
  --message '{"role": "user", "content": [{"type": "text", "text": "Write a story"}]}' \
  --model claude-sonnet-4-5-20250929 \
  --stream
```

## API Reference

See the [Anthropic API documentation](https://docs.anthropic.com/en/api/complete) for more details.
