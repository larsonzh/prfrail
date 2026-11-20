# P83 security findings remediation validation report (SECURITY-FINDINGS-REMEDIATION)

Generated: 2026-11-21. Slice: P83. Grade: `[SLICE]`, implementation type (security fix), size M-L.

**Verdict: `PARTIALLY CLOSED`** — F2 / F5 / F7 and the T9 narrowing are closed; F1 / F4 are partially implemented with the residuals deferred to P84; F3 / F6 / F8 are deferred to P84. **This must not be read as "closed" or "security-closed".**

---

## 1. Slice identity and authorization

| Item | Content |
|---|---|
| Identity | `SECURITY-FINDINGS-REMEDIATION` (remediates F1–F8 from the P79 report section 9) |
| Authorization | User designated the start on 2026-11-20; two same-day additional rulings (triage approval, one extra re-review); final ruling on 2026-11-21 |
| Basis | P79 report `docs/validation/s2-pre-independent-security-review{,_EN}.md` section 9; v66 snapshot section 4 queue item |
| Boundaries | Fix F1–F8 only; no S1 Exit gate change, no threat-model baseline change, no directive version bump, no gitee |

## 2. Step 0 anchor measurements

| # | Anchor | Expected | Actual | Result |
|---|---|---|---|---|
| 1 | HEAD | `7a63962` | `7a63962` | match |
| 2 | origin/main | `7a63962` | `7a63962` | match |
| 3 | Directive version | v1.71 | v1.71 | match |
| 4 | SELFTEST | 740/740 | 740/740 | match |
| 5 | Three scopes | all 0 | index / tree / ci all 0 | match |
| 6 | snapshot v64 | see v66 section 1 | see v66 section 1 (counted by LF) | match |
| 7 | snapshot v65 | see v66 section 1 | see v66 section 1 (counted by LF) | match |
| 8 | snapshot v66 | verified by v67 | see v66 section 1 (recorded) | recorded |

## 3. Goal and scope

Remediate the eight security findings registered by P79, in priority order F2 / F3 → F6 → F1 / F4 / F5 / F7 / F8; one falsifiable closure assertion plus at least one counterexample regression test with a mutation check per finding.

## 4. Implementation and closure status

| Id | Face | Severity | Main change locations | Closure status |
|---|---|---|---|---|
| F1 | face 1 | Medium | `internal/gates/runner.go`, `internal/evidence/bounddir_*.go` | **partially implemented**: root bound before execution, root not opened when no artifacts are declared, non-blocking open plus type check, FIFO refused; the T4 residuals (case folding, intermediate reparse, root `""` semantics) are deferred to P84 |
| F2 | face 2 | High | `internal/console/runtime_agent.go`, `internal/adapters/agent_runner_replay_dispatcher.go`, `internal/adapters/agent_runner_admission.go`, `internal/console/trusted_ledger.go` | **closed**: the ledger path is derived from the durable root, production never falls back to the in-memory copy, the clock is read after the ledger, the re-confirmation is the last external observation before launch, and a failure retains R with zero launches |
| F3 | face 2 | High | `internal/chain/review.go`, `internal/console/runtime_agent.go`, `internal/console/runtime_supervised.go` | **deferred to P84**: only the latest grant is accepted, the risk confirmation is separated from policy approval, publish-time re-reads landed; **T2 STOP** — AgentRunner has no persistable frozen policy artifact to compare `grant.PolicyHash` against |
| F4 | face 2 | Medium | `internal/console/runtime_agent.go`, `internal/console/runtime_supervised.go`, `internal/guard/lease.go` | **partially implemented**: process-level holder check, in-lock takeover with generation advance, per-generation lease retention, run-ownership check on proofs; the T5#4 residuals (batch rewrite without predecessor validation, cross-process interruption recovery) are deferred to P84 |
| F5 | face 2 | Medium | `internal/gates/runner.go`, `internal/gates/adapter.go` | **closed**: classification comes only from the frozen hook, admission runs unconditionally, a definition-digest or recovery-guarantee or external-system conflict is refused, and the adapter no longer fabricates placeholder digests |
| F6 | face 3 | High | `internal/console/ai.go`, `internal/adapters/ai_probe_copilot_executor.go`, `internal/adapters/ai_preflight.go` | **deferred to P84**: approval material carries the profile digest and the full binding, the approval source is derived from the chain-config directory, validation precedes the credential read, the approved digest plus file identity travel to the spawn; **T3 residuals** — the Windows name-to-image binding window, and Linux same-inode not meaning byte-immutable |
| F7 | face 3 | Medium | `internal/guard/process_windows.go`, `internal/console/secret_input_windows.go` | **closed**: terminator headroom, per-entry scratch and discarded arrays zeroed, terminal-input growth copies zeroed, partial failure block zeroed, each with an observer seam and a counterexample |
| F8 | face 3 | Medium | `internal/adapters/ai_probe_copilot_executor.go`, `internal/adapters/ai_probe_copilot.go` | **deferred to P84**: `--log-dir` on-disk logging disabled, logging-destination arguments refused, refusal evidence limited to the closed classification; **T8a not implemented** — BYOK is not refused when on-disk redaction cannot be proven |
| T9 | range | Medium | `internal/chain/review.go`, `internal/console/ai.go` | **closed**: the evidence-root binding applies only in `policy` mode; a host-auth profile without `baseUrl` is no longer refused by the receiver gate |

## 5. Contract impact

All eight fixes and the T9 narrowing were judged **DOES-NOT-TOUCH**: the violated sentences (`CONTRACTS` :147 / :150 / :212 / :214 / :216 / :391, `SECURITY` :23 / :60) already contain the required semantics. This slice **did not change** `docs/CONTRACTS{,_EN}.md`, any schema or any fixture.

## 6. Triage and user rulings

| Item | Ruling | Record |
|---|---|---|
| (7) budget (1 final + 3 re-reviews) | user ruling of 2026-11-20 | exceeds the section 7.5 literal standard "M-L counts as M"; a direct user ruling that may not be cited as a precedent |
| Blind review | one mandatory round (F3 / F6 touch identity semantics) | counted separately as +1, does not consume the final or re-review budget, runs before the final review |
| T8 half | "refusal evidence uses stable classification only" approved; "refuse BYOK when redaction cannot be proven" not approved and deferred | registered separately as T8a |
| One additional re-review | user ruling of 2026-11-20 (after the 3/3 budget was exhausted) | exceeds section 7.5; a direct user ruling that may not be cited as a precedent |
| T4 / T5#4 | user ruling of 2026-11-21: not fixed, deferred to P84 | F1 / F4 may not be recorded as closed |

## 7. Review chain

| Stage | Carrier | Rounds | Verdict |
|---|---|---|---|
| (5) pre-review | `DeepSeek V4 Pro (deepseek)` | 2 / 2 (exhausted) | r1 FINDINGS (1 Medium) → r2 PASS |
| (6) independent scan | `MAI Code 2 Flash (copilot)` | 2 / 3 | r1 FINDINGS (4 blocking) → r2 PASS |
| (7) blind review | `GPT-6.2 Sol (copilot)` Max | 1 | `BLIND REVIEW: FINDINGS` |
| (7) final review | `GPT-6.2 Sol (copilot)` Max | 1 / 1 | `FINAL REVIEW: FINDINGS` (9 items) |
| (7) re-reviews | `GPT-6.2 Sol (copilot)` Max | 3 / 3 (exhausted) | r1 / r2 / r3 all FINDINGS |
| (7) additional re-review | `GPT-6.2 Sol (copilot)` Max | 1 (user ruling) | `ADDITIONAL RE-REVIEW: FINDINGS` (2 Medium in scope) |

Verbatim transcripts: final review in `tmp/P83-SECURITY-FINDINGS-REMEDIATION/07-final-review.md`; the three re-reviews in `08-reviews-r1-r3.md`; the additional round in `09-additional-review.md` (temporary artifacts, not committed per section 11.2).

## 8. Items deferred to P84

| Id | Deferred content | Source |
|---|---|---|
| P84-a | Policy artifact (an AgentRunner frozen policy persistable for independent re-verification, bound to run / task / step) | T2 STOP (F3) |
| P84-b | Immutable install source (stable Windows path plus ACL, unwritable Linux install directory) | T3 residuals (F6) |
| P84-c | BYOK proof carrier (prove the child cannot write unredacted logs; refuse when unprovable) | T8a (F8) |
| P84-d | bounddir case-folding and intermediate-reparse residuals, root relative-path semantics | T4 residuals (F1) |
| P84-e | Durable proof recovery (predecessor validation in the batch rewrite, a recoverable generation-to-proof transaction) | T5#4 residuals (F4) |

All five are deferred parts of the F1–F8 remediation duty; the user names P84 at its start and rules the priority.

## 9. Documentation sync

| Object | Disposition |
|---|---|
| This report pair | adds `docs/validation/security-findings-remediation{,_EN}.md` (CN / EN strictly mirrored) |
| Ledger | `docs/t027/REMAINING_SLICES{,_EN}.md`, P83 section and completion record (per-F status, R / D, budget disclosure, deferral list) |
| SECURITY note | `docs/SECURITY{,_EN}.md` line 5 states the remediation status truthfully (F2/F5/F7 fixed; F1/F4 partial; F3/F6/F8 deferred) |
| ADR-047 | stays `COMPENSATED`, unchanged |
| Not changed | `CHECKLIST_v66{,_EN}`, `s1-exit`, `p57` gate table, the directive itself, CI, `DEV_PLAN` |

## 10. Gate results

| Check | Result |
|---|---|
| `node tools/gates/gate.js --selftest` | `SELFTEST: PASS 740/740` |
| `--all --skip-base --scope=index` | `TOTAL_FAIL=0` |
| `--all --skip-base --scope=tree` | `TOTAL_FAIL=0` |
| `--all --skip-base --scope=ci` | `TOTAL_FAIL=0` |
| `gofmt -l internal cmd` | no output |
| Encoding | added / changed `.go` are UTF-8 without BOM and LF; `.md` is UTF-8 with BOM and LF |

## 11. (9) validation results and items not executable on this host

| Item | Result |
|---|---|
| `go build ./...` | pass |
| `go vet ./...` | pass |
| `go test ./... -count=1` | all 15 packages with tests pass, no FAIL |
| `GOOS=linux go build ./...` / `go vet ./...` | pass |
| `GOOS=darwin go build ./...` | pass |
| `-race` | **not executed** (no gcc on this host) |

**Items not executable on this host (their acceptance evidence moves to the CI Ubuntu leg)**:

| Case | Reason |
|---|---|
| `TestP83Fix8T4cBoundDirNamesReadsFromDescriptor` | Linux-only (build tag) |
| `TestP83Fix8EVerifiedDescriptorInheritedOnlyByItsChild` | Linux-only (build tag) |
| `TestP83Fix7T3dDescriptorExecNeverResolvesTheReplacedPath` | Linux-only (build tag) |
| `TestP83Fix8T4bReparseArtifactIsRefused` file-symlink subtest | no `SeCreateSymbolicLinkPrivilege` on this host (the junction subtest ran and passed) |
| `TestP83Fix7T5c1RetainedSymlinkIsRefused` | same as above (the Windows link-refusal branch is covered statically and by the junction variant) |
| `TestP83Fix8T4c` / `E` mutation replay | same as above |

### CI execution record

| Stage | Commit | Run | Result |
|---|---|---|---|
| First run (red) | `ec532a4` | `44795140136` | ubuntu `Test` failed 5 cases; windows `Gates` reported `base go test ./...` ok=14 FAIL=3 |
| Fix | `c820743` | — | three test files only: platform-neutral F8 fixtures, F2 end-to-end accepting the earlier admission gate, robust fix8E sibling script |
| Re-run (green) | `c820743` | `44801690679` | both legs `success`; ubuntu ran Test/Race/Contract fixtures/Gates; windows ran Test/Gates (Race and Contract fixtures skipped by design) |
| Package-level evidence | `c820743` | `44801690679` | both legs `PASS  base go test ./... :: ok=15 FAIL=0`, `SELFTEST: PASS 740/740`, `TOTAL_FAIL=0` |

VM local (`10.0.0.199`, go1.22.12): targeted `internal/adapters` / `internal/console` / `internal/guard` all `ok`; full suite 15 packages all `ok`; the five cases all `--- PASS` in verbose mode (no SKIP); targeted `-race` all `ok`.

## 12. Environment and tiers

| Role | Carrier and tier | Set at | Paused |
|---|---|---|---|
| Controller | `DeepSeek V4.1 Flash (deepseek)` High | start of the slice | no |
| (1) architect | `deep-reasoner` + `DeepSeek V4 Pro (deepseek)` Max | plan and triage rounds | no |
| (2) implementer | `prfrail-implementer` + `DeepSeek V4.1 Flash (deepseek)` | throughout | no |
| (3) tester | `prfrail-tester` + `DeepSeek V4.1 Flash (deepseek)` | throughout | no |
| (5) pre-review | `DeepSeek V4 Pro (deepseek)` | pre-review ring | no |
| (6) scan | `MAI Code 2 Flash (copilot)` | scan ring | no |
| (7) blind / final / re-review | `GPT-6.2 Sol (copilot)` Max | tier taken 2026-11-20 | no |

## 13. Review accounting

| Item | Value | Note |
|---|---|---|
| R | 2 | one (5)-ring remediation re-run plus one (6)-ring remediation re-run; re-review rounds do not count per section 9.2 |
| D | 0 | no temporary-file discipline violation in this slice |
| Probes | 0 | no real external call; the real BYOK path was never invoked |
| (7) budget | blind 1 + final 1 + re-review 3 + additional 1 | the additional round is a direct user ruling that may not be cited as a precedent; both the report and the ledger must disclose it |
| (5) / (6) caveat | **their PASS verdicts bind the pre-remediation snapshot** | the (5) 2/2 and (6) 2/3 PASS verdicts predate several remediation rounds; the final snapshot is covered only by the (7) re-reviews, which remained FINDINGS |
| Anomaly E-1 | `docs/t027/CHECKLIST_v66.md` showed an unstaged working-tree change | the file is on the "do not touch" list; the origin could not be attributed, so it was reverted with `git checkout --` per the user ruling, restoring working tree equals index (same file set); not counted as D |

## 14. Explicitly not executed

No gitee action; no directive version bump; no new OB / ADR; no CI or `DEV_PLAN` change; no real BYOK probe; no `-race`; the Linux-only cases were not executed on a real Linux host (left to the CI Ubuntu leg).

## 15. (8)a self-check

All five deliverables are present (code and tests, this report pair, the ledger section and completion record, the SECURITY note, the evidence pack); every finding has a falsifiable assertion and a counterexample; no test was weakened or deleted; no contract surface changed beyond added internal APIs; the verdict line agrees with the per-finding status and never says "closed" or "security-closed".

## 16. Six-stage self-check (sixstage)

```sixstage
## Section A 总体结论
PARTIALLY CLOSED

## Section B 发现清单

| SEVERITY | Location | Residual |
|---|---|---|
| High | F3 | T2 STOP: no frozen policy artifact for AgentRunner |
| High | F6 | T3: Windows name-to-image binding window, Linux same-inode not byte-immutable |
| Medium | F8 | T8a: BYOK refusal not implemented |
| Medium | F1 | T4: bounddir case-folding and intermediate-reparse residuals |
| Medium | F4 | T5#4: batch predecessor validation, cross-process interruption recovery |

## Section C 逐条证据与推理
- The guards and counterexamples for F2 / F5 / F7 and the T9 narrowing are the table in section 4.
- All six (7) rounds are archived verbatim in the evidence pack under 04-reviews/ and 01-arch/.
- The (5) and (6) PASS verdicts bind the pre-remediation snapshot; the final snapshot is covered only by (7).

## Section D 可证伪性审计
- Delete any finding's main guard and its counterexample test turns red.
- Rewrite the verdict as "closed" or "security-closed" and the section 4 versus section 15 consistency check turns red.
- Remove the T2 STOP or T8a deferral record and the section 8 versus ledger completion record cross-check turns red.
- UNCOVERED: the Linux-only cases and the symlink-privilege case were not executed on this host; the real BYOK path was never invoked.


## Section E 我无法验证的部分
The Windows name-to-launch-image binding window, the Linux same-inode byte mutability and the child's own on-disk logging cannot be proven statically on this host; they are deferred to P84.

INDEPENDENT SCAN: FINDINGS
```

## 17. Artifact fence

```artifacts
docs/SECURITY.md	repo
docs/SECURITY_EN.md	repo
docs/t027/REMAINING_SLICES.md	repo
docs/t027/REMAINING_SLICES_EN.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/00-briefing.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/01-arch/remediation-plan.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/01-arch/ruling-triage.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/07-additional-review.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/07-final-review.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/04-reviews/07-re-reviews-r1-r3.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/MANIFEST.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/SHA256SUMS.txt	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/freedom-list.md	repo
docs/validation/evidence/P83-SECURITY-FINDINGS-REMEDIATION/input-list.md	repo
docs/validation/security-findings-remediation.md	repo
docs/validation/security-findings-remediation_EN.md	repo
internal/adapters/agent_runner_admission.go	repo
internal/adapters/agent_runner_replay_dispatcher.go	repo
internal/adapters/ai_preflight.go	repo
internal/adapters/ai_probe_copilot.go	repo
internal/adapters/ai_probe_copilot_executor.go	repo
internal/adapters/ai_probe_copilot_test.go	repo
internal/adapters/bounddir.go	repo
internal/adapters/bounddir_alias.go	repo
internal/adapters/p83_fixround3_tester_test.go	repo
internal/adapters/p83_fixround4_f8_tester_test.go	repo
internal/adapters/p83_fixround5_tester_test.go	repo
internal/adapters/p83_fixround6_tester_test.go	repo
internal/adapters/p83_fixround7_tester_test.go	repo
internal/adapters/p83_fixround8_tester_test.go	repo
internal/chain/engine_test.go	repo
internal/chain/p83_t9a_tester_test.go	repo
internal/chain/review.go	repo
internal/chain/review_test.go	repo
internal/chain/state_seal_counterexample_test.go	repo
internal/console/ai.go	repo
internal/console/ai_test.go	repo
internal/console/approvals.go	repo
internal/console/cli.go	repo
internal/console/cost.go	repo
internal/console/p83_f2_tester_test.go	repo
internal/console/p83_f4_fixround_test.go	repo
internal/console/p83_f6_tester_test.go	repo
internal/console/p83_fixround3_tester_test.go	repo
internal/console/p83_fixround4_tester_test.go	repo
internal/console/p83_fixround5_tester_test.go	repo
internal/console/p83_fixround5_windows_test.go	repo
internal/console/p83_fixround6_tester_test.go	repo
internal/console/p83_fixround7_tester_test.go	repo
internal/console/p83_fixround8_tester_test.go	repo
internal/console/p83_security_findings_test.go	repo
internal/console/runtime_agent.go	repo
internal/console/runtime_agent_e2e_test.go	repo
internal/console/runtime_multitask_agent_test.go	repo
internal/console/runtime_multitask_driver_test.go	repo
internal/console/runtime_noop.go	repo
internal/console/runtime_supervised.go	repo
internal/console/secret_input_windows.go	repo
internal/console/trusted_ledger.go	repo
internal/evidence/bounddir.go	repo
internal/evidence/bounddir_clean.go	repo
internal/evidence/bounddir_darwin.go	repo
internal/evidence/bounddir_escape_payload_unix_test.go	repo
internal/evidence/bounddir_escape_redbaseline_test.go	repo
internal/evidence/bounddir_escape_unix_test.go	repo
internal/evidence/bounddir_failclosed_unix_test.go	repo
internal/evidence/bounddir_fifo_linux_test.go	repo
internal/evidence/bounddir_identity_test.go	repo
internal/evidence/bounddir_symlink_storm_test.go	repo
internal/evidence/bounddir_unix.go	repo
internal/evidence/bounddir_windows.go	repo
internal/evidence/boundread.go	repo
internal/evidence/fileid_unix.go	repo
internal/evidence/fileid_windows.go	repo
internal/evidence/p83_fixround8_linux_test.go	repo
internal/evidence/p83_fixround8_windows_test.go	repo
internal/gates/adapter.go	repo
internal/gates/p83_f1_f5_tester_test.go	repo
internal/gates/p83_f5_fixround_test.go	repo
internal/gates/p83_fixround5_tester_test.go	repo
internal/gates/p83_fixround6_tester_test.go	repo
internal/gates/p83_fixround7_tester_test.go	repo
internal/gates/policy_test.go	repo
internal/gates/probe_permit_counterexample_test.go	repo
internal/gates/runner.go	repo
internal/gates/runner_artifact_bound_test.go	repo
internal/gates/runner_effect_policy_test.go	repo
internal/gates/runner_test.go	repo
internal/guard/executable_binding.go	repo
internal/guard/executable_binding_error.go	repo
internal/guard/lease.go	repo
internal/guard/lease_open_unix.go	repo
internal/guard/lease_open_windows.go	repo
internal/guard/lease_test.go	repo
internal/guard/p83_fixround5_tester_test.go	repo
internal/guard/p83_fixround5_windows_test.go	repo
internal/guard/p83_fixround6_tester_test.go	repo
internal/guard/p83_fixround6_windows_test.go	repo
internal/guard/p83_fixround7_linux_test.go	repo
internal/guard/p83_fixround7_tester_test.go	repo
internal/guard/p83_fixround7_windows_test.go	repo
internal/guard/p83_fixround8_linux_test.go	repo
internal/guard/p83_fixround8_windows_test.go	repo
internal/guard/process.go	repo
internal/guard/process_env_block_windows_test.go	repo
internal/guard/process_linux.go	repo
internal/guard/process_unsupported.go	repo
internal/guard/process_windows.go	repo
```
