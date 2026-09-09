# T017 Node 24 Actions GPT-5.3 Codex 独立主审存档

[English](t017-node24-actions-codex-review_EN.md)

日期：2026-09-09。历史 verifier：`18c9340c60deb30049554491f8905df51d0a3180`。受审的新 verifier：`6024526893fd7acc47de44ee862572ac9ec33303`。主审：GPT-5.3 Codex。最终结论：`PASS FOR NEW VERIFIER BASELINE`。

本结论批准 `6024526893fd7acc47de44ee862572ac9ec33303` 作为 Node 24 Actions 升级后的新 verifier baseline。它不撤销历史 T017 `COMPLETE`，也不等于新基线的 T017 运行矩阵已经完成，更不是发布授权或发布者身份认证。

## Findings

未发现可复现的 `Critical`、`High`、`Medium` 或 `Low` 阻断项。

1. verifier/candidate SHA 隔离、`main` 上的显式 `workflow_dispatch`、双平台真实 runner、只读仓库权限与 checkout 凭据禁持久化约束均保持不变。
2. candidate-build、candidate-probe、selfhost 的依赖链保持不变；可信 verifier、固定 bootstrap seed、候选元数据及候选二进制仍位于独立目录。
3. 外部 oracle、固定 fixture SHA-256、seed hash、selfhost、release metadata 与最终 verify 的失败传播路径未被削弱。
4. 三个可信 artifact 下载点均显式设置 `digest-mismatch: error`，与 `actions/download-artifact@v8.0.1` 的 fail-closed 摘要校验语义一致。
5. `workflow_test.go` 锁定全部 Action SHA 和参数；`workflow_mutation_test.go` 继续拒绝 Action、脚本、结构、参数和 YAML 绕过变异。测试同步未删除或弱化既有反例。

## Action 供应链核验

以下官方 GitHub Actions 均固定为完整 commit SHA；对应 release 页将该 SHA 标识为版本 tag 的已验证提交，对应 `action.yml` 均声明 `runs.using: node24`：

| Action | 版本 | 固定 SHA |
|---|---|---|
| `actions/checkout` | `v7.0.1` | `3d3c42e5aac5ba805825da76410c181273ba90b1` |
| `actions/setup-go` | `v7.0.0` | `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` |
| `actions/upload-artifact` | `v7.0.1` | `043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` |
| `actions/download-artifact` | `v8.0.1` | `3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` |

官方来源：

1. <https://github.com/actions/checkout/releases/tag/v7.0.1>
2. <https://github.com/actions/setup-go/releases/tag/v7.0.0>
3. <https://github.com/actions/upload-artifact/releases/tag/v7.0.1>
4. <https://github.com/actions/download-artifact/releases/tag/v8.0.1>

## 评审范围与证据边界

1. 提交 `6024526893fd7acc47de44ee862572ac9ec33303` 仅修改 `.github/workflows/ci.yml`、`internal/release/workflow_test.go` 和 `internal/release/workflow_mutation_test.go`。
2. workflow 的实质变化仅为四类 Action 的完整 SHA/版本升级，以及三个 download-artifact 步骤新增 `digest-mismatch: error`；对应契约测试同步锁定这些值。
3. 普通 push CI run `34314888161` 已成功，证明新 Action 组合可完成常规 Windows/Linux CI，但该 run 没有执行完整 T017 信任矩阵，不能作为矩阵完成证据。
4. 本次主审为只读评审：未修改代码、Git 历史或 workflow，未触发发布、签名或 attestation。

## 尚待完成的运行时证据

1. 基于新 verifier 创建一个不同 SHA 的 docs-only candidate；candidate 不得修改受审 workflow、release verifier 或其契约/变异测试。
2. 从 `main` 手工触发完整 `workflow_dispatch`，并确认以下 8 个 job 全部成功：test（Ubuntu/Windows）、candidate-build（Ubuntu/Windows）、candidate-probe（Ubuntu/Windows）、selfhost（Ubuntu/Windows）。
3. 归档 run URL、verifier/candidate SHA、各 job 结果、六个 artifact 的名称/大小/摘要，以及 dispatch 输入摘要。
4. 在上述矩阵全部成功并完成证据归档前，不得宣称新 verifier 基线的 T017 运行验证 `COMPLETE`。

本存档只记录主审结论；历史 verifier `18c9340c60deb30049554491f8905df51d0a3180`、candidate `1e7af676e8e84028c7ecfef1ccde91213738c20c` 与 run `34277671704` 的历史 T017 `COMPLETE` 证据保持有效且不被改写。