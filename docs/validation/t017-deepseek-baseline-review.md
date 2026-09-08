# T017 Verifier Baseline DeepSeek 独立主审存档

[English](t017-deepseek-baseline-review_EN.md)

日期：2026-09-09。结论：`PASS FOR BASELINE REVIEW`。未发现可复现的 `BLOCKER` 或 `HIGH`。仍待独立 verifier 基线落库、不同 candidate SHA，以及首次真实 Windows/Linux GitHub CI。

## 评审来源与边界

1. 本页归档用户从折叠聊天窗口手工恢复的 DeepSeek V4 Pro 主审报告。原报告中的两个表格按复制出的行列内容重建；其余内容保持原意。
2. 该主审采用与 [GPT-5.3 Codex verifier baseline 独立主审](t017-baseline-review.md)相同的严格只读要求，最终 verdict 同为 `PASS FOR BASELINE REVIEW`。
3. 两份一致结论相互印证受审 baseline 已通过独立模型复核，但都只是会话评审记录，不是签名 receipt、attestation 或可独立认证的模型身份记录。
4. 文中代码位置是受审快照的定位线索；后续文件变化可能使行号漂移。

## Findings

**未发现 BLOCKER/HIGH。**

1. `MEDIUM`（已知、已披露的信任边界）：脚本摘要门禁是内容漂移检测，不是 PowerShell 语义证明。若攻击者同时修改 workflow `run` 与测试中的期望摘要，本地测试可重新变绿。

	- 证据：`workflow_test.go:234`、`workflow_test.go:241`、`s1-selfhost.md`“本轮整改边界”第 4 条、`t017-baseline-review.md` Findings 第 1 条。
	- 现有测试为何未阻止：这是门禁的固有结构属性，无法用同源测试自证。
	- 缓解（已落实并归档）：摘要来自独立常量、禁止从待验 workflow 自动刷新、baseline 必须独立评审落库。维持不阻断。

2. `LOW`：candidate 的 `actual` 结果是候选进程的自我报告。`candidate-probe` 必须执行候选二进制，恶意候选可伪造符合格式的 JSON 输出（甚至在 runner 上自行改写 `actual.json`），与外部 oracle 的相等性检查只检测**漂移**，不是对恶意候选的安全证明。

	- 证据：`ci.yml#L136`（probe 执行候选）、`probe.go:47`、`verify.go:82`（仅比较相等）。
	- 现有测试为何未阻止：black-box 探针无法消除自证；`fake-oracle-write` 变异只覆盖脚本层，不覆盖候选进程自身行为。
	- 建议（后续增强，不阻断）：在 `s1-selfhost*.md` 增补一句“actual 为候选自报结果，external oracle 用于防漂移而非防恶意候选”；长期依赖 T018/S2 的进程级沙箱与可观测性。

3. `LOW`：本地工具链 Go 1.27.0 不等于 CI 固定 Go 1.22.12；本地全绿不能替代目标 runner 结果。

	- 证据：`ci.yml#L48` 附近 `GO_VERSION: "1.22.12"`；本地 `go1.27.0 windows/amd64`。

4. `LOW`：生产代码离线约束是 import/字符串静态扫描，不是 OS 网络沙箱；命令词表非完备。

	- 证据：`network_test.go:18`、`policy.go:193`（`requireCapabilities` 对 deny/loopback/allowlist 在能力缺失时返回 `ErrPolicyUnavailable`）。文档已如实表述（`s1-selfhost.md` 实现范围第 8 条），维持不阻断。

## 离线门禁退出码（全部 0）

| 命令 | 退出码 |
|---|---:|
| `go build ./...` | 0 |
| `go vet ./...` | 0 |
| `go test -count=1 ./...` | 0（全部包 ok） |
| 聚焦 workflow 契约与变异测试 | 0 |
| `git diff --check` | 0 |

## 变异子测试：实际通过 540 个，分组

| 分组 | 通过数 |
|---|---:|
| `TestWorkflowActionMutationsAreRejected` | 208 |
| `TestWorkflowActionParametersAreLocked` | 78 |
| `TestWorkflowScriptMutationsAreRejected` | 144 |
| `TestWorkflowScriptsRejectCommentedCommands` | 12 |
| `TestWorkflowStructureMutationsAreRejected` | 82 |
| `TestWorkflowYAMLMutationsAreRejected` | 16 |

## 重点复审对照

- **A 完整脚本门禁：通过。** 契约与全部变异测试共用 `decodeWorkflow`（`workflow_test.go:79`）/ `validateWorkflowContract`（`workflow_test.go:147`）；12 个摘要为独立常量（`workflow_test.go:241`）；注释、here-string、未调用 scriptblock、不可达分支、`exit 0`、变量替换、函数遮蔽、伪造 oracle 写入均被拒绝；摘要同时绑定 job/步骤名、顺序、`shell`、`if`、工作目录；“同改即变绿”边界已由文档诚实限定为独立评审。
- **B Action 与工作流结构：通过。** 11 个 Action 按 job/角色锁定仓库 + 完整 SHA（`workflow_test.go:303`），全部 `with` 参数锁定（`workflow_test.go:317`）；8 个输入均为 required string 无 default；permissions/工具链/矩阵/fail-fast 显式 false/needs/步骤顺序全部校验；YAML 节点级检查拒绝未知/重复字段、多文档、null、anchor/alias/merge，并保留字段存在性。
- **C 三 runner 信任链：通过。** candidate-build 只构建封存候选；probe 从 verifier-source 构建 harness 并只上传不可信 oracle；selfhost 全新 runner 重建 verifier/seed、下载同一 binary/oracle、绝不执行候选（candidate metadata checkout 仅 `git show --format=%cI` 读取提交时间）；artifact 名绑定 `runner.os` + candidate SHA；三个信任 job `cache: false` + `GOTOOLCHAIN=local`，不可信状态无法经 Go 缓存进入 verifier 构建。
- **D verifier 与发布证据：通过。** `VerifySelfHost` 绑定 seed/candidate/oracle 所有权与内容；realpath、symlink、目录重叠、路径穿越、trailing JSON、重复 case 均 fail-closed；SHA256SUMS 完整覆盖并拒绝多余/缺失/篡改；SBOM 绑定 source commit 且与许可证清单逐模块一致；`PublisherIdentityAuthenticated: false` 为常量输出（`verify.go:202`）；bootstrap 不含新 verifier，流程明确要求先独立评审落库再 dispatch 不同候选。
- **E 运行时离线：通过。** 独立扫描 0 处 `net` 导入；5 处 `os/exec` 全部在批准边界；三种网络模式在能力缺失时执行前 fail-close（`network_test.go:77`），且未把静态扫描表述为沙箱证明。

## 报告项

- 改动 Go/YAML/JSON 为 UTF-8 无 BOM + LF，Markdown 为 UTF-8 BOM + LF：**30 个改动文件全部通过**。
- `tmp/`：仅 `.gitkeep`，无残留。
- 本地 Go 版本：`go1.27.0 windows/amd64`，**不等于** CI 固定 1.22.12。
- 首次真实 GitHub CI：**未执行**；独立 verifier baseline：**未落库**（HEAD 仍为 `fc99f553...`，未 commit/push/tag/dispatch）；真实 runner 攻击演练：**未执行**。
- T017 状态：`IMPLEMENTED_AWAITING_CI`，DEV_PLAN 未勾选，如实。

本轮为严格只读：未修改任何文件、未触发任何 Git 写操作或 GitHub workflow，sessbridge 未触碰。