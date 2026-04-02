# messages count-tokens

Count tokens in a message without generating a response.

## Synopsis

```
ant messages count-tokens [flags]
```

## Description

Counts the number of input tokens for a message request. Useful for estimating costs before making a request.

## Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--message` | Yes | - | Input message (JSON format) |
| `--model` | Yes | - | Model ID to use |
| `--system` | No | - | System prompt |
| `--tool-choice` | No | - | Tool choice configuration |
| `--tools` | No | - | Tools available to the model |

## Examples

### Count tokens for a message

```bash
ant messages count-tokens \
  --message '{"role": "user", "content": [{"type": "text", "text": "Hello, world!"}]}' \
  --model claude-sonnet-4-5-20250929
```

### With system prompt

```bash
ant messages count-tokens \
  --system "You are a helpful assistant." \
  --message '{"role": "user", "content": [{"type": "text", "text": "Explain quantum computing"}]}' \
  --model claude-sonnet-4-5-20250929
```

## Response Format

```json
{
  "input_tokens": 15
}
```

## Use Cases

- **Cost estimation**: Calculate tokens before batch processing
- **Rate limit planning**: Understand your token usage patterns
- **Optimization**: Identify opportunities to reduce token count

## API Reference

See the [Anthropic API documentation](https://docs.anthropic.com/en/api/messages-count-tokens) for more details.
