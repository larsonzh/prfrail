# P78 RELEASE-ZIP-BINDING 方案书（r1）

> 锚点：HEAD = origin/main = `c93571352f456a275d6e71ce3e084c9f96360e86`（准则 v1.69；SELFTEST 737/737；三作用域 TOTAL_FAIL=0）。用户 2026-11-17 起跑消息同轮授权：v0.1.0 Release 创建、CI `workflow_dispatch` 触发、主提交 push、§11.2 回填提交 push——**push 仅 origin、合计 2 次**，不做任何 gitee 动作。
> 本方案书中的记号：`CANDIDATE` = `c93571352f456a275d6e71ce3e084c9f96360e86`；`<RUN-ID>` = selfhost 所在 workflow run 的 run id（②b 拿到后回填，见 Q2.2/Q3.3）；`<VERIFIER>` = 本片主提交（仅文档）推送后 main HEAD 的完整 SHA（②a 提交后用 `git rev-parse HEAD` 取得）。

## Q1 构建链路

### Q1.1 `workflow_dispatch` 八个 input 的取值与依据

`ci.yml` 的 `workflow_dispatch.inputs` 实读共 **8 个**（全部 required、type string）。起跑消息写「9 个」系引用误差，如实披露：实为 candidate/verifier/bootstrap 三个 commit 类 input + 五个 sha256 类 input；本方案按 8 个执行，`ci.yml` 零改动。

| # | input | 取值 | 依据 |
|---|---|---|---|
| 1 | `candidate_commit` | `c93571352f456a275d6e71ce3e084c9f96360e86` | 步骤 0 锚点 HEAD；产品二进制源码 commit；`ci.yml:136` 据此 checkout 并构建 `./cmd/prfrail` |
| 2 | `verifier_commit` | `<VERIFIER>`（主提交完整 SHA） | `ci.yml:171-176`、`:283-288`：必须 == `GITHUB_SHA` 且 ≠ candidate；dispatch 前主提交已 push ⇒ 取 `git rev-parse HEAD` |
| 3 | `bootstrap_commit` | `fc99f553c3066b24230a093d14bc17af0c6a398e` | `ci.yml:376-379`：必须 == `testdata/selfhost/bootstrap.json` 的 `seedCommit`（固定值） |
| 4 | `bootstrap_manifest_sha256` | `Get-FileHash testdata/selfhost/bootstrap.json` 小写 64 hex | `ci.yml:373-375`：selfhost job 与 `verifier-source/testdata/selfhost/bootstrap.json` 逐字节比对 |
| 5 | `expected_oracle_sha256` | `testdata/selfhost/expected.json` 的 SHA-256 | 同上（`ci.yml:373-375`） |
| 6 | `license_policy_sha256` | `testdata/selfhost/licenses-policy.json` 的 SHA-256 | 同上（`ci.yml:373-375`） |
| 7 | `noop_chain_sha256` | `testdata/selfhost/noop-chain.json` 的 SHA-256 | `ci.yml:238-239`：candidate-probe job 校验 |
| 8 | `executable_chain_sha256` | `testdata/selfhost/executable-chain.json` 的 SHA-256 | 同上（`ci.yml:238-239`） |

五个 sha256 在 ②a 于本机用 `Get-FileHash -Algorithm SHA256`（`.Hash.ToLowerInvariant()`）计算。合法性与稳定性：主提交仅增改 `.md` 文档，`testdata/selfhost/` 与 `ci.yml` 在 candidate 与 verifier 两处树中逐字节相同 ⇒ 从当前工作树计算的五个值对两个 commit 同时有效。dispatch 命令（八个 `-f`，见 Q3.3）在报告登记原始输出。

### Q1.2 关键约束求解

约束：① `verifier_commit` == dispatch 时 main HEAD（`GITHUB_SHA`）；② `candidate_commit` ≠ `verifier_commit`。

**选定组合**：candidate = `CANDIDATE`（`c93571352f456a275d6e71ce3e084c9f96360e86`）；verifier = 本片主提交（仅文档，完整 SHA 待 ②a 提交后取得）。

合法性：主提交是新提交 ⇒ 与 `CANDIDATE` 不同 ⇒ ②成立；主提交 push 后 dispatch ⇒ main HEAD = `<VERIFIER>` ⇒ ①成立。

语义论证：
- 二进制源码 = `CANDIDATE`：产品代码自 `c93571352` 起零改动（本片不改 `.go`/`go.mod`），candidate-build 从 `candidate-source`（`CANDIDATE`）构建 `./cmd/prfrail`，`-ldflags "-X main.version=0.1.0+<CANDIDATE>"` ⇒ 二进制版本串 = `0.1.0+c93571352f456a275d6e71ce3e084c9f96360e86`，与 `INSTALLATION.md:144` 已冻结断言 `version.data.version -eq "0.1.0+$CandidateCommit"` 一致。
- verifier = 主提交（仅文档）：其 `verifier-source` 树（probe harness 源码 + `testdata/selfhost/` fixtures）与 `CANDIDATE` 逐字节相同 ⇒ 由 verifier 构建的可信 probe/selfhost 工具链对 candidate 的校验不因文档提交而改变；candidate-probe 用 noop/executable 两条冻结链实测候选二进制并比对 oracle（`ci.yml:247` 期望版本 = `0.1.0+<CANDIDATE>`）。
- 该组合正是 `ci.yml` 设计的用法（`verifier_commit` 的 input 描述原文：`Reviewed main commit containing this workflow and verifier`）。

### Q1.3 路径判定：CI `candidate-build` + `selfhost`（本机同参数构建仅作等价性对照）

**选定 CI 路径**。理由：
1. ADR-006 发布信任冻结为「固定 commit SHA + GitHub CI + SHA256SUMS + SBOM + 许可证清单」——GitHub CI 是**必需发布证据**，不可替代。
2. 只有 `selfhost` job 产生完整的发布证据链：seed 树哈希 → candidate probe oracle 比对（`actual.json` vs `expected.json`）→ `prfrail-release generate`（产出 `licenses.json` / `sbom.spdx.json` / `SHA256SUMS`）→ `prfrail-release verify` 校验 → 整目录打包为 artifact `prfrail-Windows-<CANDIDATE>`（`ci.yml:406-408`）。本机无法复刻 verifier/candidate 双身份分离、外部 oracle 与 artifact 不可变上传。
3. 绑定四元组需要 CI run id；本机构建无 run，无法满足 `INSTALLATION.md:206` 的 run 绑定。

**本机 `go build -trimpath -ldflags "-X main.version=0.1.0+<CANDIDATE>"` 仅为等价性对照**：在 `tmp/P78-RELEASE-ZIP-BINDING/10-build/` 构建，断言 `version --json` == `0.1.0+<CANDIDATE>` 且 `go build ./...`/`go vet ./...`/`go test ./...` 全绿；**不**断言与 CI 产物逐字节相同（工具链/环境差异允许），**不**作为发布证据，只作本地行为对照并登记进报告。

### Q1.4 `cmd/prfrail-release` 的角色：不自建打包脚本，直接复用 CI 产物

本片**不**自建打包脚本。`selfhost` job 已内置 `prfrail-release generate` / `verify` 完成四个发行文件的生产与校验（`internal/release/generate.go:116/143/146` 写 `licenses.json` / `sbom.spdx.json` / `SHA256SUMS`；`cmd/prfrail-release/main.go:85` 硬编码三文件名校验）。本片动作 = `gh run download <RUN-ID> --name prfrail-Windows-<CANDIDATE>` 取得 artifact 目录（内容 = `generations/candidate/release/*` 全量，恰四文件），随后组装 ZIP 并复验。理由：不自建可避免与已冻结、已受 CI 校验的生成链产生第二套证据源；且本片边界禁止改任何产品代码与 `ci.yml`。

### Q1.5 最终 ZIP 组装（Windows）

解压根目录只含 `prfrail/`、其内恰四文件；命名 `prfrail-Windows-<CANDIDATE>.zip`。`Compress-Archive` 是二进制打包，**不涉及 BOM/文本编码**；内部四个文件由 `Copy-Item` 字节复制，字节不变。

```powershell
$root   = 'tmp\P78-RELEASE-ZIP-BINDING\10-build'
$commit = 'c93571352f456a275d6e71ce3e084c9f96360e86'
# ① 下载 artifact（gh run download 自动解压到 --dir）
gh run download <RUN-ID> --repo larsonzh/prfrail --name "prfrail-Windows-$commit" --dir "$root\artifact"
# ② 组装前断言：$root\artifact\prfrail 恰含 prfrail.exe / SHA256SUMS / sbom.spdx.json / licenses.json
# ③ 组装
New-Item -ItemType Directory -Force "$root\zip-root\prfrail" | Out-Null
Copy-Item "$root\artifact\prfrail\*" "$root\zip-root\prfrail"
Compress-Archive -LiteralPath "$root\zip-root\prfrail" -DestinationPath "$root\prfrail-Windows-$commit.zip"
# ④ 结构断言（解到新目录后：根恰一个 prfrail\ 目录，恰四文件）
Expand-Archive -LiteralPath "$root\prfrail-Windows-$commit.zip" -DestinationPath "$root\zip-check"
```

字节关系：Release ZIP 与 Actions artifact ZIP（REST/浏览器下载得到的同名 zip）**不要求逐字节相同**——两者压缩元数据/打包实现可不同；要求的是**内容集相同**（恰四文件且字节一致）与 **SHA256SUMS 自洽**（对 ZIP 内每个文件逐行校验通过）。报告登记两个 ZIP 各自的 SHA-256 为观察值，不声称等价字节。

### Q1.6 构建产物与日志落盘

全部落 `tmp/P78-RELEASE-ZIP-BINDING/`（不入库、⑧c 清理）：`10-build/`（artifact 下载、ZIP 组装、Release notes 草稿 `release-notes.md`、本机等价对照二进制、构建/下载日志）、`20-reverify/`（Q4 四项工作目录与断言日志）。`tmp/` 根零文件。

## Q2 发行记录绑定

### Q2.1 四元组取值规则

| 元 | 取值 |
|---|---|
| 版本 | `v0.1.0`（Release tag 与产品版本口径；二进制内版本串 = `0.1.0+<CANDIDATE>`） |
| commit | candidate commit 全 40 位小写：`c93571352f456a275d6e71ce3e084c9f96360e86` |
| run | dispatch 后 `selfhost` 所在 workflow run 的 run id（数字；`gh run view <RUN-ID>` 取得，绑定该 run 下 candidate-build / candidate-probe / selfhost 双腿 job 结论全 success） |
| 文件名 | `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip` |

### Q2.2 三处落点的精确目标文本（CN + EN 成对）

`<RUN-ID>` 位置在第二次提交（回填提交）补齐，共 4 处（见 Q3.3）。commit 直接写完整 40 位，不用占位符。

**落点 1：`docs/validation/s1-exit{,_EN}.md:66`（§剩余风险与下一步第 1 条；非门表区；门表 9 行零改动）**

```
CN  file: docs/validation/s1-exit.md:66
1. 手工解压便携 ZIP、用户自选独立目录、显式路径运行且不修改 PATH 的模型已冻结；正式 ZIP 文件名和下载入口已由 P78 `RELEASE-ZIP-BINDING` 完成发行绑定（2026-11-17）：版本 `v0.1.0`、candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86`、CI run `<RUN-ID>`、文件名 `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`。EOL 通知与例外已由 ADR-010 获产品所有者批准（2026-11-16）。

EN  file: docs/validation/s1-exit_EN.md:66
1. The manually extracted portable-ZIP model, user-selected independent directory, explicit-path invocation, and no-PATH policy are frozen. The formal ZIP filename and download entry point completed release binding in P78 `RELEASE-ZIP-BINDING` (2026-11-17): version `v0.1.0`, candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86`, CI run `<RUN-ID>`, filename `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`. The EOL notice and exception policy was approved by the product owner via ADR-010 (2026-11-16).
```

**落点 2：`docs/S1_SUPPORT_MATRIX{,_EN}.md`（文件末尾新增一节，加法式；平台清单与支持层级表零改动；`:5` 不改，理由见 Q5）**

```
CN  file: docs/S1_SUPPORT_MATRIX.md（文末新增）
## 已绑定发行记录

v0.1.0 候选发行（P78 `RELEASE-ZIP-BINDING`，2026-11-17）已绑定：版本 `v0.1.0`；candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86`；CI run `<RUN-ID>`（candidate-build / candidate-probe / selfhost 双腿全绿）；文件名 `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`。发行信任边界仍以本页「发行信任边界」节为准；本绑定不改变平台清单与支持层级。

EN  file: docs/S1_SUPPORT_MATRIX_EN.md（文末新增）
## Bound Release Record

The v0.1.0 candidate release (P78 `RELEASE-ZIP-BINDING`, 2026-11-17) is bound: version `v0.1.0`; candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86`; CI run `<RUN-ID>` (candidate-build / candidate-probe / selfhost green on both legs); filename `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`. The release trust boundary still follows the "Release Trust Boundary" section of this page; this binding does not change the platform list or support tiers.
```

**落点 3：`docs/S1_RELEASE_NOTES{,_EN}.md`——`:5` 状态行替换 + 「发布前剩余门禁」段（L22–L25，磁盘实读；起跑消息引 :26-31 系行号漂移，同 8-vs-9 一并披露）+ 文末 P78 切片同步引语**

```
CN  file: docs/S1_RELEASE_NOTES.md:5
日期：2026-09-09；更新于 2026-11-17（P78 `RELEASE-ZIP-BINDING`）。状态：`CANDIDATE / RELEASED`。v0.1.0 候选发行已创建并绑定：版本 `v0.1.0`、candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86`、CI run `<RUN-ID>`、文件名 `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`；下载入口 = GitHub Release `v0.1.0`。支持起始日期与 S1 exit 整体批准仍待 S2 前冻结；本文件不授权 tag、GitHub Release、发布或安装。

EN  file: docs/S1_RELEASE_NOTES_EN.md:5
Date: 2026-09-09; updated 2026-11-17 (P78 `RELEASE-ZIP-BINDING`). Status: `CANDIDATE / RELEASED`. The v0.1.0 candidate release is created and bound: version `v0.1.0`, candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86`, CI run `<RUN-ID>`, filename `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`; download entry = GitHub Release `v0.1.0`. The support start date and whole S1 exit approval remain pending before S2; this file does not authorize a tag, GitHub Release, publication, or installation.
```

```
CN  file: docs/S1_RELEASE_NOTES.md（L22–L25 整段替换）
1. 对最终 ZIP 复跑核验、首次运行、升级、回滚和卸载。——已于 2026-11-17 由 P78 `RELEASE-ZIP-BINDING` 完成（证据：`validation/release-zip-binding.md`）。
2. 将版本、artifact 文件名、下载位置、EOL 通知方式和支持日期绑定到本说明及支持矩阵。——版本 `v0.1.0`、文件名 `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip` 与下载位置（GitHub Release `v0.1.0`）已绑定；EOL 通知方式已由 ADR-010 冻结（2026-11-16）；支持起始日期待 S1 exit 批准时冻结。
3. 获得独立评审结论和产品所有者明确 S1 exit/发行批准。——发行部分已由产品所有者批准（2026-11-17，P78 起跑消息同轮；ADR-046 Q10=(C) 第一步）；S1 exit 整体批准待 S2。
4. T027-B 已完成（部分不适用，ADR-045）；AT-23 阶段 1 条件闭合为永久事实（ADR-045）；S-G1 / S-G8 仍 blocked；S1 Exit 整体批准待 S2。

EN  file: docs/S1_RELEASE_NOTES_EN.md（same block）
1. Rerun verification, first run, upgrade, rollback, and uninstall against the final ZIP. — Done on 2026-11-17 by P78 `RELEASE-ZIP-BINDING` (evidence: `validation/release-zip-binding_EN.md`).
2. Bind version, artifact filenames, download location, EOL notification channel, and support date to these notes and the support matrix. — Version `v0.1.0`, filename `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`, and the download location (GitHub Release `v0.1.0`) are bound; the EOL notification channel was frozen by ADR-010 (2026-11-16); the support start date awaits S1 exit approval.
3. Obtain the independent-review verdict and explicit product-owner S1 exit/release approval. — The release part was approved by the product owner (2026-11-17, the P78 kickoff message in the same turn; ADR-046 Q10=(C) step one); whole S1 exit approval remains pending S2.
4. T027-B is complete (partly not applicable, ADR-045); the AT-23 stage-1 condition closure is a permanent fact (ADR-045); S-G1 / S-G8 stay blocked; whole S1 Exit approval remains pending S2.
```

```
CN  file: docs/S1_RELEASE_NOTES.md（P57 引语之后新增）
> 切片同步（P78 `RELEASE-ZIP-BINDING`）：GitHub Release `v0.1.0` 已创建并绑定发行记录（版本 `v0.1.0` / candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86` / CI run `<RUN-ID>` / 文件名 `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`）；复验四项通过；本段门禁清单已同步更新。

EN  file: docs/S1_RELEASE_NOTES_EN.md（after the P57 quote）
> Slice sync (P78 `RELEASE-ZIP-BINDING`): GitHub Release `v0.1.0` was created and its release record bound (version `v0.1.0` / candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86` / CI run `<RUN-ID>` / filename `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`); all four reverification items passed; the gate list above is updated accordingly.
```

**需在回填提交补齐 `<RUN-ID>` 的位置（共 4 处）**：① `s1-exit{,_EN}:66`；② `S1_SUPPORT_MATRIX{,_EN}` 新节；③ `S1_RELEASE_NOTES{,_EN}:5`；④ `S1_RELEASE_NOTES{,_EN}` 文末 P78 引语。其余目标文本不含 run 号，无回填。⑩ 前 grep `<RUN-ID>` 必须零命中。

### Q2.3 Release 附件策略

**结论：上传 4 个资产 = ZIP + `SHA256SUMS` + `sbom.spdx.json` + `licenses.json`，不上传独立 `prfrail.exe`。**
理由：① ZIP 是冻结安装模型的唯一正式下载载体（§2 解压根目录 = `prfrail/` 四文件），作为主资产；② `SHA256SUMS`/SBOM/许可证清单是 ADR-006 发布信任证据，须在 Release 页**无需解压即可见**，供供应链核验；③ 独立 `prfrail.exe` 会制造第二个二进制下载面，与 ZIP 内字节存在漂移风险，故排除（即「打散上传四个文件」不成立——四个文件中的二进制不打散上传）。三个元数据文件与 CI 产物逐字节一致（不自建、不重生成）。

### Q2.4 绑定落点结论

**是**：四元组同时落 ① 三处文档落点（Q2.2）；② 本片报告对 `docs/validation/release-zip-binding{,_EN}.md`（②e 生成，真实 run id 成文）；③ GitHub Release notes 正文（`--notes-file`，见 Q3.2）。三处文档落点是长期可核对权威，报告与 Release notes 是执行证据与对外入口。

## Q3 Release 创建

### Q3.1 `gh release create` 精确参数

```powershell
gh release create v0.1.0 `
  --repo larsonzh/prfrail `
  --target c93571352f456a275d6e71ce3e084c9f96360e86 `
  --title "v0.1.0 (P78 RELEASE-ZIP-BINDING)" `
  --notes-file "tmp\P78-RELEASE-ZIP-BINDING\10-build\release-notes.md" `
  "tmp\P78-RELEASE-ZIP-BINDING\10-build\prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip" `
  "tmp\P78-RELEASE-ZIP-BINDING\10-build\SHA256SUMS" `
  "tmp\P78-RELEASE-ZIP-BINDING\10-build\sbom.spdx.json" `
  "tmp\P78-RELEASE-ZIP-BINDING\10-build\licenses.json"
```

- **tag 名**：`v0.1.0`（ADR-010 冻结的首版候选版本号）。创建前断言 `git ls-remote --tags origin v0.1.0` 为空。
- **`--target` = candidate commit（`CANDIDATE`），不是主提交**。理由：Release 指向的源码必须与二进制源码逐字节同源；主提交只是为满足 verifier 身份约束而晚于 candidate 的纯文档提交，把 tag 指向它会制造「发行源码 ≠ 二进制源码」的第二身份，并污染四元组中 commit 的唯一性。
- **notes 文件来源**：`tmp/P78-RELEASE-ZIP-BINDING/10-build/release-notes.md`（②b 起草；此时 run id 已拿到，直接用真实值，无占位）。
- **附件清单**：见 Q2.3 四个资产。
- **`--draft` = 否、`--prerelease` = 否（直接发布，非预发布标记）**。理由：① v0.1.0 是 ADR-010 冻结的首版候选，Release 创建已获产品所有者同轮批准（Q10=(C) 第一步，`CHECKLIST_v61` 与 `REMAINING_SLICES` 均载「创建 v0.1.0 Release」）；② draft 状态会阻止公开下载，而绑定后的 `INSTALLATION.md` §2 下载入口必须可即时核验；③「候选状态、非已发布 SLA、S1 exit 整体待 S2」这些限定已由 `S1_RELEASE_NOTES:5` / `S1_SUPPORT_MATRIX:5` / `ADR_REGISTER:62` 如实承载，不必再借 GitHub 的 prerelease 标记重复表达（该标记语义更粗，且会与「正式 Release」的口径互相干扰）。

### Q3.2 Release notes 骨架（`release-notes.md`，双语；②b 以真实 run id 成文）

```
# ProofRail v0.1.0（候选发行 / Candidate Release）

> 状态：候选发行（CANDIDATE / RELEASED）。本 Release 由产品所有者批准创建（2026-11-17，P78 起跑消息同轮；ADR-046 Q10=(C) 第一步）。S1 exit 整体批准待 S2；非已发布 SLA。
> Status: candidate release; creation approved by the product owner (2026-11-17, P78 kickoff, same turn; ADR-046 Q10=(C) step one). Whole S1 exit approval remains pending S2; not a published SLA.

## 发行绑定 / Release Binding
- 版本 / version: v0.1.0
- candidate commit: c93571352f456a275d6e71ce3e084c9f96360e86
- CI run: <RUN-ID>（candidate-build / candidate-probe / selfhost 双腿全绿 / green on both legs）
- 文件名 / filename: prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip

## 能力摘要 / Capability Summary
（引用 S1_RELEASE_NOTES「候选能力」节，逐条转写）

## 安装与校验入口 / Install and Verify
- 下载 = 本 Release 资产；核验 = docs/INSTALLATION.md §2（SHA256SUMS 逐行 / SBOM / licenses / version --json == 0.1.0+c93571352f456a275d6e71ce3e084c9f96360e86）

## 发行信任边界 / Release Trust Boundary
- 完整 commit SHA + GitHub CI + SHA256SUMS + SPDX SBOM + 许可证清单；无签名/attestation；校验和不认证发布者身份。

## EOL 公告义务判定 / EOL Notice Obligation
- v0.1.0 是首版 minor，无前一 minor ⇒ 无 EOL 公告义务（ADR-010 义务仅在「新 minor 发布时公告前一 minor」触发）；如实记录。

## 支持与已知限制 / Support and Known Limitations
- 引用 S1_SUPPORT_MATRIX 平台清单与发行信任边界；SLA 未发布；S1 exit 整体批准待 S2。
```

### Q3.3 时序论证（≤2 次 push，run 号先拿到）

顺序：**① push 主提交**（Q5 全部文档目标文本，`<RUN-ID>` 标记 4 处）→ 观察该 push 触发 CI 双腿 green → **② dispatch**（`gh workflow run ci.yml -f candidate_commit=… -f verifier_commit=<VERIFIER> -f bootstrap_commit=… -f bootstrap_manifest_sha256=… -f expected_oracle_sha256=… -f license_policy_sha256=… -f noop_chain_sha256=… -f executable_chain_sha256=…`）→ **③ 等 selfhost 完成拿 run id**（`gh run watch <RUN-ID>`；记录四 job 双腿结论）→ **④ `gh run download` 取 artifact → 组装 ZIP（Q1.5）→ 复验四项（Q4）→ ⑤ 创建 Release（Q3.1，notes 用真实 run id）→ ⑥ 回填提交（第二次 push）**：把 4 处 `<RUN-ID>` 替换为真实 run id，并落报告对 `docs/validation/release-zip-binding{,_EN}.md`、证据包追加件与 `REMAINING_SLICES{,_EN}` 完成记录。

论证：
- run 号在 ③ 拿到 ⇒ 绑定句的最终成文（⑥）与 Release 创建（⑤）都发生在 run 号之后，满足「run 号先拿到再绑定与创建 Release」。
- push 计数 = 主提交 1 次 + 回填提交 1 次 = **2 次**，不超授权。Release 创建走 `gh` API（tag 由 API 建立，非 git push），不消耗 push 额度；Release notes 若有整改用 `gh release edit`（不占 push）。
- 主提交内容在 ⑤⑥⑦ 审查前已 push（verifier 约束使然）⇒ ⑤⑥⑦ 对落库 diff（主提交 + 回填提交 = 全片 diff）的整改**一律并入回填提交**；任何需要第 3 次 push 的情形 ⇒ 停下升级用户裁决。

## Q4 复验矩阵（全部对最终 ZIP，非候选 artifact 目录）

工作目录全部落 `tmp/P78-RELEASE-ZIP-BINDING/20-reverify/`；每项断言原始输出存档 `tmp/P78-RELEASE-ZIP-BINDING/10-build/`。任一项失败 ⇒ **停下升级**（不修码、不自愈、不换方案），如实报告实况。

| # | 项 | 状态路径 | 操作 | 成功断言 | 失败处置 |
|---|---|---|---|---|---|
| 1 | 安装 | `20-reverify/install/` | 解压最终 ZIP → 结构断言 → §2 复验块 → §3 五步快速上手（isolated workspace） | 解压根目录恰含 `prfrail/`、其内恰四文件；`SHA256SUMS` 逐行校验通过；`sbom.spdx.json` / `licenses.json` 可解析；`version --json` == `0.1.0+c93571352f456a275d6e71ce3e084c9f96360e86`；`validate` 通过；`preview` 零命令/网络/模型/凭据/版本探测调用；noop `run` 达 `COMPLETED` 且 `report` 可读 | 停下升级 |
| 2 | 升级 | `20-reverify/upgrade/` | dirA = 安装目录（先在 dirA 生成历史 run）；dirB = 新目录解压同一最终 ZIP；确认无写者 → dirB 完整路径只读 `version` / `validate` / `preview` → 后续命令显式改用 dirB 路径 → dirA 保留 | dirB 版本 == `0.1.0+<CANDIDATE>`；dirA 未被覆盖；dirA 历史 run 仍可 `report --run-dir` 读取 | 停下升级 |
| 3 | 回滚 | `20-reverify/rollback/` | 停止写者 → 路径切回 dirA → dirA `version` + `report --run-dir <历史 run>` 只读核验 | 历史 run 报告内容与升级前一致；dirA 数据未被 dirB 写入；未知格式保持只读 | 停下升级 |
| 4 | 卸载保留数据 | `20-reverify/uninstall/` | 确认无 prfrail 进程 → 仅删除 dirB 二进制目录 | dirB 已删；证据/配置/run/store 目录仍存在；仍可用 dirA 完整路径 `report --run-dir` 读取保留数据 | 停下升级 |

升级口径说明：v0.1.0 是首版 ⇒ **不存在旧正式版本可做真实跨版本升级**；本复验按 `INSTALLATION.md` §5 的切换机制口径执行（同一最终 ZIP 解压到第二个独立目录 = 步骤 2「新 ZIP 解压到新的 commit 目录、不得覆盖旧目录」；显式切换完整路径 = 步骤 4；旧目录保留为回滚点 = 步骤 4/5），并把「历史 run 仍可读取」作为升级后兼容断言；真实跨版本升级演练留待首个 patch/minor 发布。此口径逐条满足 §5 五步语义（停止写者与备份在演练中以「无运行写者 + 全量工作目录落在独立 tmp 路径」体现）。注：起跑消息引 `INSTALLATION.md:210-216` 升级段——磁盘实读该文件共 206 行，§5 升级与回滚五步位于 L190–L200（如 L195 = 步骤 2），§6 卸载 L204–L206；`:210-216` 系行号漂移，按 §5/§6 语义执行，如实披露。

## Q5 文档同步判定（三层归属）

| # | file:line | 归属 | 判定与理由 | 目标文本 |
|---|---|---|---|---|
| 1 | `INSTALLATION{,_EN}.md:5` | 第一层 | 「当前尚无正式 GitHub Release」现役假句，翻转 | 见下 ① |
| 2 | `INSTALLATION{,_EN}.md` §2（`:60` 段后） | 第一层 | 下载入口需改为 Release 资产（加法式一句；§2 脚本模板与 L144 `0.1.0+$CandidateCommit` 冻结断言不动） | 见下 ② |
| 3 | `INSTALLATION{,_EN}.md:206` | 第一层 | 绑定句翻转为已完成（含四元组） | 见下 ③ |
| 4 | `S1_RELEASE_NOTES{,_EN}.md:5` | 第一层 | 状态行替换为 `CANDIDATE / RELEASED` + 四元组 | Q2.2 |
| 5 | `S1_RELEASE_NOTES{,_EN}.md` L22–L25 | 第一层 | 门禁段按 P78 结果同步（非历史改写，是本文件的现役门禁清单） | Q2.2 |
| 6 | `S1_RELEASE_NOTES{,_EN}.md` 文末 | 第一层 | P78 切片同步引语 | Q2.2 |
| 7 | `S1_SUPPORT_MATRIX{,_EN}.md` 文末 | 第一层 | 加法式绑定节 | Q2.2 |
| 8 | `S1_SUPPORT_MATRIX{,_EN}.md:5` | **不触** | 「待正式发行载体和 S1 exit 批准后生效」仍为真（S1 exit 未批准）；「当前不是已发布 SLA」仍为真且与 `ADR_REGISTER:62` 不翻一致 | — |
| 9 | `S1_SUPPORT_MATRIX{,_EN}.md:7-15` | **不触** | 平台清单与支持层级表（硬边界） | — |
| 10 | `validation/s1-exit{,_EN}.md:66` | 第一层（绑定落点） | 剩余风险第 1 条加法式翻转；门表 9 行（L52–L60）零改动 | Q2.2 |
| 11 | `ADR_REGISTER{,_EN}.md:19` | 第二层 → **触** | ADR-009 状态列是活记录（同 ADR-010 由 P77 翻转的先例），闭括号 | 见下 ④ |
| 12 | `ADR_REGISTER{,_EN}.md:62` / `:63` | **不触（默认不翻）** | 「当前仍不是已发布 SLA」：候选发布 ≠ SLA 生效，该句待 S1 exit 通过或用户裁定再翻；本方案不主张翻，无需停下升级 | — |
| 13 | `ADR_REGISTER{,_EN}.md:104` / `:105` | **不触** | T001 决策单 = 2026-09-09 时点批准记录；「最终 ZIP 仍需绑定发行记录并复验；不授权发布」记录的是该时点授权范围（§11.1 历史不改写）；闭括号只落 :19 状态列 | — |
| 14 | `INSTALLATION_PLAN{,_EN}.md:32` / `:38` / `:50` | 第二层 → **触**（三处半句翻转） | 「正式下载入口仍待冻结 / 尚未发布 / 仍须绑定」三处现役假句 | 见下 ⑤ |
| 15 | `OPERATIONS{,_EN}.md:5` | 第二层 → **触** | 「当前没有 ProofRail 正式发行包」现役假句 | 见下 ⑥ |
| 16 | `DEV_PLAN{,_EN}.md:96` / `:97` | 第二层 → **触**（加法式注记，不改历史句） | 「最终 ZIP 绑定/复验、EOL 通知方式…仍缺失」现役假句；T018 勾选状态不变（S1 exit 批准仍缺失） | 见下 ⑦ |
| 17 | `PRODUCT_REQUIREMENTS{,_EN}.md` | **不触** | 三关键词零命中；无发布入口类待落实句；§7 独立安全评审指针由 ADR-047（P76）负责，不在本片绑定面 | — |
| 18 | `docs/t027/REMAINING_SLICES{,_EN}.md` P78 段 | 第一层 | 步骤 0c 台账段登记 + ⑧b 完成记录（只写报告路径 + R/D） | ⑧b 时成文 |
| 19 | `docs/validation/release-zip-binding{,_EN}.md` | 第一层 | 新建报告对（含 S-G3 收口章与 `sixstage` 围栏） | ②e |
| 20 | `docs/validation/evidence/P78-RELEASE-ZIP-BINDING/` | 第一层 | 证据包（②e/⑨ 追加 `01-arch/` / `02-batch/` / `04-reviews/` / `06-gates/` / `MANIFEST.md` / `SHA256SUMS.txt`） | ②e |

目标文本 ①–⑦：

```
① CN  INSTALLATION.md:5（句内替换末句）
正式 GitHub Release `v0.1.0` 已创建（2026-11-17，P78 `RELEASE-ZIP-BINDING`）；最终 ZIP 已按 §2 下载并复验通过，与最终发行物结构相同。

① EN  INSTALLATION_EN.md:5
The formal GitHub Release `v0.1.0` was created on 2026-11-17 (P78 `RELEASE-ZIP-BINDING`); the final ZIP was downloaded and reverified per section 2 and matches the final carrier layout.

② CN  INSTALLATION.md §2（:60 段后新增一段）
正式下载入口为 GitHub Release `v0.1.0` 的资产 `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`（`https://github.com/larsonzh/prfrail/releases/tag/v0.1.0`）；下方 Actions run 下载路径保留为 CI 证据下载路径。

② EN  INSTALLATION_EN.md §2（after the :60 paragraph）
The formal download entry is the GitHub Release `v0.1.0` asset `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip` (`https://github.com/larsonzh/prfrail/releases/tag/v0.1.0`); the Actions-run download path below remains the CI-evidence path.

③ CN  INSTALLATION.md:206（整句替换）
本页 §2 命令已于 2026-11-17 用最终 ZIP `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip` 复验通过，并在 [S1 exit 报告](validation/s1-exit.md)、[支持矩阵](S1_SUPPORT_MATRIX.md)和[发行说明](S1_RELEASE_NOTES.md)中绑定：版本 `v0.1.0`、candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86`、CI run `<RUN-ID>`、文件名 `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`。

③ EN  INSTALLATION_EN.md:206
The section 2 commands were rerun and verified on 2026-11-17 against the final ZIP `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`, and the [S1 exit report](validation/s1-exit_EN.md), [support matrix](S1_SUPPORT_MATRIX_EN.md), and [release notes](S1_RELEASE_NOTES_EN.md) now bind: version `v0.1.0`, candidate commit `c93571352f456a275d6e71ce3e084c9f96360e86`, CI run `<RUN-ID>`, filename `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip`.

④ CN  ADR_REGISTER.md:19（状态/完成条件列尾部替换）
2026-09-09 产品所有者明确批准；最终 ZIP 已绑定发行记录并复验（P78，2026-11-17）

④ EN  ADR_REGISTER_EN.md:19
Product owner explicitly approved 2026-09-09; the final ZIP was bound to the release record and rerun (P78, 2026-11-17)

⑤ CN  INSTALLATION_PLAN.md:32（句中替换）
state/store 默认位置仍待冻结；正式下载入口已绑定为 GitHub Release `v0.1.0` 资产；EOL 通知方式已由 ADR-010 获产品所有者批准（2026-11-16）。
⑤ CN  INSTALLATION_PLAN.md:38（句替换）
正式下载入口与最终 ZIP 已发布（GitHub Release `v0.1.0`），并已用最终载体重跑。
⑤ CN  INSTALLATION_PLAN.md:50（句尾替换）
正式发行文件名和下载地址已绑定：`prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip` 与 GitHub Release `v0.1.0`。

⑤ EN  INSTALLATION_PLAN_EN.md:32
The default state/store location remains open; the formal download entry point is bound to the GitHub Release `v0.1.0` asset; the EOL notification channel was approved by the product owner via ADR-010 (2026-11-16).
⑤ EN  INSTALLATION_PLAN_EN.md:38
The formal download entry point and final ZIP now exist (GitHub Release `v0.1.0`), and the drill was rerun against the final carrier.
⑤ EN  INSTALLATION_PLAN_EN.md:50
the formal release filename and download location are bound: `prfrail-Windows-c93571352f456a275d6e71ce3e084c9f96360e86.zip` on GitHub Release `v0.1.0`.

⑥ CN  OPERATIONS.md:5（首句替换）
正式发行 `v0.1.0` 已创建（2026-11-17，候选状态、非已发布 SLA）；手工解压 Windows amd64 便携 ZIP、显式路径运行且不修改 PATH 的 S1 安装模型已获批，安装步骤见 [安装指南](INSTALLATION.md)。

⑥ EN  OPERATIONS_EN.md:5
The formal release `v0.1.0` was created (2026-11-17, candidate status, not a published SLA). The S1 model of manually extracting a Windows amd64 portable ZIP, invoking it by explicit path, and not modifying PATH is approved; see the [installation guide](INSTALLATION_EN.md).

⑦ CN  DEV_PLAN.md:96（T018 状态句后追加）
**P78 注记（2026-11-17）**：最终 ZIP 绑定/复验与 EOL 通知方式已落实（P78 `RELEASE-ZIP-BINDING` / ADR-010）；T018 勾选状态不变（产品所有者 S1 exit 批准仍缺失）。

⑦ EN  DEV_PLAN_EN.md:97
**P78 note (2026-11-17)**: final-ZIP binding/rerun and the EOL notification channel are done (P78 `RELEASE-ZIP-BINDING` / ADR-010); the T018 checkbox stays unchanged (the product-owner S1 exit approval is still missing).
```

升降级披露：第二层初始分类中，ADR_REGISTER（:19 触、:62/:104 不触）、INSTALLATION_PLAN（触 3 处）、OPERATIONS（触 1 处）、DEV_PLAN（触 1 处加法式）由实读升级为「触」；PRODUCT_REQUIREMENTS 实读判不触；无第三层 → 第二层升级冲突边界的情形（升级目标均在【七】边界允许面内）。

## Q6 组合 grep 排查

三个独立 grep（本片 Windows 环境执行，作用域 `docs\**\*.md`）：

```powershell
Get-ChildItem docs -Recurse -Filter *.md | Select-String -Pattern 'RELEASE-ZIP'    -List | Select-Object Path
Get-ChildItem docs -Recurse -Filter *.md | Select-String -Pattern 'GitHub Release' -List | Select-Object Path
Get-ChildItem docs -Recurse -Filter *.md | Select-String -Pattern 'v0\.1\.0'       -List | Select-Object Path
```

**实测结论（本方案书成文时已实跑，命中清单按文件归类）**：

- `RELEASE-ZIP`：`docs/t027/CHECKLIST_v60{,_EN}`、`docs/t027/CHECKLIST_v61{,_EN}`、`docs/t027/REMAINING_SLICES{,_EN}`、`docs/validation/adr010-approval{,_EN}`、`docs/validation/evidence/P77-*/P78-*`（证据归档）。
- `GitHub Release`（prfrail 仓库内）：`docs/ADR_REGISTER{,_EN}`（:20 ADR-010 行、:62/:63 支持策略段）、`docs/INSTALLATION{,_EN}:5`、`docs/S1_RELEASE_NOTES{,_EN}:5`、`docs/S1_SUPPORT_MATRIX{,_EN}:21`、`docs/t027/CHECKLIST_v60{,_EN}` / `CHECKLIST_v61{,_EN}`、`docs/validation/s1-install-drill{,_EN}:5`、`docs/validation/s1-selfhost{,_EN}:100`、`docs/validation/t017-baseline-review{,_EN}:46`、证据归档（P72/P76/P77/P78）。另命中 `sessbridge\docs\SESSBRIDGE_README{,_EN}`——属另一仓库（sessbridge），不属 prfrail，不触。
- `v0\.1\.0`：`docs/ADR_REGISTER{,_EN}`（:20 / :57-58 ADR-046 行 / :62-63 支持策略段）、`docs/S1_SUPPORT_MATRIX{,_EN}:18`、`docs/t027/CHECKLIST_v60{,_EN}` / `CHECKLIST_v61{,_EN}`、`docs/t027/REMAINING_SLICES{,_EN}`、`docs/validation/adr010-approval{,_EN}`、`docs/validation/s1-exit-closeout-assessment{,_EN}:111`、证据归档。注：`INSTALLATION{,_EN}` §2 脚本含 `0.1.0+$CandidateCommit`（非 `v0.1.0` 前缀），为已冻结核验断言，保持不触。

**升降级规则**：① 命中文件属第三层冻结快照（`CHECKLIST_v60/v61`）、历史报告（`s1-install-drill` / `s1-selfhost` / `t017-baseline-review` / `adr010-approval` / `s1-exit-closeout-assessment`）或证据归档 ⇒ 不触，登记理由（历史不改写）；② 命中文件为现役指南/矩阵/说明且句子现役为假 ⇒ 升级为第二层触（已并入 Q5 表格）；③ 命中为第一层台账（`REMAINING_SLICES`）⇒ 按第一层执行；④ 拟触文件与【七】边界冲突 ⇒ §1.7.2 停下升级；⑤ sessbridge 命中 ⇒ 不属 prfrail，不触。本片 grep 结果**无新增升级项**——全部命中均已归入上述分类。

## 不触清单

- 产品代码：`.go`（含 `internal/release/**`、`cmd/prfrail-release/**`）、`go.mod`、`ci.yml`、`.github/**` 其余文件。
- 门表：`docs/validation/p57-closeout-assessment{,_EN}` §4.2、`docs/validation/s1-exit{,_EN}` 门表（L52–L60）。
- 准则与契约：`docs/DELIVERY_DIRECTIVE{,_EN}`、`docs/CONTRACTS{,_EN}`、`docs/SECURITY{,_EN}`、`docs/ARCHITECTURE{,_EN}`、`docs/TEST_STRATEGY{,_EN}`、`docs/BUSINESS_WORKFLOWS{,_EN}`、`docs/CODING_CONVENTIONS{,_EN}`、`docs/DOCUMENTATION_PLAN{,_EN}`、`docs/RFC-*.md`。
- 冻结快照与历史：`docs/t027/CHECKLIST_v60{,_EN}`、`CHECKLIST_v61{,_EN}`、v50–v61 快照、`docs/validation/` 下历史报告（`s1-install-drill` / `s1-selfhost` / `t017-baseline-review` / `adr010-approval` / `s1-exit-closeout-assessment` 等）与既有证据归档。
- ADR：除 `ADR_REGISTER{,_EN}:19` 状态列闭括号外，ADR-001~047 全部不触（含 ADR-010 行 :20 已由 P77 关闭、ADR-046/047 行、`:62/:63` SLA 句、`:104/:105` T001 决策单行）；不新增 ADR / OB。
- `S1_SUPPORT_MATRIX{,_EN}` 平台清单与支持层级表（:7–:15）及 :5。
- schema、fixtures、`PRODUCT_REQUIREMENTS{,_EN}`、sessbridge 仓库全部文件。
- 无代码签名；无 Linux 正式发布；无 gitee 任何动作。

## 风险与停下升级触发条件

1. **`<RUN-ID>` 标记的短暂存在**：主提交含 4 处 `<RUN-ID>` 标记，回填提交补齐。若在主提交后、回填前任何环节失败停片 ⇒ 仓库短暂保留「已绑定」句与标记，处置由用户裁决（回滚主提交或继续）；⑩ 前 grep 清零为硬条件。
2. **主提交先于 ⑤⑥⑦ 审查推送**（verifier 约束使然）：⑤⑥⑦ 对落库 diff 的整改并入回填提交；Release notes 侧整改用 `gh release edit --notes-file`（不占 push）。**任何第 3 次 push 需求 ⇒ 停下升级用户裁决**，不自行追加。
3. **构建/校验失败**：dispatch 任一 job 红（含五个 sha256 不匹配、probe 版本不匹配、`bootstrap_commit` 不符）⇒ 按 §6.4 取证、复跑一次、如实登记；仍红 ⇒ **停下升级**。不修 `ci.yml`、不修产品代码、不以本机构建顶替。
4. **复验四项任一断言失败** ⇒ 停下升级（Q4 矩阵「失败处置」列）。
5. **Release 创建失败**（tag 已存在、资产上传失败、`gh` 权限变化）⇒ 停下升级；不重复创建、不删除重建。
6. **`ADR_REGISTER:62` SLA 句**：默认不翻。若复验/绑定过程发现必须翻的情形 ⇒ 停下升级用户裁决，不自行改。
7. **门禁红线**：④ 三作用域 + SELFTEST、build/vet/test 任一红 ⇒ 停下；⑥ 必做、⑦ 额度 1+2 按【五】执行；G6-4 / E-1 / D=2 / OB-139 纪律全片适用。
8. **artifact 有效期**：Actions artifact 约 90 天过期，正式入口 = Release 资产（`INSTALLATION.md` §2 已声明「Actions artifact 可能过期，它不是长期正式下载入口」）。

DESIGN: DONE
