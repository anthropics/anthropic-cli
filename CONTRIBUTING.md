## Contributing to documentation

The documentation for the CLI lives at [platform.claude.com/docs/en/api/sdks/cli](https://platform.claude.com/docs/en/api/sdks/cli). To suggest changes, open an issue.

## Setting up the environment

To set up the repository, run:

```sh
$ ./scripts/bootstrap
$ ./scripts/lint
```

This will install all the required dependencies and build the CLI.


## Modifying/Adding code

Most of the CLI's command code is generated from the API spec. Modifications to generated code will be persisted between generations, but may result in merge conflicts between manual patches and changes from the generator.

## Running the CLI locally

Use the `scripts/run` script to build and run the CLI from source:

```sh
$ ./scripts/run messages create --help
```

To produce a standalone binary instead, run `./scripts/build`, which places an `ant` binary in the repository root.

## Linking different Go SDK versions

The CLI is built on the [Anthropic Go SDK](https://github.com/anthropics/anthropic-sdk-go). You can link the CLI against a different version of the SDK using the `./scripts/link` script.

To link to a specific version from a repository (version can be a branch, git tag, or commit hash):

```sh
$ ./scripts/link github.com/org/repo@version
```

To link to a local copy of the SDK:

```sh
$ ./scripts/link ../path/to/anthropic-go
```

If you run the link script without any arguments, it will default to `../anthropic-go`. Run `./scripts/unlink` to undo.

## Running tests

Tests run against a mock server set up from the OpenAPI spec. Most of the time the test script manages the mock server for you:

```sh
$ ./scripts/test
```

To run the mock server yourself (for example, to run the CLI against it manually), use:

```sh
$ ./scripts/mock
```

## Formatting

This repository uses the standard gofmt code formatter:

```sh
$ ./scripts/format
```

## Re-recording the README demo

The demo GIF at the top of the README is recorded with [VHS](https://github.com/charmbracelet/vhs) from [`.github/demo.tape`](.github/demo.tape). After changing the tape or the examples it shows, re-record it with:

```sh
$ ./scripts/record-demo
```
