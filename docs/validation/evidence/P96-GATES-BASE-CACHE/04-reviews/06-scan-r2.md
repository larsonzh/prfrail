# ⑥ 独立扫描 r2（逐字原文）

载体 `independent-reviewer` + `MAI Code 2 Flash (copilot)`。切片 P96 `GATES-BASE-CACHE`。日期：2026-11-27。

---

证据定位基本齐了；我现在汇总成严格的 C/P/O 判定，并把“条款文本与实现细节”的差异与运行风险分开，不混成同一结论。

## Section A
`INDEPENDENT SCAN: FINDINGS`

## Section B
- MEDIUM | `tools/gates/lib/baseCache.js:56,225-237` | M-1 未完全关闭：任意 `ok=<digits> FAIL=0` 都能通过形态校验；读取阶段无来源绑定，伪造记录仍可能回放 | 最小修复：给条目增加当前进程可验证的来源摘要，并补 exact-shape 反例
- MEDIUM | `tools/gates/lib/baseCache.js:79-95,167-168` | 已删除跟踪文件先被丢弃，不计入变更路径上限；M-2 对删除形态仍不完整 | 将所有跟踪状态变化计入独立路径计数；仅对仍存在路径读取字节
- MEDIUM | `tools/gates/lib/baseCache.js:48,203-212`；`docs/DELIVERY_DIRECTIVE.md:646-647` | `GOTOOLCHAIN` 是明确的工具链反例，但不在六变量键中，可能热遍掩盖切换工具链冷遍 | 键加入 `GOTOOLCHAIN`，并补 T157 或独立变异
- LOW | `tools/gates/lib/baseCache.js:20-23,203-212`；`docs/DELIVERY_DIRECTIVE.md:646` | 顶部注释及条款未列 ignored `.go` 与 gofmt 二进制段，存在文档漂移 | 同步键段完整清单

## Section C
- H-1：**CLOSED**。`baseCache.js:178-184` 枚举并摘要 ignored `.go`；`:203-210` 将摘要入键；T165 位于 `selftest.js:8304`。
- H-2：**CLOSED**。`baseCache.js:187-198` 摘要 gofmt 二进制；GOROOT 同时支持键值型与纯值型；T165 位于 `selftest.js:8310`。
- M-1：**PARTIAL**。类型校验已收紧为固定 label/detail 形态，但仍缺来源绑定；T165 只验证任意 `"forged"` 被拒绝，未覆盖 plausibly exact-shape。
- M-2：**PARTIAL**。修改型跟踪文件已限制至 200，T165 位于 `selftest.js:8316`；删除型跟踪路径仍未计入变更上限。
- M-3：**CLOSED / 可接受设计边界**。注入式 T157–T165 位于 `selftest.js:7903-7907`，不 spawn 真实工具链是速度与离线约束；冷/热、三作用域证据由 r2 门禁读数覆盖（`gate-readings-r2.txt:7-11`）。这不是条款违例，缺口由 ⑨ 原生验证补位。

## Section D
可证伪性成立：M-1 可用精确形态伪造记录复验；M-2 可用 201 个删除状态复验；GOTOOLCHAIN 可变更其 go env 值复验；键顺序可机械重算。当前 r2 已证实 M2 修改型、H1/H2、GOROOT 双形态及冷热字节等价。

## Section E
本轮只读未跑真实工具链或写盘；无法现场复验 GOTOOLCHAIN 切换、条目篡改及删除型跟踪路径的完整边界。
