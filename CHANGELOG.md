# Changelog

## 1.5.0 (2026-04-30)

Full Changelog: [v1.4.0...v1.5.0](https://github.com/anthropics/anthropic-cli/compare/v1.4.0...v1.5.0)

### Features

* support passing path and query params over stdin ([6f7d931](https://github.com/anthropics/anthropic-cli/commit/6f7d9319d698f71b87452b60c747fa4f2ab14399))


### Bug Fixes

* **cli:** correctly load zsh autocompletion ([9c22b63](https://github.com/anthropics/anthropic-cli/commit/9c22b6376753a581feaa494cff95baf059894d13))
* flags for nullable body scalar fields are strictly typed ([0a9ce73](https://github.com/anthropics/anthropic-cli/commit/0a9ce739679889d996ae331b4e8852bdebb2b74d))

## 1.4.0 (2026-04-28)

Full Changelog: [v1.3.2...v1.4.0](https://github.com/anthropics/anthropic-cli/compare/v1.3.2...v1.4.0)

### Features

* **api:** improve Managed Agents APIs ([234e3e7](https://github.com/anthropics/anthropic-cli/commit/234e3e7a32b7d8f3d54d5b2d2b84e031dcf89cd0))

## 1.3.2 (2026-04-23)

Full Changelog: [v1.3.1...v1.3.2](https://github.com/anthropics/anthropic-cli/compare/v1.3.1...v1.3.2)

## 1.3.1 (2026-04-23)

Full Changelog: [v1.3.0...v1.3.1](https://github.com/anthropics/anthropic-cli/compare/v1.3.0...v1.3.1)

## 1.3.0 (2026-04-23)

Full Changelog: [v1.2.1...v1.3.0](https://github.com/anthropics/anthropic-cli/compare/v1.2.1...v1.3.0)

### Features

* **api:** CMA Memory public beta ([926d861](https://github.com/anthropics/anthropic-cli/commit/926d861abb348b8279019e2c95822592b2dd68f7))
* **cli:** add `--raw-output`/`-r` option to print raw (non-JSON) strings ([534b5fa](https://github.com/anthropics/anthropic-cli/commit/534b5faf893852a5e998e7cd271046125e013c03))
* **cli:** default to interactive explore format for retrieve/list commands when connected to TTY ([929e904](https://github.com/anthropics/anthropic-cli/commit/929e90484002d1559481c214370889c9ecc141b4))
* **cli:** send filename and content type when reading input from files ([7ee2c07](https://github.com/anthropics/anthropic-cli/commit/7ee2c07de4b7f18a41b63610565d1e0be961e782))


### Bug Fixes

* **api:** restore missing features ([55d84d1](https://github.com/anthropics/anthropic-cli/commit/55d84d11a16c9cb4520922d94561d04867680a35))
* correctly serialize --file parameter ([fb0f8bb](https://github.com/anthropics/anthropic-cli/commit/fb0f8bb44c219aba1fdc09be9189a4c65c8c5223))


### Chores

* **ci:** add github env support for goreleaser ([9064b6f](https://github.com/anthropics/anthropic-cli/commit/9064b6f8a8a75e3bc3dc7bc75b2c477f3e8e944c))
* **client:** config fixes ([a89b08b](https://github.com/anthropics/anthropic-cli/commit/a89b08b0be1e6b33398621d05bff17ba303031e4))
* **cli:** use `ShowJSONOpts` as argument to `formatJSON` instead of many positionals ([a9ac5e1](https://github.com/anthropics/anthropic-cli/commit/a9ac5e1e6a65abfbb1001aabdb7d804e723644d4))
* **internal:** more robust bootstrap script ([a453fd0](https://github.com/anthropics/anthropic-cli/commit/a453fd0a7e437049a2a88d5632cbd50743587c46))
* **tests:** bump steady to v0.22.1 ([7aca529](https://github.com/anthropics/anthropic-cli/commit/7aca529007551b5cca9bb20b916e3d238bfb0342))

## 1.2.1 (2026-04-16)

Full Changelog: [v1.2.0...v1.2.1](https://github.com/anthropics/anthropic-cli/compare/v1.2.0...v1.2.1)

### Bug Fixes

* **goreleaser:** correct pull request config ([5d6c9ae](https://github.com/anthropics/anthropic-cli/commit/5d6c9aeb2f5552d9bf698b9a7f5b9e8c36c001bd))


### Chores

* **ci:** support manually triggering release workflow ([9f53a96](https://github.com/anthropics/anthropic-cli/commit/9f53a96d298248c45822057a6c53ec6dd3e3f768))

## 1.2.0 (2026-04-16)

Full Changelog: [v1.1.0...v1.2.0](https://github.com/anthropics/anthropic-cli/compare/v1.1.0...v1.2.0)

### Features

* **api:** add claude-opus-4-7, token budgets and user_profiles ([df20ce9](https://github.com/anthropics/anthropic-cli/commit/df20ce924956647bf362c198dbcadb426c686662))
* **api:** manual updates ([07273ef](https://github.com/anthropics/anthropic-cli/commit/07273ef2e27993e452db24cfdb59088989349c9f))
* **cli:** alias parameters in data with `x-stainless-cli-data-alias` ([991b8e9](https://github.com/anthropics/anthropic-cli/commit/991b8e972802e2ec3ca5663ab0c6fb31ead8a4df))


### Bug Fixes

* **cli:** fix incompatible Go types for flag generated as array of maps ([ced5845](https://github.com/anthropics/anthropic-cli/commit/ced58459c9d668fdde293adeb4ed676e5c73b800))
* fix for failing to drop invalid module replace in link script ([ad79ded](https://github.com/anthropics/anthropic-cli/commit/ad79ded899364b5e8cb288d90597fd4b7984e538))
* use correct multipart array format ([326a8b5](https://github.com/anthropics/anthropic-cli/commit/326a8b5ae00259c439cf0ea613d57fd41babc602))


### Chores

* add documentation for ./scripts/link ([d1a18e2](https://github.com/anthropics/anthropic-cli/commit/d1a18e23681a821cd3d626bc73d9ad2750e465ab))
* **ci:** remove release-doctor workflow ([2c92e20](https://github.com/anthropics/anthropic-cli/commit/2c92e20fdd01bb42f6051c668cdb7be544ade2d7))
* **cli:** additional test cases for `ShowJSONIterator` ([9c94055](https://github.com/anthropics/anthropic-cli/commit/9c94055e3e651cc383e1022ab3cc1c5474d46167))
* **cli:** fall back to JSON when using default "explore" with non-TTY ([cd58bd2](https://github.com/anthropics/anthropic-cli/commit/cd58bd23c08c9716aa7c73d789b3cbe1662ed9cf))
* **cli:** switch long lists of positional args over to param structs ([4373b01](https://github.com/anthropics/anthropic-cli/commit/4373b01d3c239fdb6dd2bcbc2620da89065fb4ad))
* **internal:** codegen related update ([8ea4789](https://github.com/anthropics/anthropic-cli/commit/8ea4789ee4e687b3f82609f4e8ba90d667fde294))
* **internal:** codegen related update ([48aff04](https://github.com/anthropics/anthropic-cli/commit/48aff040e5a6b166e7f4d0f9073e15dbab875a3d))


### Documentation

* update examples ([3213488](https://github.com/anthropics/anthropic-cli/commit/3213488ea69fab6b47e2cef8c807b26961d857ee))

## 1.1.0 (2026-04-09)

Full Changelog: [v1.0.0...v1.1.0](https://github.com/anthropics/anthropic-cli/compare/v1.0.0...v1.1.0)

### Features

* **api:** manual updates ([0563971](https://github.com/anthropics/anthropic-cli/commit/0563971f7ecbb7a0abe9c7ad4131ce0ec7891b2b))


### Chores

* **cli:** let `--format raw` be used in conjunction with `--transform` ([4748f25](https://github.com/anthropics/anthropic-cli/commit/4748f255fd1e151019115e8e2ed37e0c7a56a607))

## 1.0.0 (2026-04-08)

Full Changelog: [v0.0.1-alpha.0...v1.0.0](https://github.com/anthropics/anthropic-cli/compare/v0.0.1-alpha.0...v1.0.0)

### Features

- Initial release of the `ant` CLI.
