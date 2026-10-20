# Slice `GATES-PATCH` validation report (directive v1.40)

> Slice: `GATES-PATCH`; directive version v1.40; tier L; this file is the **⑧a draft** (the ⑨ and ⑩ readings are added by the ⑧b finalisation).
> Base HEAD = `2482971` (read at ⑧a start); index and worktree consistency is re-checked by the controller at ⑧b.
> Bilingual: this file is the CN authority; `docs/validation/directive-gates-patch_EN.md` is the section-by-section mirror; both sides carry an equal number of lines.

## 1 Change summary

Before ⑧a the change set held 13 paths; after ⑧a it holds 16 paths. Five criterion files under `tools/**` are code changes; the rest are documents and ledger entries.

| Path | + | − | Notes |
|---|---|---|---|
| `docs/ADR_REGISTER.md` | 1 | 0 | ADR-027 registry row |
| `docs/ADR_REGISTER_EN.md` | 1 | 0 | `_EN` in step |
| `docs/DELIVERY_DIRECTIVE.md` | 43 | 21 | v1.40 directive body |
| `docs/DELIVERY_DIRECTIVE_EN.md` | 43 | 21 | `_EN` strict mirror |
| `docs/DEV_PLAN.md` | 1 | 1 | wording reshape of an expiring literal |
| `docs/DEV_PLAN_EN.md` | 1 | 1 | `_EN` in step |
| `docs/t027/REMAINING_SLICES.md` | 2 | 2 | P2 status flipped to done plus the residual-transfer note |
| `docs/t027/REMAINING_SLICES_EN.md` | 2 | 2 | `_EN` in step |
| `tools/gates/gate.js` | 2 | 0 | registers G5-d |
| `tools/gates/lib/criteria/g5d.js` | 143 | 0 | new G5-d |
| `tools/gates/lib/criteria/g6.js` | 203 | 13 | continuation set gains `:` and G6-7 |
| `tools/gates/lib/criteria/g9.js` | 87 | 6 | G9-e |
| `tools/gates/selftest.js` | 494 | 10 | assertions 552 → 588 |

**artifacts fence** (16 paths, all declared `repo`):

```artifacts
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-gates-patch.md	repo
docs/validation/directive-gates-patch_EN.md	repo
docs/validation/evidence/GATES-PATCH-freedom-list.md	repo
tools/gates/gate.js	repo
tools/gates/lib/criteria/g5d.js	repo
tools/gates/lib/criteria/g6.js	repo
tools/gates/lib/criteria/g9.js	repo
tools/gates/selftest.js	repo
```

## 2 Execution pipeline

The controller dispatched the steps serially; the readings of each step are below.

| Step | Role / carrier | Model | Conclusion | Rounds | Notes |
|---|---|---|---|---|---|
| ① design | controller and `deep-reasoner` | `DeepSeek V4 Pro (deepseek)` | design complete | 1 | four work groups plus the SELFTEST increment table |
| ② implementation | `prfrail-implementer` | controller model | complete | 3 | criterion code, documents and two fix rounds |
| ③ test | `prfrail-tester` | controller model | complete | 1 | three-scope readings and 5 Low items |
| ⑤ pre-review r1 | `prfrail-tester` | `GLM-5.3 (glm)` | PASS WITH FIXES | 1 | 2 Medium + 4 Low |
| ② fix | `prfrail-implementer` | controller model | complete | 1 | two Medium closed + 2 assertions + 6 number syncs |
| ⑤ pre-review r2 | `prfrail-tester` | `GLM-5.3 (glm)` | PASS | 1 (narrowed) | both Medium CLOSED, zero new findings |
| ⑥ independent blind scan | `quick-verifier` | `MAI-Code-1.1-Flash (copilot)` | PASS | 1 | blind input, zero Medium+ |
| ⑦ terminal review r1 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | PASS | 1 | all five groups pass, zero Medium+ |
| ⑦ re-reviews r2 - r3 | `independent-reviewer` | `GPT-5.3-Codex (copilot)` | r2 FINDINGS (1 Medium) -> r3 PASS | 2 (narrowed) | the freedom list missed the §5.1 `*.raw.txt` priority order; fix = F12 appended |
| ⑧a draft | `prfrail-documenter` | controller-selected | this file | 1 | three files written |

- Rework convention: ② had 3 dispatches (initial plus 2 fix rounds); ⑤ had 2 (quota 1+1 exhausted); ⑦ had 3 (terminal r1 plus narrowed r2 and r3).
- The duration field is missing (historical): this slice stored no start or end timestamps, and none is invented.
- The ⑦ r2 / r3 disposition is in §5 (back-filled by ⑧b).

## 3 Exceptions and fallbacks

| # | Event | Impact | Disposition |
|---|---|---|---|
| D1 | two mechanical normalisations in step 0 (a one-line rewrite in §12.5, blank-line normalisation in §5.1) | text shape changes, meaning unchanged | Resolved |
| D2 | ① expected 40 new assertions (→592); actual 34 (→586); r2 added 2 (→588) | the gap comes from narrowing after the T88 pre-gate downgrade | Recorded |
| D3 | the 5 Low items from ③ (D-1 to D-5) | nothing blocking | Recorded |
| D4 | the two candidate shapes were downgraded on "sampled 4, 2 unsatisfiable" | the claim is conservative against the independent reviewer's full-population measurement | Recorded |
| D5 | the controller once wrote the ⑤ package build script to the `tmp/` root | moved to `tmp/GATES-PATCH/` immediately; this is the second instance behind OB-100 | Resolved |
| D6 | after a fix round ② deleted the whole `tmp/GATES-PATCH/` (including the ⑤ r1 input package and the replay script) | beyond its own product scope; the controller rebuilt the r2 package from the original records | Fixed |
| D7 | ② used rendered text as `oldString` in the first fix round and missed | switched to a deterministic script over the staged index that rewrites the whole line | Resolved |
| D8 | the ⑤ quota 1+1 was exhausted | r2 re-ran in place with narrowed input; no third run | Recorded |
| D9 | the controller first built the r2 and ⑥ packages with three embedded `--all` runs | the background process hung and was killed; changed to log-first then light assembly, readings unaffected | Resolved |
| D10 | the freedom list was produced later than ⑦ r1 (§1.7 requires after ⑥ and before ⑦) | the ⑦ r1 package does not contain the list | Recorded |
| D11 | `DECL_ID_RE` only added a leading boundary (`OB-97x` still yields `OB-97`) | both sides share one constant, so it is symmetric and produces no false green | Recorded |

- The D6 deviation type is now the second occurrence; the precedent is item 2 of section 4 in `docs/validation/evidence/G9-C-freedom-list.md`.
- The independent reviewer extended D4's sample to every commit and measured 13 unequal pairs, so this slice's two-sample basis is conservative.

## 4 Environment and tier

| Item | Value |
|---|---|
| OS | Windows (local `d:\LZProjects\prfrail`) |
| Node.js | `v24.17.0` |
| Go | `go1.27.0 windows/amd64` |
| Base HEAD | `2482971` |
| Tier | L |
| ⑥ | mandatory (criterion body change, so §7.2.1 item ③ does not hold) |
| ⑦ quota | 1+3 |
| ⑤ quota | 1+1 (exhausted) |
| Batch helper | ADR-027 (`GLM-5.3-FlashX (glm)`, limited to four task classes) |
| Probes | 0 |

## 5 Review conclusions (⑤⑥⑦)

| Round | Id | Severity | Conclusion / handling | Disposition |
|---|---|---|---|---|
| ⑤ r1 | M1 | Medium | `DECL_ID_RE` lacked a leading word boundary, so `ADR-027` was cut to `DR-027` | Fixed |
| ⑤ r1 | M2 | Medium | the "scripted status" paragraph omitted G5-d / G6-7 / G9-e | Fixed |
| ⑤ r1 | L1–L4 | Low | EN verb closure, the `翻` word boundary, the `TMP_REF_RE` boundary, two back-filled ledger rows absent from the enumeration | Recorded |
| ⑤ r2 | — | — | both Medium CLOSED, zero new findings | Closed |
| ⑥ r1 | — | — | blind scan, zero Medium+ | Recorded |
| ⑦ r1 | — | — | all five groups pass, zero Medium+ | Recorded |
| ⑦ r2 | — | Medium | narrowed re-review of the freedom list: it missed the priority order of `*.raw.txt` inside the frozen evidence package (§5.1, new in v1.40) | Fixed (F12 appended) |
| ⑦ r3 | — | — | the narrowed re-check judged that Medium closed, with no new finding | Closed |

- **⑤ r1**: the two Medium items are in the table above; after the fix r2 re-ran in place with narrowed input and judged them CLOSED.
- **The 4 Low items kept from ⑤ r1**: the EN verb closure is all past participles, so the v1.40 row using the active `lands` yields 0 EN declarations; `翻` has no word boundary (`翻倍` hits); `TMP_REF_RE` has no leading boundary (`attmp/x` hits); the two historical back-filled ledger rows `OB-101` and `OB-102` do not appear in the deliverable enumeration.
- **The 5 Low items kept from ③**: see D-1 to D-5 in §7.4.
- **⑥ and ⑦**: both used blind input; ⑦ r1 passed every one of the five groups and reported no Medium+.
- **The ⑦ narrowing chain**: r2 checked the freedom list only and raised 1 Medium (the list missed the `*.raw.txt` priority order of §5.1) ⇒ fix = F12 appended ⇒ r3 judged that Medium closed, with no new finding.
- **⑦ quota used = 3 / 3**: terminal r1 + narrowed r2 + narrowed r3.

## 6 Review input packages (including blind-isolation proof)

| Package | Location | Contents | Manifest |
|---|---|---|---|
| ⑤ r1 | missing (historical; lost in the D6 over-reach cleanup) | brief / slice definition / design / full diff / gate readings | missing (historical) |
| ⑤ r2 | `tmp/GATES-PATCH/review-5r2/` | brief / hunks for the two Medium items / count sync / self-test hunks / readings / three-scope gate output | `MANIFEST.md` |
| ⑥ r1 | `tmp/GATES-PATCH/review-6/` | brief / full diff / numstat / status / three-scope gate / self-test / readings | `MANIFEST.md` |
| ⑦ r1 | missing (historical; not readable by ⑧a) | same shape as ⑥, with an adversarial brief | missing (historical) |

- The ⑥ and ⑦ briefs were dispatched blind: no ⑤ conclusions and no earlier findings attached.
- The packages contain no earlier verdict file; reviewer scratch directories are kept apart from the input packages.
- ⑥ and ⑦ r1 did not carry the freedom list (it was produced after ⑦ r1, see D10).
- Input sanitisation: the packages carry only the slice definition, design, test report, diff, numstat, gate readings, self-test output and criterion sources.

## 7 Falsifiability

### 7.1 Boundary probes (9)

| # | Probe | Expected | Actual source |
|---|---|---|---|
| E1 | `lang:zh:en` consumed as one token | `UsageError` (exit 2) | T84 positive case |
| E2 | `lang:zh:` consumed as one token | `UsageError` (exit 2) | T84 positive case |
| E3 | legal `lang:zh` still resolves to zh | pass | T84 negative case |
| E4 | `ADR-027` yields no `DR-027` | no `DR-027` declaration | new T85 assertion |
| E5 | `登记 ADR-028 与 OB-103` | only `OB-103` is produced | new T85 assertion |
| E6 | `OB-97x` with only the leading boundary | yields `OB-97` (symmetric on both sides, D11) | ③ manual probe |
| E7 | G5-d with a missing key | FAIL | T86 assertion |
| E8 | G5-d with a duplicate key, an empty value or no frontmatter | FAIL | T86 assertion |
| E9 | G9-e pointing at a non-existent path | INFO only | G9-e implementation and gate readings |

### 7.2 Historical replay

- G6-7 full historical replay: 50 commits, 96 file-path occurrences, 0 findings, 0 errors.
- The replay covers every historical commit, not a sample; the replay script was lost in the D6 cleanup, so the readings are the controller's existing record.

### 7.3 Mutation payloads

| Assertion block | Count | Payload |
|---|---|---|
| T84 | 7 | 2 positive, 3 negative, 2 mutation cases |
| T85 | 2 (new in this slice) | both directions of the leading boundary |
| T86 | 8 | the five G5-d key shapes: missing, duplicate, empty, no frontmatter; extra keys are not judged |
| T88 | 3 | after the pre-gate downgrade only the record-state payload assertions remain |

- Mutation actions and the asserted conclusions are unchanged; payload details are in M1 and M3 of the freedom list.

### 7.4 Coverage boundaries and known Low items

| Id | Content | Owner |
|---|---|---|
| D-1 | `DECL_ID_RE` accepts ASCII hyphens only; U+2011 can bypass it | Recorded |
| D-2 | `TMP_REF_RE` is case-sensitive and does not accept backslashes; OB-77's own ledger row produces 2 INFO items | Recorded |
| D-3 | `--check=G5-d` also runs the G5 module and emits one extra INFO | Recorded |
| D-4 | `--check=G6-7` exits 2 on a fixture without `lang:` | Recorded |
| D-5 | nested `.github/agents/sub/*.agent.md` files are not judged | Recorded |

- G9-e emits INFO only and does not take part in the verdict (F5).
- G6-7 mechanises the ID subclass only; wording-level reconciliation without IDs stays with ⑤⑥⑦ and G6-5.

## 8 Gate results

| Command | Result |
|---|---|
| `node tools/gates/selftest.js` | `SELFTEST: PASS 588/588` |
| `node tools/gates/gate.js --all --scope=index` | `TOTAL_FAIL=1` (G5(f) only: three ⑧a paths are untracked) |
| `node tools/gates/gate.js --all --scope=ci` | `TOTAL_FAIL=0` |
| `node tools/gates/gate.js --all --scope=tree` | `TOTAL_FAIL=1` (G7-b only: three ⑧a paths are not staged) |
| `G6-7` on the real change set | `scanned 2, declarations 3` |
| `G9-e` index and tree | before ⑧a `scanned 92 added line(s), 4 reference(s), 2 unresolvable`; after ⑧a tree `scanned 640, 16 reference(s), 14 unresolvable`, all INFO |
| `G5-d` this round | `no changed .github/agents/*.agent.md in scope` |
| Role-file census | among 14 role files only `sol-orchestrator.agent.md` lacks `user-invocable`, a pre-existing deviation |
| base `go test ./...` | run inside `--all` |

⑧a per-check readings:

| Criterion | Reading |
|---|---|
| G7-a | PASS |
| G1-a | PASS |
| G2 | PASS |
| G3-a | PASS |
| G6-4 | PASS |
| G8-a | PASS |

- The `TOTAL_FAIL=1` in the index and tree scopes comes from G5(f) (three untracked paths) and G7-b (three new paths declared `repo` but unstaged) respectively; ⑧a runs no `git add`, so the controller's staging clears both at once. The first `--all` run made G4b report one unregistered character (code point U+53F3); the whitelist was not touched, and rewording to a synonymous form brought the re-run to `unknown=[]`.
- The three-scope re-run is three independent runs, not reused results from one process.

**CN and `_EN` six shapes (`docs/DELIVERY_DIRECTIVE.md` vs `_EN`):**

| Shape | CN | EN |
|---|---|---|
| lines | 1845 | 1845 |
| headings | 122 | 122 |
| `**` | 4438 | 4438 |
| `|` | 2947 | 2947 |
| blank lines | 447 | 447 |
| backticks | 4528 | 4528 |

- Both sides are UTF-8 with BOM and LF.

### 8.1 Known gap (fact registered by ⑨)

`gate.js --list` lists `base` as the 19th id, but `gate.js --check=base` returns exit 2 with `unknown criterion id: base` ⇒ the id sets of `--list` and `--check` disagree; `base` only runs under `--all`. This is a pre-existing defect, not introduced by this slice; it is registered here and not fixed in this slice.

## 9 ⑨ native verification and commit state

**⑨ executed**: 2026-10-20, native verification by the controller; every reading below is measured.

| Item | Reading |
|---|---|
| per-criterion `--check=<ID> --scope=index` | all 18 exited 0: `G1-a` / `G1-b` / `G1` / `G2` / `G3-a` / `G3-b` / `G4a` / `G4b` / `G5-a` / `G5-b` / `G5-c` / `G5-d` / `G5` / `R2.5` / `G6` / `G7` / `G8` / `G9`; `G5-d` is INFO (no `.github/agents/*.agent.md` change in scope) and the other 17 are PASS |
| three-scope `--all` | `index` / `tree` / `ci` all report `TOTAL_FAIL=0` at exit 0 |
| `selftest.js` | `SELFTEST: PASS 588/588` (not below this slice's bar of 586) |
| clean-worktree replay | `git worktree add --detach tmp/GATES-PATCH/pristine HEAD` (HEAD = `2482971`) → `git diff --cached --binary` exports the staged patch → `git apply --index` → inside the pristine tree `gate.js --all --scope=index` gives `TOTAL_FAIL=0` and `selftest.js` gives `588/588` ⇒ the change set is self-contained and depends on neither the worktree nor untracked files; then `git worktree remove --force`, the patch file is deleted, and `git worktree list` shows the main worktree only |
| exit-code matrix | a legal criterion id → 0; `--check=<missing id>` → **2** (`unknown criterion id`); `--scope=<illegal value>` → **2** (`invalid --scope value`) |
| `G5(d)` fence closure | `G5(d) staged is a subset of declared :: 16 path(s) declared` ⇒ the artifacts fence and the staged set close in both directions |

**⑩ not executed**: the stop point waits for user authorisation; this slice is neither committed nor pushed (`git log -1 --oneline` is still `2482971`).

**CI**: not triggered (no push).

## 10 Cost and accounting

| Position | Count | Model / notes |
|---|---|---|
| ① | 1 | `DeepSeek V4 Pro (deepseek)` |
| ② | 3 | controller model |
| ③ | 1 | controller model |
| ⑤ | 2 | `GLM-5.3 (glm)`, quota 1+1 exhausted |
| ⑥ | 1 | `MAI-Code-1.1-Flash (copilot)`, quota 3 |
| ⑦ | 3 | `GPT-5.3-Codex (copilot)`, quota 1+3 (terminal r1 plus narrowed r2 / r3) |
| Batch helper | 1 | `GLM-5.3-FlashX (glm)`, ADR-027 |
| Probes | 0 | no probe quota used |
| ⑧ | 1 | ⑧a (⑧b counted separately) |

- Under tier 1+3, ⑦ actually used 1 round; this round reports no credits figure, only the checkable model and round facts.

## 11 Explicitly not executed

- No `git add`, no `git commit`, no `git push`; gitee untouched.
- ⑨ native verification and the ⑩ stop point are not run; the ⑧b finalisation and ⑧c close-out cleaning are not run.
- All three ⑦ reads ran (terminal r1, narrowed r2, narrowed r3); quota used = 3 / 3.
- CONTRACTS, schema, fixtures, workflows and `internal/**` are untouched.

## 12 Next steps

1. The controller stages the 16 paths one by one per the artifacts fence, then re-runs the three scopes of `--all`.
2. ⑨ reproduces the change set in a pristine worktree and back-fills the readings into §9.
3. ⑧c cleans `tmp/GATES-PATCH/` item by item after the records are on disk and proves it by mechanical read-back.

## 13 Execution freedom

Quoted one by one from `docs/validation/evidence/GATES-PATCH-freedom-list.md` (the list is authoritative; this table does not change its meaning).

| # | Choice | Landing point |
|---|---|---|
| F1 | the leading boundary character class of `DECL_ID_RE` is `[A-Za-z0-9_-]`, plus 2 new T85 assertions | `g6.js`, `selftest.js` |
| F2 | the `G6-7` verb closure takes 9 CN words and 7 EN words, all past or perfective | `g6.js`, §6.3 |
| F3 | `G6-7` takes the whole §12.4 plus appendix C area as evidence, and the judging line is never its own evidence | `g6.js` |
| F4 | `G6-7` does not judge when a document has neither §12.4 nor appendix C | `g6.js`, §6.3 |
| F5 | `G9-e` emits INFO only, defines no whitelist path and does not take part in the verdict | `g9.js`, §6.3 |
| F6 | `G5-d` judges only `.github/agents/*.agent.md`; extra keys are not judged | `g5d.js` |
| F7 | the T88 self-test block shrinks from the designed 9 assertions to 3 | `selftest.js` |
| F8 | the two candidate shapes are downgraded to record state | §6.3, §12.5 |
| F9 | the downgrade is carried by a "record-state ruling" paragraph inside the G3 row of §6.3, with no new section number | §6.3 |
| F10 | version packaging is carried by inline notes, with no new paragraph | §6.3 |
| F11 | "group 1 to group 4" is an organising unit for the report and the dispatch and does not change the criterion registry structure | this report, dispatch |

- The per-item registration of deviations D1 to D11 is in §3 and section 4 of the list; D12 (the ⑦ r2 Medium, a missing list entry) is registered in section 4 of the list.

DOCS: DONE
