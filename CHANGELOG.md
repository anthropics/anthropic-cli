# Changelog

## 1.0.0 (2026-04-08)

Full Changelog: [v0.0.1-alpha.0...v1.0.0](https://github.com/anthropics/anthropic-cli/compare/v0.0.1-alpha.0...v1.0.0)

### Features

* add default description for enum CLI flags without an explicit description ([b037e4b](https://github.com/anthropics/anthropic-cli/commit/b037e4bd14761d89d94afad93406dc0c32414dc4))
* allow `-` as value representing stdin to binary-only file parameters in CLIs ([9ec723c](https://github.com/anthropics/anthropic-cli/commit/9ec723c65d819753e41490c036be6c27bc9f1a1e))
* **api:** add support for Claude Managed Agents ([a349104](https://github.com/anthropics/anthropic-cli/commit/a34910465acfaf9aba5accb2d4772c9956dd1a41))
* **api:** Add support for claude-mythos-preview ([4e24fe0](https://github.com/anthropics/anthropic-cli/commit/4e24fe0609b6f6f283164543a46e5a70f357d07a))
* **api:** fix(cli): allow plain string for --system flag in CLI ([c8f18fb](https://github.com/anthropics/anthropic-cli/commit/c8f18fb7e268327f39c5c1d5211ea90383661b01))
* **api:** GA thinking-display-setting ([bd11386](https://github.com/anthropics/anthropic-cli/commit/bd1138666aefc971433765d4a0ed6c576517b5b8))
* better error message if scheme forgotten in CLI `*_BASE_URL`/`--base-url` ([8e14a7a](https://github.com/anthropics/anthropic-cli/commit/8e14a7a9aef550b1ecf54bb015adc46d51c180c0))
* binary-only parameters become CLI flags that take filenames only ([b6f6dfb](https://github.com/anthropics/anthropic-cli/commit/b6f6dfbacc6418793190d0d422b918e041042a8f))
* **client:** add max_per_page to CLI pagination ([2e335c6](https://github.com/anthropics/anthropic-cli/commit/2e335c688079e8c2605f3d60cf948d5003a80c59))
* set CLI flag constant values automatically where `x-stainless-const` is set ([a3861cf](https://github.com/anthropics/anthropic-cli/commit/a3861cffcfa8759c90745380e485cde13b67d074))
* **tests:** update mock server ([5f28d76](https://github.com/anthropics/anthropic-cli/commit/5f28d76bf358465f6d7543f42b79605aa997898f))


### Bug Fixes

* avoid reading from stdin unless request body is form encoded or json ([dcb5ae2](https://github.com/anthropics/anthropic-cli/commit/dcb5ae2ce236409a7ecdf668a5dde1537513ba1b))
* better support passing client args in any position ([3983a85](https://github.com/anthropics/anthropic-cli/commit/3983a858e5c7b5f616ec0d0d899eeba16918ac2a))
* cli no longer hangs when stdin is attached to a pipe with empty input ([40af100](https://github.com/anthropics/anthropic-cli/commit/40af10034e4940b3228204cbd47d9b18f1907a8d))
* fall back to main branch if linking fails in CI ([edb41a0](https://github.com/anthropics/anthropic-cli/commit/edb41a08f7d4b4db3c3d1a7db790896ba1bd9cca))
* fix for encoding arrays with `any` type items ([5324c86](https://github.com/anthropics/anthropic-cli/commit/5324c86ae85a43bd640c58bdfcdcc05b5acae7a1))
* fix for off-by-one error in pagination logic ([1944a1d](https://github.com/anthropics/anthropic-cli/commit/1944a1d5c045c78c4099f1095fac346710a46b0e))
* fix for test cases with newlines in YAML and better error reporting ([0275df0](https://github.com/anthropics/anthropic-cli/commit/0275df0ae6471dee6c8e01ca4e638c79e6781cb4))
* fix quoting typo ([fe96c27](https://github.com/anthropics/anthropic-cli/commit/fe96c278a5b170805be7de2c0572638562b4d4c3))
* **format:** run go format ([874b72d](https://github.com/anthropics/anthropic-cli/commit/874b72d93ec57b0aaa18560ac1075814a4d6993f))
* handle empty data set using `--format explore` ([a8c535e](https://github.com/anthropics/anthropic-cli/commit/a8c535e2f3c081f6661efd03c8813e4329065c7c))
* improve linking behavior when developing on a branch not in the Go SDK ([ea78613](https://github.com/anthropics/anthropic-cli/commit/ea786137c464e1913d2a7256f25843ad5fd1f3fe))
* improved workflow for developing on branches ([d11462b](https://github.com/anthropics/anthropic-cli/commit/d11462b8a5ea7fbbcaa5a1ad97c44ad20043d836))
* no longer require an API key when building on production repos ([efd44d0](https://github.com/anthropics/anthropic-cli/commit/efd44d0a7e89ef41fe6c434426a5ae2854f18731))
* only set client options when the corresponding CLI flag or env var is explicitly set ([11ac26c](https://github.com/anthropics/anthropic-cli/commit/11ac26c393578e28b53bd7814a9bc3940efcf992))
* use `RawJSON` when iterating items with `--format explore` in the CLI ([1a48e03](https://github.com/anthropics/anthropic-cli/commit/1a48e031144b42a37072bac31c0fc63e715f45e6))


### Chores

* **ci:** run builds on CI even if only spec metadata changed ([62117e9](https://github.com/anthropics/anthropic-cli/commit/62117e98cb1829d7c0980ecb4907e1be264d21c0))
* **ci:** skip lint on metadata-only changes ([4916eef](https://github.com/anthropics/anthropic-cli/commit/4916eef051a5a0e4adb21d5cf331eb6014104aad))
* **cli:** Claude Developer Platform -&gt; Claude Platform ([34be728](https://github.com/anthropics/anthropic-cli/commit/34be728d1b2bd96fd366645b9b01fb6f73dfd8f8))
* **client:** update anthropic-go dependency ([7b311e7](https://github.com/anthropics/anthropic-cli/commit/7b311e76825f362af28ce972f4de5461d5502535))
* **internal:** codegen related update ([d6bd0cd](https://github.com/anthropics/anthropic-cli/commit/d6bd0cd2f68f434f44130a665db9108f0db1d531))
* **internal:** codegen related update ([2162128](https://github.com/anthropics/anthropic-cli/commit/2162128dfcf39ee5c193a1204483b8b95fd2da5b))
* **internal:** codegen related update ([f3a3f4e](https://github.com/anthropics/anthropic-cli/commit/f3a3f4e880cf663ea16dd34c4cb69315bb0a2686))
* **internal:** codegen related update ([0628e7c](https://github.com/anthropics/anthropic-cli/commit/0628e7c3371f528625f3074acb3ec26c804c0299))
* **internal:** codegen related update ([964b642](https://github.com/anthropics/anthropic-cli/commit/964b64259b3642ad755986349cf0a656aaa67f6c))
* **internal:** codegen related update ([51d6aea](https://github.com/anthropics/anthropic-cli/commit/51d6aeabaff93e396d3b18c4edb4d3a613719a30))
* **internal:** codegen related update ([9099f84](https://github.com/anthropics/anthropic-cli/commit/9099f846da2a2531bde46640a46b61bfff4529ce))
* **internal:** codegen related update ([c9e4450](https://github.com/anthropics/anthropic-cli/commit/c9e4450c7a02c98fdac0a318a76981218a96e3a7))
* **internal:** codegen related update ([fdd332c](https://github.com/anthropics/anthropic-cli/commit/fdd332cf32ad9d512e335a4dc7a636c611661702))
* **internal:** regenerate SDK with no functional changes ([913c2ac](https://github.com/anthropics/anthropic-cli/commit/913c2ac9c41200bb37a40b493cf68c9c3ac52caa))
* **internal:** tweak CI branches ([109ba91](https://github.com/anthropics/anthropic-cli/commit/109ba911f37909263884daddf3c899041892171d))
* **internal:** update gitignore ([389d78b](https://github.com/anthropics/anthropic-cli/commit/389d78bd72d8cdaa8718c7ae380a8092d520da84))
* **internal:** update multipart form array serialization ([2196f64](https://github.com/anthropics/anthropic-cli/commit/2196f64bfcf2ff52b5128d74333045fe5b52761f))
* mark all CLI-related tests in Go with `t.Parallel()` ([0cc4ad3](https://github.com/anthropics/anthropic-cli/commit/0cc4ad3c17fe8ba4b6eff3a193d53a92b1b1c249))
* modify CLI tests to inject stdout so mutating `os.Stdout` isn't necessary ([9a08ea8](https://github.com/anthropics/anthropic-cli/commit/9a08ea8a5a3bda1a54468cb280de7e5ae41ba419))
* omit full usage information when missing required CLI parameters ([e357ef3](https://github.com/anthropics/anthropic-cli/commit/e357ef3734887a107316f2bc418f450a8c625f83))
* switch some CLI Go tests from `os.Chdir` to `t.Chdir` ([5f89a21](https://github.com/anthropics/anthropic-cli/commit/5f89a2186990a385dd02c1bd8df0295a94b12539))
* sync repo ([59921d0](https://github.com/anthropics/anthropic-cli/commit/59921d0776c15cd7fff1abc5f4b23248896d34e6))
* **tests:** bump steady to v0.19.4 ([eedbcb0](https://github.com/anthropics/anthropic-cli/commit/eedbcb09708031e77f1d3de0625e6d7b15d25c3c))
* **tests:** bump steady to v0.19.5 ([85f4306](https://github.com/anthropics/anthropic-cli/commit/85f4306aa9016c9ca0324ae3e43dc99252cd5b94))
* **tests:** bump steady to v0.19.6 ([adc7010](https://github.com/anthropics/anthropic-cli/commit/adc70107587e7140467b921b210aa5b722468e97))
* **tests:** bump steady to v0.19.7 ([91c5fb6](https://github.com/anthropics/anthropic-cli/commit/91c5fb60fe7e841da6f8955a46a74b0c466e7173))
* **tests:** bump steady to v0.20.1 ([158978f](https://github.com/anthropics/anthropic-cli/commit/158978f65ab0eedfe91b123f5a7330d0cdb40489))
* **tests:** bump steady to v0.20.2 ([fc343ef](https://github.com/anthropics/anthropic-cli/commit/fc343effce3cdbcc5eea731f4278c8e05671c2d5))
* **tests:** unskip tests that are now supported in steady ([f27a15a](https://github.com/anthropics/anthropic-cli/commit/f27a15af932b3cb084a3a300cb0ba7899c2df484))
* update SDK settings ([7879944](https://github.com/anthropics/anthropic-cli/commit/7879944394728be533ea283a8748fdeab200d2f6))
* zip READMEs as part of build artifact ([387bc97](https://github.com/anthropics/anthropic-cli/commit/387bc973fa098390c3f78d72bba3e97167d735ae))
