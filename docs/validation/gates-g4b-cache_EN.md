# G4b-CACHE slice validation report (business slice: cross-process corpus cache for G4b)

> Slice: `P27 / G4b-CACHE` (directive v1.49, not a governance slice, M tier). Date: 2026-10-25. Baseline `HEAD` = `7b84528` (after the P26 delivery). Change landing points in §2, review chain in §3, gate readings in §4.

## 1 Slice overview

| Field | Content |
|---|---|
| Slice id | `P27 / G4b-CACHE` |
| Date | 2026-10-25 |
| Size and tier | M (business slice, performance, not governance) |
| Baseline and directive | `HEAD` = `7b84528` (after the P26 delivery, directive v1.48), v1.49 once this slice lands |
| Goal | add a **cross-process cache** to the G4b baseline-corpus CJK character set: the whole `--scope=tree` run drops from a cold median of **68539 ms** to a warm median of **7226 ms** |
| Boundary | caches the corpus character set only; never the verdict, the added-line characters or the whitelist; changes no criterion judgement logic |

**Five design principles** (ruled by ①, implemented by ②, independently verified by ③ and ⑥):

1. store the corpus character set only: an entry value is always `{v:1,corpus:string[]}` and carries no judgement information.
2. never cache the verdict: `out.rec('G4b', …)` runs unconditionally every time, so a hit must still recompute it.
3. six key segments: version, resolved sha, `ctx.kind`, whitelist source, whitelist digest, criterion source digest.
4. recompute the key before every read: the key is never persisted alone and is always rebuilt before the hit test.
5. fail open on any anomaly: a failure to create the directory, read, parse, write or `rev-parse` silently degrades to the full scan.

## 2 Change landing points

| File | Size | Role |
|---|---|---|
| `tools/gates/lib/criteria/g4b.js` | **280 physical lines**, numstat `241/34` | cache implementation: key, store, fail-open and comments |
| `tools/gates/selftest.js` | adds T108–T113 for **23** assertions in total (674 → **697**), numstat `354/0` (zero deleted lines) | cache key, fail-open, dirty index and load-bearing mutations |
| `docs/DELIVERY_DIRECTIVE.md` | **2034 physical lines**, numstat `18/2` | header v1.49, the §0.3 row, the §6.3 G4 cache clause plus the **version-bump duty**, **OB-120** in §12.4, appendix **C.0aj** |
| `docs/DELIVERY_DIRECTIVE_EN.md` | **2034 physical lines**, numstat `18/2` (symmetric with CN) | the mirror side, aligned sentence by sentence |
| `docs/t027/REMAINING_SLICES.md` | numstat `21/0` | the P27 row and detail block, including the observation / next-slice advice row and the known-limitation row |
| `docs/t027/REMAINING_SLICES_EN.md` | numstat `21/0` (symmetric with CN) | the mirror side |
| `tools/gates/cjk-newwords.txt` | **6** new registrations (**37** → **43**) | whitelist entries for the characters the verbatim transcription and the withdrawn wording require (each carrying a reason and a context excerpt) |

Additional notes (landing form and ordering self-checks):

- **(1)(a) Landing form**: the cache lands across processes at `repoRoot/tmp/gates-cache/g4b/`; it is NOT placed under `tmp/G4b-CACHE/` because that directory is a slice temp artifact and ⑧c clears it - clearing would destroy the cache, contradicting the C.0aj requirement that "the cache directory is runtime state" stays self-consistent.
- **③ ordering self-check**: ③ = `L68115`, ⑤ = `L69717` => ③ < ⑤ holds; ④ is a master-side step (not in the dispatch table) and runs after ③ per §5.3.
- **⑥ before ⑦**: ⑥ r3 = `L71059`, ⑦ r1 = `L71123` => ⑥ < ⑦ holds.
- **R7.1 itemised metering**: see the 26 rows of the §3.1 dispatch table; `retryOf` semantics = ②-2 through ②-7 each re-run the previous ② after remediation; ⑥ r2/r3 are the format re-issue and the remediation re-run of r1; ⑦ r2 is the remediation re-review of r1.

## 3 Review chain

### 3.1 Dispatch table (mechanical read)

| Stage | Role (`agentName`) | Carrier model | Quota | Conclusion |
|---|---|---|---|---|
| ① | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 3 rounds | first version, incremental re-check and the ⑦-M ruling, all `DESIGN: DONE` |
| ② | `prfrail-implementer` | `DeepSeek V4.1 Flash (deepseek)` | 6 rounds | implementation, condition ③ removal, comment precision, four write-backs, OB-120 in §12.4, the restored wording then withdrawn |
| ③ | `prfrail-tester` | `DeepSeek V4.1 Flash (deepseek)` | 3 rounds | all three `TEST: DONE` (Medium+ zero) |
| ⑤ | `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | 1+1 exhausted | r1 `PRE-REVIEW: PASS WITH FIXES`, r2 the same |
| ⑥ | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | 3 exhausted | r1 void for contract deviation, r2 `FINDINGS`, r3 `PASS` |
| ⑦ | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | 1+2 | r1 `FINDINGS` (1 High), r2 `RE-REVIEW: PASS` |
| batch help | `deep-reasoner` | `GLM-5.3-Flash (glm)` | 3 calls | enumeration, key facts and a test draft |

- the three batch-help calls do read-only preparation only: they produce no verdict and change no quota count; ⑥ r1 was void because its output contract did not conform (see §3.5), carrying no review value while still consuming quota.

### 3.2 ① four key rulings

| Item | Conclusion |
|---|---|
| Cache value purity | stores `{v:1,corpus:string[]}` only; the verdict is recomputed every time (`out.rec('G4b', …)` in `g4b.js` runs unconditionally) |
| Key completeness | the three counterexamples are carried by the sha segment, the whitelist-digest segment and the `sourceDigest()` segment respectively; `ctx.kind` is a **load-bearing segment** (range uses `ls-tree`, tree and index use `ls-files`) |
| fail-open | `mkdirSync`, `get`, `parse`, `put` and `rev-parse` all degrade silently; the only `UsageError` is `whitelist.error`, raised before any cache work |
| **Condition ③ removal** | the first design required "the cache is enabled only when the tree and index are clean"; ① ruled that argument unsound in the incremental re-check (the corpus reads the HEAD object store only); after removal `ctx.kind` became load-bearing, and **T113 plus a mutation** pin it against resurrection |

### 3.3 ③ three independent rounds

| Round | Conclusion | Highlights |
|---|---|---|
| r1 | `TEST: DONE` (Medium+ = 0) | assertion count 674 → 691; the three counterexamples and fail-open reproduced by independent probes; every key segment proven load-bearing by mutation |
| r2 | `TEST: DONE` (Medium+ = 0) | verification of the write-backs that took 691 → 697; `sha256` restored bit for bit after repeated mutations of the production file; the three tree timing groups and the tree↔ci attribution |
| r3 | `TEST: DONE` (Medium+ = 0) | adds that `sourceDigest` is unchanged at `10d7ec758d494d7570d7e98e227802c5732ef6db0a59b6f604f02dbf46b5e21f` (the split sha256 of `run.toString()` and `cjkChars.toString()` is byte-identical); the six-step real-git end-to-end run; a median re-check of the tree group |

What the three rounds share: assertion count 674 → 691 → **697** (a net **23** assertions); the three counterexamples (a HEAD change, a whitelist change, a criterion-source change) and fail-open are reproduced by **independent probes**; every key segment is proven **load-bearing** by mutation; the production file `sha256` is restored bit for bit; and the six-step real-git "clean write → dirty read still hits" run uses the temporary repository `repo-clean` (`git write-tree` identical before and after at `c9bfc762…`).

### 3.4 ⑤ pre-review (quota 1+1, both rounds used)

- **r1 `PRE-REVIEW: PASS WITH FIXES`**: 1 Medium (the `G4B-CORPUS-DECOUPLE` observation and next-slice advice registered in **no** staged document) + 3 Low (T108–T112 omit T113; §6.3 lacks the version-bump duty; a `'none'`-key TOCTOU in the second `rev-parse` inside `keyString`) + 3 Info.
- **r2 `PRE-REVIEW: PASS WITH FIXES`**: the 1M and 3L of r1 are all closed; 1 Low added (the P27 status row still reads "① pending dispatch") + 3 Info (the §6.3 version-bump duty has no mechanical enforcement; #4 and #5 await the report; editor buffer hygiene); ⑤ ruled explicitly that the master handling - "restore 彻底 and register it in the whitelist" - was **correct**.

### 3.5 ⑥ independent scan (three rounds, with verbatim transcription)

- r1 concluded `PASS` but **violated the output contract** (a `## Section A — 结论` heading shape, a Section D that is not a falsifiability audit, a final line that is not a conclusion line) ⇒ void.
- r2 `FINDINGS`: 1 Medium (the cache key was not mechanically bound to the semantics of `ctx.listAllMarkdown` and `ctx.showAt`).
- r3 `PASS`: after the ① ruling and ⑥ own re-check that Medium is **withdrawn** and downgraded to a known limitation.

The following is the complete r3 output, transcribed **verbatim** (nothing rewritten); the two conclusion lines inside the fence (the first line under Section A and the block-final line) carry equal values at legal positions, which matches the `G7-f` rule in `§6.3`.

```sixstage
## Section A
INDEPENDENT SCAN: PASS

## Section B
无

## Section C
- 命令 1: cd /d/LZProjects/prfrail; git rev-parse --short HEAD; git status --porcelain -uall
  关键输出摘要: HEAD=7b84528；状态为 6 个 M、0 个 ??，没有未跟踪残留。
- 命令 2: git log --oneline -- tools/gates/lib/ctx.js
  关键输出摘要: 仅 2 条提交；即 ctx.js 的语义扩展历史极短。
- 命令 3: git show --stat --oneline 8a2841f -- tools/gates/lib/ctx.js tools/gates/lib/criteria/g4b.js tools/gates/selftest.js
  关键输出摘要: 8a2841f 同时改了 3 个文件，包含 ctx.js、g4b.js、selftest.js。
- 命令 4: git show --stat --oneline 025158f -- tools/gates/lib/ctx.js tools/gates/lib/criteria/g4b.js tools/gates/selftest.js
  关键输出摘要: 025158f 同时改了 ctx.js 和 selftest.js；两次历史都未显示“只改 ctx.js 且不联动 g4b.js/selftest.js”。
- 命令 5: git grep -nE 'actions/cache|gates-cache|tmp/' -- .github/workflows/ci.yml
  关键输出摘要: 无命中；CI 工作流未配置 actions/cache，也未使用 gates-cache 或 tmp/ artifact 复用。
- 命令 6: node tools/gates/gate.js --selftest
  关键输出摘要: SELFTEST: PASS 697/697；TOTAL_FAIL=0；G4b unknown=[]；whitelist=37；退出码 0。
- 命令 7: 读取 .github/workflows/ci.yml 的 Gates 步骤和 checkout 片段
  关键输出摘要: `Checkout source` 配置 fetch-depth: 0；Gates 步骤位于全新 checkout 之后；与用户给出的主控事实一致，CI 是冷跑，不存在本地 tmp/gates-cache 复用路径。

## Section D
(i) 失效方向：部分成立。它的“危险方向”确实是“旧缓存命中导致少报（false negative）”，但它并不能越过 CI；这是 ① 的核心判断。这里的关键不是“可能有少报”，而是“少报是否能通过冻结提交门禁”。在真实工作流里，CI 每次都是全新 checkout，且 `.github/workflows/ci.yml` 中没有 actions/cache、没有 gates-cache、没有 tmp/ artifact 下载；因此本地旧条目无法成为 CI 证据源。结论：该论点只证明“本地预检可能瞬时不一致”，不能证明“提交门禁会被绕过”。
(ii) 概率：成立。三条件连乘确实极低：ctx 语义变更 × 未升 CACHE_VERSION × 本地存在旧条目；而且 git 历史非常短：ctx.js 只有 2 条提交，且两个关键变更都同步改了 g4b.js/selftest.js。`keyString` 也把 resolved sha、whitelist 摘要、sourceDigest、CACHE_VERSION 一起绑定，且 `selftest.js` 的 T109–T111 已明确验证 HEAD 改动、白名单改动、判据源码改动都必须 miss。于是“语义漂移而未更新缓存键”的可发生性被强力压低。
(iii) 可检测性：部分成立。这里确实没有一个“在语义漂移时自动变红”的断言，T108–T113 主要覆盖的是缓存键、fail-open、索引脏、不变更/变异手段，而不是“真实 ctx.js 改了但没升级缓存键”的现场监测。这个事实是真，但它不等于“能绕过 CI”，因为结构性冷跑已经把风险边界压在本地、而不是提交门禁。也就是说，检测性缺口存在，但它只影响本地快速反馈，不改变 CI 的有效性。
(iv) 现行义务：成立。`docs/DELIVERY_DIRECTIVE.md` 中已经写明“升版义务”：ctx.listAllMarkdown / ctx.showAt / git 枚举或读取语义变化时，必须同步升级缓存格式版本前缀，否则旧条目一律视为陈旧。并且 `OB-120` 作为候选修法已登记，说明本问题已被纳入治理项，非无视。这里真正的问题不是“没有规则”，而是“该风险在当前工作流中不构成 Medium 级别的 CI 间接威胁”。
结论：本题不成立为 Medium；至少应归为 Low/已知限制，并且更稳妥的裁决是撤回 Medium。没有发现一个“能同时骗过本地命中与 CI 冷跑”的可发生反例序列；这类风险在当前仓库结构下只剩本地预检瞬时不一致，不足以支撑 Medium 级别结论。

## Section E
- r2 的格式问题：已修复。本轮严格按照模板输出，且最终行与 Section A 完全一致。
- 无法核验项：没有对“未来某次 ctx.js 语义改动 + 未升 CACHE_VERSION + 本地旧缓存被复用”的特定历史分支做真实演练；若存在该类实测记录，结论会被推翻。
- 可能推翻 PASS 的观测：若未来实测显示 CI 中确实存在 persistent cache 复用，或出现“ctx.js 改动未同步升级 g4b.js / selftest.js 但能通过新 checkout 结论”，则 Medium 可能重新成立。
INDEPENDENT SCAN: PASS
```

### 3.6 ⑦ final review (blind r1 plus re-review r2)

- r1 concluded with `FINAL-REVIEW: FINDINGS`: 1 **High** (`SELFTEST: FAIL 695/697`, two T87 assertions `unresolvable`).
- the master ran the **§9.1** three-step root-cause arbitration: `CAND_COUNT=0` before the sweep (the only node process was PID 4820, the MCP host); **7 single runs all `PASS 697/697`** (`EXIT=0`, `FAIL_LINES=0`); **both concurrent A and B `PASS 697/697`**; both branches of `resolveRepoArtifact` reported `ok:true`; `git ls-files --error-unmatch -- tmp/.gitkeep` exited 0.
- ruling: both hypotheses - "a defect in this slice" and "selftest cannot tolerate concurrency" - are **falsified**.
- two root-cause candidates stand side by side: **A (a one-off concurrent run, most likely the leftover selftest process that ⑥ r2 itself reported; its PID was never recorded and its lifetime cannot be reconstructed - `Kernel-Process/Analytic` holds no record on this machine)** and **B (a one-off `spawnGit` spawn failure)**; neither is 100% mechanically established.
- **Key qualification**: `spawnGit` collects stderr into the result object instead of printing it, so the absence of git error noise in `codex-st.txt` **cannot** distinguish "git reported an error" from "the spawn failed".
- r2 concluded `RE-REVIEW: PASS`: the High is not a code defect; the same-position assertions in the 7 single runs and in concurrent A and B all pass, and both branches of `res-probe` report `ok:true`; it also notes that this slice changes only the G4b cache and T108–T113 and never touches the T87 or G9 logic.

## 4 Gate readings

### 4.1 Readings at freeze (six files)

| Scope | EXIT | Last line | G1-b | G3(a) | G4a | G4b | G6-4/5/6/7 | G9 |
|---|---|---|---|---|---|---|---|---|
| `--all --scope=index` | 0 | `TOTAL_FAIL=0` | OK | OK | OK | OK | OK | OK |
| `--all --scope=tree` | 0 | `TOTAL_FAIL=0` | OK | OK | OK | OK | OK | OK |
| `--all --scope=ci` | 0 | `TOTAL_FAIL=0` | OK | OK | OK | OK | OK | OK |
| `--selftest` | 0 | `SELFTEST: PASS 697/697` | — | — | — | — | — | — |

Criteria readings: `G1-b` reports `header=v1.49 last=v1.49 rows=50 OK`; `G3(a)` is symmetric at `18/2` for the directive pair and `21/0` for the ledger pair; `G4a` reports `6 file(s) OK`; `G4b` reports `newCjk=401 unknown=[]` (whitelist=37, the state before this report landed; index and tree legs) and `newCjk=148 unknown=[]` for the ci leg (`endpoint:HEAD` - the ci leg covers `HEAD..base` only); `G9` reports `6 path(s) OK`. Also `SELFTEST: PASS 697/697` (7 single runs plus 2 concurrent runs all green).

### 4.2 Readings at freeze once this report lands

Once this report and the freedom list land, the change set is **10 paths** (7 `M` and 3 `A`, with 0 untracked `??`); the ⑨ frozen readings are:

| Scope | EXIT | Last line | changed | G4a | G4b | G9 |
|---|---|---|---|---|---|---|
| `index` | 0 | `TOTAL_FAIL=0` | 10 | `10 file(s) OK` (1 excluded: evidence tree) | `newCjk=658 unknown=[]` (`whitelist=43`) | `10 path(s) OK` (1 excluded) |
| `tree` | 0 | `TOTAL_FAIL=0` | 10 | same as above | same as above | same as above |
| `ci` | 0 | `TOTAL_FAIL=0` | 4 | `4 file(s) OK` | `newCjk=148 unknown=[]` (`whitelist=37`, `endpoint:HEAD`) | `4 path(s) OK` |

- `SELFTEST: PASS 697/697` (EXIT 0, `^FAIL` = 0).
- Note: this table holds **closeout-moment** readings; the number of cache entries under `tmp/gates-cache/g4b/` varies with runs (it was 9 at closeout) and does not affect any criterion reading.
- `G7-f`: `2 changed sixstage fence(s) complete (A→B→C→D→E, closed-set conclusion, SEVERITY consistency)`.
- `G9-e` registers INFO for the runtime logs this report references because they are **not staged** (informational, never red); this is by design - temporary artifacts are not committed, so referencing them only yields a hint and never a failure.

### 4.3 `SELFTEST` count change

`674` → `697` (a net **23** assertions): T108 3, T109 3, T110 3, T111 3, T112 5 and T113 6. Appendix reading: `git diff --numstat HEAD -- tools/gates/selftest.js` reports `354 0` (zero deleted lines).

## 5 Low-level finding dispositions

| # | Source | Severity | Location | Disposition |
|---|---|---|---|---|
| 1 | ① | — | `g4b.js:33` is **83 columns** (① preferred verbatim text) | accept the width overflow: the single-line comment it replaced was 87 columns (a net improvement), and the 81 columns at `g4b.js:137` match the file existing maximum, so the file **column ceiling is not raised** |
| 2 | ② | — | line-count caliber | Closed: ② once reported **281 physical lines**; the actual count is **280 physical lines** (an off-by-one in `split('\n').length`) |
| 3 | ③ | Low | a zero-byte untracked file named `({)` | Closed (temp-artifact discipline, self-healing branch): ③ found it, and after authorization `Remove-Item -LiteralPath '({)'` removed it, the `??` entries of `git status --porcelain -uall` went to zero, and `--scope=index` moved from `TOTAL_FAIL=1` to `TOTAL_FAIL=0`. Origin: mtime 05:22:34 (inside this slice window); a full transcript scan shows no recorded `> ({` redirect in that window and the neighbouring matches are all JS `=> ({` false hits, so the exact command cannot be located (registered as-is, never invented); the mechanism family is PowerShell treating an expression fragment as a redirect target. It is also a positive example for G5(f): the file made `--scope=index` report `TOTAL_FAIL=1`, proving that criterion really catches a stray untracked file |
| 4 | ⑤ r1 | Low | a `'none'`-key TOCTOU in the second `rev-parse` inside `keyString`, plus five T112 faces that do not directly cover a failing `mkdirSync`, a throwing `put` and invalid JSON syntax | Recorded: not fixed in this slice, already registered in the ledger known-limitation row; the direction is fail-safe and these paths are isomorphic to the tested ones |
| 5 | added by this slice | Low | the silent stderr of `spawnGit` in `selftest.js` | Observed: the candidate fix appends the stderr tail to `detail` on a spawn failure (no verdict change, diagnostics only); the trigger is the next unexplained failure of any git-related assertion, assessed by ①; not changed in this slice |

## 6 Boundaries and non-retroactivity

- `git diff HEAD -- lib/ctx.js tools/gates/gate.js .gitignore internal .github docs/DEV_PLAN.md` prints nothing.
- no change to §2.2, §2.4, §5.2, §7.5, any other criterion body, CONTRACTS, the schemas, the fixtures or the workflows.
- `tools/gates/cjk-newwords.txt` gains 6 registrations against `HEAD` (**37** → **43**, each carrying a reason and a context excerpt); encoding gate: `.md` is BOM plus LF, `.js` is no BOM plus LF.

## 7 Remaining items and follow-ups

### 7.1 Known limitation (adopted verbatim)

> ⑦ r1 independently reported a High (`SELFTEST FAIL 695/697`, two T87 `unresolvable`). The master ran the three-step root-cause arbitration: no gate or selftest process before the sweep; all 7 single runs `PASS 697/697`; both concurrent A and B `PASS 697/697`; both branches of `resolveRepoArtifact` and `git ls-files --error-unmatch -- tmp/.gitkeep` behave normally. ⇒ both hypotheses - "a code defect" and "selftest itself cannot tolerate concurrency" - are falsified. At least two root-cause candidates remain, neither 100% mechanically established: (A) a one-off concurrent run at the ⑦ execution moment, most likely the leftover selftest process that ⑥ r2 itself reported - that PID was never recorded and its lifetime cannot be reconstructed (`Kernel-Process/Analytic` holds no record); (B) a one-off `spawnGit` spawn failure in that run - `codex-st.txt` carries no git error noise, but `spawnGit` collects stderr into the result object instead of printing it, so "no noise" cannot distinguish "git reported an error" from "the spawn failed". The ⑦ reading is real; the attribution is corrected to a transient environment state (candidate A or B, undecided). Guard: no lock is added in this slice (concurrent A and B do not reproduce, and a lock would cover candidate A only; if the real cause is B a lock is useless); the observation is registered: if it recurs and the sweep confirms no leftover process ⇒ raise a separate slice that adds process-level mutual exclusion to the selftest entry (candidate A) or diagnostics plus a retry to `spawnGit` (candidate B). ⑦ also committed two process deviations (§7.3② executing a command, §7.3③ creating `tmp/G4b-CACHE/codex-st.txt`); both are registered and neither changes the validity of the readings.

### 7.2 Design deviations (four-element form)

- **D-1 (admitted by ①)**: activation condition ③ in the first design (for tree and index, enable the cache only when `git diff --cached --quiet HEAD` exits 0) rests on an incomplete argument; the correct invariant is "the cached corpus is always a subset of the true HEAD corpus and a hit is inherently fail-safe", independent of index cleanliness. The condition is removed (`ctx.kind` became a load-bearing key segment), and T113 plus a mutation pin it against resurrection.
- **D-2 (admitted by ①)**: the `cacheUsable` docblock delivered by ① once read "the full-scan path computes exactly the same index-dependent value"; that equality holds only while the index state has not changed since the write, and after a staged `.md` deletion the cold corpus shrinks and no longer equals the cached value (③ measured `HIT_VS_COLD_SAME_VERDICT=false`). It is rewritten as "the cached value is the HEAD corpus under the index state at write time, always a subset of the true corpus, later shrinkage only makes a cold run stricter, and the hit direction is fail-safe".
- **D-3 (execution-freedom expansion, withdrawn)**: during a write-back the master changed the OB-120 wording in §12.4 from 完全解耦 to 彻底解耦 and added a `彻` whitelist entry; a mechanical query shows both 完 and 全 already exist in the HEAD corpus (165 files contain 完), so 完全 was legal and the change was a **reverse operation**. the wording change and the whitelist entry are withdrawn and the whitelist is back to **37** entries (the 6 entries this report registers are separate, see §2 and §6). Lesson: an execution-freedom expansion (a wording change outside the reviewer findings) must pass a mechanical test first and never be optimised by intuition; "self-directed optimisation" is a risk source in a project with strict governance.
- **D-4 (temp-artifact discipline, self-healing branch)**: see item 3 of §5 (`({)`).
- **D-5 (one-sided claim annotation)**: ② once reported `EXIT=0` for `--selftest` as "derived from the code path" rather than captured in the same round; ③ then captured `EXIT=0` twice in its own rounds and that doubt went to zero.
- **D-6 (reviewer process deviation, ⑦ r1)**: ⑦ r1 executed `node … --selftest` (§7.3②) and created `tmp/G4b-CACHE/codex-st.txt` (§7.3③); the file timestamps 08:28:11–08:34:49 (about 6.5 minutes, matching one complete selftest) support this. **It does not change the validity of the readings**; a repeat in ⑦ r2 would be escalated per §7.3. r2 did not repeat it.
- **D-7 (same-turn temporary authorisation exception)**: the §5.1 pre-⑩ terminal read-back ② ③ does not pass for this slice (`tmp/` recursive file count > 1, with `gates-cache` as a top-level entry). Disposition = a same-turn authorisation exception granted by the user on 2026-10-25 (four elements in §7.5); this is **not a rule violation**. The runtime cache entries under `tmp/gates-cache/g4b/` are kept (the count varies with runs; it was 9 at closeout) and ⑧c does not clear them; this item must not become a precedent.

### 7.3 Observations and next-slice advice

- **Observation**: for tree and index `listAllMarkdown` uses `git ls-files` (reads the index), so a staged `.md` deletion shrinks a cold corpus relative to the true HEAD corpus and the cold run reports strictly more; inside a "write then staged deletion of `.md`" window a cache hit and a cold run can disagree (measured warm=PASS, cold=FAIL), in the fail-safe direction (the cached value is always a subset of the true corpus and never whitewashes a violation). **Candidate OB `G4B-CORPUS-DECOUPLE`**: the impact surface is verified - `listAllMarkdown` is **consumed by G4b only** (`g4b.js:223`); G4a uses `changedPaths` plus `readBytes`, G7 uses `changedPaths` and `exists`, G8 uses `listTreeFiles` (already on `ls-tree <headRef>`), and G9, G1 and G5 use `changedMarkdown`. Not executed in this slice; OB-120 is registered in §12.4 (first state word `待议`).
- **Candidate OB (performance)**: the source of the roughly **24–37 s** frozen cost of the ci scope is undetermined and the G4b cache cannot solve it (`--check=G4b` alone measures ci cold 59435 ms and warm 3090 ms, a 94.80% drop; the whole ci run is 91579 ms cold and 27428–39727 ms warm; `gate.js:147` is `await Promise.all(jobs)`, so the whole run is the concurrency critical path).
- **Candidate OB (diagnostics)**: the silent stderr of `spawnGit` (see item 5 of §5).
- **Candidate OB: the C.0aj ↔ §5.1 rule conflict** - C.0aj allows the cache directory to sit under `tmp/` and be kept, while the §5.1 pre-⑩ terminal read-back requires the `tmp/` recursive file count to be 1 with `.gitkeep` as the permanent top-level name; the two cannot hold at once. Candidate fixes: (i) move the cache directory out of `tmp/` altogether (e.g. `.gates-cache/` at the repository root) and add a pointer sentence in §5.1; (ii) add a "runtime cache directory exception" clause to §5.1 ② ③. Not fixed in this slice; handed to the next governance slice.

### 7.4 Three performance groups (median, `--scope=tree`)

| Group | 3 raw readings (ms) | Median | Log line |
|---|---|---|---|
| A miss (cache cleared every time) | 68539 / 66608 / 75815 | **68539** | `miss 88364fcfdfda corpus=1333` |
| B hit (no clearing) | 7232 / 7226 / 6808 | **7226** | `hit 88364fcfdfda corpus=1333` |
| C rebuild after invalidation (run after clearing) | 81089 / 72017 / 63941 | **72017** | `miss 88364fcfdfda corpus=1333` |

- the hit-to-miss ratio is **10.54%**, a **89.46%** drop (threshold ≥30%); group C is close to group A (a 5% gap, inside the cold run own variation); caliber: the `tmp/gates-cache` directory is cleared before a cold run, at least 3 runs give the median, a hit still recomputes the key and the verdict (non-zero cost), the data comes from this Windows machine and scales with the `.md` corpus size.
- the fixed-overhead reference (`--check=G1-a --scope=tree`) is 1846 ms; after deduction `(7226−1846)/(68539−1846)` = **8.07%**, a **91.9%** drop.

### 7.5 Same-turn temporary authorisation exception (§1.7.4 four elements)

| Element | Content |
|---|---|
| Decided by | the user |
| Date | 2026-10-25 |
| Reason | C.0aj (the cache directory is runtime state, kept as runtime state) and the §5.1 pre-⑩ terminal read-back (`tmp/` recursive file count = 1 with `.gitkeep` as the only top-level name) conflict internally; this slice neither fixes the rule nor changes code and takes the temporary authorisation exception instead |
| Basis | the user's reply in this turn, verbatim "选 D" |

- Nature: it exceeds the literal standard of the §5.1 pre-⑩ terminal read-back ② ③ and is a direct user ruling.
- Scope: only the runtime cache entries under `tmp/gates-cache/g4b/` are kept (the count varies with runs; it was 9 at closeout); §5.1 ② ③ do not apply to this slice.
- It must not become a precedent; the rule conflict is registered as a candidate OB (see §7.3) for a later governance slice.

## 8 Commit and CI

This report and the freedom list are the ⑧a artifacts; ⑨ re-runs the three scopes plus `--selftest` and appends the frozen readings to §4.2. The commit and CI evidence follow (a run triggered by a write-back commit after this one is an observation record only and gets no new commit of its own per §11.2).

| Item | Value |
|---|---|
| Commit | `64d1d58f82446bd9c1675df23523411785f5f29d` |
| Push | `origin` `7b84528..64d1d58` (never gitee) |
| CI run | `40183876643` (`https://github.com/larsonzh/prfrail/actions/runs/40183876643`) |
| Conclusion | `completed / success` |

- Legs: `Go windows-latest` = `success`, `Go ubuntu-latest` = `success`; the other 3 jobs are `skipped` (the candidate and release chains are not triggered, as expected).
- `Gates` step command (identical on both legs, repo-relative): `node gate.js --selftest && node gate.js --all --scope=ci`.
- Windows leg: `SELFTEST: PASS 697/697`, `TOTAL_FAIL=0`, `ok github.com/larsonzh/prfrail/internal/gates` (`0.244s`).
- Ubuntu leg: `SELFTEST: PASS 697/697`, `TOTAL_FAIL=0`, `ok github.com/larsonzh/prfrail/internal/gates` (`0.140s`).
- Independent corroboration: on a fresh checkout with no local cache, the 697 assertions (including T87 and T108-T113) are green on both Windows and Ubuntu => independent counter-evidence for the T87 concern raised in ⑦ r1.

```artifacts
tools/gates/lib/criteria/g4b.js	repo
tools/gates/selftest.js	repo
tools/gates/cjk-newwords.txt	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/gates-g4b-cache.md	repo
docs/validation/gates-g4b-cache_EN.md	repo
docs/validation/evidence/G4B-CACHE-freedom-list.md	repo
```
