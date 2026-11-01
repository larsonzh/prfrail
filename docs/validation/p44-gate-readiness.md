# ProofRail 验证报告 — P44-T027B-GATE-READINESS

| 字段 | 内容 |
|---|---|
| 切片 | `P44-T027B-GATE-READINESS` |
| 类型 | 实现类 `[SLICE]` |
| 准则版本 | v1.60 |
| 基线 HEAD | `a87b078` |
| 日期 | 2026-11-01 |

## 目标与边界

- **目标**：闭合 C7 与 N4——readiness 严格分离探针就绪、真实验收就绪与生产发布就绪三级，探针结论不得推进验收；`ProbePermit` 一次性消费，调用前原子预留，消费后失效，失败、超时与重启均不重置。
- **口径**：不得宣称 C7/N4 端到端闭合，只能宣称门控就绪与许可路径已满足当前切片验收前提。
- **收口形态**：`internal/gates` 新增 `readiness.go`（`ReadinessGoal` 三值 + sealed `Readiness` + 三个具体类型 + `EvaluateT027BReadiness`）与 `probe_permit.go`（`ProbePermit` 三态机 + `Require` / `Complete` / `LoadProbePermit`）。
- **范围**：`internal/gates/**`。
- **边界**：`internal/console/**` 零改动（未端到端接线，登记 R-3）；不碰 `internal/chain/**`、`internal/adapters/**`、schema、fixtures 与 `ci.yml`；不推 gitee；`commit` / `push` 停 ⑩ 等授权。

## 审查链结论

- **① 架构**：`deep-reasoner` + DeepSeek V4 Pro；《C7+N4 门控就绪与一次性许可方案》，`DESIGN: DONE`；§八 四条修订注记 R-1～R-4 为 ⑤ r1 裁定后补记。
- **③ 红测**：`RED-CONFIRMED`；变异验证 11 项（M1–M11）全部 RED，原始输出归档于 `03-tests/mutation-evidence.txt`。
- **② 实现**：落地 `internal/gates` 四个文件（两实现 + 两测试）；本包 P44 用例 20 例 + 9 子例。
- **④ 门禁**：`TOTAL_FAIL=0`；读数见 §门禁读数。
- **⑤ r1**：`PRE-REVIEW: PASS WITH FIXES`（2 Medium：规格漂移未记录 / 无生产接线；另 3 Low + 4 Info）。
- **⑤ r2**：`PRE-REVIEW: PASS WITH FIXES`（仅 3 Low，r1 各项逐条闭合）。
- **⑥**：`INDEPENDENT SCAN: PASS`（Section B 零缺陷；Section D 列 5 处未覆盖面），本片额度 1/3。
- **⑦ r1**：`FINAL REVIEW: BLOCKED`（FR-01 High 崩溃窗口 / FR-02 Medium 缺红测 / FR-03 Low 证据未落盘）。
- **⑦ r2**：`BLOCKED`（NR-01 Medium MANIFEST 三方不自洽 / NR-02 Low D9 未登记）。
- **⑦ r3**：`FINAL REVIEW: PASS`（0 新增缺陷），本片额度 3/3。
- **R / D**：R = 0、D = 0（不继承上片）。

## 六段扫描转录

来源 `docs/validation/evidence/P44-T027B-GATE-READINESS/06-scan/review.md`；围栏内逐字转录、零改写，范围自 `## A 结论` 起、至结论行止。

```sixstage
## Section A 结论

INDEPENDENT SCAN: PASS

## Section B 缺陷

| 编号 | 严重度 | file:line | 缺陷描述 | 修法 |
|---|---|---|---|---|
|  |  |  |  |  |

## Section C 核验结论

1. 成立 — C7 三级不可混用，且包外无法实现 `Readiness`。`Readiness` 是 sealed interface，只有同包的私有 tag 方法才允许满足；`GoalProbe`/`GoalAcceptance`/`GoalProduction` 只能通过 `EvaluateT027BReadiness` 返回，见 [internal/gates/readiness.go](internal/gates/readiness.go#L69-L212)。`GoalAcceptance` 会在 `Seal == nil` 时先拒绝探针证据，`GoalProduction` 会在 `Seal == nil` 时直接返回 `ErrReadinessAcceptanceMissing`，而不看 `ProbeEvidence`，这正是 R-1 的 fail-closed 行为，见 [internal/gates/readiness.go](internal/gates/readiness.go#L173-L212)。`ReadinessInput` 以哈希串和 `Seal` 绑定表达 L2/L3，未形成跨层类型别名；真正的隔离由 sealed interface + 三个方法集差异保持，见 [internal/gates/readiness.go](internal/gates/readiness.go#L81-L138) 与 [internal/gates/readiness_test.go](internal/gates/readiness_test.go#L20-L48)。

2. 成立 — N4 是严格一次性消费，且 `Require`/`Complete` 不存在“重置额度”路径。`Require` 仅可从 `Unused -> Claimed` 走 CAS；任何失败、超时、重启、O_EXCL 冲突都落到 `ErrProbePermitInUse`/`ErrProbePermitConsumed`，并把状态推进到 `Spent`，见 [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L129-L203)。`Complete` 仅能从 `Claimed -> Spent` 或 `Spent -> Spent`，且拒绝未要求过的 `Unused`，没有“Spent 复活”路径，见 [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L216-L243)。`LoadProbePermit` 对“缺失/存在/损坏”做 fail-closed：缺失返回 `Unused`；存在/损坏返回 `Spent`，不允许发放第二次配额，见 [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L129-L171)。空 `path`、`ctx` 取消、并发 CAS 失败、`O_EXCL` 冲突四类边界在实现中均有自洽语义，见 [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L184-L285)。

3. 成立 — `EvaluateT027BReadiness` 是纯函数。它只读取 `goal` 和 `input`，并在 `GoalAcceptance`/`GoalProduction` 调用 `SealBinding.validate` 与正则校验，没有任何文件 I/O、没有写状态、没有消费 permit，见 [internal/gates/readiness.go](internal/gates/readiness.go#L157-L212)。测试也做了目录快照对比，断言评估前后无副作用，见 [internal/gates/readiness_test.go](internal/gates/readiness_test.go#L171-L214)。这也验证了 R-2 对“纯函数/哈希串签名”的实现方向是成立的。

4. 成立 — 哨兵完整性与 `errors.Is` 可达。`ErrProbePermitMissing`、`ErrProbePermitInUse`、`ErrProbePermitConsumed` 在 [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L55-L67) 定义；`ErrReadinessProbeNotAcceptance`、`ErrReadinessEvidenceIncomplete`、`ErrReadinessAcceptanceMissing`、`ErrReadinessProductionAT23Required` 在 [internal/gates/readiness.go](internal/gates/readiness.go#L19-L49) 定义。`GoalAcceptance` 和 `GoalProduction` 的返回都用 `%w` 包装，未出现“把 permit 错误包成 readiness 错误”或反向错误，见 [internal/gates/readiness.go](internal/gates/readiness.go#L173-L212)。

5. 成立 — 测试强度是有效的，且没有“镜像实现”式的无意义断言。编译期 pin 在 [internal/gates/readiness_test.go](internal/gates/readiness_test.go#L20-L48) 明确锁定了 sealed interface 与三类访问器；`TestReadinessLevelsCannotBeMixed` 仍有一个运行期断言，但仅用于正向 L1 分支的校验，见 [internal/gates/readiness_test.go](internal/gates/readiness_test.go#L108-L140)。它不构成跨级混用通道；真正的“不可混用”由编译期 pin 保证。该测试文件中也没有把实现逻辑复制进断言。

6. 成立 — 边界和污染控制成立。`git status --short --untracked-files=all` 只枚举了本 slice 允许文件：四个 gates 实现/测试文件、两个台账文件以及证据目录中的载荷。未出现 any sibling package drift；没有出现 [internal](internal) 之外的越界改动。允许集合包括 [internal/gates/readiness.go](internal/gates/readiness.go)、[internal/gates/probe_permit.go](internal/gates/probe_permit.go)、[internal/gates/readiness_test.go](internal/gates/readiness_test.go)、[internal/gates/probe_permit_test.go](internal/gates/probe_permit_test.go)、[docs/validation/evidence/P44-T027B-GATE-READINESS/00-briefing.md](docs/validation/evidence/P44-T027B-GATE-READINESS/00-briefing.md)、[docs/validation/evidence/P44-T027B-GATE-READINESS/freedom-list.md](docs/validation/evidence/P44-T027B-GATE-READINESS/freedom-list.md)、[docs/validation/evidence/P44-T027B-GATE-READINESS/MANIFEST.md](docs/validation/evidence/P44-T027B-GATE-READINESS/MANIFEST.md)、[docs/validation/evidence/P44-T027B-GATE-READINESS/SHA256SUMS.txt](docs/validation/evidence/P44-T027B-GATE-READINESS/SHA256SUMS.txt)、[docs/validation/evidence/P44-T027B-GATE-READINESS/01-arch-design/review.md](docs/validation/evidence/P44-T027B-GATE-READINESS/01-arch-design/review.md)、[docs/t027/REMAINING_SLICES.md](docs/t027/REMAINING_SLICES.md)、[docs/t027/REMAINING_SLICES_EN.md](docs/t027/REMAINING_SLICES_EN.md)。R-3 的“console 接线后置”结论也成立：本 slice 未宣称端到端闭合，且边界守住了“只读实现 + 后续装配切片”，见 [docs/validation/evidence/P44-T027B-GATE-READINESS/01-arch-design/review.md](docs/validation/evidence/P44-T027B-GATE-READINESS/01-arch-design/review.md#L131-L196)。

7. 成立 — 编码硬规则满足：四个新 Go 文件均为 UTF-8 无 BOM + LF。已直接检查文件字节，结论为 `BOM=False`、`LFonly=True`、`CRLF=False`，且 `go test ./internal/gates -count=20 -cpu 1,2,4` 返回 `ok ... 25.351s`。这覆盖了本 slice 的关键并发/语义验证；R-4 仍是环境限制，未在本机执行 `-race`，但不是实现缺陷。

## Section D 未覆盖面

- `SealBinding.validate` 的分支缺少单元覆盖：`Attempt < 1`、`TaskKey` 为空、`SealHash` 非法值的三条边界都未被单独断言；这些是 L2 质量门禁的重要失败路径。
- `LoadProbePermit` 对 `os.ReadFile` 的非 `NotExist` 失败路径未覆盖，例如权限错误或目录路径；这会决定“读取失败时是否继续返回 permit”。
- `Require` 的 `ctx == nil` 分支和 `Complete` 的“未知状态”分支尚未被直接测试；这属于强制 fail-closed 的边界条件。
- `GoalProduction` 对非空但非法格式的 `AT23Hash`（如 `abc` 或 63/65 位）未覆盖；其正则拒绝与 `ErrReadinessEvidenceIncomplete` 构成 L3 门禁。
- `GoalAcceptance` 对非法 `SealBinding`（带值但 hash 非法）未单独覆盖；目前仅验证了“没有 seal”与“没有 at23”，未专门锁定 `validate` 的格式篡改门。

## Section E 可证伪性

本 PASS 结论可被以下反例推翻：

- 任何跨包实现 `Readiness` 的类型能通过编译，或 `GoalAcceptance`/`GoalProduction` 在缺 L2 时返回 nil error，说明 sealed interface 或 fail-closed 判定失效。
- `Require` 在失败、超时、重启后仍能重新返回有效 claim；或 `Complete` 能把 `Spent` 状态复活；或 `LoadProbePermit` 对损坏标记返回 `Unused`，则 N4 一次性约束被否定。
- `EvaluateT027BReadiness` 产生任何文件写入、调用 `os.*`，或 `SealBinding.validate`/hash regex 允许伪造值，说明纯函数性不成立。
- 任何包外代码能绕过 `GoalProduction` 的 `Seal == nil` 语义并拿到生产 readiness，也会推翻 R-1 与 C7 语义。

INDEPENDENT SCAN: PASS
```
- 转录口径：标题前缀归一为 `## Section A`～`## Section E`（原文为 `## A 结论`～`## E 可证伪性`），以满足 G7-f 的 A→B→C→D→E 判据；正文逐字未改。
- 判定词归一化：原文末行的 `§A INDEPENDENT SCAN: PASS` 在围栏内按块末单行 `INDEPENDENT SCAN: PASS` 转录，与原文 `## A 结论` 下的同名结论行取值相同。

## 项闭合

严重度列取审查报告原文口径，不加粗；⑦ 结论闭集为 `FINAL REVIEW: PASS`（r3，0 新增缺陷）。

| # | 项 | 严重度 | 位置 | 处置 |
|---|---|---|---|---|
| ⑤ r1 B1 | `GoalProduction` 缺 L2 材料的判据与 ① 字面不一致，属未记录的规格漂移 | Medium | `readiness.go` 的 `GoalProduction` | 已闭合：按 §八 R-1 固化，判据取 `Seal == nil`，不看 `ProbeEvidence` |
| ⑤ r1 B2 | C7/N4 只有纯 API，全仓无生产调用点，台账与工作树不符 | Medium | `internal/gates/**` 与台账 | 已闭合：按 R-3 登记接线后置，并明写不得宣称端到端闭合 |
| ⑤ r1 B3 | 包级注释声称每个返回错误恰包一枚哨兵，与 ctx 路径不符 | Low | `probe_permit.go` 包级注释 | 已闭合：注释改为 permit 状态错误恰包一枚哨兵，ctx 取消原样透传 |
| ⑤ r1 B4 | 两个负向类型断言为永真断言 | Low | `readiness_test.go` | 已闭合：永真断言删除，改由编译期方法集钉承担 |
| ⑤ r1 B5 | 缺三处负例：损坏标记、双实例竞争、重启后幂等 | Low | `probe_permit_test.go` | 已闭合：四例均已补入 |
| ⑤ r2 r2-1 | MANIFEST 元数据未闭环，SHA256SUMS 只覆盖两件载荷 | Low | 证据包元数据 | 已闭合：本次归档 `07-final/review.md` 并重算全部载荷哈希 |
| ⑤ r2 r2-2 | readiness 包级注释与 GoalProbe 分支返回的跨域哨兵不符 | Low | `readiness.go` 注释 | 已闭合：注释补例外句，GoalProbe 分支包 `ErrProbePermitMissing` |
| ⑤ r2 r2-3 | `Require` 持久化失败的包装与「原样透传」注释不符 | Low | `probe_permit.go` 注释 | 已闭合：注释限定持久化失败包 `ErrProbePermitConsumed` |
| ⑦ r1 FR-01 | CAS 与标记落盘之间存在崩溃窗口，存在双准入 | High | `probe_permit.go` 的 `persistClaim` | 已闭合：线性化点为标记文件 `O_CREATE` 与 `O_EXCL` 组合创建且 `Sync` 成功，`Require` 返回 nil 当且仅当标记已在盘 |
| ⑦ r1 FR-02 | `SealBinding.validate` 三条边界缺红测 | Medium | `readiness_test.go` | 已闭合：三条边界补测（⑥ Section D 记为已补） |
| ⑦ r1 FR-03 | 变异验证证据未落盘 | Low | `03-tests/` | 已闭合：复算落盘 `mutation-evidence.txt`，11 项 M1–M11 全部 RED；登记 D10 |
| ⑦ r2 NR-01 | MANIFEST 条目集、SHA256SUMS 与磁盘载荷三方不自洽 | Medium | 证据包元数据 | 已闭合：三方集合两两差集为空，抽样复算哈希一致 |
| ⑦ r2 NR-02 | D9 未在 MANIFEST 登记 | Low | `MANIFEST.md` | 已闭合：D9 与 D10 收入过程偏差登记表 |
| ⑦ r3 | 复查 0 新增可执行缺陷，给出放行条件 5 条 | Info | `07-final/review.md` | 已闭合：`FINAL REVIEW: PASS` |

## 交付

- `ReadinessGoal` 三值：`GoalProbe` / `GoalAcceptance` / `GoalProduction`，全部经 `EvaluateT027BReadiness` 产出。
- sealed `Readiness`：标签方法未导出，包外无法实现第四级；三个具体类型 `ProbeReadiness` / `AcceptanceReadiness` / `ProductionReadiness` 的方法集互斥，取错级别为编译红。
- `EvaluateT027BReadiness` 为纯函数：只读 `goal` 与 `input`，无文件读写、无状态写入、不消费 permit。
- 判定链顺序：`ErrProbePermitMissing` → `ErrReadinessProbeNotAcceptance` → `ErrReadinessEvidenceIncomplete` → `ErrReadinessAcceptanceMissing`（判据为 `Seal == nil`，不看 `ProbeEvidence`）→ `ErrReadinessProductionAT23Required`；哈希接受 `sha256:` 前缀或裸 64 位小写十六进制。
- `ProbePermit` 三态 `Unused` → `Claimed` → `Spent`，`Spent` 为 terminal，无回边；线性化点为标记文件 `O_CREATE` 与 `O_EXCL` 组合创建成功且 `Sync` 成功；`Require` 返回 nil 当且仅当标记已在盘；失败、超时与重启均不重置。
- `LoadProbePermit` 见标记即判 `Spent`（标记损坏亦 fail-closed）；`Complete` 幂等且拒绝 `Unused`；标记文件删除等同于重新签发，生产路径不得删除；落盘前 ctx 取消不烧额度。
- 测试：`internal/gates` 内 P44 用例 20 例 + 9 子例，全绿；变异验证 11 项（M1–M11）全部 RED。

## 偏差与登记

| 编号 | 内容 | 来源 | 判定 |
|---|---|---|---|
| D1 | `GoalProduction` 判据取 `Seal == nil`，不取 ① 的合取式 | R-1 / ⑤ r1 裁定 | 接受：L3 的不变量是不得缺 L2 材料，探针证据存在与否不改变该事实；已由 §八 R-1 固化 |
| D2 | `ReadinessInput` 用哈希串替代 chain 类型，① 的 L3 构造签名强制先有 L2 值这一类型层绑定未落地，改为运行时判据 | R-2 / ⑤ r2 判 Low | 接受：sealed interface 与三个方法集差异仍阻止跨层取用，残余为构造签名层；登记为已知限制 |
| D3 | `internal/console/**` 零改动，C7/N4 未端到端接线 | R-3 / ⑤ r1 B2 | 接受：接线后置，不得宣称端到端闭合；触发条件为 `AgentRunnerAssemblyOptions` 注入形态定稿 |
| D4 | 本机无 cgo，`-race` 未执行；替代证据为 `-count=50 -cpu 1,2,4` 与 `-count=20`；CI Race 腿不含 `internal/gates` | R-4 / ⑤ r1 裁定 3 | 接受：登记为已知限制，Race 腿扩展属独立切片 |
| D5 | ② 载体清理时误删本片 `tmp/P44-T027B-GATE-READINESS/` 中间件 | D9 / ⑦ r2 NR-02 | 接受：仅影响未入库一次性文件，最终门禁读数在 ⑨ 重取 |
| D6 | ③ 首轮变异验证证据未落盘 | D10 / ⑦ r1 FR-03 | 接受：已复算补录 `03-tests/mutation-evidence.txt`，11 项全 RED |

## 模型与计量

| 环节 | 载体 | 模型 |
|---|---|---|
| ① 架构 | `deep-reasoner` | DeepSeek V4 Pro |
| ② 实现 | `prfrail-implementer` | DeepSeek V4.1 Flash |
| ③ 测试 | `prfrail-tester` | DeepSeek V4.1 Flash |
| ⑤ 预审 | `deep-reasoner` | DeepSeek V4 Pro |
| ⑥ 扫描 | `independent-reviewer` | MAI-Code-1.1-Flash |
| ⑦ 终审 | `independent-reviewer` | GPT-5.3-Codex |
| ⑧ 文档 | `prfrail-documenter` | DeepSeek V4.1 Flash |

- **读数来源**：①⑤⑥⑦ 取自各自归档文件头的 `runSubagent` 声明；② 与 ③ 的实际载体与模型以 ⑩ 回执派单表的机械读取结果为准。
- **主控模型**：`DeepSeek V4.1 Flash`（High 档）。

## 门禁读数

- 命令：`node tools/gates/gate.js --all --skip-base`；tree 与 index 两作用域均为 `TOTAL_FAIL=0`。
- `G4b`：`newCjk=524`，`unknown=[]`。
- `G6-4`：`no expiring literal (scanned 462)`。
- 上列读数取自本片 ⑧ 末次实跑；⑨ 定稿读数以同一命令重取为准。

## 提交与 CI

- 提交 `9061b67` 已推送 origin；CI `41253798647` 首跑 ubuntu 腿红（`internal/adapters` 既有并发用例 flake，已登记 DR-13），同 run 复跑双腿全绿；包级 ok 行两腿均取到 `internal/gates`、`internal/adapters`、`internal/chain`、`internal/console`，Race 腿取到 `internal/adapters`、`internal/chain`。

## 后续片候选登记

- console 接线（R-3）：`runtime_agent.go:103`、`runtime_supervised.go:248`、`runtime.go:33` 与 `ai.go:39-90` 注入 `ReadinessInput`，解除端到端口径限制。
- Race 腿扩展（R-4）：把 `./internal/gates/...` 纳入 CI Race 腿，属 `ci.yml` 变更。
- 类型层绑定（R-2）：`ProductionReadiness` 构造签名接收 L2 值，恢复 ① 的类型层约束。
- ⑥ Section D 未覆盖面补测：`LoadProbePermit` 非 `NotExist` 读失败、`Require` 的 `ctx == nil` 分支、`Complete` 未知态、真跨进程双写。
- AT-23 签收工件：为 L3 提供不可伪造的验收摘要来源，避免自证。

## 已知限制

- **R-2**：L3 以 `Seal` 绑定与 `AT23Hash` 为输入，不含 ① 要求的构造签名层绑定；跨层取用仍由 sealed interface 与三个方法集差异阻挡。
- **R-3**：C7/N4 未端到端接线，只能宣称门控就绪与许可路径已满足当前切片验收前提，不得宣称端到端闭合。
- **R-4**：`-race` 未在本机执行，并发结论依赖 `-count=50 -cpu 1,2,4` 与 `-count=20` 替代证据；CI Race 腿不含 `internal/gates`。
- **⑥ Section D 已补覆盖**：`SealBinding.validate` 三边界（`Attempt < 1`、`TaskKey` 为空、`SealHash` 非法）已在 FR-02 补测，记为已覆盖。
- **⑥ Section D 仍未覆盖（登记，非缺陷）**：`LoadProbePermit` 非 `NotExist` 读失败、`Require` 的 `ctx == nil` 分支、`Complete` 未知态、真跨进程双写。
- **D9 / D10**：`tmp/` 中间件误删与首轮变异证据未落盘均已登记，补救记录见证据包 `MANIFEST.md` 的过程偏差登记表。

```artifacts
docs/ARCHITECTURE.md	repo
docs/ARCHITECTURE_EN.md	repo
docs/CONTRACTS.md	repo
docs/CONTRACTS_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P44-T027B-GATE-READINESS/00-briefing.md	repo
docs/validation/evidence/P44-T027B-GATE-READINESS/01-arch-design/review.md	repo
docs/validation/evidence/P44-T027B-GATE-READINESS/03-tests/mutation-evidence.txt	repo
docs/validation/evidence/P44-T027B-GATE-READINESS/05-prereview/review.md	repo
docs/validation/evidence/P44-T027B-GATE-READINESS/06-scan/review.md	repo
docs/validation/evidence/P44-T027B-GATE-READINESS/07-final/review.md	repo
docs/validation/evidence/P44-T027B-GATE-READINESS/08-ci/ci-evidence.txt	repo
docs/validation/evidence/P44-T027B-GATE-READINESS/MANIFEST.md	repo
docs/validation/evidence/P44-T027B-GATE-READINESS/SHA256SUMS.txt	repo
docs/validation/evidence/P44-T027B-GATE-READINESS/freedom-list.md	repo
docs/validation/p44-gate-readiness.md	repo
docs/validation/p44-gate-readiness_EN.md	repo
internal/gates/probe_permit.go	repo
internal/gates/probe_permit_test.go	repo
internal/gates/readiness.go	repo
internal/gates/readiness_test.go	repo
```
