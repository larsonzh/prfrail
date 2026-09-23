# DR-1 validation report: the availability probe's CLI credit floor

> Slice `DR-1` (business slice; includes one directive-ledger registration) | class `[SLICE]` | size S | date 2026-09-24.
> Status: **stopped - incomplete** (⑩ stop point; commit / push need same-turn authorization).

## 1. Change summary (including the execution-freedom list)

The value `prfrail ai check` passes to the candidate CLI as `--max-ai-credits` changed from "the product request budget" to `max(30, maximumRequests)`, so the product availability probe can produce an `available` record on the pinned candidate for the first time.

- Implementation: `internal/adapters/ai_probe_copilot.go` gains the constant `copilotCLICreditFloor = 30` (`:23`), and the argument site moved from `:111` to `:117`; `internal/console/ai.go` is **untouched**.
- Tests: the argv assertion changed from `"1"` to `"30"`; two new cases (three sub-cases plus a CLI-rejection fail-closed case).
- Documentation and ledger: the v1.14 ledger registration (OB-46-50), ADR-014, the DR-1 slice definition, the chain config, the evidence pack and its checksum list.
- Freedom list: `docs/validation/evidence/DR-1-freedom-list.md` (F-1 to F-6 plus the user-ruling items D-1 to D-6).

## 2. Execution pipeline

| Stage | Role / model | Outcome | Quota |
|---|---|---|---|
| ① architecture | **no standalone plan** (baseline = user rulings R1-R5) | see §3 item 1 | - |
| ②a contract-first | not triggered (R1=A froze record fields and schema) | - | - |
| ②b code | master | constant plus argument site | - |
| ③ tests | master | two cases plus the mutation check | - |
| ④ gates | master | see §9 | - |
| ⑤ pre-review / re-review | `deepseek-v4-pro` (Max) | `PASS WITH FIXES` → remediation → `PASS` | 2/2 |
| ⑥ independent scan | `mai-code-1.1-flash` (blind) | `INDEPENDENT SCAN: FINDINGS` | 1/3 |
| ⑦ final review | `deepseek-v4-pro` (**downgraded**, ADR-014) | `FINAL REVIEW: PASS` | 1/1 |
| ⑨ native validation | master | one live probe ⇒ `available` | 1/5 |

## 3. Anomalies and fallbacks (item by item)

1. **No standalone ① plan** (a deviation from §1.2): the user rulings R1-R5 are the plan baseline, recorded in `tmp/dr1/rulings.md` and as freedom-list item D-1.
2. **Execution-order deviation** (F-6): the live probe ⑨ ran first and this report was written afterwards, so that the report never contains a placeholder form that G2 forbids; no rule semantics changed.
3. ⑤ Medium① (mutation evidence not in the input pack): closed after landing `tmp/dr1/mutation-check.txt`.
4. ⑤ Low② (dangling §12.4 reference): re-pointed to ADR-014 in `docs/ADR_REGISTER{,_EN}.md` per the user's ruling, with the ledger pointer corrected.
5. The ⑥ layer did not declare a pass: per `docs/DELIVERY_DIRECTIVE.md` §7.2 the alternate path routes closure to ⑦, closed by ⑦'s verdict (user ruling, option b).
6. Two ⑦ Lows: the §7.2 branch has no literal counterpart (recorded), and the `DEV_PLAN{,_EN}` write-back plus the `:111 → :117` reference drift (see §13).
7. **⑧b write-back script defect (self-detected)**: the first script read the same ledger file twice and the second write overwrote the first ⇒ the defect-bullet rewrite and the closeout line were lost; it surfaced when `git diff --numstat` showed the ledger at `+17/-0` only, and a patch script asserting "old text present, new text absent" re-applied them, verified as `+20/-1`; it affected no code, criterion or conclusion.

## 4. Environment and tiers

- Platform: local Windows, Go 1.22 toolchain; pin = `copilot.exe` 1.0.83 (sha256 `d3f3bb7b…0f671ee2`, re-verified byte-identical before the probe).
- Tiers: master at standard High; ⑤ and ⑦ = `deepseek-v4-pro` (Max, not charged to the Copilot budget); ⑥ = `mai-code-1.1-flash`.
- **⑦ downgrade**: this slice's ⑦ ran on a non-allowlist model, so **independence is reduced**; recorded independently as ADR-014 per the OB-37 meta-rule (four elements plus "not a precedent" plus the three annotations), handed to ⑧c for review.
- **Pilot deferral number 3**: this slice validates a business path, not directive governance ⇒ it is not a pilot; four elements: ruling party = the user; date = 2026-09-24; reason = business-path validation; basis = the user's ruling text from the same round; **not a precedent**.

## 5. Review outcome summary

- ⑤ first round `PASS WITH FIXES` (0 Critical / 0 High / 1 Medium / 5 Low); re-review `PASS` (all six items closed, no new deviation).
- ⑥ `INDEPENDENT SCAN: FINDINGS` (High ×1 plus Medium ×1; theme = the CLI effective allowance differs from the product budget and the record carries no cost dimension).
- ⑦ `FINAL REVIEW: PASS` (all four areas pass; 2 Low plus 2 Note, all downstream obligations).
- Quota: ⑤ 2/2, ⑥ 1/3, ⑦ 1/1, probe 1/5 (four returned).

## 6. §7.2 disposition record (the alternate path for the ⑥ FINDINGS)

- Both ⑥ findings point at one difference: the CLI effective allowance of 30 is not the product request budget of 1, and the record's `requestsUsed` is only a request count with no cost dimension; their "minimal fix" (write the effective allowance into the record, or require `unavailable` in that situation) directly conflicts with **R1=A**.
- Disposition (user ruling, option b): **keep R1=A, do not fix**; the ⑥ layer **does not declare a pass**; per §7.2 the alternate path routes closure to ⑦, closed by `FINAL REVIEW: PASS`.
- Positive signal: ⑤ (with the rulings) and ⑥ (blind) flagged **the same difference independently** ⇒ the difference is real, visible and not invented by us.

## 7. Review input pack summary (including blind isolation)

- ⑤ and ⑦ inputs: the slice definition, the R1-R5 rulings digest, the raw diff (`tmp/dr1/review-input-r2.diff`, 9 files), raw gate and focused-test output, the mutation evidence, the chain-config self-proof and the verbatim contract sentences; **excluding** the ① conversation history and the ②b implementation reasoning.
- ⑥ first round: **only** the diff, the contract sentences and the read-only constraint, **with no list from anyone else**.
- All three layers carried the hard constraint "create, modify or delete no file; make no claim of running `node` / `go` / `gate` commands, citing only paths the master wrote to disk".

## 8. Falsifiability (mechanical checks plus the mutation check)

- Mutation: an assertion-based script required exactly one matching site (otherwise it throws) ⇒ after mutation **three tests failed** (`ai_probe_copilot_test.go:80`, `:318`, `:342`) ⇒ byte-identical restore ⇒ the sha256 matched before and after (`B805EA42…A4EB7`) ⇒ green again.
- Uncovered invariants (registered, not merged into this slice): the existing `invalid FailureCode → unknown` and `maximumRequests < 1 → error` branches have no direct test case.

## 9. Gate results (final value before ⑩)

- `gofmt -l` empty, `go build ./...` clean, `go vet ./...` exit 0.
- focused `go test -count=1 -run CopilotCLIProbe ./internal/adapters/` ⇒ ok; full `go test -count=1 ./...` ⇒ 14 packages ok, FAIL 0.
- `node tools/gates/gate.js --all --scope=tree` ⇒ `TOTAL_FAIL=0`; G4a encoding criterion green (its own output carries the file count); G6-5 scanned 2 / refs 4 / unresolved 0; G6-6 scanned 10; G3(a) four pairs (directive `+7/-1`, ADR `+1/-0`, DEV_PLAN `+1/-1`, ledger `+20/-1`); G1-b header and log both v1.14.
- **Mirror scope (stated as measured)**: `DELIVERY_DIRECTIVE{,_EN}` is a **strict positional mirror** (1309 / 1309, zero per-line mismatches); `ADR_REGISTER{,_EN}` and `REMAINING_SLICES{,_EN}` are **not** whole-file positional mirrors (HEAD already carries 37 / 226 positional differences, not introduced here), so they are held to G3(a)'s "symmetric changed-line counts" plus block-level shape equality (the two ADR-014 rows are 14 bold / 5 bars on both sides; the DR-1 definition block has equal line shapes throughout).
- **Not executed**: `gate.js --selftest` did not run in this slice (no criterion was changed ⇒ no trigger), and will run before the commit and in CI.

## 10. ⑧c wrap-up cleanliness record

- Whole-repository per-file encoding: `tracked=2923 text=2834 binary=89 violations=3`, all three inside **frozen evidence packs** (the OB-27 temporary exemption: one CRLF file in `b2-2026-09-17` and two BOM files in `b3b-2026-09-21`); zero violations among this slice's new files.
- Evidence-pack self-consistency: `packs=5 entries=2331 missing=0 mismatch=0 not-self-checkable=1`; this slice's pack moved from not-self-checkable to `entries=4 missing=0 mismatch=0 -> PASS`.
- `tmp/` was **not cleaned**: per the user's standing instruction (it belongs to `DIRECTIVE-LEDGER-ARCHIVE` and OB-49), this slice's drafts and raw captures are kept and declared `local-only`.

## 11. Known difference and reserved candidate OBs

- **Known difference (R1=A, not fixed)**: between the contract sentence at `docs/CONTRACTS.md:191` ("within an explicit cost/request budget") and the CLI's actual allowance of 30 there is a cost-dimension difference, and the record does not carry that dimension; annotated in three places - the code comment (`:19-22`), the slice-definition boundary row and this report.
- **Candidate OB-52** (to be registered next slice): that difference (merging ⑥'s High and Medium) **combined with** ⑥'s Low④ (a CLI floor change would collapse to `unknown` with only hashes kept, so it has no mechanical discoverability) - two sides of one difference.
- **Reserved candidate OB-51**: the read-only confirmation set for BYOK and secretStorage, deferred until after DR-1 closes; if a future `DIRECTIVE-MODEL-LISTS` slice needs a number, it takes **OB-53**.

## 12. Cost and metering

- Paid model calls total **4** (⑤ 2 + ⑥ 1 + ⑦ 1); of these ⑥ is charged to the Copilot budget.
- Probe **1/5** (target 1, four returned); a successful probe is booked as **one Copilot premium** under the worst case; **"BYOK attribution unproven" is stated as such** (this slice did not use BYOK).
- No ① call (no standalone plan); ⑧ and ⑨ ran in the master, zero paid calls.

## 13. Not-executed items and follow-ups

- Not executed: commit and push (need same-turn authorization); `gate.js --selftest` (before the commit and in CI); ⑥'s verbatim rows were not given to ⑦ (deliberate isolation); the BYOK and secretStorage read-only confirmation (after wrap-up).
- Follow-ups: the `DEV_PLAN{,_EN}` B1-paragraph write-back and the `:111 → :117` reference fix landed in this slice; the next slice registers OB-52 and OB-51 (if the BYOK item is still open).

## 14. Artifact list

```artifacts
internal/adapters/ai_probe_copilot.go	repo
internal/adapters/ai_probe_copilot_test.go	repo
docs/DELIVERY_DIRECTIVE.md	repo
docs/DELIVERY_DIRECTIVE_EN.md	repo
docs/ADR_REGISTER.md	repo
docs/ADR_REGISTER_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/DEV_PLAN.md	repo
docs/DEV_PLAN_EN.md	repo
docs/validation/dr-1-availability-probe-floor.md	repo
docs/validation/dr-1-availability-probe-floor_EN.md	repo
docs/validation/evidence/DR-1-freedom-list.md	repo
docs/validation/evidence/dr-1-2026-09-24/proofrail.chain.json	repo
docs/validation/evidence/dr-1-2026-09-24/availability-1083-available.json	repo
docs/validation/evidence/dr-1-2026-09-24/probe-01-ai-check.json	repo
docs/validation/evidence/dr-1-2026-09-24/probe-01-ai-verify.json	repo
docs/validation/evidence/dr-1-2026-09-24/README.md	repo
docs/validation/evidence/dr-1-2026-09-24/SHA256SUMS.txt	repo
tmp/dr1/**	local-only
```

**Disposition of `tmp/dr1/**`**: it holds this slice's review input pack, the mutation transcript, the chain-config self-proof, the raw probe captures and the drafts; per the user's standing instruction it is **not cleaned** and is left to `DIRECTIVE-LEDGER-ARCHIVE` (dual theme: ledger archival and the `tmp/` lifecycle).

## 15. Appendix A: ⑨ native-validation record and comparison

| Item | B1 pre-fix baseline | After this slice's fix |
|---|---|---|
| Record file | `docs/validation/evidence/b1-20260917/availability-1083-unknown.json` | `docs/validation/evidence/dr-1-2026-09-24/availability-1083-available.json` |
| `status` | `unknown` | `available` |
| `reason` | `unknown` | absent (`available` forbids `reason`) |
| `requestsUsed` | 0 | 1 |
| Exit code | 1 | 0 |
| `profileConfigHash` | `sha256:d529c357…e436` | the same value (byte-identical) |
| `profileId` / `channel` | `copilot-cli-host` / `agent-runner-cli` | the same values |

- Additional probe evidence: `ai verify` against the same record returned `ok=true` and `exitCode=0` (freshness and budget checks pass, zero model calls).
- The record is 745 bytes with SHA-256 `2DC3EDBD…4418`; the pin binary re-hashed to `D3F3BB7B…1EE2`.

## 16. Appendix B: ⑤ / ⑥ / ⑦ verdict lines and ⑥'s two findings verbatim

- ⑤ first round: `PRE-REVIEW: PASS WITH FIXES`; re-review: `PRE-REVIEW: PASS`.
- ⑥: `INDEPENDENT SCAN: FINDINGS`.
- ⑦: `FINAL REVIEW: PASS`.
- ⑥'s findings verbatim (not rewritten):

`High | internal/adapters/ai_probe_copilot.go:19-23 | 把候选 CLI 的平台最小 credits（30）强制覆盖原始 maximumRequests，但没有把 effectiveProbeCredits 与 declared budget 一并写进 record；available 可在 maximumRequests=1 时出现，违反"显式费用/请求预算内"的契约。 | 把 effectiveProbeCredits 与 declared budget 一并写进 record`

`Medium | internal/adapters/ai_probe_copilot_test.go:292-342 | 测试只断言 RequestsUsed == 1 和 Status == "available"，没有模拟真实 credit floor 与 declared budget 的冲突，也没有验证 effectiveProbeCredits <= budget | 增加一条"maximumRequests=1 且真实 effective credits=30 时必须返回 unavailable"的反例`
