# Slice validation report — `DIRECTIVE-PREREVIEW-GUIDANCE` (v1.29, MINOR)

Date: 2026-10-13. Status: **completed** (⑧a / ⑧b finalized; the ⑨ and ⑧c results are in §9 and §10). Slice id: `DIRECTIVE-PREREVIEW-GUIDANCE`. Tier: `[SLICE]`. Size: M.

## 1. Change summary

- **D-1**: a new `**4.2.1 预审员的最低核验动作（五条）**` subsection is added after the §4.2 table (five items plus a closing sentence).
- **D-3a**: the §4.2 table pre-reviewer row's "Prohibited" cell becomes a **pure pointer**: `**改文件**；其余见 §4.2.1 / §4.2.2`.
- **D-3b**: a new `**4.2.2 预审员的禁止项（权威措辞）**` subsection (forbids "using test code or documentation changes themselves as the basis for judging implementation correctness"; "reading the above artifacts to verify report claims" is not restricted).
- **D-2**: the Critical / High rows of §4.3 each gain a rating criterion in their "Meaning" cell; a `> **定级边界注记**` paragraph is added after the §4.3 table.
- **D-4**: the Appendix B.5 prohibited line is synchronized to the same core wording string and points to §4.2.1.
- **Step 0**: three §12.4 rows — OB-87 gains `**n=2 追加（2026-10-13）**`; OB-88's disposition flips to `**已处置**`; **OB-89** is added (the pilot deferral continues for the 7th time).

## 2. Execution pipeline

| Stage | Role / carrier (qualified name) | Result | Note |
|---|---|---|---|
| Step 0 | controller (assertion script) | done | G4b pre-check first reported 2 unregistered characters ⇒ reworded, whitelist not extended |
| ① architecture | `deep-reasoner` / `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE` | artifacts `design-r1.md` plus amendment A |
| ② documentation implementation | `prfrail-documenter` / `DeepSeek V4.1 Flash (deepseek)` | `IMPLEMENT: DONE` | 12/12 landing points; controller re-verified 29/29 PASS |
| ③ documentation consistency | `prfrail-documenter` | `TEST: DONE` | 37 assertions plus 6 counterexamples and 2 boundary cases (in-memory mutation) |
| ④ controller integration + gates | controller | done | evidence persisted as `diff-step0-2.diff` and `gates-index-step0-2.txt` |
| ⑤ pre-review | `deep-reasoner` / `GLM-5.3 (glm)` | r1 `FINDINGS` → r2 `PASS` | see §3 and §5 |
| ⑥ independent scan | `independent-reviewer` / `MAI-Code-1.1-Flash (copilot)` | `PASS` | 0 findings; Sections A–E complete |
| ⑦ final review | `independent-reviewer` / `GPT-5.3-Codex (copilot)` | first round `PASS` plus re-review `PASS WITH FIXES` | see §5 and §6 |

## 3. Exceptions and fallback record

- **⑤ r1 `FINDINGS` (2 High plus 3 Medium)**: root cause = the **design document r1 handed to ⑤ lagged behind the user's 2026-10-13 F12 ruling**, and its change-point list did not include the three §12.4 rows from step 0 ⇒ disposition = **amend the design document (amendment A, A1–A7, a factual record)**, then rerun ⑤ per §5.3.
- **Controller self-correction D-1**: §1.7.3 requires the freedom list to be persisted "after ⑥ / before ⑦" and included in ⑦'s review input package; this slice missed it ⇒ it was persisted immediately at the fixed path, and the gap was closed using ⑦'s re-review slot (`RE-REVIEW: PASS WITH FIXES`, F-1..F-6 all in-bounds, no Medium+).
- **⑤ r1's single Low (the evidence file's "transcription noise") is not reproducible by measurement**: `before === fresh`, `+---` count 0, context ` ---` count 2, sha256 `c2086d25…` ⇒ recorded per the **OB-6 precedent** as "a reviewer finding overturned by measurement", not an implementation defect.

## 4. Environment and effort tier

Controller High; ①⑤ Max; ⑦ `GPT-5.3-Codex (copilot)` plus **Extra High (fixed)**, **not downgraded ⇒ OB-37 not triggered**. Probes **0**.

## 5. Review conclusions (⑤⑥⑦)

| Stage | Conclusion line | Findings |
|---|---|---|
| ⑤ r1 | `PRE-REVIEW: FINDINGS` | 2 High plus 3 Medium |
| ⑤ r2 | `PRE-REVIEW: PASS` | all 5 deviations CLOSED; 1 Low (not reproducible by measurement) |
| ⑥ | `INDEPENDENT SCAN: PASS` | 0 |
| ⑦ first review (blind) | `FINAL REVIEW: PASS` | 0 plus 3 "checked, not a problem" items |
| ⑦ re-review | `RE-REVIEW: PASS WITH FIXES` | no Medium+; F-2 and the report-consistency duty are deferred to ⑧b as planned |

### Findings already disposed of (each carries a disposition)

| Severity | Finding | Disposition |
|---|---|---|
| **High** | ⑤ r1 D-2a "frozen text not landed verbatim" | Closed: amendment A1 shows the source was a stale design document; the implementation matches the user's F12 ruling verbatim |
| **High** | ⑤ r1 D-2b "comment sentence rewritten" | Closed: same root cause (amendments A1 and A2) |
| **Medium** | ⑤ r1 "'定向复核' has no anchor in §4.3" | Closed: the user's ⑩ ruling chose (a); the note's parenthetical now takes disposition names verbatim from the §4.3 table (folded into this slice's commit) |
| **Medium** | ⑤ r1 "the OB-87 addition exceeds the design document's list" | Closed: amendment A3 corrects the change-set scope (step 0 included) |
| **Medium** | ⑤ r1 "OB-88's landing note lacks the D-2 content" | Closed: amendment A4 confirms the timing is ⑧b |
| **Low** | ⑤ r1 "evidence file transcription noise" | Rejected: not reproducible by measurement (OB-6 precedent) |

## 6. Review input package summary (including the blind-review isolation proof)

- ⑦ **first round is blind**: the input is only "slice definition plus final diff plus verbatim contract text plus read-only constraint"; **not attached** are the ⑤ / ⑥ reports, the ③ coverage matrix, the freedom list or the design document.
- ⑦ **re-review round**: the **freedom list** is added (the object required by §1.7.3).
- Input sanitisation: material handed to reviewers is repository `.md` only, containing no IP / host name / credential / token (§7.3).

## 7. Freedoms exercised (consistent with `docs/validation/evidence/DIRECTIVE-PREREVIEW-GUIDANCE-freedom-list.md`)

F-1 the note uses the blockquote form; F-2 the pointer wording (**the user has ruled it be changed to 「其余规则见」 at ⑧b**); F-3 the new subsections use bold line-leading labels; F-4 a numbered list plus a standalone closing sentence; F-5 the version package is deferred to ⑧b; F-6 mutations are written to in-memory copies only. The ⑧b finalization must write this list back; the two must agree.

## 8. Falsifiability (mechanical checks plus mutation testing)

- All positive assertions are green; 6 counterexamples plus 2 boundary cases (core string once / three times) are each judged red, executed by **③ on in-memory copies** without touching the target files.
- The controller independently re-verified 29 assertions (`tmp/prg/verify-impl.js`), all green.
- Key invariants: the core wording string appears **exactly twice** in CN / EN; the §4.2 pointer cell **does not** contain the core string; the CN / EN five-tuple is equal (line count 1556 / `**` 3489 / `|` 2368 / blank 392).

## 9. Gate results

`gate.js --all --scope=index` (3 files) → **`TOTAL_FAIL=0`**: G1-a, G1-b, G2, G3(a) (`+21/−6` symmetry plus all five shapes equal), G3(b), G4a (3 file(s) OK; one evidence-tree case excluded per §6.3), G4b (`newCjk=366 unknown=[]`), G5-a, G5-b, R2.5, G6-1 through G6-6, G8-a (1 path OK), G8-d and the base gates (`ok=14 FAIL=0`) all PASS.

## 10. ⑧c wrap-up cleanup record

Repository-wide per-file encoding verification (`git ls-files`, 3434 entries = 3345 text plus 89 binary) ⇒ **0 violations outside the frozen evidence tree**, with 9 frozen-tree files listed separately under the §5.1 ⑧c exemption (one `.raw.txt` CRLF in B2, two `pin.json` BOM in B3b, six under `glm-eval/archive`); `gate.js --all --scope=tree` and `--scope=index` both `TOTAL_FAIL=0`; `--selftest` PASS 335/335; residue check = only the `tmp/prg/` process files (already declared `local-only` in this report's fence), no undeclared residue. **Commit and CI**: commit `89f63f1` (push range `4fc1117..89f63f1`) has been pushed to `origin/main` (**gitee not pushed**); CI run `38368822132` (head `89f63f1`) was **green on its first run** - ubuntu 2m33s / windows 2m59s, with the `Gates` step on both legs reading `SELFTEST: PASS 335/335`, `GATE REPORT scope=ci changed=5` and `TOTAL_FAIL=0`; the ubuntu leg also ran the Contract fixtures, corroborated by package-level Go `ok` lines; the three dispatch-only jobs were skipped (normal).

## 11. Cost and metering

⑤ 2 calls (1 pre-review plus 1 post-fix re-review); ⑥ 1 call; ⑦ 2 calls (1 blind final review plus 1 re-review); ① 1 call; ② 1 call; ③ 1 call; controller ④ / ⑧a / ⑨ / ⑧b / ⑧c 1 call each; probes 0; the metering source of truth is the session call receipts (§7.5's R7.1).

## 12. Explicitly not executed

The GLM guided-arm comparison (not scheduled, per the ruling); the §7.7 "1 in 4 slices" blind review (the window count has no ledger anchor, OB-82 pending ⇒ judged not drawn, and ⑦'s first round was already executed blind).

## 13. Next-step recommendations

Next-slice recommendations: controller High; ①⑤ Max; ⑦ = `GPT-5.3-Codex (copilot)` regular tier with Extra High fixed (not downgraded when rule semantics are touched); **this slice's ⑦ budget is now exhausted at 2/2** ⇒ any further ⑦ round must be escalated to a user ruling under §7.8.3; ⑥ has 2 calls left; commit / push still rest on the same-turn authorization at ⑩.

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/validation/evidence/DIRECTIVE-PREREVIEW-GUIDANCE-freedom-list.md	repo
docs/validation/directive-prereview-guidance.md	repo
docs/validation/directive-prereview-guidance_EN.md	repo
tmp/prg/diff-step0-2.diff	local-only
tmp/prg/gates-index-step0-2.txt	local-only
tmp/prg/design-r1.md	local-only
tmp/prg/design-r1-amendment-A.md	local-only
tmp/prg/pre-review-r1.md	local-only
tmp/prg/slice-definition.md	local-only
tmp/prg/verify-impl.js	local-only
tmp/prg/land-impl.js	local-only
tmp/prg/preflight-cjk.js	local-only
tmp/prg/land-step0.js	local-only
tmp/prg/dump-evidence.js	local-only
tmp/prg/reexport-diff.js	local-only
tmp/prg/norm.js	local-only
```
> **Archive note (2026-10-13, A1-a)**: the `tmp/` items cited by this report's fence have been archived per §1.5.1 into `docs/validation/evidence/directive-history/prereview-guidance/` (with `MANIFEST.md` and `SHA256SUMS.txt`); the original paths are not rewritten, and the correspondence is recorded in that unit's `MANIFEST.md`. This is a follow-on to the v1.23 archiving execution (run at the `[ROUTINE]` tier, without bumping the directive version).

