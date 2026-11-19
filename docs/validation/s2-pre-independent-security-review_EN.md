# P79 Verification Report: Pre-S2 Independent Security Review (Deferred Supplement)

Date: 2026-11-19. Status: in progress (stage (2), documentation implementation). Slice: P79 `S2-PRE-INDEPENDENT-SECURITY-REVIEW` (`[SLICE]`, assessment type).

## 1. Slice Identity and Authorization

- Identifier: P79 `S2-PRE-INDEPENDENT-SECURITY-REVIEW`; tier `[SLICE]`; type = assessment (independent security review supplement); scale M-L (already frozen in the v65 snapshot §4; this slice does not change the snapshot).
- Controller: `DeepSeek V4.1 Flash (deepseek)` High. The judgment load of this slice sits with (7); the controller only performs orchestration / materialization / transcription / gates / receipt, does not supplement or soften any review judgment from its own reading, and never issues a face conclusion on behalf of (7).
- Trigger basis: ADR-047 (`docs/ADR_REGISTER.md:58`). Three obligation elements = no supplement inside S1, exemption deferred to S2, and one independent security review with a registered conclusion before S2 starts; due date = the earlier of the S2 start date and 2027-03-31, and if the supplement is missing at the due date, S2 measured runs must not start.
- Review object: the "independent security review conclusion" gate of `PRODUCT_REQUIREMENTS.md` §7 (`:109` blocking gate sentence); its historical state is "never executed", evidenced by `SECURITY.md:5`. That object differs from the "independent review" row of the `docs/validation/s1-exit.md:59` gate table (T017 baseline); the evidence is not interchangeable and this slice does not change the s1-exit gate table.
- Four authorization elements: adjudicator = user; date = 2026-11-19; reason = supplement before the ADR-047 compensation obligation falls due; basis = this start message + the v65 snapshot §4.
- (7) budget ruling: tier L = 1 final review + 3 re-reviews. That ruling **exceeds the literal standard of §7.5 ("M-L counts as M"), is a direct user ruling, and must not be used as a precedent**, and must be recorded truthfully in the metering section of the report.
- Probes: 0 (no external real invocation in this slice).
- Blind review: the three semantics that §7.7 mandates sampling (ownership / stop / identity) are not touched by this slice, so none is sampled; the reason is registered in step 0.
- Boundary: no change to the S1 Exit eight-gate criteria; no change to the `SECURITY{,_EN}.md` threat-model baseline (conclusion archiving + annotation only); no change to the directive / `CONTRACTS` / schema / fixtures / product code / CI; no new OB / ADR (unless the user rules so); no network calls, no external commands, no reading `%TEMP%`, no reading files outside the repository.

Model and review chain (set once before start; no mid-flight tier upgrade):

| Stage | Role | Model | Tier / budget |
|---|---|---|---|
| (1) | Architect (`deep-reasoner`) | `DeepSeek V4 Pro (deepseek)` | Max; produces the review design |
| (5) | Pre-reviewer | `DeepSeek V4 Pro (deepseek)` | budget 1 + 1; pre-review of the report draft |
| (6) | Independent scan | `MAI Code 2 Flash (copilot)` | ≤3; output-contract conformance scan of the draft |
| (7) | Independent final review | `GPT-6.2 Sol (copilot)` | Max; 1 final review + 3 re-reviews, replacing the routine final-review tier of the tier-L table |
| Batch assist | Auxiliary execution | `GLM-5.3-Flash (glm)` | ≥2 runs, strictly serial; limited to the four kinds enumeration / excerpt-and-number checks / bilingual mirror checks / MANIFEST drafting; no cherry-picking |

## 2. Step-0 Anchor Measurements

| Item | v65 snapshot value | Measured value | Conclusion |
|---|---|---|---|
| `HEAD` | `9fd8cd6` | `f211bca` | matches (v65's own landing commit; snapshot §1 already stated that HEAD advances once v65 itself lands and the live state governs) |
| `origin/main` | `9fd8cd6` | `f211bca` | matches (same wording as above) |
| Directive version | `v1.71` (CN / EN headers both `v1.71`) | `v1.71` (CN / EN headers identical) | matches |
| `SELFTEST` | `node tools/gates/gate.js --selftest` | — | `PASS 740/740` (exit 0); readout `gate-selftest-final.txt` |
| Whole-tree encoding | `node tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/enc-scan.js` | — | every violation sits in a frozen evidence pack under `docs/validation/evidence/**` (non-frozen main tree is zero); the encoding hard rule is section 6.1 |
| Three scopes | `TOTAL_FAIL=0` | `index` / `tree` / `ci` all `TOTAL_FAIL=0` | matches |
| v63 whole-file line count | 222 (CN / EN symmetric) | 222 (split on LF, CN / EN symmetric) | matches |
| v64 whole-file line count | 224 (CN / EN symmetric) | 224 (split on LF, CN / EN symmetric) | matches |
| v65 whole-file line count | re-checked by v66 | HEAD advances after this version lands | not adjudicated under the snapshot wording |
| `tmp` top level | `gates-cache` + `.gitkeep` | `.gitkeep` + `gates-cache` + `P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW` (added by this slice, legitimate) | matches |
> **Timing caliber (section 11.3, v1.53)**: the readings in this section are the **step-0 baseline** (`changed=0`) values; the current readings for the change-set gates are the ones in the gate-results chapter and the (10) receipt (before the stage-6 transcription lands, `G7-f` is a known pending item and must turn green before (9)).

## 3. Review Scope and Three-Face Mapping

The three faces are taken verbatim from the ADR-047 "uncovered scope" and must not be added to or removed from:

| Face | Review object | Criterion source | Conclusion tier |
|---|---|---|---|
| (1) Threat-model face | `SECURITY.md` §1 (assets and trust boundaries, body :11), §2 (write-permission matrix table, :15–:21) plus the referenced misuse-scenario carrier §3 STRIDE table (SEC-01…SEC-15, :29–:44) and §3.1 (C6 reserved slot, :50–:54); `SECURITY_EN.md` has the same structure. Review action = independent re-check: completeness of the asset list, consistency between the §1 trust claims and the row/column permissions of the §2 matrix, whether the STRIDE "protection and responsibility" column can be delivered by existing components, and whether the misuse scenarios are closed | `SECURITY.md:5` (REQ-004/005/007/013/019/021-026/028); `PRODUCT_REQUIREMENTS.md:109` (gate sentence); RFC §12, §16.7–16.8 (design source); `TEST_STRATEGY.md` §2 AT-01…AT-24 table (AT number resolution, `SECURITY.md:46` table-tail sentence); `ADR_REGISTER.md:58` uncovered scope (1) | Report "threat-model face" subsection: that face's Section A–E block + face conclusion line `FACE 1: PASS\|FINDINGS` (template in §6) |
| (2) Hard security gate face | Criterion side: `docs/validation/sg6-gatelevel-consolidate.md`, `tools/gates/gate.js`, `tools/gates/lib/criteria/{g1,g5a,g7,g9}.js`. Implementation side: `internal/gates/policy.go` `Prepare`(:30)/`validateRequest`(:48)/`validateHookEffectPolicy`(:96); `internal/gates/runner.go` `Run`(:41)/`assess`(:115)/`scanOutput`(:172); `internal/gates/readiness.go` level closed set(:70–83)/`ProbeReadiness`(:113)/`AcceptanceReadiness`(:118); `internal/gates/probe_permit.go` marker state machine(:74–91)/`NewProbePermit`(:126); `internal/guard/lease.go` `Claim`(:82)/`Takeover`(:104)/`validTakeoverProof`(:120); `internal/guard/sandbox.go` `VerdictStatus`(:21–34); `internal/chain/postflight.go` `requirePostflightQualifiedReview`(:308); `internal/evidence/verify.go` `VerifyManifestObjects`(:37). Review action = independent re-check of "the independence of fail-closed preconditions": whether the precondition of each rejection path can be satisfied or bypassed by the same subject being guarded; whether "modify the policy" and "independently approve" are separated | `SECURITY.md:23` (two sentences in one paragraph: "the restriction must be enforced by runner/OS capability; prompts alone are not enough" and "the executor must not hold the combined permission to modify the policy and independently approve the same change"), §2 matrix row semantics; `CONTRACTS.md` §4–5 (Prepare/Inspect/Validate/Promote, handoff), §7 (admission is pure preflight, dispatch re-confirms that the authorization is not revoked and the reservation is outstanding, and any unsatisfied item is fail-closed); `TEST_STRATEGY.md` AT-05 (self-approval rejection) /AT-06/AT-18 (hard gates are not waivable) /AT-19; `ADR_REGISTER.md:58` uncovered scope (2) | Report "hard security gate face" subsection: Section A–E block + `FACE 2: PASS\|FINDINGS` |
| (3) Secrets and credentials face | `internal/adapters/ai_secret_windows.go` `NewPlatformAISecretManager`(:48)/`Set`(:52)/`Resolve`(:77)/`Exists`(:89)/`Delete`(:101)/`windowsCredentialTarget`(:139)/`containsASCIIControl`(:150); `internal/adapters/ai_provider.go` `ValidateAIProviderProfile`(:47)/`validHTTPSBaseURL`(:81)/`validDNSName`(:101)/`validAISecretReference`(:115); `internal/adapters/ai_probe_copilot.go` `Probe`(:60)/`copilotCLIProbeEnvironment`(:134)/`clearCopilotCLIProbeEnvironment`(:150); `internal/console/ai.go` `executeAICheck`(:33)/`configuredAIProfile`(:230); `docs/OPERATIONS.md` (`ai check` reads the bound SecretStore credential, `secret set` echoes nothing). Review action = independent re-check: secretRef opacity, credential injection/clearing lifecycle, whether credentials can reach configuration / summaries / logs / errors / evidence, whether the BYOK channel (Copilot CLI probe) passes credentials only through the child-process environment and clears them afterwards, and non-Windows fail-closed behaviour | `SECURITY.md` §4 (:56–:58: secrets never pass through models, are never written to configuration/logs, secret-direct, outbound scanning); `CONTRACTS.md` §7 secretRef section ("the profile stores only the provider type, base URL, model, transport options and an opaque secretRef; API keys/tokens must not enter configuration, summaries, logs, errors or evidence, and at runtime they can only be resolved by the system SecretStore and injected temporarily into the child-process environment"), the `windows-credential:<target>` section (`prfrail secret set` only after confirming the local console echoes nothing, `status` returns only the reference and existence, the adapter receives only a transient copy that is cleared afterwards, non-Windows platforms are fail-closed), "`config explain` must not show secretRef / credential existence / resolved value", "static validate/explain/preview do not read the SecretStore"; `TEST_STRATEGY.md` AT-10 (secret-direct is terminal-only, the secret canary never enters evidence/models) /AT-06 (a failed secret scan blocks); `DEV_PLAN.md:119` (implementation-complete sentence); `ADR_REGISTER.md:58` uncovered scope (3) | Report "secrets and credentials face" subsection: Section A–E block + `FACE 3: PASS\|FINDINGS` |

Line numbers inside the table are on-disk locating clues; when (7) cites them it must follow the §7.3 verbatim-location assertion (block header line number + verbatim transcription), and if either the location or the verbatim match fails, the item is treated as a "clue" and must not be relied on.

## 4. Review Method and Input Package

- The input package is registered in five groups A / B / C / D / E; item, line and byte counts follow the measurements in `input-list.md` (this report registers neither item nor byte totals; the per-item numbers are asserted by that list).
- Additional items: `00-briefing.md` (start message verbatim), `freedom-list.md`, the contractual fact columns of design §2 / §6, and `01-arch/package-manifest.md` (sanitization registration, see below).

| Group | Content | Items |
|---|---|---|
| A | Security baseline and requirements: `SECURITY{,_EN}`, `PRODUCT_REQUIREMENTS{,_EN}`, `ADR_REGISTER{,_EN}`, `DEV_PLAN{,_EN}`, `TEST_STRATEGY{,_EN}` | 10 |
| B | Execution directive and snapshots: `DELIVERY_DIRECTIVE{,_EN}`, `CHECKLIST_v65{,_EN}`, `REMAINING_SLICES{,_EN}` | 6 |
| C | Criterion sources: `sg6-gatelevel-consolidate.md`, `s1-exit{,_EN}.md`, `gate.js`, `g1/g5a/g7/g9.js` | 8 |
| D | Product-side hard security gate implementation: 4 files under `internal/gates/*`, 3 under `internal/guard/*`, `internal/chain/postflight.go`, 3 under `internal/evidence/*` | 11 |
| E | Secrets and credentials face: `ai_secret_windows.go`, `ai_provider.go`, `ai_probe_copilot.go`, `internal/console/ai.go`, `docs/OPERATIONS.md` | 5 |

Sanitization rules (verbatim filtering; "not found" is never a substitute for scanning):

1. Filter IP literals (IPv4 / IPv6, including private-range shapes).
2. Filter user names, host names, SSH target strings (of the form `git@<host>:`), `%TEMP%` and absolute paths outside the repository.
3. Filter credential / token / key / `password=` shapes.
4. Replace with `<redacted:category>` placeholders and keep the line count unchanged (same counting convention as `input-list.md`, so line-based locations do not drift).
5. Keep in-repository relative paths `docs/...` and `internal/...` (listed explicitly in §7.3); keep `tools/...` and `examples/...` after the controller confirms they are non-sensitive.
6. Land the sanitization result in `01-arch/package-manifest.md`: register per item the pre / post SHA-256, line count, byte count and filter hit count.

Materialization requirements ((7) has read-only tools and must not claim to have run commands):

- Every reading that (7) is asked to re-check must be pre-materialized by the controller as a static item inside the input package; (7) answers only from static items.
- Line-numbered excerpts: `SECURITY{,_EN}` §1–§2 / §4 and the `CONTRACTS.md` §7 secretRef section.
- Function-level line-numbered slices: implementation files are given as function headers + key rejection branches, never as whole-file pastes.
- Mechanical-reading static items: `go test ./...` / `go vet ./...` output, SELFTEST counts and the like are written to disk as `.txt` by the controller after the corresponding step runs and then packed.
- Citation convention: (7) may only write "per `<path>` inside the input package" and must not write "ran some command".

## 5. Independence Evidence, Output Contract and Conclusion Tiers

- Single judging subject: every review judgment is made by the (7) deep tier `GPT-6.2 Sol (copilot)` Max; the controller only orchestrates / materializes / transcribes / gates / receipts and never supplements, softens or rewrites any finding or conclusion line.
- Pipeline order and anchoring isolation: (5) pre-review (1 + 1) → (6) independent scan (≤3) → (7) blind final review; strictly serial, and the (5) / (6) lists never enter the first-round input of (7).
- (7) first-round blindness: the first-round input = the start message verbatim + the verbatim text of the ADR-047 row + the three-face definitions + the sanitized input package + the §2 criterion-source column and the §6 conclusion-tier template of the design; no (5) / (6) lists, no controller working notes, and no other design sections that carry directional content; only re-review rounds may attach the lists, labelled "read after completing the independent scan".
- Transcription contract: the final message of (7) is transcribed verbatim by the controller into the fenced block of the conclusion-summary section, with the Section A–E headings, face conclusion lines and composed conclusion line reproduced word for word; no rebuilding from bullet points, no blank-line reordering and no wording rework.
- Evidence trail: the input-package summary section of the report carries the blindness proof (whether lists were attached in round one, whether the input was sanitized, the MANIFEST pointer); metering registers every field from a single source of truth, and the (7) budget carries the 1 + 3 user-ruling annotation.

(7) emits one independent Section A–E block per face plus one composed block that carries no new findings:

| Section | Content |
|---|---|
| Section A | The face conclusion line `FACE <i>: PASS\|FINDINGS` (i = 1 / 2 / 3; two-valued only) |
| Section B | Findings list, each line fixed as `SEVERITY \| file:line \| violated sentence (verbatim quote) \| minimal fix`; `SEVERITY` ∈ {`Critical`,`High`,`Medium`,`Low`}; `file:line` must be locatable on disk; empty when there are no findings |
| Section C | Item-by-item evidence and reasoning, citing only static input-package material and contractual sentences |
| Section D | Falsifiability audit: for each item, "if this guard were deleted / reverted, which test or assertion turns red and which does not"; uncovered invariants are marked with a standalone `UNCOVERED` sentence, not counted as findings and not affecting the conclusion line |
| Section E | "What I cannot verify" + "what would overturn my PASS" |

Conclusion tier definition and composition rules:

- Each face conclusion tier = that face's Section A–E block + the face conclusion line `FACE <i>: PASS\|FINDINGS`; `PASS` only records "this supplement found no defect violating the criteria" and does not prove security; `FINDINGS` = at least one finding, each with `file:line` + minimal fix.
- First-read two-value constraint: this slice is an assessment-type first read with no fix round, so face conclusion lines are limited to `PASS\|FINDINGS`; if (7) emits an out-of-table form (such as `PARTIAL`), the controller does not rewrite it, transcribes the original wording and escalates to the user through the upgrade path.
- Consistency hard constraint: no findings in Section B ⇒ that face conclusion line must be `PASS`; findings in Section B ⇒ that face conclusion line must be `FINDINGS`.
- Composition rule (mechanical): all three faces `PASS` ⇒ `FINAL REVIEW: PASS`; any face `FINDINGS` ⇒ `FINAL REVIEW: FINDINGS`.
- Two qualifications must be registered together with the composition rule: **a PASS on any single face is not an overall PASS**; **an overall PASS is not "S2 may start"** (S2 start remains bound by other gates, and an overall PASS is not "security proven", with `UNCOVERED` items reproduced as-is).
- Format-deviation handling: missing five-section shape / missing face conclusion line ⇒ resend once (counted against the (7) budget); a further deviation ⇒ stop and escalate to the user, without self-repairing the format.
- Stage-6 transcription: the stage-6 final message is transcribed verbatim by the controller into the `sixstage` fence of the conclusion-summary chapter (Section A–E headings plus the block-final conclusion line); outside the fence only controller notes and pointers.

## 6. 审查结论摘要（⑥ 逐字转录）

```sixstage
## Section A 总体结论
INDEPENDENT SCAN: PASS

## Section B 发现清单
无

## Section C 逐条证据与推理
- CN/EN 报告均在 §6 建立 `sixstage`围栏：[CN :114](docs/validation/s2-pre-independent-security-review.md#L114)、[EN :114](docs/validation/s2-pre-independent-security-review_EN.md#L114)。A–E 标题各仅出现一次，Section A 与块末结论均为 `INDEPENDENT SCAN: PASS`。
- 行号问题已闭合：CN/EN 主文均引用 §4 `:56–:58`；原始错误的 `:37–40` 仅残留于 r1 发现复述，不是现行引用。
- 字节口径已闭合：`REMAINING_SLICES.md` 在两份静态件中均为 `381800`，未发现旧值 `381794`：[input-list :30](docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/input-list.md#L30)、[package-manifest :25](docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/01-arch/package-manifest.md#L25)。
- 消毒口径已闭合：逐件登记含行数列，并明确未改写时消毒前后 SHA 相同：[manifest :9](docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/01-arch/package-manifest.md#L9)、[说明 :60](docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/01-arch/package-manifest.md#L60)。
- 静态门禁读数报告 `G5 PASS`，`TOTAL_FAIL=0`：[gate readout :1](tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/gate-g7f2.txt#L1)、[汇总行 :6](tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/gate-g7f2.txt#L6)。

## Section D 可证伪性审计
- 删除或遗漏任一份 `sixstage`围栏：CN/EN 形状复核立即产生发现。
- 在主文重新引用 `SECURITY.md §4 :37–40`：当前行号锚点复核立即红。
- 将输入清单恢复为 `381794` 或使两份静态件漂移：字节一致性复核立即红。
- 删除行数列或消毒前后 SHA口径说明：逐件完整性复核立即红。
- 删除门禁读数或恢复非零 `TOTAL_FAIL`：门禁结论立即红。
- UNCOVERED：未读取 `04-reviews/`，未覆盖运行时、网络调用及真实生产环境。

## Section E 我无法验证的部分 + 什么会推翻我的结论
无法验证文档之外的运行面、真实生产安全结论及未读取的预审意见；静态门禁件未逐项列出 `G7-f/G7-a`。若后续证实围栏缺失、行号错误、字节值漂移、消毒登记缺失或门禁非零，本结论应改为 `INDEPENDENT SCAN: FINDINGS`。

INDEPENDENT SCAN: PASS
```

## 7. ⑦ 独立终审逐字转录（安全评审本体）

This section is a verbatim transcription of the (7) first-round final message; the original archive is `docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/07-security-r1.md`, and the controller reworded nothing.

```text
定点核对完成。三面都有能定位到源码和契约的发现；结论会区分已接线的 BYOK 路径和仍被上游门禁挡住的装配路径，并如实注明没有运行测试。

## 面 1 威胁模型
### Section A 结论
FACE 1: FINDINGS
### Section B 发现清单

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| Medium | [internal/gates/runner.go:208](internal/gates/runner.go#L208) | [docs/SECURITY.md:23](docs/SECURITY.md#L23)：「限制必须由 runner/OS 能力落实；不能仅靠提示词。」 | 将产物读取绑定到授权工作区目录句柄，在打开时拒绝越界；无法安全打开即拒绝收集，不入证据库。 |

### Section C 逐条证据与推理
读取 [docs/SECURITY.md:11](docs/SECURITY.md#L11)、[docs/SECURITY_EN.md:11](docs/SECURITY_EN.md#L11) 和 [docs/TEST_STRATEGY.md:29](docs/TEST_STRATEGY.md#L29)：双语基线一致区分了目录分离与进程隔离、内容完整性与身份认证。STRIDE 已列出路径逃逸，缺口在于 SEC-02 承诺的打开时检查未覆盖所有读取入口。读取 [internal/gates/runner.go:194](internal/gates/runner.go#L194)：产物收集按路径展开，只对末端做 `Lstat`，再按路径重新打开。它既未约束中间目录解析，也未把检查与读取绑定到同一对象。现有 `GuardExecutor` 缺必需能力而拒绝执行，这限制了当前可达性，不补足该读取边界。
### Section D 可证伪性审计
- 删除 **BoundDir** 的中间路径限制，`TestBoundDirMiddleSymlinkRefusedWithClosedErrno` 会失去预期拒绝；它不经过 Gate 产物收集。删除证据对象摘要校验，`TestVerifyManifestObjectsFailsClosed` 的 `tampered` 子项会红；删除事件归属检查，`TestVerifyManifestStateEventBinding` 的跨 run 子项会红。
- 删除 C6 的魔数检查或未提交槽判定，`TestC6SlotValidationOrderBadMagicIsCorruption` 或 `TestC6SlotValidationOrderUncommittedPayloadIsAbsent` 会红。基线已明示装配记录丢失后的漂移检测降级，不把它记为完整性已证明。
- UNCOVERED：未见 Gate 产物读取跨目录边界时同时断言拒绝、证据库无越界对象、任务不推进。现有 AT 编号不能代替这些断言。
### Section E 我无法验证的部分 + 什么会推翻我的判断
未验证原生平台上的完整 STRIDE 运行证据。若能指出所有 `Runner` 产物读取调用都由不可绕开的根目录绑定保护，并有针对这条读取路径的反例，才会推翻本项。

## 面 2 硬安全门禁
### Section A 结论
FACE 2: FINDINGS
### Section B 发现清单

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| High | [internal/console/runtime_agent.go:1245](internal/console/runtime_agent.go#L1245) | [docs/CONTRACTS.md:212](docs/CONTRACTS.md#L212)：「生产 admission 是纯 preflight 与不可变绑定：不得消费 request replay 身份、不得发布 R/C 或写 replay store、不得启动进程；首次 dispatch 资格与 replay 身份消费由 dispatch 路径在发布 R 之后依据 replay store 判定；admission 的结论是评估时刻的 preflight verdict，不是锁，可能在 dispatch 前失效，dispatch 在发布 R 之后、允许 launch 之前必须重新确认授权未被撤销且预算 reservation 仍 outstanding，任一不满足即 fail-closed，不得把先前的 admission 通过当作 launch 许可。」 | 在 launch 前向可信账本读取当前授权与预留；读取失败、撤销或已结算均保留 R，拒绝启动。 |
| High | [internal/chain/review.go:116](internal/chain/review.go#L116) | [docs/CONTRACTS.md:147](docs/CONTRACTS.md#L147)：「评审者与候选变更产生者是否同一身份、`evidenceRootHash`/`policyHash` 引用是否真实存在及 waiver 到期后不得复用，均由独立 checker 判定。」 | 发布前从执行者不可写的可信来源验证策略批准与候选证据根；没有有效绑定则拒绝，或停在待人工批准。 |
| Medium | [internal/console/runtime_agent.go:731](internal/console/runtime_agent.go#L731)、[internal/console/runtime_supervised.go:1642](internal/console/runtime_supervised.go#L1642) | [docs/CONTRACTS.md:150](docs/CONTRACTS.md#L150)：「`promotion-receipt.schema.json` 将父/候选 snapshot、evidence root、review receipt、停机与租约证据绑定为单个发布事务。」 | 发布前核对当前工作区的排他租约或等价写入隔离证明，将其真实记录绑定到本次 attempt；缺失或不一致则拒绝。 |
| Medium | [internal/gates/runner.go:73](internal/gates/runner.go#L73) | [docs/CONTRACTS.md:391](docs/CONTRACTS.md#L391)：「只读/本地可丢弃/外部写分离；S1 外部写拒绝，未知副作用暂停对账；诊断不变更锁/日志」 | 执行前从冻结 hook 策略导出有效副作用分类；分类缺失、矛盾或不允许时不能进入 `Execute`。 |

### Section C 逐条证据与推理
- **撤销独立性。** 读取 [internal/console/cli.go:562](internal/console/cli.go#L562) 和 [internal/adapters/agent_runner_replay_dispatcher.go:150](internal/adapters/agent_runner_replay_dispatcher.go#L150)：dispatcher 在正确的位置调用复核，但生产回调返回 CLI 启动时读取的授权副本。复核可发现随时间过期，不能发现初读后另行持久化的撤销。预算也只接到已装载对象；未证明可观察跨进程结算。
- **策略与批准。** 读取 [internal/console/runtime_agent.go:678](internal/console/runtime_agent.go#L678)、[internal/console/runtime_supervised.go:1598](internal/console/runtime_supervised.go#L1598) 和 [internal/chain/engine.go:254](internal/chain/engine.go#L254)：两条装配路径在 postflight 合格后本地计算策略摘要并返回 `approve`。同名拒绝是真守卫，但摘要合形不证明策略事先获独立批准；postflight 与批准不能互作原件。
- **租约与副作用。** 读取 [internal/guard/lease.go:104](internal/guard/lease.go#L104) 和 [internal/console/runtime_supervised.go:1491](internal/console/runtime_supervised.go#L1491)：租约库缺独立 verifier 会拒绝接管，但发布端的 `LeaseEvidence` 只由 run ID 派生，不证明当前独占。归还 hook 接线未填 `Effect`，而 runner 仅在有分类记录时调用准入。现版 CLI 另有 hook 定义和 runner 能力拒绝，不能把这个缺口说成现版已经执行了外部写入。
### Section D 可证伪性审计
- 删除能力检查，`TestPrepareRejectsUnavailablePolicy`、`TestRunnerDoesNotExecuteWhenPolicyIsUnavailable` 会红；删除启用后的沙箱拒绝，`TestAgentRunnerRunSandboxRefusesBeforePermit`、`TestSupervisedRunSandboxRefusesBeforePermit` 会红。契约明订的默认 `off` 不算绕过启用后的检查。
- 删除同路径独占或损坏 marker 拒绝，`TestProbePermitCrossProcessSingleClaim`、`TestProbePermitDamagedMarkerStaysSpent` 分别会红。删除就绪级别阻断，`TestReadinessProbeConclusionCannotAdvanceAcceptance`、`TestReadinessProductionRequiresAT23` 会红；这些是级别和格式断言，不是证据签发者认证。
- 删除同名自批、过期 waiver、postflight 事实绑定或接管 verifier 检查，分别由 `TestRejectSelfReview`、`TestRejectExpiredOrWrongCandidateWaiver`、`TestPostflightGateRefusesAnUnprovenPass`、`TestLeaseTakeoverFailsClosedWithoutVerifier` 反证。
- UNCOVERED：生产接线缺少这四项反例：初读后持久化撤销且 **R 保留、零启动**；策略未获独立批准且 **零发布**；租约证明缺失且 **无完成回执**；副作用分类缺失且 **零执行**。账本读取失败测试和显式外部写分类测试覆盖的是相邻分支。
- UNCOVERED：探针测试证明同一路径只认领一次，不证明独立签发或不同路径不能代表同一授权；`--probe-permit` 的调用方信任边界不能用 marker 的存在代证。
### Section E 我无法验证的部分 + 什么会推翻我的判断
未验证独立控制面是否在这两条装配路径外实施了不可绕过的实时账本、策略批准、排他写入或 hook 分类。能反驳上述判断的是实际接线与对应的零启动、零发布或零执行断言。未交付的路径仍应保持阻断。

## 面 3 秘密与凭据
### Section A 结论
FACE 3: FINDINGS
### Section B 发现清单

| SEVERITY | file:line | 被违反的句子（逐字引用） | 最小修复 |
|---|---|---|---|
| High | [internal/console/ai.go:94](internal/console/ai.go#L94) | [docs/SECURITY.md:60](docs/SECURITY.md#L60)：「网络默认拒绝，例外授权须列 host/port/purpose、提供方数据策略及期限。」 | 解析凭据前，由可信本机控制路径校验已批准的 profile、凭据引用、接收端与子进程身份绑定；不符即不读库、不启动。 |
| Medium | [internal/guard/process_windows.go:457](internal/guard/process_windows.go#L457) | [docs/CONTRACTS.md:216](docs/CONTRACTS.md#L216)：「adapter 只接收解析后的短暂字节副本，使用后清零；系统错误只返回稳定分类，不得拼接 credential blob。」 | 对 Windows 启动所用的原生环境块及编码中间缓冲保留所有权；成功或失败均及时清零，并核对凭据读取与终端输入的临时副本。 |
| Medium | [internal/adapters/ai_probe_copilot_executor.go:53](internal/adapters/ai_probe_copilot_executor.go#L53) | [docs/CONTRACTS.md:214](docs/CONTRACTS.md#L214)：「profile 只保存 provider 类型、base URL、模型、传输选项及不透明 `secretRef`；API key/token 不得进入配置、摘要、日志、错误或证据，运行时只能由系统 SecretStore 解析并临时注入子进程环境。」 | 只保留经过校验的结构化事实，不对未脱敏流生成持久证据；关闭未证明脱敏的子进程落盘日志，无法保证则拒绝 BYOK 探针。 |

### Section C 逐条证据与推理
- 读取 [internal/adapters/ai_provider.go:115](internal/adapters/ai_provider.go#L115)、[internal/adapters/ai_secret_windows.go:139](internal/adapters/ai_secret_windows.go#L139) 和 [internal/console/ai.go:66](internal/console/ai.go#L66)：profile 验证只检查引用形状和 HTTPS 语法，Windows 后端按当前用户的目标名解析。`ai check` 的配置选择和可执行文件路径没有独立核对凭据与接收方的批准绑定。单请求授权、HTTPS 和配置摘要均不能证明接收方有权持有凭据。
- 读取 [internal/adapters/ai_preflight.go:57](internal/adapters/ai_preflight.go#L57)、[internal/adapters/ai_probe_copilot.go:80](internal/adapters/ai_probe_copilot.go#L80) 和 [internal/guard/process_windows.go:221](internal/guard/process_windows.go#L221)：SecretStore 返回的 Go 切片、探针环境副本及 executor 的字节缓冲有清零；真正传给 Windows 的 UTF-16 环境块及转换副本没有。`CredFree` 调用只能据实记为释放，不能据此认定擦除。
- 读取 [internal/adapters/ai_probe_copilot_executor.go:92](internal/adapters/ai_probe_copilot_executor.go#L92) 和 [internal/adapters/ai_secret_unsupported.go:22](internal/adapters/ai_secret_unsupported.go#L22)：ProofRail 未把解析值放进探针 argv，非 Windows SecretStore 拒绝解析，结构化原因也不转录原始错误。但子进程可写临时日志，原始 stdout/stderr/usage 未经秘密检查即入摘要，目录清理不能证明从未落盘。未观察到实际泄漏。
### Section D 可证伪性审计
- 去掉 SecretStore 返回副本或探针调用环境副本的清零，`TestAIAvailabilityCheckUsesOnlyRequestedChannelAndClearsSecret`、`TestCopilotCLIProbeBuildsIsolatedDeepSeekSmoke` 分别会红。去掉缺秘密时的拒绝，`TestAIAvailabilityCheckBlocksMissingSecretWithoutProbe` 会红。去掉 explain 遮蔽，`TestAIConfigRoundTripAndExplainRedactsSecretReference` 会红。
- UNCOVERED：未见生产入口对未批准的凭据/接收方绑定断言零解析、零启动；也未见覆盖原生环境块的成功及失败清零、真实平台后端拒绝后零启动，以及合成秘密在子进程输出、日志及最终证据中的端到端不泄漏检查。
### Section E 我无法验证的部分 + 什么会推翻我的判断
未验证 Windows 原生调用和子进程自身的日志行为。可推翻相应发现的是生产入口已有不可绕过的批准绑定、原生缓冲的完整清零路径，或能证明日志与摘要只消费脱敏事实的端到端反例。

## 合成
FINAL REVIEW: FINDINGS

读取 [docs/ADR_REGISTER.md:58](docs/ADR_REGISTER.md#L58) 与 [docs/PRODUCT_REQUIREMENTS.md:109](docs/PRODUCT_REQUIREMENTS.md#L109)：三面各有不符合判据之处，不能登记为安全门通过，也不能以此替代 S1 Exit 评审。以上是只读判断，未执行命令、未变更文件。方案书定位回显夹带了许可节之外的片段，未采信；未读取 ⑤/⑥ 清单或 P79 报告结论。
```

## 8. 三面结论与合成

| Face | Conclusion tier | Findings | Highest severity |
|---|---|---|---|
| Face 1 threat model | `FINDINGS` | 1 | `Medium` |
| Face 2 hard security gates | `FINDINGS` | 4 | `High` |
| Face 3 secrets and credentials | `FINDINGS` | 3 | `High` |
| **Composed** | `FINAL REVIEW: FINDINGS` | 8 | `High` |

`FINDINGS` on any face ⇒ the composed result is `FINDINGS`; an overall `FINDINGS` does not mean "S2 must not start", but the ADR-047 compensation duty is discharged by "review completed and conclusion registered", see section 9.

**Re-review chain**: re-review round 1 `RE-REVIEW: FINDINGS` (1 Medium: the tense on line 7 of `SECURITY_EN` was not restricted to a historical state) → fix (the English line 7 became `Before P79 … had never been executed`, and the Chinese line 7 gained the same "before P79" qualifier) → re-review round 2 (narrowed input) `RE-REVIEW: PASS` (0 new findings); stage-7 budget 1 final review + 3 re-reviews, actual 1 + 2 (the third re-review was not used).

## 9. 发现处置（本片不改产品代码）

- This slice is assessment type; its boundary explicitly forbids changing product code and the directive.
- All 8 findings (3 High + 5 Medium) are archived as-is, each with a "suggested follow-up" column.
- Whether to open fix slices, whether to revise the directive, and whether to flip ADR-047 to "compensated" are all ruled by the user; the controller must not decide them on its own.
- No code was changed, so "whether the findings were fixed" is outside this slice's acceptance scope; acceptance = review completed + conclusion registered truthfully.

| Id | Face | Severity | file:line | Suggested follow-up |
|---|---|---|---|---|
| F1 | Face 1 | `Medium` | `internal/gates/runner.go:208` | separate fix slice: bind artifact reads to an authorized workspace directory handle and refuse collection when it cannot be opened safely |
| F2 | Face 2 | `High` | `internal/console/runtime_agent.go:1245` | separate fix slice: read current authorization and reservation from a trusted ledger before launch; on read failure, revocation or settlement keep R and refuse to start |
| F3 | Face 2 | `High` | `internal/chain/review.go:116` | separate fix slice: verify policy approval and candidate evidence root from a trusted source the executor cannot write |
| F4 | Face 2 | `Medium` | `internal/console/runtime_agent.go:731`, `internal/console/runtime_supervised.go:1642` | separate fix slice: check the exclusive lease or equivalent write-isolation proof and bind it to this attempt |
| F5 | Face 2 | `Medium` | `internal/gates/runner.go:73` | separate fix slice: derive a valid side-effect classification from the frozen hook policy before execution |
| F6 | Face 3 | `High` | `internal/console/ai.go:94` | separate fix slice: validate the profile, credential reference, recipient and child-process identity binding before resolving credentials |
| F7 | Face 3 | `Medium` | `internal/guard/process_windows.go:457` | separate fix slice: own the native environment block and encoding buffers and clear them on both success and failure |
| F8 | Face 3 | `Medium` | `internal/adapters/ai_probe_copilot_executor.go:53` | separate fix slice: disable child-process disk logs not proven sanitized, or refuse the BYOK probe when that cannot be guaranteed |

## 10. 文档同步

| Object | Disposition | Basis |
|---|---|---|
| `SECURITY{,_EN}` note line | updated to "completed + date + reviewer" | deliverable (2) of this slice; the annotation duty in the ADR-047 affected surface |
| ADR-047 row of `ADR_REGISTER{,_EN}` | P79 completion note added, **not flipped to "compensated"** | user ruling; a status-word change exceeds this slice's authorization |
| section 7 of `PRODUCT_REQUIREMENTS{,_EN}` | **no annotation added** | listed as a check item in the start message; the user rules on whether to annotate |
| P79 section of `REMAINING_SLICES{,_EN}` | step-0 facts registered and the completion record appended | deliverable (4) of this slice |

## 11. 审查输入包摘要与盲审隔离证明

- Input package = five groups A / B / C / D / E plus additional items; per-item line counts, bytes, SHA-256 and filter hit counts are asserted by `input-list.md` and `01-arch/package-manifest.md`.
- The sanitization conclusion points to `01-arch/package-manifest.md`: per-item pre / post SHA-256 registration, with identical SHA for items that were not rewritten.
- (7) first-round blindness = no (5) / (6) lists attached and no controller working notes attached; the first-round input = the start message, the ADR-047 original text, the three-face definitions, the sanitized input package, and design sections 2 and 6.
- (7) was read-only and ran no commands: every reading is cited as "per `<path>` inside the input package".
- (7) transcription is verbatim: the section 7 fence matches the archive `04-reviews/07-security-r1.md` line by line.

## 12. 可证伪性与变异检验

- This slice is documentation / assessment type with no product code change; the falsifiable objects are documented facts and gate readings.
- Version metadata: if the directive version, SELFTEST count or three-scope readings disagree with reality, the step-0 anchor re-check and the (10) receipt re-run turn red immediately.
- Bilingual shape: if CN / EN line counts, `##` heading counts or per-line pipe and bold-marker counts diverge, the G3 mirror-symmetry check turns red immediately.
- Fences: if the `sixstage` / `text` / `artifacts` fence is missing or unclosed, the fence check turns red immediately.
- Artifact fence: if a fence path does not match a file on disk or omits a slice artifact, the artifact-fence check turns red immediately.
- Reference resolvability: if a cross-reference (sections 7 to 9, sections 14 to 15) points at a nonexistent section, the reference check turns red immediately.
- Encoding: if any file is not UTF-8 without BOM plus LF, the encoding gate turns red immediately.
- No product-level mutation testing was performed: this slice changes no product code, so there is no guard or assertion to mutate.

## 13. 门禁结果

Time point: the readings in this chapter are taken when stage (8a) finalizes; the final change-set gate readings are the ones in the (10) receipt.

| Scope | Command | TOTAL_FAIL | Note |
|---|---|---|---|
| `index` | `node tools/gates/gate.js --all --scope=index` | 0 | red items cleared (G5(d) / G5(f) / G7-a / G7-f / G8-b all PASS); readout `gate-index-final.txt` |
| `tree` | `node tools/gates/gate.js --all --scope=tree` | 0 | same caliber; readout `gate-tree-final.txt` |
| `ci` | `node tools/gates/gate.js --all --scope=ci` | 0 | same caliber as CI (degrades to the previous-commit range when no CI variables are present locally) |
| `SELFTEST` | `node tools/gates/gate.js --selftest` | — | all green (exit 0) |

## 14. 环境与档位

| Role | Carrier | Model | Tier | Set at | Paused |
|---|---|---|---|---|---|
| Controller | `sol-orchestrator` | `DeepSeek V4.1 Flash (deepseek)` | High | set before start | no |
| (1) Architect | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | Max | set once before start | no |
| (5) Pre-reviewer | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | Max | set once before start | no |
| (6) Independent scan | `independent-reviewer` | `MAI Code 2 Flash (copilot)` | default | set once before start | no |
| (7) Independent final review | `independent-reviewer` | `GPT-6.2 Sol (copilot)` | Max | set once before start | no |
| (8) Documenter | `prfrail-documenter` | `DeepSeek V4.1 Flash (deepseek)` | High | set once before start | no |

The four fields (section 8.4) = set-at time, tier value, whether paused, and gate results before / after a pause; this slice was not paused, so those gate results are not applicable. The section 8.3 same-model restriction does not apply ((7) and the controller use different models).

## 15. 成本与计量

| Role | Model | Rounds | Note |
|---|---|---|---|
| (1) Architect | `DeepSeek V4 Pro (deepseek)` | 1 | design `01-arch/design-r1.md` |
| (5) Pre-reviewer | `DeepSeek V4 Pro (deepseek)` | 2 | budget 1 + 1 exhausted |
| (6) Independent scan | `MAI Code 2 Flash (copilot)` | 2 | budget 3 used 2; r2 `INDEPENDENT SCAN: PASS` |
| (7) Independent final review | `GPT-6.2 Sol (copilot)` | 1 | budget 1 final review + 3 re-reviews; 1 final review used in this round |
| (8) Documenter | `DeepSeek V4.1 Flash (deepseek)` | 1 | this report pair |
| Batch assist | `GLM-5.3-Flash (glm)` | 2 | enumeration / excerpt-and-number checks / bilingual mirror checks / MANIFEST drafting |

The (7) budget = 1 final review + 3 re-reviews; that budget exceeds the section 7.5 literal standard "M-L counts as M", is a direct user ruling of 2026-11-19, and **must not be used as a precedent**. No blind review was sampled: the three semantics that section 7.7 mandates (ownership / stop / identity) are not touched by this slice. Batch assist ran twice; probes = 0.

## 16. 明确未执行事项

- No product code changed.
- No real external call executed.
- No Ubuntu VM run.
- No ATT-level runtime security experiment.
- None of the findings listed in section 9 was fixed.
- No annotation added to `PRODUCT_REQUIREMENTS` (the user rules on that).
- ADR-047 not flipped to "compensated" (the user rules on that).

## 17. 先例禁用核对

This slice cites the user ruling on "(7) budget 1 + 3 re-reviews" and records truthfully that it must not be used as a precedent; it cites neither the OB-11 exception nor any other clause marked as not citable as a precedent.

## 18. ⑧a 自查

- (1) Mechanical audit of the cost table: present; evidence = the section 15 table and the dispatch line of each round under `docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/`.
- (2) "Disposed" `file:line`: absent; no item in this slice is marked "disposed".
- (3) Old-caliber grep: present; command and result are in `tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/old-caliber-grep.txt`.
- (4) Every Info / Low item of each (6) / (7) round lands in the report: present; cross-checked against `04-reviews/06-scan-r1.md`, `06-scan-r2.md` and `07-security-r1.md`, none of which carries an Info / Low item.

## 19. A-NN 编号完整性

This report uses no `A-NN` numbering; F1–F8 are the finding identifiers of section 9, defined in the section 9 table.

## 20. ⑧c 收尾清洁记录

- Whole-tree per-file encoding check: every violation sits in a frozen evidence pack under `docs/validation/evidence/**`, and the non-frozen main tree is zero (same caliber as the last row of section 13; frozen packs are checked for self-consistency per section 5.1 and are never re-encoded).
- `tmp/` assertion: the scripts and intermediate readouts inside `tmp/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/` were removed item by item per section 5.1, and afterwards the `tmp/` top level holds only `.gitkeep` plus the runtime-cache container `gates-cache` (the exemption set of section 5.1).
- Untracked items: the `??` entries of `git status --porcelain -uall` are 0; the whole change set is in the index.
- Runtime-cache exemption: `tmp/gates-cache/g4b/` is an enumerated member of the section 5.1 exemption set and is not cleared by (8)c.

## 21. Exceptions and Fallback Record (R / D)

| ID | Stage | Root cause | Handling | Counted |
|---|---|---|---|---|
| R-1 | (5) | Pre-review r1 reported 2 Medium + 4 Low (missing `sixstage` contract, undisclosed gate trajectory, line-number drift, etc.) | Fixed the design and the report, then re-ran (5) in place to r2 | R = 1 |
| R-2 | (6) | Independent scan r1 reported 2 High + 3 Medium (missing fence, line-number drift, byte mismatch between static artifacts, incomplete sanitisation register) | Fixed, then re-scanned (6) to r2 `INDEPENDENT SCAN: PASS` | R = 2 |
| — | (7) | Re-review round 1 reported 1 Medium (`SECURITY_EN` tense); a re-review is an in-budget obligation round and is not counted as R under section 9.2 | Fixed, then narrowed re-review returned `PASS` | not counted |

- **R measured = 2** (one for the stage-5 round, one for the stage-6 round; the stage-7 re-review round is not counted under section 9.2). This reaches the `R>=2` threshold of section 9.2, which calls for marking the slice "high risk"; the marker was not set explicitly, **recorded here as a deviation of this slice** (every stage required by the high-risk row of section 5.2 — (1)(2)(3)(5)(6)(7)(9) — was actually executed, so no stage was skipped in substance), submitted to the user for the (10) review ruling.
- **D measured = 1**: after the (8)c clean-up, the controller mistakenly created `tmp/P79-fix/` (not `tmp/<slice ID>/`) to hold the scripts for this correction; the whole directory was deleted immediately with zero residue (`tmp` top level = `.gitkeep` + `gates-cache`; the recursive file count excluding `gates-cache/g4b` is 1). Registered as **D = 1** under the section 5.1 self-heal grade; since the controller itself was the violator, the **self-heal verdict awaits the same-turn independent review layer or the user**, and is not deemed established before that confirmation.

## 工件围栏

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/PRODUCT_REQUIREMENTS.md	repo
docs/PRODUCT_REQUIREMENTS_EN.md	repo
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/00-briefing.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/01-arch/design-r1.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/01-arch/package-manifest.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/05-prereview-r1.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/05-prereview-r2.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/07-re-review-r1.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/07-re-review-r2.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/04-reviews/07-security-r1.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/06-gates/gate-enc-scan.txt	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/06-gates/gate-index-final.txt	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/06-gates/gate-selftest-final.txt	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/06-gates/gate-tree-final.txt	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/MANIFEST.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/SHA256SUMS.txt	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/freedom-list.md	repo
docs/validation/evidence/P79-S2-PRE-INDEPENDENT-SECURITY-REVIEW/input-list.md	repo
docs/validation/s2-pre-independent-security-review.md	repo
docs/validation/s2-pre-independent-security-review_EN.md	repo
```
