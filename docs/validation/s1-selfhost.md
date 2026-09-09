# T017 自托管与可信发布验证报告

[English](s1-selfhost_EN.md)

日期：2026-09-08，更新于 2026-09-10。结论：`COMPLETE`。修复版经双主审并通过 verifier baseline 的 Windows/Linux push CI；针对不同 candidate SHA 从 `main` 手工 dispatch 的 Windows/Linux 完整矩阵也已通过，AT-13/14 证据均已归档。

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
5. 本轮整改不等于重新完成本机两代演练或真实 runner 攻击演练；这些边界不因 CI 通过而扩大。修复版已通过双主审、落库并通过 verifier baseline 的双平台 push CI，随后不同 candidate SHA 的手工 dispatch 完整矩阵也已通过；T017 状态为 `COMPLETE`。

## 当前门禁结果

1. `CGO_ENABLED=0 go build ./...`：通过。
2. `CGO_ENABLED=0 go vet ./...`：通过。
3. `CGO_ENABLED=0 go test -count=1 ./...`：通过；`internal/release` 包含外部信任 pin、防直接联网和发布证据聚焦测试。
4. `cd tools/contracts && npm test`：通过，2/2 测试、82 个独立 fixture 全部匹配。
5. Windows 11 本机完整两代演练：通过；固定 seed 哈希、外部 oracle、SPDX/许可证/SHA256SUMS 均复核成功，临时 worktree 和 `tmp/` 产物已清理。
6. GitHub push CI run `34253473081`：2026-09-09 首次运行，Windows/Linux `test` job 均失败，不构成 AT-14 证据。失败暴露 Go 1.22/原生 runner 差异：进程终止测试未及时回收 child、Unix probe 测试的输出目录与候选二进制同名、许可证测试使用 test binary 导致主模块路径为空，以及 Windows 对不存在目标路径的 overlap 比较未统一解析已有父目录。
7. 上述修复在本地当前工具链通过聚焦及全量测试；[DeepSeek V4 Pro 独立复审](t017-ci-remediation-deepseek-review.md)与 [GPT-5.3 Codex 独立复审](t017-ci-remediation-codex-review.md)均给出 `PASS FOR BASELINE REVIEW`。修复已作为 verifier baseline `18c9340c60deb30049554491f8905df51d0a3180` 落库；GitHub push CI [run 34268076912](https://github.com/larsonzh/prfrail/actions/runs/34268076912) 的 [Ubuntu job](https://github.com/larsonzh/prfrail/actions/runs/34268076912/job/102202411754) 与 [Windows job](https://github.com/larsonzh/prfrail/actions/runs/34268076912/job/102202412040) 均成功。
8. 维护者 `larsonzh` 从 `main` 手工触发 GitHub Actions [run 34277671704](https://github.com/larsonzh/prfrail/actions/runs/34277671704)（CI #3，`workflow_dispatch`）；run 绑定 verifier `18c9340c60deb30049554491f8905df51d0a3180`，输入 candidate `1e7af676e8e84028c7ecfef1ccde91213738c20c`，二者不同。run 结论为 `success`，8/8 jobs 成功：[Go Ubuntu](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102234603547)、[Go Windows](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102234603526)、[candidate build Ubuntu](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102234603542)、[candidate build Windows](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102234603184)、[candidate probe Ubuntu](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102234809504)、[candidate probe Windows](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102234809369)、[bootstrap/release Ubuntu](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102235086119) 和 [bootstrap/release Windows](https://github.com/larsonzh/prfrail/actions/runs/34277671704/job/102235086172)。

## 首次手工 CI artifact

以下为 run `34277671704` 公开 API 返回的 6 个未过期 GitHub Actions artifact；摘要是 GitHub 对 artifact archive 计算的 digest，不替代归档内部经 verifier 检查的 SHA256SUMS：

| Artifact | 大小（字节） | GitHub archive digest |
|---|---:|---|
| `candidate-binary-Linux-1e7af676e8e84028c7ecfef1ccde91213738c20c` | 2,247,948 | `sha256:7e4d569e6754809ce3da97cdebcfad1db058eef9b51578f3dfc9cb56601b8d0e` |
| `candidate-binary-Windows-1e7af676e8e84028c7ecfef1ccde91213738c20c` | 2,397,137 | `sha256:ab8edb01ddb7dafc54ee78bf6942b755fffe5437ec185617416c90c6d4983e66` |
| `candidate-oracle-Linux-1e7af676e8e84028c7ecfef1ccde91213738c20c` | 261 | `sha256:44465973567919c34f39ba869a1b8030fa9611a0107d98ece948e6bb37195caf` |
| `candidate-oracle-Windows-1e7af676e8e84028c7ecfef1ccde91213738c20c` | 261 | `sha256:fdecda0b17ee81a782f5d348f3e697b28ff56f11b2bb24f31ad09f5ae57b58d1` |
| `prfrail-Linux-1e7af676e8e84028c7ecfef1ccde91213738c20c` | 2,249,050 | `sha256:43a99ee59076e3eae96c9ed61016159511f9ac9309169f170befead286d1d41a` |
| `prfrail-Windows-1e7af676e8e84028c7ecfef1ccde91213738c20c` | 2,398,244 | `sha256:7adc8f21c81df88db31d18ea5b83e892d5e7d3a9416088c51b0ec12b15df6ca7` |

## Node 24 Actions 替换基线

2026-09-09，GitHub Actions 全部升级为原生 Node 24 版本并固定完整 commit SHA。GPT-5.3 Codex 与 DeepSeek V4 Pro 分别给出 `PASS FOR NEW VERIFIER BASELINE`；verifier `6024526893fd7acc47de44ee862572ac9ec33303` 随后从 `main` 手工 dispatch，不同的 docs-only candidate 为 `603e8dc03e9fb74ce9c53243d38e704fd8469c6a`。

[run `34325758548`](https://github.com/larsonzh/prfrail/actions/runs/34325758548) 的 `event=workflow_dispatch`、`head_sha=6024526893fd7acc47de44ee862572ac9ec33303`、`conclusion=success`；Go、candidate-build、candidate-probe、bootstrap/release evidence 在 Ubuntu/Windows 共 8/8 jobs 成功。6 个未过期 artifact 如下：

| Artifact | 大小（字节） | GitHub archive digest |
|---|---:|---|
| `candidate-binary-Linux-603e8dc03e9fb74ce9c53243d38e704fd8469c6a` | 2,247,947 | `sha256:cd7523336b44636362290d9532e2a6e4ab45775ae8bbe6e5a7207dc66131e58a` |
| `candidate-binary-Windows-603e8dc03e9fb74ce9c53243d38e704fd8469c6a` | 2,397,139 | `sha256:12d90fd5e60403b249455c15eaf1aafd4eceef8ed6713abc4b26eb73c6db067d` |
| `candidate-oracle-Linux-603e8dc03e9fb74ce9c53243d38e704fd8469c6a` | 261 | `sha256:dfbdf4a319ef0a7fb436eae0c151b350ba394ae167aa42fda7114f84764b19f1` |
| `candidate-oracle-Windows-603e8dc03e9fb74ce9c53243d38e704fd8469c6a` | 261 | `sha256:2f7d7bc568fc67389f6f59a67df9eb3e6df8de8d6b2f2e2202b44f9af9a6bcb8` |
| `prfrail-Linux-603e8dc03e9fb74ce9c53243d38e704fd8469c6a` | 2,249,054 | `sha256:24ee9025608def05341f5e3c037067d73287289898e1ad71c856eb7e9c5f97bf` |
| `prfrail-Windows-603e8dc03e9fb74ce9c53243d38e704fd8469c6a` | 2,398,249 | `sha256:38cb242a023f676bccd2b1ef64a35c62c7a2dc6a951162d4a2323d6951768fdf` |

这些 digest 是 GitHub artifact archive 摘要，不替代归档内 verifier 核验的 SHA256SUMS。Node 24 替换没有扩大 T017 的信任声明。

## T018 最终候选可信运行（嵌套目录变更前）

2026-09-10，维护者从 `main` 手工触发 [run `34358467083`](https://github.com/larsonzh/prfrail/actions/runs/34358467083)。run 的 verifier/head SHA 为 `5f1222752ec3c66ad3b6a1e9034c520ba3e12ee2`，输入 candidate 为不同的 `bcc235496e0bece825047e95f3165a36585a59c0`；Windows/Linux 的 Go、candidate-build、candidate-probe 和 bootstrap/release 共 8/8 jobs 成功，并生成 6 个 artifact。

Windows verified artifact ID 为 `10106867086`，名称为 `prfrail-Windows-bcc235496e0bece825047e95f3165a36585a59c0`，大小 `2,398,240` 字节，未过期，GitHub archive digest 为 `sha256:10d902164194a8b1c6d1121adba5363bb69f829ac3516f7b891c6f5f41d08e06`。其离线核验、升级、回滚、PATH/进程无副作用和五步快速上手结果见[安装演练](s1-install-drill.md)。

该运行早于 verified artifact 的 `prfrail/` 单一顶层目录变更，只证明旧的平铺四文件布局。目录变更修改了受审 workflow 脚本及摘要，必须建立新的 verifier baseline，并以不同 candidate 重跑完整可信矩阵后，才能形成最终载体证据。

## 首次 CI 输入

维护者须从已独立评审的 commit/fixture 计算并手工填写以下必填输入，不得直接信任候选生成值；workflow 必须从 `main` 分支启动，其他 ref 的 self-host job 会被拒绝：

| 输入 | 当前评审值 |
|---|---|
| `candidate_commit` | `1e7af676e8e84028c7ecfef1ccde91213738c20c`；与 verifier 不同的 docs-only candidate，父提交为 verifier baseline |
| `verifier_commit` | `18c9340c60deb30049554491f8905df51d0a3180`；双主审通过且 push CI run `34268076912` 的 Windows/Linux job 均成功 |
| `bootstrap_commit` | `fc99f553c3066b24230a093d14bc17af0c6a398e` |
| `bootstrap_manifest_sha256` | `ae0d5db85d6a45de19dc75d40460705e0a7a0dea20de2a11632c25bf8b4543fc` |
| `expected_oracle_sha256` | `388f88842710abc90676ed0a85f4bd935b594fc552146577c304b0bbbd850aaf` |
| `license_policy_sha256` | `ff9d243c67d516a20147cb5b8684741c3f5ddb77a761af7aa1aab1df4a398013` |
| `noop_chain_sha256` | `fbeee2b64ab174637c77d9cdd1b04bd7fd6c681912f0b73f7a4847bc19a7b31c` |
| `executable_chain_sha256` | `dbc5fc415917398c88a6ed669971c417ed8ee6884a2c8ebfc1271807d757935b` |

## 边界与后续

1. `12fa05d1f56c5e9a03a0964fce9d9f842bb85397` 的 push CI 失败后，修复版经双主审并以 `18c9340c60deb30049554491f8905df51d0a3180` 建立新的绿色 verifier baseline；之后须由该较早 verifier commit 审核不同的 candidate commit。这是一次性的信任引导，不以候选运行结果替代人工审查。
2. 当前实现不发布安装包，也不授权 `git commit`、`git push`、tag、签名或 GitHub Release。
3. 固定 seed 只读使用；失败时丢弃 candidate，不能覆盖或回滚 seed。
4. 新 verifier baseline 的 push CI 与后续手工 dispatch 完整矩阵均已双平台通过；verifier/candidate/run 与 artifact 证据已补录，T017 标记为 `COMPLETE`。
5. T018、安装/升级/卸载实测和最终 S1 exit 仍未完成。