# 第二次 Sol 红队审查原文（P54 CONTRACTS-S7-DURABILITY · design-r2 收窄复核）

载体：`independent-reviewer` + `GPT-6.1 Sol (copilot)`（§1.7.4 临时授权，仅本片有效、不得作为先例；ADR-029 已登记；2026-11-04，只读）。以下为红队返回文本逐字归档。

**1. 三个 High 的闭合质量**
- **High 1：未完全闭合。** 对照首审 [red-team-sol-review.md:8](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/red-team-sol-review.md#L8)，r2 已补明未提交中间态、校验顺序、三时点及锁内截断重试；但丢失识别只在预建记录仍有效存在时成立，补充 A 的历史全集保证存在反例，见第 3 项。
- **Medium｜允许集仍不一致。** [design-r2.md:60](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md#L60) 同段写“该轮盘点读者可见 ∈ {absent, dispatched-unknown} 均落允许集”，又写“头写后恒 dispatched-unknown”。前者允许头写后读到 absent，后者拒绝，验收器可据此实现成不同判据。**最小修复：**区分进程 kill 与断电恢复，统一各切点允许集；把“头已写”与“恢复时头有效可见”分开，不改三判据本体。
- **High 2：无剩余发现，设计层闭合。** 对照首审 [red-team-sol-review.md:15](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/red-team-sol-review.md#L15)，[design-r2.md:52](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md#L52) 明写“ownership/launch/terminal 继续直呼 `writeReplayRecordNoReplace`”；A2 也明确仅 R/C 使用槽容器，canonical JSONL 由 payload 保持。实现后的调用域仍待验证。
- **High 3：无剩余发现，闭合。** 对照首审 [red-team-sol-review.md:29](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/red-team-sol-review.md#L29)，[design-r2.md:11](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md#L11) 已逐字写出“本裁定推翻 r1 的 High 6 判断”，并明确旧解析作废。生产次序确为 [声明集校验](internal/console/runtime_agent.go#L148) 先于 [store 构造](internal/console/runtime_agent.go#L181)，支持 Q1=B。

**2. 七条边界映射**
**无映射留白。** [design-r2.md:33](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md#L33) 起的七条映射均能在 A2 找到对应原句：二.1 装配调用点与不改签名；二.2 全集、不得分批；二.3 任一失败整 run 停止；二.4 存在槽只核验及缺槽分类；二.5 非生产入口不强制；二.6 装配层相对 bootstrap 的偏离登记；二.7 落记录、集合比较、格式与发布路径。

这是**文本覆盖成立**，不代表二.4／二.7 的恢复保证已经证明，问题如下。

**3. 补充 A**
- **High｜当前声明自洽不能证明历史全集不变。** [design-r2.md:80](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md#L80)：“补建出的集合与丢失前逐项相同，漂移检测的丢失无实际安全后果”。反例：同一 run 原全集为 `{a}`，预建记录丢失；重开提供另一组合法声明 `{b}`。本次校验与确定性计算不证明 `{b}={a}`。按现稿，记录缺失分支可核验／补建 `{b}` 并发布新记录，历史集合漂移不再被拒绝。保留下来的 R/C 绑定、orphan 与 runID 校验均不能证明两个全集相等。**最小修复：**补出独立、跨重开绑定原全集的证据；无法证明时停止并报裁。若要接受漂移检测保证一并丢失，须明确登记该降级并取得裁定，不能称为“无实际安全后果”。

- **(a) 分支结论：条件成立，不是历史判别。** [design-r2.md:75](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md#L75) 的两分支，仅在“槽缺失，且记录已有效解码或明确不存在”的限定下互斥、穷尽；不覆盖记录损坏、读取失败等错误。这些错误应按 A-4 拒绝。“记录缺失”只能**按策略视作**从未预建，不能证明历史上确实未预建。

- **(b) 局部安全性：无反例。** no-replace 创建空槽、既存槽只核验，可以避免覆盖**仍存在**的已发布字节；空槽也不会凭空形成 C 或终局。但这不足以推出历史身份、计划与派发防重保证全部不受损。

- **(c) 与既有实现：仅缺失态口径对齐。** [State:294](internal/adapters/agent_runner_replay_store.go#L294) 在 R/C 均缺失时确实返回 `absent`；[有界重读:579](internal/adapters/agent_runner_replay_store.go#L579) 对该组合直接返回，只有 R 缺、C 在才继续重读。因此“不能区分两种来历”属实，不能由此推出历史全集未变。另外，[碰撞路径:513](internal/adapters/agent_runner_replay_store.go#L513) 一旦确认碰撞便只重读，不重新写入；不能把可观察缺失概括成无条件重新取得身份的许可。

- **Low｜顺序论证的适用域过宽。** [design-r2.md:95](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md#L95)：“而且是唯一安全序”，并以“无首派发、无 R/C 已发布”为前提。记录丢失后重开可以仍有已发布 R/C，这个前提不普遍成立。**最小修复：**将该前提限定为首次预建；将结论限定为“当前单阶段记录格式下，避免未完成预建被误判 corruption 的顺序”。程序执行先后也不等于断电后的持久化顺序；反序拒绝运行主要是误停问题，不能不加限定称为不安全。

- **Low｜自身哈希口径冲突。** [design-r2.md:73](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md#L73)：“不参与任何 recordHash”；但 [design-r2.md:93](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md#L93) 明定自身 recordHash 域。**最小修复：**改成“不参与 R/C 的 recordHash；自身 recordHash 按 A-4 校验”。

**4. 补充 B**
- **nil receiver 与默认无操作论证：无 Go 语义问题，条件性成立。** 指针 receiver 可以为 nil；若方法先判 nil 接缝、不访问 receiver 字段、不调用其他逻辑，返回 nil 可以成立。共享调用也不必改变步骤 4–7 的持久状态与数据输入。此结论不是源码、指令或调度时序相同的证明。
- **Medium｜BA3 会把正常状态迁移判成不等价。** [design-r2.md:109](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md#L109)：“`State`/`RecordRequest`/`RecordCompletion` 在调用前后对任意 requestId 的返回完全一致（等值重放）”。在同一 store 上依次探测，第一次 `RecordRequest` 可返回 [first-dispatch](internal/adapters/agent_runner_replay_store.go#L348)，第二次本就应返回 [unknown-block](internal/adapters/agent_runner_replay_store.go#L350)；completion 首写与重放布尔值也不同。**最小修复：**从两个等价初始状态分别运行对照臂与新增调用臂，比较同一操作序列的返回和落盘结果；不要把会写入的接口当作无副作用探针。
- **BA1、BA2、BA4、BA5：无可证伪性缺口。** 分别可以用参数／nil receiver 用例、目录快照、构建检查及构建文件赋值审计验证；目前是待执行断言，不是已通过证据。
- **Low｜API 形态与执行路径混写。** [design-r2.md:105](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md#L105)：“未在既有代码路径上插入条件判断”；但装配路径仍新增调用与错误检查。[design-r2.md:118](docs/validation/evidence/CONTRACTS-S7-DURABILITY/01-arch/design-r2.md#L118) 的“指令层零调用”也不能由“调用点不变、改调空实现”推出。**最小修复：**分别说明新增 API、装配调用及不可达错误分支；只主张限定观察面上的行为等价，删除指令级保证。

**5. 引用纪律**
- **Low｜判定表注记仍缺要求的完整括注。** [contracts-s7-durability.md:66](docs/validation/contracts-s7-durability.md#L66)：“本表 C2 行的 proven 为候选级、前提 P′ 内，非平台级”。语义限制虽已解释，但按本轮“逐处携带完整括注”的字面要求，此处仍未携带指定字符串。**最小修复：**在 `proven` 后补“（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证）”，或删除重复等级断言，仅指向已带括注的表行。
- 其余所审新增 C2 等级引用：**无发现**。A2、F 及判定表两处等级行均带完整括注；首审归档中明确作为缺陷原文保存的引用，不应反向改写。

**6. 边界合规**
**无新增越界发现。** 所审设计没有要求修改三判据本体、ADR-013 结论、构造器签名或提升 Windows 平台等级；共享无操作调用可按限定的 Unix 可观察行为等价理解。[Windows 分支:5](internal/adapters/agent_runner_replay_store_windows.go#L5) 当前仍返回 `unproven`。报告 §1④已披露历史独立登记缺口，本轮不将其改写为历史授权齐备。

全程只读，未产生 diff；未运行 build/test/lint 或故障注入，也未将既有门禁读数当作本轮新证据。以上是设计反例，不声称当前生产已发生错误派发。

RED-TEAM: FINDINGS — Q1=B、七条映射及接缝收窄已闭合，但历史全集绑定缺口、允许集冲突与 BA3 验证错误仍需修订后复核。
