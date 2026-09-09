# S1 Exit Validation Report

[中文](s1-exit.md)

Date: 2026-09-09. Verdict: `BLOCKED`. T018 passes the frozen read-only whois comparison and minimal generic/C/Go loops required by AT-15, and a candidate lifecycle drill for trusted Windows artifacts also passes. The product owner approved manual portable-ZIP extraction without PATH modification, with bilingual installation, support-matrix, and release-note drafts now present. The final ZIP is not yet bound and rerun, the EOL notification channel remains open, and the product owner has not approved final S1 exit, so partial success must not be reported as S1 completion.

## AT-15 Evidence

1. This session received read-only authorization limited to non-secret, redacted fixtures under `D:\LZProjects\whois\testdata\`. No active start files, logs, `out/`, or `release/` content was read; whois was neither modified nor executed.
2. `examples/whois-shadow/` exports nine fixed cases from `testdata/cidr_matrix_cases_draft.tsv`. Source SHA-256 is `8a6abecb0178a970677b84b3e5403ce8830f281e1a197d792ae54089b3e9efb0`; the source has 1,981 bytes, 10 lines, and zero secret-marker matches.
3. A stdlib-only comparator aligns cases by ID and compares input, result, and failure classification. All 9/9 positive cases pass; separate input, result, and failure-classification mutations fail closed.
4. `examples/{generic-minimal,c-minimal,go-minimal}/` declare `generic-standard`, `c-standard`, and `go-standard` language scopes. All three CLI `validate → preview → run → report` loops write only to temporary run directories, reach `COMPLETED`, make zero command/network/model/credential/version-probe calls, and leave config bytes unchanged.
5. Boundary: the whois evidence is an offline projection of an authorized frozen fixture, not a live query, production execution, or cutover. The language examples exercise the current noop-only CLI and do not claim executable-hook or compiler execution.

## AT-01–AT-22 Summary

| AT | Status | Primary evidence |
|---|---|---|
| AT-01 | PASS | `DEV_PLAN_EN.md` T002/T003: 29 schemas, 82 fixtures, JCS vectors |
| AT-02 | PASS | [T006 snapshot](t006-snapshot_EN.md) |
| AT-03 | PASS | [T007 process/lease](t007-process-lease_EN.md), [T009 chain](t009-chain_EN.md) |
| AT-04 | PASS | [T005 evidence](t005-evidence_EN.md), [T008 applier](t008-applier_EN.md) |
| AT-05 | PASS | [T011 review/publish](t011-chain_EN.md) |
| AT-06 | PASS | [T010 gates](t010-gates_EN.md) |
| AT-07 | PASS | [T012 tickets/repair](t012-tickets-repair_EN.md) |
| AT-08 | PASS | [T013 adapters](t013-adapters_EN.md) |
| AT-09 | PASS | [T014 handoff](t014-handoff_EN.md) |
| AT-10 | PASS | [T014 handoff](t014-handoff_EN.md) |
| AT-11 | PASS | [T015 harness/docs](t015-taskdef-documentation_EN.md) |
| AT-12 | PASS | [T016 CLI](t016-cli-console_EN.md) |
| AT-13 | PASS | [T017 self-host](s1-selfhost_EN.md), including Node 24 run `34325758548` |
| AT-14 | PASS | [T017 self-host](s1-selfhost_EN.md), eight Windows/Linux jobs and six artifacts |
| AT-15 | PASS | `examples/whois-shadow/`, all three `*-minimal` language examples, and focused checks in this report |
| AT-16 | PASS | [T019 preview](t019-preview_EN.md) |
| AT-17 | PASS | [T020 export](t020-export_EN.md) |
| AT-18 | PASS | [T021 approvals](t021-approvals_EN.md) |
| AT-19 | PASS | [T022 effects/diagnostics](t022-effects-diagnostics_EN.md) |
| AT-20 | PASS | [T023 cost ledger](t023-cost-ledger_EN.md) |
| AT-21 | PASS | [T024 lifecycle](t024-lifecycle_EN.md) |
| AT-22 | BLOCKED | T025 is not implemented: `operator-interaction` Schema/goldens, embedded SessionBridge IPC, TUI interaction inbox, and recovery tests are absent |

These PASS verdicts cover the existing deterministic AT evidence; they do not automatically satisfy stage documentation and release gates.

## RFC Section 14 S1 Exit

| Gate | Status | Evidence/gap |
|---|---|---|
| Core task chain and fail-closed scenarios | BLOCKED | Existing AT-01–AT-14 and AT-16–AT-21 pass; new T025/AT-22 is not implemented |
| Seed self-host and real platforms | PASS | T017 Windows/Linux trust matrix and native Windows two-generation drill |
| Read-only whois shadow | PASS | 9/9 fixed cases and three drift counterexamples; no production cutover |
| generic/C/Go harnesses | PASS | All three harness definitions and minimal user examples pass the read-only CLI loop |
| Formal install-package tutorial | PARTIAL | The [portable ZIP installation guide](../INSTALLATION_EN.md) and model are frozen; the final ZIP is not yet bound and rerun |
| Install/upgrade/rollback/uninstall exercises | PARTIAL | The [Windows candidate portable-package drill](s1-install-drill_EN.md) passes; the final carrier still requires rerun |
| Support matrix and release notes | PARTIAL | The [candidate support matrix](../S1_SUPPORT_MATRIX_EN.md) and [release-note draft](../S1_RELEASE_NOTES_EN.md) exist but are not bound to a formal release |
| Independent review | PASS | Both T017 baseline and Node 24 replacement received Codex/DeepSeek reviews |
| Final product-owner S1 approval | BLOCKED | This session's fixture-read authorization is not release or stage approval |

## Residual Risks and Next Steps

1. The manually extracted portable-ZIP model, user-selected independent directory, explicit-path invocation, and no-PATH policy are frozen. The formal ZIP filename and download entry point still require release binding. ADR-010 now proposes the EOL notice and exception policy, but product-owner approval remains outstanding.
2. Windows candidate artifacts now pass first-run, upgrade, rollback, and evidence-preserving uninstall; the final carrier must rerun the drill. Linux remains core CI only and is not S1 production support.
3. The generic/C/Go examples prove language scopes and noop orchestration only. Executable CLI steps currently fail closed.
4. Complete T025 and pass AT-22 to prove the TUI operator-interaction loop without `@sbr-review`; then rerun the tutorial against the final ZIP, bind the support matrix/release notes, and obtain explicit product-owner S1 exit approval. Until then, T018 and S1 remain incomplete.