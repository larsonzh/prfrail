# P95 GATES-CRITERIA-ASYNC validation report (tooling slice · assessment-type closure)

- Slice: P95 `GATES-CRITERIA-ASYNC` (M-tier tooling `[SLICE]`, original goal = asynchronising G6 / G7 / G8); lead = `DeepSeek V4 Pro (deepseek)` Max (user ruling 2026-11-27, not a precedent).
- Date: 2026-11-27. Directive version: v1.77 (**no version bump**). R = 0, D = 0 (start baseline). Probes 0.
- Evidence pack: `docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/` (start message, (1) design, the ⑤ / ⑥ / ⑦ review transcripts, batch-helper records, the freedom list, the input list, `MANIFEST.md` and `SHA256SUMS.txt`).
- Shape: **assessment-type closure** - the (1) benefit gate ruled against doing it, and on 2026-11-27 the user ruled for candidate A (that ruling **must not serve as a precedent**); criterion code untouched.

## 1. Conclusion

This slice does not asynchronise G6 / G7 / G8. The (1) impact assessment and benefit gate returned a rejection: the ctx read APIs are fully synchronous, so under a single Node process asynchronisation cannot overlap synchronous I/O, and the theoretical upper bound of the concurrency benefit is ≈0 while the practical bound is ≤250 ms out of a 5.8 s wall clock. The user ruled for candidate A, so the slice degrades to an assessment-type closure: **criterion code untouched, directive not bumped (v1.77 unchanged)**; P88 report section 7 recommendation 4 is registered in this slice as "**assessed, benefit does not hold, not implemented**"; the two `CHECKLIST_v78` files are untouched and will be absorbed by the v79 snapshot. Current state = pushed to `origin/main`, CI run `45939437180` green on both legs (see section 10); after the backfill commit it parks at step (10) (origin only).

## 2. Change set and evidence

- Zero changes on the criterion implementation surface: `tools/gates/lib/criteria/g6.js`, `tools/gates/lib/criteria/g7.js`, `tools/gates/lib/criteria/g8.js`, `tools/gates/gate.js`, `tools/gates/lib/ctx.js` and `tools/gates/selftest.js` are all byte-identical to HEAD `22370fe`, and the slice change set contains none of them.
- Ledger pair: `docs/t027/REMAINING_SLICES.md` and `docs/t027/REMAINING_SLICES_EN.md`; the P95 entry moves from "not started" to "in progress" and to "closed" on close-out with a completion record.
- Report pair: this file and its `_EN` mirror.
- Evidence pack: 9 payloads plus the two metadata files `MANIFEST.md` and `SHA256SUMS.txt` (entry set equal to the latter).
- Untouched items: the directive pair, `ci.yml`, `internal/**`, CONTRACTS, `schemas/`, `testdata/`, workflows, the G5-b closed set and model whitelist, and the ADR register - none of them is in this slice's change set.

## 3. Benefit argument (transcribed from (1))

### 3.1 Serial composition (measured table plus net self-cost estimate)

Reading caliber: Windows with 8 logical cores, Node v24.17.0, warm corpus, two rounds of a single sample.

| Criterion / probe | `scope=tree` ms | `scope=index` ms |
|---|---|---|
| G1-a (fixed-overhead probe) | 1247 | 1177 |
| G4a | 1675 | - |
| G4b | 1596 | - |
| G6 | 1282-1493 | 796 |
| G7 | 610-642 | 471 |
| G8 | 1803-1905 | 1372 |
| `--all --skip-base` total | 5799 | 3812 |

Net self-cost estimate (after deducting the fixed probe overhead, clamped at 0 when below the probe): the three criteria cost ≈0.6-0.9 s on tree and ≈0.2 s on index, i.e. 10-15% / ≈5% of their respective wall clocks. Under `--all` with base the share is below 1% (P88 warm base ≈94 s). Nature of the work: G6 / G7 are pure-CPU regular-expression line scans and G8 is sha256 hashing plus set comparison; the hot path still uses synchronous git child processes (one per line in G7-b, one batched ignore check per enumeration in G8), which asynchronisation cannot overlap.

### 3.2 Concurrency benefit ceiling (the ≈0 theoretical upper bound)

Theoretical upper bound ≈0: every ctx read API is synchronous (`readFileSync` / `spawnSync`), so no asynchronous I/O is available to overlap; under a single-threaded Node process asynchronisation only interleaves cooperatively and total CPU is conserved, so wall clock = sum of self-cost + await microtask overhead. Even if the ctx read APIs were also moved to async, the bound = async-I/O share of the three criteria x (sum - max): most of the tree net self-cost is regex and sha256 CPU time, leaving at most ≈250 ms overlappable, and ≈0 on index. The practical saving is ≤0.3-4% of the `--all --skip-base` wall clock and ≈0 for `--all` including base. Negative possibilities: microtask scheduling overhead, interleaving destroying warm-cache locality (tens of ms either way), and fast failure altering the current synchronous sequential-abort behaviour.

### 3.3 Concurrency cost (shared-cache concurrent-read safety / report-order preservation / new failure shapes)

- **Lazy-cache concurrent-read safety**: `bytesCache` / `existsCache` / `trackedCache` / `numstatCache` / `showCache` / `ignoreCache` / `hunkCache` / `treeCache` are all function-closure Maps whose writes complete atomically inside synchronous calls, and ctx has no async function, so no half-written window can straddle an await; written values are idempotent (same key, same value).
- **Shared Buffer**: the byte reader returns one and the same Buffer reference and every current caller only reads it (BOM counting, sha256, decode-to-string). Safe today, but it is an unhardened assumption; after any concurrency change the comment discipline "callers must not mutate the Buffer" must be kept.
- **WeakMap**: the G8 pack-root cache is keyed by ctx and used by G8 alone, so there is no cross-criterion sharing problem.
- **Report order**: the scheduler already back-fills each criterion's own output buffer in plan order; the three `run` functions write only their own `out`, with zero `console` / `process.stdout` inside criteria; the single-criterion path already awaits the `run` return value and is Promise-compatible.
- **No cross-criterion side-effect dependency**: the criteria directory assigns nothing onto ctx, the only disk write belongs to the G4b cache (outside this slice), and G6's module-level reuse of G7 is read-only constants and functions.
- **New failure shapes**: an async throw propagates through fast failure and changes the existing synchronous sequential-abort shape (exit-code semantics unchanged); plus microtask scheduling overhead and lost warm-cache locality.

### 3.4 Conclusion and alternatives

- **Conclusion**: do not do it. First, the theoretical upper bound is ≈0; second, introducing new failure shapes and deep-water review cost for ≈0 benefit fails the directive's benefit standard; third, the real wall-clock bulk is the four base items (P88 report section 7 recommendation 2) and the G4b cold scan (OB-113 sub-item 3), not G6 / G7 / G8.
- **Alternatives**: keep the status quo; register P88 report section 7 recommendation 4 as "assessed, benefit does not hold, not implemented" (closed by this slice); if real parallelism is ever needed, the candidate path is worker_threads (cost: git child processes moved into workers, cross-thread cache consistency, change-set serialisation; high risk with a ceiling of only ≈0.9 s, so it does not hold either).

## 4. Impact determination (CONTRACTS/schema/fixtures/internal/**/ci.yml all untouched; under candidate A the section 6.3 concurrency sentence is untouched and no version bump is needed)

| Item | Touched if implemented | Basis |
|---|---|---|
| CONTRACTS | No | gate scripts are outside that domain and `docs/CONTRACTS.md` has no gate clause |
| `schemas/` | No | no gate-related definition |
| `testdata/` and fixtures | No | not involved; `selftest.js` is code, not a fixture |
| `internal/**` | No | the change stays on the criterion implementation surface, and earlier precedents all excluded `internal/**` |
| `ci.yml` / workflows | No | the CI invocation surface and exit-code contract are unchanged |

- Rule text: the criterion concurrency-scheduling sentence in section 6.3 and the OB-113 disposition note in section 12.4 would have to be updated with the slice and bumped to MINOR v1.78 **only on the implementation branch**; under candidate A neither is touched, so **the directive stays at v1.77**.
- (7) tier: the implementation branch would fall into the deep-water round, while candidate A returns to the standard tier - the user's ruling is `GPT-5.3-Codex (copilot)` with a quota of 1 final review + 1 re-review.
- No blind review draw; probes 0; no new ADR; no gitee.

## 5. Review chain and closure

### 5.1 (5) pre-review

(the (5) verbatim output is archived at `docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/05-pre-review-r1.md`.)

(5) r1 = `PRE-REVIEW: PASS` (0 High / 0 Medium / 2 Low + 3 Info). Both Lows are fixed: line 24 gains the clamp wording "below the probe counts as 0", and the head lead line plus the ledger quota row gain the "user ruling, not a precedent" marker. The 3 Info items are recorded (`~0.9s` vs `≈0.9 s` glyph form, the 9-payload count, the P88 calibre-incomparability note). No Medium+, no re-review needed.

### 5.2 (6) independent scan

(the (6) verbatim output is archived at `docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/06-scan-r1.md`; the fence below is a structured transcript: text preserved verbatim, the verdict line placed at the fence end per section 11.3.)

```sixstage
### Section A
INDEPENDENT SCAN: PASS

### Section B
| SEVERITY | file:line | 被违反的句子或主张 | 最小修复 |
|---|---|---|---|
| — | 无可执行发现 | 无可执行发现 | 无可执行发现 |

### Section C
- S1 PASS — 暂存变更集 15 项，全为 `docs/**`；`git diff --cached -- tools/gates` 无输出，无判据或准则实现变更。[中文报告:12-16](docs/validation/gates-criteria-async.md#L12-L16)、[英文报告:12-16](docs/validation/gates-criteria-async_EN.md#L12-L16)。
- S2 PASS — 两份报告保持“不做”、理论上界 ≈0、五项影响全否、候选 A 不升版 v1.77；用户裁定的 ⑦ 常规档 1+1、⑥ 1 次及不得作先例均守恒。[结论证据](docs/validation/gates-criteria-async.md#L10)、[影响证据](docs/validation/gates-criteria-async.md#L56-L67)、[台账口径](docs/t027/REMAINING_SLICES.md#L2539-L2553)、[① 原文](docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/01-design.md#L42-L48)。
- S3 PASS — 1247/1177、5799/3812、0.6–0.9 s、≈0.2 s、10–15%/≈5%、<1%、≤250 ms、≤0.3–4%、≈0.9 s、≈94 s 均与 ① 一致；P88 冷基线 20865/12525/10311 与 P88 baseline TSV逐项一致。[报告表](docs/validation/gates-criteria-async.md#L28-L40)、[P88 基线](docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/baseline.tsv#L26-L28)、[⑤ 核验](docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/05-pre-review-r1.md#L19)。
- S4 PASS — 双语均为 UTF-8 BOM、LF 无 CR；LF 数 133、分割行数 134（含末尾空行），行数相等；`artifacts`围栏逐字相等，末行均为闭集令 `SLICE: READY_FOR_REVIEW`。[中文末行](docs/validation/gates-criteria-async.md#L133)、[围栏起点](docs/validation/gates-criteria-async.md#L115-L130)。证据包为 9 载荷，MANIFEST 与 SHA 列表条目及实际载荷集一致；重算 9 个 SHA 均匹配。[MANIFEST](docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/MANIFEST.md#L8-L19)、[SHA 列表](docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/SHA256SUMS.txt#L1-L9)。
- S5 PASS — 中文/英文净自耗句均已夹钳低于探针者按 0 计；主控头部及台账已标用户裁定“不/不得作先例”。[中文口径](docs/validation/gates-criteria-async.md#L36)、[中文主控](docs/validation/gates-criteria-async.md#L3)、[台账口径](docs/t027/REMAINING_SLICES.md#L2553)、[⑤ 两条 Low 原文](docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/05-pre-review-r1.md#L11-L12)。
- S6 PASS — 未发现禁用占位形态；双语豁免行理由非空，且当前证据显示⑥尚未执行。[中文豁免](docs/validation/gates-criteria-async.md#L131)、[英文豁免](docs/validation/gates-criteria-async_EN.md#L131)、[契约核验](docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/05-pre-review-r1.md#L22)。

### Section D
以下观察会推翻 PASS：暂存集出现 `tools/**`、`internal/**`、CONTRACTS、schema、fixture、CI/workflow；报告翻转“不做”结论或准则 v1.77；任一实测数字漂移；报告对 LF 行数、围栏或末行闭集不一致；MANIFEST/SHA/载荷集或任何载荷 SHA 不一致；低耗钳位或主控不得作先例标记缺失；豁免理由为空、与实际状态矛盾，或出现禁用占位形态。

### Section E
我无法独立复核⑦终审、批量辅助两次实派及⑥记录落盘后的最终内容；逐判据绝对值实测原始捕获件不在本次审查输入中，故其真值主要由①算术分解及⑤原文交叉核验。

INDEPENDENT SCAN: PASS
```

(end of transcript; the verdict line is preserved verbatim at the fence end.)

### 5.3 (7) final review and re-review

(7) r1 = `FINAL REVIEW: FINDINGS` (1 Medium + 1 Low). Medium = the report phase state was not synchronised with the completed rings; Low = the "verbatim transcript" wording does not match the verdict-line placement. Both are fixed (§1 phase sentence, §5.1/§6 back-fill, the transcript note reworded to "structured transcript"). Per §5.3 the re-review r2 runs with a narrowed input (fix diff + previous findings + mechanical readings).

(7) r2 = `RE-REVIEW: PASS` (both the Medium and the Low are CLOSED, no new Medium+); the 1+1 quota is fully used and convergence is reached. The r1 / r2 verbatim outputs are archived at `docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/07-final-review-r1.md` and `07-final-review-r2.md`.

## 6. Metering

| Role | Model / carrier | Calls |
|---|---|---|
| (1) architect | `DeepSeek V4 Pro (deepseek)` Max | 1 |
| (2) documentation engineer | `DeepSeek V4.1 Flash (deepseek)` | 1 |
| (3) consistency check | `DeepSeek V4.1 Flash (deepseek)` | 1 |
| (4) gates | lead, mechanical | several rounds (index/tree) |
| (5) pre-review | `DeepSeek V4 Pro (deepseek)` Max | 1 |
| (6) independent scan | `MAI Code 2 Flash (copilot)` | 1 |
| (7) final review + re-review | `GPT-5.3-Codex (copilot)` (user ruling) | 1 + 1 (fully used) |
| batch helper | `GLM-5.3-Flash (glm)` (carrier `deep-reasoner`) | 2 |

- R = 1, D = 0 (final values); the R calibre is in §7; probes 0; no blind review draw.
- Single source of truth for metering = the session call receipts.

## 7. Anomaly handling

- **Step 0 first-round index gate G8-a flagged red**: the evidence pack manifest was still missing `MANIFEST.md` and `SHA256SUMS.txt`, a genuine gap; both were added, the rerun cleared to zero. This is step-0 self-healing, not a deviation, and is not counted into R.
- **The (4) gate red sequence (each fixed in place and rerun to zero; per the P88 precedent not counted into R, each registered)**: G3(b) status glyphs on the CN side only -> reworded to drop the glyphs; G7-e last line `SLICE: IN_PROGRESS` outside the closed set -> `SLICE: READY_FOR_REVIEW`; G7-f new report without a fence -> a non-empty-reason skip declaration first, replaced by the verbatim fence after (6); G7-f conclusion line must sit at the fence end -> appended there; G4a double BOM (the normalising script did not strip the old BOM first) -> strip then prepend one; G2 placeholder shapes in two places -> reworded; G4b unknown char -> reworded away (`cjk-newwords.txt` untouched).
- **(5) two Lows (inside a PASS verdict)**: fixed (§3.1 calibre sentence, lead-ruling markers), no Medium+, no re-review; per the P88 precedent not counted into R.
- **(7) r1 one Medium + one Low**: fixed plus the r2 re-review converged => **R = 1** (the final R of this slice).
- Temp-file D count: 0 - the only temp file is `tmp/P95-GATES-CRITERIA-ASYNC/03-consistency-check.md` (the (3) output, inside the declared location), removed at (8c).

## 8. Major-change recommendations (handed to later slices)

1. **P88 report section 7 recommendation 4 (criterion asynchronisation)**: **assessed, benefit does not hold, not implemented** (closed by this slice); no longer a candidate.
2. **P88 report section 7 recommendation 2 (cross-scope base reuse)**: still a genuine candidate; it sits in the same sensitive family as "never cache verdicts" and needs a clause before any implementation.
3. **OB-113 sub-item 3 (G4b cold-scan inner loop)**: still a genuine candidate; it must be handled together with the cache-version bump and key migration, with the benefit argued first.
4. **worker_threads parallel path**: the cost evaluation does not hold (git child-process migration, cross-thread cache consistency, change-set serialisation; ceiling only ≈0.9 s), so it is not filed as its own slice.

## 9. Artifact fence

```artifacts
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/00-briefing.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/01-design.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/04-reviews/07-final-review-r2.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/05-batch-assist/glm-mirror-check.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/05-batch-assist/glm-numbers-check.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/MANIFEST.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/SHA256SUMS.txt	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/freedom-list.md	repo
docs/validation/evidence/P95-GATES-CRITERIA-ASYNC/input-list.md	repo
docs/validation/gates-criteria-async.md	repo
docs/validation/gates-criteria-async_EN.md	repo
```

## 10. Submission and CI

- Main commit `238e45c` (`docs: record the P95 async-criteria assessment (not implemented)`), pushed to `origin/main` (range `22370fe..238e45c`; **gitee untouched**).
- CI run **`45939437180`** (head `238e45c`) ended `completed / success`; both legs succeeded (`Go ubuntu-latest` / `Go windows-latest`), the other candidate jobs are skipped by the CI definition.
- Hard evidence from both legs: `SELFTEST: PASS 770/770`, `GATE REPORT scope=ci changed=16`, `TOTAL_FAIL=0` (on both legs); package-level `ok` lines ubuntu 21 / windows 15, package-level `FAIL` 0; real failure lines verified = 0 (the self-test's expected-failure assertion descriptions are not counted).
- This section was added by a docs-only commit after both CI legs turned green (§11.2); the backfill commit's own CI run is observation-only.

SLICE: READY_FOR_REVIEW
