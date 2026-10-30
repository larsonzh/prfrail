# 证轨 P39 报告：C6 预算落盘（准则 v1.60）

切片 `P39-T027B-DURABLE-BUDGET`（实现类 `[IMPL]`）；执行日 2026-10-31；准则 v1.60；R = 4、D = 0。

## 1 执行摘要
- **落地内容**：新增持久预算账本与跨进程共享锁，`Reserve` / `Settle` / `RequireOutstandingReservation` 全部改为锁内重读后决策。
- **修复范围**：C6（预算未落盘）——共享上限不可跨进程绕过、提交失败不签发确认、`unknown` 占用重启后仍计入。
- **验证结果**：`go build ./...` 与 `go vet ./...` 通过；`go test ./...` 全绿；新增反例矩阵六条腿全部通过。
- **核心判断**：① 设计 → ③ 红测 → ② 转绿 → ⑤ 首轮 FAIL 修 5 项 BLOCKER → ⑤ 复审 PASS → ⑥ CLEAR → ⑦ PASS。

## 2 前置与实现
### 2.1 ⑤ 预审（两轮）
- **首轮 FAIL**：五项 BLOCKER——只读查询可返回未确认尾部、提交三态不可区分、只读路径创建锁文件、账本缺失重发初始额度、矩阵多条腿未真实执行。
- **修复与复审**：② 逐项修复并新增分阶段故障注入面与可取消锁等待；⑤ 复审 PASS，另记 4 项非阻塞发现。
### 2.2 ②③ 实现与测试
- **测试先行**：③ 先写红测（引用未定义接口）→ ② 实现转绿。
- **矩阵覆盖**：双进程并发上限、锁身份稳定与别名、历史漂移拒绍、只读不产生确认、逐阶段提交故障、账本缺失拒绍重发。
- **平台登记**：`TestDurableHelperReserve` 与 symlink 子例的 Skip 为子进程守卫与平台能力限制，已逐条登记。

## 3 独立扫描（⑥）
```sixstage
## Section A 扫描范围与方法
- 范围：实现文件与测试文件；方法：只读静态扫描 + 契约对照。

## Section B 发现
| # | SEVERITY | 问题 | 定位 |
|---|---|---|---|

## Section C 覆盖与遗漏
- 覆盖：并发锁序、原子提交、错误分类、测试可信度、契约一致性。
- 遗漏：无。

## Section D 与前一环节结论的差异
- 与 ⑤ 复审结论一致，无新增分歧。

## Section E 建议
- 建议在契约修订片升级 Windows 目录屏障等级。

INDEPENDENT SCAN: PASS
```

## 4 终审（⑦）与待裁定
- **终审结论**：⑦ 判 PASS——C6 判定为已解决，未发现死锁、长期得不到锁、双写、确认伪造或平台语义错误。
- **待裁定**：Windows 目录屏障仍按契约维持 `unproven`，是否在 `CONTRACTS-S7-DURABILITY` 片提升耐久等级。
- **观察项候选**：⑤ 预审的四项非阻塞发现登记于此，留待 `CONTRACTS-S7-DURABILITY` 或后续重构片处理——其一，持锁进程崩溃后的锁释放语义无测试腿；其二，`rename` 成功后报错的可见分支无测试腿；其三，`newDurableAck` 的 `confirmed` 字段命名与语义不符；其四，矩阵「确认前取消」腿未落地。

## 5 流程声明
- **载体与模型**：① `deep-reasoner` + `GPT-6.1 Sol (copilot)`（用户特批）；②③⑧ `DeepSeek V4.1 Flash (deepseek)`；⑤ `DeepSeek V4 Pro (deepseek)`；⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`；⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`。
- **偏离留痕**：② 首次派单误用 `DeepSeek V4 Pro (deepseek)`，产物四件已作废并留证（见 §7）。

## 6 制品登记
变更集包括实现文件、一个测试文件、证据包、台账对与 v28 清单。

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/CHECKLIST_v28.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P39-T027B-DURABLE-BUDGET/01-arch-design/review-batch-b.md	repo
docs/validation/evidence/P39-T027B-DURABLE-BUDGET/01-arch-design/review.md	repo
docs/validation/evidence/P39-T027B-DURABLE-BUDGET/MANIFEST.md	repo
docs/validation/evidence/P39-T027B-DURABLE-BUDGET/SHA256SUMS.txt	repo
docs/validation/evidence/P39-T027B-DURABLE-BUDGET/freedom-list.md	repo
docs/validation/p39-durable-budget.md	repo
docs/validation/p39-durable-budget_EN.md	repo
internal/tickets/durable_cost_fault.go	repo
internal/tickets/durable_cost_ledger.go	repo
internal/tickets/durable_cost_ledger_test.go	repo
internal/tickets/durable_cost_lock_shared.go	repo
internal/tickets/durable_cost_lock_unix.go	repo
internal/tickets/durable_cost_lock_windows.go	repo
docs/validation/evidence/P39-T027B-DURABLE-BUDGET/02-implement/voided-v4pro-waste.md	repo
```

## 7 异常处理 / 事实包
1. **模型绑定偏离**：② 首次派单误用 `DeepSeek V4 Pro (deepseek)`，违反 §2.2 角色-模型绑定；止损动作 = 立即停止、四件产物全部作废并记sha256、改由 `DeepSeek V4.1 Flash (deepseek)` 测试先行重做。
2. **⑤ 首轮 FAIL**：属流程内正常回路，五项 BLOCKER 已闭环，未计入偏离。
3. **平台限制**：Windows 无 `chmod` 语义，原只读目录近似腿已改为故障注入腿。
4. **无其它异常**：无越权调用、无未登记临时件、无契约改动。
