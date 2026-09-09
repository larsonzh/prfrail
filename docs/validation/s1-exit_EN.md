# S1 Exit Validation Report

[中文](s1-exit.md)

Date: 2026-09-09. Verdict: `BLOCKED`. T018 has started, and both the frozen read-only whois comparison and the minimal non-C Go loop required by AT-15 pass. The formal install-package tutorial, runnable generic/C user examples, install/upgrade/rollback/uninstall exercises, and final product-owner S1 approval remain incomplete, so partial success must not be reported as S1 completion.

## AT-15 Evidence

1. This session received read-only authorization limited to non-secret, redacted fixtures under `D:\LZProjects\whois\testdata\`. No active start files, logs, `out/`, or `release/` content was read; whois was neither modified nor executed.
2. `examples/whois-shadow/` exports nine fixed cases from `testdata/cidr_matrix_cases_draft.tsv`. Source SHA-256 is `8a6abecb0178a970677b84b3e5403ce8830f281e1a197d792ae54089b3e9efb0`; the source has 1,981 bytes, 10 lines, and zero secret-marker matches.
3. A stdlib-only comparator aligns cases by ID and compares input, result, and failure classification. All 9/9 positive cases pass; separate input, result, and failure-classification mutations fail closed.
4. `examples/go-minimal/` declares a standalone `go-standard` language scope. Its CLI `validate → preview → run → report` loop writes only to a temporary run directory, reaches `COMPLETED`, makes zero command/network/model/credential/version-probe calls, and leaves the config bytes unchanged.
5. Boundary: the whois evidence is an offline projection of an authorized frozen fixture, not a live query, production execution, or cutover. The Go example exercises the current noop-only CLI and does not claim executable-hook support.

## AT-01–AT-21 Summary

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
| AT-15 | PASS | `examples/whois-shadow/`, `examples/go-minimal/`, and focused checks in this report |
| AT-16 | PASS | [T019 preview](t019-preview_EN.md) |
| AT-17 | PASS | [T020 export](t020-export_EN.md) |
| AT-18 | PASS | [T021 approvals](t021-approvals_EN.md) |
| AT-19 | PASS | [T022 effects/diagnostics](t022-effects-diagnostics_EN.md) |
| AT-20 | PASS | [T023 cost ledger](t023-cost-ledger_EN.md) |
| AT-21 | PASS | [T024 lifecycle](t024-lifecycle_EN.md) |

These PASS verdicts cover the existing deterministic AT evidence; they do not automatically satisfy stage documentation and release gates.

## RFC Section 14 S1 Exit

| Gate | Status | Evidence/gap |
|---|---|---|
| Core task chain and fail-closed scenarios | PASS | AT-01–AT-14 and AT-16–AT-21 reports |
| Seed self-host and real platforms | PASS | T017 Windows/Linux trust matrix and native Windows two-generation drill |
| Read-only whois shadow | PASS | 9/9 fixed cases and three drift counterexamples; no production cutover |
| generic/C/Go harnesses | PARTIAL | Harness definitions and a Go user example exist; runnable generic/C user examples do not |
| Formal install-package tutorial | BLOCKED | `INSTALLATION_PLAN_EN.md` remains explicitly a plan; no formal release package exists |
| Install/upgrade/rollback/uninstall exercises | BLOCKED | T024 provides library-level lifecycle behavior but no exercised installer entry point |
| Support matrix and release notes | PARTIAL | Platform/version policy is frozen but no final documents are bound to a formal release |
| Independent review | PASS | Both T017 baseline and Node 24 replacement received Codex/DeepSeek reviews |
| Final product-owner S1 approval | BLOCKED | This session's fixture-read authorization is not release or stage approval |

## Residual Risks and Next Steps

1. The release carrier, install location, and PATH policy remain open; do not invent install commands or publish an unauthorized artifact.
2. The final artifact requires native Windows install, first-run, upgrade, rollback, and uninstall exercises. Linux remains core CI only and is not S1 production support.
3. Runnable generic/C user examples must use isolated temporary directories and prove source-tree immutability. Executable CLI steps currently fail closed.
4. After those artifacts and the final support matrix/release notes exist, the product owner must explicitly approve S1 exit. Until then, T018 and S1 remain incomplete.