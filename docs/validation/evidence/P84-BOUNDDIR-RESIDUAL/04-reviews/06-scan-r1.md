# P84 ⑥ 独立扫描 r1 逐字转录 · 2026-11-21

- 载体：`independent-reviewer`（只读）；模型 `MAI Code 2 Flash (copilot)`。
- 额度：⑥ ≤3（本轮 = 第 1 轮）。
- 消毒声明（§7.3）：原文未含凭据 / token / IP / 私密主机名。

---

INDEPENDENT SCAN: FINDINGS

**Section A**
总体结论：发现 2 项 Medium；T4 原生反例已通过，T5#4 仍存在批量投影与恢复边界风险。

**Section B**
- Medium — [runtime_supervised.go](internal/console/runtime_supervised.go#L1225-L1237)：读取 store 的 `current` 后仍回写调用方传入的 `refreshed`，未校验二者相等；并发代次推进时可能留下旧代次 proof。批量注释声称“以 store 为权威”，但输出仍是参数，不构成线性化投影。
- Medium — [runtime_agent.go](internal/console/runtime_agent.go#L2186-L2243)：保留历史 walk 固定上限 64。超过上限的真实 retained chain 会返回 depth 错误，即使该 proof 的 hash、claim 和 retained token 都匹配；当前 C2 只覆盖 2 代。

**Section C**
- T4：句柄相对走链位于 [bounddir_windows.go](internal/evidence/bounddir_windows.go#L572-L642)，NTSTATUS 映射位于 [bounddir_windows.go](internal/evidence/bounddir_windows.go#L719-L726)，Ordinal 比较位于 [bounddir_windows.go](internal/evidence/bounddir_windows.go#L737-L745)。枚举 class=10 位于 [bounddir_windows.go](internal/evidence/bounddir_windows.go#L270-L315)，布局断言可拒绝 ABI 不匹配。
- T5：前驱链逐项校验 claim、代次、PreviousHash 和 retained 原件，位于 [runtime_agent.go](internal/console/runtime_agent.go#L2213-L2251)。异主接管仍保留 PID/start-token 与旧 holder 停机证明。
- 工作树无 schema、CLI 或 go.mod/go.sum 越界改动；P84 两组新增测试当前均为未跟踪文件。

**Section D**
- 可证伪：本地重跑 P84W1/W2/W3 与 P84 C1/C2/C3 均通过；C1 外 claim/hash、C2 多代、C3 跨进程恢复均被真实执行。
- P52 内部 junction 放行、界外/refeparse 回退拒绝均绿；P83 末级 reparse、异主拒绝和 exact-predecessor 条件未见放宽。
- 两项 Medium 缺少并发代次推进与 `>64` retained 代次反例，不能由现有测试证伪。

**Section E**
- 本轮重跑了 P84 聚焦测试；未重跑全仓全套，盘面 [test-all-post3.txt](tmp/P84-BOUNDDIR-RESIDUAL/test-all-post3.txt) 全包通过，门禁读数 [gate-tree-post3b.txt](tmp/P84-BOUNDDIR-RESIDUAL/gate-tree-post3b.txt) 为 `TOTAL_FAIL=0`。
- 未验证 Windows 386/ARM、跨文件系统真实 64 MiB+ 目录枚举、停机与目录 fsync 的 power-loss 窗口。
