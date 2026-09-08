# T017 自托管与可信发布验证报告

[English](s1-selfhost_EN.md)

日期：2026-09-08。结论：`IMPLEMENTED_AWAITING_CI`。本机 Windows 验证通过；[GPT-5.3 Codex 主审](t017-baseline-review.md)与 [DeepSeek V4 Pro 主审](t017-deepseek-baseline-review.md)均给出 `PASS FOR BASELINE REVIEW`。T017 仍需首次 GitHub CI 的 Windows/Linux 原生矩阵通过后才能标记 `COMPLETE`。

## 实现范围

1. `.github/workflows/ci.yml`：固定 Go 1.22.12、`CGO_ENABLED=0`、`GOTOOLCHAIN=local` 及 GitHub Actions 完整 commit SHA，在 Windows/Linux 分别执行 build/vet/test；自托管矩阵只由维护者从 `main` 显式 `workflow_dispatch` 启动，并在独立目录构建固定 seed、候选和可信 verifier。三个信任 job 的 setup-go 均显式 `cache: false`，不通过 Go 缓存恢复或保存跨 runner 的构建状态。
2. `internal/release/` 与 `cmd/prfrail-release/`：验证完整 40 位 commit SHA、seed 不变、目录隔离、外部 oracle、SHA256SUMS、SPDX 2.3 SBOM 和实际链接模块许可证。
3. `testdata/selfhost/`：固定 bootstrap commit `fc99f553c3066b24230a093d14bc17af0c6a398e`、四项外部 oracle、noop/可执行链和许可证策略。
4. `cmd/prfrail/main.go`：版本变量允许 CI 通过链接参数绑定固定 candidate commit。

## AT-13/14 证据

1. seed 与 candidate 使用不同源码、构建和运行目录；真实路径重叠、oracle 所有权反转、seed 变更均 fail-close。
2. 外部 oracle 固定 `version`、`validate-noop`、`run-noop`、`reject-executable` 四项结果；candidate 只生成 actual，不能生成 expected。actual 是候选进程的自报结果，external oracle 的相等性检查用于检测行为漂移，不构成候选无恶意的安全证明。维护者在 dispatch 时独立输入 candidate、verifier、bootstrap commit 与五个 fixture SHA-256，workflow 无 candidate 控制的默认值。
3. verifier commit 必须等于承载当前受审 workflow 的 `main` commit，candidate commit 必须是另一个完整 SHA。`candidate-build` job 只构建并封存 candidate binary；全新 `candidate-probe` runner 下载该不可变 artifact 并执行探针，只上传不可信 oracle；最终 `selfhost` runner 再下载同一 binary 与 oracle，从 verifier-source 构建工具且不执行 candidate，只用 verifier commit 中的 oracle/策略完成 `hash-tree/selfhost/generate/verify`。
4. 本机两代演练从固定 seed commit 构建 seed，从当前 candidate 构建候选，四项 oracle 一致；seed、candidate 和 oracle 哈希均由独立发布工具复算。
5. 发布元数据从 Go 二进制 build info 提取实际链接模块；SPDX 文档绑定 source commit，并与许可证清单逐模块核对路径、版本和许可证。未批准许可证、空/不一致 SPDX 包清单、符号 revision、缺失/篡改 artifact、符号链接和不完整 SHA256SUMS 均拒绝。
6. `PublisherIdentityAuthenticated=false` 是固定输出。SHA256SUMS 只证明内容与清单一致，不证明发布者身份；签名/attestation 是后续增强。
7. GitHub workflow 只上传 CI artifact，不创建 Release、不签名、不提交或推送代码。CI 的 checkout、setup-go、npm 与 artifact 上传属于 GitHub 托管构建过程，不属于 ProofRail 运行时。
8. `cmd/` 与 `internal/` 的生产 Go 文件受测试约束，不得直接导入 `net` 或其子包，也不得硬编码常见下载器/联网 Git 子命令；当前 process runner 不声明网络隔离能力，所有网络模式在执行前 fail-close。SPDX namespace 只是离线标识符，不触发网络请求。

## 本轮整改边界

1. 移除以 `strings.Contains` / `strings.Index` 判断脚本安全性的断言。12 个 `run` 脚本按 job/步骤绑定独立常量中的完整 SHA-256，并校验 `shell`、`if`、工作目录和 Action/脚本字段互斥。注释化、字符串伪装、未调用代码块、提前成功退出、变量替换或追加伪造 oracle 写入均不能保留原摘要。
2. 正式工作流契约测试与内存变异测试共用 `decodeWorkflow` 和 `validateWorkflowContract`。结构校验覆盖全部 job 的输入映射、原生矩阵、依赖和步骤顺序；Action 校验精确覆盖仓库、SHA、全部 `with` 参数及执行条件，包括 test checkout 的 `persist-credentials: false`、Setup Go 的版本和缓存策略，以及最终 artifact 上传。
3. YAML 解码拒绝未知/重复字段、多文档、显式默认值、别名/锚点及非许可的 null；保留原始字段存在性，不能通过 `uses` 搭配空 `run`/`shell` 隐藏混合步骤。负例只在内存修改 YAML/结构，不执行替代脚本或修改真实 workflow。
4. 固定摘要是内容漂移门禁，不是 PowerShell 语义证明、发布者身份认证或运行时沙箱。脚本变更即使只改注释也必须重新审查完整脚本及执行上下文，再显式更新基线常量；禁止从待验 workflow 自动刷新期望摘要。修改测试自身仍需独立代码评审。
5. 本轮整改不等于重新完成两代演练，也不代表首次 GitHub CI 或真实 runner 攻击演练已通过。修复版已通过双主审，仍待作为新的 verifier baseline 落库并通过后续 CI；T017 状态保持 `IMPLEMENTED_AWAITING_CI`。

## 当前门禁结果

1. `CGO_ENABLED=0 go build ./...`：通过。
2. `CGO_ENABLED=0 go vet ./...`：通过。
3. `CGO_ENABLED=0 go test -count=1 ./...`：通过；`internal/release` 包含外部信任 pin、防直接联网和发布证据聚焦测试。
4. `cd tools/contracts && npm test`：通过，2/2 测试、82 个独立 fixture 全部匹配。
5. Windows 11 本机完整两代演练：通过；固定 seed 哈希、外部 oracle、SPDX/许可证/SHA256SUMS 均复核成功，临时 worktree 和 `tmp/` 产物已清理。
6. GitHub push CI run `34253473081`：2026-09-09 首次运行，Windows/Linux `test` job 均失败，不构成 AT-14 证据。失败暴露 Go 1.22/原生 runner 差异：进程终止测试未及时回收 child、Unix probe 测试的输出目录与候选二进制同名、许可证测试使用 test binary 导致主模块路径为空，以及 Windows 对不存在目标路径的 overlap 比较未统一解析已有父目录。
7. 上述修复在本地当前工具链通过聚焦及全量测试；[DeepSeek V4 Pro 独立复审](t017-ci-remediation-deepseek-review.md)与 [GPT-5.3 Codex 独立复审](t017-ci-remediation-codex-review.md)均给出 `PASS FOR BASELINE REVIEW`。修复尚未作为新的 verifier baseline 落库，push CI 和首次手工 dispatch 的 Windows/Linux 完整矩阵仍未执行。AT-14 与 T017 保持未完成。

## 首次 CI 输入

维护者须从已独立评审的 commit/fixture 计算并手工填写以下必填输入，不得直接信任候选生成值；workflow 必须从 `main` 分支启动，其他 ref 的 self-host job 会被拒绝：

| 输入 | 当前评审值 |
|---|---|
| `candidate_commit` | 待修复后的 verifier baseline 落库后创建：必须与 verifier commit 不同的完整 candidate SHA |
| `verifier_commit` | 待重新独立评审并落库；`12fa05d1f56c5e9a03a0964fce9d9f842bb85397` 因 push CI 失败不能作为最终绿色 baseline |
| `bootstrap_commit` | `fc99f553c3066b24230a093d14bc17af0c6a398e` |
| `bootstrap_manifest_sha256` | `ae0d5db85d6a45de19dc75d40460705e0a7a0dea20de2a11632c25bf8b4543fc` |
| `expected_oracle_sha256` | `388f88842710abc90676ed0a85f4bd935b594fc552146577c304b0bbbd850aaf` |
| `license_policy_sha256` | `ff9d243c67d516a20147cb5b8684741c3f5ddb77a761af7aa1aab1df4a398013` |
| `noop_chain_sha256` | `fbeee2b64ab174637c77d9cdd1b04bd7fd6c681912f0b73f7a4847bc19a7b31c` |
| `executable_chain_sha256` | `dbc5fc415917398c88a6ed669971c417ed8ee6884a2c8ebfc1271807d757935b` |

## 边界与后续

1. `12fa05d1f56c5e9a03a0964fce9d9f842bb85397` 已完成独立评审和首次落库，但 push CI 失败。修复会改变 verifier 相关代码，必须重新独立评审并建立新的绿色 verifier baseline；之后由该较早 verifier commit 审核不同的 candidate commit。这是一次性的信任引导，不以候选运行结果替代人工审查。
2. 当前实现不发布安装包，也不授权 `git commit`、`git push`、tag、签名或 GitHub Release。
3. 固定 seed 只读使用；失败时丢弃 candidate，不能覆盖或回滚 seed。
4. 新 verifier baseline 的 push CI 与后续手工 dispatch 完整矩阵均须双平台通过；补录 verifier/candidate/run 证据后才可将 T017 标记 `COMPLETE`，任一失败均保持阻断。
5. T018、安装/升级/卸载实测和最终 S1 exit 仍未完成。