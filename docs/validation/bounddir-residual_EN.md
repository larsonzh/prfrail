# P84 bounddir residual remediation validation report (BOUNDDIR-RESIDUAL)

Date: 2026-11-21. Status: `PARTIALLY CLOSED` (partial closure). Slice identity: P84 BOUNDDIR-RESIDUAL

**Verdict: `PARTIALLY CLOSED`** (partial closure) — T5#4 (F4 residual) is fully closed; T4 (F1 residual) is largely closed, with the residual = N13-UNC (the two canonical spellings of a UNC share root are rejected as unequal; a deterministic in-bound false rejection, fail-closed in direction, an availability defect); plus one Low evidence-attribution item (the W10 mutation annotation contradicts the r3 transcript). **This must not be read as "closed".**

---

## 1. Slice identity and authorization

| Item | Content |
|---|---|
| Identity | `BOUNDDIR-RESIDUAL` (remediates P84-d and P84-e from the P83 report section 8; closes F1 / F4) |
| Authorization | User designated the start on 2026-11-21; three scheduling rulings the same day; one narrow re-review as a separate budget exception |
| Basis | P83 report section 8 deferral list; the T4 / T5#4 transcripts in the P83 evidence pack; v67 snapshot section 4 queue item |
| Boundaries | Fix T4 / T5#4 only; do not touch F3 / F6 / F8 (P85 range); no S1 Exit gate change; no threat-model baseline change; no directive version bump; no gitee |

## 2. Step 0 anchor measurements

| # | Anchor | Reading | Result |
|---|---|---|---|
| 1 | HEAD | `71d03c3` | match |
| 2 | origin/main | `71d03c3` | match |
| 3 | Directive version | v1.71 | match |
| 4 | SELFTEST | 740/740 | match |
| 5 | Three scopes | index / tree / ci | all `TOTAL_FAIL=0` |
| 6 | Rolling window v65 (LF count) | 220 | match |
| 7 | Rolling window v66 (LF count) | 216 | match |
| 8 | Rolling window v67 (LF count) | 216 | match |
| 9 | Mechanical readings | ADR 47 / OB 139 / DR 15 / pilot deferral 25 | match |

## 3. Goal and scope

Take the two residuals from the P83 report section 8 verbatim, neither widened nor narrowed: T4 makes artifact enumeration and descendant opens still provable for ownership on Windows; T5#4 makes the update from a lease generation to the task proof a recoverable transaction. One falsifiable closure assertion, at least one counterexample regression test and one mutation check per item; the (1) stage rules the contract impact before any code moves.

## 4. Implementation and closure status

| Item | Face | Main change locations | Closure status |
|---|---|---|---|
| T4-a | F1 | `openDescendantRelative` in `internal/evidence/bounddir_windows.go` (component-by-component chain rooted at the bound handle) | **closed**: W2 counterexamples on both sides (outside refused, inside admitted) |
| T4-b | F1 | `isUnderDirectory` and `boundDirVerify` in the same file (Ordinal replaces EqualFold) | **closed**: W1 a same-named directory differing only in case is refused, and so is the end-to-end path |
| T4-c | F1 | `reparseStoredTarget` and `boundRelativeSuffix` in the same file (the stored target is the sole authority, in-bound is guaranteed by construction) | **closed**: W5b / C12 / W8 / W9 counterexamples |
| T4-d | F1 | `boundRootName` in `internal/evidence/bounddir_clean.go` and both call sites | **closed**: W3 two root spellings agree and are refused with the same error |
| T4-e | F1 | the existing short-circuit in `internal/gates/runner.go` (no code change in this item) | **closed (pinned)**: the existing no-glob counterexample stays green |
| T4-f | F1 | `readDirNames` and the root enumeration in the same file (reads the bound handle directly) | **closed**: W6 root enumeration performs no name resolution |
| T4-g | F1 | `openRelativeComponent` and the component length guard in the same file | **closed**: W7 an oversize component is refused before the open |
| T4-h | F1 | the embedded-NUL guard in `reparseStoredTarget` in the same file | **closed**: W11 dual parse counterexample (the narrow re-review judges N14 CLOSED) |
| T4-i | F1 | the equality form in `boundRelativeSuffix` (a drive root keeps its trailing separator) | **residual N13-UNC**: the drive-root equality is fixed, the two canonical UNC share-root spellings are still rejected |
| T5#4-a | F4 | `persistRenewedPublishLease` in `internal/console/runtime_supervised.go` and `verifyPublishLeaseHistory` in `runtime_agent.go` | **closed**: C1 a foreign proof is not rewritten and the batch errors, C2 multi-generation history is caught up |
| T5#4-b | F4 | `RunLocked` in `internal/guard/lease.go` and the three recovery branches in `runtime_agent.go` | **closed**: C3 cross-process re-entry recovers, C6 a concurrent advance refuses the whole batch, C7 the chain depth bound is fail-closed |
| T5#4-c | F4 | `verifyPublishLeasePredecessor` in `runtime_agent.go` (pinned as is) | **closed**: C4 / C10 shape-matching but mismatching proofs are still refused |
| T5#4-d | F4 | `publishLeaseBelongsToRun` and `assertPublishLeaseCurrent` in `runtime_agent.go` | **closed**: C15 a foreign consumer is refused, C16 the proof is written inside the lock, C18 the entry holds the lock |
| Lock base | F4 | `internal/guard/lease_lock_windows.go`, `lease_lock_unix.go`, `lease_lock_unsupported.go` | **closed**: C9 the kernel releases a dead holder's lock, C11 a live holder keeps it, C17 a legacy lock directory is reported separately |
| Narrow re-review remediation | F1 | the equality form in `boundRelativeSuffix` and the embedded-NUL guard in `reparseStoredTarget` (remediation after the (7) narrow re-review) | N13 = **partial** (drive-root equality fixed; the UNC share-root false rejection is deferred, characterised as a deterministic in-bound false rejection, fail-closed in direction, an availability defect); N14 = **closed** |

## 5. Contract impact

| Item | Ruling | Rationale |
|---|---|---|
| T4 (F1 residual) | **DOES-NOT-TOUCH / IN-BOUNDARY** | the quoted sentence fixes a requirement, not a mechanism; handle-relative chaining, Ordinal comparison and stored-target parsing are `internal/evidence` implementations of it; `ntdll!NtCreateFile` is the same kind of call as the existing `kernel32` LazyDLL use, an OS API rather than a third-party package, and `go.mod` is untouched |
| T5#4 (F4 residual) | **DOES-NOT-TOUCH / IN-BOUNDARY** | the quoted sentences are already implemented by the store + proof architecture; the proof describes itself as outside the recordHash set and the lease files are the store-local layout of `internal/guard`; only internal chain validation and OS lock wrappers are added, with no new schema, CLI or dependency |

Neither item touches a contract, so `docs/CONTRACTS{,_EN}.md`, the schemas, the fixtures and `go.mod` are unchanged. The conduct rule held: the implementation introduced no new schema record, no new CLI command, no new dependency and no third-party package.

## 6. Triage and user rulings

| Item | Ruling | Record |
|---|---|---|
| Controller carrier | switched to `DeepSeek V4 Pro (deepseek)` High, replacing Flash | user ruling of 2026-11-21; recorded in section 12 |
| (7) budget | 1 final + 3 re-reviews | exceeds the section 7.5 literal "M-L counts as M"; a direct user ruling that may not be cited as a precedent |
| Blind review | one mandatory round, counted separately as +1, runs before the final review | T4 touches identity semantics and T5#4 touches ownership semantics (section 7.7) |
| (7) tier | `GPT-6.2 Sol (copilot)` Max | same-round tier authorization; section 8.4 records the set time and the pause state |
| One narrow re-review | allowed, input limited to the N13 / N14 remediation | budget exception, exceeds section 7.5 and this slice's budget ruling, may not be cited as a precedent |
| R accounting | remediation re-runs inside the re-review ring do not count as R; R = 3 | section 9.2 accounting |
| Start-message "stop at R 2" | conflicts with section 9.2, the directive prevails | handled under section 13.2.11 discipline F, registered, not recorded as behaviour drift |

## 7. Review chain

| Stage | Carrier | Rounds | Verdict |
|---|---|---|---|
| (1) design | `deep-reasoner` + `DeepSeek V4 Pro (deepseek)` Max | Rev 1–6 | contract ruling DOES-NOT-TOUCH / IN-BOUNDARY |
| (2) implementation | `prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)` | 11 rounds | item-by-item landing, including the Rev 3b lock-base rework |
| (3) testing | `prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)` | 8 rounds | counterexamples W1–W11 and C1–C18; every mutation turned red and was restored with SHA256 |
| (5) pre-review | `DeepSeek V4 Pro (deepseek)` | 2 / 2 (exhausted) | r1 PASS WITH FIXES → r2 PASS |
| (6) independent scan | `MAI Code 2 Flash (copilot)` | 2 / 3 | r1 FINDINGS → r2 PASS |
| (7) blind review | `GPT-6.2 Sol (copilot)` Max | 1 | FINDINGS (2 High + 3 Medium + 1 Low) |
| (7) final review | `GPT-6.2 Sol (copilot)` Max | 1 / 1 | FINDINGS |
| (7) re-reviews | `GPT-6.2 Sol (copilot)` Max | 3 / 3 (exhausted) | r1 / r2 / r3 all FINDINGS |
| (7) narrow re-review | `GPT-6.2 Sol (copilot)` Max | 1 (budget exception) | FINDINGS (N13-UNC residual) |

Every transcript is archived under `docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/`; design revisions 1–6 are archived under `01-arch/` in the same pack. The (7) narrow re-review is the last read, and its FINDINGS (N13-UNC residual) are handled under the P83 report fallback close-out (b) precedent: the residual is recorded as is and deferred to a later slice, not restated as PASS.

## 8. Items deferred to P85

| Id | Deferred content | Source |
|---|---|---|
| none added | this slice adds no P85 payload; N13-UNC is an F1 residual while the P85 boundary holds F3 / F6 / F8 only, so whether to fold it in needs a user ruling | (7) narrow re-review |

N12 (updating the evidence-pack MANIFEST description) is not a deferral; it lands with the ⑧b close-out.

## 9. Documentation sync

| Object | Disposition |
|---|---|
| This report pair | adds `docs/validation/bounddir-residual{,_EN}.md` (CN / EN strictly mirrored) |
| Ledger | the P84 section and completion record in `docs/t027/REMAINING_SLICES{,_EN}.md` (lands with ⑧b) |
| Installation docs | `docs/INSTALLATION{,_EN}.md` gains 2 lines each: the legacy Mkdir lock upgrade boundary and the manual ordering |
| `SECURITY` note | the line 5 note in `docs/SECURITY{,_EN}.md` states this report's verdict truthfully (lands with ⑧b) |
| `CONTRACTS` / schemas | unchanged (both items judged DOES-NOT-TOUCH) |
| ADR | none added |
| Not changed | `CHECKLIST_v67{,_EN}`, `s1-exit`, the `p57` gate table, the directive itself, CI, `DEV_PLAN` |

## 10. Gate results

| Check | Result |
|---|---|
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 740/740` |
| `--all --skip-base --scope=index` | `TOTAL_FAIL=0` |
| `--all --skip-base --scope=tree` | `TOTAL_FAIL=0` |
| `--all --skip-base --scope=ci` | `TOTAL_FAIL=0` |
| `go test ./... -count=1` | EXIT=0, all 15 packages ok (measured by (8)b on the Windows host) |
| `gofmt -l internal cmd` | no output |
| IDE diagnostics | zero warnings on the files this slice changed |
| Cross-platform diagnostics | 6 entries (2 each in guard / console / evidence under the js / wasm targets), no effect on the windows / linux / darwin builds |
| Encoding | changed `.go` are UTF-8 without BOM and LF; added `.md` are UTF-8 with BOM and LF |

The table above holds the (4)-gate-time readings, re-checked by the (8)b measurement (`go test` 15 packages ok, `SELFTEST` 740/740); after (8)a wrote the pair and again before (8)b staged anything, the tree scope was re-run and the only FAIL was G7-b (declared artifact not staged), which is inherent before `git add` (the (8)b re-run readings are in `tmp/P84-BOUNDDIR-RESIDUAL/gate-tree-postNarrow.txt` and `tmp/P84-BOUNDDIR-RESIDUAL/gate-tree-8b.txt`); after (8)b stages the paths the controller re-runs the index / tree scopes and they go to zero - that expectation is not an absolute assertion and the re-run output governs.

## 11. (9) validation results and items not executable on this host

| Item | Result |
|---|---|
| Windows host `go build` / `go vet` / `go test ./...` | `go build` / `go vet` ok; `go test ./... -count=1` EXIT=0, all 15 packages ok |
| `-race` | not executed (no gcc on this host); carried by the CI ubuntu Race step |
| Ubuntu VM (`10.0.0.199`, go1.22.12, `bash -lc`, `GOPROXY=https://goproxy.cn,direct GOSUMDB=off`) | `go build` / `go vet` ok; `go test ./... -count=1` EXIT=0, all 15 packages ok; targeted `-race` on 4 packages (evidence / guard / console / adapters) ok |
| CI ubuntu Race step | first run after the push (backfilled under section 11.2 once (10) authorises it) |

Every reading in this table is backfilled by (8)b from measurement; (8)a states no estimate. If any of the three layers reports FAIL, the verdict must be restated from that measurement instead of keeping this section's conclusion. Items not executable on this host are registered: the Linux-tagged cases (bounddir unix case sensitivity and reparse, lease transaction recovery) are covered by the VM layer; the W1 end-to-end case-sensitive directory is created with `fsutil` and runs successfully on this host without being skipped.

## 12. Environment and tiers

| Role | Carrier and tier | Set at | Paused |
|---|---|---|---|
| Controller | `DeepSeek V4 Pro (deepseek)` High | 2026-11-21, slice start | no |
| (1) architect | `deep-reasoner` + `DeepSeek V4 Pro (deepseek)` Max | 2026-11-21 | no |
| (2) implementer | `prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)` High | 2026-11-21 | no |
| (3) tester | `prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)` High | 2026-11-21 | no |
| (5) pre-review | `DeepSeek V4 Pro (deepseek)` | 2026-11-21 | no |
| (6) scan | `MAI Code 2 Flash (copilot)` | 2026-11-21 | no |
| (7) blind / final / re-review / narrow | `GPT-6.2 Sol (copilot)` Max | 2026-11-21, the (7) stage of this slice | no |
| Batch assist | `GLM-5.3-Flash (glm)`, strictly serial, four task classes only | 2026-11-21 | no |

## 13. Review accounting

| Item | Value | Note |
|---|---|---|
| R | 3 | one (5)-ring remediation re-run plus one (6)-ring remediation re-run plus one (7)-final remediation re-run; re-review-ring remediation does not count as R |
| D | 0 | no temporary-file discipline violation; the `%TEMP%` probe file that the editor buffer re-wrote after deletion is corrected in the revision record and not counted as D |
| Probes | 0 | no real external call |
| Batch assist GLM | 2 times | extract checking plus bilingual mirroring, all PASS |
| (7) budget | blind 1 + final 1 + re-review 3 + narrow re-review 1 (budget exception) | four elements: decider = the user, date = 2026-11-21, reason = N13 is a deterministic in-bound false rejection whose fix direction is clear and small, one narrow re-review is worth it, basis = the user's ruling message; **may not be cited as a precedent** |
| High-risk mark | enabled (section 9.2) | from R 2 onward every ring runs the section 5.2 high-risk row |
| (5) / (6) caveat | their PASS binds the pre-remediation snapshot | the final snapshot is covered only by the (7) narrow re-review, which is still FINDINGS |
| Rule conflict | the start-message "stop at R 2" conflicts with section 9.2 | handled under section 13.2.11 discipline F with the directive prevailing, registered |

## 14. Explicitly not executed

No gitee action; no directive version bump; no new ADR or OB; no CI or `DEV_PLAN` change; no `-race` (no gcc on this host); the Linux-only cases were not executed on a real Linux host (left to the Ubuntu VM and the CI ubuntu leg); N13-UNC is not fixed (a residual whose disposition needs a user ruling); of the (9) three-layer readings, the Windows host and VM layers are backfilled (section 11) and the CI ubuntu leg's readings are backfilled after the push (section 11.2).

## 15. (8)a self-check

(1) Every "closed" carries file:line or test-name evidence (the section 4 table maps row by row onto the case names in `internal/evidence/p84_bounddir_residual_test.go` and `internal/console/p84_t5h4_residual_test.go`). (2) No superseded wording survives: the start-message claim that "the root is spelled as an empty relative path while an explicit `.` is refused" is restated as the shared-classifier behaviour that the code actually implements. (3) Readings carry their time point: sections 2 and 10 are (4)-gate-time readings and section 11 was a placeholder that (8)b has now backfilled with measured readings.

## 16. Six-stage self-check (sixstage)

```sixstage
## Section A 总体结论
PARTIALLY CLOSED（T4 residual N13-UNC, 1 item; all else closed）

## Section B 发现清单

| SEVERITY | Location | Residual |
|---|---|---|
| Medium | N13-UNC | the two canonical spellings of a UNC share root are rejected as unequal (deterministic in-bound false rejection, fail-closed in direction, an availability defect) |
| Low | W10 evidence attribution | the mutation annotation "pre-round-11 original form" contradicts the r3 transcript; the two must be labelled separately |

## Section C 逐条证据与推理
- The four T5#4 branches and the lock base are pinned by the C1–C18 counterexamples, item by item in the section 4 table.
- The T4 chain, Ordinal, root semantics and NUL guard are pinned by W1–W11, item by item in the section 4 table.
- Every (7) transcript is archived under `04-reviews/` of the evidence pack; design revisions 1–6 under `01-arch/` of the same pack.

## Section D 可证伪性审计
- Restore the TrimRight form of the `boundRelativeSuffix` equality: the root-equality leg of W10 turns red immediately.
- Delete the embedded-NUL unit scan: both direct-parse subtests of W11 go from refusal to a successful decode.
- Fall the chain validation in `persistRenewedPublishLease` back to a `ResourceID` comparison: C1 turns red immediately.
- Delete the catch-up branch or the Reclaim path in `renewPublishLease`: C3 turns red immediately.
- Rewrite this report's verdict as "closed": the mutually-checking residual row in section 4 and this section's conclusion turn red.
- UNCOVERED: the reachability of the UNC share root at a production entry is not measured; the per-item mutation replay and the (9) three-layer readings are backfilled by ⑧b.

## Section E 我无法验证的部分
The reachability of the UNC share root at a production entry, the per-item mutation replay and the (9) three-layer readings are all outside what can be verified statically on this host.

INDEPENDENT SCAN: FINDINGS
```

## 17. Artifact fence

```artifacts
docs/INSTALLATION.md	repo
docs/INSTALLATION_EN.md	repo
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/bounddir-residual.md	repo
docs/validation/bounddir-residual_EN.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/00-briefing.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/01-arch/01-arch-design.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/05-pre-review-r1.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/05-pre-review-r2.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/06-scan-r1.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/06-scan-r2.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-blind-review.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-final-review-r1.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-narrow-re-review.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-re-review-r1.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-re-review-r2.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/04-reviews/07-re-review-r3.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/MANIFEST.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/SHA256SUMS.txt	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/freedom-list.md	repo
docs/validation/evidence/P84-BOUNDDIR-RESIDUAL/input-list.md	repo
internal/console/p83_fixround7_tester_test.go	repo
internal/console/p84_t5h4_residual_test.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_supervised.go	repo
internal/evidence/bounddir_clean.go	repo
internal/evidence/bounddir_unix.go	repo
internal/evidence/bounddir_windows.go	repo
internal/evidence/p84_bounddir_residual_test.go	repo
internal/guard/lease.go	repo
internal/guard/lease_lock_unix.go	repo
internal/guard/lease_lock_unsupported.go	repo
internal/guard/lease_lock_windows.go	repo
tmp/P84-BOUNDDIR-RESIDUAL/	local-only
```
