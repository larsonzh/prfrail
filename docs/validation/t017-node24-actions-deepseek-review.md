# T017 Node 24 Actions DeepSeek V4 Pro 独立主审存档

[English](t017-node24-actions-deepseek-review_EN.md)

日期：2026-09-09。历史 verifier：`18c9340c60deb30049554491f8905df51d0a3180`。受审的新 verifier：`6024526893fd7acc47de44ee862572ac9ec33303`。主审：DeepSeek V4 Pro。最终结论：`PASS FOR NEW VERIFIER BASELINE`。

本结论批准该提交作为 Node 24 Actions 升级后的新 verifier baseline；它不是新基线 T017 运行矩阵的完成证明、发布授权或发布者身份认证，也不改写历史 T017 `COMPLETE` 证据。

## Findings

未发现 `Critical`、`High`、`Medium` 阻断项；无阻断级 `Low` 发现。

不构成阻断的观察：

1. 普通 `test` job 的 setup-go 使用 `cache: true`，三个信任 job 均使用 `cache: false`；该差异是既有隔离设计并由契约测试锁定。
2. upload-artifact v7 默认 `overwrite: false`；三个上传名称均包含平台与 candidate SHA，同一 run 内不会冲突。
3. artifact 下载不保留可执行位；Linux probe 在执行 candidate 前显式运行 `chmod +x`，既有缓解未被升级移除。

## 独立核验结论

1. 区间 `18c9340..6024526` 仅有两个提交：`971417e` 只修改四份 T017 完成证据文档；`6024526` 只修改 `.github/workflows/ci.yml`、`internal/release/workflow_test.go` 和 `internal/release/workflow_mutation_test.go`。
2. `6024526` 的实质变化仅为四类 Action 的版本/完整 SHA 升级、三个 download 步骤新增 `digest-mismatch: error`，以及对应契约测试同步；信任拓扑、触发器、权限与脚本体未改变。
3. 官方 release/tag/commit/action.yml 核验确认以下完整 SHA 均属于官方 `actions/*` 仓库的对应版本 tag，release 页标记为已验证提交，且 `action.yml` 均声明 `runs.using: node24`：

| Action | 版本 | 固定 SHA |
|---|---|---|
| `actions/checkout` | `v7.0.1` | `3d3c42e5aac5ba805825da76410c181273ba90b1` |
| `actions/setup-go` | `v7.0.0` | `b7ad1dad31e06c5925ef5d2fc7ad053ef454303e` |
| `actions/upload-artifact` | `v7.0.1` | `043fb46d1a93c77aae656e7c1c64a875d1fc6a0a` |
| `actions/download-artifact` | `v8.0.1` | `3e5f45b2cfb9172054b4087a40e8e0b5a5461e7c` |

4. breaking changes 对本工作流无负面影响：checkout v7.0.1 的修复与依赖更新不改变所用输入；setup-go v7 的 ESM/缓存依赖升级不改变 `go-version`/`cache` 语义；upload-artifact v7 新增 `archive` 直传但本工作流未启用，仍按默认压缩上传；download-artifact v8 按 Content-Type 处理解压并将摘要不匹配默认改为 `error`，属于 fail-closed 强化。GitHub-hosted runner 不受 v7 对 self-hosted runner 最低版本要求的影响。
5. download-artifact v8 的 `digest-mismatch` 合法值为 `ignore`、`info`、`warn`、`error`，默认值为 `error`；三个可信下载点均显式设为 `error`。upload v7 记录服务端 artifact digest，download v8 按该 digest 校验，两者协议兼容。
6. verifier 等于 `GITHUB_SHA`、candidate 不等于 verifier、仅从 `main` dispatch、candidate-build → candidate-probe → selfhost 依赖、双平台矩阵、`contents: read`、`persist-credentials: false`、信任 job `cache: false`、外部 oracle 与固定 fixture SHA-256 均保持不变。
7. 契约测试继续锁定 Action 名称、完整 SHA、步骤角色、参数、权限、触发器、输入、job 依赖、双平台矩阵和脚本摘要，并拒绝 foreign action、changed SHA、role swap、参数删除/修改、跳过步骤、提前成功、伪造 oracle 及 YAML anchor/merge/null/duplicate/trailing-document 等变异；同步修改未删除或弱化既有反例。

官方来源：

1. <https://github.com/actions/checkout/releases/tag/v7.0.1>
2. <https://github.com/actions/setup-go/releases/tag/v7.0.0>
3. <https://github.com/actions/upload-artifact/releases/tag/v7.0.1>
4. <https://github.com/actions/download-artifact/releases/tag/v8.0.1>
5. <https://github.com/actions/upload-artifact/releases/tag/v7.0.0>
6. <https://github.com/actions/download-artifact/releases/tag/v8.0.0>

## 已有 CI 证据边界

普通 push CI run `34314888161` 已成功，只能证明 Node 24 Action 组合可在 GitHub-hosted Windows/Linux runner 上完成常规 build、vet、test 与契约夹具。该 run 中三个信任 job 受 dispatch 条件限制而未执行，不能证明 T017 信任矩阵完成。

## 尚待完成的运行时证据

1. 基于新 verifier 创建一个不同 SHA 的 docs-only candidate；candidate 不得修改受审 workflow、release verifier 或其契约/变异测试。
2. 从 `main` 手工 dispatch 完整矩阵，确认 test、candidate-build、candidate-probe、selfhost 各自在 Ubuntu/Windows 上成功，共 8 个 job。
3. 归档 run URL、verifier/candidate SHA、各 job 结果、六个 artifact 的名称/大小/摘要与 dispatch 输入摘要。
4. 矩阵成功并归档前，不得宣称新 verifier 基线的 T017 运行验证 `COMPLETE`。

本轮主审严格只读：未修改或创建文件，未 commit/push，未触发 workflow、Release、签名或 attestation。历史 verifier `18c9340c60deb30049554491f8905df51d0a3180`、candidate `1e7af676e8e84028c7ecfef1ccde91213738c20c` 与 run `34277671704` 的 T017 `COMPLETE` 证据保持有效。