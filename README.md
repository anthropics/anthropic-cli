# Claude Platform CLI

The official CLI for the [Claude Developer Platform](https://platform.claude.com/docs/en/api).

<!-- x-release-please-start-version -->

## Installation

### Installing with Homebrew

```sh
brew install anthropics/tap/ant
```

### Installing with Go

To test or install the CLI locally, you need [Go](https://go.dev/doc/install) version 1.22 or later installed.

```sh
go install 'github.com/anthropics/anthropic-cli/cmd/ant@latest'
```

Once you have run `go install`, the binary is placed in your Go bin directory:

- **Default location**: `$HOME/go/bin` (or `$GOPATH/bin` if GOPATH is set)
- **Check your path**: Run `go env GOPATH` to see the base directory

If commands aren't found after installation, add the Go bin directory to your PATH:

```sh
# Add to your shell profile (.zshrc, .bashrc, etc.)
export PATH="$PATH:$(go env GOPATH)/bin"
```

<!-- x-release-please-end -->

### Running Locally

After cloning the git repository for this project, you can use the
`scripts/run` script to run the tool locally:

```sh
./scripts/run args...
```

## Usage

The CLI follows a resource-based command structure:

```sh
ant [resource] <command> [flags...]
```

```sh
ant messages create \
  --api-key my-anthropic-api-key \
  --max-tokens 1024 \
  --message '{content: [{text: x, type: text}], role: user}' \
  --model claude-sonnet-4-5-20250929
```

For details about specific commands, use the `--help` flag.

### Environment variables

| Environment variable   | Required | Default value |
| ---------------------- | -------- | ------------- |
| `ANTHROPIC_API_KEY`    | no       | `null`        |
| `ANTHROPIC_AUTH_TOKEN` | no       | `null`        |

### Global flags

- `--api-key` (can also be set with `ANTHROPIC_API_KEY` env var)
- `--auth-token` (can also be set with `ANTHROPIC_AUTH_TOKEN` env var)
- `--help` - Show command line usage
- `--debug` - Enable debug logging (includes HTTP request/response details)
- `--version`, `-v` - Show the CLI version
- `--base-url` - Use a custom API backend URL
- `--format` - Change the output format (`auto`, `explore`, `json`, `jsonl`, `pretty`, `raw`, `yaml`)
- `--format-error` - Change the output format for errors (`auto`, `explore`, `json`, `jsonl`, `pretty`, `raw`, `yaml`)
- `--transform` - Transform the data output using [GJSON syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md)
- `--transform-error` - Transform the error output using [GJSON syntax](https://github.com/tidwall/gjson/blob/master/SYNTAX.md)

### Passing files as arguments

To pass files to your API, you can use the `@myfile.ext` syntax:

```bash
ant <command> --arg @abe.jpg
```

Files can also be passed inside JSON or YAML blobs:

```bash
ant <command> --arg '{image: "@abe.jpg"}'
# Equivalent:
ant <command> <<YAML
arg:
  image: "@abe.jpg"
YAML
```

If you need to pass a string literal that begins with an `@` sign, you can
escape the `@` sign to avoid accidentally passing a file.

```bash
ant <command> --username '\@abe'
```

#### Explicit encoding

For JSON endpoints, the CLI tool does filetype sniffing to determine whether the
file contents should be sent as a string literal (for plain text files) or as a
base64-encoded string literal (for binary files). If you need to explicitly send
the file as either plain text or base64-encoded data, you can use
`@file://myfile.txt` (for string encoding) or `@data://myfile.dat` (for
base64-encoding). Note that absolute paths will begin with `@file://` or
`@data://`, followed by a third `/` (for example, `@file:///tmp/file.txt`).

```bash
ant <command> --arg @data://file.txt
```

## Linking different Go SDK versions

You can link the CLI against a different version of the Anthropic Go SDK
for development purposes using the `./scripts/link` script.

To link to a specific version from a repository (version can be a branch,
git tag, or commit hash):

```bash
./scripts/link github.com/org/repo@version
```

To link to a local copy of the SDK:

```bash
./scripts/link ../path/to/anthropic-go
```

If you run the link script without any arguments, it will default to `../anthropic-go`.
