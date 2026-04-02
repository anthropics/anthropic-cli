# Output Formats

The Anthropic CLI supports multiple output formats to suit different use cases — from human-readable pretty printing to machine-parseable JSON.

## Available Formats

| Format | Description | Best For |
|--------|-------------|----------|
| `auto` | Automatically select best format | General use (default) |
| `json` | Raw JSON | Scripting, APIs, piping |
| `yaml` | YAML formatted | Configuration files, readability |
| `pretty` | Colorized, formatted | Interactive terminal use |
| `raw` | Unmodified response | Debugging, raw API access |
| `gjson` | GJSON transformed | Data extraction, filtering |

## Format Details

### auto (Default)

Automatically selects the most appropriate format based on context:

- Terminal with TTY: `pretty`
- Piped output: `json`
- Non-interactive: `json`

```bash
ant models list
# Uses pretty format in terminal
```

### json

Raw JSON output without formatting:

```bash
ant models list --format json
```

Output:
```json
{"data":[{"id":"claude-sonnet-4-5-20250929","display_name":"Claude Sonnet 4.5","created_at":"2025-09-29T00:00:00Z"}],"has_more":false}
```

Perfect for scripting:

```bash
# Extract model IDs
ant models list --format json | jq -r '.data[].id'

# Check if a specific model exists
ant models list --format json | jq -e '.data[] | select(.id == "claude-sonnet-4-5-20250929")'
```

### yaml

YAML formatted output for improved readability:

```bash
ant models list --format yaml
```

Output:
```yaml
data:
  - id: claude-sonnet-4-5-20250929
    display_name: Claude Sonnet 4.5
    created_at: "2025-09-29T00:00:00Z"
has_more: false
```

### pretty

Colorized, formatted output optimized for terminal viewing:

```bash
ant models list --format pretty
```

Features:
- Syntax highlighting
- Color coding by data type
- Tree-like structure visualization
- Pager support for long output

### raw

Unmodified API response:

```bash
ant messages create ... --format raw
```

Useful when you need the exact bytes returned by the API.

### gjson

[GJSON](https://github.com/tidwall/gjson) transformed output for data extraction:

```bash
# Extract specific field
ant models list --format json --transform "data.0.id"
# Output: "claude-sonnet-4-5-20250929"

# Extract array of values
ant models list --format json --transform "data.#.display_name"
# Output: ["Claude Sonnet 4.5", "Claude Opus 4", ...]

# Use wildcards
ant models list --format json --transform "data.*.id"
```

## GJSON Path Syntax

Common GJSON patterns:

| Path | Description |
|------|-------------|
| `name` | Access object field |
| `data.0` | First array element |
| `data.#` | Array length |
| `data.#.name` | Map over array |
| `data.*` | All array elements |
| `data.?(@.type=="text")` | Filter by condition |
| `data.0.name\|@pretty` | Pretty print |

## Error Formatting

Control error output separately with `--format-error`:

```bash
ant messages create ... \
  --format json \
  --format-error yaml
```

This outputs successful responses as JSON and errors as YAML.

## Transformations

Apply GJSON transformations to extract specific data:

```bash
# Get just the content from a message response
ant messages create \
  --max-tokens 100 \
  --message "Say hello" \
  --model claude-sonnet-4-5-20250929 \
  --format json \
  --transform "content.0.text"

# Get token counts
ant messages count-tokens ... --transform "input_tokens"

# Batch results
ant messages:batches results <batch-id> --transform "request_counts"
```

## Piping and Scripting

JSON format works seamlessly with other tools:

```bash
# Get all model IDs
ant models list --format json | \
  jq -r '.data[].id' | \
  while read model; do
    echo "Model: $model"
  done

# Save response to file
ant messages create ... --format json > response.json

# Pretty print saved response
cat response.json | jq .
```

## Interactive JSON Explorer

For exploring JSON responses interactively:

```bash
# Open in interactive JSON viewer
ant models list --format json | fx

# Or using fzf for selection
ant models list --format json | jq -r '.data[].id' | fzf
```

## Format Selection Tips

| Use Case | Recommended Format |
|----------|-------------------|
| Interactive use | `auto` or `pretty` |
| Shell scripts | `json` |
| Configuration files | `yaml` |
| Data extraction | `json` + `--transform` |
| Debugging | `raw` or `json` |
| API integration | `json` |

## See Also

- [Global Flags](./global-flags.md)
- [GJSON Documentation](https://github.com/tidwall/gjson)
- [jq Documentation](https://stedolan.github.io/jq/)
