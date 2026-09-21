# SW-3 · Role-set versioning and encoding normalisation — validation report

Date: 2026-09-21. Status: `IN_REVIEW` (⑥ and ⑦ in §6; commit and push in §8).

## 0. Honesty rules

1. This slice does **exactly three things**: **encoding normalisation** of 11 pre-existing role files (no-BOM+CRLF → BOM+LF), a **uniform header note across 14 files**, and an **amendment to directive R2.5**. It changes **no** `model` / `tools` value and **no** body text of duties or constraints.
2. Every byte count, reconciliation result and gate exit code below is **taken from actual command output**; anything not executed is listed in §7.
3. **Known-unfixed item**: **8** of the pre-existing files default to a §2.4 blacklisted `model` (Terra 3 / Luna 2 / Gemini 1 / GPT-5.4 1 / GPT-5.6-Sol 1) ⇒ registered as **OB-8**; this slice **does not touch it** (no scope creep).

## 1. Scope and boundaries

- Authorisation: on 2026-09-21 the user ruled for **option A** (normalise + version everything) and supplied the header-note text and the new R2.5 wording.
- Boundaries respected: no `model` / `tools` value changed; no duty/constraint body text changed; the `sol-orchestrator/` toolchain was not touched.

## 2. Deliverables

| # | Deliverable | Content |
|---|---|---|
| 1 | 11 pre-existing role files | BOM+LF normalisation + one header-note line; **all brought under version control** |
| 2 | 3 `prfrail-*` role files | header note unified to the user-specified text (replacing the earlier wording) |
| 3 | `docs/DELIVERY_DIRECTIVE.md` | R2.5 reworded (user's wording) + header-note hard rule with a mechanical criterion + OB-8 + C.0c status; version v1.4 → **v1.5** |
| 4 | `tools/gates/cjk-newwords.txt` | 4 new legitimate characters registered (诱/横/幅/犯) plus 2 pre-existing ones (刀/月) |

## 3. Falsifiability (the two decisive assertions of this slice)

| # | Assertion | Mechanical criterion | Result |
|---|---|---|---|
| 1 | **Only EOL/BOM/header-note changed** (no semantic change) | byte reconciliation: `orig bytes − line breaks + 3(BOM) + 124(note) + 2 = new bytes` | **11/11 exact** ✅ |
| 2 | Same (independent corroboration) | comparison against the generator's own copies `.github/agents/sol-orchestrator/generated/sub/_*.agent.md`: **after removing the header-note block (note line + adjacent blank line) and unifying line endings, the full string must be byte-identical** (no blank-line collapsing, no trimming) | **10/10 byte-identical** ✅ (the 11th, `sol-orchestrator.agent.md`, has no `sub/` copy and relies on assertion 1 only) |
| 3 | All 14 files carry the header note | per-file `includes(note text)` | **14/14** ✅ |
| 4 | Everything versioned and not ignored | `git check-ignore -v` | `prfrail-*` and `deep-reasoner` print **nothing**; `sol-orchestrator/_config.yaml`, `*.managed.json` **still hit** ✅ |
| 5 | `prfrail-*` carry no `model` key (G5-a) | frontmatter regex | **0 hits** ✅ |

> **Counter-inspection**: the guard for assertion 1 is "the normalisation script both strips `\r` and prepends the BOM"; if the script only added the note without touching EOL, the reconciliation would yield `orig bytes + 124 + 2` (off by 21/95 and 3) ⇒ **it goes red**.
> **Two traps in the reconciliation (measured)**: ① the BOM's **3** bytes must be counted; ② 22 `split('\n')` elements mean **21** line breaks (`sol-orchestrator` is 96/95). Missing either produces a "false FAIL".

## 4. Gate results

| Gate | Result |
|---|---|
| `gofmt -l .` | empty ✅ |
| `go build ./...` / `go vet ./...` / `go test ./...` | exit 0 / 0 / 0; `ok=14`, `FAIL=0` ✅ |
| G1-a table structure (changed lines) | PASS ✅ |
| G1-b version-metadata consistency | PASS (header v1.5 = last log row v1.5; 6 rows, no duplicates) ✅ |
| G2 placeholder residue | PASS (0) ✅ |
| G4a encoding (16 files) | PASS ✅ |
| G4b anomalous characters (⊖ whitelist) | PASS (0) ✅ |
| G5-a role frontmatter hard gate | PASS ✅ |

## 5. Consistency write-back

- The directive gained a **v1.5** row in §0.3; Appendix **C.0c** moved from "awaiting authorisation" to **✅ done**, and its deliverable/acceptance wording was synchronised with the actual scope (including the header note), so the old "status contradicts scope" defect cannot recur.
- R2.5 now also records the **header-note hard rule** (verbatim text + a `grep -L` mechanical criterion) and the **OB-8** known hazard.

## 6. Review

### 6.1 ⑥ ③.5 independent scan (SW-3's own budget, call 1)

**Verdict line**: `INDEPENDENT SCAN: PASS` (0 findings). **No list was attached in round 1**; Sections A–E were format-compliant, Section D gave the per-guard "which check goes red" audit, and Section E listed "what I cannot verify" and "what would overturn my PASS". Input: `tmp/sw3-review.diff` + read-only repo. Metering: SW-3's own budget **1/3**.

### 6.2 ⑦ independent final review (round 3, user authorisation to break the §7.5 cap once)

**Verdict line**: `RE-REVIEW: FINDINGS` (3 Medium + 1 Low). **R1–R8 re-check: 7 closed, R7 (G5-a robustness) partially closed**; this round's 4 findings **were each verified as valid and remediated**:

| # | Sev | Finding (summary) | Disposition |
|---|---|---|---|
| T1 | Medium | the "9 pre-existing files default to a blacklisted `model`" **count was wrong**; the real count is **8**, and no family breakdown was given | Fixed: R2.5 / OB-8 / this report all say **8**, with the breakdown added (Terra 3 / Luna 2 / Gemini 1 / GPT-5.4 1 / GPT-5.6-Sol 1) |
| T2 | Medium | assertion 2's "10/10 line-identical" wording did not match the measurement - the **header-note block** (note line + adjacent blank line) must be removed, not just the note line | Fixed: wording now reads "after removing the note block and unifying EOLs, the full string is byte-identical (no collapsing, no trimming)", with the copy path pinned; **re-run gives 10/10** |
| T3 | Medium | G5-a matched only a fixed shape and covered only `prfrail-*` ⇒ whitespace variants and the other 11 files slipped through | Fixed: G5-a hardened into **two criteria** (① `prfrail-*` carry no `model` key; ② **all** `.agent.md` carry no empty `model`, covering `^\s*model\s*:\s*$` / `model: ""` / `model: ''`) |
| T4 | Low | the header-note criterion only checked "present in the file", not position ⇒ moving it into the body still passed | Fixed: the criterion is now **position-anchored** (the first non-empty line after the frontmatter must equal the note); and the 3 `prfrail-*` **note blocks were moved above the H1**, making all 14 files uniform (measured **14/14**) |

**Round 3 authorisation and cost record (explicitly requested by the user)**: **`⑦ round 3 is an exception authorisation: §5.3 and §7.5 conflict, and the user ruled for closure discipline; this is a one-off break and does not constitute a new cap`** - the §7.5 ⑦ allowance is 2 (final review 1 + re-review 1) but this slice used **3** (including the authorised round 3); that precedence applies **only here and is not written into the directive** (so it cannot be read as a new cap).

**Round 4**: per the user's ruling it is **not auto-authorised**; this round's 4 remediations are **self-verified mechanically only** (no further independent re-review) ⇒ that residue is left to the user's ruling.

## 7. Explicitly not executed

1. **The 8 pre-existing files whose default `model` is blacklisted were not fixed** (OB-8, awaiting the user's ruling).
2. **The generator was not re-run**, and `_config.yaml` was not brought under version control (R2.5 already states the precondition "if the generator is ever re-enabled").
3. **`DELIVERY_DIRECTIVE_EN.md` was not generated** (user's ruling: generate it after ⑦ round 3 passes).

## 8. Commit and push

- Commit: explicit per-file `git add` (**`-A` forbidden**).
- Push: **`origin` only** (gitee untouched); per the user's ruling, **⑦ first, then push**.
- CI: observe after pushing and **report both legs green** (§6.4).
