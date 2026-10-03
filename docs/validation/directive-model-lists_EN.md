# Slice validation report - `DIRECTIVE-MODEL-LISTS` (v1.18)

> Per the serial pipeline of `docs/DELIVERY_DIRECTIVE.md` §5.1 and the evidence duty of §6.2. **The CN file is
> authoritative**; this file is its positional mirror. The self-compliance artifacts fence is at the end.

## 1. Change summary (including the execution freedom list) **Goal**: converge the blacklist to the complete ID `kimi-k3`, add `gpt-6.1-sol` (deep tier) and `gpt-6-luna` (parallel alternative) to §2.2, add the ⑦ three-model tier table to §5.2, add the temporary-authorization five elements as §1.7.4, turn G5-b into a whitelist closed set (directive sentence, `tools/gates` implementation and self-test fixtures), and leave ledger / ADR records. **Execution freedom list**: the master wrote it as an orchestration artifact at `docs/validation/evidence/directive-model-lists-freedom-list.md` (F-1 to F-10, T-1 to T-3); this section agrees with it: F-1 the §6.3 G6-5 endpoint wording completion; F-2 the appendix B.6 carrier line rewrite and the removal of "most expensive"; F-3 the §2.5 old G5-b description now pointing at §6.3; F-4 the §2.2 readiness sentence dropping the stale "seven" count; F-5 `OB-53` not flipped to resolved (items ④ / ⑤ unlanded); F-6 the binding of the depth-tier trigger set to `gpt-6.1-sol` (to be confirmed at ⑩); F-7 the choice of OB-65 to OB-69 numbers; F-8 the Luna probe carrier, the display-name retry, and the extra Sol minimal call; F-9 `R2.2` left unchanged with OB-67 registered instead; F-10 the three extra G5-b boundary fixtures and T38 to T40.

## 2. Deliverables and their landing points

| # | Deliverable | Landing point |
|---|---|---|
| ① | two new §2.2 rows + the tier pointer | `docs/DELIVERY_DIRECTIVE.md` §2.2 (`modelId` column gains `gpt-6.1-sol` / `gpt-6-luna`) |
| ② | §2.4 blacklist convergence | same document §2.4: only `kimi-k3` + the fallback sentence + the closed-set reading |
| ③ | §5.2 ⑦ three-model tier table + cost/effort boundary note | same document §5.2 |
| ④ | §6.3 G5-b as a whitelist closed set | same document §6.3; implementation `tools/gates/lib/criteria/g5b.js` |
| ⑤ | §1.7.4 temporary-authorization five elements | same document §1.7 |
| ⑥ | criterion implementation and self-test fixtures | `g5b.js`, `tools/gates/selftest.js` (T29 to T42), 13 files under `tools/gates/testdata/fixtures/g5b-*.txt` |
| ⑦ | version bundle and appendix ledger entry | landed by this slice at ⑧b: file head version, the §0.3 row, appendix C.0j |
| ⑧ | the Luna tool-chain measurement | this slice L1 to L5 probe (readings in ADR-019 and section 6) |
| ⑨ | records | ADR-019 in `docs/ADR_REGISTER{,_EN}.md`; OB-65 to OB-69 in §12.4 |

## 3. Incidents and fallbacks

- **L1 rejected on the first attempt**: calling `runSubagent` with the `modelId` form `gpt-6-luna` made the tool answer "not found" together with its list of available models; switching to the display-name string from that list made L1 to L5 all green. The difference is registered as **OB-67** (to be discussed).
- **Unclosed finding after the ⑤ quota ran out**: ⑤ round two reported R2-1 (a stale G5-b description left in §2.5) and it was fixed in place; because the ⑤ quota (1 + 1 re-review) was exhausted no further round was bought, and by the §7.8.3 "quota exhausted" record-keeping requirement it is registered as **OB-69** (its disposition cell opens with to-do; the content is folded into the ⑦ input of this slice) and covered by this slice ⑦ re-review.
- **Red lights the freedom list itself produced**: ⑦ final review showed that list triggering G4b (a new Chinese character absent from the corpus) and G6-4 (a line-count literal shape); the wording was rewritten and the gate re-run gave `TOTAL_FAIL=0`.

## 4. Execution pipeline ① architecture (`deep-reasoner` / DeepSeek V4 Pro) ⇒ ② implementation (master + landing scripts) ⇒ ③ document
consistency (32 mechanical checks) ⇒ ④ gates (tree / selftest / index) ⇒ ⑤ pre-review, two rounds (DeepSeek
V4 Pro) ⇒ ⑥ independent scan (MAI-Code-1.1-Flash) ⇒ freedom list ⇒ ⑦ blind final review and re-review
(GPT-5.3-Codex) ⇒ ⑧ report and write-back ⇒ ⑨ native validation ⇒ ⑩ stop.
**Quota actually used**: ① 1/1, ⑤ 2/2 (1 + re-review), ⑥ 1/3, ⑦ 2/2 (final + re-review), probes 0 (this slice
declared no probe quota).

## 5. Environment and tier

Windows host; Go toolchain `go1.27.0 windows/amd64`; Node v24.17.0; repository
`github.com/larsonzh/prfrail`, branch `main`. This slice changes no Go source and touches none of
`internal/**`, `CONTRACTS`, `schemas`, `fixtures`, `workflows` or `go.mod`.

## 6. Review outcome summary (⑤ ⑥ ⑦)

- **⑤ round one** (DeepSeek V4 Pro): `PRE-REVIEW: FINDINGS` with 1 High (F1 version-bundle sequencing) +
  2 Medium (F2 four stale ⑦ Codex labels, F3 depth-tier binding to be confirmed) + 5 Low; every item
  disposed of (F1 citation fixed, F2 registered as OB-68, F3 moved to the ⑩ confirmation, F4/F5/F6/F7 fixed
  in place, F8 recorded as a freedom-list duty).
- **⑤ round two**: `PRE-REVIEW: FINDINGS` with one new Medium (R2-1) and three Low records; R2-1 was
  remediated (see section 3).
- **⑥ independent scan** (MAI-Code-1.1-Flash): `INDEPENDENT SCAN: FINDINGS`, but the report **carried no
  file:line-level defect** (it does not meet the §4.1 output contract), so no ⑥ fix loop was opened; its
  three substantive claims map to OB-67 (already registered) and to declared retained text (the `Kimi K3`
  configuration quote in §2.4 and the family breakdown paragraph in §2.5). Verdict text and the master
  ruling are in freedom list T-3.
- **⑦ final review (blind, no list)**: `FINAL REVIEW: FINDINGS` with 2 High (H1 the freedom-list red lights,
  H2 the `R2.2` versus OB-67 conflict) + 2 Medium (M3 the tier table versus §7.5, M4 no positive fixture for
  the new members) + 1 High (H5 the evidence chain of the previous slice report).
- **⑦ re-review**: `RE-REVIEW: PASS WITH FIXES` - H1 / M4 / H5 closed; M3 closed as "disclosed, awaiting a
  later ruling"; **H2 not closed**: ⑦ explicitly requires **Option B** (this slice does not change contract
  semantics on its own; keep `待议` and escalate to the user), and the three options (change `R2.2` / add a
  `modelId`-to-call-string mapping / keep it and state the difference in §2.2) must be put to the user at the
  ⑩ stop point; **until that ruling the semantics must not be called closed**.
- **Review input packs and blind isolation**: the ⑦ final review ran blind (no ⑤/⑥ report, no freedom list)
  because this slice touches model **identity** semantics (the allow-list / blacklist is model admission
  identity), which §7.7 makes a mandatory blind draw; the ⑦ re-review used a narrowed input (this round's
  remediation diff + the round-one findings and dispositions + the contract sentences + the read-only
  constraint).

## 7. Falsifiability (mechanical checks and mutation testing)

- **G5-b criterion**: 14 assertions (T29 to T42) cover closed-set members, a missing key, blacklist residue,
  case, quoting, collection of every key, the empty-value division of labour (anchored both ways with
  G5-a ②), ignoring anything outside the frontmatter, trimming, an inline comment, and positive fixtures for
  both members this slice added; the **mutation anchors T37 / T42** show that widening the set and dropping a
  member each turn the suite red, and assert a byte-exact restore of `g5b.js`.
- **CJK and expiring literals**: this slice tripped G4b / G6-4 twice on its own freedom list; both were
  handled by rewording, never by widening the whitelist.
- **A11 symmetric spot check**: ⑥ issued no "remove guard X ⇒ test Y goes red" claim this round, so there was
  nothing to spot-check; ⑦'s mutation-type claim is carried directly by the T37 / T42 assertions.

## 8. Gate results

| Scope | Result |
|---|---|
| `--scope=tree` | exit 0; `changed=24`; `TOTAL_FAIL=0` (G1-a, G1-b, G2, G3(a)/G3(b), G4a, G4b, G5-a, G5-b, G6-1 to G6-6, G7-a to G7-d) |
| `--selftest` | `SELFTEST: PASS 264/264` (16 more assertions than the 248 at slice start) |
| `--scope=index` | see ⑨ native validation (run after staging) |
| basic gates | `go build ./...` / `go vet ./...` / `go test ./...` in ⑨ |

## 9. Known divergences and coverage gaps

- **The `R2.2` call-string reading disagrees with the tool** (OB-67, to be discussed): this slice does not
  change contract semantics; it registers and escalates.
- **The §5.2 tier table is out of sync with the "⑦ Codex" rows in §7.5 / §8.2** (OB-68, to be discussed): a
  cost / effort boundary note was added under the table while the four original sentences stay untouched.
- **The depth-tier trigger set bound to `gpt-6.1-sol`** (F-6): new semantics introduced by this slice, put to
  the user for confirmation at the ⑩ stop point.
- **The ⑥ report does not meet the §4.1 output contract**: no file:line defect, so no fix loop was opened;
  the verdict text and the ruling are in freedom list T-3.
- **The ⑧b version bundle is not reviewed by ⑦**: the conditional requirement of ⑤ round two, stated
  explicitly at the ⑩ stop point for the user to veto.
- **Residual historical citations**: the historical `附录 C.0–C.0i` lines (the v1.13 row in §0.3, OB-12,
  OB-44) are left as they are (§11.1 does not rewrite history).

## 10. Candidate OBs (landed in §12.4 by this slice)

OB-65 (third pilot deferral, direct user ruling), OB-66 (counting basis for deferrals, a next-slice
candidate), OB-67 (the `R2.2` call-string reading versus the tool), OB-68 (the §5.2 tier table versus four
stale ⑦ Codex labels), OB-69 (unclosed finding from ⑤ round two, folded into the ⑦ input).

## 11. Open items and next steps

| Item | Nature | Where it goes |
|---|---|---|
| the `R2.2` reading | contract semantics | user ruling at the ⑩ stop point (one of three options; ⑦ requires Option B first) |
| the depth-tier binding | new-semantics confirmation | user confirmation at the ⑩ stop point |
| cost and effort definitions in §7.5 / §8.2 | rule gap | OB-68, next slice |
| the chain-config side of `model: "auto"` | contract gap | OB-53 item ⑤, next slice |
| the deferral counting basis | ledger wording | OB-66, next slice |
| whether ⑤ gets a third round | quota | user ruling at the ⑩ stop point |

## 12. ⑩ Stop report

- Executed through the ⑦ re-review; the ⑩ stop point is **awaiting authorisation to commit**: this slice does
  **not commit, push, or push to gitee**, and waits for the user authorisation in the same turn.
- State: `DIRECTIVE-MODEL-LISTS: STARTED | awaiting user ruling (H2 / F-6 / T-1 / ⑤ quota)`.
- Four items to rule on: (1) the three `R2.2` options; (2) confirmation of the depth-tier binding; (3) whether
  the version bundle landing at ⑧b without ⑦ review is accepted; (4) whether ⑤ gets a third round.
- **Write-back check** (OB-63 wording): validation report = written back (this file); `DEV_PLAN` = **not
  applicable** (this is a directive-governance slice whose authoritative ledger is appendix C, see §11.2);
  ledger = written back (OB-65 to OB-69 in §12.4 and appendix C.0j); ADR = written back (ADR-019).

## 13. Artifact list

```artifacts
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/t027-assembly-min.md	repo
docs/validation/t027-assembly-min_EN.md	repo
docs/validation/directive-model-lists.md	repo
docs/validation/directive-model-lists_EN.md	repo
docs/validation/evidence/directive-model-lists-freedom-list.md	repo
tools/gates/lib/criteria/g5b.js	repo
tools/gates/selftest.js	repo
tools/gates/testdata/fixtures/g5b-body-only.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-luna.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal-sol.txt	repo
tools/gates/testdata/fixtures/g5b-closed-set-legal.txt	repo
tools/gates/testdata/fixtures/g5b-empty-value.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-case.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-inline-comment.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-kimi.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-legacy.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-multi.txt	repo
tools/gates/testdata/fixtures/g5b-illegal-quoted.txt	repo
tools/gates/testdata/fixtures/g5b-legal-nokey.txt	repo
tools/gates/testdata/fixtures/g5b-legal-trailing-space.txt	repo
tmp/dml/design-1-v4pro.md	local-only
tmp/dml/change-set.diff	local-only
tmp/dml/contract-excerpt.md	local-only
tmp/dml/consistency-out4.txt	local-only
tmp/dml/gate-tree7.txt	local-only
tmp/dml/selftest7.txt	local-only
tmp/dml/ci-run-37073985144.txt	local-only
tmp/dml/luna-probe/target.txt	local-only
```

## 16. Ruling receipt after the ⑩ stop point (2026-10-03, factual record)

- The user issued four rulings at the ⑩ stop point: (1) the `R2.2` reading takes **(C)** (keep the `R2.2` body, add a measured-difference note in §2.2, and register OB-70); (2) **confirmed** the binding of the depth-tier trigger set to `gpt-6.1-sol`; (3) **accepted** the version bundle landing at ⑧b without ⑦ review (it is a factual record meeting the OB-11 exception standard and does not trigger a ⑦ re-run); (4) ⑤ gets **no** third round (R2-1 is closed through ⑦).
- Deferral counting basis: **accepting "the §12.4 ledger as the sole source of truth"**; the oral counts in receipts do not count, this one is recorded as the third (OB-65), and OB-66 is written back by the next slice.
- Sections §11 / §12 above keep their ⑩-stop-point wording (`awaiting ruling`); **no conclusion line is rewritten**: this record is the only addition, and ⑦ was not re-run (the quota is exhausted and the user ruled out a fourth round).
