# Changelog

## 1.12.2 (2026-06-24)

Full Changelog: [v1.12.1...v1.12.2](https://github.com/anthropics/anthropic-cli/compare/v1.12.1...v1.12.2)

### Chores

* **api:** add support for sending User Profile ID in request headers ([2a9ec20](https://github.com/anthropics/anthropic-cli/commit/2a9ec20801996eb19686496cebfbdfb3a9420d37))
* **internal:** codegen related update ([70616f5](https://github.com/anthropics/anthropic-cli/commit/70616f54b54e6a48daab449acb3a9e99c42508da))


### Documentation

* **model-release:** model-ID-launch checklist ([62562f3](https://github.com/anthropics/anthropic-cli/commit/62562f3cf05adb984e844d831d2db8d461157775))

## 1.12.1 (2026-06-10)

Full Changelog: [v1.12.0...v1.12.1](https://github.com/anthropics/anthropic-cli/compare/v1.12.0...v1.12.1)

### Bug Fixes

* **api:** add `frontier_llm` refusal category ([#10](https://github.com/anthropics/anthropic-cli/issues/10)) ([7144269](https://github.com/anthropics/anthropic-cli/commit/71442696ce0fe6ad3da8128aa1edab006d7dc6e5))

## 1.12.0 (2026-06-09)

Full Changelog: [v1.11.0...v1.12.0](https://github.com/anthropics/anthropic-cli/compare/v1.11.0...v1.12.0)

### Features

* **api:** add support for Managed Agents deployments and environment variable credentials ([b36672a](https://github.com/anthropics/anthropic-cli/commit/b36672a68b8b8b5962e4a56c1ac8be6f7325da22))
* **api:** manual updates ([d94edec](https://github.com/anthropics/anthropic-cli/commit/d94edec20eef87f0e1c98611f1c857379667574f))

## 1.11.0 (2026-06-09)

Full Changelog: [v1.10.0...v1.11.0](https://github.com/anthropics/anthropic-cli/compare/v1.10.0...v1.11.0)

### Features

* **api:** add support for claude-mythos-5 and claude-fable-5, with support for server-side fallbacks on refusal ([8302a45](https://github.com/anthropics/anthropic-cli/commit/8302a45f5d5a6df0289f9061e9cd6318d7530c63))
* **api:** manual updates ([c5d792d](https://github.com/anthropics/anthropic-cli/commit/c5d792d52248a85aa9adf7d33f8ae7bd88ae949a))
* **api:** small updates to Managed Agents types ([242c693](https://github.com/anthropics/anthropic-cli/commit/242c69330b3c195defc83ffded8eeb634c136968))
* **client:** adds client-side fallbacks middleware for API providers that do not support server-side fallbacks ([8302a45](https://github.com/anthropics/anthropic-cli/commit/8302a45f5d5a6df0289f9061e9cd6318d7530c63))


### Chores

* **internal:** fix artifact url ([20689c3](https://github.com/anthropics/anthropic-cli/commit/20689c394a7eae8edc47a2070d6b65259bd9f1d5))
* **internal:** fix branch names ([2168a5a](https://github.com/anthropics/anthropic-cli/commit/2168a5a12630129de4cb05e75447ad9d08185d05))
* **internal:** update private repo name ([c4f0954](https://github.com/anthropics/anthropic-cli/commit/c4f09545a37a0401ee14fee042585ca0e1fefe80))


### Documentation

* add animated terminal demo ([7573585](https://github.com/anthropics/anthropic-cli/commit/75735853716522fe7363f6df7efe45e89bf76df9))
* move development docs into CONTRIBUTING.md ([ab8aed7](https://github.com/anthropics/anthropic-cli/commit/ab8aed756a4edeb242f5470f643eb30348e0a12c))
* point security reports to Anthropic's HackerOne program ([#5](https://github.com/anthropics/anthropic-cli/issues/5)) ([f5c8f0f](https://github.com/anthropics/anthropic-cli/commit/f5c8f0faa5eeee57e5c644cf09465a0d48852aad))
* release-branch readiness sweep (model-release Step 5) + coverage-audit refinements ([2454db4](https://github.com/anthropics/anthropic-cli/commit/2454db4a3c1a2d77761c48a3285c0850ff547eea))
* replace README demo GIF with webm and poster frame ([46d7f93](https://github.com/anthropics/anthropic-cli/commit/46d7f93d32750879f84a15be69e6b674891a7f05))
* restore README demo GIF ([46e0b71](https://github.com/anthropics/anthropic-cli/commit/46e0b71e0cd581cc3764ff1970dd02d32caf6c00))
* restructure README to match SDK conventions ([e22f218](https://github.com/anthropics/anthropic-cli/commit/e22f218a4c1b5f6edd955a37e31eb59922c1a7eb))
* update product name to Claude Platform ([fa89021](https://github.com/anthropics/anthropic-cli/commit/fa89021825e94af94ae5df03ce9a1e0fbdc3ad2a))
* use claude-opus-4-8 in examples and lead with the ant name ([8ecc617](https://github.com/anthropics/anthropic-cli/commit/8ecc617d62bd23f1df0c6705b8dd289d2cb32cf1))

## 1.10.0 (2026-05-28)

Full Changelog: [v1.9.3...v1.10.0](https://github.com/anthropics/anthropic-cli/compare/v1.9.3...v1.10.0)

### Features

* **api:** Add support for claude-opus-4-8, mid-conversation system blocks, and usage.output_tokens_details ([1f9b079](https://github.com/anthropics/anthropic-cli/commit/1f9b0794209b7781e36977b970110518e6cd8ac0))


### Documentation

* replace literal newlines ([eddcf54](https://github.com/anthropics/anthropic-cli/commit/eddcf5414b6ec53c557c9f2c5188b7286307eee3))

## 1.9.3 (2026-05-26)

Full Changelog: [v1.9.2...v1.9.3](https://github.com/anthropics/anthropic-cli/compare/v1.9.2...v1.9.3)

### Bug Fixes

* stop clearing org/workspace from profile config on logout ([#406](https://github.com/anthropics/anthropic-cli/issues/406)) ([a9999c8](https://github.com/anthropics/anthropic-cli/commit/a9999c8777fabf5dbb484a556580e134ed1c775c))

## 1.9.2 (2026-05-26)

Full Changelog: [v1.9.1...v1.9.2](https://github.com/anthropics/anthropic-cli/compare/v1.9.1...v1.9.2)

### Bug Fixes

* clear organization/workspace on logout and honor --organization-id on login ([1562d45](https://github.com/anthropics/anthropic-cli/commit/1562d45f0a6a6b74fb1eadc4aae8b5ce96bf561f))
* treat text/plan with format: binary as raw upload ([e7f5c54](https://github.com/anthropics/anthropic-cli/commit/e7f5c54f9306899f6482f1533ed1ec236e3b58cb))

## 1.9.1 (2026-05-19)

Full Changelog: [v1.9.0...v1.9.1](https://github.com/anthropics/anthropic-cli/compare/v1.9.0...v1.9.1)

### Bug Fixes

* **runner:** skip tool calls the worker does not own ([#404](https://github.com/anthropics/anthropic-cli/issues/404)) ([2fbfd7e](https://github.com/anthropics/anthropic-cli/commit/2fbfd7e0c847ea5aa45fb7b5d0b74d3846e4cfbd))

## 1.9.0 (2026-05-19)

Full Changelog: [v1.8.0...v1.9.0](https://github.com/anthropics/anthropic-cli/compare/v1.8.0...v1.9.0)

### Features

* **client:** Add support for self-hosted sandboxes in CMA with sandbox helpers ([cbfbd92](https://github.com/anthropics/anthropic-cli/commit/cbfbd923c900f7d5ff842fe08adab568908ff84c))

## 1.8.0 (2026-05-13)

Full Changelog: [v1.7.1...v1.8.0](https://github.com/anthropics/anthropic-cli/compare/v1.7.1...v1.8.0)

### Features

* **api:** Add support for cache diagnostics beta ([14f198e](https://github.com/anthropics/anthropic-cli/commit/14f198eb9cf5dd742a1e1bef20d8d4f45f47c5bb))


### Chores

* **internal:** codegen related update ([cd25682](https://github.com/anthropics/anthropic-cli/commit/cd256829ab1a17815d87b995301e6410331555f3))

## 1.7.1 (2026-05-09)

Full Changelog: [v1.7.0...v1.7.1](https://github.com/anthropics/anthropic-cli/compare/v1.7.0...v1.7.1)

### Chores

* redact api-key headers in debug logs ([b479ed1](https://github.com/anthropics/anthropic-cli/commit/b479ed15a5b6780d0f58618941f4900ae270c5f7))

## 1.7.0 (2026-05-06)

Full Changelog: [v1.6.0...v1.7.0](https://github.com/anthropics/anthropic-cli/compare/v1.6.0...v1.7.0)

### Features

* **api:** add support for Managed Agents multiagents and outcomes, webhooks, vault validation ([6f6a472](https://github.com/anthropics/anthropic-cli/commit/6f6a472f51c025b3445591c428146156f884df50))


### Chores

* **client:** update go dependency ([298152f](https://github.com/anthropics/anthropic-cli/commit/298152fd53bb9d0d927e71ab623cd2a9787a2ce6))

## 1.6.0 (2026-05-05)

Full Changelog: [v1.5.0...v1.6.0](https://github.com/anthropics/anthropic-cli/compare/v1.5.0...v1.6.0)

### Features

* **client:** allow targeting a workspace for OIDC federation token exchange ([a53595b](https://github.com/anthropics/anthropic-cli/commit/a53595b13775bfd4a2b8676355a8a065cd1dc1b6))

## 1.5.0 (2026-05-04)

Full Changelog: [v1.4.0...v1.5.0](https://github.com/anthropics/anthropic-cli/compare/v1.4.0...v1.5.0)

### Features

* **client:** add Workload Identity Federation, interactive OAuth, and auth profiles ([41be2c4](https://github.com/anthropics/anthropic-cli/commit/41be2c4e3921177c773b214973969df547e90c41))
* support passing path and query params over stdin ([6f7d931](https://github.com/anthropics/anthropic-cli/commit/6f7d9319d698f71b87452b60c747fa4f2ab14399))


### Bug Fixes

* **cli:** correctly load zsh autocompletion ([9c22b63](https://github.com/anthropics/anthropic-cli/commit/9c22b6376753a581feaa494cff95baf059894d13))
* flags for nullable body scalar fields are strictly typed ([0a9ce73](https://github.com/anthropics/anthropic-cli/commit/0a9ce739679889d996ae331b4e8852bdebb2b74d))


### Chores

* **internal:** ignore ambiguous schema diagnostics ([9b999de](https://github.com/anthropics/anthropic-cli/commit/9b999de22123a6a9861d49ec6852c02af612e708))

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
