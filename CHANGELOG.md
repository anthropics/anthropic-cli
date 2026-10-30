# Changelog

## [2.5.0](https://github.com/anthropics/anthropic-cli/compare/v2.4.0...v2.5.0) (2026-10-30)

### Features

* **api:** fill RBAC and spend limit request IDs from piped output


### Bug Fixes

* **cli:** preload the explorer by terminal height, not width


### Chores

* **deps:** build with Go 1.26.10
* **tests:** fix flaky flag tests that shared state

## [2.4.0](https://github.com/anthropics/anthropic-cli/compare/v2.3.0...v2.4.0) (2026-10-26)

### Features

* **api:** add user_profile_id to External Keys and Workspaces


### Bug Fixes

* **client:** honor ANTHROPIC_CUSTOM_HEADERS in API requests


### Chores

* bump go sdk to v2.4.0
* **ci:** remove unused build upload steps
* **docs:** update example user ID in Admin API examples

## [2.3.0](https://github.com/anthropics/anthropic-cli/compare/v2.2.0...v2.3.0) (2026-10-23)

### Features

* **api:** add Claude Tag endpoints for Claude Enterprise organizations


### Bug Fixes

* **apply:** ignore top-level keys that start with `$`, such as `$schema`, in resource files and included data files; a vault file may carry only `$schema` and `$comment`
* **cli:** escape terminal control sequences in rendered output


### Chores

* bump go sdk to v2.3.0
* **client:** stop sending query params hardcoded in endpoint paths
* **docs:** update descriptions of the per-user usage and cost reports
* **internal:** check the repository by owner and public name in workflows
* **internal:** read the Go SDK's repo name from a CI variable
* **internal:** remove preview build upload from CI


### Documentation

* update the README and the example request

## [2.2.0](https://github.com/anthropics/anthropic-cli/compare/v2.1.0...v2.2.0) (2026-10-23)

### Features

* **api:** add support for message threads
* **api:** add usage and cost reports by dimension, mark older reports legacy


### Bug Fixes

* **apply:** print control characters in plan output, warnings and errors as visible escapes
* **sessions:** show the full input of MCP and custom tool calls in the `--web` session viewer
* **sessions:** show the full input of MCP and custom tool calls in the terminal view, including at its approval prompt
* **sessions:** show the input fields a built-in tool's layout leaves out at the terminal view's approval prompt
* **sessions:** in the terminal view, drop bidirectional control characters from session text, and show control and invisible characters in names and tool input as escapes such as `\x1b` or `\u200b`


### Chores

* bump go sdk to v2.2.0
* **ci:** update the runner label for internal CI jobs
* **docs:** fix usage and cost report examples and ending_at docs


### Documentation

* **google-cloud:** update the product name

## [2.1.0](https://github.com/anthropics/anthropic-cli/compare/v2.0.0...v2.1.0) (2026-10-21)

### Features

* **api:** add enforcement status, actions and filter to user profiles


### Chores

* bump go sdk to v2.1.0

## [2.0.0](https://github.com/anthropics/anthropic-cli/compare/v1.44.0...v2.0.0) (2026-10-21)

### ⚠ BREAKING CHANGES

* **api:** add support for new API version 2026-10-21

### Features

* **api:** add a workspace_id filter and the usage_tier scope type to list spend limits
* **api:** add support for new API version 2026-10-21


### Chores

* bump go sdk to v2.0.0


### Documentation

* **api:** remove outdated system role note from messages description

## [1.44.0](https://github.com/anthropics/anthropic-cli/compare/v1.43.0...v1.44.0) (2026-10-20)

### Features

* **api:** add org webhook endpoints and user profile enforcement events


### Chores

* bump go sdk to v1.84.0
* **ci:** run the changelog check on a standard runner

## [1.43.0](https://github.com/anthropics/anthropic-cli/compare/v1.42.0...v1.43.0) (2026-10-19)

### Features

* **api:** add workspace rate-limit updates to the Admin API
* **api:** make User Profiles generally available
* **auth:** the stderr note printed when `ANTHROPIC_API_KEY` or `ANTHROPIC_AUTH_TOKEN` overrides a profile now names the variable, the profile and the fix


### Bug Fixes

* **auth:** `auth status` no longer reports a credential passed as `--api-key` or `--auth-token` as an environment override
* **client:** treat x-amz-security-token as a sensitive header


### Chores

* bump go sdk to v1.83.0
* **ci:** skip the changelog check on draft pull requests
* **ci:** update publish workflow
* **docs:** add credential expiry descriptions to the Managed Agents API reference
* **docs:** correct plugin manifest_version and reach descriptions
* **docs:** improve Managed Agents memory store descriptions
* **docs:** improve Managed Agents vault and credential descriptions

## [1.42.0](https://github.com/anthropics/anthropic-cli/compare/v1.41.0...v1.42.0) (2026-10-15)

### Features

* **api:** add the speed parameter and usage.speed to Messages


### Chores

* bump go sdk to v1.82.0
* **docs:** describe headers on Managed Agents and telemetry endpoints


### Documentation

* **api:** note that some models need a beta header to use speed

## [1.41.0](https://github.com/anthropics/anthropic-cli/compare/v1.40.0...v1.41.0) (2026-10-14)

### Features

* **api:** add execution_identity to Managed Agents agents
* **api:** add organization_type_not_supported error to spend limit lists
* **api:** add telemetry destination test endpoint
* **api:** add telemetry destinations to the Admin API
* **apply:** resolve file paths in a deployment's vault_ids
* **client:** add native support for OpenTelemetry tracing

  Spans are exported when the `OTEL_EXPORTER_OTLP_ENDPOINT` or `OTEL_EXPORTER_OTLP_TRACES_ENDPOINT` environment variable is set, otherwise tracing is off. Learn more in the [documentation](https://platform.claude.com/docs/en/manage-claude/opentelemetry/sdk).


### Bug Fixes

* **apply:** accept null for a field that takes references
* **apply:** clear a deployment's resources when its line is deleted
* **apply:** clear a deployment's vault_ids when its line is deleted
* **client:** treat proxy-authorization as a sensitive header


### Chores

* bump go sdk to v1.81.0
* **docs:** correct prefill, thinking, fallback and model descriptions
* **docs:** document the scope required to assign the managed role
* **docs:** update model examples
* include the license files in the published packages
* **internal:** generate each resource's type beside its fields
* **internal:** generate resource fields as a slice
* **internal:** generate tables of resource fields
* **internal:** narrow the tripwire test rule in the contributor guides


### Documentation

* **api:** clarify when custom tool results can be sent

## [1.40.0](https://github.com/anthropics/anthropic-cli/compare/v1.39.1...v1.40.0) (2026-10-09)

### Features

* **api:** add types for the Chat and Cowork unified analytics metrics
* **api:** add workflows, multiagent configuration and thread status filtering to Managed Agents
* **api:** declare error codes on the Files operations


### Bug Fixes

* **apply:** resolve file paths in an agent's predefined agent lists


### Chores

* bump go sdk to v1.80.0

## [1.39.1](https://github.com/anthropics/anthropic-cli/compare/v1.39.0...v1.39.1) (2026-10-08)

### Bug Fixes

* **auth:** refuse to send login codes or refresh tokens over cleartext http (#152)


### Chores

* bump go sdk to v1.79.1
* **ci:** build against the Go SDK branch of a pull request's base
* **ci:** check that pull requests update the changelog
* **internal:** resolve api-key and auth-token credentials without the root flags (#148)

## 1.39.0 (2026-10-07)

Full Changelog: [v1.38.0...v1.39.0](https://github.com/anthropics/anthropic-cli/compare/v1.38.0...v1.39.0)

### Features

* **api:** add claude-haiku-5-5 and typed computer and browser toolset tool calls ([9a19671](https://github.com/anthropics/anthropic-cli/commit/9a19671db8b1ae9223acb815871360fa114511d8))
* **api:** add include_default parameter to list workspaces ([5d26503](https://github.com/anthropics/anthropic-cli/commit/5d26503b5388e89b90cc778899a019bcc647d00b))
* **api:** add lifecycle stage fields and filter to /v1/models ([86c6973](https://github.com/anthropics/anthropic-cli/commit/86c6973da12ac32f9fe383f60ab89bd21ad38d60))
* **auth:** add ant auth login --admin to request the org:admin scope ([#149](https://github.com/anthropics/anthropic-cli/issues/149)) ([8f7ce5f](https://github.com/anthropics/anthropic-cli/commit/8f7ce5f4d71e3aa35c4cfcd9d6aef562b7672241))
* **cli:** support Claude Platform on Google Cloud ([#151](https://github.com/anthropics/anthropic-cli/issues/151)) ([5f16257](https://github.com/anthropics/anthropic-cli/commit/5f162570bcbdca7279ba00fc9c6ab14b3ccae9f1))


### Bug Fixes

* add help text and types to flags specific to one union variant ([5f16257](https://github.com/anthropics/anthropic-cli/commit/5f162570bcbdca7279ba00fc9c6ab14b3ccae9f1))
* **cli:** accept a bare model id for agents --model, show object flag syntax and sub-flags in help ([5f16257](https://github.com/anthropics/anthropic-cli/commit/5f162570bcbdca7279ba00fc9c6ab14b3ccae9f1))
* **client:** keep a piped value in the request body when a path parameter has the same name ([f0b8001](https://github.com/anthropics/anthropic-cli/commit/f0b8001156b5819e12adb4cb690ccd178d6b42fe))
* **client:** send array header params as a single comma-separated header ([f50220f](https://github.com/anthropics/anthropic-cli/commit/f50220ff88fa9f56f9a8a701e8b1fbb4026f04e5))


### Chores

* **docs:** describe a federation rule's target by its type ([5f16257](https://github.com/anthropics/anthropic-cli/commit/5f162570bcbdca7279ba00fc9c6ab14b3ccae9f1))
* **docs:** fix example IDs in sessions, agents and vault credentials ([5f16257](https://github.com/anthropics/anthropic-cli/commit/5f162570bcbdca7279ba00fc9c6ab14b3ccae9f1))
* **docs:** update Managed Agents multiagent and thread descriptions ([5f16257](https://github.com/anthropics/anthropic-cli/commit/5f162570bcbdca7279ba00fc9c6ab14b3ccae9f1))
* **internal:** add REVIEW.md with review instructions ([5f16257](https://github.com/anthropics/anthropic-cli/commit/5f162570bcbdca7279ba00fc9c6ab14b3ccae9f1))
* **internal:** group imports in two blocks in all files ([5f16257](https://github.com/anthropics/anthropic-cli/commit/5f162570bcbdca7279ba00fc9c6ab14b3ccae9f1))
* **tests:** run tests skipped for a query param bug that is now fixed ([4d8e250](https://github.com/anthropics/anthropic-cli/commit/4d8e25047606f421b8c33609fd5ffb4232b38329))
* **tests:** run tests that an older mock server could not serve ([dc98ed1](https://github.com/anthropics/anthropic-cli/commit/dc98ed15baec9bc8e8ae702c46f2074132e41da0))
* **tests:** run the skipped skill create tests ([#150](https://github.com/anthropics/anthropic-cli/issues/150)) ([5f16257](https://github.com/anthropics/anthropic-cli/commit/5f162570bcbdca7279ba00fc9c6ab14b3ccae9f1))
* **tests:** run the tests skipped for incomplete example responses ([c80fcf5](https://github.com/anthropics/anthropic-cli/commit/c80fcf518c5a60a08310f59099a94b33f62afb7b))


### Documentation

* **api:** state the agent tools limit as 256 ([57cb228](https://github.com/anthropics/anthropic-cli/commit/57cb22820abdf989cb991a290d4a8fc1b0388fe7))

## 1.38.0 (2026-09-30)

Full Changelog: [v1.37.0...v1.38.0](https://github.com/anthropics/anthropic-cli/compare/v1.37.0...v1.38.0)

### Features

* **api:** add list spend limits endpoint ([6ccede8](https://github.com/anthropics/anthropic-cli/commit/6ccede8da03b13efc15a114dc8354faf3b5b9cca))


### Bug Fixes

* **apply:** key the remote skill cache by full commit SHA ([#147](https://github.com/anthropics/anthropic-cli/issues/147)) ([5cfdca2](https://github.com/anthropics/anthropic-cli/commit/5cfdca2d4fa1d470698dec1d283eec68c6c39ab2))

## 1.37.0 (2026-09-30)

Full Changelog: [v1.36.0...v1.37.0](https://github.com/anthropics/anthropic-cli/compare/v1.36.0...v1.37.0)

### Features

* **api:** add Claude Enterprise analytics, spend limits, and RBAC groups and roles to the Admin API ([141981a](https://github.com/anthropics/anthropic-cli/commit/141981ab5659111e1e70fd479c7f4fa4edb84f67))
* **api:** add per-user usage and cost reports to the Admin API analytics ([1ef485b](https://github.com/anthropics/anthropic-cli/commit/1ef485b54a62ef30eafee37dcb5a76a431344bc9))
* **api:** add Plugins and Plugin Marketplaces to the Admin API ([a988455](https://github.com/anthropics/anthropic-cli/commit/a988455d725147df0f9a4652486a0e6b9f55ed42))
* **api:** allow removing a plugin's org-wide installation setting ([e68aebc](https://github.com/anthropics/anthropic-cli/commit/e68aebc39a73e0cd96aab11208d2e83394b4ebf5))
* **api:** Organization API endpoints are now GA ([a96980c](https://github.com/anthropics/anthropic-cli/commit/a96980c82bb0309f0c49ba35cba20e558720e5e3))

## 1.36.0 (2026-09-28)

Full Changelog: [v1.35.0...v1.36.0](https://github.com/anthropics/anthropic-cli/compare/v1.35.0...v1.36.0)

### Features

* **api:** add between_tools thinking type ([bad6c55](https://github.com/anthropics/anthropic-cli/commit/bad6c5502abe9735f08b222487c58312c5e0a938))
* **api:** add claude-sonnet-5-5 ([a91c723](https://github.com/anthropics/anthropic-cli/commit/a91c723af8ca474f9a67f042137a2d205f93402f))
* **api:** add include_inherited and source to workspace rate limits ([15f557f](https://github.com/anthropics/anthropic-cli/commit/15f557f314cff42d340b37eb9d14f4854cbcedf1))
* **api:** add typed event type values to the Managed Agents events list filter ([eb3ae68](https://github.com/anthropics/anthropic-cli/commit/eb3ae68266aa0a3c1f738900445f4273b9dee5a2))
* **api:** cache diagnostics GA — diagnostics on Message / MessageCreateParams ([ef7f12e](https://github.com/anthropics/anthropic-cli/commit/ef7f12efde5413759d7c71110a3514f20f11dcf9))


### Bug Fixes

* **sessions:** read event results through the event variant ([#129](https://github.com/anthropics/anthropic-cli/issues/129)) ([d6c0f7d](https://github.com/anthropics/anthropic-cli/commit/d6c0f7dc2fd2a68e4862d97e0a2b62e591ffb5f8))


### Chores

* bump go sdk to v1.76.0 ([3bcd2e9](https://github.com/anthropics/anthropic-cli/commit/3bcd2e94c2b6d96f7b8576bc5198ebb02ee7487b))
* **ci:** build only the shipped targets in CI ([d0032df](https://github.com/anthropics/anthropic-cli/commit/d0032df1b34a9169ffcded14ca6be6fa79f82180))
* **ci:** choose the CI runner by repository ([9e3bcc9](https://github.com/anthropics/anthropic-cli/commit/9e3bcc90e8678f740aabac3ff83efd4e4f39dd79))
* **docs:** clarify that stream: true returns the raw event stream ([e8d753f](https://github.com/anthropics/anthropic-cli/commit/e8d753f58ce150ad61c3aa9bdb656f60781d9272))
* **internal:** run bootstrap once, after linking the Go SDK, in CI ([54499e2](https://github.com/anthropics/anthropic-cli/commit/54499e21c7f47ddbc6d4f6d6b708847b98ccdedc))


### Documentation

* **api:** prefer each field's own description over its shared type's ([9b6c678](https://github.com/anthropics/anthropic-cli/commit/9b6c6783bc789783a00dd314e25765a62285859f))

## 1.35.0 (2026-09-22)

Full Changelog: [v1.34.0...v1.35.0](https://github.com/anthropics/anthropic-cli/compare/v1.34.0...v1.35.0)

### Features

* **api:** add support for claude-opus-5-5, inline tool definitions and MCP tool-list pinning (beta) ([86bc26c](https://github.com/anthropics/anthropic-cli/commit/86bc26c16459d4ae928964990ab2957e49588888))


### Chores

* bump go sdk to v1.75.0 ([#132](https://github.com/anthropics/anthropic-cli/issues/132)) ([72ce520](https://github.com/anthropics/anthropic-cli/commit/72ce520636792818833927bd43b7d5550e1bbda0))
* **docs:** add descriptions to the Dreams API reference ([69c0904](https://github.com/anthropics/anthropic-cli/commit/69c09042a2a5f50fa36d6bb2b7096ba8454d2184))
* **docs:** add descriptions to the User Profiles API reference ([69c0904](https://github.com/anthropics/anthropic-cli/commit/69c09042a2a5f50fa36d6bb2b7096ba8454d2184))
* **docs:** add memory store descriptions to the Managed Agents API reference ([69c0904](https://github.com/anthropics/anthropic-cli/commit/69c09042a2a5f50fa36d6bb2b7096ba8454d2184))
* **docs:** improve descriptions in the Dreams API reference ([31c4638](https://github.com/anthropics/anthropic-cli/commit/31c4638dbc39e47d0d394daa6d580d8b159b4f62))

## 1.34.0 (2026-09-18)

Full Changelog: [v1.33.0...v1.34.0](https://github.com/anthropics/anthropic-cli/compare/v1.33.0...v1.34.0)

### Features

* **apply:** manage vaults ([#105](https://github.com/anthropics/anthropic-cli/issues/105)) ([8c9396f](https://github.com/anthropics/anthropic-cli/commit/8c9396f5aaf0a8377c8ef1ca94523671529f7d2d))


### Chores

* bump go sdk to v1.74.0 ([#128](https://github.com/anthropics/anthropic-cli/issues/128)) ([3ed8b96](https://github.com/anthropics/anthropic-cli/commit/3ed8b96e154f636cc394fde0a127a596b054326f))
* **ci:** update authentication for fetching the Go SDK in CI ([e2b84f4](https://github.com/anthropics/anthropic-cli/commit/e2b84f46838d85bbd661288a8e37d9725e3dd25e))
* **docs:** remove placeholder parameter descriptions ([62aa6a1](https://github.com/anthropics/anthropic-cli/commit/62aa6a118ca1b4c5929e7ffcd864ea02d0642d8f))
* **internal:** skip clearing the Go VCS cache when its path is empty ([575920d](https://github.com/anthropics/anthropic-cli/commit/575920d38bc126f5c09ba8ac0bc75c0675a0f3bc))


### Documentation

* add path and header parameter descriptions ([09c6a7f](https://github.com/anthropics/anthropic-cli/commit/09c6a7f5ea8187994e1c1c6a22b6d2e493e3b079))

## 1.33.0 (2026-09-15)

Full Changelog: [v1.32.0...v1.33.0](https://github.com/anthropics/anthropic-cli/compare/v1.32.0...v1.33.0)

### Features

* **api:** add auto mode tool permissions for Managed Agents ([8952120](https://github.com/anthropics/anthropic-cli/commit/895212025b4e28022efd3b20fb72b4cc72489e8c))
* **api:** add compaction parameter and signed compaction blocks (beta) ([2b5931e](https://github.com/anthropics/anthropic-cli/commit/2b5931e7ee641beee8211fd61519b69246122ece))
* **api:** add url_sources to the web fetch tool ([81f6bb5](https://github.com/anthropics/anthropic-cli/commit/81f6bb5eb2f0bd213fa26a21a43f284cd01894a9))
* **api:** add workspace_id parameter to user profiles methods ([6c4c959](https://github.com/anthropics/anthropic-cli/commit/6c4c9593553ad8c0ebaca94567d534b2c79f88aa))


### Bug Fixes

* **client:** ignore invalid Retry-After values and validate maxRetries ([8952120](https://github.com/anthropics/anthropic-cli/commit/895212025b4e28022efd3b20fb72b4cc72489e8c))
* **client:** retry connection errors in the async client and stop blocking in the coroutine retry loop ([8952120](https://github.com/anthropics/anthropic-cli/commit/895212025b4e28022efd3b20fb72b4cc72489e8c))
* **client:** use the default backoff when Retry-After is out of range ([b9162c8](https://github.com/anthropics/anthropic-cli/commit/b9162c8275efc99d2bd409500fe84873a684d4a0))
* **worker:** use Go SDK's toolset to enable memory stores ([5199245](https://github.com/anthropics/anthropic-cli/commit/51992458e909515cca20f4866b44681f78040b41))


### Chores

* bump go sdk to v1.73.0 ([#104](https://github.com/anthropics/anthropic-cli/issues/104)) ([526ad85](https://github.com/anthropics/anthropic-cli/commit/526ad85d277b87d7324978e4cbab2e4314f00f75))
* **ci:** pass branch name to link step via env var ([8952120](https://github.com/anthropics/anthropic-cli/commit/895212025b4e28022efd3b20fb72b4cc72489e8c))
* **deps:** pin a test-only indirect dependency for reproducible go.sum ([564ada5](https://github.com/anthropics/anthropic-cli/commit/564ada5fb991cfd3c0772838fadbd9b1c5910e30))
* **docs:** correct the compaction beta's parameter descriptions ([d02eeb2](https://github.com/anthropics/anthropic-cli/commit/d02eeb282c001fd7de5714aad5f3baf726cc600d))
* **tests:** stop the mock server without failing a passing test run ([bec0d1a](https://github.com/anthropics/anthropic-cli/commit/bec0d1a57a3f9ec27bfcd031b273240759ac5b31))


### Documentation

* stop documenting unions with their first variant's description ([9db5cd4](https://github.com/anthropics/anthropic-cli/commit/9db5cd4db9facdef54790099e44fb404a9884d0b))

## 1.32.0 (2026-09-10)

Full Changelog: [v1.31.0...v1.32.0](https://github.com/anthropics/anthropic-cli/compare/v1.31.0...v1.32.0)

### Features

* **api:** add auto mode tool permissions for Managed Agents ([0428ab8](https://github.com/anthropics/anthropic-cli/commit/0428ab889cb7298e84afa3ac8989dfbcae348eb2))
* **api:** add content_too_large web_fetch tool error code ([a138600](https://github.com/anthropics/anthropic-cli/commit/a138600048210b9dd8e31aac489f77e9093bc451))
* **api:** add the user-profiles-2026-09-04 beta value and external_user_details to user profiles ([92c17de](https://github.com/anthropics/anthropic-cli/commit/92c17ded09d933f2aa45b796a10bc4cbb374c890))
* **api:** support mounting public GitHub repositories without an authorization_token in Managed Agents sessions ([a07d5fb](https://github.com/anthropics/anthropic-cli/commit/a07d5fb29650575c72f7d7e92df6cf1e3213d5e2))


### Bug Fixes

* **cli:** quote filenames in shell completions ([eaaec6d](https://github.com/anthropics/anthropic-cli/commit/eaaec6deb633c8223cb3d82d00cfa30412284c4c))


### Chores

* **docs:** correct the environment scope field description ([108b488](https://github.com/anthropics/anthropic-cli/commit/108b488d5f79503ea9c05fc0a179636ad1e808a7))
* **internal:** remove generated file header comments ([97a39e5](https://github.com/anthropics/anthropic-cli/commit/97a39e5ea0056bcc8fa0a172193cf59b2b5458df))
* **internal:** restore package version ([0ccdc5b](https://github.com/anthropics/anthropic-cli/commit/0ccdc5b663cf0e345dec6d14f59e62e72b4c2883))
* **internal:** restore package version ([cee583f](https://github.com/anthropics/anthropic-cli/commit/cee583f5538db3fec6ae06b3fa4fb89953235ad7))
* **internal:** stop stamping the package version into generated files ([859c4e8](https://github.com/anthropics/anthropic-cli/commit/859c4e8ebe980599adf9b33c87759d37be43af70))

## 1.31.0 (2026-09-04)

Full Changelog: [v1.30.0...v1.31.0](https://github.com/anthropics/anthropic-cli/compare/v1.30.0...v1.31.0)

### Features

* **api:** add named types for organization compliance settings state ([4aeee85](https://github.com/anthropics/anthropic-cli/commit/4aeee858da0cbeb92c25db9b1cfb33ee4d83467c))


### Chores

* **ci:** cache Go build output in the build job ([da8ff8c](https://github.com/anthropics/anthropic-cli/commit/da8ff8c877e6a500d4569e12df4363ff81e84053))
* **tests:** reword the skip reason on the path-level query param tests ([fa58532](https://github.com/anthropics/anthropic-cli/commit/fa58532c62b35785be016ce35581e70171eb7702))

## 1.30.0 (2026-09-03)

Full Changelog: [v1.29.0...v1.30.0](https://github.com/anthropics/anthropic-cli/compare/v1.29.0...v1.30.0)

### Features

* add `ant apply` for managing agents, skills, environments, memory stores and deployments as code ([#97](https://github.com/anthropics/anthropic-cli/issues/97)) ([d69ffe9](https://github.com/anthropics/anthropic-cli/commit/d69ffe9911738c1ed6fb6a9d5bed104a9c687e9a))
* **api:** add Claude Tag category and user breakdowns to usage reports ([9ec28f9](https://github.com/anthropics/anthropic-cli/commit/9ec28f93d41ccdde28e2a19fc97032385d31c2e0))
* **api:** add support for sending a workspace ID on more endpoints ([f0cab61](https://github.com/anthropics/anthropic-cli/commit/f0cab612ce1a74c8a015d31db704b34cb56d8842))


### Bug Fixes

* apply dotted sub-flags to union array flags ([c9ad752](https://github.com/anthropics/anthropic-cli/commit/c9ad7520b125468acc8eeb05199524de832cc61b))


### Chores

* bump go sdk to v1.70.1 ([5fc9be6](https://github.com/anthropics/anthropic-cli/commit/5fc9be6c88703f2b7edc97c93b1c239f46a78373))
* **internal:** narrower codeowners scope ([c32c35e](https://github.com/anthropics/anthropic-cli/commit/c32c35e4474e961e1e56ca6ef6728d8dcedeac4a))
* **internal:** revert codeowners change ([58e6f5e](https://github.com/anthropics/anthropic-cli/commit/58e6f5e9ac231275852b2ea617f9045e03db9282))


### Documentation

* **api:** update a few doc strings ([465f12e](https://github.com/anthropics/anthropic-cli/commit/465f12e578147fc4ebacf0f8ee5b47e312460f59))

## 1.29.0 (2026-09-01)

Full Changelog: [v1.28.0...v1.29.0](https://github.com/anthropics/anthropic-cli/compare/v1.28.0...v1.29.0)

### Features

* **api:** beta user profiles: add external_user_onboarded_at, remove relationship in favor of access_type ([b37e2bd](https://github.com/anthropics/anthropic-cli/commit/b37e2bdee6af31d8bf4f47d5f48ea0a2a0b69cd0))
* **api:** manual updates ([f615a64](https://github.com/anthropics/anthropic-cli/commit/f615a64ec3fa42a28fe044f256e95913a23ec9fe))
* **api:** organization compliance settings, user-profile order_by, memory-store and toolset schema updates ([e9f08fe](https://github.com/anthropics/anthropic-cli/commit/e9f08fe262388e185a1d3d94008f988466b2d30f))


### Bug Fixes

* **ci:** update link step to use env var ([3adb1ee](https://github.com/anthropics/anthropic-cli/commit/3adb1eef74a53d5dfa0d171d8d84f3bd8a5f3b30))


### Chores

* bump go dep ([#84](https://github.com/anthropics/anthropic-cli/issues/84)) ([f5bd036](https://github.com/anthropics/anthropic-cli/commit/f5bd0368c559d2c04ed9ae8a5eb78e135a53d306))
* **internal:** bump go dep ([#81](https://github.com/anthropics/anthropic-cli/issues/81)) ([87f2b9c](https://github.com/anthropics/anthropic-cli/commit/87f2b9c70ef0ba5f5bd6d2d3a919c17836102956))

## 1.28.0 (2026-08-27)

Full Changelog: [v1.27.0...v1.28.0](https://github.com/anthropics/anthropic-cli/compare/v1.27.0...v1.28.0)

### Features

* **api:** beta files/skills namespaces use GA shapes; drop dated beta header pins ([a7666a4](https://github.com/anthropics/anthropic-cli/commit/a7666a4455c153456d9bf38ec526bd89c91c33ff))


### Chores

* **internal:** codegen related update ([81209db](https://github.com/anthropics/anthropic-cli/commit/81209dbce25786dc84c4591ab85d9b7f212cff75))
* **internal:** codegen related update ([bf6cc85](https://github.com/anthropics/anthropic-cli/commit/bf6cc857337b83cad96ca0b8d8453daaa73e865c))


### Documentation

* **api:** clarify pagination on the organization rate-limit list endpoints ([32714ce](https://github.com/anthropics/anthropic-cli/commit/32714ce9787dd226faf80b3ca5565b0db0c6fdf0))

## 1.27.0 (2026-08-26)

Full Changelog: [v1.26.1...v1.27.0](https://github.com/anthropics/anthropic-cli/compare/v1.26.1...v1.27.0)

### Features

* **api:** add missing anthropic-beta values ([17d560f](https://github.com/anthropics/anthropic-cli/commit/17d560f04e5fc374347456157a2c10c93162bdeb))
* **api:** add support for Organization API endpoints ([05d89e3](https://github.com/anthropics/anthropic-cli/commit/05d89e37ce1c823aa90225462e628721f3b05afe))
* **auth:** add --api-key-stdin / --auth-token-stdin and deprecate credential argv flags ([#64](https://github.com/anthropics/anthropic-cli/issues/64)) ([c19647d](https://github.com/anthropics/anthropic-cli/commit/c19647d73fb3fba56b7b5ba113bc50687ab3c1dd))


### Bug Fixes

* **auth:** allow login without a workspace id ([cc27104](https://github.com/anthropics/anthropic-cli/commit/cc2710474b52d0df50ed5150007bdb7e5ce390a4))
* **auth:** read global credential/workspace flags from the root command ([#73](https://github.com/anthropics/anthropic-cli/issues/73)) ([e832332](https://github.com/anthropics/anthropic-cli/commit/e832332f555701569b64f0d37301f007b5683043))
* **client:** keep typed repeatable flags for nullable query and header arrays ([0bba330](https://github.com/anthropics/anthropic-cli/commit/0bba330886fa461e50554ab3e5b271b8793740b4))


### Chores

* **client:** sync Go SDK version ([1088d08](https://github.com/anthropics/anthropic-cli/commit/1088d0882e2b2fd76d5fb0cb53d99c42276f5171))
* **docs:** clarify skill version `latest` support and memory version retention ([22fb67a](https://github.com/anthropics/anthropic-cli/commit/22fb67a5a4530f22b89806d77e2983b51e6038bf))
* **internal:** cancel superseded CI runs ([cba857e](https://github.com/anthropics/anthropic-cli/commit/cba857eae253957da40b7a5e61200cb160579848))
* **internal:** minor workflow reformat ([fad8f52](https://github.com/anthropics/anthropic-cli/commit/fad8f528d0926e6b45e1eb086837f91a603e0376))

## 1.26.1 (2026-08-19)

Full Changelog: [v1.26.0...v1.26.1](https://github.com/anthropics/anthropic-cli/compare/v1.26.0...v1.26.1)

### Chores

* bump sdk version ([#66](https://github.com/anthropics/anthropic-cli/issues/66)) ([acdc28b](https://github.com/anthropics/anthropic-cli/commit/acdc28b15a9665e35dde952146d62c57bb6a2d82))

## 1.26.0 (2026-08-19)

Full Changelog: [v1.25.0...v1.26.0](https://github.com/anthropics/anthropic-cli/compare/v1.25.0...v1.26.0)

### Features

* **api:** managed agents web search config and self hosted sandbox memory ([8abb4e9](https://github.com/anthropics/anthropic-cli/commit/8abb4e9528d19640935e17a72217ca7e20ab1294))

## 1.25.0 (2026-08-19)

Full Changelog: [v1.24.0...v1.25.0](https://github.com/anthropics/anthropic-cli/compare/v1.24.0...v1.25.0)

### Features

* **api:** Files and Skills APIs are now GA; add computer use and browser use toolsets ([db63120](https://github.com/anthropics/anthropic-cli/commit/db63120f2f9dee76023f2a5d46a419e8436d8b80))
* **api:** manual updates ([fe7e0a5](https://github.com/anthropics/anthropic-cli/commit/fe7e0a5244d1dc45fdde0bb7b0d71277ee2470d3))

## 1.24.0 (2026-08-18)

Full Changelog: [v1.23.0...v1.24.0](https://github.com/anthropics/anthropic-cli/compare/v1.23.0...v1.24.0)

### Features

* **api:** additions to files and memory stores ([addfe40](https://github.com/anthropics/anthropic-cli/commit/addfe403dc474b5cc859d589c9886fc409714e22))
* **api:** updates to skill, files, and user profiles ([3ebd6ce](https://github.com/anthropics/anthropic-cli/commit/3ebd6ce2b9ab09d069f154f887a7c26c80ebf17c))


### Bug Fixes

* **api:** remove unsupported mid_conv_system content block ([26f42f9](https://github.com/anthropics/anthropic-cli/commit/26f42f9a6b7f980fb81de62a3a65df0f55b48724))


### Chores

* **internal:** remove leftover prism references ([1032305](https://github.com/anthropics/anthropic-cli/commit/10323058cbda781f1047bad1c309926d0e990b2f))

## 1.23.0 (2026-08-13)

Full Changelog: [v1.22.1...v1.23.0](https://github.com/anthropics/anthropic-cli/compare/v1.22.1...v1.23.0)

### Features

* **api:** add output_behavior to dream creation (create a new memory store or update the input store in place) ([c9223e4](https://github.com/anthropics/anthropic-cli/commit/c9223e447d68fa37c77599b570618301d1ef3b9f))


### Bug Fixes

* **ci:** upload ant binaries instead of nonexistent cdp paths ([99bc95e](https://github.com/anthropics/anthropic-cli/commit/99bc95ef086da0b8b25745ab58df8344fd47dde0))


### Chores

* **internal:** codegen related update ([90f3c73](https://github.com/anthropics/anthropic-cli/commit/90f3c7340dcf489df7cf7cb84a00df578c84a7bd))
* **internal:** codegen related update ([adcaf9f](https://github.com/anthropics/anthropic-cli/commit/adcaf9f6a19c3033289c195bcfa7163ab2a1dee2))
* **internal:** migrate legacy openapi transforms to core commands ([2dc0ae4](https://github.com/anthropics/anthropic-cli/commit/2dc0ae435accf3935cc108215682cb4e82dd45cb))


### Documentation

* **api:** clarify that user profile name is optional for resold profiles ([f0619ca](https://github.com/anthropics/anthropic-cli/commit/f0619ca2956ef227f74fb772d9500b6c7d4f0598))

## 1.22.1 (2026-08-07)

Full Changelog: [v1.22.0...v1.22.1](https://github.com/anthropics/anthropic-cli/compare/v1.22.0...v1.22.1)

### Chores

* bump sdk version ([4afa682](https://github.com/anthropics/anthropic-cli/commit/4afa682626075a34cf080c3ba58cf9fd0505eb74))

## 1.22.0 (2026-08-06)

Full Changelog: [v1.21.0...v1.22.0](https://github.com/anthropics/anthropic-cli/compare/v1.21.0...v1.22.0)

### Features

* **api:** add support for session budgets, advisor tool, pinned inference location and skills auto-loading from GitHub ([8aeb8e7](https://github.com/anthropics/anthropic-cli/commit/8aeb8e73454d4d9c54be1ed412cebd050109ce55))

## 1.21.0 (2026-07-24)

Full Changelog: [v1.20.0...v1.21.0](https://github.com/anthropics/anthropic-cli/compare/v1.20.0...v1.21.0)

### Features

* **api:** add claude-opus-5 model ([3e67b3d](https://github.com/anthropics/anthropic-cli/commit/3e67b3d753f87f813fbf6cca58277c723bd15e74))
* **api:** add tool addition/removal blocks and tool_change events ([3e67b3d](https://github.com/anthropics/anthropic-cli/commit/3e67b3d753f87f813fbf6cca58277c723bd15e74))
* **api:** expand client-side fallback credit token types and add server-side fallbacks default option ([3e67b3d](https://github.com/anthropics/anthropic-cli/commit/3e67b3d753f87f813fbf6cca58277c723bd15e74))


### Bug Fixes

* **client:** honor ANTHROPIC_BASE_URL in API requests ([#43](https://github.com/anthropics/anthropic-cli/issues/43)) ([0dcaa4f](https://github.com/anthropics/anthropic-cli/commit/0dcaa4f2badc95cc11f29a8773f13ceb188fbbeb))


### Chores

* bump dependency ([8358acd](https://github.com/anthropics/anthropic-cli/commit/8358acd45ce1313ff15bf4f15cbd14d5b5fe86fc))
* **ci:** cap macOS notarize timeout at 40m; add notary-status debug workflow ([#42](https://github.com/anthropics/anthropic-cli/issues/42)) ([3aa7a80](https://github.com/anthropics/anthropic-cli/commit/3aa7a80a5adf820bdb826ff28beed836e3d08f3a))

## 1.20.0 (2026-07-23)

Full Changelog: [v1.19.0...v1.20.0](https://github.com/anthropics/anthropic-cli/compare/v1.19.0...v1.20.0)

### Features

* **api:** add new stop reason 'model_context_window_exceeded' ([cb2f59e](https://github.com/anthropics/anthropic-cli/commit/cb2f59e48d4b205496534bd7fa1d8f29487fb7a4))
* **api:** manual updates ([8f87e6a](https://github.com/anthropics/anthropic-cli/commit/8f87e6aac8b898b8a37d50cb13d0a568158a98a8))

## 1.19.0 (2026-07-22)

Full Changelog: [v1.18.0...v1.19.0](https://github.com/anthropics/anthropic-cli/compare/v1.18.0...v1.19.0)

### Features

* **api:** add support for Managed Agents model effort, initial session events, and threads delta streaming ([1d61ed6](https://github.com/anthropics/anthropic-cli/commit/1d61ed66d31c4cff85ec9fa64bdd0a9f8755b8ac))


### Chores

* **internal:** codegen related update ([fd751e4](https://github.com/anthropics/anthropic-cli/commit/fd751e48b1f39ba9dacee90fc211b51a8b8a00e0))

## 1.18.0 (2026-07-16)

Full Changelog: [v1.17.0...v1.18.0](https://github.com/anthropics/anthropic-cli/compare/v1.17.0...v1.18.0)

### Features

* **api:** add support for MCP Tunnels ([25b0012](https://github.com/anthropics/anthropic-cli/commit/25b0012f97a040e2f9f7f4c05c87ac2ade31988a))
* **api:** manual updates ([f40ef8a](https://github.com/anthropics/anthropic-cli/commit/f40ef8a3abdf4a651fa0447ec4ef73f9cf15dbc1))

## 1.17.0 (2026-07-10)

Full Changelog: [v1.16.0...v1.17.0](https://github.com/anthropics/anthropic-cli/compare/v1.16.0...v1.17.0)

### Features

* **api:** add support for dreaming ([9ed8015](https://github.com/anthropics/anthropic-cli/commit/9ed8015b7afee21bd0a6e780fa89957978564e77))


### Chores

* **docs:** update model example ([4a5da72](https://github.com/anthropics/anthropic-cli/commit/4a5da72740a9f0aba52292c416547b78995b0806))
* **docs:** updates to descriptions and examples ([f81cd30](https://github.com/anthropics/anthropic-cli/commit/f81cd30432a93fdb628d29017080e5adf9b9f54a))

## 1.16.0 (2026-07-02)

Full Changelog: [v1.15.0...v1.16.0](https://github.com/anthropics/anthropic-cli/compare/v1.15.0...v1.16.0)

### Features

* **api:** add agent-memory-2026-07-22 beta header ([4789498](https://github.com/anthropics/anthropic-cli/commit/4789498f89818f216b03d99ce30da1a5b1af88f3))

## 1.15.0 (2026-07-01)

Full Changelog: [v1.14.1...v1.15.0](https://github.com/anthropics/anthropic-cli/compare/v1.14.1...v1.15.0)

### Features

* **api:** manual updates ([c54265c](https://github.com/anthropics/anthropic-cli/commit/c54265c572d7b20bc750a97c21031f4c5de1c68d))


### Chores

* **api:** remove some nonfunctional types from the SDKs ([c297011](https://github.com/anthropics/anthropic-cli/commit/c2970117b64c724fe1dc81cd0c66190ed0711936))

## 1.14.1 (2026-07-01)

Full Changelog: [v1.14.0...v1.14.1](https://github.com/anthropics/anthropic-cli/compare/v1.14.0...v1.14.1)

### Chores

* bump go sdk version ([c55fa73](https://github.com/anthropics/anthropic-cli/commit/c55fa7383667a03205b149c873a2e4ae9ca62d2d))
* sync ([aa1cf31](https://github.com/anthropics/anthropic-cli/commit/aa1cf31e02dcc4c8fc2696bb9d3be917a2976fa3))

## 1.14.0 (2026-06-30)

Full Changelog: [v1.13.0...v1.14.0](https://github.com/anthropics/anthropic-cli/compare/v1.13.0...v1.14.0)

### Features

* **api:** add support for Managed Agents event delta streaming, agent overrides, reverse pagination, vault credential injection scoping, and agent and deployment webhook events ([9dfd032](https://github.com/anthropics/anthropic-cli/commit/9dfd032a96cdcee357c5a49fa64c33d6fa9a8e70))

## 1.13.0 (2026-06-30)

Full Changelog: [v1.12.2...v1.13.0](https://github.com/anthropics/anthropic-cli/compare/v1.12.2...v1.13.0)

### Features

* add global --workspace-id flag to set the anthropic-workspace-id header ([#6](https://github.com/anthropics/anthropic-cli/issues/6)) ([2479219](https://github.com/anthropics/anthropic-cli/commit/2479219db8ba2eeb9713f6e16f11911a5c9204d1))
* **api:** add support for 20260318 web fetch and support tools ([7bb4c88](https://github.com/anthropics/anthropic-cli/commit/7bb4c8890dababac8de194af8547341784e19c1c))
* **api:** add support for claude-sonnet-5 ([f558ef8](https://github.com/anthropics/anthropic-cli/commit/f558ef8f66d82b4271eb9e6a1c2aadc1b0256cca))
* **api:** manual updates ([05f8243](https://github.com/anthropics/anthropic-cli/commit/05f82439094b1b28bcbca25d8ef31e4743a059a6))
* **api:** manual updates ([f2230c3](https://github.com/anthropics/anthropic-cli/commit/f2230c38f272faec4964a7967a5928baddf5ef3d))


### Chores

* **api:** accept user profile ID's when counting tokens ([ac6bf4b](https://github.com/anthropics/anthropic-cli/commit/ac6bf4b7c0b702d1886394a9bc4c61ab441c7354))
* **docs:** updates to descriptions and example values ([7708932](https://github.com/anthropics/anthropic-cli/commit/77089320c047ef10a0b5a81d00d4bbfd546bcc7a))

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
