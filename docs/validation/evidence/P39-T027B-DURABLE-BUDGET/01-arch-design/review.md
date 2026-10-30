## 1 锁与作用域
- 保留 [internal/tickets/cost.go](internal/tickets/cost.go) 的 `CostLedger` 为纯规则模型；按蓝图新增 `DurableCostLedger`，持久写入口统一包住 `AppendAllocation`、`Reserve`、`Settle`，不向生产调用方暴露可写内存模型。
- `budgetHash` 是 `reservationHash`，不能反推 scope，更不能作为独立锁键。宿主先绑定可信的 `ScopeKind/ScopeHash`，再在该账本核验 reservation；锁粒度为整个共享 scope，所有 reservation 共用一把锁。
- 宿主固定唯一预算根；规范路径由 `(ScopeKind, ScopeHash)` 派生，hash 路径分量只取合法十六进制部分。同 scope 不得按 run、request 或 caller 路径另开账本；路径别名必须归到同一文件身份，无法判定即拒绝。
- 采用稳定旁路文件锁：Windows 对固定锁文件用 `LockFileEx` 锁定固定字节区间；Unix 后端用 `flock`。锁文件不随快照替换、不按超时删除，句柄不继承给子进程；进程崩溃由内核释放锁。
- 进程内按同 scope 串行，再取得跨进程排他锁；锁覆盖重读、规则执行、快照提交及耐久确认。只读检查同样参与锁协议，但不得创建缺失账本或锁文件。
- [internal/console/cli.go](internal/console/cli.go) 的 `executeAgentRunnerRun` 改为打开宿主规范存储；`ExecuteAgentRunnerRun` 将同一持久预算端口注入 admission、dispatcher、reconfirmer、publisher，禁止各自加载独立可写快照。
- 预算锁与 replay/registry 锁绝不嵌套：store 方法返回并释放其锁后才访问预算；预算操作返回并释放锁后才发布下一记录。锁内不调用 store、registry、spawn 或外部等待，重确认仍只是时点判断。

## 2 锁内重读协议
- 每次取得 scope 锁后，重新打开当前规范账本，禁止沿用原子替换前的文件句柄；执行 `DecodeCostLedgerRecord` → `LoadCostLedger`，由既有校验重算摘要、顺序、索引、allocation 链与 summary。
- 将磁盘记录与宿主绑定的账本身份核对：`LedgerID/CreatedAt/ScopeKind/ScopeHash/Currency/PricingMode/PricingVersion` 必须一致；解码成功不等于选中了正确共享账本。
- 本实例保存最后一次已确认快照，设其条目数为 `n`；磁盘允许合法追加，但必须不少于 `n`，且前 `n` 条逐条 canonical 相等。合法尾部追加不能仅因 `LedgerHash` 变化而拒绝。
- 磁盘缩短、旧前缀变化或身份漂移均以 `ErrInvalidCostLedger` fail-closed；不合并分叉、不覆盖赢家、不从旧内存重建。摘要自洽也不能豁免此前缀检查。
- 只在刚重放的局部 `CostLedger` 上执行本次操作，再取 `Snapshot`；提交前再次验证新记录完整保留刚读取的磁盘前缀，业务条目只能合法追加。确认失败不得把局部模型提升为“已确认缓存”。
- 拟改的只读 `RequireOutstandingReservation(ctx, hash)` 每次锁内重读，再沿用 `validateAgentRunnerBudget` 的 run/request/authorization 绑定检查；所需 reservation 还必须属于本实例已耐久确认的前缀，读取新尾部本身不能制造确认。
- 首次打开及显式重新确认在 admission 外完成；未确认的新 reservation 只能先在该边界确认再准入。初始 allocation 也须走同一提交协议；恢复时账本缺失不得解释为可重新初始化。

## 3 耐久确认序列
- 拟新增持久端口：`Reserve(ctx,input)` 返回 `ReservationOutcome, DurableAck, error`；`Settle(ctx,input)` 返回 `SettlementOutcome, DurableAck, error`。`DurableAck` 字段私有，绑定 scope、快照摘要与对应条目摘要，零值无效，不进 wire。
- 锁内以 `Snapshot` 和 `evidence.EncodeCanonical` 生成完整记录，使用 UTF-8 无 BOM、LF；在账本同目录独占新建临时文件，检查所有字节写足，禁止直接截断规范账本。
- 严格顺序：**写完整快照 → 临时文件 `Sync` 成功 → 检查 `Close` 成功 → 同目录经验证的原子替换提交 → 父目录同步屏障成功并检查目录句柄关闭 → 签发确认 → 释放锁**；新建目录链也必须先完成必要父目录屏障。
- 写入、`Sync` 或提交前 `Close` 失败即不执行原子提交，返回错误和无效确认；局部规则已经返回 outcome 或 `Appended=true` 均不能冒充已落盘。
- 原子提交返回错误时仍在锁内重读规范路径：目标快照已可见则返回“可见但未确认”；仅旧快照可见且可证明未提交则返回“未提交”；无法判定则返回“提交结果不明”。三类为拟定内部错误分类，均保留底层原因并可程序化识别。
- 原子提交成功后，目录同步、目录关闭失败，或取消导致确认未完成，必须返回“可见但未确认”、无效确认及已见 outcome；保留规范快照，不回滚、不释放额度。崩溃或确认回执丢失则由恢复执行同键重新确认。
- 幂等命中也必须针对锁内最新完整快照重新执行确认序列；`Appended=false` 不免确认。`settleTerminalIntent` 仅使用原 intent 的计划、时间、键与 R/C 证据，核对拟新增 `ack.RequireEntry(settled.EntryHash)` 成功后才允许 `RecordCompletion`，随后 closure。
- 原子替换及目录屏障后端必须独立具备适用平台/前提的证明，不能借用 replay 的新 link 结论。Windows 仍保持 §7 的 `unproven` 拒绝；预算确认不升级发布等级，新的终局链在 intent/结算之前拒绝，已有完整链仅无写重放。

## 4 unknown 语义
- 沿用 `normalizeSettlementEntry`：`CostSettlementUnknown` 的 `ChargedAmountMicros=nil`，`ObservedCalls=nil`、`ObservedTokens=nil`；这是未知，不是零费用、零调用或零 token，原 reservation 的预留值保持不变。
- unknown **结束 outstanding 资格，但不释放预算占用**：`RequireOutstandingReservation` 返回 `ErrCostReservationSettled`；`OutstandingReservations()` 不再计它，`UnknownHoldReservations()` 仍计它。
- 重启通过 `LoadCostLedger` 重放原 reservation 与 unknown settlement；`outstandingUsage` 继续累计原预留金额、调用数和 token，`validateSharedCaps` 继续强制共享上限，不能只统计“尚无 settlement”的 reservation。
- `UnknownReservedAmountMicros` 按现有口径覆盖未结算及 unknown 的预留金额；金额不可确定时保持空值，不补零。有限金额上限下，无法计算占用继续返回 `ErrCostUnknownCommittedUsage`。
- 保留 `stableTerminalSettlementKey(requestId, completionHash)`：同键同载荷返回原条目并重新确认；同键异载荷为 `ErrCostSettlementConflict`，同 reservation 异键沿用 `ErrCostReservationSettled` 拒绝，不追加第二次结算，本批不提供覆盖 unknown 或释放接口。
- 任一 settlement 已存在、预算耐久未确认、仅 R、已有 launch-intent，或已有终局，均不得据此 relaunch；identity 缺失不是“未启动”证明。恢复只用原 terminal-intent 幂等补齐结算与发布，不从账本反推 completion，也不把 unknown 伪造为完成。

## 5 失败矩阵
| 失败点 | 现象 | 必须的返回值 | 禁止的行为 |
|---|---|---|---|
| 锁、scope、解码、前缀、规则或后端前置检查 | 无法取得可信当前状态，或操作不合法 | 原始锁错误、既有校验/预算/冲突错误或后端拒绝；无有效确认 | 换路径分锁、重置账本、覆盖旧前缀、带病准入 |
| 完整写入、`Sync`、提交前 `Close` | 规范账本尚未提交本次快照 | “未提交”错误及无效确认 | 继续原子提交，以内存 outcome 放行 R 或 C |
| 原子提交报错 | 可能已替换，系统调用结果不足以判定 | 锁内重读后返回“可见但未确认”“未提交”或“提交结果不明”；均无有效确认 | 直接声称失败无副作用、删除可见条目、重用额度 |
| 提交后目录屏障、关闭或确认中断 | 新快照可见，但调用方没有完整确认 | 可返回时返回“可见但未确认”；崩溃/回执丢失后按原键重读并重新确认 | 返回成功、回滚账本、在未获确认时新发布 C/closure、盲目 relaunch |