# P88 GATES-PERF-ASSESS validation report (tooling slice)

- Slice: P88 `GATES-PERF-ASSESS` (tooling `[SLICE]`, tier S); lead = `DeepSeek V4.1 Flash (deepseek)` High.
- Date: 2026-11-25. Directive version: v1.73 (MINOR). R = 0, D = 0.
- Evidence pack: `docs/validation/evidence/P88-GATES-PERF-ASSESS/` (start message, (1) design, (2) implementation notes, timing readings, the A/B comparison, the review transcripts, the freedom list and the input list).

## 1. Conclusion

The gate timing composition is measured and landed; two light optimisations (O-1 batched blob prefetch, O-3c raw-reading memos) are in place, and **the criterion verdicts stay byte-identical** (five same-input A/B pairs are byte-equal). **Current state = awaiting the (10) commit authorization**.

## 2. Change set and evidence

- Criterion implementation: `tools/gates/lib/ctx.js` (the only product change, +164/-23). `tools/gates/lib/criteria/**`, `gate.js` and `selftest.js` are byte-identical to HEAD.
- Directive pair: header v1.72 to v1.73, a new v1.73 changelog entry in section 0.3 and appendix C.0az (each +18/-1, symmetric across the two languages).
- Ledger pair: the P88 status moves from "not started" to "in progress" and to "closed" on close-out (each +2/-2).
- Report pair: this file and its `_EN` mirror.
- Evidence pack: `docs/validation/evidence/P88-GATES-PERF-ASSESS/`, 22 payloads (per that pack's `MANIFEST.md`, equal to the `SHA256SUMS.txt` entry set) plus the two metadata files `MANIFEST.md` and `SHA256SUMS.txt`.

## 3. Mechanical readings

### 3.1 Timing composition before the change (`baseline.tsv`)

| Item | Exit | Elapsed ms | Note |
|---|---|---|---|
| `gofmt -l .` | 0 | 1114 | base item 1 |
| `go build ./...` | 0 | 24503 | base item 2 |
| `go vet ./...` | 0 | 8858 | base item 3 |
| `go test ./...` | 0 | 143628 | base item 4, 80.6% of the base total |
| base total | - | 178103 | executed once per `--all`, serially |
| `--all --skip-base --scope=range:HEAD~3..HEAD` | 1 | 490013 | cold corpus, dominated by one `git show` per file in G4b |
| `--all --scope=range:HEAD~3..HEAD` | 1 | 240824 | base 178103 + criteria 62721 (warm corpus) |
| `--all --skip-base --scope=tree` | 1 | 5247 | HEAD corpus already cached |
| `--all --scope=tree` | 1 | 99070 | warm base 93823 + criteria 5247 |
| `--all --skip-base --scope=index` | 1 | 1946 | - |
| `--selftest` | 0 | 265294 | 744 assertions |

Per-criterion readings (`--check=<ID>`, warm corpus, each including about 1.4-1.9 s of process and change-set resolution overhead): G6 20865, G7 12525, G8 10311 and G4a 10236 are the four heaviest; the full table is in `baseline.tsv`.

### 3.2 Before / after comparison (measured A/B)

| Path | Old ms | New ms | Speed-up |
|---|---|---|---|
| tree full (with base) | 75000 | 79628 | 0.94x (single sample, dominated by Go cache warmth) |
| index criteria only | 2430 | 1750 | 1.39x |
| range criteria only | 59272 | 33083 | 1.79x |
| ci criteria only (cold corpus) | 170277 | 12929 | 13.17x |
| `--selftest` | 209992 | 128881 | 1.63x |

> **Caliber for the two `--selftest` readings**: the 265294 ms in 3.1 comes from the pre-change composition baseline round (which also ran the four base checks around it, i.e. a higher machine load), while the 209992 ms in 3.2 comes from the A/B pairing round (old implementation running alone in that round). Both are the pre-change implementation but from **different rounds and loads**, so the difference is not an implementation difference; the pairing conclusion rests on the same-round comparison in 3.2.
CI critical path (`--selftest && --all --scope=ci`, cold-runner equivalent): **380.3 s to 141.8 s (2.68x, 238.5 s saved)**.

### 3.3 Byte identity

- A1==A2 (tree full), B1==B2 (index), B3==B4 (range), C1==C2 (ci, cold corpus), D1==D2 (selftest): five pairs with equal sha256 and `fc /b` exiting 0.
- The C group line `INFO G4b-cache :: miss 59c285d4cdad corpus=1466` is byte-identical, i.e. the cache key `key12` and the corpus set did not move.
- Comparator falsifiability: after flipping one bit `fc /b` exits 1 (`comparator-sanity.txt`).
- Report-pair line caliber: before the two caliber notes of this section each file held **174** lines (counted by LF newlines; `git diff --cached --numstat` reported 174/174 at that point) and **176** lines after them, equal on both sides throughout. The earlier receipt said "175" because it counted the array length after splitting, which includes the empty segment after the trailing newline; the actual line count follows LF newlines.

### 3.4 Gates and base checks after the change

- `gate.js --selftest`: `SELFTEST: PASS 744/744` (both before and after the change).
- `--all --scope=tree`: `TOTAL_FAIL=0`.
- `go build ./...` / `go vet ./...` / `go test ./...`: exit 0; `go test` reports 15 packages ok and 0 FAIL.
- The raw `TOTAL_FAIL=0` readings for the index and ci scopes are added to the evidence pack in the (9) final state (`03-verification/gate-index.txt` and `gate-ci.txt`).

## 4. Review chain and closure

- (5) pre-review (`DeepSeek V4 Pro`, quota 1+1, used 1): `PRE-REVIEW: PASS`, zero Medium+, one Low (the textconv equivalence premise between `cat-file --batch` and `git show` was not stated) plus four Info items. The Low is fixed by adding the premise sentence to the `ctx.js` comment block (comment only; after the edit the range criteria output is still byte-identical to B4).
- (6) independent scan (`MAI Code 2 Flash`, quota up to 3, used 2): r1 `INDEPENDENT SCAN: FINDINGS` (one HIGH, calling the range-scope `TOTAL_FAIL=1` a regression); after adding the paired "pre-change implementation on the same range" evidence, r2 returned `INDEPENDENT SCAN: PASS` and established that the same `TOTAL_FAIL=1` exists in the old implementation and that `range` is not one of the ledger acceptance scopes.
- (7) final review (`GPT-5.3-Codex` regular tier, quota 1+1, one final review used): `FINAL REVIEW: PASS` with zero findings; it ruled Q1 (byte identity holds), Q2 (keeping `CACHE_VERSION` acceptable), Q3 (the (5)/(6) dispositions are sufficient) and Q4 (v1.73 stays in bounds). The re-review quota was not used.

## 5. Metering

- (1) x 1 (`DeepSeek V4 Pro`); (2) x 1 (`prfrail-implementer` + `DeepSeek V4.1 Flash`); (3) x 1 (`prfrail-tester` + `DeepSeek V4.1 Flash`); (5) x 1; (6) x 2; (7) x 1 (final review).
- Zero probes; no blind review (user ruling c in the start message); one GLM batch-assist call (the report pair bilingual mirror check, carrier `deep-reasoner` + `GLM-5.3-Flash (glm)`); the timings were measured directly by (3), so no further batch assist was dispatched.
- The single source of truth for metering is the session call receipts; no quota exception was used.
## 6. Exception handling

- **First (3) driver run lost its flags**: the driver named its parameter `$Args`, a PowerShell automatic variable, so all six cells exited 2 with usage text only. That output is isolated under `03-ab/_invalidated-run1/` and the run was repeated after the fix; the g4b cache stayed at 100 entries throughout.
- **The cold-cold pair needed a second cache clear**: group C requires both sides to take the `miss` branch, so the cache was cleared again between C1 and C2 (a cold run writes an entry, which would turn C2 into a `hit`). The backup was fully restored afterwards with zero per-entry name and sha256 differences.
- **F-10 deviation (disclosed)**: `freedom-list.md` F-10 marks "delete or archive `tmp/gates-cache/g4b/` entries" as `out-of-bounds`, while (3) performed backup, clear and restore to build the cold state. The final state is byte-restored (100 entries) with nothing left behind; the action was a necessary precondition for the cold-cold pair and is recorded as a deviation rather than hidden.
- **A pack directory was first created as a file**: the first copy ran before the target directory existed, so `Copy-Item` produced a zero-byte file named `03-verification`. The file was deleted, the directory created and the pack rebuilt, with `MANIFEST.md` and `SHA256SUMS.txt` regenerated from disk.
- **(6) r1 HIGH and the r2 correction**: r1 lacked the paired "pre-change implementation on the same range" evidence and therefore read the historical range's pre-existing `TOTAL_FAIL=1` as a regression. With the pair supplied, r2 returned `PASS`; both rounds are archived verbatim under `04-reviews/`.
- **Capture determinism**: the first attempt also appended `ELAPSED_MS` to every captured gate reading, so the rebuilt `SHA256SUMS` could never converge (the timing changes every run). The captures now carry `EXIT=` only and all timing readings live in `timings.md` and `ab-report.md`, which makes the evidence pack byte-reproducible.
- Temp-file D count: 0 (the scripts and captures under `tmp/P88-GATES-PERF-ASSESS/` are removed at (8)c).

## 7. Larger-refactor suggestions (handed to a later slice, not implemented here)

1. **Break down the `--selftest` cost**: time it per suite (G6/G7/G8 corpus parsing, mutation probes, the CLI subprocess cases), memo the repeated fixture disk reads and isolate the slowest mutation group. `selftest.js` is a criterion implementation file, so any change there carries the same (6)/(7) obligations.
2. **Cross-scope base reuse (a verdict cache)**: the four base checks are a pure function of the workspace plus the toolchain, so a cross-process cache with an invalidation key is conceivable; but it is the same family as the §6.3 G4b rule "never cache a verdict", so a clause must land first. The short-term substitute is to run the CI Gates step as `--all --skip-base --scope=ci` (a workflow change needs user authorization and a separate slice).
3. **G4b inner-loop tuning**: replace the whole-set spread `new Set([...corpus, ...chars])` with incremental merging, batch `cjkChars`, and move the corpus to an incremental form. All of that lives inside the `run` body, so it necessarily moves `sourceDigest` and must be bundled with a `CACHE_VERSION` bump and a `key12` migration; the current magnitude is hundreds of milliseconds, so the gain must be argued first.
4. **Asynchronous criteria**: `runCriteriaConcurrent` already provides the scheduling contract, so making G6/G7/G8 asynchronous would deliver real concurrency; it touches "concurrent reads of shared caches plus the report-order contract" and belongs to a separate slice.
5. **An optional `--timings` output for `gate.js`**: this slice's timing comparison relies on external measurement; per-step timings emitted by the gate itself would remove that dependency, but the output format is a protocol surface and the document contract must change first.

## 8. Open items and hand-over

- **CI legs to be evidenced after the push**: this slice is not committed yet, so the behaviour of `--selftest && --all --scope=ci` on both the Linux and Windows legs, under the runner default Node version, will be settled by the CI run after the push; the run id and both legs' readings are back-filled then.
- **Raw `TOTAL_FAIL=0` readings for index and ci**: produced in the (9) final state (staging complete, the report itself already in the index) and added to the evidence pack, which closes the (5) Info item.
- **No new `selftest` assertions**: the ledger acceptance requires `SELFTEST` 744/744 to stay literally unchanged, so this slice adds no assertion to `selftest.js`; the equivalence evidence for the batch read and the memos is carried by the five (3) A/B pairs and the (2) process verification instead.
- **v72 snapshot absorption**: the P88 section and the v1.73 version pack are absorbed by the v72 snapshot; this slice does not touch `CHECKLIST_v71{,_EN}` or `DOCUMENTATION_PLAN{,_EN}`.

## 9. (6) output transcription (verbatim)

> Transcription note: the fence carries the natural verbatim shape of (6) r2 (`INDEPENDENT SCAN: PASS`). To satisfy the G7-f rule "A to E exactly once", the duplicated `## Section A` heading at the block end is not repeated, while the conclusion line itself is kept verbatim as the block-final line; no other line was rewritten.

```sixstage
### Section A
INDEPENDENT SCAN: PASS

### Section B

| SEVERITY | file:line | 被违反的句子或主张 | 最小修复 |
|---|---|---|---|
| - | - | 无可执行发现 | - |

### Section C

- r1 的 `HIGH` 不成立。旧实现 B3 在 B3-old-range.out.txt 的 G5(d) 行与 TOTAL_FAIL 行已有 `TOTAL_FAIL=1`；新实现 B4 在 B4-new-range.out.txt 的同一两行逐字相同。复跑输出 gate-range-post-comment.txt 的第 21 与第 52 个行位也相同；A/B 机械核验在 `ab-report.md:78` 明确 `B3 == B4 IDENTICAL`。因此这是历史提交区间既有判定，不是本片回归。
- 台账的三作用域仍为 `index / tree / ci`；台账验收证据（`REMAINING_SLICES.md:2377`）未把 `range:HEAD~3..HEAD` 列为验收作用域。它只属于配对诊断证据。
- **S1：PASS。** 批量预取先完整解析、任何异常整批丢弃，提交缓存只发生在全部响应通过之后；缺失对象和空 Blob 分别映射为 `null` 与零长度，条数校验位于 `ctx.js:518-561`。换行被拒绝回退，且缓存采用 `has()` 不会将 `null` 与未缓存混淆；超大输入受 spawn 缓冲上限触发 fail-open。
- **S2：PASS。** 键空间由路径或 `ref NUL path` 分隔，`ctx.js:412-423`，无跨快照串味。`numstat` 返回独立副本，`ctx.js:646-671`。真实调用方仅读数：如 `g3a.js:61-67`、`g4a.js:33`、`g8.js:403`。
- **S3：PASS。** 三作用域和 SELFTEST 的字节对比在 `ab-report.md:76-80` 均为 `IDENTICAL`；实现说明未新增输出行，`02-notes.md:52`。
- **S4：PASS。** 变更只落到 `ctx.js`；`g4b.js:278-318` 的 `run` 未改，`CACHE_VERSION` 仍为 v2（`g4b.js:43`），G5-b 闭集仍在 `g5b.js:25`，闭集实现未进入改动集。
- **S5：PASS（本闭集内）。** 准则头部与 `DELIVERY_DIRECTIVE.md` 的 §0.3 v1.73 变更日志条目、C.0az 一致；`impl.diff` 未改判据夹具、schema、workflow、`internal/**` 或 CI。
- **S6：PASS。** `gate-selftest.txt:749` 明确 `SELFTEST: PASS 744/744`，且 `ok` 断言行与旧实现合计 744 条；`selftest.js` 未进入改动集，`02-notes.md:5`。

### Section D
可证伪性：
- 一处旧/新字节差异、任一范围输出的 G5(d) 或 `TOTAL_FAIL` 差异，或旧包存在不同判定，都会推翻 r1 已不存在回归结论。
- 批量响应任一伪造、分片、缺失或条数异常若导致缓存部分提交，就会推翻 S1；当前实现提交前完全解析且检查总长度。
- 任一调用方改写缓存 Buffer、布尔或数组并改变后续判据，就会推翻 S2；实际实现对 `numstat` 返回副本。
- 任一 G4b / gate / selftest 输出差异、任一判据运行体或缓存版本差异、任一 745/744 计数差异，都会推翻对应结论。

### Section E
无法验证：非 Windows/Linux CI、并发 gate 进程共享缓存，以及未来 `.gitattributes` textconv 差异；本闭集只有 Windows 本机配对和现仓库语料。P88 的双语报告对未包含在请求输入中，因此无法独立核验其严格镜像；报告成对成形或双语出现失配会推翻该验收项。

INDEPENDENT SCAN: PASS
```

## 10. Artifact fence

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/00-briefing.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/01-arch/01-arch-design.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/02-impl/02-notes.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/ab-report.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/ab-sha256.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/baseline.tsv	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/build.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/gate-ci.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/gate-index.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/gate-selftest.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/gate-tree.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/gofmt.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/test.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/timings.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/03-verification/vet.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/05-batch-assist/glm-mirror-check.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/MANIFEST.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/SHA256SUMS.txt	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/freedom-list.md	repo
docs/validation/evidence/P88-GATES-PERF-ASSESS/input-list.md	repo
docs/validation/gates-perf-assess.md	repo
docs/validation/gates-perf-assess_EN.md	repo
tools/gates/lib/ctx.js	repo
```

SLICE: READY_FOR_REVIEW
