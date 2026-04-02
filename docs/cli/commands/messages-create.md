# messages create

Create a message using Claude models.

## Synopsis

```
ant messages create [flags]
```

## Description

Creates a message request to the Claude Messages API. This is the primary interface for interacting with Claude.

## Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--max-tokens` | Yes | - | Maximum number of tokens to generate |
| `--message` | Yes | - | Input message (JSON format) |
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

### Basic message

```bash
ant messages create \
  --max-tokens 1000 \
  --message '{"role": "user", "content": [{"type": "text", "text": "What is machine learning?"}]}' \
  --model claude-sonnet-4-5-20250929
```

### Multi-turn conversation

```bash
ant messages create \
  --max-tokens 2000 \
  --message '{"role": "user", "content": [{"type": "text", "text": "Previous message content"}]}' \
  --model claude-sonnet-4-5-20250929
```

### With streaming

```bash
ant messages create \
  --max-tokens 500 \
  --message '{"role": "user", "content": [{"type": "text", "text": "Tell me a joke"}]}' \
  --model claude-sonnet-4-5-20250929 \
  --stream
```

## Response Format

Successful responses include:

```json
{
  "id": "msg_01X...",
  "type": "message",
  "role": "assistant",
  "model": "claude-sonnet-4-5-20250929",
  "content": [
    {
      "type": "text",
      "text": "Response content here..."
    }
  ],
  "stop_reason": "end_turn",
  "stop_sequence": null,
  "usage": {
    "input_tokens": 10,
    "output_tokens": 50
  }
}
```

## API Reference

See the [Anthropic Messages API documentation](https://docs.anthropic.com/en/api/messages) for more details.
