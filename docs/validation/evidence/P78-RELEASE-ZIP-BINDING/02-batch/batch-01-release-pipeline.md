# P78 发布链路现状 — 事实枚举（只读，未修改任何文件）

## 1. `.github/workflows/ci.yml` 的 job 清单

`jobs:` 块从 `ci.yml:50` 开始，文件共 408 行，共 4 个 job。

**`test`**（`ci.yml:51-94`）
- `ci.yml:51-52`：`  test:` / `    name: Go ${{ matrix.os }}`
- `if:`：未找到（该 job 无 `if:`）；`needs:`：未找到
- `ci.yml:55-56`：`      matrix:` / `        os: [ubuntu-latest, windows-latest]`
- steps（按出现顺序）：`ci.yml:59` `- name: Checkout source`、`ci.yml:67` `- name: Setup Go`、`ci.yml:73` `- name: Build`（`run: go build ./...`）、`ci.yml:76` `- name: Vet`、`ci.yml:79` `- name: Test`、`ci.yml:82` `- name: Race`（`ci.yml:83` `        if: runner.os == 'Linux'`）、`ci.yml:86` `- name: Contract fixtures`（`ci.yml:87` `        if: runner.os == 'Linux'`）、`ci.yml:93` `- name: Gates`（`ci.yml:94` `        run: node tools/gates/gate.js --selftest && node tools/gates/gate.js --all --scope=ci`）

**`candidate-build`**（`ci.yml:96-144`）
- `ci.yml:96-97`：`  candidate-build:` / `    name: Candidate build (${{ matrix.os }})`
- `ci.yml:98`：`    if: github.event_name == 'workflow_dispatch' && github.ref == 'refs/heads/main'`
- `needs:`：未找到；`ci.yml:101-102`：`      matrix:` / `        os: [ubuntu-latest, windows-latest]`
- `ci.yml:104-105`：`    env:` / `      CANDIDATE_COMMIT: ${{ inputs.candidate_commit }}`
- steps：`ci.yml:107` `- name: Validate candidate pin`、`ci.yml:114` `- name: Checkout candidate`、`ci.yml:121` `- name: Setup Go`、`ci.yml:127` `- name: Build candidate without executing it`、`ci.yml:139` `- name: Upload immutable candidate binary`

**`candidate-probe`**（`ci.yml:146-254`）
- `ci.yml:146-147`：`  candidate-probe:` / `    name: Candidate probe (${{ matrix.os }})`
- `ci.yml:148`：`    needs: candidate-build`
- `ci.yml:149`：`    if: github.event_name == 'workflow_dispatch' && github.ref == 'refs/heads/main'`
- `ci.yml:152-153`：`      matrix:` / `        os: [ubuntu-latest, windows-latest]`
- `ci.yml:155-159`：env 映射 `CANDIDATE_COMMIT`/`VERIFIER_COMMIT`/`NOOP_CHAIN_SHA256`/`EXECUTABLE_CHAIN_SHA256`
- steps：`ci.yml:161` `- name: Validate repository trust pins`、`ci.yml:187` `- name: Checkout trusted verifier`、`ci.yml:194` `- name: Setup Go`、`ci.yml:200` `- name: Build trusted probe harness`、`ci.yml:211` `- name: Download immutable candidate binary`、`ci.yml:218` `- name: Probe untrusted candidate`、`ci.yml:249` `- name: Upload untrusted probe oracle`

**`selfhost`**（`ci.yml:256-408`）
- `ci.yml:256-257`：`  selfhost:` / `    name: Bootstrap and release evidence (${{ matrix.os }})`
- `ci.yml:258`：`    needs: candidate-probe`
- `ci.yml:259`：`    if: github.event_name == 'workflow_dispatch' && github.ref == 'refs/heads/main'`
- `ci.yml:262-263`：`      matrix:` / `        os: [ubuntu-latest, windows-latest]`
- `ci.yml:265-271`：env 映射 `CANDIDATE_COMMIT`/`VERIFIER_COMMIT`/`BOOTSTRAP_COMMIT`/`BOOTSTRAP_MANIFEST_SHA256`/`EXPECTED_ORACLE_SHA256`/`LICENSE_POLICY_SHA256`
- steps：`ci.yml:273` `- name: Validate repository trust pins`、`ci.yml:296` `- name: Checkout trusted verifier`、`ci.yml:303` `- name: Checkout fixed bootstrap seed`、`ci.yml:310` `- name: Checkout candidate metadata`、`ci.yml:317` `- name: Setup Go`、`ci.yml:323` `- name: Build trusted verifier and seed`、`ci.yml:341` `- name: Download immutable candidate binary`、`ci.yml:348` `- name: Download untrusted probe oracle`、`ci.yml:355` `- name: Verify evidence and generate release metadata`、`ci.yml:403` `- name: Upload verified CI artifact`

## 2. workflow_dispatch inputs

`ci.yml:8` 为 `    inputs:`，`ci.yml:9-40` 共定义 **8 个** input（用户提纲写 9 个；**第 9 个：未找到**）。全部 `required: true`、`type: string`。

| # | input key | description 原文 | required | type | 行号 |
|---|---|---|---|---|---|
| 1 | `candidate_commit` | `Full lowercase commit SHA to validate` | `true` | `string` | 9-12 |
| 2 | `verifier_commit` | `Reviewed main commit containing this workflow and verifier` | `true` | `string` | 13-16 |
| 3 | `bootstrap_commit` | `Fixed lowercase 40-character bootstrap commit` | `true` | `string` | 17-20 |
| 4 | `bootstrap_manifest_sha256` | `SHA-256 of testdata/selfhost/bootstrap.json` | `true` | `string` | 21-24 |
| 5 | `expected_oracle_sha256` | `SHA-256 of testdata/selfhost/expected.json` | `true` | `string` | 25-28 |
| 6 | `license_policy_sha256` | `SHA-256 of testdata/selfhost/licenses-policy.json` | `true` | `string` | 29-32 |
| 7 | `noop_chain_sha256` | `SHA-256 of testdata/selfhost/noop-chain.json` | `true` | `string` | 33-36 |
| 8 | `executable_chain_sha256` | `SHA-256 of testdata/selfhost/executable-chain.json` | `true` | `string` | 37-40 |

## 3. inputs 的取值来源

**candidate_commit / verifier_commit 的关系约束（candidate-probe 与 selfhost 各有一份相同校验）**
- `ci.yml:171-176`（candidate-probe，位于 `ci.yml:161` 的 `Validate repository trust pins` 步骤内）：
  ```powershell
  if ($env:VERIFIER_COMMIT -ne $env:GITHUB_SHA) {
    throw 'verifier commit must equal the reviewed workflow commit on main'
  }
  if ($env:CANDIDATE_COMMIT -eq $env:VERIFIER_COMMIT) {
    throw 'candidate and verifier commits must be different'
  }
  ```
- `ci.yml:283-288`（selfhost，同样文本）。另有 40-hex 格式校验：`ci.yml:167` `if ($value -cnotmatch '^[0-9a-f]{40}$') {` → `ci.yml:168` `throw "$_ must be a full lowercase commit SHA"`（selfhost 对应 `ci.yml:279-280`，candidate-build 单独校验 `ci.yml:110-112`）。
- Go 侧同款模式：`internal/release/verify.go:25` `commitPattern      = regexp.MustCompile(`+"`"+`^[0-9a-f]{40}$`+"`"+`)`。
- input 语义描述：`ci.yml:10` `Full lowercase commit SHA to validate`、`ci.yml:14` `Reviewed main commit containing this workflow and verifier`。除 CI 校验外，仓库内未找到这两个值的枚举清单文件。

**bootstrap_commit ↔ bootstrap.json.seedCommit**
- `ci.yml:376-379`：
  ```powershell
  $bootstrap = Get-Content -Raw verifier-source/testdata/selfhost/bootstrap.json | ConvertFrom-Json
  if ($bootstrap.seedCommit -ne $env:BOOTSTRAP_COMMIT) {
    throw 'bootstrap manifest and CI seed commit differ'
  }
  ```
- `testdata/selfhost/bootstrap.json:3`：`  "seedCommit": "fc99f553c3066b24230a093d14bc17af0c6a398e"`（同一文件 `:4-7` 另有 `establishedBy`/`candidatePolicy`/`oracle`/`rollback` 字段）。
- bootstrap_commit 还用于 seed 检出与 seed 二进制版本注入：`ci.yml:306` `          ref: ${{ env.BOOTSTRAP_COMMIT }}`、`ci.yml:338`（见 §4）。

**5 个 sha256 input → 文件映射（yml description + Assert-Sha256 调用处）**
- `Assert-Sha256` 函数定义（两处重复定义）：`ci.yml:232-237` 与 `ci.yml:367-372`，其中 `ci.yml:233` `            $actual = (Get-FileHash -LiteralPath $Path -Algorithm SHA256).Hash.ToLowerInvariant()`，`ci.yml:235` `              throw "trusted fixture hash mismatch: $Path"`
- `noop_chain_sha256` → `verifier-source/testdata/selfhost/noop-chain.json`：`ci.yml:238` `          Assert-Sha256 verifier-source/testdata/selfhost/noop-chain.json $env:NOOP_CHAIN_SHA256`
- `executable_chain_sha256` → `verifier-source/testdata/selfhost/executable-chain.json`：`ci.yml:239` `          Assert-Sha256 verifier-source/testdata/selfhost/executable-chain.json $env:EXECUTABLE_CHAIN_SHA256`
- `bootstrap_manifest_sha256` → `verifier-source/testdata/selfhost/bootstrap.json`：`ci.yml:373` `          Assert-Sha256 verifier-source/testdata/selfhost/bootstrap.json $env:BOOTSTRAP_MANIFEST_SHA256`
- `expected_oracle_sha256` → `verifier-source/testdata/selfhost/expected.json`：`ci.yml:374` `          Assert-Sha256 verifier-source/testdata/selfhost/expected.json $env:EXPECTED_ORACLE_SHA256`
- `license_policy_sha256` → `verifier-source/testdata/selfhost/licenses-policy.json`：`ci.yml:375` `          Assert-Sha256 verifier-source/testdata/selfhost/licenses-policy.json $env:LICENSE_POLICY_SHA256`
- 注意：5 处校验全部针对 `verifier-source/` 检出副本（`ci.yml:190` `          ref: ${{ env.VERIFIER_COMMIT }}`、`ci.yml:191` `          path: verifier-source`）。64-hex 格式校验：`ci.yml:182-183` / `ci.yml:291-292` `throw "T017_$_ must be a lowercase SHA-256 digest"`。

**各 fixture 本体内容（供核对）**
- `testdata/selfhost/expected.json:3-8`：四个 case `{ "id": "version"/"validate-noop"/"run-noop"/"reject-executable", "outcome": "passed" }`
- `testdata/selfhost/licenses-policy.json:2-12`：`modules` 数组含 `github.com/larsonzh/prfrail` → `MIT`、`golang.org/x/text` → `BSD-3-Clause`
- `testdata/selfhost/noop-chain.json:4`：`    "id": "selfhost-noop"`；`testdata/selfhost/executable-chain.json:4`：`    "id": "selfhost-reject"`（其 step `:14` `          "kind": "build"`）

## 4. 构建与打包路径

| 动作 | 原文 | file:line |
|---|---|---|
| candidate 二进制构建（版本注入） | `go build -trimpath -ldflags "-X main.version=0.1.0+$env:CANDIDATE_COMMIT" -o (Join-Path $releaseRoot "prfrail$suffix") ./cmd/prfrail` | ci.yml:136 |
| candidate artifact 上传 | `name: candidate-binary-${{ runner.os }}-${{ inputs.candidate_commit }}` / `path: generations/candidate/release/prfrail*` | ci.yml:142-143 |
| probe harness 构建 | `go build -trimpath -o (Join-Path $generations "prfrail-release$suffix") ./cmd/prfrail-release` | ci.yml:209 |
| probe 调用 | `& $releaseTool probe ` （后续 6 行 flag） | ci.yml:241-247 |
| probe oracle 上传 | `name: candidate-oracle-${{ runner.os }}-${{ inputs.candidate_commit }}` / `path: generations/candidate/actual.json` | ci.yml:252-253 |
| verifier release 工具构建 | `go build -trimpath -o $releaseTool ./cmd/prfrail-release` | ci.yml:335 |
| seed 二进制构建（版本注入） | `go build -trimpath -ldflags "-X main.version=bootstrap-$env:BOOTSTRAP_COMMIT" -o $seedBinary ./cmd/prfrail` | ci.yml:338 |
| seed 哈希计算 | `$seedHash = (& $releaseTool hash-tree --root $seedRoot).Trim()` | ci.yml:381 |
| selfhost 证据校验调用 | `& $releaseTool selfhost ` + `--commit $env:CANDIDATE_COMMIT` / `--seed $seedRoot` / `--seed-hash $seedHash` / `--candidate $candidateRoot` / `--expected verifier-source/testdata/selfhost/expected.json` / `--actual (Join-Path $candidateRoot 'actual.json')` | ci.yml:382-388 |
| generatedAt 来源 | `$generatedAt = (git -C candidate-source show -s --format=%cI $env:CANDIDATE_COMMIT).Trim()` | ci.yml:390 |
| generate 调用 | `& $releaseTool generate ` + `--commit $env:CANDIDATE_COMMIT` / `--root $releaseRoot` / `--binary "prfrail$suffix"` / `--license-policy verifier-source/testdata/selfhost/licenses-policy.json` / `--generated-at $generatedAt` | ci.yml:391-396 |
| verify 调用 | `& $releaseTool verify --commit $env:CANDIDATE_COMMIT --root $releaseRoot` | ci.yml:397 |
| 打包目录 | `$packageRoot = Join-Path $candidateRoot 'artifact/prfrail'` → `New-Item ... -Force | Out-Null` → `Copy-Item -Path (Join-Path $releaseRoot '*') -Destination $packageRoot` | ci.yml:399-401 |
| 最终 artifact 上传 | `name: prfrail-${{ runner.os }}-${{ inputs.candidate_commit }}` / `path: generations/candidate/artifact/*` / `if-no-files-found: error` | ci.yml:406-408 |

## 5. `cmd/prfrail-release` 的功能面

用法行：`cmd/prfrail-release/main.go:23` `		return fmt.Errorf("usage: prfrail-release <hash-tree|probe|selfhost|generate|verify>")`

| 子命令 | flags（定义行） | 调用的 internal/release 函数 |
|---|---|---|
| `hash-tree` | `root`（main.go:28） | `release.HashTree`（main.go:32） |
| `probe` | `binary`/`noop-chain`/`executable-chain`/`work-root`/`output`/`expected-version`（main.go:40-45） | `release.ProbeCandidate`（main.go:49） |
| `selfhost` | `commit`/`seed`/`seed-hash`/`candidate`/`expected`/`actual`（main.go:52-57） | `release.VerifySelfHost`（main.go:61） |
| `generate` | `commit`/`root`/`binary`/`license-policy`/`generated-at`（main.go:65-69）；`main.go:73` `		timestamp, err := time.Parse(time.RFC3339, *generatedAt)` | `release.GenerateBundleMetadata`（main.go:77） |
| `verify` | `commit`/`root`（main.go:80-81） | `release.VerifyReleaseBundle`，且文件名硬编码在 main.go:85：`		result, err := release.VerifyReleaseBundle(release.BundlePlan{SourceCommit: *commit, ArtifactRoot: *root, ChecksumManifest: "SHA256SUMS", SBOMPath: "sbom.spdx.json", LicenseManifestPath: "licenses.json"})` |

## 6. `internal/release` 的产物集

**`GenerateBundleMetadata`**（`internal/release/generate.go:58` `func GenerateBundleMetadata(plan GeneratePlan) error {`）写出 3 个文件：
- `generate.go:116`：`	if err := writeJSON(filepath.Join(root, "licenses.json"), licenseOutput); err != nil {`
- `generate.go:143`：`	if err := writeJSON(filepath.Join(root, "sbom.spdx.json"), sbom); err != nil {`
- `generate.go:146`：`	return writeChecksumManifest(root, "SHA256SUMS")`
- 前置校验：`generate.go:82-83` `	if policy.SchemaVersion != "1.0.0" || len(policy.Modules) == 0 {` / `		return fmt.Errorf("%w: incomplete license policy", ErrInvalidBundle)`；`generate.go:110` `			return fmt.Errorf("%w: no approved license for module %q", ErrInvalidBundle, module.Path)`

**`VerifyReleaseBundle`**（`internal/release/verify.go:140` `func VerifyReleaseBundle(plan BundlePlan) (BundleEvidence, error) {`）校验：
- `verify.go:148`：`	manifestPath, err := resolveBundlePath(root, plan.ChecksumManifest)`（manifest 名由调用方传入 `"SHA256SUMS"`，见 §5）
- `verify.go:161`：`		return BundleEvidence{}, fmt.Errorf("%w: checksum manifest does not cover every artifact", ErrInvalidBundle)`
- `verify.go:174`：`			return BundleEvidence{}, fmt.Errorf("%w: checksum mismatch for %q", ErrInvalidBundle, path)`
- `verify.go:194-195`：`	if !equalLicenseModules(sbomModules, licenseModules) {` / `		return BundleEvidence{}, fmt.Errorf("%w: SBOM and license manifest modules differ", ErrInvalidBundle)`
- SBOM 与源提交绑定：`verify.go:366`（`validateSBOM` 起于 verify.go:354）：`	if document.Name != "prfrail-"+sourceCommit || document.DocumentNamespace != "https://proofrail.dev/spdx/"+sourceCommit {`
- license manifest 校验：`verify.go:392` `func validateLicenseManifest(path string) (string, []licenseModule, error) {`
- 返回证据：`verify.go:202` `		PublisherIdentityAuthenticated: false,`
- 辅助：`verify.go:289` `func readChecksums(path string) (map[string]string, string, error) {`（`:310` `		return nil, "", fmt.Errorf("%w: empty SHA256SUMS", ErrInvalidBundle)`）；`verify.go:319` `func regularBundleFiles(root, checksumManifest string) ([]string, error) {`
- 另一入口 `VerifySelfHost`（verify.go:82）：seed 哈希 `verify.go:114-115` `	if seedHash != plan.ExpectedSeedHash {` / `		return SelfHostEvidence{}, fmt.Errorf("%w: seed changed", ErrInvalidSelfHost)`；oracle 比对 `verify.go:125-126` `	if !equalOracle(expected, actual) {` / `		return SelfHostEvidence{}, fmt.Errorf("%w: candidate result differs from external oracle", ErrInvalidSelfHost)`
- `HashTree`（verify.go:206 `func HashTree(root string) (string, error) {`）：`verify.go:207` `	inventory, err := evidence.BuildReleaseInventory(root, time.Time{})`

## 7. `tools/` 下可复用的工具（仅列已摘到用途出处的）

| 路径 | 用途（摘录出处） |
|---|---|
| `tools/gates/gate.js` | `DELIVERY_DIRECTIVE section 6.3 gate runner - single entry point.`（gate.js:4；用法 `node tools/gates/gate.js --all [--skip-base] [--scope=tree|index|range:<A>..<B>|ci]` 见 gate.js:5-9） |
| `tools/contracts/`（README.md） | `# ProofRail contract oracle`（README.md:1）；`This directory is an independent S0 validation tool.`（README.md:3）；末段注明 `All repository JSON consumed here must be UTF-8 without BOM and LF-only.`（README.md 末节） |
| `tools/validation/c1-linux/run-c1-faces.sh` | `# C1 - Linux AgentRunner native validation.`（run-c1-faces.sh:3；`# It is a validation harness only: it never edits a source file, never commits, # never pushes and never touches the network (GOPROXY=off + GOTOOLCHAIN=local).` 见 :9-11） |
| `tools/agent-probe/Invoke-AgentProbe.ps1` | `.SYNOPSIS` `Controlled T026 CLI Agent capability probe runner (single invocation, fail-closed).`（Invoke-AgentProbe.ps1:2-3） |
| `tools/agent-probe/Invoke-CandidateProbe.ps1` | `.SYNOPSIS` `Scenario-level live probe runner for the T027 AgentRunner candidate (GitHub Copilot CLI).`（Invoke-CandidateProbe.ps1:3-4） |
| `tools/agent-stub/main.go` | `// Command agent-stub is the deterministic offline CLI used by the ProofRail // AgentRunner experiments.`（main.go:1-2） |
| `tools/b4-e2e/run-faces.ps1` | `.SYNOPSIS` `B4 (AT-23 mechanism leg) E2E driver: builds the three artifacts, prepares tmp/b4, ... aggregates tmp/b4/pack-verdict.json and writes tmp/b4/SHA256SUMS.txt.`（run-faces.ps1:2-5） |

存在但本次未摘到文件头/README 用途的目录：`tools/agent-probe/` 其余文件（CopilotCli.ps1、supervisor/、enforcement-proxy/ 等）、`tools/b3a-rig/`、`tools/b3b-rig/`、`tools/b4-e2e/` 其余 Go 文件、`tools/gates/lib/`、`tools/gates/selftest.js` —— 用途**未摘录**（未读取对应头注释）。编码相关自检散见于 b3 系列脚本头注释，如 `tools/b3a-rig/Test-B3aRig.ps1:19` `# PowerShell 5.1 compatible. Encoding: UTF-8 with BOM + LF. All output English.`

## 8. 事实性观察（仅陈述文件内容）

1. `candidate-build` / `candidate-probe` / `selfhost` 三个 job 的 `if:` 均为 `github.event_name == 'workflow_dispatch' && github.ref == 'refs/heads/main'`（ci.yml:98、149、259）；`test` job 无 `if:`（push/pull_request/workflow_dispatch 均触发，见 ci.yml:3-7）。
2. 依赖链：`candidate-probe` `needs: candidate-build`（ci.yml:148）；`selfhost` `needs: candidate-probe`（ci.yml:258）。
3. 全部 4 个 job 均为 `os: [ubuntu-latest, windows-latest]` 双矩阵（ci.yml:56、102、153、263）。
4. workflow_dispatch inputs 共 8 个，全部 `required: true`、`type: string`；提纲所述「9 个 input」中第 9 个未找到。
5. `verifier_commit` 与 `candidate_commit` 的约束在 candidate-probe 与 selfhost 各重复实现一次且文本一致：verifier 必须等于 `GITHUB_SHA`、candidate 必须不等于 verifier（ci.yml:171-176、283-288）。
6. 三个发布相关 job 的 commit/SHA-256 格式校验以 pwsh `-cnotmatch` 正则实现（ci.yml:110、167、182、279、291）；Go 侧 `internal/release/verify.go:25-26` 有同样的 `^[0-9a-f]{40}$` 与 `^([0-9a-f]{64})  ([^\\]+)$` 模式。
7. 5 个 sha256 input 只被用来校验 `verifier-source/testdata/selfhost/` 下检出的 fixture 副本（ci.yml:238-239、373-375）；当前 `bootstrap.json` 的 `seedCommit` 值为 `fc99f553c3066b24230a093d14bc17af0c6a398e`（testdata/selfhost/bootstrap.json:3）。
8. candidate 二进制经 artifact `candidate-binary-<os>-<candidate_commit>` 上传并以 `digest-mismatch: error` 下载两次（ci.yml:142-144、211-216、341-346）；probe oracle 经 `candidate-oracle-<os>-<candidate_commit>` 传递（ci.yml:252-254、348-353）。
9. 版本注入有两处不同前缀：candidate 为 `0.1.0+<candidate_commit>`（ci.yml:136，probe 校验同一值 `--expected-version "0.1.0+$env:CANDIDATE_COMMIT"` 见 ci.yml:247），seed 为 `bootstrap-<bootstrap_commit>`（ci.yml:338）。
10. bundle 文件名固定为 `licenses.json` / `sbom.spdx.json` / `SHA256SUMS`（generate.go:116、143、146；main.go:85 硬编码传入 verify）；`BundleEvidence.PublisherIdentityAuthenticated` 恒为 `false`（verify.go:202）。
11. selfhost 最终上传路径为 `generations/candidate/artifact/*`，artifact 名 `prfrail-<os>-<candidate_commit>`（ci.yml:406-408）；打包内容为 `generations/candidate/release/` 的整体拷贝（ci.yml:399-401）。
12. `expected.json` 含 4 个 case，与 `probe.go:57-61` 中的四个 probe case（version / validate-noop / run-noop / reject-executable）同名（probe.go:57-61、testdata/selfhost/expected.json:3-8）。
13. CI 的 `test` job 中 `Race` 与 `Contract fixtures` 两步仅 Linux 执行（ci.yml:83、87）。
