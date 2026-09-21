# T027 · B3b · 候选注入矩阵与耐久三档定案 — 验证报告

日期：2026-09-21 ｜ 状态：**COMPLETE（报告成文；未提交、未推送）** ｜ 执行主控：DeepSeek V4.1 Flash（⑥ 段主控档位 = 深度/Max）
配套权威：`docs/t027/B3B_DESIGN.md`（设计与整改单一权威）｜ `docs/t027/B3A_RIG_DESIGN.md`（冻结装置）｜ 证据包：`docs/validation/evidence/b3b-2026-09-21/`

> **级别与前提必须一起读**：本报告只使用三态等级 `proven` / `unresolved` / `disproven`；
> 任何 `proven` 都**只在前提 P 下**成立，**不得**读作长期保证、普遍耐久、跨平台结论或 AT-23 证据。

---

## 0. 诚实准则（本报告所有结论的前提）

1. **前提 P 之外的任何推论都不成立**（§1）；结论不覆盖物理硬件、其它 hypervisor/存储栈、其它文件系统/载荷尺寸、宿主自身掉电。
2. **等级三态**：`proven` / `unresolved` / `disproven`。`unresolved` **不是**"失败"，也**不是**"不耐久"；它表示在前提 P 与装置分辨率下**未能建立主张**。
3. **`proven` 的五条件**（设计 §7）：① 全部 5 格 × N=5 有效轮落在允许集；② **决定格机制证据齐备**；③ 区分力门成立；④ 全局不变量成立；⑤ journal 全链与 plan/cut/reconcile 完整。
4. **测试不能自证**：自检与断言只证明"删掉守卫会变红"这类可证伪性（§6/§9），不构成结论本身；结论的唯一来源是 ⑤ 的原生证据。
5. **不得使用**长期保证类措辞（`durable` / `guarantee` / `survives power loss`）；只能用 §5 的 `claimScope` 句式逐字表达。

## 1. 范围与边界

**前提 P（冻结，逐字见 `matrix/verdict.json` 的 `premise`）**

- hypervisor **VirtualBox 7.2.18 r175117**；guest **Windows 11 25H2 build 26200.8037**（零售、未激活）；文件系统 **NTFS**。
- 存储：SATA/AHCI，**`useHostIOCache=false`**，两块盘 `--nonrotational on`；注入卷 **`R:` = INJECT**（`data.vdi`，7.98 GB，4096 B 群集）。
- 注入方式：**guest 硬断电**（`VBoxManage controlvm poweroff`）；基线快照 `pristine` = `065a25e1-d95c-4f33-8081-4f06216ecce8`（本片不重取）。
- `nic1=none`（客机无网络栈）。**信令域声明**：`C:\prfrail-prep\` 为阶段标记的第二写域，**不参与测量**；测量域只有 `R:` 上的 replay store 根，每轮仅**一次**被测发布。
- 载荷：记录级 **~1 KiB**（1082 B）；cut 前冻结延时 **2.0 s**（沿用 B3a 冻结标定，不重标定）；每格 **N=5**。

**不覆盖（不得外推）**

- 物理硬件/固件写回与真实掉电（本片是 hypervisor 级硬断电）；其它 hypervisor、其它文件系统、SMB/共享卷/云盘（U5）；宿主自身掉电（U3）。
- `.tmp` 残留清理、并发多写发布、resume/结算等其它链路。
- **AT-23**：本片**不构成** AT-23 证据，不解锁 Windows `first-dispatch`，不解除其它平台 `unproven`。
- 统计分辨率：格内 0 丢失的单侧 95% 上界 p<0.45 ⇒ 只能读作"该前提下、该装置分辨率内未观测到反例"。

**装置保真度（本片未证明，逐字保留）**：本片发现过一次"两次盘点可见性晚到"的现象（未链接临时件 size 0→1082，见 R33），
证据偏向"首次枚举读到了未定型元数据视图"的测量伪影，但**不能**仅凭该证据排除宿主/虚拟盘层在 cut 之后继续可见化的保真度问题。
因此本报告**不得**被读作"本装置已证明等价于物理断电"。

## 2. 候选与装置

| 候选 | 语义 | 落点 |
| --- | --- | --- |
| **C1**（参照臂） | 现行基线：`os.Link` no-replace + 平台父目录步（Windows 上为空操作） | 默认路径本身；**永不 `proven`** |
| **C2** | `MoveFileExW` + `MOVEFILE_WRITE_THROUGH (0x8)` 的 no-replace move 取代 `os.Link` | 发布原语钩子 `replayStorePublishPrimitiveHook` |
| **C3** | 目录句柄 `FlushFileBuffers` 取代 Windows 上的空操作父目录步 | 父目录钩子 `replayStoreParentSyncHook` |
| **C2C3** | C2 + C3 组合 | 两个钩子 |
| **positive-control** | 基线原语 + 目录/记录句柄 flush（无卷句柄权限时回退并单独标记） | 两个钩子；**不主张决定格**，只标定装置 |

- 两个接缝**均为 nil 默认、零语义**：生产路径逐字节不变（③/④ 均复核过）。
- **D3（门禁绕过，已声明）**：写入器为实验进程，内置 `publicationDurability=proven` 并直呼 `writeReplayRecordNoReplace`——只为到达被测发布点，不改变生产门禁。
- **变异臂（装置自证伪，人工臂，不计入 103 轮）**：编译期常量 `b3bVariantDegrade`（`-tags "b3bnative b3bvariant"`）把所有候选降级为基线原语，而 trace 如实记录降级后的原语。见 §5.4。

**冻结协议（设计 §2/§4，B3a 复用）**

- 5 个发布阶段：`temp-written` → `temp-synced` → `temp-closed` → `record-linked` → `parent-synced`。
- **决定格**：C2 = {record-linked, parent-synced}；C3 = {parent-synced}；C2C3 = {parent-synced}；C1/正对照 = ∅。
- **阶段归属以 cut 前归档的 `marker-transcript.txt` 为主**；逃逸 ⇒ `sessionBlock`（finalize exit 6）。
- **fail-closed 规则（本片判定的关键）**：决定阶段的有效轮若 `mechanismOk=false`，该格**不得** `eligible`，理由记 `mechanism evidence missing in rounds N`。
  含义：**机制未被证明执行过的"存活"不能兑现成耐久主张**——这是本片装置自证伪能力的正式落点。

## 3. 执行结果

### 3.1 门禁（Windows 本机，提交前终跑）

| 项 | 结果 |
| --- | --- |
| 装置自检 `tools/b3b-rig/Test-B3bRig.ps1` | **`checks=103 failures=0`** |
| `Invoke-ScriptAnalyzer -Path tools/b3b-rig -Recurse` | **0** |
| `gofmt -l`（触碰文件 + `internal`） | 空 |
| `go build ./...` / `go vet ./...` / `go vet -tags b3bnative ./internal/adapters/` | 全 ok |
| `go test -count=1 ./...` | **14 包全 ok、0 FAIL** |
| 原生套件 `go test -tags b3bnative -count=1 ./internal/adapters/` | **exit 0**（25 s） |
| 编码与行尾 | 本片触碰的 `.md`/`.ps1` = UTF-8 **with BOM** + LF；`.go`/`.json` = **without BOM** + LF |
| **GitHub Actions（提交后观察）** | run **`35562559401`**（head `abfcecc`）= **success**：Ubuntu 腿 Build/Vet/Test/**Race**/**Contract fixtures** 逐步全 success；Windows 腿全 success；`workflow_dispatch` 专用作业（Candidate build/probe、Bootstrap）按设计 skipped |

### 3.2 主矩阵（会话 `2026-09-21`，写入器 sha256 `06de8a72…`）

- **17 格 × 5 = 85 轮，全部有效：`ok=85 void=0`**；运行失败 0（唯一次运行失败 r07 发生在修复前的旧代码上，已由 R34 修复后清零，见 §9.4）。
- 标定/负对照（判别力门）：正对照 `positive-control:parent-synced` **5/5 存活**；负对照 `c1:temp-written` **5/5 丢失**（另 `c1:temp-synced`、`c1:temp-closed` 各 5/5 丢失）。
- 逐格结果（存活/丢失）：

| 格 | survived | lost | 机制令牌齐备 |
| --- | --- | --- | --- |
| `positive-control:parent-synced` | 5 | 0 | 5/5 |
| `c1:temp-written` / `temp-synced` / `temp-closed` | 0 / 0 / 0 | 5 / 5 / 5 | — |
| `c1:record-linked` | 4 | **1** | 5/5 |
| `c1:parent-synced` | 5 | 0 | 5/5 |
| `c2:temp-written` / `temp-synced` / `temp-closed` | 0 / 0 / 0 | 5 / 5 / 5 | — |
| `c2:record-linked` | **5** | 0 | **5/5** |
| `c2:parent-synced` | **5** | 0 | **5/5** |
| `c3:temp-written` / `temp-synced` / `temp-closed` | 0 / 0 / 0 | 5 / 5 / 5 | — |
| `c3:record-linked` | 5 | 0 | 0/5（该格非其决定格） |
| `c3:parent-synced` | 5 | 0 | **4/5**（见 §5.2） |
| `c2c3:parent-synced` | 5 | 0 | **4/5**（见 §5.2） |

- **盘点读取次数实测**：`reads=2` 或 `reads=3`，**未出现 4 读**（R33 的额外代价平均 1 次盘点 ≈27 s/轮）。

### 3.3 U4 排序子矩阵（`-Paired`，轮号 ≥500）

- **6 格 × 3 = 18 轮，全部有效：`ok=18 void=0`**；18 轮 `targetKind` 均为 `completion`（finalize 的强制项）。
- **真孤儿 C 计数 = 0**（全 6 格）。

### 3.4 变异臂（D1，人工臂，独立 rig root `b3b-runs-variant`）

- 2 轮 `c2:parent-synced`：`survived` **且** `stageHit=true`，但 `mechanismOk=false`（缺 `MoveFileExW(WRITE_THROUGH,no-replace)`），
  trace 逐字为 `… os.Link … platform-parent-sync …` ⇒ 装置**拒绝**把这两轮"存活"兑现为机制成立的证据。
- 离线互证：变体构建下 `TestB3bNativeCandidateTraceNamesTheRealPrimitive` **FAIL**（诊断逐字给出降级原语）；默认构建下同一测试 **PASS**、原生套件 exit 0。
- 产物并入证据包 `variant/`（含自撰 `README.md`）。

### 3.5 证据包

`docs/validation/evidence/b3b-2026-09-21`：`files=2157 rigMismatches=0 missingRequired=0`；
`optionalSectionsAbsent=calibration,u4,guest`（**偏离 D4**：本装置所有格——含正对照与 U4 配对轮——均写入 `rounds/`，配对轮以 `paired=true` 标记；**不得**把可选分区缺席读成"标定/U4 未执行"）。

## 4. 结论账本

### 4.1 三档定案（`-Finalize` exit 0，2026-09-21 12:08:21；数据源 `matrix/verdict.json`）

| 臂 | 等级 | 依据（机器输出） |
| --- | --- | --- |
| **C2** | **`proven`** | `decisionSets=[record-linked,parent-synced]`，两格**全 `eligible`**（5/5 在允许集且机制证据齐备） |
| C1 | `unresolved` | 参照臂无决定格（`decisionSets=[]`）；`infoViolations=["record-linked"]`（该格 `disproven`，见 4.2） |
| C3 | `unresolved` | `decisionSets=[parent-synced]` → `notEligible`，理由逐字 **`mechanism evidence missing in rounds 77`** |
| C2C3 | `unresolved` | 同上，理由逐字 **`mechanism evidence missing in rounds 81`** |
| positive-control | （标定） | 5/5 存活，`eligible`；不主张耐久 |

**判别力门（防假 `proven`）**：`ok=true`，负对照（c1 阶段 1）**5 轮 5 丢失**、正对照 **5 轮 0 丢失**（`requiredN=5`）。

**U4（排序）**：`unresolved`，`validRounds=18`，`orphans=0`，注记逐字：
`ordering can only be falsified by a counterexample, never proven by a clean sample`。

**结论口径（`claimScope`，逐字引用）**：

> within premise P (VirtualBox Win11B3 / NTFS on SATA-AHCI, useHostIOCache=false, record-level ~1 KiB payload, 2 s cut window, N=5 per cell): no counterexample observed, candidate mechanism demonstrably executed, device discrimination demonstrated at stage temp-written only; not a long-term guarantee and not a claim about other hypervisors, filesystems, payload sizes or the host losing power

**必须逐句声明的限定**：

- (a) 正对照实际执行的是**卷/记录句柄 flush**（trace 逐字记录，无卷句柄权限时回退并单独标记）；区分力门只能读作"在该次刷盘下装置可观测到干净结果"，**不得**读作"卷级设备缓存已证明被排空"。
- (b) `c2c3` 与 `c3` 在 `parent-synced` 同一决定点、该阶段组合原语与 C3 等价 ⇒ `c2c3` 只能作**诊断臂/重复臂**报告，**不得**当作独立的耐久主张。
- (c) 阶段证明取自 **cut 前**转录，而非断电后拷回件；`C1` 在阶段 4 的 1/5 丢失**只**作为装置区分力的证据，不外推为生产行为结论。

### 4.2 逐格拒绝/证伪项原文

| 格 | 等级 | 理由（机器输出逐字） |
| --- | --- | --- |
| `c1/record-linked` | `disproven` | `round 21: lost-absent is outside the allowed set for record-linked`（该轮 `stageHit && mechanismOk`，故可归属为反例） |
| `c3/parent-synced` | `unresolved` | `mechanism evidence missing in rounds 77` |
| `c2c3/parent-synced` | `unresolved` | `mechanism evidence missing in rounds 81` |

## 5. 关键判定说明（防误读）

### 5.1 为什么 C2 可以 `proven` 而 C3/C2C3 不能

C2 的两个决定格共 10 轮全部落在允许集、机制令牌齐备、装置有区分力、journal 全链完整；
C3/C2C3 的决定格各有 **1/5 轮**断电落在机制真正执行之前（trace 中无 `FlushFileBuffers(dir)` 等令牌，
见 r77/r81）。装置按 fail-closed **拒绝**把"机制未执行过的存活"算作该臂的证据——
若算作证据，就等于让候选认领它并未产生的耐久性。

### 5.2 由此产生的**不得读取**结论

- **不得**写"c3/c2c3 不耐久"或"c3/c2c3 被证伪"：它们只是**未能建立主张**（`unresolved`）。
- **不得**写"U4 已证明有序"：0 孤儿只能"未被证伪"。
- **不得**写"候选在决定格优于基线"以外的因果（除 C1 阶段 4 的 1/5 丢失这一装置区分力证据外）。
- **不得**声称"journal 链可检出轮工件篡改"：链只验 journal 行；工件完整性靠 plan+inventory 重派生与逐文件哈希比对；对工件与 journal 的**整体重写**只能靠证据包校验和与 S0 pin 这两个外部锚点发现。
- **不得**声称"证据包自动可复算"：只能说可由 `rounds/` + `journal.jsonl` 复算。

## 6. 改动点与测试点（可证伪性）

- **装置**：`tools/b3b-rig/` 共 10 个脚本（新写）+ 证据包构建器；`tools/b3a-rig/` 4 个文件**原样复用**（会话 pin 记录其哈希）。
- **Go 侧**：`internal/adapters/agent_runner_replay_store.go`（两个 nil 默认接缝 + 派发器）、
  `agent_runner_replay_store_b3b_test.go`、`agent_runner_b3b_candidates_windows_test.go`、
  `agent_runner_native_b3b_windows_test.go`、变异臂两文件（`agent_runner_b3b_variant_default_test.go` / `_degraded_test.go`）。
- **断言总数 103**（`checks=103 failures=0`），其中 6 个篡改拒裁夹具 + 1 个截断容忍夹具 + 6 条入库证据包断言 + 3 条变异臂接线断言。
- **反证实验（删守卫/回退逻辑 ⇒ 断言必红）**，全部实际执行：

| # | 变异 | 期望红灯 | 实测 |
| --- | --- | --- | --- |
| 1 | `-Cell` 解析正则退回不支持连字符 | `control-arm-cell-is-schedulable` | **红** ✓ |
| 2 | 盘点稳定对退回"只看首对" | `inventory-settles-on-the-second-pair` | **红** ✓ |
| 3 | finalize 空会话守卫禁用（`-eq -1`） | `finalize-refuses-a-session-with-no-rounds` | **红** ✓ |
| 4 | 去掉预检取锁调用（**A11 抽查项**） | `preflight-takes-the-session-lock` | **红** ✓ |
| 5 | 变体构建跑真原语断言 | `TestB3bNativeCandidateTraceNamesTheRealPrimitive` | **红** ✓（默认构建同测试绿） |

## 7. 已知边界与后续

- **装置区分力风险**（设计 §10.1）：记录级 ~1 KiB 载荷下阶段 4 只被证伪 1/5；更大载荷/更窄窗口未见丢失时的主张上限仍是 `unresolved`。
- **阶段 4/5 在盘点观察上等价**（A7 CP4≡CP5）⇒ 标记通道失效即只能 `unresolved`（本片 r77/r81 即该类的实例）。
- **保真度未证明**（§1 末）：不得声称装置等价物理断电；如需该主张必须另有独立检验（本片未做）。
- **时间口径**：单轮实测 ≈5.3–7.3 min（非设计早期估算的 3.8–4.6 min）；按每轮 `inventory-reads.txt` 实测分布回算，不得沿用"两次盘点"旧口径。
- **D4**：证据包 `calibration/`・`u4/` 分区缺席（见 §3.5）。
- **CI 抖动观察（一次）**：docs 提交 `d2f508c` 的 run `35563800665` 首跑 Ubuntu `Test` 步骤失败于 `tools/agent-probe/enforcement-proxy` 的 `TestConnectAuditTrailFailureDeniesAndDoesNotTunnel`（`AUDIT TRAIL BROKEN (write …/proxy.jsonl: file already closed)`；`main_test.go:125: target must have been reached while the audit trail was healthy`）；**同一代码**在前一 run `35562559401` 该包为 success，复跑 `35563800665` 双腿全绿，docs 回写提交 `dd36b9b` 的 run `35564805991` 双腿全绿且该包未复现 ⇒ 判定**偶发**、指向该测试自身句柄/时序竞态；**本片不跨切片修改 B2 期工具**（已登记为 **DR-2**，见切片清单的已知缺陷表）。
- 后续：`proven` 若要用到产品契约，需独立的契约修订切片；Windows `unproven`、`first-dispatch` 拒绝、AT-23 结论**均未改变**。
- **台账缺口修复状态（2026-09-21，用户授权后）**：`REMAINING_SLICES{,_EN}.md` 完成记录表的 **A5/A6/A7/B2/B3a 五行已补录**——依据各切片**自身验证报告**、逐行标注来源、仅补可核验字段（日期/提交号/运行号/证据链接/一句话摘要），无法核验者标注 **历史缺失**（如 B2 的 CI 运行号）；补录后双语行数对称（11/11）。同轮修复了我自己引入的一处缺陷：**EN 台账的 A4 行曾在补录编辑中被误删**（oldString 含 A4、newString 未回填），已回填并复核——该缺陷由本轮自检发现并在同一轮修复，**如实留痕**。**仍待回填**：`DEV_PLAN{,_EN}.md` 的 B3a 段落（B3b 段已补）。

## 8. 成本与计量

- **探针额度：0/10 使用**（本片全离线，无付费探针、无真实模型调用）。
- **⑤ 墙钟**：`01:25:20`（主会话 S0）→ `12:08:21`（finalize）= **10 h 43 min**；证据包 `12:09:13`。
  其中：标定 5 轮 26 min；c1（含 1 次运行失败整轮重试 ×2）61 min；主矩阵其余 14 格 ≈7 h 22 min；U4 18 轮 ≈1 h 35 min；变异臂（含变体预检）≈12 min。
- **③ V4 Pro**：① 前置分析 ×1；③ 预审 ×1（`changes-required`，F1/F2/F3 阻断级 → R1–R12）；收口复审 **PASS**（轮次口径见设计说明 §13 表）；
  另有 **1 次调用因提供方连接失败**（不消耗额度，按 §3.9 可用性预检口径登记）。
- **③.5 MAI-Code-1.1-Flash**：首轮 ×1（M1/M2）→ 整改（R13/R14）→ 收口 **PASS**（合计 2 次，未超 §3.10 的 ≤3）。
- **④ GPT-5.3 Codex：6 次**（**超 §2.3 的 1+1 上限，按 D2 裁决如实记账**）。逐次目的与发现数：

| # | 目的 | 结论 | 发现 |
| --- | --- | --- | --- |
| 1 | 首轮终审（不附任何清单） | `FINDINGS` | 4 项（C1–C4） |
| 2 | 对 C1–C4 整改的复审 | `FINDINGS` | 5 项（MAJOR#1–#5） |
| 3 | 收口 #1（含 P1–P4 与 5 项整改） | `FINDINGS` | 5/5 历史项关闭 + **1 项新 MAJOR（R29）** + 1 处覆盖缺口（R30） |
| 4 | 收口 #2（R29/R30 整改后） | **`PASS`（明确 GO）** | 0 |
| 5 | 定向 delta 复审（P5：`-Cell` 解析器） | **`PASS`** | 0（附 1 处可证伪性补强建议 → 已补断言） |
| 6 | 定向 delta 复审（P7：盘点稳定对） | **`PASS`** | 0（附 1 处证据落盘建议 → 已补 `inventory-mismatch-*`） |

  **超预算原因（非执行者违规）**：用户明确要求"整改后必须重跑本级复审"（闭环铁律），
  且 ⑤ 期间新出现两个审查点（P5 控制臂不可调度、P7 盘点稳定对）——每次整改后按铁律重跑 ④ 是**义务**而非选择。
- **token**：由用户从用量 UI 回填后计入本行（主控不自行估计）。
- **门禁墙钟**：自检 ≈20 s/次（本片累计约 10 次）；原生套件 25 s；`go test ./...` ≈40 s。
- **双语一致性检查（§四.10）**：`REMAINING_SLICES{,_EN}.md` 的 B3b 状态行/勾选项/完成记录行/标题符号两侧**同构**（各 13 处 `B3b` 标记，两侧均 `✅`）；`ADR_REGISTER{,_EN}.md` 的 ADR-013 行与 `DEV_PLAN{,_EN}.md` 的 B3b 段落同样双侧镜像；本报告与 `_EN` 的数字/命令/结论逐项一致。**不一致项：0**。

## 9. 审查记录（③ / ③.5 / ④）与收口闭环

### 9.1 ③ V4 Pro

- 预审（实现完成后）：**`changes-required`**（F1/F2/F3 为阻断级）→ 整改 R1–R12 → 收口复审 **`PASS`**。
- 预审侧重点：决定格与机制令牌自相矛盾（R1）、正对照无驱动器致区分力门结构性永假（R2）、
  正对照 flush 语义过强（R3）、阶段归属读 cut 后拷回件（R4）等，均已闭环并落到断言。

### 9.2 ③.5 MAI（独立扫描，单层）

- 首轮 **`INDEPENDENT SCAN: FINDINGS`**：2 项（M1/M2，见设计说明 §13）→ 整改 R13/R14 → 收口 **`PASS`**。
- **A12 独立佐证**：MAI 报 PASS，故 ⑤ 必须至少一项专项实验独立证明关键不变量 ⇒ 本片以**变异臂**（装置自证伪，§3.4）
  承担该角色：它对"机制证据是否真的在把关"给出了与 MAI 无关的独立判据（降级原语下装置拒绝兑现，且离线可复现）。

### 9.3 ④ Codex（6 次，逐次见 §8）

- 审查覆盖：安全（绕过/注入/零值边界）、架构一致性（切片目标与契约）、完整性（正例/反例/边界）、
  **测试反例审查**（薄弱测试、只 happy path、未覆盖竞态）、可证伪性审计（Section D）。
- **A11 抽查（对称抽样，实际执行）**：
  - ④ 声明「去掉预检取锁调用 ⇒ `preflight-takes-the-session-lock` 红」→ **实测红** ✓（本报告 §6 表 #4）。
  - ④ 声明「盘点稳定对退回只看首对 ⇒ `inventory-settles-on-the-second-pair` 红」→ **实测红** ✓（§6 表 #2）。
  - **缺位声明**：③.5 MAI 的收口轮为 `PASS`，未给出可执行的"删守卫→变红"条目，故抽样无法落在 MAI 侧；
    记为已知局限（抽样只提供**存在性佐证**，不构成覆盖率证明）。
- **收口闭环记录**（用户质询后补做）：

| 级 | 收口结果 | 备注 |
| --- | --- | --- |
| ③ V4 Pro | **PASS** | 无新实质发现 |
| ③.5 MAI | **PASS** | — |
| ④ 收口 #1 | `FINDINGS` | 5/5 历史项关闭 + R29 + R30 |
| ④ 收口 #2 | **PASS（明确 GO）** | — |
| ④ delta（P5） | **PASS** | 附可证伪性补强建议，已采纳 |
| ④ delta（P7） | **PASS** | 附证据落盘建议，已采纳 |

- **审查结论落地校验（§3.11）**：三轮审查给出的 `file:line` 均在落地前逐条定位复核；无法定位者按线索处理（本片无此类条目）。

### 9.4 ⑤ 就地修复（实机运行中才发现，逐项已入设计说明 R31–R35 / P5–P8）

| # | 现象 | 整改 | 反证 |
| --- | --- | --- | --- |
| P1 | 硬断电后残留 `VBoxHeadless` 会话致 `startvm` 报 `E_FAIL` | `Wait-B3bVmPowerOff` + `Clear-B3bStaleVmSession`（词边界匹配 + fail-closed）+ 有界重试 | 实机 `startvm ok=True attempts=1` |
| P2 | 装置无并发保护（两实例同操作同一 VM） | 会话锁（pid + owner + `sessionDir` 指纹；仅"活 pid 且同目录"才拒绝） | 双断言（阻挡/跨目录拷贝不阻挡） |
| P3 | 冻结写入器拒绝空串 ⇒ 断电后未捕获异常 | `Write-B3bArtifactText`（`[AllowEmptyString]`）替换所有可为空的写入点 | 实机单轮 exit 0 |
| P4 | 矩阵丢弃非预期退出码的轮 stderr | 记录轮 stderr（该补丁直接定位了 P3） | — |
| P5 | 控制臂名含连字符 ⇒ `-Cell` 校验正则拒绝 ⇒ 正对照标定格**无法从 CLI 调度** | 抽出共享解析器 `Resolve-B3bCellSpec` + 行为式断言 | 回退正则 ⇒ 断言红 |
| P6 | 会话号缺省回落"当天日期" ⇒ `-Finalize` 漏传时会**零轮 exit 0 写空判定** | finalize 零轮 fail-closed + 自检夹具会话号显式贯通 | 禁用守卫 ⇒ 断言红 |
| P7 | 盘点契约实现过严（要求"头两次"一致，而本意是"枚举态已定"） | `Get-B3bSettledInventoryPair`（向前扫描首个相邻一致对，上界 4 读）+ 每次读取存档 + 稳定对入证据 | 退回"只看首对" ⇒ 断言红 |
| P8 | 单次盘点读取失败即整轮作废并重跑整轮 | 轮内重试 + 退出码/stderr 入日志 | 三条静态断言 |
| P9 | 变异臂文件名尾段 `arm` 被 Go 当作隐式 `GOARCH=arm` ⇒ 该文件**永不参与构建** | 改名 `…_degraded_test.go` + 用 `IgnoredGoFiles` 核验归属 | `go list` 两视图核验；坑位记入仓库记忆 |

### 9.5 残差-抽查（§3.11）

**本片触发必抽条件**：本片改动涉及**停机语义**（硬断电注入、断电后重启与盘点）与装置所有权语义（会话锁）。
按 §3.11「必抽」要求，应在 ④ 之前追加 1 次 Codex 盲审。**实际执行方式**：本片 ④ 的**首轮调用即不附任何 ③/③.5 清单**（§3.10 的锚定硬隔离），
即为一次事实上的盲审，故**未重复追加**第二次盲审——该合并判断按 §3.11 要求在此显式记录理由：
两次盲审的输入集合（切片定义 + diff + 契约原句 + 只读约束）完全同构，重复调用只增加成本而不增加判别力。
**盲审发现数**：4 项（C1–C4）；**是否推翻 ③.5 的 PASS**：否（两者发现集合交集为空，均由后续收口闭环关闭）。

## 10. 明确未执行事项

- **已 `commit` 并 `push`（2026-09-21，用户同轮显式授权）**：提交 **`abfcecc`**（精确暂存 2183 文件），推送 `origin/main`（`74cbd68..abfcecc`），**未推 gitee**；CI run **`35562559401`** 双腿全绿（Ubuntu 含 Race/Contract fixtures、Windows）。以下各项仍未执行。
- **未做真实模型调用 / 付费 probe**（本片全离线）。
- **未做装置保真度的独立检验**（即未用外部手段比对 cut 前后虚拟盘字节，故 §1 末的限定成立）。
- **未跑长档压力**（如 200 轮同刻竞态）：本片固定 N=5/格 + U4 3 轮/格，符合设计预算。
- 变异臂是**人工臂**，不属 103 轮预算，也未纳入任何臂的三档定案。

## 11. 下一步建议（§7.4）

> 下一片按依赖图应为 **B4 · Windows T027 端到端（AT-23 前置）**（或在契约层先行处理 B3b 的 `proven` 是否被采纳）；
> 建议档位——主控 **深度（Max）**（仅设计与结论成文步骤；实现步骤可回标准）、
> V4 Pro **①③ 必调**（B4 触及停机/身份/所有权语义）、③.5 **启用**（行为变更切片）、Codex **Extra High（④ 必调 + 按 §3.11 必抽盲审）**。
> 需要启动时告诉我，我会按 7 步协议 + 闭环纪律执行（先 V4 Pro 前置分析，协议先行）。

- **准则待议（D2 裁决要求登记，本轮不改准则）**：硬门切片（A7/B2/B3/B4）的 ④ 调用上限是否需要**单列**？
  本片实测 6 次（1 终审 + 1 复审 + 2 收口 + 2 delta），其中 4 次源于"整改后必须重跑本级复审"的闭环铁律。
  建议在下一轮启动前评估：为硬门切片显式放宽上限（如 ≤4），或在准则中把"闭环重跑"明确为**不占上限**。
