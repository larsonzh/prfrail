# SW-3 · Role-set versioning and encoding normalisation — validation report

Date: 2026-09-21. Status: `IN_REVIEW` (⑥ and ⑦ in §6; commit and push in §8).

## 0. Honesty rules

1. This slice does **exactly three things**: **encoding normalisation** of 11 pre-existing role files (no-BOM+CRLF → BOM+LF), a **uniform header note across 14 files**, and an **amendment to directive R2.5**. It changes **no** `model` / `tools` value and **no** body text of duties or constraints.
2. Every byte count, reconciliation result and gate exit code below is **taken from actual command output**; anything not executed is listed in §7.
3. **Known-unfixed item**: **8** of the pre-existing files default to a §2.4 blacklisted `model` (Terra 3 / Luna 2 / Gemini 1 / GPT-5.4 1 / GPT-5.6-Sol 1) ⇒ registered as **OB-8**; this slice **does not touch it** (no scope creep). The gap was later closed by **v1.7** under the user's ruling.
4. **Terminology**: the layer identifier at the time of writing was `③.5` (inherited from the retired directive); as of directive v1.7 it is renamed **`⑥` (independent scan layer)**. The headings above keep the original identifier to preserve the record.

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

### 6.3 Cost and metering (including the independent-verification gap)

| role | model | stage | calls | note |
|---|---|---|---|---|
| master | `DeepSeek V4.1 Flash (deepseek)` | ⑧ docs / gates / assembly | throughout | bootstrap-authored SW-3 artifacts + mechanical gates |
| independent scan | `MAI-Code-1.1-Flash (copilot)` | ⑥ | **1** (SW-3's own budget 1/3) | no list attached in round 1; verdict `INDEPENDENT SCAN: PASS` |
| independent final review | `GPT-5.3-Codex (copilot)` | ⑦ round 3 | **1** (exception authorisation, see §6.2) | verdict `RE-REVIEW: FINDINGS` (3 Medium + 1 Low) |

- **Fallback count = 0**; **⑦ round 4 was not run** (the user ruled against auto-authorising it).
- **Recorded verbatim at the user's request**: "`T1–T4 整改仅经自验证 + ⑥ 独立扫描，未再跑 ④ 复审；依据用户规则；若后续准则修订在本区域发现问题，优先检查 T3/T4 判据改动。`"
- **Pending item (registered as OB-9)**: **the first action of the next directive revision** is to make **the diff of the T3/T4 criterion changes** a **focused review point** for ⑦ **and include it in that round's input** (not as a standalone extra round); it is recommended to also include this v1.7 region (G5-b / identifier disambiguation).

## 7. Explicitly not executed
1. **OB-8 (the 8 blacklisted default `model` values) was not fixed inside SW-3** (no scope creep); the gap was closed by **the later v1.7** under the user's ruling: the field-removal probe passed ⇒ the 8 files had their `model` field deleted + **G5-b** was added to prevent regression (see directive §0.3 v1.7 and §12.4 OB-8).
2. **The generator was not re-run**, and `_config.yaml` was not brought under version control (R2.5 already states the precondition "if the generator is ever re-enabled").
3. **`DELIVERY_DIRECTIVE_EN.md` was not generated inside SW-3** (the ruling then was to generate it after ⑦ round 3 passes); that mirror was produced by **the later v1.7** and committed separately.

## 8. Commit and push

- Commit: explicit per-file `git add` (**`-A` forbidden**).
- Push: **`origin` only** (gitee untouched); per the user's ruling, **⑦ first, then push**.
- CI: observe after pushing and **report both legs green** (§6.4).
