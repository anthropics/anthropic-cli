# messages:batches

Manage batch processing of messages.

## Synopsis

```
ant messages:batches <command> [flags]
```

## Description

Batch processing allows you to submit multiple message requests in a single batch for efficient, high-throughput processing.

## Commands

| Command | Description |
|---------|-------------|
| `create` | Create a new batch |
| `retrieve` | Get batch details |
| `list` | List all batches |
| `delete` | Delete a batch |
| `cancel` | Cancel a processing batch |
| `results` | Get batch results |

## Workflow

1. Prepare a JSONL file with requests
2. Create a batch with `messages:batches create`
3. Monitor progress with `messages:batches retrieve`
4. Download results with `messages:batches results`

## Examples

### Create a batch

```bash
ant messages:batches create \
  --requests requests.jsonl \
  --endpoint /v1/messages
```

### Check batch status

```bash
ant messages:batches retrieve <batch-id>
```

### List all batches

```bash
ant messages:batches list
```

### Get results

```bash
ant messages:batches results <batch-id>
```

## See Also

- [messages:batches create](./messages-batches-create.md)
- [messages:batches retrieve](./messages-batches-retrieve.md)
- [messages:batches list](./messages-batches-list.md)
- [messages:batches results](./messages-batches-results.md)
