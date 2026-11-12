# P69 ⑥ 独立扫描 r1 归档（逐字原文）

载体：`runSubagent(agentName='independent-reviewer', model='MAI-Code-1.1-Flash (copilot)')`（只读；盲态，未附任何他人清单）。结论行：**INDEPENDENT SCAN: PASS**（Section B 无发现）。以下为扫描最终消息逐字转录。

## Section A

INDEPENDENT SCAN: PASS

## Section B

(no findings)

## Section C

- The contract is matched exactly by the level-1 gate semantics in [internal/gates/readiness.go](internal/gates/readiness.go#L141-L170): `EvaluateT027BReadiness` is a pure function, and for `GoalProbe` it checks presence only, while the actual quota consumption is delegated to `ProbePermit.Require`. That matches the contract sentence: “EvaluateT027BReadiness is the only construction entry and is a pure function; zero I/O; no permit consumption.”
- The one-shot lifecycle is enforced in [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L1-L95) and [internal/gates/probe_permit.go](internal/gates/probe_permit.go#L176-L244): the state machine is `Unused → Claimed → Spent`, the linearization point is the `O_CREATE|O_EXCL` marker file persisted with `Sync`, and `Complete` is idempotent only after a valid claim. This is the direct contract anchor for the “fail-closed” rule.
- The production wiring sits in the allowed scope only: [internal/console/readiness_gate.go](internal/console/readiness_gate.go#L1-L24), [internal/console/cli.go](internal/console/cli.go#L282-L395), [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L120-L190), and [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go#L259-L320). The flow is consistent with the design: CLI resolves and validates the permit before routing, AgentRunner evaluates before preconditions and consumes only after preconditions pass, and supervised starts enforce a fresh-start ledger gate before reusing the permit.
- The negative-path test bundle in [internal/console/readiness_gate_test.go](internal/console/readiness_gate_test.go#L34-L205) is the real behavioral proof for the five acceptance conditions: nil permit denied, spent permit denied, missing CLI flag denied, incomplete record set kept fail-closed, and successful assembly consumes the permit exactly once.
- The mutation evidence in [docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt](docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt) is especially important: M1–M5 are all red when the guard is removed or weakened, then restored green after the original state is restored. That is exactly the correct “can this be falsified?” check.
- The recorded gate readouts are all green in the evidence bundle: [tmp/P69-T027B-E2E-WIRING/06-gates/p4b-07-selftest.txt](tmp/P69-T027B-E2E-WIRING/06-gates/p4b-07-selftest.txt), [tmp/P69-T027B-E2E-WIRING/06-gates/p4b-08-index.txt](tmp/P69-T027B-E2E-WIRING/06-gates/p4b-08-index.txt), [tmp/P69-T027B-E2E-WIRING/06-gates/p4b-09-tree.txt](tmp/P69-T027B-E2E-WIRING/06-gates/p4b-09-tree.txt), and [tmp/P69-T027B-E2E-WIRING/06-gates/p4b-10-ci.txt](tmp/P69-T027B-E2E-WIRING/06-gates/p4b-10-ci.txt). The index/tree/ci reports all say `TOTAL_FAIL=0`, and the selftest report is `737/737`.

## Section D

- If the CLI presence gate in [internal/console/cli.go](internal/console/cli.go#L295-L395) is removed:
  - The check in [internal/console/readiness_gate_test.go](internal/console/readiness_gate_test.go#L167-L205) goes red: `TestRunRefusesMissingProbePermitFlag` fails because the missing `--probe-permit` path becomes a later record-read or assembly error instead of a fail-closed refusal.
  - The “which check is not red?” answer is that the existing `--request-record` and “local runtime only supports noop steps” validations still remain red; the unique thing that becomes green wrongly is the missing-permit gate itself.

- If the `Require` claim is removed from the shared helper in [internal/console/readiness_gate.go](internal/console/readiness_gate.go#L8-L24) or from the assembly entry in [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L146-L180):
  - `TestAgentRunnerAssemblyRefusesSpentProbePermit` in [internal/console/readiness_gate_test.go](internal/console/readiness_gate_test.go#L105-L166) goes red because presence-only evaluation would allow an already-spent permit to proceed.
  - The “not red” part is the initial `EvaluateT027BReadiness` presence check; that one is intentionally not enough to decide quota consumption, so it is not the guard that catches a spent permit.

- If the deferred `Complete` release is removed:
  - The mutation case M5 in [docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt](docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt) goes red; the second run would report `ErrProbePermitInUse` instead of `ErrProbePermitConsumed`.
  - This is the exact “complete must be recorded” invariant; no other check catches it.

- If the supervised fresh-start gate in [internal/console/runtime_supervised.go](internal/console/runtime_supervised.go#L269-L320) is removed:
  - A second supervised run on an existing interaction ledger can wrongly pass the permit through again; the re-entry behavior becomes fail-open.
  - The current evidence does not show a dedicated second-park mutation on a real ledger; the invariant is still logically covered by the ledger check and the design, but it is the one place where an OS-level concurrency test would be the strongest proof.

- If the stronger “reject before marker creation” ordering in [internal/console/runtime_agent.go](internal/console/runtime_agent.go#L146-L180) is removed:
  - The nil/spent permit tests in [internal/console/readiness_gate_test.go](internal/console/readiness_gate_test.go#L66-L166) are the ones that go red, because the run-root marker could be created before the permit failout.
  - The design is intentionally ordered to prevent a marker from appearing on a fail-closed refusal.

- The no-coverage invariant I would call out is the real cross-process race itself: the code relies on `O_CREATE|O_EXCL` plus `Sync`, but the artifacts do not contain a true concurrent two-process reproduction of the race. That is not a logic bug in this slice; it is the one remaining externality that would require a live two-process OS test to eliminate all doubt.

## Section E

- The things that are not directly observable from the read-only evidence:
  - A real two-process race on the same permit marker across separate OS processes.
  - A live Windows-specific agent-runner publication-durability path that depends on the actual pinned CLI and replay store behavior.
  - A full end-to-end run of the real CLI in a production environment outside the recorded artifacts.
  - The precise behavior of the Linux-side `GOOS=linux` path beyond the recorded pass states; the output files are empty, but the higher-level gate reports still confirm overall pass.

- What would overturn this PASS:
  - Any proof that a `nil` or already-spent permit is accepted into a real run.
  - Any proof that a second supervised entry with an existing ledger bypasses the fresh-start gate.
  - Any proof that the permit is consumed without `Require`, or that `Complete` is skipped on the success path.
  - Any diff outside the allowed scope that modifies the gate logic itself under [internal/gates](internal/gates) or alters the chain/adapters semantics.
  - Any real runtime failure where the mutation tests in [docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt](docs/validation/evidence/P69-T027B-E2E-WIRING/06-gates/mutation-evidence.txt) no longer hold.

- The current evidence is strong enough to say the slice is fail-closed in the intended sense and stays within the declared contract. The remaining uncertainty is only the untested OS-level race behavior, not a logical defect in the gate ordering or the permit lifecycle.
