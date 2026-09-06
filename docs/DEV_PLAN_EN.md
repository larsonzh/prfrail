# ProofRail Iterative Development Plan

[简体中文](DEV_PLAN.md)

Date: 2026-09-06; all implementation tasks remain incomplete. Inputs: [requirements](PRODUCT_REQUIREMENTS_EN.md), [architecture](ARCHITECTURE_EN.md), [contracts](CONTRACTS_EN.md), [security](SECURITY_EN.md), [ADRs](ADR_REGISTER_EN.md). This follows revised RFC sections 14/19; added S1 minimum scope requires renewed budget approval, without pulling S2/S3 enhancements into the first release.

## 1. Workflow and Constraints

Go 1.22 module compatibility, pure Go, CLI/TUI and content-addressed files; dependencies freeze in S0. No framework migration applies, so Java/.NET transformation rules are inapplicable. Read-only source, single writer, first-fail-stop, independent approval and no default Git writes have no architectural exceptions.

Handle one task per session. Split anything exceeding one verifiable change into subtests instead of generating every module. Red test, minimal implementation, same test, package tests, full gates, bilingual docs/evidence. Offline deterministic checks precede authorized model calls. Estimate dates after measuring the first two tasks, not from unmeasured assumptions.

WIP=1. Agree time/cost caps per session; no paid calls by default. Two same-class failures without new evidence pause/escalate. Reserve stronger models for contract/security ambiguity; cheaper models execute frozen small slices. Git commit/push still require same-turn authorization, origin only by default.

## 2. Iterations and Tasks

Paths not yet present are planned creation locations. T001-T004 are S0; T004 is isolated experimentation, not S1 production code. Every RFC section 17 gate and owner approval must precede T005.

### P0.1 Documents and Approval (REQ-001/014/017/023/026/028)

- [x] T001 [Plan:P0.1] Review docs/PRODUCT_REQUIREMENTS*.md, docs/SECURITY*.md and docs/ADR_REGISTER*.md; record owner/security reviewer, budgets, support/signing and shadow scope. Depends: this package. Accept: section 17.1/6/8 sign-off and license decisions. Exclude: approving for the user. Completed 2026-09-06: owner approved S0 specification freeze and ADR section 4; independent security conclusions remain a separate T003/gate 17.3 closure.

T001 status: `COMPLETE`. This approves only S0 specification freeze and allows T002 to begin. It does not authorize new paid calls, concrete whois data access, Git operations, release or S1 production coding. Signer/key, exact dependencies and independent security review remain later gates.

### P0.2 Machine Contract Freeze (REQ-001/004/005/007/008/010/012/013/017/018/020/021/022/024/025/028)

- [ ] T002 [Plan:P0.2] Revise docs/RFC-proofrail-unattended-ai-engineering-product.md first, synchronize docs/CONTRACTS*.md/ADRs, create schemas/*.schema.json; freeze ADR-002/004 fields, versions, canonical vectors, transitions, errors and queue. Depends: T001. Accept: AT-01 structure/independent review. Exclude: production code using guessed fields.
- [ ] T003 [Plan:P0.2] Create testdata/contracts/{valid,invalid}/ and tools/contracts/ independent validator/semantic fixture oracle: three tasks, four kinds, dual/triple language, handoff, docs/waiver and all receipts/manifests. Depends: T002. Accept: AT-01 and section 17.2/4/9-12. Exclude: candidate-generated expected results.

T002 status: `IN_PROGRESS`. Common version, ID, path, JCS and UTF-8 rules are frozen, with ten Draft 2020-12
structural schemas for chain, hook, target, workspace, state-event, error, adapter-envelope, adapter-receipt,
snapshot-manifest and evidence-manifest. Initial canonical vectors, single-record transitions, error categories/CLI exits,
file-queue framing/claim/result/dispatch/takeover receipt rules, snapshot boundaries and the pre-review evidence root are frozen;
cross-event/cross-record checkers, the concrete error catalog/primary ordering, signature
receipts, remaining persistent objects and ADR-004 capability probes are still pending. T003's
independent validator/semantic checker, vector fixtures and positive/negative goldens have not run, so T002, AT-01 and RFC
section 17.2 must not be marked complete.

T001 also reviews ADR-008 and PC-01–PC-06 workload/budgets. T002 freezes RFC section 19 authorization/preview/export/effect/cost/backup records; T003 adds independent positive/negative goldens for gate 17.13. This freezes contracts, not S1 runtime acceptance.

### P0.3 Technical Spikes (REQ-006/007/009/015/017/021/023/024/027)

- [ ] T004 [Plan:P0.3] Probe Go/TUI, Windows atomic replacement/process-tree stop, Linux paths and IPC/file queue in root tmp; record docs/validation/s0-spikes.md and delete experiments. Depends: T002; S1 also requires all T003/T001 gates. Accept: reproducible commands, actual platforms and explicit unavailable capabilities. Exclude: unauthorized paid calls or production claims.

### P1.1 Evidence and Snapshots (REQ-004/023/024/025/028)

- [ ] T005 [Plan:P1.1] Implement internal/evidence/{event,receipt,verify}.go and matching _test.go for canonical/event chains/references/read-only verification. Depends: all S0 gates. Accept: AT-04. Exclude: free-text PASS or unknown-version writes.
- [ ] T006 [Plan:P1.1] Implement internal/snapshot/{capture,restore,store}.go/tests for paths/exclusions/quotas/retention/GC dry-run. Depends: T005. Accept: AT-02. Exclude: Git-history recovery or referenced-object deletion.

### P1.2 Writes and Ordered Loop (REQ-002/003/004/012/018/023/024/026)

- [ ] T007 [Plan:P1.2] Implement internal/guard/{process,lease}.go, platform files/tests for tree termination, takeover and leases. Depends: T005/T004. Accept: AT-03/06. Exclude: killing user processes or PID/expiry-only takeover.
- [ ] T008 [Plan:P1.2] Implement internal/taskdef/checker.go and internal/applier/{apply,journal,recover}.go/tests: sequential memory validation, transaction and rollback. Depends: T005/T006/T007. Accept: AT-04 at every crash point. Exclude: validate-while-writing or omitted assertions.
- [ ] T009 [Plan:P1.2] Implement internal/chain/{engine,state,recover}.go/tests: three tasks, four kinds, both change modes, projection/pause/cancel, initially fake agent/runner. Depends: T008. Accept: AT-03. Exclude: candidates as downstream parents.
- [ ] T010 [Plan:P1.2] Implement internal/gates/{runner,policy,result}.go/tests and wire chain. Depends: T007/T009. Accept: AT-06. Exclude: weakening unavailable network/resource enforcement.

### P1.3 Independent Acceptance and Bounded Repair (REQ-005/014/022/023/024/025/026)

- [ ] T011 [Plan:P1.3] Implement internal/chain/{review,publish}.go/tests, bound candidate/review/receipt, rejection and valid waiver. Depends: T010. Accept: AT-05. Exclude: agent self-approval or technical-success PASSED.
- [ ] T012 [Plan:P1.3] Implement internal/tickets/{ledger,budget,fingerprint}.go and internal/repair/transaction.go/tests. Depends: T011. Accept: AT-07. Exclude: budget resets or direct formal-definition repair.

### P1.4 Agents and Handoff (REQ-006/007/008/016/021/023/024/025/026)

- [ ] T013 [Plan:P1.4] Implement internal/adapters/{filequeue,sessbridge}.go/tests with consumer-owned ports. Depends: T012. Accept: AT-08 offline then authorized host. Exclude: SessionBridge edits, GUI/auto fallback or repeated application.
- [ ] T014 [Plan:P1.4] Implement internal/chain/handoff.go and internal/console/handoff.go/tests. Depends: T013. Accept: AT-09/10. Exclude: model-mediated secrets, expiry-based success or concurrent automatic writers.

### P1.5 Harnesses and Usable CLI/TUI (REQ-009/010/011/013/015/016/018/020/022/027)

- [ ] T015 [Plan:P1.5] Create harnesses/{generic,c,go}/, internal/taskdef/documentation.go/tests for template A, impact rules, freshness and C+Go tasks. Depends: T014. Accept: AT-11. Exclude: whois core terminology or arbitrary script B.
- [ ] T016 [Plan:P1.5] Complete cmd/prfrail/main.go, internal/console/, config parsing/explanation/tests; sync README.md and docs/OPERATIONS*.md. Depends: T015. Accept: AT-12. Exclude: direct view-state writes or successful placeholder commands.

### P1.6 Trusted Release and Shadow Acceptance (REQ-001/002/006/011/014/015/019/024/025/028)

- [ ] T017 [Plan:P1.6] Create .github/workflows/ci.yml, testdata/selfhost/ and docs/validation/s1-selfhost.md; pin tooling, isolate generations, external oracle, native/signature checks. Depends: T024 and bootstrap approval. Accept: AT-13/14. Exclude: seed overwrite or unauthorized release.
- [ ] T018 [Plan:P1.6] Create examples/{whois-shadow,go-minimal}/ and docs/validation/s1-exit.md; compare fixed read-only input/result/failure classifications and summarize all AT/risks. Depends: T017 and whois input authorization. Accept: AT-15 and every RFC S1 exit gate. Exclude: production whois cutover or partial-pass completion.

### P1.7 Product Flow and Lifecycle (Before P1.6 Release Acceptance)

Keep IDs stable; execution order is T016 → T019–T024 → T017 → T018. All require S0 approval. Paths are planned; split tasks by individual negative tests rather than generating whole features at once.

- [ ] T019 [Plan:P1.7] PC-01: implement internal/taskdef/preview.go, internal/console/preview.go, tests and no-AI example. Depends: T016. Accept: AT-16. Exclude: preview commands/network/models/automatic installation.
- [ ] T020 [Plan:P1.7] PC-02: implement internal/snapshot/export.go, internal/evidence/delivery.go/tests. Depends: T019, T006/T011. Accept: AT-17. Exclude: candidate-as-accepted export, target overwrite, source/Git writes.
- [ ] T021 [Plan:P1.7] PC-03: implement internal/chain/authorization.go, internal/console/approvals.go/tests, wire guard stopping. Depends: T020, T007/T011. Accept: AT-18. Exclude: self-authorization, hard-gate waiver, expiry-based approval.
- [ ] T022 [Plan:P1.7] PC-04: implement internal/gates/effects.go, internal/evidence/diagnostics.go and chain recovery-plan tests. Depends: T021, T010. Accept: AT-19. Exclude: external-write runners, blind unknown-effect retries, diagnostic lock/log deletion.
- [ ] T023 [Plan:P1.7] PC-05: extend internal/tickets/budget.go, internal/adapters/ usage and internal/console/ cost reporting. Depends: T022, T012/T013. Accept: AT-20. Exclude: restart-released unknown holds, cross-run cap bypass, default telemetry.
- [ ] T024 [Plan:P1.7] PC-06: implement internal/snapshot/backup.go, internal/evidence/disposition.go, lifecycle controls/tests and release inventory tooling. Depends: T023. Accept: AT-21, SBOM/support documentation. Exclude: unknown-format migration, default evidence/shared-tool deletion, unauthorized signing/release.

## 3. Requirement Mapping

Evidence is planned, not implemented.

| REQ ID | Plan Items | Tasks | Acceptance/implementation evidence |
|---|---|---|---|
| REQ-001 | P0.1,P0.2,P1.6 | T001,T002,T003,T018 | AT-01/15; generic schema/non-C example |
| REQ-002 | P1.2,P1.6 | T009,T018 | AT-03/15; chain/report |
| REQ-003 | P1.2 | T009 | AT-03; ordered scheduling |
| REQ-004 | P0.2,P1.1,P1.2 | T002,T006,T008 | AT-02/04; snapshot/applier |
| REQ-005 | P0.2,P1.3 | T002,T011 | AT-05; review/publish |
| REQ-006 | P0.3,P1.4,P1.6 | T004,T013,T017 | AT-08/14; adapters/CI |
| REQ-007 | P0.2,P0.3,P1.4 | T002,T004,T013 | AT-08; two transports |
| REQ-008 | P0.2,P1.4 | T002,T013 | AT-08; context recovery |
| REQ-009 | P0.3,P1.5 | T004,T016 | AT-12; TUI |
| REQ-010 | P0.2,P1.5 | T002,T016 | AT-01/12; config explanation |
| REQ-011 | P1.5,P1.6 | T015,T018 | AT-11/15; three harnesses |
| REQ-012 | P0.2,P1.2 | T002,T007,T010 | AT-06; gate runner |
| REQ-013 | P0.2,P1.5 | T002,T015 | AT-11; template generation |
| REQ-014 | P0.1,P1.3,P1.6 | T001,T012,T018 | AT-07/15; stage report |
| REQ-015 | P0.3,P1.5,P1.6 | T004,T016,T017 | AT-12/13/14; build/TUI |
| REQ-016 | P1.4,P1.5 | T013,T016 | AT-08/12; no-IDE CLI |
| REQ-017 | P0.1,P0.2,P0.3 | T001,T002,T003,T004 | AT-01; readiness |
| REQ-018 | P0.2,P1.2,P1.5 | T003,T009,T016 | AT-01/03; four kinds |
| REQ-019 | P1.6 | T017 | AT-13/14; two generations |
| REQ-020 | P0.2,P1.5 | T003,T015 | AT-01/11; components/dual language |
| REQ-021 | P0.2,P0.3,P1.4 | T003,T004,T014 | AT-09/10; handoff |
| REQ-022 | P0.2,P1.3,P1.5 | T003,T011,T015 | AT-01/05/11; coordinated docs |
| REQ-023 | P0.1,P0.3,P1.1,P1.2,P1.3,P1.4 | T001,T004,T005,T007,T010,T011,T012,T014 | AT-05/06/07/09; security |
| REQ-024 | P0.2,P0.3,P1.1,P1.2,P1.3,P1.4,P1.6 | T002,T004,T005,T006,T008,T009,T011,T014,T017 | AT-02/03/04/09/13; recovery |
| REQ-025 | P0.2,P1.1,P1.3,P1.4,P1.6 | T002,T005,T011,T013,T014,T017 | AT-04/06/10/14; evidence |
| REQ-026 | P0.1,P1.2,P1.3,P1.4 | T001,T009,T010,T012,T013 | AT-03/06/07; budgets |
| REQ-027 | P0.3,P1.5 | T004,T016 | AT-12; terminal/JSON |
| REQ-028 | P0.1,P0.2,P1.1,P1.6 | T001,T002,T003,T005,T017 | AT-01/04/14; compatibility |

## 4. Testing and Definition of Done

Take the union of section 3 and this table. Checkpoint product-slice supplements likewise add P0.1/P0.2/P1.7 without new requirement IDs.

| PC | REQ | Plan Items | Tasks | Acceptance |
|---|---|---|---|---|
| PC-01 | REQ-009/010/016/027 | P0.1,P0.2,P1.7 | T001,T002,T003,T019 | AT-16 |
| PC-02 | REQ-002/004/025 | P0.1,P0.2,P1.7 | T001,T002,T003,T020 | AT-17 |
| PC-03 | REQ-005/021/023 | P0.1,P0.2,P1.7 | T001,T002,T003,T021 | AT-18 |
| PC-04 | REQ-012/023/024 | P0.1,P0.2,P1.7 | T001,T002,T003,T022 | AT-19 |
| PC-05 | REQ-006/026 | P0.1,P0.2,P1.7 | T001,T002,T003,T023 | AT-20 |
| PC-06 | REQ-019/025/028 | P0.1,P0.2,P1.7 | T001,T002,T003,T024 | AT-21 |

See [TEST_STRATEGY_EN.md](TEST_STRATEGY_EN.md): Go testing, real local files/processes, independent schema validation, CLI oracle and authorized host smoke; no database/browser dependency. Package tests join go test ./... and report pass/fail/skip counts. Mocks do not prove native Windows process/filesystem or actual adapter capability.

A task needs implementation/tests, requirement/contract links, bilingual affected docs, commands/exit/evidence and no scope violations. Iterations require build/vet/test plus relevant AT. S1 requires all mandatory tests and human sign-off. Plan coverage is not implementation coverage.

## 5. Later Stages and Resources

Detail S2 after S1 exit: production Linux, composed dependency ordering/three-language execution, supervised/Web/optional VS Code views, rule-based model selection and intelligent documentation hints. S3 contains script B, more editors and inception. Preserve requirement IDs, extend tasks/AT through RFC/ADR, and never silently move required S1 capabilities out.

Each iteration reports actual human time, calls/cost estimates, failures, accepted requirements and next action. Budget overruns pause for owner choice or written RFC scope change, not automatic models/services.

RFC section 19 later candidates: S2 graphical onboarding/terminal localization, redacted notifications/remote approval identity design, external PR evaluation, external-write compensation, cost trends and explicit storage migration; S3 commercial backend/cloud management requires demand evidence. These are outside T019–T024 and will be planned later. Current totals: 10 plan items, 24 pending tasks, 21 planned acceptance groups, not completion percentages.