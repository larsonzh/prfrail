# Directive switchover · SW-1 (reference rewrite) + SW-2 (role carriers) — validation report

Date: 2026-09-21. Status: `COMPLETE` (SW-2 live probe passed; SW-1 rewrite done; the ⑦ independent final review is in §7).

## 0. Honesty rules (premises of every conclusion in this report)

1. This slice only performs a **directive switchover and role-carrier construction**: no `.go`, `schema/` or fixtures change, no contract semantics change, and no existing role file changes semantically.
2. Every byte count, hit count, exit code and residue count below is **taken from actual command output**; anything not executed is listed in §9.
3. `docs/DELIVERY_DIRECTIVE_EN.md` **does not exist yet** (§11.1 requires the mirror only after the CN is frozen). All references inside `_EN` documents therefore point at the **CN authoritative file** and say so; this is a **known temporary state**, not an omission.
4. This slice had the **master author product artifacts directly** - a bootstrap-period exception to the role-readiness precheck (**anchor = directive Appendix C.2**; the waived rule is **§3.2's "no original product artifacts"**), and it does **not** waive the independent review.

## 1. Scope and boundaries

**SW-2 (role carriers)** — clears the execution deadlock (Critical) found by the Codex blind review:

- Three new role files: `.github/agents/prfrail-{implementer,tester,documenter}.agent.md` (`tools` includes `edit`).
- `.gitignore` narrowed from the directory-level `.github/agents/` to `.github/agents/sol-orchestrator/` + `.github/agents/*.managed.json` + `**/__pycache__/`.
- Boundaries respected: **no existing role file was modified** (including `sol-orchestrator.agent.md` and the generator toolchain); no code changed.

**SW-1 (reference rewrite)** — removes the dual-directive state (Codex C1): **17 files touched** (§2.2). The v3.1 bodies were **not** rewritten (header pointer only); historical validation reports were **not** rewritten (related-documents line only).

## 2. Execution pipeline

### 2.1 Stage table

| Stage | Role | Carrier + model | Result |
|---|---|---|---|
| ① Architecture | — | — | **N/A**: no architectural decision |
| ② Implementation | master (bootstrap) | session body `DeepSeek V4.1 Flash` | 3 role files, `.gitignore`, 17 reference rewrites, `tools/gates/cjk-newwords.txt` |
| ③ Tests | — | — | **Degraded** (§5.2 documentation slice): replaced by reference assertions + a live probe |
| ④ Integration + gates | master | — | §6.1 four gates + G1–G5, see §6 |
| ⑤ Pre-review | — | — | **Skipped** (§5.2: optional for a documentation slice) |
| ⑥ ③.5 scan | — | — | **Skipped**, rationale in §7.1 |
| ⑦ Independent final review | independent reviewer | `independent-reviewer` + `GPT-5.3-Codex (copilot)` | §7.2 |
| ⑧ Documentation | master (bootstrap) | session body | this report + back-write |
| ⑨ Native verification | master | — | **Degraded**: encoding gate + G1–G5 + reference resolvability |
| ⑩ Stop point | — | — | awaiting commit authorization |

### 2.2 Actual SW-1 rewrite inventory

| File | Change |
|---|---|
| `docs/t027/REMAINING_SLICES{,_EN}.md` | closure-discipline line repointed to §5.3 (it restated v3.1's ①–⑥ step numbers, incompatible with the new pipeline numbering, so it is no longer restated); observation item `§3.10`→`§7.2`, `④`→`⑦` |
| `docs/DEV_PLAN{,_EN}.md` | §1 gained a pointer to the new directive (`git grep` measured 0 hits before) |
| `docs/DOCUMENTATION_PLAN{,_EN}.md` | authority-domain table gained a "delivery execution discipline" row |
| `.github/copilot-instructions.md` | "positioning and authority" section gained a pointer |
| `docs/t027/FLASH_OPERATING_DIRECTIVE.md`, `_EN.md`, `_v3.1.md`, `_v3.1_EN.md` | header gained a supersede pointer (body untouched, +1/-0) |
| `docs/t027/FLASH_T027_BRIEFING.md`, `_v3.1.md` | header gained a supersede pointer (+1/-0) |
| `docs/t027/FLASH_OPERATING_DIRECTIVE_v3.1.html` | top banner marking it a historical render (the table did not list `.html`) |
| `docs/t027/A6_PILOT_LAUNCH.md` | header gained a historical-draft pointer (**beyond the table**; 4 hits measured) |
| `docs/validation/B2-EXTERNAL-ENFORCEMENT{,_EN}.md` | related-documents line gained the new directive and marks v3.1 superseded (**beyond the table**) |

## 3. Anomalies and fallbacks

### 3.1 Bootstrap exemption

`SW-1`/`SW-2` name the documentation engineer (⑧a) as executor, but **the product-layer carriers are exactly what `SW-2` creates** ⇒ a bootstrap dependency. The master therefore authored them directly (§3.2's "no original product artifacts" does not apply here), without waiving ⑦.

### 3.2 A-1 Master damaged 11 generated role files (restored byte-for-byte)

While normalising the change set's encoding, the master **mis-scoped the operation to the whole dirty tree** and thus modified 11 pre-existing generated `.agent.md` files (no-BOM+CRLF → BOM+LF). **Detected immediately and restored byte-for-byte**, with per-file byte checks: **11/11 `RESTORED`** (576 / 603 / 576 / 580 / 601 / 629 / 528 / 527 / 591 / 4616 / 565).

- **Root cause**: the gate scope was not written down before implementation (§6.3 G4a did not distinguish "change set vs dirty tree").
- **Fixed**: §6.3 G4a now states "**scope = the change set being committed** (≠ the whole dirty tree)".
- **Fallback count**: 1 (see §8).

### 3.3 A-2 Phantom terminal interruption

Two master terminal commands echoed `^C` without executing; the probe subagent reported two interrupted `Remove-Item -Force` calls, worked around with `[System.IO.File]::Delete()`. **Recorded** in §9.3 "phantom terminal interruption" (the row predates the directive; this is its second measured reproduction).

### 3.4 A-3 Self-referential gate defects (3, fixed in place)

G2 / G3 / G4b were all found unsatisfiable the **first time they were actually executed** ⇒ OB-2 / OB-3 / OB-5. See §7.1 and directive §12.4.

### 3.5 A-4 Write tool emits CRLF (new failure mode)

Measured by the probe: a file newly written by the write tool lands as CRLF (31 B containing `0d 0a`). **Appended to** directive §9.3 per R9.1.

## 4. Environment and effort tiers

- Master: session body `DeepSeek V4.1 Flash` (tier set by the operator in the UI; **no tier record could be obtained ⇒ marked "historically missing"**, consistent with §2.3 "model and tier self-reports are untrustworthy").
- Probe subagent: `agentName='prfrail-implementer'` + `model='DeepSeek V4.1 Flash (deepseek)'` (**model passed explicitly**, R2.1).
- ⑦ review subagent: `agentName='independent-reviewer'` + `model='GPT-5.3-Codex (copilot)'`.

## 5. Falsifiability

| # | Acceptance assertion | Mechanical command | Result |
|---|---|---|---|
| 1 | New role file is not ignored | `git check-ignore -v .github/agents/prfrail-implementer.agent.md` | **no output** ✅ |
| 2 | Generator still ignored | `git check-ignore -v .github/agents/sol-orchestrator/_config.yaml` | hits `.gitignore:18` ✅ |
| 3 | Artifacts still ignored | `git check-ignore -v .github/agents/sol-orchestrator.managed.json` | hits `.gitignore:19` ✅ |
| 4 | Existing role files not ignored | `git check-ignore -v .github/agents/independent-reviewer.agent.md` | **no output** ✅ |
| 5 | Role files are seen by version control | `git status --short -uall` | the 3 `prfrail-*` files show as `??` (new) ✅ |
| 6 | **The carrier can really write files** (SW-2 core) | **four-probe matrix** (⑦ round-1 F6 + round-2 R6) | all three carriers ran create → read back → delete, all green; a **4th probe** tested `prfrail-documenter` writing `.md`: after byte-exact normalisation 55 B, starts with `EF BB BF`, `0D=0`, `0A=2`, CJK decodes without mojibake; `tmp/` entries identical before/after, **residue 0** ✅ |
| 7 | Reference rewrite has no misses | `git grep -n "FLASH_OPERATING_DIRECTIVE\|FLASH_T027_BRIEFING" -- docs` | **21 lines / 8 files**, **without a supersede note = 0** ✅ |
| 8 | Bilingual change symmetry | pairwise `git diff --numstat` | **6/6 PASS** (`+1/-0` x3, `+2/-0` x1, `+2/-2` x1, `+1/-1` x1) ✅ |
| 9 | No leftover old **section numbers** (F7) | `git grep -n "3\.10\|3\.11\|第五章\|Section 5" -- docs` | **120 lines / 24 files**; per-file classification below, **no unclassified item** ✅ |

**F7 section-number assertion classification record**: 120 hit lines across 24 files, all classified:

- **Inside the change set (handled)**: `docs/t027/A6_PILOT_LAUNCH.md` (header pointer) + the 8 retired-directive/briefing files (C.0a whitelist + header pointer).
- **Outside the change set (frozen historical records, handled per SW-1's boundary "do not widen the rewrite to historical validation reports")**: 16 files, all under `docs/validation/` (8 historical validation reports x2 languages, 6 `evidence/b3*/` evidence files, `s0-spikes` x2 languages).

**Counter-inspection ("does removing the guard turn it red?")**:

- The guard for assertions 1/2 is the two `.gitignore` rules; reverting `.github/agents/` to directory-level ignores turns assertion 1 from "no output" into a hit (**it goes red**) ⇒ the assertion is meaningful.
- The guard for assertion 6 is the carrier's `edit` + `execute` permissions; this slice proved those permissions **do exist** with **all three carriers** (control: `quick-verifier` measured `CAN_EDIT=no`, i.e. the probe cannot be completed without them).

## 6. Gate results

### 6.1 Base gates (§6.1)

| Gate | Command | Result |
|---|---|---|
| Format | `gofmt -l .` | empty output `[]` ✅ |
| Build | `go build ./...` | exit 0 ✅ |
| Static analysis | `go vet ./...` | exit 0 ✅ |
| Unit tests | `go test ./...` | exit 0; `ok=14`, `FAIL=0`, packages without tests `4` ✅ |
| Contract fixtures | — | **N/A**: this slice touched neither `schemas/` nor fixtures |
| Encoding | per-file byte check | all 30 changed files conform (`.md` = BOM+LF; `.html`/`.gitignore`/`.txt` = no BOM+LF) ✅ |

### 6.2 Closing self-check G1–G5

| # | Check | Result |
|---|---|---|
| G1 | Status-marker consistency | **PASS**: directive Appendix C.0 ✅ / C.0b ✅ agree with the §0.3 changelog; `REMAINING_SLICES` has no `SW-*` row ⇒ N/A (the SW ledger lives in Appendix C) |
| G2 | Placeholder residue | **PASS** (0 after the fix). **The first run produced 4 false positives**, all from the G2 rule line itself ⇒ OB-2 |
| G3 | Bilingual symmetry | **PASS** (6/6 + 6/6). **The first run under the original criterion produced 3 DIFFs**, all pre-existing asymmetries (501/488, 270/315) ⇒ OB-5; clause (b) key-field alignment was added (F8). **(b)'s scope was then narrowed to "changed lines"**: whole-file comparison mis-reported 10/9 on `B2-EXTERNAL-ENFORCEMENT` (a legacy `✅`); using the `git diff -U0` added lines gives **7/7 pairs at 0/0 PASS** |
| G4 | Encoding + anomalous characters | **PASS**. (a) all changed files green on BOM/EOL; (b) 10-character difference set (飞/咨/肉/忘/百/摸/纠/诱/横/幅) **verified character by character as legitimate new words**, registered in the whitelist ⇒ difference set ⊖ whitelist = **0** |
| G5 | Change-set contract | see §10 (staged with explicit per-file `git add`; **`-A` forbidden**); also run **G5-a**: `grep -l "^model:" .github/agents/prfrail-*.agent.md` must produce **no output** (no output = PASS) |

**R6.2 trigger log**: G2/G3/G4b failed on the first run ⇒ **the commit was stopped**, the criteria were fixed / observation items registered, and the full gate was re-run green.

## 7. Review

### 7.1 The ⑤ / ⑥ trade-off (including one misjudgement corrected by ⑦)

- ⑤ **pre-review**: §5.2 states ⑤ is optional for a simple slice and this slice has **no behavioural change**, so the concurrency/crash/boundary surface is empty ⇒ skipped (⑦ round 2 did not object).
- ⑥ **③.5 scan**: **the master first skipped it as a "pure documentation slice"; ⑦ round 2 ruled that invalid** - although the slice only changed `.md` / `.gitignore` / `.txt`, it **amended the §6.3 gate criteria and the evidence model**, which triggers §7.2 condition ③ ⇒ **it may not be skipped**. **⑥ was run later** (result in §7.2b).
- **⑦ is never optional** (R5.3), and this slice **amends the directive itself** (self-referential revision) ⇒ external independent review is mandatory.

### 7.2 ⑦ independent final review (round 1: blind, no list attached)

**Verdict line**: `FINDINGS` (3 High + 5 Medium). **Our per-item triage: all 8 are valid and were remediated** (dispositions below); the reviewer also explicitly marked 4 items "checked, not a problem".

| # | Sev | Finding (summary) | Disposition |
|---|---|---|---|
| F1 | High | §2.2 "carriers not yet built" contradicts C.0b "SW-2 done" | Fixed: §2.2 now reads "all seven carriers ready"; only SW-3 awaits authorization |
| F2 | High | SW-1 trigger (needs the `_EN` mirror) contradicts its status (done) | Fixed: the execution-time interpretation is stated (CN frozen is enough; the `_EN` mirror becomes a later slice) |
| F3 | High | Role files lack a `model` key, inconsistent with §2.6; reviewer suggested adding an empty one | **Measured test refutes the remedy**: `model: ""` makes the carrier **completely unloadable** (isolated A/B) ⇒ we **fixed the directive instead**: §2.6 now requires the key to be **omitted** → OB-6 |
| F4 | Medium | SW-1's "every hit" assertion contradicts "historical drafts get a header only" | Fixed: split into **file-level + section-number-level assertions**, plus a new **C.0a historical-whitelist** |
| F5 | Medium | §6.3 never states the scripting status, so "mechanical checks" reads as fully automated | Fixed: §6.3 gained a "scripting status (only the whitelist is on disk)" paragraph |
| F6 | Medium | The write probe covered only 1 of 3 carriers | Fixed: promoted to an **all-three-carrier probe matrix**, 3/3 passed (C.0b acceptance evidence ③) |
| F7 | Medium | The reference assertion searched only old file names, so old **section numbers** could slip through | Fixed: C.0's assertion now also searches section-number keywords; this slice measured 120 lines / 24 files and classified them (§5) |
| F8 | Medium | Bilingual symmetry compared line counts only, so semantic drift could still pass | Fixed: G3 gained (b) **key-field alignment** (commit ids / run ids / status enums / verdict enums) |

**Review input package**: sealed change set `tmp/sw-blind-review.diff` (93474 B; this report was **excluded** to keep the round blind) + read-only repo access. **Blind-isolation proof**: round 1 attached **none** of our findings; no ⑤/⑥ conclusions were supplied.

### 7.2b ⑥ ③.5 independent scan (run later; 3 rounds total, cap used up)

| Round | Verdict line | Finding | Disposition |
|---|---|---|---|
| 1 | `INDEPENDENT SCAN: FINDINGS` | 1 Medium: §12.4's OB-6/OB-7 were concatenated into one row (**independently overlapping** ⑦ round-2 R3) ⇒ OB-7 | Fixed (row split) + **G1-a** table-structure criterion added |
| 2 | `INDEPENDENT SCAN: FINDINGS` | 1 Medium: the header version v1.3 was **out of sync** with the §0.3 log row v1.4 | Fixed (header raised to v1.4) + **G1-b** version-metadata criterion added |
| 3 | `INDEPENDENT SCAN: PASS` | 0 | — |

- Inputs: `tmp/sw-mai-rerun.diff` + read-only repo; **no list was attached in any of the three rounds** (including the post-fix re-runs).
- **Metering**: ⑥ used the full §7.2 cap of **3** (first run + two post-fix re-runs).
- **OB-1 data point**: all three rounds produced **format-compliant** Sections A–E with substantive D/E content ⇒ the "documentation tasks necessarily deviate from the format" attribution is **further weakened**.
- **Evidence cleanup**: the three sealed packages (`tmp/sw-*.diff`) are `tmp/` scratch; per §1.5 they were **deleted once used** (never committed; this report keeps only the commands and statistics).

### 7.3 Round 2 re-review (fix list attached, per §7.3's "lists allowed on re-review rounds" rule)

**Verdict line**: `FINDINGS` (2 High + 5 Medium + 1 Low). **Our per-item triage: all 8 are valid and were remediated.**

| # | Sev | Finding (summary) | Disposition |
|---|---|---|---|
| R1 | High | The report claims `COMPLETE` yet says "re-review verdict pending", and already counts round 2 in the metering | Fixed: this section back-fills the round-2 verdict ⇒ status and metering are consistent |
| R2 | High | The ⑥-skip justification does not hold (the slice changed `.gitignore`/`.txt` **and amended the gate criteria and evidence model**) | Fixed: **⑥ was run** (§7.2b); the **§7.2 criterion was tightened** with a **counter-example** added |
| R3 | Medium | The v1.1 and v1.2 changelog rows were concatenated into one table row (structural damage) | Fixed: split into two separate rows (**⑥ also reported this independently**) |
| R4 | Medium | C.0b still said "does not exist" alongside "done" (tenses not reconciled) | Fixed: now reads "did not exist **when SW-2 was raised**" |
| R5 | Medium | The report anchored the bootstrap exemption to §3.1, while the directive anchors it to Appendix C.2 / §3.2 | Fixed: anchor and waived rule corrected in both languages |
| R6 | Medium | The three-carrier probe only created `.txt`, so `documenter`'s BOM+LF fidelity on `.md` was unverified | Fixed: a **4th probe** (`.md` fidelity) was added; it measured that the toolchain **also emits no-BOM+CRLF for `.md`** ⇒ appended to §9.3 |
| R7 | Medium | F3's guard relied on a manual precheck, with no scripted hard gate | Fixed: **G5-a** added (`grep -l "^model:"` must print nothing), listed as a mandatory scripted check |
| R8 | Low | Assertion 8's breakdown (`x4`) contradicted the stated "6 pairs" | Fixed: corrected to `+1/-0` x3, `+2/-0` x1, `+2/-2` x1, `+1/-1` x1 |

**Self-found along the way**: the directive references `Section A-E` in three places (§4.2 / §7.2 / B.6) but **never defined it** (v3.1's definition was not inherited) ⇒ the definition was restored from its inherited source (v1.3).

**§7.5 ceiling status**: ⑦ has used 2 units (final review + re-review), **exactly the cap**. §5.3 demands re-running ⑦ after remediation until no Medium+ remains, but a third round would **exceed the §7.5 cap** ⇒ per §7.5 ("exceeding the cap escalates to the user") this is **escalated to the user** (see §11).

## 8. Cost and metering (R7.1)

| role | model | stage | startedAt | reason | retryOf |
|---|---|---|---|---|---|
| master | `DeepSeek V4.1 Flash (deepseek)` | ②④⑧⑨ | whole session | bootstrap authoring + gates + report | — |
| probe | `DeepSeek V4.1 Flash (deepseek)` | ② (SW-2 probe) | 2026-09-21 | verify carrier write permission (acceptance evidence ③) | — |
| independent review | `GPT-5.3-Codex (copilot)` | ⑦ | 2026-09-21 | independent final review (first round blind) | — |
| independent scan | `MAI-Code-1.1-Flash (copilot)` | ⑥ | 2026-09-21 | late ③.5 independent scan (after the ⑤/⑥ misjudgement) | — |

- **Fallback count = 1** (A-1; threshold 3, so §9.2's high-risk marker did not fire).
- **⑥ call count**: **3** (late first run + two post-fix re-runs), the §7.2 cap of 3 is **used up**.
- **⑦ call count**: **3** (round-1 blind final review + round-2 re-review + **round 3 (exception authorisation)**). The §7.5 ⑦ cap is 2 ⇒ round 3 is an **exception**; the "mandatory blind review +1" did not fire (this slice touches no ownership/shutdown/identity semantics), and §7.3 already makes round 1 list-free, so the blind first round costs no extra unit.
- **`⑦ round 3 is an exception authorisation: §5.3 and §7.5 conflict, and the user ruled for closure discipline; this is a one-off break and does not constitute a new cap`** (user ruling, 2026-09-21; the precedence is **not written into the directive**, so it cannot be read as a new cap).
- **Paid-call ceiling**: Codex **x3** (§7.5 final review + re-review + the authorised round 3); MAI **x3** (⑥ late run + two post-fix re-runs, using up the §7.2 cap); V4 Pro x0 (①⑤ skipped).
- **Historically missing**: master tier and token counts were not recorded (§2.3: self-reports are untrustworthy).

## 9. Explicitly not executed

1. **`SW-3` (role-set versioning and encoding normalisation)**: the 11 pre-existing generated role files are still no-BOM+CRLF and unversioned ⇒ **awaiting authorization** (the two options are in directive Appendix C.0c).
2. **`docs/DELIVERY_DIRECTIVE_EN.md` not generated**: §11.1 requires the mirror after the CN freeze; this slice only rewrote references inside existing `_EN` documents.
3. **`tools/gates/` scripts not implemented**: G1–G5 were run by the master **by hand**; apart from the whitelist file no script was committed (the directive §12.5 already lists this gap).
4. **No pilot slice was run** (the §12.1 metrics need a pilot to be collected).
5. **Nothing pushed, no CI observed** (awaiting authorization).
6. **`ADR_REGISTER` not back-written**: this slice produced no new decision (a directive switchover is engineering discipline, not an ADR-level decision) ⇒ N/A.

## 10. Commit and CI

- Staging: explicit per-file `git add <specific-file>` (**`-A` forbidden**); G5 three-set agreement is recorded with the commit.
- Push: **`origin` only** (gitee untouched).
- CI: observe after pushing and **report both legs green** (§6.4); a red first run is handled with the §6.4 four steps and registered as DR-N.
- **Measured CI (2026-09-21)**: this slice plus SW-3 were pushed as head `52c28af` → run **`35606581371`**: the **first attempt failed the Windows `Test` step** (`tools/agent-probe/enforcement-proxy` · `TestConnectUsesUpstreamProxy`, `main_test.go:450: target accepts = 0, want 1`; the Ubuntu leg was green and the other 14 packages were all `ok`); `gh run rerun --failed` made **attempt 2 green on both legs**, run `success`. **A green rerun must not be used to hide a red first run** ⇒ registered as **DR-4** per §6.4 (a different test from both DR-2 and DR-3); the slice changed no `.go` file ⇒ **no causal link** to this slice's changes.

## 11. Next steps

1. **Rule on `SW-3`** (recommended: normalise the 11 files and version them; conservative alternative: narrow the `.gitignore` exception to `prfrail-*.agent.md`).
2. **Generate the `DELIVERY_DIRECTIVE_EN.md` mirror** (CN is frozen; `_EN` references can then point at the mirror).
3. **Implement the `tools/gates/` G1–G5 scripts** (with this slice's corrected G2 / G3 / G4b criteria plus the whitelist).
4. **Run the pilot slice = DR-2 + DR-3 fix**, collect the 5 §12.1 metrics, and use it as **OB-1**'s code-slice control.
5. Afterwards back-write the DR-2/DR-3 report under `docs/validation/` and the ledger's B4 precondition state.
