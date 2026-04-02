# Quick Start Guide

Get up and running with the Anthropic CLI in 5 minutes.

## Prerequisites

- Anthropic CLI installed (`ant` command available)
- Anthropic API key ([get one here](https://console.anthropic.com))

## Step 1: Set Up Authentication

Export your API key:

```bash
export ANTHROPIC_API_KEY="sk-ant-api03-..."
```

Add to your shell profile for persistence:

```bash
echo 'export ANTHROPIC_API_KEY="your-api-key"' >> ~/.zshrc
source ~/.zshrc
```

## Step 2: Verify Installation

Check that everything is working:

```bash
ant --version
ant models list
```

You should see a list of available Claude models.

## Step 3: Create Your First Message

Create a simple message with Claude:

```bash
ant messages create \
  --max-tokens 1024 \
  --message '{"role": "user", "content": [{"type": "text", "text": "Hello! Can you explain what you can help me with?"}]}' \
  --model claude-sonnet-4-5-20250929
```

## Step 4: Try Different Output Formats

The CLI supports multiple output formats:

```bash
# Pretty formatted (default for terminal)
ant messages create ... --format pretty

# JSON for scripting
ant messages create ... --format json

# YAML for readability
ant messages create ... --format yaml
```

## Step 5: Extract Specific Data

Use GJSON transformations to extract exactly what you need:

```bash
# Get just the response text
ant messages create \
  --max-tokens 100 \
  --message '{"role": "user", "content": [{"type": "text", "text": "Say hello!"}]}' \
  --model claude-sonnet-4-5-20250929 \
  --format json \
  --transform "content.0.text"
```

## Common Tasks

### List Available Models

```bash
ant models list
```

### Count Tokens

Estimate token count before making a request:

```bash
ant messages count-tokens \
  --message '{"role": "user", "content": [{"type": "text", "text": "Hello, world!"}]}' \
  --model claude-sonnet-4-5-20250929
```

### Batch Processing

Process multiple messages efficiently:

```bash
# Create a batch
ant messages:batches create \
  --requests requests.jsonl \
  --endpoint /v1/messages

# Check status
ant messages:batches retrieve <batch-id>

# Get results
ant messages:batches results <batch-id>
```

## Next Steps

- Learn about [authentication](../cli/authentication.md) options
- Explore [global flags](../cli/global-flags.md) for customization
- Read about [output formats](../cli/output-formats.md)
- Check the [CLI reference](../cli/) for all commands

## Troubleshooting

### "No API key provided" error

Set the `ANTHROPIC_API_KEY` environment variable:

```bash
export ANTHROPIC_API_KEY="your-key"
```

### "Invalid API key" error

Verify your key is correct at [console.anthropic.com](https://console.anthropic.com)

### "Command not found" error

Ensure the CLI is installed and in your PATH:

```bash
which ant
```

If not found, reinstall following the [installation guide](../cli/installation.md).
