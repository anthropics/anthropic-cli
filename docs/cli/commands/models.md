# models

Retrieve and list available Claude models.

## Synopsis

```
ant models <command> [flags]
```

## Description

List and retrieve information about available Claude models.

## Commands

| Command | Description |
|---------|-------------|
| `list` | List all available models |
| `retrieve` | Get details for a specific model |

## Examples

### List all models

```bash
ant models list
```

### Get specific model details

```bash
ant models retrieve claude-sonnet-4-5-20250929
```

### Filter with jq

```bash
ant models list --format json | jq '.data[] | select(.id | contains("sonnet"))'
```

## See Also

- [models list](./models-list.md)
- [models retrieve](./models-retrieve.md)
