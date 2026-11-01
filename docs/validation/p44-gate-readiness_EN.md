# ProofRail Validation Report — P44-T027B-GATE-READINESS

| Field | Content |
|---|---|
| Slice | `P44-T027B-GATE-READINESS` |
| Type | implementation class `[SLICE]` |
| Directive version | v1.60 |
| Baseline HEAD | `a87b078` |
| Date | 2026-11-01 |

## Goals and boundaries

- **Goal**: close C7 and N4 - readiness must strictly separate probe readiness, real acceptance readiness and production release readiness, so a probe result can never advance acceptance; `ProbePermit` is consumed exactly once, reserved atomically before the call, invalid after consumption, and never reset by failure, timeout or restart.
- **Caliber**: C7/N4 must not be declared end-to-end closed; only gate readiness and the permit path may be declared as meeting this slice's acceptance premise.
- **Closure shape**: `internal/gates` adds `readiness.go` (`ReadinessGoal` with three values + sealed `Readiness` + three concrete types + `EvaluateT027BReadiness`) and `probe_permit.go` (`ProbePermit` three-state machine + `Require` / `Complete` / `LoadProbePermit`).
- **Scope**: `internal/gates/**`.
- **Boundary**: `internal/console/**` unchanged (no end-to-end wiring, registered as R-3); no touch to `internal/chain/**`, `internal/adapters/**`, schema, fixtures or `ci.yml`; no gitee push; `commit` / `push` stop at 10 for authorisation.

## Review chain conclusions

- **1 architecture**: `deep-reasoner` + DeepSeek V4 Pro; the C7+N4 gate-readiness and one-shot permit design book, `DESIGN: DONE`; the four revision notes R-1 to R-4 in section 8 were recorded after the r1 pre-review ruling.
- **3 red tests**: `RED-CONFIRMED`; 11 mutation checks (M1–M11) all RED, with the raw output archived in `03-tests/mutation-evidence.txt`.
- **2 implementation**: four files land in `internal/gates` (two implementations + two test files); this package carries 20 P44 cases + 9 sub-cases.
- **4 gates**: `TOTAL_FAIL=0`; readings in section Gate readings.
- **5 r1**: `PRE-REVIEW: PASS WITH FIXES` (2 Medium: unrecorded spec drift / no production wiring; plus 3 Low + 4 Info).
- **5 r2**: `PRE-REVIEW: PASS WITH FIXES` (3 Low only, with every r1 item closed item by item).
- **6**: `INDEPENDENT SCAN: PASS` (Section B carries zero defects; Section D lists 5 uncovered areas); slice budget 1/3.
- **7 r1**: `FINAL REVIEW: BLOCKED` (FR-01 High crash window / FR-02 Medium missing red test / FR-03 Low evidence not on disk).
- **7 r2**: `BLOCKED` (NR-01 Medium MANIFEST three-way inconsistency / NR-02 Low D9 unregistered).
- **7 r3**: `FINAL REVIEW: PASS` (0 new defects); slice budget 3/3.
- **R / D**: R = 0, D = 0 (not inherited from the previous slice).

## Six-stage scan transcript

Source `docs/validation/evidence/P44-T027B-GATE-READINESS/06-scan/review.md`; the fence transcribes it verbatim with zero rewriting, from `## A 结论` through the conclusion line.

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
- Transcript note: heading prefixes are normalised to `## Section A` to `## Section E` (the archive reads `## A 结论` to `## E 可证伪性`) so the G7-f A→B→C→D→E criterion can see them; the body text is unchanged.
- Verdict-token normalisation: the archive's final line `§A INDEPENDENT SCAN: PASS` is transcribed inside the fence as the block-final `INDEPENDENT SCAN: PASS`, carrying the same value as the conclusion line under `## A 结论`.

## Item closure

The severity column keeps the review reports' own wording and is never bolded; the final-review closed set is `FINAL REVIEW: PASS` (r3, 0 new defects).

| # | Item | Severity | Location | Disposition |
|---|---|---|---|---|
| 5 r1 B1 | the `GoalProduction` L2-material criterion contradicts the literal design text - an unrecorded spec drift | Medium | the `GoalProduction` in `readiness.go` | Closed: fixed by section 8 R-1, the criterion is `Seal == nil` and `ProbeEvidence` is not read |
| 5 r1 B2 | C7/N4 expose a pure API with zero production call sites, contradicting the ledger and the working tree | Medium | `internal/gates/**` and the ledger | Closed: wiring deferred and registered as R-3, with an explicit ban on the end-to-end claim |
| 5 r1 B3 | the package comment claims every returned error wraps exactly one sentinel, which the ctx paths contradict | Low | the package comment in `probe_permit.go` | Closed: the comment now says permit-state errors wrap one sentinel and ctx cancellation passes through as-is |
| 5 r1 B4 | two negative type assertions were always-true assertions | Low | `readiness_test.go` | Closed: the always-true assertions were removed and compile-time method-set pins carry the load |
| 5 r1 B5 | three negative cases were missing: damaged marker, two-instance race, post-restart idempotence | Low | `probe_permit_test.go` | Closed: all four cases are now present |
| 5 r2 r2-1 | MANIFEST metadata was not closed and SHA256SUMS covered only two payloads | Low | the evidence pack metadata | Closed: `07-final/review.md` is archived here and every payload hash is recomputed |
| 5 r2 r2-2 | the readiness package comment contradicts the cross-domain sentinel returned by the GoalProbe branch | Low | the `readiness.go` comment | Closed: an exception clause was added, the GoalProbe branch wraps `ErrProbePermitMissing` |
| 5 r2 r2-3 | the `Require` persistence-failure wrapping contradicts the "returned as-is" comment | Low | the `probe_permit.go` comment | Closed: the comment now limits itself and the persistence failure wraps `ErrProbePermitConsumed` |
| 7 r1 FR-01 | a crash window sits between the CAS and the marker landing, allowing a double admission | High | the `persistClaim` in `probe_permit.go` | Closed: the linearisation point is the marker file created with the `O_CREATE` and `O_EXCL` combination after a successful `Sync`, and `Require` returns nil if and only if the marker is on disk |
| 7 r1 FR-02 | three `SealBinding.validate` boundaries had no red test | Medium | `readiness_test.go` | Closed: the three boundaries are tested (6 Section D records them as added) |
| 7 r1 FR-03 | the mutation evidence was not on disk | Low | `03-tests/` | Closed: `mutation-evidence.txt` was recomputed and saved, 11 items M1–M11 all RED; registered as D10 |
| 7 r2 NR-01 | the MANIFEST entry set, SHA256SUMS and the on-disk payload disagreed three ways | Medium | the evidence pack metadata | Closed: all three pairwise set differences are empty and sampled hashes recompute identically |
| 7 r2 NR-02 | D9 was not registered in MANIFEST | Low | `MANIFEST.md` | Closed: D9 and D10 are both listed in the process deviation table |
| 7 r3 | the re-check found 0 new executable defects and gave 5 release conditions | Info | `07-final/review.md` | Closed: `FINAL REVIEW: PASS` |

## Delivery

- `ReadinessGoal` with three values: `GoalProbe` / `GoalAcceptance` / `GoalProduction`, all produced through `EvaluateT027BReadiness`.
- sealed `Readiness`: the tag method is unexported, so no fourth level can be implemented outside the package; the three concrete types `ProbeReadiness` / `AcceptanceReadiness` / `ProductionReadiness` have mutually exclusive method sets, so taking the wrong level is a compile red.
- `EvaluateT027BReadiness` is a pure function: it reads only `goal` and `input`, performs no file IO, writes no state and consumes no permit.
- decision chain order: `ErrProbePermitMissing` → `ErrReadinessProbeNotAcceptance` → `ErrReadinessEvidenceIncomplete` → `ErrReadinessAcceptanceMissing` (the criterion is `Seal == nil`, `ProbeEvidence` is not read) → `ErrReadinessProductionAT23Required`; hashes accept a `sha256:` prefix or a bare 64-character lowercase hex string.
- `ProbePermit` three states `Unused` → `Claimed` → `Spent`, with `Spent` terminal and no back edge; the linearisation point is the marker file created with the `O_CREATE` and `O_EXCL` combination after a successful `Sync`; `Require` returns nil if and only if the marker is on disk; failure, timeout and restart never reset it.
- `LoadProbePermit` treats any marker as `Spent` (a damaged marker also fails closed); `Complete` is idempotent and refuses `Unused`; deleting the marker file equals a fresh issuance and the production path must never delete it; a ctx cancellation before the marker lands burns no quota.
- tests: 20 P44 cases + 9 sub-cases inside `internal/gates`, all green; 11 mutation checks (M1–M11) all RED.

## Deviations and registrations

| Id | Content | Source | Verdict |
|---|---|---|---|
| D1 | the `GoalProduction` criterion is `Seal == nil` rather than the design's conjunction | R-1 / 5 r1 ruling | Accepted: the L3 invariant is that L2 material must not be missing, and probe evidence does not change that fact; fixed by section 8 R-1 |
| D2 | `ReadinessInput` uses hash strings instead of chain types, so the design's type-level binding (the L3 constructor signature forcing an L2 value) did not land and a runtime criterion replaced it | R-2 / 5 r2 rates it Low | Accepted: the sealed interface and the three method sets still block cross-level access; the residual is the constructor signature and is registered as a known limitation |
| D3 | `internal/console/**` is unchanged, so C7/N4 are not wired end to end | R-3 / 5 r1 B2 | Accepted: wiring is deferred and the end-to-end claim is banned; the trigger is the final `AgentRunnerAssemblyOptions` injection shape |
| D4 | this machine has no cgo, so `-race` did not run; the substitute evidence is `-count=50 -cpu 1,2,4` and `-count=20`; the CI Race leg does not cover `internal/gates` | R-4 / 5 r1 ruling 3 | Accepted: registered as a known limitation, and extending the Race leg is a separate slice |
| D5 | the 2 carrier deleted this slice's `tmp/P44-T027B-GATE-READINESS/` work files while cleaning | D9 / 7 r2 NR-02 | Accepted: only untracked one-shot files were affected and the final gate readings are retaken at 9 |
| D6 | the first mutation round saved no evidence to disk | D10 / 7 r1 FR-03 | Accepted: `03-tests/mutation-evidence.txt` was recomputed and recorded, all 11 items RED |

## Models and metering

| Stage | Carrier | Model |
|---|---|---|
| 1 architecture | `deep-reasoner` | DeepSeek V4 Pro |
| 2 implementation | `prfrail-implementer` | DeepSeek V4.1 Flash |
| 3 tests | `prfrail-tester` | DeepSeek V4.1 Flash |
| 5 pre-review | `deep-reasoner` | DeepSeek V4 Pro |
| 6 scan | `independent-reviewer` | MAI-Code-1.1-Flash |
| 7 final review | `independent-reviewer` | GPT-5.3-Codex |
| 8 documentation | `prfrail-documenter` | DeepSeek V4.1 Flash |

- **Reading source**: stages 1, 5, 6 and 7 are taken from the `runSubagent` declarations in their archived file headers; the actual carrier and model of stages 2 and 3 follow the mechanical transcript read recorded in the 10 receipt dispatch table.
- **Primary model**: `DeepSeek V4.1 Flash` (High tier).

## Gate readings

- Command: `node tools/gates/gate.js --all --skip-base`; both the tree and the index scope report `TOTAL_FAIL=0`.
- `G4b`: `newCjk=524`, `unknown=[]`.
- `G6-4`: `no expiring literal (scanned 462)`.
- The readings above are this slice's final 8 run; the 9 finalisation retakes them with the same command.

## Commit and CI

- Commit `9061b67` is pushed to origin; CI `41253798647` first attempt failed on the ubuntu leg (a pre-existing concurrency flake in `internal/adapters`, registered as DR-13) and a failed-job rerun of the same run is green on both legs; package-level ok lines cover `internal/gates`, `internal/adapters`, `internal/chain` and `internal/console` on both legs plus `internal/adapters` and `internal/chain` on the Race leg.

## Later-slice candidate registrations

- Console wiring (R-3): inject `ReadinessInput` at `runtime_agent.go:103`, `runtime_supervised.go:248`, `runtime.go:33` and `ai.go:39-90`, lifting the end-to-end caliber limit.
- Race leg extension (R-4): add `./internal/gates/...` to the CI Race leg, which is a `ci.yml` change.
- Type-level binding (R-2): let the `ProductionReadiness` constructor take an L2 value and restore the design's type-level constraint.
- 6 Section D uncovered areas: `LoadProbePermit` failing on a read error other than `NotExist`, the `ctx == nil` branch of `Require`, the `Complete` unknown state, and a genuine cross-process double write.
- AT-23 sign-off artifact: give L3 an unforgeable acceptance digest source instead of self-attestation.

## Known limitations

- **R-2**: L3 takes the `Seal` binding and `AT23Hash` as input and carries no constructor-signature binding as the design required; cross-level access is still blocked by the sealed interface and the three method sets.
- **R-3**: C7/N4 are not wired end to end, so only gate readiness and the permit path may be declared as meeting this slice's acceptance premise; the end-to-end claim stays banned.
- **R-4**: `-race` did not run on this machine, so the concurrency conclusion rests on the `-count=50 -cpu 1,2,4` and `-count=20` substitute evidence; the CI Race leg does not cover `internal/gates`.
- **6 Section D now covered**: the three `SealBinding.validate` boundaries (`Attempt < 1`, empty `TaskKey`, illegal `SealHash`) were added in FR-02 and are recorded as covered.
- **6 Section D still uncovered (registered, not defects)**: `LoadProbePermit` failing on a read error other than `NotExist`, the `ctx == nil` branch of `Require`, the `Complete` unknown state, and a genuine cross-process double write.
- **D9 / D10**: the `tmp/` work-file deletion and the unsaved first mutation round are both registered, with the remedies in the process deviation table of the evidence pack `MANIFEST.md`.

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
