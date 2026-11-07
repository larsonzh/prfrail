# Task execution checklist (cross-session snapshot v50 · P60 closed + R-08/R-09 ruled)

Generated: 2026-11-07, HEAD = origin/main = `756c033` (the v49 note blank-line fix commit), directive v1.64, CI run `42359958580` green on both legs. **This snapshot and its sibling commit are the next commit after `756c033`; its hash is unknown at generation time**; once that commit lands, HEAD advances one step, and `git rev-parse --short HEAD` gives the live truth.
Supersedes: fully replaces v49 and any earlier snapshot (historical snapshot files are kept, not deleted). Backfill commits (docs-only) only advance HEAD and do not change this section's anchor semantics.

Changes this version (vs v49): (1) **P60 `DIRECTIVE-JUDGE-BLINDSPOTS` closed** (`e381898` + `dafe004`, runs `42323986343` / `42328605588` green on both legs; the J2×J3 ruling landed at directive v1.64; G1(e2)/G5(e2) added; OB-138 flipped `已处置`; SELFTEST 703 → 722); (2) **R-08/R-09 ruled and landed** (ruled 2026-11-06, `docs/validation/r08-r09-ruling{,_EN}.md`; landed in `5164734`, backfill/note commits through `756c033`; section 6's two R-table rows flipped); (3) **directive v1.64** (authorized by the J2 ruling, not a self-made in-slice decision); (4) section 1 anchors updated (HEAD/origin = `756c033`, v1.64, SELFTEST 722/722, plus a new v49-line-count anchor); (5) section 2's fact pack switches to P60 and P59 moves to section 3; (6) layer 1 item 8 unblocked (R-08/R-09 ruled), layer 2 item `G-4` closes and dequeues via P60; (7) section 5 observation list updated (OB-138 leaves); (8) section 6 R-08/R-09 flips to ruled (inherited from v49, no timestamp notes); (9) section 8 gains the discipline-Z candidate; (10) section 9 gains gate pitfalls #23 / #24 (normalization-script pit + G6-4 pit); (11) v49's three timestamp notes do not carry over (v50 itself reflects the ruled status).

## 1. Current anchors (check these 6 numbers first)

| Item | Value | Check command |
|---|---|---|
| HEAD | `756c033` | `git rev-parse --short HEAD` |
| origin/main | `756c033` | `git rev-parse --short origin/main` |
| Directive version | **v1.64** | header of `docs/DELIVERY_DIRECTIVE.md` |
| SELFTEST | **722/722** | `node tools/gates/gate.js --selftest` |
| Three scopes | **all TOTAL_FAIL=0** | `gate.js --all --scope=index`, `tree`, `ci` |
| v49 line count | **233 (CN/EN symmetric)** | line count of `docs/t027/CHECKLIST_v49{,_EN}.md` |

If ≥2 items mismatch ⇒ the repository reality prevails; report the final readings first, then continue.

Fixed `tmp` assertion: top level = `.gitkeep` + `gates-cache`; `gates-cache` holds only `g4b`; P60's `tmp/P60-DIRECTIVE-JUDGE-BLINDSPOTS/` was cleared at 8c; the P60 ⑩ D-2 residual `gate-v49-tree.txt` was cleared at this version's generation.

## 2. P60 fact pack

| Field | Value |
|---|---|
| Slice / size | `DIRECTIVE-JUDGE-BLINDSPOTS` (governance-class `[SLICE]`, layer 2 item `G-4`) / M (J2-led; the J1/J2/J3 ruling) |
| Goal | Close the OB-138 criterion blind spot (in-block completion-record field rows are outside the G1(e)/G5(e) coverage); ruling = J2×J3 combination, J1 not chosen |
| Change surface | directive `DELIVERY_DIRECTIVE{,_EN}` v1.64 (header / §0.3 / §6.3 G1/G5 / §12.4 OB-138); `tools/gates/lib/criteria/g1.js` / `g5.js` (G1(e2)/G5(e2)); `tools/gates/selftest.js` (T119–T124); `REMAINING_SLICES{,_EN}` (P60 block + status-table row); `ADR_REGISTER{,_EN}` (ADR-034); report pair and evidence pack |
| Core output | **v1.64** (the G1(e2)/G5(e2) sub-criteria; OB-138 flipped `已处置`); SELFTEST 703 → **722** (+19, T119–T124 positive/negative fixtures) |
| Review chain | ① `DESIGN: DONE`; ② two rounds (r2 fixes F-1/F-2) → ③ two rounds `TEST: DONE` → ⑤ r1 `PASS WITH FIXES` (M-1/M-2) → r2 `PASS`; ⑥ `INDEPENDENT SCAN: PASS`; ⑦ `FINAL REVIEW: PASS` + blind review r1 `FINDINGS` (1 Medium = the G7-f expected item, closed by the ⑧ transcription) → r2 `FINDINGS` (3 wrap-up-timing items, fixed) → r2b `PASS` (final, doubles as the post-fix ⑦ re-review); ⑦ metering = final-review 1/1 + blind-review 1/1 + re-review 2/3, within quota |
| Commit / CI | `e381898` + `dafe004` (⑩ clarification) runs `42323986343` / `42328605588` green on both legs: both-leg Gates `SELFTEST: PASS 722/722` + `TOTAL_FAIL=0`; Test 36 packages ok; green on the first attempt, no DR registered |
| Metering | `R = 0`, temp-artefact D = 2 (D-1 blind-review order deviation from §7.7, accepted per ruling (a); D-2 tmp top-level residual `gate-v49-tree.txt`, kept for v50 per ruling (b)); every `runSubagent` passed `agentName` + `model` explicitly; probes and paid external calls 0 |
| Rulings | J2×J3 combination (J1 not chosen); the v1.64 bump is authorized by the J2 ruling; ⑦ counting classification (final 1/1 + blind 1/1 + re-review 2/3); D-1 ruling (a); the r2b equivalence statement; ADR-034 (controller V4 Pro, step 0, per the OB-37 meta-rule) |
| Document sync | directive v1.64 four sites; ledger P60 closed + status-table row; ADR-034; report pair and evidence pack (16 payloads, MANIFEST / SHA256SUMS self-consistent) |
| Capability boundary | no product code / CI / CONTRACTS / schema / fixture change; `cjk-newwords.txt` registered only the char `稽` (⑥ verbatim, reason in the P60 report); no ruling on the user's behalf for J1/J2/J3; no D-6 class overreach (no self bump / no forged authorization / no over-quota dispatch) |
| Not closed | v50 generation (this snapshot); AT-23 start-up (layer 1 item 8 unblocked, awaiting the user's briefing) |

## 3. Completed slices (recent baseline)

| Version | Slice | Main commit | Conclusion |
|---|---|---|---|
| v1.60 | P49 | `75990a8` + backfill commit | closed (R = 1, D = 2; run `41445657595` green on the first attempt) |
| v1.60 | P50 | `81f9dc9` + backfill commit | closed (R = 1, D = 1; run `41507722672`) |
| v1.60 | P51 | `2b0af2d` + backfill commit | closed (R = 1, D = 2; run `41569353460`) |
| v1.60 | P52 | `c636dd8` + backfill commit | closed (R = 0, D = 1; run `41642139932`; GAP-77 breach evidence) |
| v1.60 | P53 | `7802819` + backfill commit | closed (R = 0, D = 0; run `41707160563`; GAP-77 fixed) |
| v1.60 | P54 | `7c0763a` + `f38e0ce` | closed (R = 1, D = 1; run `41840557379`; C6 decision path + degradation ruling) |
| v1.60 | P55 | `63e2e24` + `c76c6ce` | closed (R = 0, D = 1; runs `41925623071` + `41928327517` both green; AT-13 two-sided evidence chain) |
| v1.61 | P56 | `504caff` | closed (R = 2, D = 1; run `41978355255` both green; G4b corpus enumeration decoupled, SELFTEST 697 → 703) |
| v1.61 | P57 | `7a43dfa` | closed (assessment class; R = 0, D = 0; run `41991695581` both green; P5 no slice + P6 eight-gate checklist) |
| v1.62 | P58 | `f798cea` + `8d30e66` | closed (assessment class; R = 0, D = 0; runs `42062391497` + `42069357349` both green; P7 no slice, row kept) |
| v1.63 | P59 | `75b90a8` | closed (governance class; R = 0, D = 0; run `42121119502` green on the first attempt; R-12 ruled A2×B1×Y1 landed) |
| v1.64 | P60 | `e381898` + `dafe004` | closed (governance class; R = 0, D = 2; runs `42323986343` + `42328605588` both green; J2×J3 landed, SELFTEST 722/722) |

> **Note on the D column in §3**: each "D = n" in this table is shorthand for the temp-artefact discipline-violation count (§5.1), synonymous with "temp-artefact D" in §2 / §11; the deviation-table D-1..D-n in a report's "anomaly handling" chapter is an independent sequence (§5.1 "naming distinction"), and the two are not interchangeable. P58's "D = 0" means no temp-artefact violation; its 6 deviations are in the "metering" row of §2 and P58 report §8.

## 4. Execution queue

**Ordering basis of this version**: layer 0 is dequeued, layer 2's first item (P7) was dequeued by P58, its second item (G-3 / P59) by P59, and its third item (G-4 / P60) by P60; the layer-1 hard gate (item 8 AT-23) is unblocked by the R-08 / R-09 rulings (2026-11-06) and can start; layer 2 is fully dequeued; layer 3 stays unstarted. The in-layer order is fixed and layers do not jump.

### Layer 0 pre-work (✅ done, dequeued)

| Seq | Item | Content | Size | Depends on | Boundary |
|---|---|---|---|---|---|
| 0-1 | P5 `T027-B` close-out checklist | ✅ done: **no new slice**; remainder trifurcates R-1 / R-2 / R-3 | assessment action | — | dequeued; see P57 report §3 |
| 0-2 | P6 `S1-EXIT-DOC` close-out checklist | ✅ done: S1 Exit eight-gate checklist (S-G1 to S-G8) | assessment action | — | dequeued; see P57 report §4 |

> Calibre: both items were closed and dequeued by P57; their conclusions are the **only source** for layer 1 and the S1 Exit, and no same-type assessment action may be raised again.

### Layer 1 main-line hard gate (advance immediately once unblocked)

| Seq | Slice | Identifier and content | Size | ⑦ quota | Depends on | Boundary |
|---|---|---|---|---|---|---|
| 8 | AT23-RESTRICTED-CANDIDATE | AT-23 restricted real candidate and sign-off | M | Codex 1 final review | **R-08 / R-09 ruled (2026-11-06) ⇒ unblocked** | no sign-off without a real-candidate authorization |

### Layer 2 routine governance (debt clearance)

| Seq | Item | Content | Start condition |
|---|---|---|---|
| G-1 | P7 governance-side split-out OB roll-up | ✅ done and dequeued: closed by P58, verdict = no slice, keep the P7 row ⏸ | dequeued |
| G-3 | `DIRECTIVE-SCOPE-GUARD` (P59) | ✅ done and dequeued: closed by P59 (2026-11-06), R-12 ruled A2×B1×Y1 landed at v1.63; conclusion in the report pair | dequeued |
| G-4 | `DIRECTIVE-JUDGE-BLINDSPOTS` (P60) | ✅ done and dequeued: closed by P60 (2026-11-07), J2×J3 landed (G1(e2)/G5(e2) + OB-138 `已处置`; v1.64) | dequeued |

### Layer 3 low-priority alternatives

| Seq | Item | Content | Start condition |
|---|---|---|---|
| L-1 | Five low-priority tasks | not started without a trigger | the user must name them explicitly |

> Layer boundaries: layer 0 and every layer-2 item are all dequeued and no assessment action is produced any more; the layer-1 block is lifted by the R-08 / R-09 rulings; layer 2 / layer 3 must not be read as "the main line resumed".

## 5. Observations

- Observing (13): `OB-1 / 7 / 10 / 33 / 87 / 99 / 100 / 104 / 105 / 107 / 111 / 113 / 135` (**OB-137 was flipped `已处置` by P59 and OB-138 by P60; both leave the observing list**).
- To be discussed (10): `OB-19 / 21 / 52 / 53 / 55 / 57 / 58 / 59 / 61 / 75`. **OB-120 was dispositioned by P56 (`已处置`)**.
- To do (1): `OB-69`. Non-final total **24**.
- GAP registry: **GAP-77 fixed (closed by P53)**; GAP-78/79/80 candidates (registered in P52 report §2b and the ledger).
- **P59 empirical observation (not an OB; registered in the P59 report; adopted here)**: the A2×B1×Y1 combination on its first governance slice — P59 showed no D-6 class overreach (no self bump / no forged authorization / no over-quota dispatch), n=1, an observation item, later governance slices keep comparing.
- **P60 empirical observation (not an OB; registered in the P60 report; adopted here)**: G1(e2)/G5(e2) took effect on this slice's own in-block completion row (the J2 self-coverage proof); the ⑦ counting classification (final 1/1 + blind 1/1 + re-review 2/3) stayed within quota.
- Mechanical readings: max OB in `§12.4` = **OB-138**; max deferral = the 25th; ADR max = **34**; DR max = **13**.

## 6. Awaiting user rulings (R-01 to R-12)

| # | Item | Status |
|---|---|---|
| R-01 | P46 strategic review project initiation | done |
| R-02 | `CONTRACTS-S7-DURABILITY` pace and ownership | done (closed by P54) |
| R-03 | P45 follow-up slice ownership | done (`RACE-LEG-EXPAND` closed) |
| R-04 | DR-13 "fixed + observe" calibre | pending |
| R-05 | whether DR13-FIX-D3 is a standalone project | done (initiated standalone and closed) |
| R-06 | OB-55 ownership and the T027-B pace | pending |
| R-07 | MVA, S1 moved-out items and the S2 commitment narrowing | pending |
| R-08 | candidate, Windows premise P and probe authorization | **ruled (2026-11-06): BYOK main path + local runs + conditional closure + stage progression; see `docs/validation/r08-r09-ruling{,_EN}.md`** |
| R-09 | AT-24 `unknown` acceptance boundary | **ruled (2026-11-06): relaxed to mostly acceptable, two bottom lines kept (public network blocked, external side effects explicitly confirmed); see `docs/validation/r08-r09-ruling{,_EN}.md`** |
| R-10 | overall budget and per-slice model caps | pending (blocks scheduling) |
| R-11 | scope of the historical corrections | done (tier B plus per-item differentiation) |
| R-12 | disposition form for the D-6 class overreach and OB-137 | **ruled (2026-11-06; A2×B1×Y1; landed at v1.63; ruling text in the P59 evidence pack `01-arch/ruling-2026-11-06-r12.md`)** |

> P58 ruling note (2026-11-05): v1.62 **ratified**; P7 stays `⏸` and the P7 row is unchanged; the 5 quota deviation confirmed and registered as ADR-031 (never a precedent); OB-137 / OB-138 registered as observing; the D-3 wording confusion goes to the next slice. R-04 / R-06 / R-07 / R-08 / R-09 / R-10 stay pending.
> P59 ruling note (2026-11-06): R-12 ruled A2×B1×Y1 and landed at v1.63 (the bump is authorized by R-12, not a self-made in-slice decision); ⑦ metering = 1 final-review + 1 re-review (within 1+3); no separate blind review (§7.7 conditions not triggered); ADR-032 / 033 registered; v49 generated. The P7 row stays unchanged.

## 7. Start-up and step-0 general rules

1. Check the anchors → step 0: register this slice's ledger row (`R/D = 0` baseline, temp location `tmp/<sliceID>/`), and write `input-list.md` and `freedom-list.md`.
2. Every fix slice must give **regression tests** and register "unrun tests/lint" explicitly.
3. ⑨ three scopes green + `SELFTEST 722/722` + `go build/vet/test ./...` green before committing.
4. Model bindings (hard): ①⑤ `deep-reasoner` + `DeepSeek V4 Pro (deepseek)`; ②③⑧ `DeepSeek V4.1 Flash (deepseek)`; ⑥ `independent-reviewer` + `MAI-Code-1.1-Flash (copilot)`; ⑦ `independent-reviewer` + `GPT-5.3-Codex (copilot)`; strategic-review rounds 1/3 `independent-reviewer` + `GPT-6.1 Sol (copilot)`; batch assist `deep-reasoner` + `GLM-5.3-Flash (glm)`. **Every `runSubagent` must pass both `agentName` and `model`.**
5. **Execution-rhythm discipline**: routine keeps going — ④ green ⇒ dispatch ⑤ directly; ⑤/⑥/⑦ reporting Medium+ ⇒ fix and re-run that ring in place; **⑩ is the only routine stop point**.
6. Cross-version / cross-platform triage: when CI and local diverge, first check toolchain versions and build-tag coverage; after cross-building re-check `go env GOOS GOARCH`; before deleting a "locally unused" symbol run at least one opposite-`GOOS` `go vet`.
7. Batch assist: each slice actually dispatches `deep-reasoner` + `GLM-5.3-Flash (glm)` **≥ 2 times** (strategy and output cap in section 8); **strictly serial**, no parallel subagents.
8. Controller model strategy: default `DeepSeek V4.1 Flash`; switch to Pro on bottom-layer concurrency/state-machine/handle-passing complexity, a single design document over 10,000 characters, deep root-cause analysis of large CI logs, or two consecutive Flash behavior drifts; **any downgrade must be independently registered per the OB-37 meta-rule** (P59 precedent: controller V4 Pro registered as ADR-032).
9. **Role carriers are not replaceable**: ② must be `prfrail-implementer` and ③ must be `prfrail-tester`; the controller only orchestrates, writes to disk, gates and receipts. **Exception**: pure evidence-supplement / log-grab slices may be self-executed by the controller, but need the user's same-turn explicit statement.
10. **The briefing message is an artifact**: step 0 must write the user's briefing verbatim to `docs/validation/evidence/<sliceID>/00-briefing.md`, put it into `MANIFEST.md` and `SHA256SUMS.txt`, and declare it in the `artifacts` fence; missing it ⇒ ① may not start.
11. **Linearization-point discipline**: for persistent state machines, "persisted successfully" must be the linearization point; the in-memory CAS is only a same-process fast rejection cache.
12. **Evidence-pack and fence sync discipline**: every new evidence payload must add its `MANIFEST.md` entry, recompute `SHA256SUMS.txt`, and add its fence row in the same round.
13. **Reviewer archival discipline**: ⑤⑥⑦ outputs must be archived verbatim (`04-reviews/`); the controller must not rebuild them from bullet points.
14. **Assembly-verification-slice discipline**: assertions must rest on falsifiable mechanical criteria; constant identities must declare they are not sentinels; the report must state whether production code changed (zero change).
15. **Assessment-class slices must stop at ⑩ after ⑨**: no implementation object; stop at ⑩ for the same-turn authorization after ⑨ is green.
16. **A single output must not exceed 10,000 characters**: batch before the cap; no self-restart; read as needed.
17. **Compression ends the slice**: at 100K context or on a compression notice ⇒ end the current slice immediately; one conversation runs one L-tier slice.
18. **The briefing must contain a "document sync update" section**: if the slice changes project phase or contract capability, stage ⑧ must update the authoritative documents; omitting it counts as a structural defect.
19. **Environment pre-check and reuse iron rule**: declare external environments at step 0; probe for missing tools (PATH, config, `find`) before a minimal install; **remote SSH must use `bash -lc` to load the login environment and probe user-level directories**; pin the version to CI parity and verify checksums.
20. **Evidence-capture replayability and fidelity**: the capture's `COMMAND:` line must match the real argv; if only a display abbreviation can be written, file an erratum in `anomalies.md` and put the faithful command + a `commandErrata` pointer into the derived `index.json`; drift/constructed inputs need reproducible rebuild steps in `replay.md`.

## 8. Cross-slice collaboration disciplines (the briefing must land them)

| Element | Content |
|---|---|
| Carrier table | ①⑤⑥⑦ use `deep-reasoner` / `independent-reviewer` (read-only); ②③⑧ use `prfrail-implementer` / `prfrail-tester` / `prfrail-documenter` (writable); `prfrail-tester` is forbidden on review-layer rings; batch assist must be `deep-reasoner`, `quick-researcher` forbidden |
| Dispatch-table mechanical read | each ⑩ dispatch-table `agentName` must be read verbatim from `transcript.arguments.agentName` with the script stdout attached; **if the session debug log carries no transcript (as in P56 / P58 / P59), the receipt must explicitly state "not mechanically verified" and register it as a known fact** |
| Dispatch serial | `runSubagent` blocks serially, no parallel subagent calls; pass `agentName` and `model` every time |
| Commit message | always English, `<type>: <summary>`; if it contains any non-ASCII character, use `git commit -F <UTF-8 file>` |
| Independently checkable vs not | independently checkable outputs (enumeration, fact comparison, citation-point lists) are not voided by model errors; non-independently-checkable ones (judgement, conclusions, review opinions) must be redone |
| ⑥ before ⑦ | ⑥ must finish before ⑦; when ⑥ is redone and the documents changed, the ⑦ conclusion must be re-evaluated |
| Temp location and cleanup | always `tmp/<sliceID>/`; only 8c cleans it, and only this slice's directory — never recursively delete a parent; measure character counts or record intermediate processes with inline scripts |
| Long-command redirection | no PowerShell `>`; use `\| Tee-Object -FilePath <path>` or `cmd /c "..."` |
| Gate execution | mid-slice uses `--all --skip-base`; green everywhere before commit |
| Discipline E · index-scope order | precise `git add` first → then `--scope=index`; the receipt must report `G6-1 scanned N` with N greater than 0 |
| Discipline F · spec-over-example boundary | when a briefing example conflicts with the spec, follow the spec but register the conflict in the receipt |
| Discipline G · the fence is the landing snapshot | `artifacts` fences sorted by path, TAB-separated, CN/EN byte-identical; touching a new file requires a new fence row |
| Discipline H · opposite-GOOS re-check | before deleting or changing constants/functions, run at least one opposite-`GOOS` `go vet` |
| Discipline O · output-length defense | a single output must not exceed 10,000 characters; batch before the cap; no self-restart |
| Discipline P · controller model-switch trigger | see section 7 item 8; a downgrade is registered per OB-37 |
| Discipline Q · batch-assist output cap | the dispatch must state the character cap; the receipt must report the actual characters |
| Discipline R · stale-tense batch clearance | ledger-vs-fact mismatches go through a one-shot clearance slice; code and commits are the fact source |
| Discipline S · session-length and compression | at 100K context or a compression notice ⇒ end the current slice immediately |
| Discipline T · quantifier shape | the ledger and reports must not write an "identifier directly followed by a quantifier" shape: G6-4 reads it as a line-count literal; write "that entry" or "that record" instead |
| Discipline U · no verbatim trigger strings | reports must not quote the strings that trigger a criterion verbatim |
| Discipline V · split large tasks, avoid half-written pollution | split large tasks into small serial dispatches; after a timeout always inspect `git status` / `git diff` for partial output and roll back precisely with `git checkout` when needed |
| Discipline W · bilingual report structure | a report pair must share order and shape: the same number of `##` headings, byte-identical `artifacts` fences, and `sixstage` fences with all five sections and the same final conclusion |
| **Discipline X · report before self-delete (candidate, v46)** | a role that discovers a stray temp file must **not delete it immediately**; it reports the path first and waits until the controller has taken the pre-snapshot readings (`tmp` recursive count and `git status --porcelain -uall`), then runs `Remove-Item`; violations are handled under the **non-self-heal tier** (section 5.1) |
| **Discipline Z · blank line before blockquote notes (candidate, v50)** | blockquote notes (`> **…**`) must carry a preceding blank line, otherwise renderers merge the note into the table/paragraph above; after writing, mechanically assert that every `> **` note line has an empty line above (count must be 0) — v49's three timestamp notes were once flagged for missing blank lines |

> **Discipline Y note**: discipline Y (no self-bump inside a slice) was promoted by P59 into a body clause of directive §1.7.2 (v1.63, carried by the three bullets); this table no longer repeats it (§11.2: one fact described in one place).

**GLM batch-auxiliary strategy**: `GLM-5.3-Flash (glm)` has ample capacity, so "at least twice" is only a floor; mechanical checks, long-document excerpts, bilingual mirror checks, dispatch drafts and log extraction should be delegated whenever possible; complex judgement and ranking are forbidden.

## 9. Checklist meta-rules

**When citing, use the opening sentence, not "item N".**

| Rule (opening sentence) | Content |
|---|---|
| "Regenerate when: every 2 slices / pending-ruling change / directive bump / queue change / user request / ≥2 anchor mismatch" | the regeneration trigger set |
| "Pending-ruling items must be listed one by one" | section 6 must list them one by one |
| "Queue items must carry the six elements" | identifier / content / size tier / ⑦ quota / depends on / boundary |
| "Rule obligation ≠ gate obligation" | rule text is not read as a mechanical criterion |
| "No physical-clock assertions" | observation items trigger on events |
| "Gates green before commit" | hard gate |
| "Dispatch-table agentName mechanical read" | read verbatim from `transcript.arguments.agentName`; state it truthfully when unavailable |
| "Governance slices do not write back DEV_PLAN" | authoritative ledger = appendix C.0–C.0x |
| "A slice name never changes once in the ledger" | names are stable |
| "Policy prevails when the authorized quantity conflicts with it" | never judge over-authorization or refuse policy execution on that ground |
| "One fact described in one place" | body or pointer, pick one |
| "Controller registrations may be reworded; reviewer texts must be verbatim" | two different standards |
| "Fact errors are independent of design decisions" | ledger/report mutual contradictions are fact errors |
| "The ⑩ receipt must carry the fact pack for the next checklist" | fields include HEAD, CI run, ledger additions, new observations, new discipline candidates, max OB / deferrals, ADR/DR max, tmp assertion, differences vs the previous version, unclosed items |
| "[ROUTINE] small fixes take no new slice number but keep the discipline" | same ⑨ gates + CI + fence rows |
| "The briefing message is an artifact" | `00-briefing.md` on disk + in MANIFEST and SHA256SUMS + fence declaration |
| "Carrier deviation = highest-grade independence downgrade" | register a D item with the independence-impact assessment |
| "Mechanical criteria outrank witnessed statements" | zero-model/zero-side-effect assertions must be falsifiable |
| "Concurrency assertions cover survival and release" | beyond the failure face, assert the winner survives and the slot releases for re-dispatch |
| "Layer boundaries must not be read wider" | capabilities a layer does not support must be declared |
| "Assessment-class slices must stop at ⑩ after ⑨" | no self commit/push |
| "A single output must not exceed 10000 characters" | batch before the cap; no self-restart |
| "Compression ends the slice" | at 100K context or a compression notice ⇒ end the current slice immediately |
| "The batch-assist floor cannot be replaced by controller self-checks" | if the actual dispatch count is short, top it up in the same round; "scripts already checked" never exempts |

### Gate pitfalls (P48 six + P49 seven + P52 two + P55 four + P56 one + P57 zero + P58 one + P59 one + P60 one)

> The P48-through-P55 entries follow v45, the P56 entry follows v46, P57 added none; **v49 adds #22 (from P59)**; **v50 adds #23 / #24 (#23 from the v49 D-calibre mini-fix evidence; #24 from P60)**.

20. **A subagent must report to the controller before self-deleting** (candidate discipline, P56 evidence): a role that creates a stray temp directory outside `tmp/<sliceID>/` deletes it immediately, making the self-heal tier's "reading-2 baseline" unobtainable. Correct order: **self-report → controller takes the pre-snapshot (`tmp` recursive count + `git status --porcelain -uall`) → self-delete → zero-residue re-read**; violations fall under the non-self-heal tier.
21. **PowerShell decoding UTF-8 stdout as GBK** (P58 evidence): `Select-String` / `Out-File -Encoding utf8` mangles the Chinese output of gates and `gh` — producing both false negatives (a constant 0 scan) and double transcoding. Correct: `node -e` with `execFileSync(cmd, args, {encoding:'utf8'})` reading and `fs.writeFileSync` writing as-is, or make scripts print ASCII-only summaries.
22. **A reviewer's conclusion-line value must strictly belong to the §4.1 closed set** (P59 evidence): the ⑦ first round printed `FINAL REVIEW: CHANGES_REQUESTED`, a value outside the set ⇒ handled as a §7.2 format deviation; a same-stage resend **counts toward the quota cap** (the first round consumed 1 final-review) — never write "invalid round, no quota". Controller dispatch texts should include the legal value enumeration.
23. **A normalization script re-adds the BOM to an already-BOM file** (v49 D-calibre mini-fix evidence): `readFileSync(f,'utf8')` keeps the leading U+FEFF, so prepending `\uFEFF` again doubles the BOM (G4a reports `bom=2`). Correct: first strip every leading BOM in a loop, then add one back; verify BOM=1 / double-BOM=0 / CR=0 after writing.
24. **G6-4 misreads 「a version number directly followed by 行」 as a line-count assertion** (P60 evidence): `\d{2,}\s*行` only matches "digits immediately followed by 行"; a version number followed by 行 trips it. Safe rewrite: replace 行 with a wording like 条目 (inserting a word between the digits and 行 is safe).

## 10. One-sentence opening for a new session

Here is my project-state snapshot (v50): please check the 6 anchors of section 1 first; if they match, continue with the section-4 queue, otherwise the repository reality prevails.

【Paste this full checklist v50】

【Supplementary documents】`REMAINING_SLICES.md` (P60 closed), `DELIVERY_DIRECTIVE.md` (v1.64), `docs/validation/p60-directive-judge-blindspots{,_EN}.md` (P60 report, `e381898` + `dafe004`, runs `42323986343` / `42328605588`), `docs/validation/r08-r09-ruling{,_EN}.md` (R-08/R-09 rulings, `5164734`).

Current slice: P60 is closed and layer 0 plus every layer-2 item are dequeued; next is layer-1 item 8 (AT-23), unblocked by the R-08 / R-09 rulings and ready to start, or another slice the user names.

## 11. Current state in one sentence

P60 closed (`e381898` + `dafe004` pushed to origin, runs `42323986343` / `42328605588` green on both legs; J2×J3 landed at directive v1.64; OB-138 flipped `已处置`; R-08 / R-09 ruled (`5164734`, see section 6); R = 0, temp-artefact D = 2); next = layer-1 item 8 (AT-23) unblocked and ready to start.
