# models list

List all available Claude models.

## Synopsis

```
ant models list [flags]
```

## Description

Returns a list of all available Claude models with their capabilities and metadata.

## Flags

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `--limit` | No | 20 | Maximum number of results |
| `--after-id` | No | - | Pagination cursor |

## Examples

### List all models

```bash
ant models list
```

### Limit results

```bash
ant models list --limit 10
```

### Get model IDs only

```bash
ant models list --format json --transform "data.#.id"
```

## Response Format

```json
{
  "data": [
    {
      "id": "claude-sonnet-4-5-20250929",
      "display_name": "Claude Sonnet 4.5",
      "created_at": "2025-09-29T00:00:00Z",
      "type": "model"
    }
  ],
  "has_more": false,
  "first_id": "...",
  "last_id": "..."
}
```

## API Reference

See the [Anthropic API documentation](https://docs.anthropic.com/en/api/models-list) for more details.
