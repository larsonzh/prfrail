# ProofRail Test and Validation Strategy

[简体中文](TEST_STRATEGY.md)

Date: 2026-09-07; test and acceptance baseline, not a passing test report. This document is the normative authority for ATs, test tiers, platforms, and fault-injection acceptance rules. Inputs: [requirements](PRODUCT_REQUIREMENTS_EN.md), [contracts](CONTRACTS_EN.md), and [security](SECURITY_EN.md). [Project proposal](RFC-proofrail-unattended-ai-engineering-product.md) sections 14 and 17 retain stage-acceptance provenance.

## 1. Baseline and Tiers

Application: local CLI/TUI with file/process orchestration. T005–T008 now have core-package _test.go coverage, but no usable-CLI product E2E exists; package tests do not replace product acceptance. Canonical test entry is go test ./...; new Go tests must be discoverable there. Build/static gates are go build ./... and go vet ./.... Observed Go is 1.27.0; the module declares 1.22, which still needs separate minimum-toolchain verification.

| Tier | Tools/data | Evidence | When |
|---|---|---|---|
| Documents | Offline links, BOM/LF, bilingual IDs, REQ-task-AT checks | Navigability and scope coverage | Every doc change |
| Schema | Independent JSON Schema validator plus semantic fixture oracle | Correct-layer positive/negative decisions | S0 freeze and protocol changes |
| Unit | Go testing, tables, fake clock/agent/runner | Transitions, budgets, checker, hashes | Each small task |
| Component integration | t.TempDir, real files, Go subprocess helper | Journal, paths, stopping, recovery, leases | File/process changes |
| CLI E2E | Built candidate, temporary workspace/store, fixed oracle | Three-task chains and US1-US6 | Iteration/S1 exit |
| Host integration | Installed SessionBridge v0.1.1, controlled instance | silent capability, routing, restart, complete response | After offline contracts; paid calls need authorization |
| Native platform | Windows 11 amd64; Linux amd64 core CI | Filesystem/signal/path behavior | Every candidate |
| Self-hosting | Seed N, candidate N+1, external oracle | Two-generation trust and clean-room replay | S1 exit/release |

No database, Docker service or browser is required for S0/S1; do not install them. The infrastructure tier uses real local files/processes, not in-memory mocks as durability proof. Browser tier is inapplicable; use terminal probes for TUI and introduce browser E2E only with S2 Web. Windows race detection may need a C toolchain: run on capable CI and disclose local skips, distinct from CGO_ENABLED=0 releases.

## 2. Acceptance Matrix

AT IDs are planned test groups, not existing functions. Expand each into positive/negative cases and report counts, PASS/FAIL/SKIP, versions and evidence paths.

| ID | Requirements | Scenario/assertion | Owner/plan |
|---|---|---|---|
| AT-01 | REQ-001/010/017/018/020/022/028 | Valid schemas; reject unknown fields/versions, duplicate IDs, empty steps, invalid noop, overlapping targets and dependency/generation cycles; dual/triple-language, handoff and doc fixtures | taskdef; T002/T003 |
| AT-02 | REQ-004/024 | Uncommitted source capture remains unchanged; reserved/case/symlink/reparse/hardlink/long-path/concurrent-change/quota cases | snapshot; T006 |
| AT-03 | REQ-002/003/018/024/026 | Three mixed-step tasks, zero-process noop, T2 failure blocks T3; pause/cancel/restart, runaway children and disk pressure | chain/guard; T007/T009 |
| AT-04 | REQ-004/024/025/028 | Missing/reordered events, corrupt objects, torn tails, crashes at journal/receipt/projection writes; only completed acceptance is replayable | evidence/applier; T005/T008 |
| AT-05 | REQ-005/022/023 | Reject self-review, rejection, expired/wrong-candidate waiver; valid independent approval publishes once, code/docs together | chain; T011 |
| AT-06 | REQ-012/023/025/026 | argv injection, environment secrets, escaped cwd, unavailable network enforcement, timeout/unlimited output/scan failure/missing artifact | gates/guard; T007/T010 |
| AT-07 | REQ-014/023/026 | Fingerprint/attempt/time/cost exhaustion, fake repairs, stale candidates, interrupted Promote; budgets cannot reset | tickets/repair; T012 |
| AT-08 | REQ-006/007/008/016 | IPC/file-queue parity, busy/wrong ID/stale results/timeouts/restarts/cache loss/truncated history; no GUI fallback or repeated application | adapters; T013 |
| AT-09 | REQ-021/023/024 | Successful return, escaped edits, lease conflict, absence/disconnection, abort/request-agent and crash; gates rerun on return | chain/guard; T014 |
| AT-10 | REQ-021/025 | Undeclared prompt stops; terminal-only secret-direct; unavailable safe input pauses; canaries absent from evidence/models | gates/console; T014 |
| AT-11 | REQ-011/013/020/022 | generic/C/Go, locked template A; self-approval/manual generated edits/missing docs block; C+Go whole-task recovery | taskdef/gates; T015 |
| AT-12 | REQ-009/010/015/016/027 | No-IDE init/validate/run/report/config explain; keyboard/no-color/narrow/--json; unimplemented commands cannot report completion | console/cmd; T016 |
| AT-13 | REQ-015/019/024 | Separate seed/candidate, external oracle, clean-room success/failure preserves seed | harness/release; T017 |
| AT-14 | REQ-006/015/019/025/028 | Native Windows/Linux core CI, minimum/fixed release Go, historical read-only/unknown-version rejection, fixed commit plus SHA256SUMS/SBOM/license verification; checksums never claim publisher identity | release; T017 |
| AT-15 | REQ-001/002/011/014 | Frozen read-only whois shadow input/result/failure parity; small non-C Go loop, no production-tree changes | harness; T018 |
| AT-16 | REQ-009/010/016/027 | PC-01: no-AI preview; zero malicious-hook/version-probe/credential/network calls, visible unknowns, unchanged source/store, narrow/JSON fact parity | taskdef/console; T019 |
| AT-17 | REQ-002/004/025 | PC-02: offline accepted-package verification; reject candidates/secrets/missing objects/tampering/existing targets/alias overlap; interrupted writes lack completion receipt, source unchanged | snapshot/evidence; T020 |
| AT-18 | REQ-005/021/023 | PC-03: revocation/expiry/wrong hash block next effects/acceptance; failed in-flight stopping stays paused; no hard-gate waiver, approval inbox survives restart | chain/guard/console; T021 |
| AT-19 | REQ-012/023/024 | PC-04: reject external writes/unenforceable restrictions, no unknown-effect redispatch; read-only redacted diagnostics/recovery plans without lock deletion/log mutation | gates/chain/evidence; T022 |
| AT-20 | REQ-006/026 | PC-05: crashes before/after reservation, dispatch and settlement; timeout holds, deduplicated settlement, two runs cannot overallocate shared cap; no precise subscription-bill claim | tickets/adapters; T023 |
| AT-21 | REQ-019/025/028 | PC-06: missing closure objects/active writers/unknown formats block; new-store restore preserves originals; uninstall retains evidence/shared tools, deletion/retention conflict, offline revocation unknown | snapshot/evidence/release; T024 |

AT-16–AT-21 are RFC section 19 S1 minimum acceptance, using deterministic fixtures without paid calls first. S0 T002/T003 freeze records/fixtures only; parse success is not runtime evidence. T018 summarizes all 21 groups. US7–US10 cover trial, delivery, revocation and retirement recovery; local measurements establish value without assumed savings.

S2/S3 remain governed by RFC section 14: real three-language integration, supervised disconnection, Web, additional adapters and generated script B. S0 config parsing is not runtime acceptance.

## 3. Fixtures, Faults and Evidence

Use small fixed non-secret fixtures with stable IDs/content/expected hashes, isolated by random run IDs and t.TempDir. Never use the current repository for destructive tests. Redacted read-only whois exports require permission; absent input means SKIP/BLOCKED, not invented equivalence.

Inject before/after temporary writes, rename, journal durability, partial target writes, review, publication receipts, events before projections, GC reference updates and operator return. On restart assert unchanged source, preserved accepted references, no repeated application and pause on uncertainty. Terminate helper subprocesses, never the IDE or user processes.

Reports contain toolchain/OS/dependency versions, input hashes, command/cwd, exit codes, case counts, failure evidence, skips and executed-versus-planned distinction. t.TempDir cleans itself; manual temporary artifacts belong in root tmp and are removed afterward. Planned docs/validation/ stores reviewed acceptance summaries, not raw secret-bearing logs.

## 4. Gates and Debugging

Run the affected package first, then build/vet/test. Protocol changes also need independent validation/compatibility goldens; process/path changes need native platforms; adapter changes need offline contracts then host tests. Missing evidence, mandatory security failure or critical skips prevent stage completion.

Debug in order: reproduce input, locate first failing layer, inspect error/evidence hash, isolate package/case, make one small repair, rerun that case, run all gates. Two same-class failures without new evidence stop model calls and escalate. Never widen permissions, reset budgets or disable gates for green results.

Use no-AI fixtures before a single authorized host smoke. Do not ask multiple models to reanalyze the same passing result. Performance follows the fixed dataset in the requirements; no measured conclusion exists yet.