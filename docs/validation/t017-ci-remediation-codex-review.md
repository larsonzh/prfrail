# T017 CI 修复版 GPT-5.3 Codex 独立主审存档

[English](t017-ci-remediation-codex-review_EN.md)

日期：2026-09-09。受审基线：`12fa05d1f56c5e9a03a0964fce9d9f842bb85397`。受审对象：当前未提交工作树相对该基线的完整 T017 CI 修复差异。主审：GPT-5.3 Codex。最终结论：`PASS FOR BASELINE REVIEW`。

本结论表示修复可以作为新的 verifier baseline 候选；它不是 T017 完成证明、Go 1.22.12 或 GitHub Windows/Linux CI 成功证明，也不是发布授权。

## Findings

未发现可复现的 `BLOCKER`、`HIGH` 或 `MEDIUM`。

1. `LOW`：`internal/guard/process_test.go` 从 `waitDone` 读取结果时没有独立超时。若未来平台实现出现“身份探测已消失但 `Process.Wait` 未返回”的异常，测试会等待到整体超时。当前主流程仍 fail-closed，且本地与重复实测稳定通过，因此不阻断 baseline。最小加固方式是使用带独立超时的 `select`。

## 静态审查结论

1. `StopProcessIdentity` 对无 deadline context 使用有界等待，`grace <= 0` 回落到 1 秒；父 context 已取消或已有 deadline 时不会放宽父约束。
2. 终止不确定仍返回 `ErrTerminationUncertain`，未把存活、PID 复用或身份不确定错误报告为 `stopped`。
3. 子进程 `Wait` 结果通过 `errors.As(..., *exec.ExitError)` 判定，不依赖平台错误文本。
4. probe 路径修改只消除 Unix 无扩展名二进制与输出目录冲突，未改变 production 行为。
5. 许可证测试构建真实 `cmd/prfrail` 二进制，保留缺失许可证 fail-closed 断言。
6. export 路径比较解析最深已存在父路径并拼回缺失后缀；非 `IsNotExist` 错误继续 fail-closed，双向 overlap 检查保持不变。
7. 文档准确记录 GitHub Actions run `34253473081` 双平台失败；T017 仍为 `[ ]` 和 `IMPLEMENTED_AWAITING_CI`。

## 实际评审范围

主审核对到相对 `12fa05d1f56c5e9a03a0964fce9d9f842bb85397` 恰好 11 个文件：5 个 Go 修复、4 个状态文档和 2 个 DeepSeek V4 Pro 修复版评审档案。无额外改动；`sessbridge` 未修改；`prfrail/tmp/` 仅 `.gitkeep`。

## 本地验证证据

1. `git diff --check`、`git diff --stat` 和目标文件完整差异检查：退出码 0。
2. `gofmt -l`：退出码 0，无输出。
3. `CGO_ENABLED=0 GOTOOLCHAIN=local go build ./...`：退出码 0。
4. `CGO_ENABLED=0 GOTOOLCHAIN=local go vet ./...`：退出码 0。
5. `CGO_ENABLED=0 GOTOOLCHAIN=local go test -count=1 -timeout=2m ./...`：退出码 0，全部包通过。
6. guard 停止测试重复 10 次、release 两项修复测试重复 3 次、snapshot overlap 测试重复 5 次：均通过。
7. 11 个评审文件编码门禁通过，编辑器诊断无错误。

本地工具链为 Go 1.27.0 windows/amd64，未在 Go 1.22.12 或 GitHub Windows/Linux runner 上复验。

## 未完成门禁

1. 修复 commit 尚未落库，新 verifier baseline 的 Windows/Linux push CI 尚未运行。
2. 仍须创建不同 candidate SHA，并从 `main` 手工 dispatch 完整 candidate-build/probe/selfhost 矩阵。
3. Windows/Linux 全部信任 job 成功并归档 run URL、verifier/candidate SHA、平台结果和 artifact 摘要后，T017 才能标记 `COMPLETE`。

本轮主审严格只读：未修改、格式化、暂存或删除文件，未执行 Git 写操作，未触发 GitHub workflow、Release、签名或 attestation，未联网或安装依赖。