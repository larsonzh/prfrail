# ⑤ 预审记录（两轮原样归档）

> 载体：runSubagent(agentName="deep-reasoner", model="DeepSeek V4 Pro (deepseek)")。r1 transcript 偏移 117179，r2 偏移 117574。

## r1（首轮）

## A 结论

`PRE-REVIEW: PASS WITH FIXES`（无 Critical/High；3 项需处置后放行，详见 B/C）

## B 发现

| # | 严重度 | file:line | 问题 | 修法 |
|---|---|---|---|---|
| B1 | Medium | `readiness.go:169-171` vs `readiness_test.go:126-133` | `GoalProduction` 缺 L2 判据实现为仅 `Seal==nil`（测试钉死 `ErrReadinessAcceptanceMissing`），与 ① 所述「`Seal==nil` 且 `ProbeEvidence` 为空」不一致；① 方案尚未归档（`01-arch-design/` 为空），无法对文核验 | 非 C7 漏洞（实现 fail-closed）；按自由度清单优先序（① 方案含修订注记 > 既有测试）由用户/② 落一条修订注记固化语义，或改实现+测试回 ① 字面；不得静默保留 |
| B2 | Medium | `internal/gates/` 整体 | C7/N4 仅提供纯 API，无任何生产调用点（全仓 `git grep` 0 命中）；台账「规模 = console 门控入口 + gates 判定实现」但工作树零 console 改动，C7 活路径（`ai.go:116→195` verify 链）未接新判定 | 本片明确接线 console 门控入口，或显式登记「接线后置切片」并在完成记录如实写，且不得宣称 C7 已端到端闭合 |
| B3 | Low | `probe_permit.go:31-32` | 注释声称「every returned error wraps exactly one of them」与实现矛盾：`Require` 两处 ctx 路径（`:145`、`:161`）与 `LoadProbePermit` 读错误（`:113`）不包任何哨兵 | 改注释为「permit-state 错误恰包一枚哨兵；ctx 取消原样透传（`errors.Is(err, context.Canceled)`）」；可选：claim 后取消用 `errors.Join(ctx.Err(), ErrProbePermitInUse)` |
| B4 | Low | `readiness_test.go:119-122` | 两个负向类型断言（`*ProbeReadiness` 不可断言为 `*AcceptanceReadiness`/`*ProductionReadiness`）为永真断言，违反文件头自述的可证伪规则 | 删除或改为编译期钉（现有 compile-time pins 已覆盖）；勿计入覆盖证据 |
| B5 | Low | `probe_permit_test.go` | 缺三处负例：标记存在但 JSON 损坏的 fail-closed；同路径两个 `NewProbePermit` 实例的 O_EXCL 竞争（仅单实例 40 goroutine 覆盖 CAS）；`Complete` 后 `LoadProbePermit` 仍拒 | 各补 1 例（损坏体→Spent；双实例恰 1 胜；重启后 Complete 幂等） |
| B6 | Info | `readiness.go:153-154` | `GoalProbe` 只查 `Permit != nil`，不查其状态：已 Spent 的 permit 仍可产出 `ProbeReadiness`（消费与判定分离属设计，但需写明否则易被误用为额度绕过） | 注释明示「判定不消费、不校验状态；消费由 `Require` 承担」 |
| B7 | Info | `probe_permit.go:105-126` | `LoadProbePermit` 对损坏 JSON 静默返回 Spent，调用方无法区分「已消费」与「标记损坏」；fail-closed 正确但缺测试（同 B5） | 补测试；可选返回包装错误区分 |
| B8 | Info | `probe_permit.go:83-91` | `NewProbePermit` 不校验 id/path；空 path 首调 `Require` 即烧额度（persist 失败→Spent→`ErrProbePermitConsumed`）且无持久标记；CAS 到标记落盘之间崩溃窗口内获胜者未返回 nil（无双准入），但标记文件删除后重启额度即重置、标记生命周期未定义 | 构造函数校验 path 或文档化 fail-closed 语义；报告登记标记生命周期为已知限制 |
| B9 | Info | `00-briefing.md` vs `MANIFEST.md` | briefing 自检行写 `HEAD=8a46b01`，MANIFEST 与 git 实况为 `a87b078`（v33 清单提交在其后） | 完成记录注明起跑后 HEAD 前进，非缺陷 |

## C 三处裁定项的判定

1. **规格与测试冲突**：判 **② 的语义正确、不构成 C7 漏洞**。`GoalProduction` 的不变量是「L3 不得缺 L2 材料」，判据应为 `Seal==nil`（探针证据存在与否不改变「L2 缺失」事实）；① 的合取式只是把该输入归为 `ErrReadinessProbeNotAcceptance`，两者同为 fail-closed，仅诊断哨兵不同。② 版更贴近 C7 核心（探针证据永远不能充当任何更高层的材料，`readiness.go:156-164` 已对 acceptance 支单独区分）。但这是**未记录的规格漂移**（D 类），且 ① 方案未归档致无法对文——必须落修订注记或改测试回 ①，否则不得进入 ⑦。
2. **非哨兵错误路径**：判 **可接受、不削弱 N4**。ctx 取消与 permit 状态机正交，保留 `errors.Is(err, context.Canceled)` 是标准 Go 惯例；额度在 claim 前后均不被重置（claim 前未取得、claim 后 fail-closed）。**但 B3 的注释声称必须更正**——现状是文档说谎而非设计缺陷。
3. **`-race` 替代**：判 **足以支撑本代码的并发结论，但须登记为已知限制**。共享状态仅 `atomic.Uint32`（构造后其余字段不可变），无数据竞争可能；跨进程由 `O_EXCL`（OS 语义）承担，与内存模型无关；`-count=5` 实测稳定、40 goroutine 恰 1 胜。注意 CI 的 Race 腿目前仅覆盖 `adapters`/`chain`（不含 `gates`），报告须写明并候选 OB 扩展 Race 腿。

## D 交叉核对

1. **哨兵无冲突**：既有 27 枚（console 18 + gates 9，与 ② 普查一致）；7 枚新哨兵与全部新 API 名在已跟踪 `.go` 中 0 命中 ⇒ 无命名/语义冲突，且证明无生产调用点（佐证 B2）。
2. **机械门禁**：`go build ./...`、`go vet ./internal/gates/`、`gofmt -l` 全绿；`go test -count=1/-count=5 ./internal/gates/` 全绿；② 的 `gate-mid.txt` `TOTAL_FAIL=0`。
3. **编码硬规则**：4 个新 `.go` 实测 `BOM=False CRLF=0` ✓；台账新增行含 `SELFTEST 697/697` 未触 G6-4（后无「行」字）。
4. **边界**：工作树仅 `internal/gates/**` 4 新文件 + 台账对 + 证据包；未触 `chain`/`adapters`/`evidence`/`tickets` ✓。
5. **密封性**：`Readiness` 未导出标签方法 ⇒ 包外不可实现（含嵌入与同名方法均无法满足，方法身份含包路径）；别名不产生新层级（动态类型不变）；三个 accessor 各自挂本层类型，无跨层读取通道 ✓。
6. **额度状态机**：无任何 `Store(Unused)`/回退 CAS；`Complete` 对 Spent 幂等、对 Unused 拒（`ErrProbePermitInUse`）；`errors.Is` 在 readiness 全部 8 条返回路径完整（`readiness.go` 每支恰包一枚）。

## E 可证伪性

以下任一成立即推翻本结论：① ① 方案归档后逐字要求合取式且无修订注记 ⇒ B1 升为流程阻断；② 未来接线把 `ProbeEvidence` 传入 `Seal` 槽位 ⇒ C7 复发；③ CI Race 腿扩展至 `internal/gates` 后报 data race ⇒ 裁定 3 作废；④ 40-goroutine 测试在 CI 复现 flake。未验证项：③ 测试先行的真实时序（须 ⑩ transcript 佐证载体分离）。**过程披露**：审查期间一次命令重定向误在 `tmp/P44-T027B-GATE-READINESS/` 创建 `sentinels.txt`（违反只读约束），未删除（删除同样被禁），请 ② 在 ⑧c 一并清理。

`PRE-REVIEW: PASS WITH FIXES`

---

## r2（复审）

# P44 ⑤ 预审 r2 报告

## A 结论

`PRE-REVIEW: PASS WITH FIXES`（3 项 Low 待修，无 Medium/High；本轮零文件操作）

## B 发现

| 编号 | 严重度 | file:line | 问题 | 修法 |
|---|---|---|---|---|
| r2-1 | Low | `docs/validation/evidence/P44-T027B-GATE-READINESS/MANIFEST.md`、`SHA256SUMS.txt` | B1 归档**元数据未闭环**：`01-arch-design/review.md` 已落盘，但 MANIFEST 条目清单仍把该件列在「待追加」、SHA256SUMS 仅含 2 件（00-briefing、freedom-list），与自述「条目集三方相等」矛盾；将来 ⑧c 若按 SHA256SUMS 校验会失败 | 按轮更新 MANIFEST 条目清单并以一次性脚本重算 SHA256SUMS（MANIFEST 已承诺该流程，补执行即可） |
| r2-2 | Low | `readiness.go:37` vs `readiness.go:157` | 注释「Every error returned by EvaluateT027BReadiness wraps exactly one of them（4 枚 readiness 哨兵）」与 GoalProbe 分支返回 `ErrProbePermitMissing` 字面矛盾——该哨兵不在四枚之列，恰是测试 `TestReadinessProbeGoalRequiresPermit` 钉住的跨域语义 | 注释补例外句：「except the goal-probe branch, which wraps ErrProbePermitMissing（N4 域）」 |
| r2-3 | Low | `probe_permit.go:44` vs `probe_permit.go:184-190` | 注释「Plain I/O failures … returned as-is」未限定范围：对 `LoadProbePermit` 读失败成立（无哨兵透传），但 `Require→persistClaim` 的 I/O 失败（MkdirAll/OpenFile/Write/Sync）被包成 `ErrProbePermitConsumed` 且底层错误用 `%v`（`errors.Is` 不可达），与「as-is」字面张力 | 包级注释补限定句：「except Require 的持久化失败：额度作废属状态错误，包 ErrProbePermitConsumed」 |

（变异验证 M1–M6 原始输出未落盘，仅 `review.md:152` 一句引用——归入 §C-B5/D，不另立条目。）

## C r1 各项闭合核验

- **B1 — 已闭合**。`review.md` §八 如实登记 R-1/R-2/R-3/R-4；R-2 原文「类型层绑定**未落地**，改为运行时判据」，未弱化为「仅命名变更」，且点名 ⑥⑦ 复核义务。`readiness.go:177-186` 代码注记指向同文件。瑕疵仅元数据不同步（见 r2-1），不影响登记实质。
- **B2 — 已闭合**。grep 全仓 `EvaluateT027BReadiness`/`NewProbePermit`/`LoadProbePermit` 仅出现于 `internal/gates`（实现+测试），console 零调用；R-3 明写「C7/N4 不得宣称端到端闭合」，台账边界列明「console 零改动」。
- **B3 — 部分闭合**。「permit 状态错误恰包一枚哨兵」逐路径成立：`LoadProbePermit` 空 path→Missing；`Require` nil permit/nil ctx→Missing，CAS 败→refusal 三分支均包 InUse/Consumed/InUse，persistClaim 败→Consumed；`Complete` nil→Missing、Unused→InUse、unknown→Consumed。「ctx 取消/超时原样透传」两处 `ctx.Err()` 成立。「I/O 错误原样透传」仅对 LoadProbePermit 成立，persistClaim 例外未在注释限定（r2-3）。
- **B4 — 已闭合**。`readiness_test.go:26-54` 仅存编译期钉：3×`_ Readiness = (*X)(nil)` + 3×方法表达式（`(*ProbeReadiness).Permit` 等）；`TestReadinessLevelsCannotBeMixed` 内注释明说原永真类型断言已删且说明理由。
- **B5 — 已闭合**。4 例俱在：`probe_permit_test.go:145-183`（DamagedMarkerStaysSpent）、`:187-207`（SecondInstanceLosesOnly）、`:210-228`（CompleteOnUnusedRejected）、`:231-255`（RestartThenCompleteIsIdempotent）。**证据链缺口**：③ 变异验证 M1–M6 的原始输出未落盘，仅 review.md:152 一句引用，证据包无对应文件——建议 ⑧c 前补录或明示删除。
- **B6 — 已闭合**。`readiness.go:146-153` GoalProbe 分支注释与实现一致：判定不消费不校验，消费归 `Require`，测试 `TestReadinessProbeDoesNotConsumePermit` 双向钉住。
- **B7 — 已闭合**。`readiness_test.go:227-236`（UnknownGoalRejected 钉 default 分支）、`:241-266`（ProbeDoesNotConsumePermit）均存在且断言语义正确。
- **B8 — 已闭合**。`probe_permit.go:96-109`（NewProbePermit doc：不校验 id/path；空 path→Spent+Consumed；删除标记=重新签发）与 `:184-192`（Require doc 同义）齐备，实现（persistClaim 空 path 报错→包 Consumed）一致。
- **B9 — 已闭合（时点未到）**。HEAD 现为 `a87b078`，起跑 `8a46b01` 已前进，与台账「依赖 P43 已闭环(8a46b01)」一致；台账 P44 详情「完成记录」字段仍空，注明动作属 ⑧ 落笔时点，当前不构成缺口。

## D 残余风险

- **R-2 类型层绑定缺失 → 可接受（Low，不构成 C7 Medium/High 缺口）**。残余类型系统约束判据：① sealed interface `readinessTag` unexported（readiness.go:75-78、111-114）包外无法伪造第四级；② 方法集差异（116-140）跨层取用编译红；③ `ProductionReadiness` 值内嵌 `ReleaseBinding{Seal SealBinding}` **值类型**（91-96、199-204），L3 值在类型层必然携带一份 L2 材料。缺失的仅是「L3 构造函数签名接收 L2 值」；但三构造点全仓唯一（readiness.go:159/172/199，grep 证实无第二构造点、字段 unexported），运行时判据（Seal==nil→ErrReadinessAcceptanceMissing + validate()）集中单点且被 `productionWithoutAcceptance`/`RequiresAT23` 多例双向钉住 ⇒ 降级幅度有限。**⑥ 复核点**：确认包内无第二构造点、`ProductionReadiness` 无法经零值路径逃逸。
- **R-3 接线后置 → 可接受**。明写不宣称端到端闭合 + 台账边界登记 + ①§七风险 1 预登记 12 处装配回归面；风险完整转交后续接线切片，本片可验证范围未虚报。
- **R-4 race 缺失 → 有条件接受**。本机 `CGO_ENABLED=0` 环境限制如实登记；替代证据（-count=50/20 + 变异验证）有据；CI Race 腿确证不含 gates（`ci.yml:84` 仅 `./internal/adapters/... ./internal/chain/...`）。遗留：并发结论（CAS 线性化）race 证据真空，建议后续切片将 `./internal/gates/...` 纳入 Race 腿（超本片边界）。

## E 可证伪性

- `TestProbePermitDamagedMarkerStaysSpent`（probe_permit_test.go:145-183）：**可红**。删除 `LoadProbePermit` 中 `state.Store(spent)` 一行（损坏 marker 退化为 unused）→ `damaged.Require` 返回 nil → 断言红；控制组（无 marker 必须可认领）同时封堵「Load 恒 Spent」的假实现——双向可证伪。硬载入分支（返回 error 且 permit==nil）被显式允许，不产生额度双发，可接受。
- `TestReadinessProbeDoesNotConsumePermit`（readiness_test.go:241-266）：**可红**。若 GoalProbe 分支误调 `Require`（消费额度）→ `os.Stat` 见 marker 或后续 `Require` 返回 Consumed → 红；仅写 marker 不消耗内存态的变异同样被 stat 断言捕获。
- 附带：`TestProbePermitConcurrentRequireSingleWinner` 可红（40 者恰 1 赢家，CAS 破坏即红）。

**闭集判定：`PRE-REVIEW: PASS WITH FIXES`**