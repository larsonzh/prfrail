# T027 · B3c · CONTRACTS durability-upgrade clause revision — validation report

Date: 2026-09-21. Status: `COMPLETE` (① V4 Pro pre-analysis / ③ V4 Pro pre-review / ③.5 independent scan / ④ Codex final review all in §9; **no commit and no push yet**, awaiting same-turn authorisation).

## 0. Honesty rules (the premise of every conclusion here)

1. This slice **only freezes the contract wording**. It does not lift the Windows `unproven` state, does not relax `first-dispatch`, does not change T027's `BLOCKED / NOT IMPLEMENTED`, and is not AT-23 evidence. Every `proven` statement must trace back to "which platform + which premise + which evidence".
2. This slice contains **no `.go` / `.json` / schema / fixtures change**; the gates are accepted on an "unchanged means green" basis.
3. This slice does **not re-run** the B3a/B3b physical rig and does not touch `tools/agent-probe/enforcement-proxy` (DR-2 is its own slice).
4. Every cited field (file, line, count, mutation result) comes from actual command output and this slice's review records; anything not executed is listed in §10.

## 1. Scope and boundaries

**Goal**: freeze B3b's **tiered** verdict (C2 `proven`; C1/C3/C2C3/U4 `unresolved`) into a **decidable contract wording**, so that "a platform's publication durability rises from `unproven` to `proven`" becomes a normative path with a premise, criteria, and boundaries, giving B4 a stable base.

**Gap closed here**: the bilingual CONTRACTS stated "publication durability is unproven" as a **platform constant** (on Windows, "new directory-entry persistence cannot currently be forced or proved" ⇒ `unproven`) and defined neither the upgrade condition nor what an upgrade does and does not unlock.

**Boundaries (respected)**:
- **This slice only freezes the wording and does not lift Windows `unproven`**; the `first-dispatch` refusal, the AT-23 conclusions and every other platform's conclusion are **unchanged**.
- No Go change, no schema/fixtures change, no rig re-run, no touching the B2-era tool.
- DR-2 (the CI flake) is not handled across slices; only a third piece of run evidence is appended with this documentation change.

## 2. Change list (8 `.md` files = the whole working-tree diff)

| # | File | Nature of the change |
|---|---|---|
| 1 | `docs/CONTRACTS.md` | new `**发布耐久等级与升级条件**` paragraph in §7 (L209); the three frozen sentences parameterised (L207 / L219 / L225) |
| 2 | `docs/CONTRACTS_EN.md` | English mirror of §7 (new paragraph L203; sentences 1/2/3 = L201 / L213 / L219) |
| 3 | `docs/ADR_REGISTER.md` | ADR-013 row completion condition and the §5 adoption precondition now point at the new clause (B3c) |
| 4 | `docs/ADR_REGISTER_EN.md` | same, English mirror |
| 5 | `docs/t027/REMAINING_SLICES.md` | new B3c slice definition (scope, goal, gap, deliverables, acceptance evidence, boundaries), status table, dependency graph (`B3b → B3c → B4`), 2026-09-21 revision note; B4 dependency gains `B3c` |
| 6 | `docs/t027/REMAINING_SLICES_EN.md` | same, English mirror |
| 7 | `docs/validation/t027-b3b-durability.md` | §7 CI-flake bullet gains a third run |
| 8 | `docs/validation/t027-b3b-durability_EN.md` | same, English mirror |

**CI back-fill (folded into this slice's documentation change, no extra commit)**: the docs back-fill commit `dd36b9b` produced CI run **`35564805991`** = **success**, with `Go ubuntu-latest` and `Go windows-latest` **both green**; the previously intermittent `tools/agent-probe/enforcement-proxy::TestConnectAuditTrailFailureDeniesAndDoesNotTunnel` did **not** recur, consistent with DR-2's verdict ("intermittent, pointing at a handle/timing race inside that test").

## 3. The new clause (`docs/CONTRACTS.md` §7 / `docs/CONTRACTS_EN.md` §7)

1. **Three-state**: `proven` / `disproven` / `unresolved` (`observed`/`falsified` are evidence annotations, never tiers), as in A7.
2. **Unit of judgment = (platform × premise P)**: every tier statement must carry its platform, premise and evidence; a tier earned under a premise must never be restated as a premise-free platform-level tier (a platform-level tier may only state the `unproven` default); every `proven` must trace back to platform + premise + evidence.
3. **Three criteria for recording `proven` (all required under the same premise)**:
   (a) reproducible injection or control evidence (valid rounds inside the allowed set, an evidence pack recomputable from the round records and the journal, no true orphans);
   (b) a **discrimination gate** (under the same rig, the same writer and the same record-level payload, negative controls show at least one loss and positive controls none; **a failed gate caps the verdict at `unresolved`**);
   (c) **mechanism evidence** (the claimed persistence mechanism demonstrably ran in that round).
   **Missing any one of them leaves the answer `unresolved`** (and `unresolved` means only that no claim was established — it must **never** be read as "not durable" or "already falsified"); an attributable counterexample (a loss where the stage was hit and the mechanism demonstrably ran) is `disproven`.
4. **Forbidden wording**: "durable forever", "long-term guarantee", "proven power-loss safe", premise-spanning universal durability; a `proven` may only be phrased as "no counterexample was observed under this premise at this device's resolution, with the mechanism demonstrably run and the device demonstrably discriminating".
5. **Object separation**: a candidate mechanism's tier ≠ the **platform publication-durability tier**; the platform tier is decided only by the evidence for the mechanism the platform's **production implementation actually adopts**, under a premise covering the target deployment environment — closing the short-circuit reading "candidate `proven` ⇒ platform `proven`".
6. **Runtime mapping**: `unresolved` and `disproven` behave as `unproven` (fail-closed behaviour unchanged).
7. **Upgrade unlock boundary**: after a (platform × premise) reaches `proven`, only that platform, within the deployment environments covered by that premise, has the corresponding refusals lifted; **nothing** is unlocked for any other premise or platform, **AT-23 is not waived** (it must still pass), **independent review and protocol-first revision are not waived**, and no gate beyond `first-dispatch` is relaxed. A platform-tier upgrade requires a protocol-first amendment of this clause, traceable evidence, and an independent review.
8. **Existing instance**: on a local Unix filesystem, "a new link was created and its parent directory was `fsync`-ed successfully" counts under the existing contract wording as an instance satisfying every criterion (its premise is the local-Unix-filesystem class).

## 4. Parameterisation of the three frozen sentences

| Sentence | Location | Change |
|---|---|---|
| 1 | CN L207 / EN L201 (replay-store paragraph) | "new directory-entry persistence **cannot currently be forced or proved** ⇒ `unproven`" → "the tier of new directory-entry persistence is **not a platform constant**: determined per platform and per premise under this section's clause… (`unproven` while not every criterion is met)"; plus "per the B3b verdict ledger (`docs/validation/t027-b3b-durability.md` §4) the Windows verdict is **tiered** (C2 `proven`, the others `unresolved`) and not a platform-wide upgrade, so the Windows publication-durability tier **remains** `unproven`" |
| 2 | CN L219 / EN L213 (launch-receipt paragraph) | "on platforms where publication durability is unproven…" → "on a platform whose publication-durability tier is `unproven` (determined per platform and per premise; see this section's clause)…"; plus "a tier upgrade changes the value of this rule only **for that platform under that premise**" |
| 3 | CN L225 / EN L219 (terminal publication and settlement paragraph) | same parameterisation plus the same upgrade-boundary clause |

**Refusal strength unchanged**: sentence 2 still rejects pre-write, never spawns, and never returns `first-launch`; sentence 3 still writes no intent, settles nothing, publishes no C and writes no closure, and only an already-complete chain may be replayed write-free.

## 5. Falsifiability (mechanical checks + mutation experiment)

Temporary script `tmp/b3c-checks.js` (gitignored, deleted after use):

- **Baseline: 34 mechanical checks, all GREEN**, covering: the three-state sentence, the unit of judgment, the tightened discrimination gate, mechanism evidence, the `unresolved` anti-misreading anchor, the upgrade-boundary clause (2 places), "AT-23 not waived", the runtime mapping, protocol-first platform upgrades, the tiered anchor, "Windows stays `unproven`", clause-title reference counts (CN 4 / EN 4), the verdict-ledger citation, **zero residual old platform-constant wording**, zero unconditional-unlock wording, zero character-level defects, B4 depending on B3c (both languages), the dependency-graph edge, zero stale line anchors, and the third CI run recorded (one per language).
- **Mutation experiment, 4 cases (A12 dedicated independent corroboration; each "break it and watch it go red")**:
  | Mutation | Check expected to go red | Observed |
  |---|---|---|
  | M1 delete "不是整体升级" from sentence 1 | C19 tiered anchor | **RED** (no collateral failures) |
  | M2 weaken the gate back to "负对照可丢失" | C05 tightened discrimination gate | **RED** |
  | M3 delete the per-platform scoping from EN sentence 2 | C12 upgrade-boundary count | **RED** |
  | M4 write an unconditional upgrade/unlock | C28 unconditional unlock = 0 | **RED** (with C21 "Windows stays `unproven`") |
- **Restore integrity**: after all four mutations every file is **byte-identical** to its pre-mutation state (SHA-256 unchanged), and the 34 checks are ALL GREEN again.
- **Coverage limits of the mechanical checks**: the three items listed as "not mechanically verifiable" — "the evidence pack is recomputable from rounds and the journal" (needs the pack to be re-run), "the premise covers the target deployment environment" (needs a deployment baseline), "independent review is sufficient" (needs the review-process record) — are **not** claimed as verified by this slice.

## 6. Gates

| Item | Result |
|---|---|
| `gofmt -l .` | empty output |
| `go build ./...` | exit 0 |
| `go vet ./...` | exit 0 |
| `go test ./...` | exit 0 (`ok` packages **14**, `FAIL` count **0**) |
| Encoding / EOL (8 `.md` files) | UTF-8 **with BOM** + **LF**: **PASS** (per-file byte check) |
| Contract fixtures | unaffected (no schema/fixtures change) |

## 7. Known boundaries and follow-ups

- **This slice is not durability evidence**: it only states what an upgrade requires; it produces no `proven` of its own.
- **Windows stays `unproven`, `first-dispatch` stays refused, T027 stays `BLOCKED / NOT IMPLEMENTED`, AT-23 has not passed** — all four are **unchanged**.
- **DR-2 (the `enforcement-proxy` intermittent flake) remains open**: this slice only appends a third piece of run evidence (no recurrence) and **does not fix a B2-era tool across slices**; the disposition stays "its own slice before B4".
- **Runtime behaviour of `unresolved` and `disproven`**: both equal `unproven` (stated in the clause); distinguishing them later would be a new slice.
- **Residual**: the ③.5 independent scan produced **no actionable finding in either of its two calls** (both deviated in format); it is recorded honestly in §9.2 and re-runs were stopped at the §3.10 cap. Its "no Medium+" outcome **is not proof of review sufficiency**.
- **Observation item (2026-09-21 user ruling: register it but do not change the directive)**: **the effectiveness of ③.5 on documentation-only slices is to be assessed**. Facts: MAI output is normal on code slices (A6), while on this slice (pure `.md`) both calls deviated from the format contract and produced no actionable finding; the sample size is n=1, **too small to amend the directive**. Handling: use the next slice, **B4 (a code slice)**, as the control - normal output on B4 attributes this slice's deviation to "documentation-type work"; the same deviation on B4 attributes it to the model. If another 2-3 documentation-only slices behave the same way, then consider adding a §3.10 exception ("a contract revision that only edits `.md` and touches no executable semantics may skip ③.5 with ④ confirming"). **The directive is not changed by this slice.**
- **Cleanup confirmation (2026-09-21, same-turn user authorisation)**: the **11** leftover items under `tmp/b3b/` were removed (8 rig console logs + 3 writer build products, **26,883,260 B** in total). **They were first confirmed unrelated to the evidence pack**: `docs/validation/evidence/b3b-2026-09-21/SHA256SUMS.txt` (2157 entries) contains **0** `.exe` entries, and none of the 11 filenames matches anything in the pack (the pack's own logs are `session/session.log` and `variant/session.log`, different paths and names); `writer.exe`'s identity is already on record via this report's §3.2 sha256 `06de8a72…` (measured `06de8a728ae984318a7955d01853443f30749a621c81d4caa47044c5b7017a82`) and the file is a build product reproducible with `go test -c -tags b3bnative`. ⇒ judged **unrelated to the evidence pack** and **deleted directly** as instructed (sha256 recorded first: `writer.exe` `06de8a72…`, `writer-variant.exe` `5b0d06d8…`, `writer-rebuild-check.exe` `880655ce…`; the eight logs were hashed individually in this round's terminal output).

## 8. Cost and accounting

- **Probe budget: 0/10 used** (fully offline: no paid probes, no real model calls).
- **① V4 Pro**: one pre-analysis (rulings U1–U14 all adopted as recommended).
- **③ V4 Pro**: one pre-review (`PASS WITH FIXES`: 1 Medium + 1 Low mandatory, 3 Low advisory) → fixes → one re-review (`RE-REVIEW: PASS`, 4 Low) → one delta re-review (`DELTA-REVIEW: PASS`, 1 Low).
- **③.5 MAI-Code-1.1-Flash**: 2 calls (first + one anomaly resend), neither producing an actionable finding (see §9.2); within the §3.10 cap of ≤3.
- **④ GPT-5.3 Codex**: 3 calls — first final review (no lists attached, effectively a blind audit) `FINDINGS` (1 Medium + 1 Low) → fixes → re-review `FINDINGS` (the only blocker being "no raw diff hunks supplied") → a closing round after the raw diff was written to disk: **`PASS`**.
- **§3.11 mandatory blind audit**: this slice's ④ first call **attached no ③/③.5 list**, which is an effective blind audit; as with B3b, a second blind audit was not added because the input sets are isomorphic and a repeat would add cost without discriminating power — the rationale is recorded here explicitly.
- **⑤ wall clock**: the mechanical checks plus four mutations take about 2 s; the gate suite about 1 min.
- **token**: filled in by the user from the usage UI (the main controller never estimates it).

## 9. Review record (① / ③ / ③.5 / ④) and closure loop

### 9.1 ① V4 Pro pre-analysis

**All of U1–U14 adopted as recommended**, notably: the clause sits as its own paragraph after the replay-store paragraph (U1); granularity is platform × premise (U2); the discrimination gate and mechanism evidence become criteria (U3); independent review is positioned as an **upgrade governance precondition**, not a tier criterion (U4); "the evidence pack is recomputable" folds into criterion (a) (U5); the `disproven` trigger rule is added (U6); the runtime mapping is stated (U7); **no schema change** (U8); the ADR keeps only the ledger and the criteria have a single authority in CONTRACTS (U10); a new validation report (U11); bilingual ledger write-back (U12); ④ keeps its final review (U13); the tiered anchor lands in two places only (U14).

### 9.2 ③ V4 Pro pre-review and the ③.5 independent scan

- **③ pre-review**: `PASS WITH FIXES`. **Medium (mandatory)**: A1 the discrimination-gate wording was weaker than the single authority in `docs/t027/B3B_DESIGN.md` §7 and omitted the "gate failed" consequence; E1/E2 the ADR-013 row completion condition and the §5 stale-sentence citation had not been written back per U10 (dangling reference). **Low**: A2 give the Unix instance its premise, A3 narrow the premise-free-tier prohibition, A4 rename and locate the source document, C1 add the `unresolved` anti-misreading anchor. All fixed.
- **③ re-review**: `RE-REVIEW: PASS` (A1–A4 / C1 / E1 / E2 all closed; no new Medium+; 4 Low: writer identity, a Chinese typo, EN bolding and label alignment).
- **③ delta re-review**: `DELTA-REVIEW: PASS` (all three deltas landed; the single Low was a bolding asymmetry, since fixed).
- **③.5 MAI independent scan**: **neither of its 2 calls produced an actionable finding**. The first deviated in structure (invented sections, no `file:line` findings, no verdict line); per §3.10 ("a resend is allowed only for truncation or anomaly, and counts against the cap") it was resent once, and the resend still deviated (still no Section A–E and no verdict line, though it did read the documents and quote real text, corroborating "nothing unlock-related found"). **Main-controller handling**: re-runs stopped at the cap, recorded as "no Medium+ finding, format deviation on record", and **four mutation experiments** carried the A12 dedicated-corroboration role in ⑤ (giving a judgement independent of MAI on whether the checks really discriminate). **That layer's outcome is not proof of review sufficiency.**

### 9.3 ④ Codex final review

| Round | Purpose | Verdict | Findings |
|---|---|---|---|
| 1 | first final review (**no ③/③.5 list attached**) | `FINDINGS` | 2: **Medium** (the B4 section's dependency line lacked `B3c`, contradicting the status table and dependency graph) + **Low** (stale line anchors in the B3c item) |
| 2 | re-review of the fixes | `FINDINGS` | the only blocker was "no raw unified-patch hunk evidence" (the subagent has no terminal); every substantive item was judged closed |
| 3 | closing round (main controller wrote the raw diff to disk for it) | **`PASS`** | 0 (with measured CN/EN isomorphism: the four clause paragraphs total **16/16** sentences and **13/13** bold markers) |

### 9.4 A self-inflicted defect found and fixed in ⑤ (honest record)

- **A character-level defect introduced by the main controller**: while writing Chinese through `\uXXXX` escapes, 「拒**绝**」(U+7EDD) was mistyped as 「拒**绍**」(U+7ECD), contaminating **5 places** in one go (`docs/CONTRACTS.md` §7 twice, `docs/t027/REMAINING_SLICES.md` three times); worse, the replace tool still **reported success** against a look-alike oldString, masking the error.
- **How it surfaced**: a ③ re-review Low (B-N1 terminology) → main-controller localisation → repository-wide check confirmed 0 occurrences in HEAD and 5 in the working tree (i.e. introduced by this slice).
- **Fix and independent verification**: all 5 were corrected to 「拒绝」; an **anomalous-character scan** then ran (the CJK character set of the `git diff --unified=0` added lines minus the whole-HEAD `.md` corpus) ⇒ **empty difference set**; the repository-wide `拒绍` count is **0**; the ④ re-review independently confirmed 0 as well.
- The lesson ("write Chinese characters literally, never through `\u` escapes; run the anomalous-character scan after every edit") is recorded in repository memory.

### 9.5 A11 spot check (symmetric sampling)

- Of ④ Codex's Section D assertions, the main controller actually executed **4** (D1 three-state / D3 discrimination gate and mechanism evidence / D4 platform × premise boundary / D8 B4 dependency closure) and turned them into executable checks; **D8 was in fact already red on the first round** (the missing B4 dependency), matching the Medium that the review had found — so the sampling both verified "delete it and it goes red" and reproduced a real defect. Items, expectations and observed values are in §5.
- **Missing-side declaration**: neither ③.5 MAI call produced an executable "delete the guard ⇒ goes red" item, so the sampling could not land on the MAI side; recorded as a known limitation (sampling provides **existence corroboration** only, never a coverage proof).

## 10. Explicitly not executed

- **No `commit` and no `push`**: per discipline the slice stops in a reviewable state awaiting explicit same-turn authorisation.
- **Windows `unproven` not lifted**, `first-dispatch` not relaxed, the AT-23 and T027 conclusions unchanged.
- **No real model calls and no paid probes**; **no B3a/B3b rig re-run**; `tools/agent-probe/enforcement-proxy` **untouched**.
- **No `.go` / `.json` / schema / fixtures change**.
- **No CI observation** (nothing has been pushed); the CI result for this working-tree change is pending the push.

## 11. Next-step recommendation (§7.4)

> Per the dependency graph the next slice should be **B4 · Windows T027 end-to-end (AT-23)**; this slice has closed the gap between "contract wording" and "tiered verdict", so B4 can take the §7 durability-upgrade clause as a direct input.
> Suggested tiers — main controller **deep (Max)** (only for the ① design and ② E2E-plan writing steps), V4 Pro **①③ mandatory** (B4 touches stop/identity/ownership semantics), ③.5 **enabled** (a behaviour-change slice; note the MAI output format is currently unstable), Codex **Extra High (④ mandatory + the §3.11 blind audit)**.
> A parallel item to clear before B4: **DR-2** (the `enforcement-proxy` intermittent flake needs its own slice).

- **Guidance item to discuss (continuing B3b's D2 ruling)**: should hard-gate slices have their own ④ call cap? This slice's ④ measured 3 calls (1 final review + 1 re-review + 1 closing round), and the closing round was caused directly by **the reviewer lacking a way to read the unified patch**. **Suggestion**: write "put the raw diff on disk" into the ④ input procedure (verified effective here) to remove format-caused wasted rounds.
