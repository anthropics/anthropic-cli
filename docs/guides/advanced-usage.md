# Advanced Usage

This guide covers advanced features and techniques for power users of the Anthropic CLI.

## Piping and Redirection

The CLI works seamlessly with Unix pipes:

### Pipe Input from Files

```bash
# Read message from file
ant messages create \
  --max-tokens 1024 \
  --message "$(cat prompt.txt)" \
  --model claude-sonnet-4-5-20250929

# Process multiple prompts
for file in prompts/*.txt; do
  echo "Processing: $file"
  ant messages create \
    --max-tokens 1024 \
    --message "$(cat $file)" \
    --model claude-sonnet-4-5-20250929 \
    --format json \
    --transform "content.0.text" \
    > "outputs/$(basename $file .txt).txt"
done
```

### Pipe Output to Other Tools

```bash
# Extract and process with jq
ant models list --format json | \
  jq '.data[] | select(.id | contains("sonnet")) | .id'

# Count responses
ant messages:batches list --format json | \
  jq '.data | length'

# Save full responses for analysis
ant messages create ... --format json | \
  tee response.json | \
  jq '.usage'
```

## Shell Scripting

### Batch Operations

```bash
#!/bin/bash

MODEL="claude-sonnet-4-5-20250929"
MAX_TOKENS=1024

process_file() {
  local file=$1
  local content=$(cat "$file")
  
  ant messages create \
    --max-tokens $MAX_TOKENS \
    --message "{\"role\": \"user\", \"content\": [{\"type\": \"text\", \"text\": \"$content\"}]}" \
    --model $MODEL \
    --format json \
    --transform "content.0.text"
}

# Process all .md files
for file in docs/*.md; do
  echo "Processing $file..."
  process_file "$file" > "summaries/$(basename $file)"
done
```

### Error Handling

```bash
#!/bin/bash

make_request() {
  local response
  local exit_code
  
  response=$(ant messages create \
    --max-tokens 100 \
    --message "$1" \
    --model claude-sonnet-4-5-20250929 \
    --format json 2>&1)
  exit_code=$?
  
  if [ $exit_code -ne 0 ]; then
    echo "Error: $response" >&2
    return 1
  fi
  
  echo "$response"
}

# Use with error handling
if ! make_request "Hello"; then
  echo "Request failed"
  exit 1
fi
```

## GJSON Transformations

### Extract Nested Data

```bash
# Get usage statistics
ant messages create ... --format json --transform "usage"

# Get multiple fields
ant messages create ... --format json --transform "{input_tokens: usage.input_tokens, output_tokens: usage.output_tokens}"

# Filter arrays
ant models list --format json --transform "data.#(display_name%\"Claude*\")#.id"
```

### Modify Output Structure

```bash
# Flatten response
ant messages create ... --format json --transform "{text: content.0.text, model: model, tokens: usage.total_tokens}"

# Extract array elements
ant messages:batches results <id> --format json --transform "body.output.content.0.text"
```

## Debug Mode

### Full Request/Response Logging

```bash
# Enable debug mode
ant --debug messages create \
  --max-tokens 100 \
  --message "Hello" \
  --model claude-sonnet-4-5-20250929 2> debug.log

# Debug output includes:
# - HTTP request details
# - Headers (with auth redacted)
# - Request/response bodies
```

## Configuration Files

Create a default configuration:

```yaml
# ~/.config/ant/config.yaml
format: yaml
debug: false
base-url: https://api.anthropic.com
```

Platform-specific paths:

| Platform | Config Path |
|----------|-------------|
| macOS | `~/.config/ant/config.yaml` |
| Linux | `~/.config/ant/config.yaml` or `~/.ant/config.yaml` |
| Windows | `%APPDATA%\ant\config.yaml` |

## Working with Files

### Upload Files (Beta)

```bash
# Upload a file
ant beta:files upload \
  --file document.pdf \
  --purpose "batch"

# List uploaded files
ant beta:files list

# Download file
ant beta:files download <file-id> --output downloaded.pdf
```

### Use File IDs in Messages

```bash
ant messages create \
  --max-tokens 1024 \
  --message '{"role": "user", "content": [{"type": "document", "source": {"type": "file", "file_id": "file-xxx"}}]}' \
  --model claude-sonnet-4-5-20250929
```

## Batch Processing

### Efficient Bulk Operations

```bash
# Prepare batch requests
jq -c '{custom_id: .id, params: {model: "claude-sonnet-4-5-20250929", max_tokens: 1024, messages: [{role: "user", content: .content}]}}' inputs.json > requests.jsonl

# Create batch
ant messages:batches create \
  --requests requests.jsonl \
  --endpoint /v1/messages

# Monitor progress
watch -n 30 'ant messages:batches retrieve <batch-id> --format json --transform "{status: processing_status, completed: request_counts.completed, total: request_counts.total}"'

# Download results
ant messages:batches results <batch-id> --output results.jsonl
```

## Performance Tips

### Use Appropriate Models

- `claude-haiku-...` - Fast, cost-effective for simple tasks
- `claude-sonnet-...` - Balanced performance and cost
- `claude-opus-...` - Best quality for complex tasks

### Token Counting

Always count tokens before batch processing:

```bash
# Estimate costs
total_input=$(cat requests.jsonl | \
  jq -r '.params.messages[0].content' | \
  while read content; do
    ant messages count-tokens \
      --message "{\"role\": \"user\", \"content\": \"$content\"}" \
      --model claude-sonnet-4-5-20250929 \
      --format json \
      --transform "input_tokens"
  done | \
  awk '{sum += $1} END {print sum}')

echo "Total input tokens: $total_input"
```

### Connection Pooling

For high-throughput applications, consider:

```bash
# Use HTTP/2 for concurrent requests
export ANTHROPIC_BASE_URL="https://api.anthropic.com"

# Run parallel requests (be mindful of rate limits)
parallel -j 4 'ant messages create ...' ::: prompts/*.txt
```

## See Also

- [Output Formats](../cli/output-formats.md)
- [Error Handling](./error-handling.md)
- [Batch Processing](./batch-processing.md)
