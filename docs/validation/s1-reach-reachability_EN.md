# S1-REACH — OB-21 reachability refresh and the OB-51/52/53 registration (validation report)

> Slice `S1-REACH` (business assessment class, not directive governance) | class `[SLICE]` | size S | date 2026-09-29.
> **Status: stopped — incomplete** (the ⑩ stop point; the incomplete list is in §13; commit and push need same-turn authorization).

[中文](s1-reach-reachability.md)

## 1. Change summary (including the execution-freedom list)

- **What this slice did**: ① refreshed the **reachability verdicts** for **OB-21, T018 and S1 exit** from the post-DR-1 facts (per gate × dependency × unlock condition, see §2); ② completed the read-only **BYOK and secretStorage five-point confirmation** (see §3); ③ registered **OB-51, OB-52 and OB-53** in `docs/DELIVERY_DIRECTIVE.md` §12.4; ④ recorded the ⑦ downgrade as **ADR-015** per the OB-37 meta-rule.
- **What this slice did not do**: no Go code, schema, CONTRACTS, fixtures or workflow change; no live probe (budget 0); no criterion change; **the historical sentence of the OB-21 row is not rewritten** (directive §11.1 historical-row wording; the refreshed verdicts live in §2 of this report and in the OB-51 row of §12.4).
- **Artifacts**: this report (both languages), the freedom list, `DELIVERY_DIRECTIVE{,_EN}` (raised to v1.15: one §0.3 row plus three §12.4 rows), `ADR_REGISTER{,_EN}` (ADR-015), `t027/REMAINING_SLICES{,_EN}` (the slice-definition block; the closeout line is added at ⑧b), `DEV_PLAN{,_EN}` (a pointer on the T018 open question). **v1.15 is a by-product inside this slice, not a standalone governance slice** (the same applies to ADR-015).
- **Execution-freedom list**: `docs/validation/evidence/S1-REACH-freedom-list.md` (the fixed path of §1.7.3; items F-1 to F-8 are listed there, and this §1 stays consistent with that file).
- **One-line verdict**: **eight of the ten are unreachable or not passed** (G-A, G-C, G-D, G-E, G-F, G-H, G-I and OB-21 itself), **one is unknown** (G-G, a human gate) and **one was executed but does not stand alone** (G-B); DR-1 removed only one sub-item of OB-21 ② (the `DR-1 未修` clause) and did **not** change the AT-23 unreachability verdict.

## 2. Reachability verdicts (per gate × dependency × unlock condition)

> The verdict set is: unreachable / not passed / executed but not standing alone / unknown. `unknown` is reserved for human gates. The unlock sentence uses the fixed template "X needs Y, and evidence Z currently shows ...", where Z must be a locatable file or evidence pack.

| Gate | Criterion (judgeable as true or false) | Evidence | Verdict | Unlock condition |
|---|---|---|---|---|
| G-A AT-23 live candidate leg | the CLI can really launch a compatibility-admitted AI candidate and bill it | `internal/console/runtime.go:33` (`ExecuteNoopRun`, the noop-only path) and `:207` to `:211` (the fail-closed `RunPostflight` stub); the T026 block of `docs/DEV_PLAN.md` (10 verified / 2 unsupported implies an incompatible candidate) and its T027 block (`BLOCKED / NOT IMPLEMENTED`) | unreachable | T027 must complete the product assembly (evidence `runtime.go:207` currently shows postflight still a fail-closed stub) and a compatible candidate must exist (evidence: the T026 block still shows two unsupported items) |
| G-B AT-23 mechanism leg | isolated workspace / timeout / missing logs / unknown recovery / rescan after exit 0 are executed on a real host | the B4 block of `docs/t027/REMAINING_SLICES.md`; `docs/validation/t027-at23-e2e.md` | executed but not standing alone | no new unlock needed (already executed); this leg enters the AT-23 verdict only once G-A holds; it must never be written as "the mechanism leg passed, so AT-23 passed" |
| G-C AT-23 as a whole | G-A and G-B both hold and real billing evidence exists | the AT-23 row of `docs/validation/s1-exit.md` (`BLOCKED`); the four separate matters in the B4 block | not passed | G-A must be unlocked first and the live E2E re-run; evidence `s1-exit.md` currently shows AT-23 still `BLOCKED` |
| G-D T028 and AT-24 | the explicit `supervised-black-box` mode is implemented | the T028 row of `docs/DEV_PLAN.md` (`PLANNED / NOT IMPLEMENTED`); the AT-24 row of `docs/validation/s1-exit.md` (`BLOCKED`) | unreachable | T028 must be implemented; evidence currently shows `PLANNED / NOT IMPLEMENTED`; this gate is **independent of candidate compatibility admission** (the black-box mode does not require a capability matrix, see ADR-012 at `docs/ADR_REGISTER.md:22`) |
| G-E final ZIP binding and re-verification | the final ZIP is bound to the release record and re-verified | `docs/ADR_REGISTER.md:19` (ADR-009: the final ZIP must still be bound and re-verified); the "formal install-package tutorial" row of `docs/validation/s1-exit.md` (`PARTIAL`) | not passed | the remaining T018 gates must close first and the binding and re-verification then be executed; evidence currently shows this item still missing |
| G-F EOL notification route | ADR-010 is approved by the product owner | `docs/ADR_REGISTER.md:20` (ADR-010 status `PROPOSED / NOT APPROVED`) | not passed | the product owner must explicitly approve ADR-010; evidence currently shows `PROPOSED / NOT APPROVED` |
| G-G product-owner S1 exit approval | the product owner grants approval (a human gate) | the "product owner final S1 approval" row of `docs/validation/s1-exit.md` (`BLOCKED`) | unknown | the product owner must grant the S1 exit approval; evidence currently shows it has not been granted, and this slice predicts no date |
| G-H T018 as a whole | the dependency fields are all satisfied and the four remaining gates are closed | `docs/DEV_PLAN.md:91` (the dependency field contains T027 and T028) and `:93` (status `IN_PROGRESS / BLOCKED`) | unreachable | T027 and T028 must complete; evidence currently shows neither has; note the **gate level and task level do not share a source** (see below) |
| G-I S1 exit overall | AT-01 to AT-24 all pass and T018 is complete and every RFC §14 exit holds | the conclusion line at the head of `docs/validation/s1-exit.md` (`BLOCKED`) and its RFC §14 table (core task chain `BLOCKED`, three `PARTIAL` rows, owner approval `BLOCKED`) | unreachable | G-C, G-D, G-E, G-F and G-G must all hold; evidence currently shows none of them does |
| OB-21 itself | whether its ① and ② sub-items still hold | ① the assembly gap (`runtime.go:33` and `:207` unchanged); ② the `DR-1 未修` sub-item is removed (DR-1 closed, `docs/validation/evidence/dr-1-2026-09-24/`) while "T026 incompatible" and "B2 §8.3 and §8.4" are unchanged (`docs/validation/B2-EXTERNAL-ENFORCEMENT.md` §8 items 3 and 4) | not passed (the disposition word stays 待议) | a candidate or assembly advance must occur to trigger its review timing; evidence currently shows neither has moved |

- **Gate level and task level do not share a source (important)**: gates G-D to G-G do **not** depend on candidate or assembly progress (T028 is a separate mode), yet the T018 **task-level dependency field** at `docs/DEV_PLAN.md:91` contains T027 and T028, so T018 **still cannot be ticked** once those four gates are empty. This report states the gate level and the task level **separately** and gives no single combined verdict.
- **How the DR-1 evidence may be cited**: the `available` record supports exactly one inference — on the same pin the probe channel can now drive the candidate and produce a minimal real response; it supports **no** inference of "compatible candidate", "launchable candidate" or "completed assembly" (per `docs/CONTRACTS.md:191`: "a valid configuration, an existing secret and an installed CLI are all insufficient to prove AI availability").
- **The historical sentence of the OB-21 row is not rewritten** (directive §11.1 historical-row wording; this slice handles the authorization ambiguity strictly per §1.7.2), so the refreshed verdicts live in §2 of this report and in the OB-51 row of §12.4, and that row's evidence cell still contains the sub-item that has since been removed.

## 3. The five-point read-only confirmation (the factual basis of OB-51)

1. **A `--model` flag and any `configuredAIProfile` default or fallback**: the `ai check` flag set in `internal/console/ai.go` does **not** contain `--model`; the only model source is the chain configuration `ai.profiles[].model`, resolved by `configuredAIProfile` (`ai.go:230`) through channel binding to profileId to profile with **no default and no fallback** (any missing link returns not-configured and fails closed); an empty `model` is rejected by `internal/adapters/ai_provider.go:51`.
2. **Whether `model` is required in a schema**: the `schemas/` directory contains **no schema covering the chain configuration `ai.profiles`** at all, so the question is **not applicable**; the required-ness is enforced by Go validation (`ai_provider.go:51`). In `schemas/chain.schema.json` the `model` field appears at `:116`, `:156` and `:185`, all at task or codeStep level, and the corresponding `required` arrays (`:108`, `:148`, `:177`) do **not** contain `model`, so it is optional there and unrelated to `ai.profiles`.
3. **Whether `model: "auto"` is forbidden by schema or code**: it is **not** forbidden — no equality comparison against `auto` exists in code and no schema covers it; the value is passed through verbatim to the CLI at `internal/adapters/ai_probe_copilot.go:114`, and repository tests treat it as a legal value (`internal/adapters/ai_provider_test.go:13`, `internal/adapters/ai_preflight_test.go:110`). The CLI-side routing semantics of `auto` are an out-of-repository fact and **cannot be verified in this slice (probe budget 0)**, so it is labelled as such.
4. **Existence of BYOK key entries** (existence and entry identifier names only, no content read): ① **the environment-variable channel is currently empty** — this host carries only `COPILOT_AGENT`, `COPILOT_DEBUG_NONCE` and `GH_TOKEN`, whereas T026's BYOK run read the key from the process environment variable `COPILOT_PROVIDER_API_KEY` or `DEEPSEEK_API_KEY` (`docs/validation/t026-agent-runner-contract.md:117`); ② **the Windows Credential Manager** holds the Copilot CLI host entry (the `LegacyGeneric:target=https://github.com` CLI entry, whose full identifier is not written down per directive §7.3) and holds **no entry with a `ProofRail` prefix**; ③ **the Copilot CLI's own store** `%USERPROFILE%\.copilot\` contains `config.json`, `session-store.db` plus `session-store.db-shm` and `session-store.db-wal`, `logs`, `session-state`, `installed-plugins` and `ide`, with **no content read** and no BYOK key file seen **at file-name level**; ④ **VS Code secretStorage cannot be enumerated** (no external enumeration API), so it is labelled as non-enumerable and **no conclusion is drawn**.
5. **prfrail's own secretStorage mechanism**: it **exists**. The backend is the current user's Windows Credential Manager (generic type, local-machine persistence) with the prefix `windows-credential:` (`internal/adapters/ai_secret_windows.go:15`), the write path sets the user-name field to `ProofRail`, and the target name is the reference with that prefix removed (`windowsCredentialTarget`, `internal/adapters/ai_secret_windows.go:139`); the CLI is `prfrail secret set`, `status`, `delete` (no echo, argv forbidden) and the contract is `docs/CONTRACTS.md:193`. The `secretRef` example at `docs/OPERATIONS.md:56` and the `secret set` command at `:65` are **documentation examples of that same mechanism** (example target `ProofRail/deepseek`) and not a second mechanism; **that example entry is currently absent** (the credential enumeration shows no `ProofRail` prefix entry). This matches DR-1's "BYOK attribution unproven" and strengthens it: **no credential entry enumerated in this slice has a BYOK-shaped target name** and the environment-variable channel has no BYOK key either, so on both enumerable channels **the BYOK credential is currently absent**.
- **Stop judgement**: **no** conclusion among the five points to a semantics change in CONTRACTS, schema or the authorization contract, so the "stop and escalate" clause of the slice definition was **not** triggered.

## 4. Execution pipeline

| Stage | Role | Status | Artifact |
|---|---|---|---|
| ① architecture | V4 Pro (`deep-reasoner`) | executed | the plan (ten-gate framework, five-point method, acceptance criteria, file-level change points, stop conditions) |
| ②a contract | — | not triggered (no contract change) | — |
| ②b implementation | master (documentation-form degradation) | executed | three §12.4 rows, ADR-015, the ledger definition block, the `DEV_PLAN` pointer |
| ③ tests | master (documentation consistency check) | executed | gate and mirror readings (see §10) |
| ④ master integration and gates | master | executed | see §10 |
| ⑤ pre-review | V4 Pro | see §7 | see §7 |
| ⑥ independent scan | to be skipped (§7.2.1 three conditions) | see §7 | confirmed by ⑦ |
| ⑦ final review | V4 Pro (downgrade, ADR-015) | see §7 | see §7 |
| ⑧a document draft | master | see §7 | this report |
| ⑨ native validation | master | see §10 and §11 | gate re-run and mirror verification |
| ⑧b document finalisation | master | see §14 | the closeout line and the completion record |
| ⑧c wrap-up cleanup | master | see §11 | see §11 |
| ⑩ stop point | master | stopped | this report |

## 5. Deviations and fallback record

1. **Structural deviation (F-1)**: this report inserts §2 "Reachability verdicts" and §3 "The five-point read-only confirmation" into the fixed §11.3 section sequence — those two sections are the main deliverable and read better before the pipeline; **no criterion and no conclusion-line semantics change**.
2. **Master fallbacks = 1**: the first ⑤ round reported one High plus several Medium and Low findings, so ⑤ was re-run after remediation per §5.3 (no ⑨ failure occurred).
3. **Line-number basis**: the `file:line` references come from the real reads of the ① and ② stages; line numbers in existing documents drift with other slices, so **filename and section number are primary and line numbers secondary**.

## 6. Environment and tier

- Platform: the Windows host (read-only enumeration and gates); no Linux native step (this slice validates no platform).
- Class: `[SLICE]` (documentation form, degraded per §5.2); master tier is the standard one; ①, ⑤ and ⑦ are `deepseek-v4-pro` (V4 Pro).
- **⑦ downgrade (ADR-015)**: this slice's ⑦ is carried by `deepseek-v4-pro`, **not** by the §2.2 whitelist model `gpt-5.3-codex`, so **this report states it plainly: ⑦ independence is reduced** and the conclusion is handed to ⑧c for review; the exception is **not generalised** and **must not be used as a precedent**.
- **Fourth pilot deferral (four elements)**: ruling party = the user; date = 2026-09-29; reason = a business assessment slice, not directive governance; basis = the ruling text of this turn. **A direct user ruling; must not be used as a precedent.** This slice **makes no pilot metric statement** (§12.1 wording) and does not count toward the pilot evaluation.
- **Probe budget = 0**: no real external call was made in this slice.

## 7. Review conclusion summary (⑤⑥⑦)

- ⑤ pre-review (`deepseek-v4-pro`), **round 1**: `PRE-REVIEW: FINDINGS` — one High (authorization class: the "correction note" on the OB-21 row) plus five Medium and seven Low. Remediation: **the correction note was reverted** (the authorization ambiguity is handled strictly per §1.7.2 and recorded as freedom item F-2); report §1 gained the v1.15 by-product sentence (M1); the §1 verdict count was corrected (M3); the §3 line references and the "file-name level" qualifier were fixed (M5, L1); the OB-53 provenance was added (L7); wording was unified (L4, L5).
- ⑤ pre-review, **round 2 (re-review)**: `PRE-REVIEW: PASS WITH FIXES` — H1's substance is removed, M1 to M5 are all removed, and L1 to L7 are removed or retained and registered; two new Low findings were raised (N1: the report says F-1 to F-7 while the list holds F-1 to F-8; N2: the OB-21 row carried one extra trailing space). Both Low findings are fixed literally: the list count is now F-1 to F-8, and the OB-21 row is byte-identical to HEAD again (the row has left the change set). **The ⑤ budget is exhausted at 2/2**, so ⑤ is not re-run and the two literal fixes are handed to ⑦'s re-review for confirmation (recorded in the closing bullet of this section).
- ⑥ independent scan: judged against the three §7.2.1 conditions — ① the change touches only documentation and ledger-class files; ② it contains no semantics change in `.go`, `.json`, schema, fixtures or script logic; ③ it involves no state-machine, gate-criterion, evidence-model, role or authorization-contract semantics, so **all three hold and it may be skipped**; per §7.2.1 ("skipping must be confirmed by ⑦ or the user") it **has been handed to ⑦**, which ruled it skippable.
- ⑦ final review (`deepseek-v4-pro`, downgraded, ADR-015), **round 1**: `FINAL REVIEW: PASS WITH FIXES` — all four areas are "essentially holds" or "holds", with three Medium findings (M-1 stale gate evidence on disk, M-2 a self-contradictory cost section, M-3 dangling §7 references) and four Low findings; the remediation and re-review conclusions are the closing bullet of this section, and **this layer's independence is reduced** (non-whitelist model), so its conclusion is handed to ⑧c for review.
- ⑦ final review, **round 2 (re-review)**: `FINAL REVIEW (re-review): PASS` — M-1 to M-3 and L-1 to L-4 are all removed, the two literal fixes confirmed on ⑤'s behalf (N1/N2) are removed, and there is no Medium or higher. One new Low was raised (the raw credential enumeration carries a user name and a token-shaped identifier, contradicting its sanitisation promise), so ⑧c deleted that file and removed its row from the artifacts block. This was ⑦'s final round (budget 2/2).

## 8. Review input-pack summary (including blind-isolation proof)

- Input pack (laid down by the master, paths given to the reviewers): `tmp/s1-reach/plan-01.md` (the ① plan), `tmp/s1-reach/review-input.diff` (the raw diff of `DELIVERY_DIRECTIVE{,_EN}`, `ADR_REGISTER{,_EN}`, `REMAINING_SLICES{,_EN}` and `DEV_PLAN{,_EN}`), `tmp/s1-reach/s1-reach-reachability.md` and `_EN.md` (this report in full, the review object; byte-identical to the `docs/validation/` files of the same name when landed at ⑧a), `tmp/s1-reach/gate-tree.txt` (the raw gate output) and the **sanitised conclusion** of the credential enumeration (the raw output carries a user name and is withheld per directive §7.3).
- **Sanitisation**: everything handed to the reviewers was filtered for user names, host names, IPs, credentials and tokens (directive §7.3); the credential enumeration contributes only the **two relevant facts** (the CLI host entry exists; no `ProofRail` prefix entry exists) and the unrelated entries are **withheld**.
- **Blind isolation**: ⑦'s first round gets **no** ⑤ or ⑥ list (only the slice definition, the diff, this report and the read-only constraint); ⑥ is to be skipped, so there is no ⑥ list to attach.
- This slice **touches no** ownership, shutdown or identity semantics, so the §7.7 mandatory blind review is **not triggered**; under the "one in four slices" rule this slice is not at a sampling position (sampling records are decided by later slices per §7.7).

## 9. Falsifiability (mechanical checks plus mutation tests)

- Mechanical checks: see §10 (G1 to G7 plus the four base gates of §6.1).
- **Documentation mutation tests** (each mutation is applied, the criterion is confirmed red, then the bytes are restored and the hash re-checked): ① delete the CN-side OB-52 row in §12.4 so that G3(a) symmetry should go red; ② change OB-51's first status word to one outside the closed set so G6-6 should go red; ③ delete this report's `artifacts` block so G7-a should go red; ④ insert a line-count or file-count literal into the v1.15 §0.3 row so G6-4 should go red; ⑤ set the header version back to v1.14 so G1-b should go red. The actual readings are in §10.

## 10. Gate results

- `node tools/gates/gate.js --all --scope=tree` ⇒ `TOTAL_FAIL=0` (`gate-exit=0`, `changed=9`; the report is not in the repository yet, so `G7-a` to `G7-d` read `INFO` and are re-run after landing).
- Per criterion: `G1-a` PASS; `G1-b` PASS (`header=v1.15 last=v1.15 2026-09-29/2026-09-29`); `G1` SKIP (not scripted); `G2` PASS; `G3(a)` PASS (four pairs symmetric: directive `+5/-1`, ADR `+1/-0`, `DEV_PLAN` `+1/-1`, ledger `+19/-0`); `G3(b)` PASS; `G4a` PASS; `G4b` PASS (`unknown=[]`); `G5-a` and `G5-b` PASS; `G5` SKIP (not scripted); `R2.5` PASS; `G6-1` to `G6-4` PASS; `G6-5` PASS (`scanned 2, refs 6`); `G6-6` PASS (`scanned 8`).
- Base gates: `gofmt -l` empty, `go build ./...` exit 0, `go vet ./...` exit 0, `go test ./...` `ok=14 FAIL=0`.
- Mirrors: the report pair and the directive pair both have zero per-line mismatches (the `_EN` files stay positionally aligned with CN); `ADR_REGISTER{,_EN}` and `REMAINING_SLICES{,_EN}` are **not whole-file positional mirrors** (they carry pre-existing positional divergence), so only the `G3(a)` changed-line symmetry applies to them.
- The ⑧c readings (including the `G7` re-run after the report lands) are in §11.

## 11. ⑧c wrap-up cleanup record

- Whole-repository per-file encoding check (`enc-tree.js`): **3** violations, **all** of them OB-27 frozen-evidence-pack exemptions (the `.raw.txt` CRLF under `b2-2026-09-17` and the two BOMs under `b3b-2026-09-21`); **this slice's new files have zero violations**.
- Gate self-test: `SELFTEST: PASS 248/248`; `gate.js --all --scope=tree` ⇒ `TOTAL_FAIL=0` (after the report landed, `G7-a` to `G7-d` are all PASS; `G7-b` covers 28 artifact rows and `G4a` prints the file count it checked itself).
- Temporary-residue inventory: `tmp/s1-reach/` retains 13 process files (the plan, drafts, assertion scripts, the raw gate output and the review-pack copies), **not cleaned** per OB-49 and the standing user instruction, with disposition `local-only`; no untracked non-temporary residue was found in the repository.
- ⑧c found and disposed of exactly one item: ⑦'s re-review Low (the raw credential enumeration carries a user name and a token-shaped identifier) ⇒ that file was deleted and its row removed from the artifacts block. This section and the ⑧b write-backs are **factual records** (tool readings, review conclusions, ledger write-back), so per the §5.3 semantic boundary they **do not trigger a re-run of ⑦**; if the user considers that they should, this slice's ⑦ budget (2/2) is exhausted and the matter must be escalated to a user ruling.

## 12. Cost and metering

| Stage | Role | Model | Calls | reason / retryOf |
| --- | --- | --- | --- | --- |
| ① | architect | `deepseek-v4-pro` | 1 | required by the slice definition (verdict framework, five-point method, change points) |
| ⑤ | pre-reviewer | `deepseek-v4-pro` | 2 | the second call is the post-remediation re-review (retryOf = the first) |
| ⑥ | independent scan | none | 0 | skippable per §7.2.1 and confirmed by ⑦ |
| ⑦ | final reviewer | `deepseek-v4-pro` | 2 | the final review (downgraded, ADR-015) and the post-remediation re-review |
- Probes: **0** (budget 0).
- Truth source for metering: the session call receipts (§7.5 hard rule R7.1); nothing is estimated from memory.

## 13. Explicitly not executed

- commit and push (need same-turn authorization); the live probe (budget 0 in this slice); a real BYOK test (if a verdict needs one, stop and escalate rather than doing it here); the ⑥ independent scan (to be skipped, pending ⑦ confirmation); any pilot metric statement (this slice does not count toward the pilot evaluation).

## 14. Next-step recommendations

- **This slice's verdicts unlock nothing**: the AT-23 live candidate leg stays unreachable and T018 and S1 exit stay blocked.
- Suggested unblocking order (a separate slice, after a user ruling): **the T027 assembly first** (the noop-only path in `runtime.go` and the fail-closed postflight stub are hard prerequisites), then a candidate-compatibility reassessment; T028 is **independent of candidate compatibility** and can be assessed in parallel with the candidate line.
- Tier advice: a T027 assembly slice touches the state machine and crash paths, so it should keep `gpt-5.3-codex`; business-assessment slices like this one can keep V4 Pro.
- Handing over open items: "VS Code secretStorage cannot be enumerated" in OB-51 is a **shared boundary** of this host and external tooling that no later slice can cross, so it should stay labelled as such long term.

```artifacts
docs/validation/s1-reach-reachability.md	repo
docs/validation/s1-reach-reachability_EN.md	repo
docs/validation/evidence/S1-REACH-freedom-list.md	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
tmp/s1-reach/plan-01.md	local-only
tmp/s1-reach/review-input.diff	local-only
tmp/s1-reach/gate-tree.txt	local-only
```
