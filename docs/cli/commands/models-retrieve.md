# models retrieve

Get details for a specific Claude model.

## Synopsis

```
ant models retrieve <model-id> [flags]
```

## Description

Retrieves detailed information about a specific Claude model including its capabilities and context window.

## Arguments

| Argument | Required | Description |
|----------|----------|-------------|
| `model-id` | Yes | The model identifier |

## Examples

### Get model details

```bash
ant models retrieve claude-sonnet-4-5-20250929
```

### Extract specific fields

```bash
ant models retrieve claude-sonnet-4-5-20250929 \
  --format json \
  --transform "{id: id, name: display_name, created: created_at}"
```

## Response Format

```json
{
  "id": "claude-sonnet-4-5-20250929",
  "display_name": "Claude Sonnet 4.5",
  "created_at": "2025-09-29T00:00:00Z",
  "type": "model"
}
```

## API Reference

See the [Anthropic API documentation](https://docs.anthropic.com/en/api/models-get) for more details.
