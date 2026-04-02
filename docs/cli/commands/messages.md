# messages

Create and count tokens for messages with Claude models.

## Synopsis

```
ant messages <command> [flags]
```

## Description

The `messages` resource provides commands for creating messages and counting tokens with Claude models.

## Commands

| Command | Description |
|---------|-------------|
| `create` | Create a message |
| `count-tokens` | Count tokens in a message |

## Examples

### Create a message

```bash
ant messages create \
  --max-tokens 1024 \
  --message '{"role": "user", "content": [{"type": "text", "text": "Hello!"}]}' \
  --model claude-sonnet-4-5-20250929
```

### Count tokens

```bash
ant messages count-tokens \
  --message '{"role": "user", "content": [{"type": "text", "text": "Hello, world!"}]}' \
  --model claude-sonnet-4-5-20250929
```

## See Also

- [messages create](./messages-create.md)
- [messages count-tokens](./messages-count-tokens.md)
