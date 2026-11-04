# input-list — `CONTRACTS-S7-DURABILITY`（P54）输入包清单

口径：OB-93（逐件字节数，不复制入包者仅登记路径与字节数）。采集时点 2026-11-04，`HEAD` = `2b8ac3e`。

## 本片直接输入

| 路径 | 字节 | 用途 |
|---|---|---|
| `docs/t027/CHECKLIST_v43.md` | 23148 | 跨会话快照（锚点、四层队列、§八纪律 A–V、§九清单元规则与门禁坑 15 条） |
| `docs/t027/REMAINING_SLICES.md` | 289469 | 台账登记落点（状态表条目 + P54 小节） |
| `docs/t027/REMAINING_SLICES_EN.md` | 324042 | 台账英文镜像落点 |
| `docs/CONTRACTS.md` | 78534 | 本片核心交付对象：§7「发布耐久等级与升级条件」与 replay store 段落修订 |
| `docs/ADR_REGISTER.md` | 38813 | ADR-013（候选 C1–C6 枚举与 B3b 结论）与本片新增 ADR 的登记落点 |
| `docs/validation/t027-b3b-durability.md` | 26392 | B3b 结论账本（C2 候选级 `proven`（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证），决定格 [record-linked, parent-synced] 5/5 齐备；C1/C3/C2C3 `unresolved`） |
| `docs/validation/t027-windows-publication-durability.md` | 35799 | A7 反证原型（C6 预建槽位定义、D4 契约冲突、候选排序） |
| `docs/validation/t027b-tactical-blueprint.md` | 6315 | P38 施工图（其 C1–C7 为 T027-B 任务编号，C6 = DurableCostLedger 已由 P39 落地；注意与 ADR-013 的 C6 预建槽位是两套编号） |
| `internal/adapters/agent_runner_replay_store.go` | 29378 | 生产发布原语（`writeReplayRecordNoReplace`、`replayStorePublishRecordNoReplace`、`replayStoreSyncParent`）与耐久门禁 |
| `internal/adapters/agent_runner_replay_store_windows.go` | 201 | Windows 平台实现（当前 `PublishDurabilityUnproven` + 父目录步空操作） |
| `internal/adapters/agent_runner_replay_store_unix.go` | 322 | Unix 平台实现（`PublishDurabilityProven` + 父目录 fsync；本片零改动） |
| `internal/adapters/agent_runner_replay_store_test.go` | 32696 | 既有测试基线（`TestAgentRunnerReplayStorePublicationDurabilityMatchesPlatform` 等） |

## 前置核对要点（步骤 0 已阅）

1. `CONTRACTS.md` §7 现行口径：发布耐久等级三态 `proven`/`disproven`/`unresolved`，判定单元 =（平台 × 前提 P）；升级需同前提 (a) 可复现注入或对照证据、(b) 区分力门、(c) 机制证据三判据齐备；平台级等级只可表述 `unproven` 默认态或携带前提与证据；Unix「link 成功 + 父目录 fsync 成功」视为满足全部判据的既有契约口径实例。
2. Windows 现状：`agent_runner_replay_store_windows.go` 返回 `PublishDurabilityUnproven`，`syncReplayStoreParentDirectory` 为空操作；`RecordRequest`/`RecordCompletion` 在 `unproven` 下 fail-closed。
3. B3b 结论：C2（`MoveFileExW` + `MOVEFILE_WRITE_THROUGH` no-replace move）的 `proven` 为候选级（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证）；C6（预建槽位单文件）未入注入矩阵（A7 判「现行契约下不可即插」）。
4. **关键开放问题（交 ① 架构裁决）**：起跑消息的「C6」与「MoveFileEx」示例分属两套候选编号——ADR-013 的 C6 = 预建槽位（结构性候选，需新证据），B3b 的 C2 = MoveFileEx（已有候选级 `proven` 证据（仅限前提 P′；P′ 覆盖物理机目标部署环境未论证））。① 必须明确本片落地的机制编号与契约修订文本，红队复核。
