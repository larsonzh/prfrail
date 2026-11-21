# P84 ⑥ 独立扫描 r2（原位复核）逐字转录 · 2026-11-21

- 载体：`independent-reviewer`（只读）；模型 `MAI Code 2 Flash (copilot)`。
- 额度：⑥ ≤3（本轮 = 第 2 轮，用 2/3）。
- 消毒声明（§7.3）：原文未含凭据 / token / IP / 私密主机名。

---

INDEPENDENT SCAN: PASS

**Section A - 范围**  
已核实：M1 修复位置、C6/C7 分支映射、前驱、异主、恢复及走链代表测试。未核实：本轮未在隔离副本实际执行 A/B 变异，但只读核验映射可信。

**Section B - M1**  
CLOSED。读取 `current` 后立即校验 token/hash：[runtime_supervised.go:1232](internal/console/runtime_supervised.go#L1232)、[runtime_supervised.go:1242](internal/console/runtime_supervised.go#L1242)。该校验先于历史核验和写入：[runtime_supervised.go:1250](internal/console/runtime_supervised.go#L1250)、[runtime_supervised.go:1253](internal/console/runtime_supervised.go#L1253)。C6 在 [p84_t5h4_residual_test.go:357](internal/console/p84_t5h4_residual_test.go#L357) 同时断言拒绝文本、零改写及已代次控制。删除校验会直接使拒绝和字节断言变红。

**Section C - M2**  
CLOSED。上限为 64：[runtime_agent.go:2186](internal/console/runtime_agent.go#L2186)，超限走 `depth check`：[runtime_agent.go:2240](internal/console/runtime_agent.go#L2240)。C7 使用真实 65代链：[p84_t5h4_residual_test.go:442](internal/console/p84_t5h4_residual_test.go#L442)，断言深度拒绝、字节不变及 retained 可读：[p84_t5h4_residual_test.go:485](internal/console/p84_t5h4_residual_test.go#L485)。

**Section D - 回归**  
C1/C2/C3/C6/C7 精确窄测试与异主、renew、再入代表测试均 PASS；完整 `./internal/console`、`./internal/guard` 测试亦 PASS。未发现前驱链、恢复三路径、异主拒绝或走链语义回退。

**Section E - 结论**  
未发现 OPEN 项。本轮未改文件；变异映射可信，但未实跑变异副本。
