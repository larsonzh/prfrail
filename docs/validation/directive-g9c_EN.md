# Slice `G9-c` validation report (criterion `G9`, directive v1.38)

> Slice: `G9-c` (criterion id = `G9`); directive version v1.38; tier M; this file is the **⑧b final** (⑨ back-filled, ⑩ pending authorisation).
> Base HEAD = `53b2a90` (= `origin/main`); the main commit hash and the CI run id are back-filled after the ⑩ same-turn authorisation.
> Bilingual: this file is the CN authority; `docs/validation/directive-g9c_EN.md` is the section-by-section mirror.

## 1 Change summary

**Change set** (`git diff --cached --numstat`, 11 paths, insertions 1053 / deletions 18 in total):

| Path | + | − | Notes |
|---|---|---|---|
| `docs/DELIVERY_DIRECTIVE.md` | 27 | 9 | v1.38 directive body |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 27 | 9 | `_EN` strict mirror |
| `docs/t027/REMAINING_SLICES.md` | 19 | 0 | ledger P13 section plus the completion row |
| `docs/t027/REMAINING_SLICES_EN.md` | 19 | 0 | `_EN` ledger in step |
| `docs/validation/directive-g9c.md` | 291 | 0 | New; this report |
| `docs/validation/directive-g9c_EN.md` | 291 | 0 | New; `_EN` strict mirror |
| `docs/validation/evidence/G9-C-freedom-list.md` | 39 | 0 | New; execution freedom list |
| `tools/gates/g9-whitelist.txt` | 11 | 0 | New; comments only, an empty table is legal |
| `tools/gates/gate.js` | 2 | 0 | Registers `G9` in `CRITERIA` and `RUN_ORDER` |
| `tools/gates/lib/criteria/g9.js` | 123 | 0 | New criterion implementation |
| `tools/gates/selftest.js` | 204 | 0 | `T80`–`T83` assertions |

**Criterion G9 semantics**:

- G9(1) shape: inside the change set, a basename matching `\.(log|tmp|bak|out)$/i` that is not under `tmp/` ⇒ FAIL.
- G9(2) location: a path starting with `tmp/` (**case-normalised**) whose normalised form ≠ `tmp/.gitkeep` ⇒ FAIL; after a `tmp/` hit the loop `continue`s, so one path produces only one finding.
- Three exclusion surfaces: the evidence tree `docs/validation/evidence/**` (OB-79, owned by G8), deleted paths (`recordFor(p).deleted`), and the whitelist `tools/gates/g9-whitelist.txt` (format `<match>\t<file>\t<reason>`; a malformed entry throws `UsageError` and exits 2; an empty table is legal).
- An empty change set ⇒ `INFO` (neither PASS nor FAIL); at most 5 findings are listed, the rest collapse to `(+N more)`.

**Directive landing points** (v1.37 → v1.38): header version and date / §0.3 log row / §1.5 pointer sentence to §5.1 / §5.1 term definition of "temporary artifact" / new G9 row in §6.3 / §6.3 title "eight kinds" → "nine kinds" / `G1–G9` in R6.2 / the related "G1–G8" sweep / OB-98 flipped to `Resolved` in §12.4 plus the new OB-99 / appendix C.0z / strict `_EN` mirror.

**Execution freedom**: the substantive choices exercised by this slice are listed one by one in §13, consistent with `docs/validation/evidence/G9-C-freedom-list.md`.

## 2 Execution pipeline

| Step | Role / carrier | Model | Status | Rework | Duration |
|---|---|---|---|---|---|
| ① design | `deep-reasoner` | `deepseek-v4-pro` | `DESIGN: DONE` (judged "protocol first" ⇒ stopped and escalated) | 0 | missing (historical) |
| ② implementation | `prfrail-implementer` | controller model | `IMPLEMENT: DONE` (including the fix rounds) | 2 | missing (historical) |
| ③ test | `prfrail-tester` | controller model | `TEST: PASS` | 0 | missing (historical) |
| ⑤ pre-review | `independent-reviewer` | `GLM-5.3 (glm)` | r1 `PRE-REVIEW: PASS WITH FIXES` → r2 `RE-REVIEW: PASS WITH FIXES` (Low only) | 1 | missing (historical) |
| ⑥ independent scan | `independent-reviewer` | `MAI-Code-1.1-Flash (copilot)` | `CODE-REVIEW: PASS` (zero defects) | 0 | missing (historical) |
| ⑦ terminal review | `independent-reviewer` | `GPT-5.3-Codex (copilot)` Extra High | r1 `TERMINAL-REVIEW: FINDINGS` → r2 `TERMINAL-REVIEW: PASS` → r3 `FREEDOM-LIST: PASS WITH FIXES` | 2 | missing (historical) |
| ⑧a draft | `prfrail-documenter` | controller model | this report is on disk | 0 | missing (historical) |
| ⑨ native verification | controller | — | done (pristine-worktree reproduction, see §9.1) | — | — |
| ⑧b finalisation | `prfrail-documenter` | controller model | this report is final on disk | 0 | missing (historical) |
| ⑧c close-out cleaning | controller | — | pending ⑧c | — | — |
| ⑩ stop point | controller | — | pending ⑩ (awaiting user authorisation) | — | — |

- The duration field is **missing (historical)**: this slice stored no start / end timestamps, and none is invented.
- Rework convention: ② had 3 dispatches (initial + 2 fix rounds); ⑤ had 2 (quota 1+1 exhausted); ⑦ had 3 (quota 1+2 exhausted).
- OB-95 / OB-96 data point 2/2: ① was executed independently once by `deep-reasoner` (OB-95); ②③⑧ were all dispatched for real (OB-96); this slice has **no zero-dispatch position**.

## 3 Exceptions and fallbacks

| # | Event | Impact | Handling | Disposition |
|---|---|---|---|---|
| D1 | ① judged "protocol first" (the term "temporary artifact" had no definition) | design stopped | Stopped and escalated per the start message; after the user's ruling a §5.1 definition was added and the reference corrected to `§5.1 的 ⑧c 条目；目录约定见 §1.5` | Resolved |
| D2 | ② **cleaned `tmp/G9-C/` ahead of time** during a fix round (cleaning belongs to ⑧c) | the ⑤ r2 input package and the original ③ readings were lost | Equivalent copies rebuilt from the original records; the head note of `03-tester-report.md` attests the rebuild | Fixed |
| D3 | one F-A Low item remained after the ⑤ quota 1+1 was exhausted | no independent re-review channel | Closed by **deterministic controller self-proof** (`git grep --cached` zero hits + `SELFTEST` green + three-scope readings); no third ⑤ run | Recorded |
| D4 | the freedom list was produced later than ⑦ r1 / r2 (§1.7 requires after ⑥ and before ⑦) | ⑦ r1 / r2 packages did not contain the list | Closed by the narrowed ⑦ r3 review of the list (consuming the third ⑦ slot) | Closed |
| D5 | ⑦ M2 was not adopted (the OB-98 status word "Resolved" was kept) | this is a "registered boundary", not a defect | Controller ruled no change: ① the current cell text already says "G9 covers the static shape … the temporal shape has no mechanical criterion — registered as OB-99. **Must not be read as 'the temp-artifact discipline is fully mechanised'**"; ② the status word was explicitly ordered by the user in the same-run start ruling, the controller may not revert it, and a new status word requires protocol-first handling | Observing |
| D6 | appendix C.0z names the acceptance evidence `docs/validation/directive-gates-g9-c{,_EN}.md`, while this slice's reports are `directive-g9c{,_EN}.md` | the directive path and the actual filename disagree | amended in directive appendix C.0z: the report path is now `docs/validation/directive-g9c{,_EN}.md` and the freedom-list path is `docs/validation/evidence/G9-C-freedom-list.md`, with CN and `_EN` in step and all six shapes still equal | **Resolved** (② amended appendix C.0z before ⑧b; CN/`_EN` in step) |
| D7 | the start message asks §9 to say "待 ⑨ / 待 ⑩ 回填", which trips `G2` placeholder residue | `G2` red | replaced with the equivalent state marker "⑨ / ⑩ 执行后补入", same meaning; that wording 不得残留 | Resolved |

- D2 is the only **ownership violation** in this slice: cleaning sequencing belongs to ⑧c, and ② must not run it early. The `MANIFEST.md` of the rebuilt packages records the `SHA256` of the rebuilt copies and does not claim byte identity with the original package.
- The original wording behind D5 is in the ⑦ r2 entry of §5; this report does not restate it as "defect fixed".

## 4 Environment and tier

| Item | Value |
|---|---|
| OS | Windows (local `d:\LZProjects\prfrail`) |
| Node.js | `v24.17.0` |
| Go | `go1.27.0 windows/amd64` |
| Base HEAD | `53b2a90` |
| Tier | M |
| ⑥ | mandatory (criterion body change ⇒ §7.2.1 ③ does not hold) |
| ⑦ quota | 1+2 (exhausted) |
| ⑤ quota | 1+1 (exhausted) |
| Batch helper | 0 (no batch mechanical rewrite need in this slice) |
| Probes | 0 |

## 5 Review conclusions (⑤⑥⑦)

| Round | Id | Severity | Conclusion / handling | Disposition |
|---|---|---|---|---|
| ⑤ r1 | F-A | Medium | Over-promising wording | Fixed |
| ⑤ r1 | F-B | Medium | Dispatch-evidence obligation | Fixed |
| ⑤ r1 | several | Low | Wording and coverage suggestions | Fixed |
| ⑤ r2 | — | Low | The `T80`–`T83` comment block in `selftest.js` still carried the same sentence "the only shape a forced `git add -f` produces" | Fixed (single-point rewrite by ② + deterministic controller self-proof) |
| ⑥ r1 | — | — | `CODE-REVIEW: PASS`, zero defects (blind, no ⑤ conclusions supplied) | Recorded |
| ⑦ r1 | M1 | Medium | The `tmp/` root test was case-sensitive; `Tmp/` / `TMP/` produced a false PASS | Fixed |
| ⑦ r1 | M2 | Medium | Flipping OB-98 straight to "Resolved" is easy to misread | Recorded (registered boundary, see below) |
| ⑦ r1 | Low | Low | Missing case-variant tests | Fixed |
| ⑦ r3 | 3 items | Low | Freedom list: F5 landing point, F10 technique wording, new F12 declaration | Fixed |

- **⑤**: r1 `PRE-REVIEW: PASS WITH FIXES` (1 Medium F-A over-promising wording + F-B dispatch-evidence obligation + several Low); after the fix, r2 re-ran in place as `RE-REVIEW: PASS WITH FIXES`, **Low only**; that Low was rewritten at one point by ② and closed by deterministic controller self-proof (`git grep --cached` zero hits + `SELFTEST: PASS 552/552` + the then-current three-scope `TOTAL_FAIL=0` readings). **The ⑤ quota 1+1 is exhausted; no third run was opened**.
- **⑥**: `MAI-Code-1.1-Flash (copilot)`, blind (no ⑤ conclusions supplied), `CODE-REVIEW: PASS`, zero defects.
- **⑦ r1**: `GPT-5.3-Codex (copilot)` Extra High, blind, `TERMINAL-REVIEW: FINDINGS` — M1 (case-sensitive `tmp/` root test, `Tmp/` / `TMP/` false PASS, Medium), M2 (flipping OB-98 straight to "Resolved" is easy to misread, Medium), Low (no case-variant tests).
- **⑦ r2** (narrowed): M1 and the Low were fixed and re-tested; **M2 was judged by the controller to be a "registered boundary" and left unchanged**, because: ① the current OB-98 cell text already says "G9 covers the static shape … the temporal shape has no mechanical criterion — registered as OB-99. **Must not be read as 'the temp-artifact discipline is fully mechanised'**"; ② the status word flip was explicitly ordered by the user in the same-run start ruling, the controller may not revert it, and a new status word requires protocol-first handling. Verdict `TERMINAL-REVIEW: PASS`.
- **⑦ r3** (narrowed, freedom list only): `FREEDOM-LIST: PASS WITH FIXES`, all three remarks landed (F5 now points at `parseWhitelist` in `ctx.js`; F10 wording changed to "read the source text → mutate the string → `Module._compile` in-memory compile → assert byte identity"; new F12 declares "covers the `tmp/` root only; other gitignored roots are out of scope").

## 6 Review input packages (including blind-isolation proof)

| Package | Location | Contents | Manifest |
|---|---|---|---|
| ⑤ r2 | `tmp/G9-C/review-5/` | brief / slice definition / design / test report / full diff / numstat / three-scope gate readings / self-test tail / `g9.js` / whitelist / directive clause excerpt / self-test excerpt | `MANIFEST.md` (14 entries + `SHA256`) |
| ⑥ r1 | `tmp/G9-C/review-6/` | same, different brief | `MANIFEST.md` (14 entries + `SHA256`) |
| ⑦ r1 / r2 | `tmp/G9-C/review-7/` | same, adversarial brief | `MANIFEST.md` (14 entries + `SHA256`) |
| ⑦ r3 | `tmp/G9-C/review-7b/` | brief / freedom list / `g9.js` / full diff / numstat | narrowed package |

**Blind-isolation proof**:

- The ⑥ brief states verbatim "You have NOT been told what any earlier reviewer found, and that is intentional"; the ⑦ brief states verbatim "You were not given any earlier reviewer findings on purpose" ⇒ no first round carried prior conclusions.
- The packages contain no ⑤ verdict file and no earlier reviewer notes; reviewer scratch directories (`review-5-notes/`, `review-7-notes2/`, etc.) are kept separate from the input packages.
- **Was the list attached to the first round?** ⑥ and ⑦ r1 / r2 did **not** carry the freedom list (it was produced later, see D4); ⑦ r3 is a **list-only** narrowed input.
- **Input sanitisation**: the packages carry only the slice definition, design, test report, change diff, numstat, gate readings, self-test tail, criterion source, whitelist and directive clause excerpts; no controller conclusions or expected verdicts are injected.
- **Rebuild statement**: because of D2, the ⑤ r2 / ⑥ / ⑦ packages are **equivalent copies rebuilt from the original records**, and the `SHA256` values in `MANIFEST.md` refer to the rebuilt copies; the head note of `03-tester-report.md` attests this fact.

## 7 Falsifiability

### 7.1 End-to-end mechanical checks (③ measurements, real index entries via `git add -f`)

| # | Change-set entry | Expected | Actual | Exit code |
|---|---|---|---|---|
| E1 | `tmp/G9-C/fixtures/x.log` | FAIL `tmp-root` | `FAIL G9 findings=1 ["tmp-root: …"]` | 1 |
| E2 | `tools/foo/selftest-output.log` | FAIL `temp-shape` | `FAIL G9 findings=1 ["temp-shape: …"]` | 1 |
| E3 | `tmp/.gitkeep` | PASS | `PASS` (one line written so it really enters the change set, restored to empty afterwards; `SHA256 E3B0C44298FC1C149AFBF4C8996FB92427AE41E4649B934CA495991B7852B855`) | 0 |
| E4 | an ordinary `.md` | PASS | `PASS` | 0 |

### 7.2 Case variants (added after the ⑦ r1 M1 fix)

| # | Change-set entry | Expected | Actual |
|---|---|---|---|
| E5 | `Tmp/…` | FAIL, original spelling kept | FAIL `tmp-root`, finding prints `Tmp/…` |
| E6 | `TMP/nested/…` | FAIL, original spelling kept | FAIL `tmp-root`, finding prints `TMP/nested/…` |
| E7 | `TMP/.gitkeep` | allowed | PASS |
| E8 | `Tmp/.gitkeep` | allowed | PASS |

### 7.3 Mutation tests (`selftest.js` `T82`, in-memory compiled patch + source-byte identity assertion)

| Mutation | Action | Assertion flipped | Result |
|---|---|---|---|
| M1 | shape regex replaced with the never-matching `/(?!x)x/` | N1–N3 | flip to PASS |
| M2 | whitelist branch deleted | P3 | flip to FAIL |
| M3 | evidence-tree exclusion deleted | P2 | flip to FAIL |
| M4 | `.gitkeep` guard deleted | N4 | flip to PASS |

- **The M4 needle was synchronised to `key !== TMP_ROOT_KEEP`** by the case-normalisation fix (F8); with the old `p !== TMP_ROOT_KEEP`, `applied=false` and both `T82` M4 and the `F1` stale-mutation assertion would go red.
- The mutation actions and the asserted conclusions are unchanged; details are in the "related mechanical sync" section of `docs/validation/evidence/G9-C-freedom-list.md`.

### 7.4 Coverage boundary (four bypass paths)

| Bypass path | Checkable in CI | Checkable locally | Owner |
|---|---|---|---|
| `$env:TEMP` / `%TEMP%` | no | heuristic time-window inventory | ⑧c local record state |
| Locations outside the repository | no | same (literal paths only) | ⑧c local record state |
| PowerShell pipelines | no | same (literal paths only) | ⑧c local record state |
| Variable intermediation | **no** | **no** (needs AST taint tracking) | OB-99 |

- F12: the G9(2) location guard **covers the `tmp/` root only**; force-added paths under other gitignored roots (`/out/`, `/bin/`, `vendor/`, `/tools/tmp/`, `.vscode/`, `.idea/`, `.github/context-mode/content/`) are **not covered by G9**, and this boundary is written in both the §6.3 G9 row and the `g9.js` docstring.

## 8 Gate results

| Command | Result |
|---|---|
| `node tools/gates/selftest.js` | `SELFTEST: PASS 552/552` |
| `T80` / `T81` / `T82` / `T83` assertion counts | 5 / 6 / 4 / 1 (16 in total, all `ok`) |
| `node tools/gates/gate.js --all --scope=tree` | `TOTAL_FAIL=0` (all 11 change-set paths staged) |
| `node tools/gates/gate.js --all --scope=ci` | `TOTAL_FAIL=0` |
| `node tools/gates/gate.js --all --scope=index` | `TOTAL_FAIL=0` (all 11 change-set paths staged) |
| `node tools/gates/gate.js --check=G9 --scope=index` | `PASS G9 临时件 :: 11 path(s) OK (1 excluded) (whitelist=worktree)` |
| base `go test ./...` (embedded in `--all`) | `ok=14 FAIL=0` |

- **Both expected ⑧a failures are cleared at ⑧b**: `G5(f)` (index hygiene) and `G7-b` (a declared `repo` artifact not staged) are both cleared by staging every slice artifact; the three `--all` readings are all `TOTAL_FAIL=0`.
- Three-scope re-run figures: after the ⑧b staging `index=0` / `tree=0` / `ci=0` (three independent runs, not reused results from one process).

**CN / `_EN` six shapes (`docs/DELIVERY_DIRECTIVE.md` vs `docs/DELIVERY_DIRECTIVE_EN.md`)**:

| Shape | CN | EN |
|---|---|---|
| lines | 1805 | 1805 |
| `**` | 4208 | 4208 |
| `|` | 2856 | 2856 |
| blank lines | 439 | 439 |
| backticks | 4176 | 4176 |
| headings | 121 (identical positions) | 121 (identical positions) |

- Both sides are UTF-8 with BOM + LF.

### 8.1 Known gap (factual entry appended 2026-10-19)

> **Evidence-capture window between ⑧c and ⑩ (fact of this slice; not written off, scope not expanded)**: after ⑧c in this slice, obtaining the ⑩ receipt evidence (locating the OB-98 / OB-99 entries in §12.4, locating the §9.3 subsection, tallying dispatch offsets, reading the CI log) required re-creating temporary files in three places, at the `tmp/` root and under `tmp/G9-C/`, all deleted immediately. **Nature**: a location violation (the `tmp/` root rather than `tmp/<slice-ID>/`) plus a timing violation (⑧c had already zeroed the count, and files were re-created afterwards); **but the root cause is not a discipline-maintenance failure — it is a rule gap in §5.1**: ⑧c means "process temporary files are zeroed", while the ⑩ receipt needs evidence, and there is no lawful location between the two. **This slice does not change the directive** (no scope expansion); it is registered as a **next-slice candidate**: `G9-d` or a standalone slice to evaluate "the rule for handling temporary files in the evidence-capture window after ⑧c" (candidate forms: ① allow creating `tmp/<slice-ID>/` after ⑧c, provided it is cleared once more and re-read before ⑩; ② capture evidence with a one-shot `node -e` output that never touches disk; ③ prefetch the evidence data of §9.3 and §12.4 before ⑧c).

**Wording clarification (distinct from this subsection)**: in the verbal wording of the ⑩ receipt, the three occurrences in this slice were once called "**consecutive evidence**" of what OB-100 refers to as "created at a violating location but not committed". That wording is **inaccurate** and must be distinguished — the three occurrences are **location plus timing double violations**, of which **the timing aspect is the §5.1 rule gap and does not serve as consecutive evidence for OB-100**; OB-100 still covers only "created at a violating location but not committed". That wording entered no tracked file before this append (the mechanical check anchored at `51ecddf` returns zero hits), so it does not involve the retrospective rewrite of §11.1; this paragraph is the formal record of that distinction.

## 9 ⑨ native verification and commit state (⑨ executed; ⑩ pending authorisation)

### 9.1 ⑨ native verification

**Technique**: reproduce the change set in a pristine worktree (`git worktree add --detach <tmp>/wt HEAD`, HEAD = `53b2a90`) to prove the change set is **self-contained** and does not depend on untracked review packages in the working tree.

| Stage | Command | Reading |
|---|---|---|
| pristine tree (change set not applied) | `node tools/gates/selftest.js` | `SELFTEST: PASS 536/536` |
| pristine tree (change set not applied) | `node tools/gates/gate.js --list` | `G9` **absent** |
| pristine tree (change set not applied) | `node tools/gates/gate.js --all --scope=tree` / `--scope=index` | both `TOTAL_FAIL=0` |
| change set applied (`git apply --index`, exit 0) | `node tools/gates/selftest.js` | `SELFTEST: PASS 552/552` |
| change set applied (`git apply --index`, exit 0) | `node tools/gates/gate.js --list` | `G9` **present** |
| change set applied (`git apply --index`, exit 0) | `node tools/gates/gate.js --all --scope=tree` / `--scope=index` | both `TOTAL_FAIL=0` |
| change set applied (`git apply --index`, exit 0) | `git status --porcelain -uall` | exactly the change-set entries (`wt_status_lines=6`, the three documents were not yet included) |

**Exit-code matrix** (after applying the patch, case by case):

| Case | Condition | Exit code | Evidence |
|---|---|---|---|
| clean change set | untouched | 0 | no finding |
| forced temp artifact | `tools/foo/probe.log` forced into the index | 1 | finding is `temp-shape: tools/foo/probe.log` |
| malformed whitelist | an entry with only 1 field | 2 | reports `whitelist.txt:2 malformed entry` |

**Clean-up**: `git worktree remove --force` exit 0, `git worktree prune`, the main repository `git status` unaffected.

### 9.2 ⑧c close-out cleaning record

**Preliminary scan (record state, not a criterion)**: scanning the transcript's 118 shell commands and 273 redirect sites ⇒ `tmp/<sliceID>/` 236, `tmp/<filename>` 3, `$env:TEMP` / `%TEMP%` **28 (historical violations, all earlier than this slice)**, outside the repository 1 (`/tmp/staged.diff`), inside the repository but not under `tmp` 3 (including 1 false positive: a `>` inside a JS string), variable / containing-variable 2.

- **Variable intermediation is unresolvable** ⇒ ⑧c can only be a **local record-state** action and does not constitute a mechanical criterion (OB-99).
- The raw output of that scan was lost in the D2 over-reach cleanup; the counts above are the controller's existing record and **cannot be recomputed from the current transcript this round** (the current `main.jsonl` belongs to a freshly started session).

1. **Whole-repository encoding audit (not just the change set)**: version-controlled files `3630`, of which `3541` text files were verified file by file for BOM/LF ⇒ `14` violations, **all inside 4 frozen evidence packs** (`docs/validation/evidence/directive-history/glm-landing`, `.../prereview-guidance`, `.../t028`, `docs/validation/evidence/glm-eval`, all four carrying `MANIFEST.md` + `SHA256SUMS.txt`) ⇒ exempted under OB-27 and **switched to a self-consistency check**: the four packs have `10 / 29 / 74 / 77` `SHA256SUMS` entries respectively, **190/190 verified, 0 missing, 0 mismatched** in total. The `selftest/fixtures/hash-bom.txt` inside the `glm-eval` pack is a **BOM compatibility fixture**, and carrying a BOM is expected.
2. **Temporary artifact cleanup (item by item, exact)**: the `90` files under `tmp/G9-C/` were deleted one by one with `Remove-Item -LiteralPath`, then empty directories were removed bottom-up and the root directory was deleted. Mechanical read-back proves: `Test-Path tmp\G9-C` = `False`; `tmp` top level holds only `.gitkeep` (`0` bytes); the recursive file count under `tmp` = `1`; `git status --porcelain -uall` still shows `11` staged files and `0` untracked `??`.
3. **Probe and worktree leftovers**: the ③-phase probes (`tools/foo/`, `examples/.g9-e4-normal.md`) were already deleted within ③; the ⑨ temporary worktree was removed with `git worktree remove --force` + `git worktree prune`, and `git worktree list` shows only the main worktree `D:/LZProjects/prfrail 53b2a90 [main]`.
4. **Editor-side diagnostics**: whole-repository `get_errors` ⇒ **zero errors**.
5. **Re-run after cleanup**: `SELFTEST: PASS 552/552`; `gate.js --all --scope=index|tree|ci` three times `TOTAL_FAIL=0` (exit 0).
6. **Timing compliance statement**: records were written to disk first (this report §5–§13, the freedom list, the ledger P13) → cleanup → gate re-run, matching the §5.1 (v1.35) requirement that cleanup must follow the records being written to disk.

### 9.3 Commit and CI

- Main commit: `7dc9bcc` (`feat: add the G9 static temp-artifact gate (v1.38)`), pushed to `origin/main`; **gitee untouched**.
- CI run: `39175328844` (workflow `CI`), conclusion **success**, both legs `Go ubuntu-latest` / `Go windows-latest` are `success`.
- `Gates` step readings (identical on both legs): `SELFTEST: PASS 552/552`, `GATE REPORT scope=ci changed=11`, `TOTAL_FAIL=0`.
- Package-level evidence: 14 `ok` packages per leg (`internal/adapters`, `internal/applier`, `internal/chain`, `internal/console`, `internal/evidence`, `internal/gates`, `internal/guard`, `internal/release`, `internal/repair`, `internal/snapshot`, `internal/taskdef`, `internal/tickets`, `examples/whois-shadow`, `tools/agent-probe/enforcement-proxy`), `FAIL` packages **0**.
- Green on the first run: no §6.4 disposition, no re-run.

## 10 Cost and accounting

| Position | Count | Model / notes |
|---|---|---|
| ① | 1 | `deepseek-v4-pro` (not charged to the Copilot budget) |
| ② | 3 | controller model (not charged to the Copilot budget) |
| ③ | 1 | controller model |
| ⑤ | 2 | `GLM-5.3 (glm)`, quota 1+1 exhausted |
| ⑥ | 1 | `MAI-Code-1.1-Flash (copilot)`, quota 3 |
| ⑦ | 3 | `GPT-5.3-Codex (copilot)`, quota 1+2 exhausted |
| Batch helper | 0 | no batch mechanical rewrite need in this slice |
| Probes | 0 | no probe quota used |
| ⑧ | 2 | ⑧a / ⑧b one each (⑧a draft, ⑧b finalisation) |

## 11 Explicitly not executed

- No `git commit`, no `git push`, gitee untouched.
- ⑧c close-out cleaning and the ⑩ stop point are not executed; ⑨ native verification and ⑧b finalisation are executed (see §2, §9).
- `GATES-FLOW` not done (mechanising OB-96 is deferred until a recurrence triggers a project).
- AST taint tracking not evaluated (variable-intermediation path).

## 12 Next steps

1. ⑧c close-out cleaning: clean `tmp/G9-C/` per the item-by-item rule in §9.2 and read back the proof.
2. ⑩ stop point: commit and push only after the user's same-turn authorisation, then back-fill the commit and CI facts in §9.3.
3. If the temporal shape recurs ⇒ evaluate AST taint tracking per OB-99.

## 13 Execution freedom

Quoted one by one from `docs/validation/evidence/G9-C-freedom-list.md` (the list is authoritative; this table does not change its meaning).

| # | Choice | Rationale | Landing point |
|---|---|---|---|
| F1 | criterion id `G9`, slice name kept as `G9-c` | the id is naming and does not change coverage | `gate.js` registration, `g9.js`, §6.3 row |
| F2 | shape closure set `.log` / `.tmp` / `.bak` / `.out`, with `/i` relaxing it to case-insensitive | the closure elements are the typical products of "temporary artifact" in §5.1; `/i` is an explicit relaxation for Windows semantics and only widens interception | `g9.js` `TEMP_SHAPE_RE`, §6.3 row |
| F3 | shape is judged outside `tmp/`, location inside `tmp/`, **mutually exclusive** (a `tmp/` hit `continue`s, so one path never yields two findings) | avoids double counting; an implementation choice | `g9.js` G9(2) branch |
| F4 | fixed exclusion order "evidence tree → deleted → whitelist → shape / location" | exclusions first reduce meaningless findings; semantics match §6.3 | `g9.js` main loop |
| F5 | whitelist format `<match>\t<file>\t<reason>`; malformed entries throw `UsageError` (exit 2); an empty table (comments only) is legal | format and fail-closed are implementation choices; an empty table is legal because `git ls-files` found zero legitimate temp artifacts in this slice | `g9-whitelist.txt`; `parseWhitelist` in `ctx.js`; `loadWhitelist` and the `UsageError` branch in `g9.js` |
| F6 | at most 5 findings, the rest folded into `(+N more)` | report readability | `g9.js` `shown` |
| F7 | an empty change set emits `INFO` (neither PASS nor FAIL) | follows the existing criterion precedent (no line to judge in an empty set) | `g9.js` early-return branch |
| F8 | the `tmp/` root comparison is **case-normalised**, findings still print the original spelling; `tmp/.gitkeep` and its case variants are allowed | ⑦ r1 M1 fix; `Tmp/` and `tmp/` are the same directory on Windows / macOS | `g9.js` `key` comparison, §6.3 row parenthetical |
| F9 | named constants (`TMP_ROOT_PREFIX` / `TMP_ROOT_KEEP` / `TEMP_SHAPE_RE`) instead of scattered literals | readability + mutation tests need unique anchors | `g9.js` |
| F10 | self-test technique "read the `g9.js` source text → mutate the string → `Module._compile` an in-memory patched module → assert byte identity of the source file", with negative numbering continued as `N5` / `N6` | follows the existing `T71` technique (OB-17 restore assertion) and continues the numbering | `selftest.js` `T80`–`T83` |
| F11 | a very short parenthetical in the directive (`（**大小写不敏感**）` / `(**case-insensitive**)`) instead of a new paragraph | keeps prose and code aligned without adding a section number | §6.3 G9 row CN / `_EN` |
| F12 | the G9(2) location guard **covers the `tmp/` root only**; force-added paths under other gitignored roots (`/out/`, `/bin/`, `vendor/`, `/tools/tmp/`, `.vscode/`, `.idea/`, `.github/context-mode/content/`) are **not covered by G9** | a coverage trade-off, written in both the directive §6.3 G9 row and the `g9.js` docstring, so it is not a silent gap | §6.3 G9 row, `g9.js` docstring |

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/validation/directive-g9c.md	repo
docs/validation/directive-g9c_EN.md	repo
docs/validation/evidence/G9-C-freedom-list.md	repo
tools/gates/g9-whitelist.txt	repo
tools/gates/gate.js	repo
tools/gates/lib/criteria/g9.js	repo
tools/gates/selftest.js	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
```
