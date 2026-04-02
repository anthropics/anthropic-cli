#!/usr/bin/env bash
# generate-docs.sh - Generate CLI documentation from source code
# Usage: ./scripts/generate-docs.sh

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
DOCS_DIR="${PROJECT_ROOT}/docs"
CLI_COMMANDS_DIR="${DOCS_DIR}/cli/commands"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Ensure commands directory exists
mkdir -p "${CLI_COMMANDS_DIR}"

# Generate completions command docs
cat > "${CLI_COMMANDS_DIR}/completions.md" << 'EOF'
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
EOF

log_info "Generated: completions.md"

# Generate messages command docs
cat > "${CLI_COMMANDS_DIR}/messages.md" << 'EOF'
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
EOF

log_info "Generated: messages.md"

# Generate messages create command docs
cat > "${CLI_COMMANDS_DIR}/messages-create.md" << 'EOF'
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
EOF

log_info "Generated: messages-create.md"

# Generate messages count-tokens command docs
cat > "${CLI_COMMANDS_DIR}/messages-count-tokens.md" << 'EOF'
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
EOF

log_info "Generated: messages-count-tokens.md"

# Generate messages:batches command docs
cat > "${CLI_COMMANDS_DIR}/messages-batches.md" << 'EOF'
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
EOF

log_info "Generated: messages-batches.md"

# Generate models command docs
cat > "${CLI_COMMANDS_DIR}/models.md" << 'EOF'
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
EOF

log_info "Generated: models.md"

# Generate models list command docs
cat > "${CLI_COMMANDS_DIR}/models-list.md" << 'EOF'
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
EOF

log_info "Generated: models-list.md"

# Generate models retrieve command docs
cat > "${CLI_COMMANDS_DIR}/models-retrieve.md" << 'EOF'
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
EOF

log_info "Generated: models-retrieve.md"

# Generate placeholder files for remaining commands
for cmd in messages-batches-create messages-batches-retrieve messages-batches-list messages-batches-delete messages-batches-cancel messages-batches-results beta-messages beta-messages-batches beta-models beta-files beta-skills beta-skills-versions; do
    if [[ ! -f "${CLI_COMMANDS_DIR}/${cmd}.md" ]]; then
        title=$(echo "$cmd" | tr '-' ' ' | sed 's/.*/\u&/')
        cat > "${CLI_COMMANDS_DIR}/${cmd}.md" << EOF
# ${title}

This command is part of the Anthropic CLI.

## Synopsis

\`\`\`
ant ${cmd//-/:} [flags]
\`\`\`

## Description

Documentation for this command is being generated. Please refer to the CLI help for detailed information:

\`\`\`bash
ant ${cmd//-/:} --help
\`\`\`

## See Also

- [CLI Overview](../README.md)
- [Global Flags](../global-flags.md)
EOF
        log_info "Generated: ${cmd}.md (placeholder)"
    fi
done

# Generate SUMMARY.md update
log_info "Documentation generation complete!"
log_info "Files created in: ${CLI_COMMANDS_DIR}"

# Count generated files
generated_count=$(find "${CLI_COMMANDS_DIR}" -name "*.md" | wc -l)
log_info "Total documentation files: ${generated_count}"
EOF

# Make script executable
chmod +x "${SCRIPT_DIR}/generate-docs.sh"

log_info "Documentation generation script created at: ${SCRIPT_DIR}/generate-docs.sh"
log_info "Run with: ./scripts/generate-docs.sh"
