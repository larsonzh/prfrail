# T017 CI 修复版 DeepSeek 独立主审存档

[English](t017-ci-remediation-deepseek-review_EN.md)

日期：2026-09-09。受审基线：`12fa05d1f56c5e9a03a0964fce9d9f842bb85397`。受审对象：当前未提交的 T017 CI 修复差异。主审：DeepSeek V4 Pro。最终结论：`PASS FOR BASELINE REVIEW`。

本结论表示修复可进入新的 verifier baseline 候选流程；它不是 T017 完成证明、Go 1.22.12 或 GitHub Windows/Linux CI 成功证明，也不是发布授权。

## Findings

未发现可复现的 `BLOCKER`、`HIGH` 或 `MEDIUM`。

1. `LOW`：`internal/guard/process_test.go` 从 `waitDone` 读取结果时没有独立超时。若未来平台实现出现“身份判定已消失但进程对象未退出”的缺陷，测试会等待至全局超时；当前 `StopProcessIdentity` 成功前已由 `waitIdentityGone` 确认身份消失，因此不阻断 baseline。后续可用带超时的 `select` 加固。

## 上轮阻断项复核

上轮 `HIGH` 已修复：子进程 `Wait` 结果不再匹配平台相关错误文本，而是通过 `errors.As(waitErr, &exitErr)` 接受 `*exec.ExitError`。该判定覆盖 Unix 信号终止与 Windows 非零退出，未扩大生产代码的成功语义。

## 实际评审范围

主审核对到相对 `12fa05d1f56c5e9a03a0964fce9d9f842bb85397` 恰好 9 个改动文件：

1. `internal/guard/process.go`
2. `internal/guard/process_test.go`
3. `internal/release/probe_test.go`
4. `internal/release/verify_test.go`
5. `internal/snapshot/export.go`
6. `docs/validation/s1-selfhost.md`
7. `docs/validation/s1-selfhost_EN.md`
8. `docs/DEV_PLAN.md`
9. `docs/DEV_PLAN_EN.md`

当时无额外改动；`sessbridge` 未修改；`tmp/` 仅 `.gitkeep`。

## 本地验证证据

1. `git diff --check`、`git diff --stat`、`git diff --name-status`：退出码 0；9 文件，79 insertions / 54 deletions。
2. `gofmt -l`：退出码 0，无输出。
3. 编码门禁：9 个文件全部通过；Markdown 为 UTF-8 BOM + LF，Go 为 UTF-8 无 BOM + LF。
4. `CGO_ENABLED=0 GOTOOLCHAIN=local go build ./...`：退出码 0。
5. `CGO_ENABLED=0 GOTOOLCHAIN=local go vet ./...`：退出码 0。
6. `CGO_ENABLED=0 GOTOOLCHAIN=local go test -count=1 -timeout=2m ./...`：退出码 0，全部包通过。
7. guard 停止测试重复 10 次、release 两项修复测试重复 3 次、snapshot overlap 测试重复 5 次：均通过。
8. `GOOS=linux GOARCH=amd64 go test -c ./internal/guard`：交叉编译通过，临时产物已清理。
9. 9 个文件编辑器诊断：无错误。

## 未验证事项与后续门禁

1. 本地工具链为 Go 1.27.0；未在 Go 1.22.12 下运行，Linux 仅完成交叉编译而非原生执行。
2. 修复 commit 尚未落库，GitHub Windows/Linux push CI 尚未复验。
3. 仍需 GPT-5.3 Codex 对包含本归档在内的最终候选差异独立复审。
4. 双主审通过后才可提交新的 verifier baseline；其 push CI 双平台通过后，还须创建不同 candidate SHA 并从 `main` 手工 dispatch 完整 candidate-build/probe/selfhost 矩阵。
5. 只有归档 run URL、verifier/candidate SHA、平台结果和 artifact 摘要后，T017 才能标记 `COMPLETE`。

本轮主审严格只读：未修改文件，未执行 Git 写操作，未触发 GitHub workflow、发布、签名或 attestation。