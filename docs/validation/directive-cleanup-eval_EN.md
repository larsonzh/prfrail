# Slice validation report - `DIRECTIVE-CLEANUP-EVAL` (ledger cleanup plus the R2.2 argument-form assessment)

> Per the serial pipeline of `docs/DELIVERY_DIRECTIVE.md` v1.18 (§5.1) and the evidence duty of §6.2.
> **The CN file is authoritative**; this file is its positional mirror. Artifacts fence at the end.
> Scale S, tier `[SLICE]`, probes 0.

## 1. Change summary (including the execution freedom list) **Goal**: (1) clear two stale ledger statements; (2) write back the OB-66 basis; (3) complete the OB-70 assessment; (4) land OB-63. **Execution freedom exercised by this slice (§1.7.1, substantive choices)**:

- **F-1 single-character correction**: the user wording "本机为事实记录" was corrected to "**本节**为事实记录" in the report §16 first sentence; basis = the subject of that sentence is a section of the report, so a verbatim copy would be factually wrong; declared here and in the receipt.
- **F-2 landing-point choice**: the OB-63 requirement landed in **§11.2 (write-back standard)** (the user offered "§5.3 or §11.3") - because OB-63's own disposition says "wording basis in §11.2", so landing it there makes that pointer resolve; declared.
- **F-3 disposition flips**: OB-63 and OB-70 moved from `registered` / `to be discussed` to `resolved` (the requirement is now in the rule body; the assessment is complete in this slice).
- **F-4 ⑦ downgrade**: this slice ⑦ is `deepseek-v4-pro` (stated by the user's slice definition), recorded under the OB-37 meta-rule as **ADR-020**.

## 2. Pre-start checklist verification (step 0, item by item)

| Item | Content | State |
|---|---|---|
| A-1 | the T027-ASSEMBLY-MIN row in `REMAINING_SLICES{,_EN}` says "committed `68a1ef0` + `d2d6370`, CI run `37073985144` green on both legs" | **Done** (completed by the previous slice step 0; re-verified verbatim here) |
| A-2 | the DIRECTIVE-MODEL-LISTS row says "committed `35e8198`, CI run `37092310920` green on both legs" | **Done** (landed by this slice) |
| B-3 | §16 first sentence of `directive-model-lists{,_EN}.md` gains "this section is a factual record and does not rewrite §11 / §12" | **Done** (landed by this slice) |
| C-4 | the OB-66 disposition states "the §12.4 ledger is the sole source of truth" plus this slice's third deferral | **Done** (landed by this slice) |
| D-5 | probe the `modelId` form with Codex / V4 Pro / Flash | **Done** (all three rejected, see §3) |
| D-6 | check the official `runSubagent` documentation | **Done** (two official pages, see §3) |
| D-7 | assess display-name stability | **Done** (see §3) |
| D-8 | judge the three options and write them into the report (without changing the `R2.2` body or the §2.2 table) | **Done** (see §3) |
| D-9 | hand the conclusion to the next slice as a chartering basis | **Done** (recommendation and chartering points at the end of §3) |
| E-10 | land the OB-63 requirement in the rule text | **Done** (landed in §11.2, see §4) |
| existing-1 | OB-68 (cost and effort definitions in §7.5 / §8.2) | **To do** (next slice; OB-68 stays `to be discussed`) |
| existing-2 | OB-53 item ⑤ (the chain-config `model: "auto"`) | **To do** (next slice; touches `internal/**`, forbidden here) |
| existing-3 | OB-67 (the first `gpt-6-luna` rejection) | **Assessment done** (same origin as OB-70; the fix ruling stays with the next slice and OB-67 keeps its `to be discussed` disposition) |
| existing-4 | the master hand-over audit (at the end of T027) | **Not applicable** (this slice is not that anchor) |
| existing-5 | the strategic review (at the end of S1) | **Not applicable** (this slice is not that anchor) |

## 3. ① Assessment: the `model` argument form of `runSubagent` (the core deliverable) **Measurement (three probes, D-5)**: calling `runSubagent` with the `modelId` form was **rejected for all three models**, each time with the same list of available models:

| Value passed (the `modelId` form) | Result | The corresponding entry in that list (qualified name) |
|---|---|---|
| `gpt-5.3-codex` | **not found** | `GPT-5.3-Codex (copilot)` |
| `deepseek-v4-pro` | **not found** | `DeepSeek V4 Pro (deepseek)` |
| `deepseek-flash` | **not found** | `DeepSeek V4.1 Flash (deepseek)` |
 Evidence from earlier slices for contrast: `gpt-6-luna` was rejected while `GPT-6 Luna (copilot)` passed (L1 to L5 all green) and `GPT-6.1 Sol (copilot)` passed ⇒ **the rejection is about the name form, not about model availability or tier availability** (in the same session the higher-tier Luna / Sol were usable).
 **Evidence retained**: the three rejection receipts are kept verbatim in `tmp/dce/probe-id-rejections.txt` (`local-only`), so the conclusion no longer rests on prose alone.

**Official documentation (D-6, two pages)**:

- the custom-agents page's `handoffs.model` row says: **Use the qualified model name in the format
  `Model Name (vendor)`**, for example `GPT-5 (copilot)` or `Claude Sonnet 4.5 (copilot)` ⇒ **the qualified
  name (display name plus vendor) is the documented way to name a model**;
- the subagents page ("Select the model for a subagent", under the **Local harness tab**, which states that the Local harness uses the `runSubagent` tool) defines the local order: (1) the explicit `model`
  parameter the main agent passes to `runSubagent`, (2) the custom agent's `model` property, (3) Auto,
  (4) the main conversation's model; it requires replacing `<model name>` with **a model available in your
  session**, and its troubleshooting table says "**A requested model doesn't run ⇒ Use one of the models
  listed in the error**";
- the same page notes that an explicit selection is **checked against the main model's cost tier**, and an
  over-tier selection does not run but reports the available models. That clause **cannot explain this slice's
  measurement on its own** (the higher-tier Luna / Sol were usable), so the cause converges on the **name form**.

**Display-name stability (D-7)**: a qualified name **carries a version and a vendor** (`DeepSeek V4.1 Flash
(deepseek)`; the documentation examples use `Claude Sonnet 4.5 (copilot)`), and the list also contains
`Auto (copilot)` whose meaning follows the setting ⇒ both version upgrades and vendor changes alter the call
string. ⇒ `R2.2`'s original intent (preventing display-name drift) is **supported by both the measurement and
the documentation**.

**Three-option judgement (D-8; the `R2.2` body and the §2.2 table are left untouched)**:

| Option | Verdict | Reason |
|---|---|---|
| (1) rewrite `R2.2` to the display-name reading | **not recommended** | it matches the documentation but imports **version / vendor drift** into the contract, against `R2.2`'s anti-drift intent |
| (2) add a mapping table (`modelId` to qualified name) | **recommended** | task packages keep the stable `modelId` while the call layer takes the qualified name from a **single mapping table**; the `R2.2` body then needs only a pointer sentence, and the table can carry a mechanical criterion (equal member sets, qualified names copied verbatim from the tool list) |
| (3) keep `R2.2` plus a note | the current state (landed by the previous slice) | zero change and zero risk, but every call needs a manual lookup and a rule living beside a note tends to be ignored |

**D-9 hand-off to the next slice**: adopting option (2) requires a separate `[SLICE]` (changing `R2.2` is a
contract-semantics change): (1) build the mapping table (§2.2's `modelId` column against the tool list's
qualified names); (2) add the pointer sentence to `R2.2`; (3) add the mechanical criterion.

## 4. Deliverables and their landing points

| Item | Landing point |
|---|---|
| A-1 | verification only (the three strings "committed `68a1ef0`", "`d2d6370`", "`37073985144`" all hit) |
| A-2 | the DIRECTIVE-MODEL-LISTS row in `docs/t027/REMAINING_SLICES{,_EN}.md` |
| B-3 | the §16 first sentence of `docs/validation/directive-model-lists{,_EN}.md` |
| C-4 | the OB-66 disposition cell in §12.4 of `docs/DELIVERY_DIRECTIVE.md` (`ruling`) |
| E-10 | the §11.2 body of `docs/DELIVERY_DIRECTIVE.md` (the OB-63 requirement); OB-63 flipped to `resolved` |
| D | OB-70 flipped to `resolved`; the assessment is in §3 of this report |
| Records | ADR-020 in `docs/ADR_REGISTER{,_EN}.md` (the ⑦ downgrade); the existing §16 of the previous report stays untouched |

## 5. Pipeline, environment and quota

① (the assessment and document implementation were carried by the master; no separate ① call) ⇒
② document implementation ⇒ ③ checklist and consistency ⇒ ④ gates ⇒ ⑤ **skipped** (optional at tier S in §5.2)
⇒ ⑥ **skipped** (pure `.md`, no code semantics, per the §7.2 skip conditions, confirmed by ⑦) ⇒ ⑦ final review
(`deepseek-v4-pro` per **ADR-020**; **reduced independence stated truthfully**) ⇒ ⑧ report and write-back ⇒
⑨ native validation ⇒ ⑩ stop. **Quota actually used**: probes 0; the three `modelId` probes were **rejected
calls that consumed no model usage**. Environment: Windows host; Node v24.17.0; Go `go1.27.0 windows/amd64`;
repository `github.com/larsonzh/prfrail`.

## 6. Review outcome

- ⑤: **skipped** (optional at tier S).
- ⑥: **skipped** (pure documentation slice, per the §7.2 skip conditions; ⑦ confirms in its verdict line).
- ⑦: round one `FINAL REVIEW: FINDINGS` (3 Medium + 9 Low, all remediated); the re-review returned `RE-REVIEW: PASS WITH FIXES` (3 residual Low findings closed at ⑧b, with ⑦ stating that no further round is needed). This slice ⑦ ran on `deepseek-v4-pro` ⇒ reduced independence, stated in ADR-020 and here.

## 7. Falsifiability and gate results

| Check | Result |
|---|---|
| `gate.js --all --scope=tree` | see §8; the first run was red on the G3(b) to-do token (0/2) and turned green after the fix |
| `gate.js --all --scope=index` | see §8 |
| `gate.js --selftest` | see §8 |
| checklist (§2) | 15 items marked individually; A-1/A-2/B-3/C-4/D-5 to D-9/E-10 are done |

## 8. Gates and native validation

See §9 (the gate readings are filled in, including the first red run and its remediation).

## 9. Gate readings

- `--scope=tree`: `TOTAL_FAIL=0`.
- `--scope=index`: `TOTAL_FAIL=0`.
- `--selftest`: `PASS 264/264` (this slice changed no criterion or fixture, so the reading should match the slice start).
- basic gates: `gofmt` empty, `go build` / `go vet` / `go test` all exit 0.

## 10. Known divergences and open items

- OB-68 (cost and effort definitions in §7.5 / §8.2) and OB-53 item ⑤ (the chain-config `model: "auto"`) are
  **outside this slice** and stay `to be discussed`, for the next slice.
- This slice ⑦ is a downgraded model (ADR-020) ⇒ its independence is below that of an allow-list model and the
  verdict must be read together with that limitation.
- The ⑥ skip rests on "pure documentation, no code semantics" and is confirmed by ⑦'s verdict line.

## 11. Candidate OBs (none added by this slice)

None added; this slice closes out OB-63 / OB-66 / OB-70.

## 12. ⑩ Stop report

- Executed through the ⑦ final review; **stopping before commit authorisation**: this slice does not commit,
  push or push to gitee, and waits for the user authorisation in the same turn.
- **Write-back check** (OB-63 wording, applied for the first time here): validation report = written back (this
  file); `DEV_PLAN` = **not applicable** (a directive-governance slice which, per its charter and the existing
  precedent, writes no `DEV_PLAN` section of its own); ledger = written back (the OB-63 / OB-66 / OB-70 disposition
  cells in §12.4 and this slice's completion row); ADR = written back (ADR-020).

- **⑦ re-review verdict**: `RE-REVIEW: PASS WITH FIXES`; the three residual Low findings (a mid-state retained index reading, the attribution wording in the v1.19 row, and a pronoun in the previous report §16) were closed at ⑧b.

## 13. Artifact list

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/directive-model-lists.md	repo
docs/validation/directive-model-lists_EN.md	repo
docs/validation/directive-cleanup-eval.md	repo
docs/validation/directive-cleanup-eval_EN.md	repo
tmp/dce/land.js	local-only
tmp/dce/peek.js	local-only
tmp/dce/gate-tree.txt	local-only
tmp/dce/gate-index.txt	local-only
tmp/dce/selftest.txt	local-only
tmp/dce/probe-id-rejections.txt	local-only
tmp/dce/gate-tree-first-fail.txt	local-only
tmp/dce/test.txt	local-only
```
> **Archiving note (v1.23)**: the `tmp/` items cited by this report were archived under `docs/validation/evidence/directive-history/` per §1.5.1 (units: `dce/`; each with `MANIFEST.md` and `SHA256SUMS.txt`); the original paths are not rewritten and the MANIFEST is the mapping of record. (Paths already cleared by earlier slices are outside this archiving and belong to the pre-existing decoupling recorded as OB-77.)
