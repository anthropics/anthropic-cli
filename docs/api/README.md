# API Overview

The Anthropic CLI provides a command-line interface to the [Anthropic API](https://docs.anthropic.com/en/api), enabling you to interact with Claude models programmatically.

## API Concepts

### Messages API

The Messages API is the primary interface for interacting with Claude. It supports:

- **Single-turn conversations**: Send a message, get a response
- **Multi-turn conversations**: Build conversation history
- **Streaming**: Real-time token-by-token responses
- **Tool use**: Enable Claude to use external tools
- **Vision**: Process images alongside text

### Models

Anthropic provides several model families:

| Model Family | Best For | Speed | Cost |
|-------------|----------|-------|------|
| Claude Opus | Complex reasoning, coding, analysis | Slower | Higher |
| Claude Sonnet | Balanced performance and cost | Medium | Medium |
| Claude Haiku | Quick tasks, high throughput | Fast | Lower |

### Authentication

All API requests require an API key passed via the `x-api-key` header:

```bash
curl https://api.anthropic.com/v1/messages \
  --header "x-api-key: $ANTHROPIC_API_KEY" \
  --header "anthropic-version: 2023-06-01" \
  --header "content-type: application/json" \
  --data '{
    "model": "claude-sonnet-4-5-20250929",
    "max_tokens": 1024,
    "messages": [
      {"role": "user", "content": "Hello, world!"}
    ]
  }'
```

## Rate Limits

| Tier | Requests/min | Tokens/min | Batch queue |
|------|---------------|------------|-------------|
| Free | 30 | 25,000 | 10 requests |
| Build | 50 | 50,000 | 100 requests |
| Scale | 1,000 | 400,000 | 10,000 requests |
| Enterprise | Custom | Custom | Custom |

## API Versioning

The API uses date-based versioning:

- Current version: `2023-06-01`
- Specify version in the `anthropic-version` header
- The CLI handles versioning automatically

## Batch Processing

For high-throughput workloads, use the batch processing API:

- Submit up to 10,000 requests per batch
- Results available within 24 hours
- 50% cost reduction compared to synchronous requests

## Error Handling

The API uses standard HTTP status codes:

| Status | Meaning |
|--------|---------|
| 200 | Success |
| 400 | Bad request |
| 401 | Unauthorized |
| 429 | Rate limited |
| 500 | Server error |
| 529 | Overloaded |

## See Also

- [Authentication](./authentication.md)
- [Rate Limits](./rate-limits.md)
- [Error Codes](./error-codes.md)
