# Slice validation report - `DIRECTIVE-SCAN-TEMPLATE`

Date: 2026-10-14. Slice id: `DIRECTIVE-SCAN-TEMPLATE`. Tier: `[SLICE]`. Size: S (`.md` only). Probe budget: 0.

## 1. Goal, scope and the three-condition ruling

- **Goal**: harden the ⑥ independent-scan prompt template in appendix B.6 - (a) Section A–E **front-loaded** with each section's content stated; (b) a **hard constraint** that the verdict line matches the Section B content; (c) `UNCOVERED` **explicitly registered** (not counted as findings, does not affect the verdict line).
- **Boundaries**: unchanged §7.2.1 three conditions; unchanged §7.2.2 output-contract body (operational restatement in B.6 only); unchanged §7.2.3; unchanged role set and carriers (the `agentName` / `model` strings are preserved verbatim); unchanged §2.2 / §2.4 / §5.2 / §7.5 / §6.3; unchanged CONTRACTS / schema / fixtures / workflows; no change under `internal/**`; no gitee push (this slice stops at ⑩, uncommitted and unpushed).
- **§7.2.1 three-condition ruling (① architecture)**: ① `.md`-only - holds; ② no `.go` / `.json` / `.yml` / schema / fixtures semantics change - holds; **③ does not hold** (the ⑥ part of B.6 is the review layer's **output-contract rule substance**, so rewriting it rewrites the operational text of a role contract; the B.7 startup card's "must not copy rule substance" constraint **does not apply** to B.6) ⇒ **⑥ is mandatory**. ⑤ and ⑦ each re-checked this independently and **agree**.

## 2. Artifacts

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-scan-template.md	repo
docs/validation/directive-scan-template_EN.md	repo
docs/validation/evidence/DIRECTIVE-SCAN-TEMPLATE-freedom-list.md	repo
tmp/dst/impl.diff	local-only
tmp/dst/gates-tree.txt	local-only
tmp/dst/gates-index.txt	local-only
tmp/dst/gates-ci.txt	local-only
```

> This report pair lists itself in the `repo` disposition rows (it must be `git add`-ed before ⑨, otherwise G7-b turns red).

## 3. Change summary

- **(a) Five sections front-loaded**: the ⑥ part of B.6 gains **9 inserted lines** after the "！！⑦ first round…" line - a contract heading line plus the five Section A–E rows (each carrying a `→§7.2.2-X` mapping annotation).
- **(b) Consistency hard constraint**: a new sentence "⑥ verdict-line consistency hard constraint (→§7.2.2-A / §4.1)" - no finding in Section B ⇒ `INDEPENDENT SCAN: PASS`; a finding ⇒ `INDEPENDENT SCAN: FINDINGS`.
- **(c) `UNCOVERED`**: the Section D row marks the not-covered invariants with `UNCOVERED`, and the consistency sentence declares that they are **not counted as findings and do not affect the verdict line**.
- **Carrier rows and all other original rows unchanged**: B.6's two carrier rows and the remaining original rows are preserved verbatim per ① R4 (these rows do not appear in the diff).
- **Version package**: header v1.30 → v1.31; a new v1.31 row in §0.3; a new **C.0s** ledger block in the appendix.
- **Ledger**: a new **P9** block (§1.2 full field set plus a status row).
- This slice **adds no G criterion** and creates no file outside the version package.

## 4. Freedom and deviations

- **Freedoms F-1 … F-10** are itemised in `docs/validation/evidence/DIRECTIVE-SCAN-TEMPLATE-freedom-list.md` and all are ruled in-bound. **F-8** is the **conservative choice** not to add a new sentence for the ⑤ re-review's strict-reading risk (adding one would exceed ⑤'s 1+1 re-review budget and introduce unreviewed text) ⇒ recorded truthfully in §6 instead.
- **A-1 (tool form)**: when the freedom list was created, `create_file` wrote it as **no BOM + CRLF**, contrary to the repository convention (BOM + LF, as in the existing freedom lists) ⇒ a **pure encoding normalisation** was performed (BOM + LF, zero content change, 0 replacement characters). The scope is that one file only.
- **A-2 (ruling basis disclosed)**: ①'s R2 rests on a **semantic entailment**, not a literal-mechanical entailment - that hinge is disclosed in the ① design book and in C.0s, and ⑤ (independently) and ⑦ (independently) **both ruled it an entailment**.

## 5. Review-chain outcomes

| Stage | Owner | Outcome |
|---|---|---|
| ① architecture | `DeepSeek V4 Pro (deepseek)` | `DESIGN: DONE` (R1–R10 rulings; the R2 risk point disclosed) |
| ③ consistency | master (mechanical scripts) | references and figures re-checked item by item; `G3(a)` five shapes all equal |
| ④ gates | master | tree / index `TOTAL_FAIL=0`; `SELFTEST: PASS 335/335` |
| ⑤ pre-review | `GLM-5.3 (glm)` | r1 `PRE-REVIEW: PASS WITH FIXES` (1 Low: L-1) ⇒ r2 `PRE-REVIEW: PASS` |
| ⑥ independent scan | `MAI-Code-1.1-Flash (copilot)` | `INDEPENDENT SCAN: PASS` (**format-compliant**: Section A–E complete, Section B no finding, verdict line consistent) |
| ⑦ final review | `GPT-5.3-Codex (copilot)` Extra High | `FINAL REVIEW: PASS` (0 findings; independently confirmed R2 is an entailment and `UNCOVERED` is "explicitly registered") |

## 6. Known gaps and candidates (truthfully registered; not to be read as pass)

- **No mechanical criterion backing**: this slice hardens the ⑥ output contract at the **prompt level only**; the real enforcement surface remains §7.2.2's failure disposition (format deviation ⇒ one resend; deviation again ⇒ that layer may not be declared passed). Registered twice, in the C.0s note and in the freedom list.
- **⑥ observation (4th data point)**: B3c = no Section A–E; v1.30 = no SCAN verdict line; AT24-WORDING = format-compliant but the verdict line contradicted the body; **this slice's ⑥ = format-compliant with a consistent verdict line**. Truthful caveat: this dispatch **inlined** the A–E and consistency rules into the task packet, so the result **cannot be attributed to the B.6 change alone** (n=1, confounders not separated) ⇒ acceptance evidence ⑤ stays in the observation period.
- **`UNCOVERED` strict-reading risk**: the ⑤ re-review noted there is no standalone line that speaks only about `UNCOVERED`; ⑦ ruled independently that "explicit registration" **is satisfied** (a Section D annotation plus a clause in the consistency sentence). If a later review adopts the strict reading, it is a wording-level correction that can go into a later slice.
- **Candidate C-3**: turn "B.6 ↔ §7.2.2 core-wording alignment" into a §6.3 soft-scan criterion (the long-term anti-drift measure raised by ① R9; this slice's boundary forbids touching §6.3).

## 7. ⑨ and ⑧c readings

- Encoding: all 7 changed artifacts in this slice are BOM + LF (including the new freedom list, normalised per A-1); the ⑧c full-tree per-file check reads `tracked=3604 text=3515 binary=89` with **0 violations in the non-frozen main tree**; the 17 exceptions inside the frozen evidence tree are handled under §5.1's ⑧c exemption (archived packs get self-consistency verification, not per-file encoding rewrites).
- Gates: `--scope=tree` (`changed=7`) and `--scope=index` (`changed=7`) both `TOTAL_FAIL=0`; `--selftest` `SELFTEST: PASS 335/335`; `G4b unknown=[]`; no `G6-4` or `G6-5` hit; `--scope=ci` locally reads the **previous commit's** change set (`changed=2`) ⇒ this slice's real ci reading is produced by CI **after the push**, and the local reading is observational only.
- Parity: the `G3(a)` five shapes are all equal for the directive pair (`+30/-1` vs `+30/-1`); the ledger pair `+18/-0` vs `+18/-0`; `G1-b` header and the last §0.3 row are both v1.31.
- Fence: this report pair carries an `artifacts` fence whose 11 rows all resolve (G7-a and G7-b rule PASS).
- Residue: this slice's `tmp/dst/` process artifacts are removed at ⑧c; untracked residue 0.
- Truthful registration: this report pair is edited once more after ⑧b writes back the ⑧c readings; the later final re-run readings are in the ⑩ stop-point receipt.
