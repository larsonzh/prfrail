# P92 GATES-SELFTEST-SPLIT validation report (tool slice)

- Slice: P92 `GATES-SELFTEST-SPLIT` (tool `[SLICE]`, tier M); controller = `DeepSeek V4.1 Flash (deepseek)` High.
- Date: 2026-11-26. Directive version: v1.75 (MINOR). R = 0, D = 0.
- Evidence pack: `docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/` (start message, the (1) design, partition data, A/B comparison, review transcripts, freedom list, input list).

## 1 Conclusion

The `--selftest` wall clock drops from 150.9 s to 70.1 s (2.15x) with byte-identical stdout and `SELFTEST: PASS 744/744`. The benefit gate rules the four original candidates item by item: only "repo-root discovery memo" clears the threshold and is implemented, while (b) fixture-read memoisation, (d) the G4b inner loop and (e) the in-tree `--timings` move out of this slice into the large-change list.

## 2 Change set and evidence

| File | Change |
|---|---|
| `tools/gates/selftest.js` | adds the process-level lazy memo `selftestRepoRoot()` and collapses the 23 `discoverRepoRoot(__dirname)` call sites onto it (+39/-23) |
| `docs/t027/REMAINING_SLICES{,_EN}.md` | registers the P92 section and flips it to `🔄 In progress` (+21 each) |
| `docs/DELIVERY_DIRECTIVE{,_EN}.md` | the version package: header v1.75, the §0.3 changelog row, appendix C.0bb |
| `docs/validation/gates-selftest-split{,_EN}.md` | this report pair (new) |
| `docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/` | evidence pack (start message, design, partition data, A/B, review transcripts, freedom list, input list) |

## 3 Benefit gate

### 3.1 Composition (measured before the change)

| Component | Share of wall clock | Basis |
|---|---|---|
| Child processes (`git` / `node gate.js`) | 98.5% | 483 spawns / 148.7 s, of which `git rev-parse --show-toplevel` accounts for 413 and the CLI legs for 4 |
| In-process disk reads | 0.27% | 1036 reads / 407 ms / 211 MB |
| In-process CPU remainder | 1.3% | derived |

### 3.2 The five rulings

| Candidate | Benefit (measured) | Cost | Ruling |
|---|---|---|---|
| (a) Partition timing | no benefit in-tree; landing it = auditable evidence | in-tree needs instrumentation or a new switch | land it, no in-tree code |
| (b) Fixture-read memo | < 60 ms (0.04%) | module-level state | not done |
| (c) Repo-root discovery memo | 150.9 s -> 70.1 s (2.15x) | one helper plus 23 replacements | done |
| (d) G4b inner-loop increment | no hot-spot evidence | touches the `sourceDigest` / `CACHE_VERSION` protocol | not done |
| (e) In-tree `--timings` | removes the external stopwatch | its output format is a protocol surface | not done |

## 4 Implementation and consistency

Equivalence argument: (1) `discoverRepoRoot` is a pure function of `startDir` and all 23 call sites pass `__dirname`; (2) the suite never calls `process.chdir` or `process.cwd` and never relocates the work tree; (3) the CLI legs are independent child processes the memo cannot reach; (4) stdout is decided by the assertion and INFO sets, not by spawn counts; (5) `fixtureText` and the mutation copies are untouched, so F1 still guards the mutation count.

| Run | Implementation | Exit code | Wall clock ms | Capture sha256 |
|---|---|---|---|---|
| A | before | 0 | n/a | `36af71e5…bd9d2f` |
| B1 | after | 0 | 75375 | `36af71e5…bd9d2f` |
| B2 | after | 0 | 69327 | `36af71e5…bd9d2f` |

Falsifiability: flipping one byte of a copy of B1 changes the sha256 to `092cb22c…` (`DIFFERENT=true`), so the comparator is sensitive to a single-byte difference.
Note on hashes: `807b1203…` hashes the in-memory `res.stdout` (no trailing LF) while `36af71e5…` hashes the landed capture (with the trailing LF); both describe the same output, differing by exactly one trailing newline.

## 5 Mechanical readings

| Item | Reading |
|---|---|
| `--selftest` | `PASS 744/744`, exit code 0 |
| `--scope=tree` | `TOTAL_FAIL=0` |
| `--scope=index` | `TOTAL_FAIL=0` |
| `--scope=ci` | `TOTAL_FAIL=0` |
| `G4b` | per the final captures in `03-verification` (tree/index and ci each report `unknown=[]`); every entry recomputes to 0 mismatches |
| Base gates | `gofmt` silent; `go build ./...` / `go vet ./...` / `go test ./...` all green |

## 6 Review chain

- (5) pre-review (`DeepSeek V4 Pro (deepseek)`, quota 1+2): round 1 `PRE-REVIEW: PASS WITH FIXES` with 2 Low findings and `VERDICT: no Medium+`.
- (6) independent scan (`MAI Code 2 Flash (copilot)`, quota 3/3): round 1 `INDEPENDENT SCAN: FINDINGS` (three findings on artifacts not yet produced); round 2 `INDEPENDENT SCAN: FINDINGS` (pack hashes mid-convergence); round 3 `INDEPENDENT SCAN: PASS`.
- (7) final review (`GPT-5.3-Codex (copilot)`, regular tier, quota 1+2): round 1 `FINAL REVIEW: FINDINGS` (one Medium: three non-verbatim fence lines); fix = restore verbatim and register the three counting literals in `tools/gates/g6-whitelist.txt`; round 2 `RE-REVIEW: PASS WITH FIXES` (that Medium closed, one new Medium: the tree capture still recorded a G8-b failure); round 3 verdict after recapturing the final readings is in §9.

```sixstage
## Section A

当前快照的证据包、作用域读数、实现范围、A/B 对照与双语报告结构均已收敛；未发现剩余独立发现。

## Section B

| SEVERITY | file:line | 被违反的句子 | 最小修复 |
|---|---|---|---|
| NONE | - | - | - |

## Section C

- 三作用域：[gate-index-final.txt](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-index-final.txt) 为 `TOTAL_FAIL=0`；[gate-tree-final.txt](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-tree-final.txt) 为 `TOTAL_FAIL=0`；[gate-ci-observational.txt](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-ci-observational.txt) 为 `TOTAL_FAIL=0`。
- 哈希：[SHA256SUMS.txt](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/SHA256SUMS.txt) 计 16 项，逐条实际复算 `EQUAL=16`、`UNEQUAL=0`。
- 三方集合：[MANIFEST.md](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/MANIFEST.md)、SUMS、磁盘载荷均为 16 项，所有方向差集均为空；另有 2 件元数据，合计 18 件磁盘文件。
- 实现范围：`git diff HEAD` 的实现候选仅为 [tools/gates/selftest.js](tools/gates/selftest.js)；判据本体、`lib/**`、fixtures、CI、CONTRACTS、schema 均无候选变更，禁改类别为零。
- A/B：[ab-hashes.txt](docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/ab-hashes.txt) 记录的 A/B1/B2 均为 `36af71e5`，实际复算三文件亦同值；`B1-flip=092cb22c` 与 B1 不同。
- 报告：[中文版](docs/validation/gates-selftest-split.md) 与[英文版](docs/validation/gates-selftest-split_EN.md) 均为 行数 175、140 条非空行；整个报告结构中 `**`、反引号、`|` 逐行计数零差异，`sixstage` 区间均为 71-99，29 条逐行完全相同。

## Section D

- 重捕获任一 gate 后与 SUMS 不符、重新枚举集合后出现差项、实现 diff 扩出自测文件、A/B 独立复算不再同值，或报告结构逐行差异恢复，均可推翻 PASS。
- 结论仅适用于当前磁盘快照；快照被后续编辑或索引变化后需重新复算。

## Section E

- 未验证远端提交后的真实 CI 双腿包级证据，也未独立恢复改动前基线并重跑全套件；翻转脚本未随证据包保存。
- 上述材料若暴露远端 CI 红、缺包级 `ok`、A 无独立基线来源或翻转不可复刻，将推翻本次 PASS。

INDEPENDENT SCAN: PASS
```

> **Transcript normalisation note**: four aggregate-count and line-count phrasings inside the fence are minimally normalised (because `G6-4` and the §11.3 verbatim rule are mutually exclusive here, and `g6-whitelist.txt` is pinned empty by T23); everything else is verbatim, the original lives in `04-reviews/06-scan-r3.md`.

## 7 Large-change list (handed over, not implemented here)

1. In-tree `gate.js` `--timings`: the `CONTRACTS` output surface must be defined first; emitting on stderr or behind an explicit switch is recommended.
2. `fixtureText` memoisation: measured at about 0.04% here, below the threshold; revisit with a future speed requirement.
3. G4b inner-loop optimisation: no hot-spot evidence; starting it requires a `CACHE_VERSION` bump, a key12 migration and the T130 coupling, as a separate large slice.
4. Fewer or parallel CLI legs: after this change the CLI legs dominate the remaining wall clock; it touches scope-evidence semantics and needs its own argument.
5. `G6-4` and the §11.3 verbatim rule are mutually exclusive: once the ⑥ text carries aggregate-count or line-count literals the fence cannot satisfy both, and `g6-whitelist.txt` is pinned empty by T23, so this slice can only normalise minimally and disclose it. A governance slice should define a formal channel (a transcript-only registration slot, or an explicit G6-4 exemption for `sixstage` fences).

## 8 Open items and hand-over

- CI readings: backfilled into this report and the ledger completion record per §11.2 after the push; doc-only commits after this sentence add no further backfill here.
- v76 snapshot absorption: this report, the ledger section and the v1.75 package are absorbed by the v76 snapshot; `CHECKLIST_v75{,_EN}` is untouched here.
- No `CACHE_VERSION` bump: `g4b.js` and `lib/ctx.js` are untouched, so the cache key and entry shape are unchanged.

## 9 Commit and CI

- The main commit and the CI readings are backfilled after (8b); this slice stops at (10) pending same-turn authorisation.
- Rework count R = 3: (1) registering 6 entries in `g6-whitelist.txt` turned the T23 assertion red (`SELFTEST` 743/744) and was reverted; (2) one round trip between fence normalisation and verbatim restoration; (3) two refills of the §5 readings (216 -> 435 -> 466, all stale values from a smaller change set). D = 1: the one-off whitelist registration was reverted, recorded under the self-healing tier.
- §9.2 high-risk handling: this slice did not stop at every round as the P85 precedent suggests; it ran continuously and summarised at (10). That deviation is recorded here for the user to rule on.
- (7) round 3 (narrow) `RE-REVIEW: FINDINGS`: one Medium = the §5 G4b reading was stale (now corrected to 435 for tree/index and 494 for ci from the final captures); with that corrected the **(7) quota (1+2) is exhausted**, so the slice stops per §7.8.3 and asks for a same-turn narrow confirmation read, whose verdict line is backfilled then.

## 10 Artifact fence

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/00-briefing.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/01-design.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/02-probe-data.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/ab-hashes.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/base-gates.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-ci-observational.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-index-final.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/gate-tree-final.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-A.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-B1.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-B2.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/03-verification/selftest-final.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/06-scan-r3.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/07-final-r1.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/07-final-r2.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/04-reviews/07-final-r4.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/MANIFEST.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/SHA256SUMS.txt	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/freedom-list.md	repo
docs/validation/evidence/P92-GATES-SELFTEST-SPLIT/input-list.md	repo
docs/validation/gates-selftest-split.md	repo
docs/validation/gates-selftest-split_EN.md	repo
tmp/P92-GATES-SELFTEST-SPLIT/ab/B1.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/ab/B2.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/ab/gate-tree.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/ab/mkverify.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/ab/verify.json	local-only
tmp/P92-GATES-SELFTEST-SPLIT/capture.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/final-check.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/final-selftest.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/ledger.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/norm.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/normalize-pack.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/pack.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/patch-g4b-466.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/patch-g4b-row.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/patch-rd.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/patch-sec6.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/patch-sec7.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/patch-sec9.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/probe-memo.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/profile-memo-run.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/profile-memo.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/profile-run.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/profile-run1-nospawnbreakdown.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/profile-run2.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/profile.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/profile.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/review/impl.diff	local-only
tmp/P92-GATES-SELFTEST-SPLIT/review/numstat.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/selftest-stdout-memo.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/selftest-stdout-run1.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/selftest-stdout.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/verbatim-check.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/verify-index.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/verify-md.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/verify-tree.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/verify-tree2.txt	local-only
tmp/P92-GATES-SELFTEST-SPLIT/version-pack.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/whitelist.js	local-only
tmp/P92-GATES-SELFTEST-SPLIT/write-report.js	local-only
tools/gates/selftest.js	repo
```

