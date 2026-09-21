# 变异臂（装置自证伪）原始证据 — 2026-09-21

## 这是什么

本目录是 B3b 权威装置之外的**人工变异臂**（设计 §6 表中"变异臂（装置自证伪）"）的原始产物，
按用户裁决 **D1**（2026-09-21）「补跑变异臂 2 轮」执行。定位是**加强证据**：不阻塞主矩阵、不改写、
也不参与主矩阵结果。主矩阵 85 轮与 U4 18 轮全部来自主会话（本目录不计入其预算）。

## 怎么做

- 写入器以**编译期常量**降级：`b3bVariantDegrade`（`-tags "b3bnative b3bvariant"` ⇒ `true`）。
  降级把**所有候选的原语替换为基线**（`os.Link` + `platform-parent-sync`），
  而 **trace 如实记录降级后真正执行的原语**（这正是要检验的那一点）。
- 变体二进制：`writer-variant.exe`，sha256 `5b0d06d87547e98f295a7adedae9d5522503e615b0152340dbb5a1a4e2416633`，8,956,416 B。
- **默认构建不受影响**（同源代码、开关为 `false`）：`go test -tags b3bnative ./internal/adapters/` exit 0；
  `go vet -tags b3bnative ./internal/adapters/` ok。默认重建产物哈希 `880655ce…` 与原 `06de8a72…` 不同，
  原因仅是构建输入集合变化（build ID），**行为未变由原生套件与下述离线互证共同证明**。
- 被测格：**`c2:parent-synced`**（c2 的决定格；诚实 c2 在同一格为主矩阵 5/5 存活），2 轮，
  独立 rig root `b3b-runs-variant`（会话号同为 `2026-09-21`，与主会话目录互不相干）。

## 结果（逐轮）

| 轮 | 候选/阶段 | outcome | stageHit | mechanismOk | mechanismMissing | trace（目标阶段条目，逐字） |
| --- | --- | --- | --- | --- | --- | --- |
| r00 | c2 / parent-synced | `survived` | `true` | **`false`** | `MoveFileExW(WRITE_THROUGH,no-replace)` | `stage:temp-written stage:temp-synced stage:temp-closed os.Link stage:record-linked platform-parent-sync stage:parent-synced` |
| r01 | c2 / parent-synced | `survived` | `true` | **`false`** | `MoveFileExW(WRITE_THROUGH,no-replace)` | 同 r00 |

## 判据与结论

1. 两轮都是 `survived` **且** `stageHit=true`：**只看存活**的话，它们"看起来"与诚实 c2 的干净格无法区分。
2. 但两轮 `mechanismOk=false`（缺 `MoveFileExW(WRITE_THROUGH,no-replace)`），trace 逐字给出降级原语
   （`os.Link`、`platform-parent-sync`）。
3. 按 `Get-B3bCellVerdict` 的 fail-closed 分支：**决定阶段**的有效轮若 `mechanismOk=false`，
   该格**不得** `eligible`（理由记 `mechanism evidence missing in rounds …`）⇒ 无论存活与否都封顶 `unresolved`。
4. ⇒ **自证伪成立**：静默降级的候选无法把"存活"兑换成耐久主张；装置对"机制未执行"有区分力。
5. **离线互证**（无需 VM）：变体构建下 `TestB3bNativeCandidateTraceNamesTheRealPrimitive` **FAIL**，
   诊断逐字显示 `got … "os.Link" …` / `"platform-parent-sync"`；默认构建下同一测试 **PASS**（exit 0）
   ⇒ 降级属实，且开关不改变默认路径。

## 边界（不得读取的结论）

- 本臂**不得**被读作"c2 在 `parent-synced` 不耐久"：降级原语**不是** c2，它只证明装置的机制检查有效。
- 本臂**不得**被读作 c2 主格结论的一部分；c2 的主张只取主会话 `rounds/` 的 5/5 轮。
- 本臂只有 2 轮，**不构成**任何比例的覆盖率证据（与 §3.10 A11/A12 的口径一致：抽样只提供存在性佐证）。
- 可复算边界：本目录含 `rounds/`（2 轮完整工件）、`journal.jsonl`、`session.log`、`pin.json`；
  其校验和由外层证据包的 `SHA256SUMS.txt` 统一覆盖。
