# ⑥ 独立扫描结论原文（逐字归档）

载体：runSubagent(agentName="independent-reviewer", model="MAI-Code-1.1-Flash (copilot)")。只读审查，未修改任何文件。以下为原文逐字（含其 ```text 围栏形态，主报告中的 sixstage 转写见主报告 §⑥）。

```text
READ ONLY: no files modified; only staged diff and source/docs inspection were used.
## Section A
Scope: staged review covers the 22-file package for this slice; under internal/ the staged diff is exactly 7 added counterexample test files, with no non-_test.go changes. The audit focused on whether the new tests truly exercise the claimed fail-closed paths, whether skip logic is honest, and whether the assertions are tied to the real contract rather than incidental errors.

## Section B
Finding: no material defect in the slice as a fail-closed acceptance-security counterexample package.
- The skip branches are not blanket passes; they are conditional on environment limits or a known vulnerability backlog.
- The Linux escape case is intentionally skipped before the vulnerable call and is tracked as a known gap, while the red-path evidence is kept in the separate attack record rather than hidden.
- The Windows junction tests skip only when the environment cannot create a junction; they are not silently passing when the platform supports the primitive.
- The child-process helper skips only when the environment variable is absent; that is harness gating, not an evasion.
- The non-sentinel assertions are anchored to the documented contract: either the exact boundary wording is required, or the error must not be one of the known sentinels. These are not loose text matches on arbitrary failures.
- The cases are written as falsification checks: remove the corresponding guard and the test becomes red.

## Section C
Evidence reviewed:
1. Staged diff confirms the internal scope: 7 new test files and no production changes under internal/.
2. The skip is placed before the vulnerable call in the Unix escape case, so it does not automatically pass the attack path; it exits before entering the known vulnerable branch. The separate attack record documents the actual exploit outcome: ATTACK-SUCCEEDED with SECRET-OUTSIDE-BOUNDARY.
3. The sentinel-heavy tests bind directly to concrete sentinels via errors.Is, matching the contracts for closed handles, over-bound growth, invalid regular files, missing permits, incomplete readiness, and incomplete seals.
4. The non-sentinel tests are not arbitrary; they assert the exact documented behavior for a missing sentinel or a boundary message, such as the directory-identity mismatch and the outside-the-boundary junction rejection.
5. The replay test is not mock-only: it persists a real seal to the file-backed store and then drives the engine replay path that compares the stored context key to the derived one, so the failure is tied to the real acceptance logic and not a fake object.
6. The environment gates are narrow: second-writer failure, junction creation unavailable, helper env absent, and known vulnerability backlog. None of them are broad CI exemptions.

## Section D
Falsifiability audit:
- The suite is shaped to reject the underlying defense: if any corresponding guard is removed, the relevant case should fail.
- The sentinel assertions are precise and one-to-one with the contract, which is the strongest anti-false-positive pattern.
- The documented asymmetries (empty context key and trimmed AT23) are explicitly noted as current behavior, not hidden defects. They are record-keeping at the contract boundary rather than silent acceptance.
- The only deliberate exception is the known Linux middle-symlink gap, which is acknowledged as a tracked fix item and is not presented as a passing check. This is a legitimate gap marker, not a fake pass.

## Section E
Unverified or non-fully-validated parts:
- The Linux middle-symlink exploit is not run in the normal unit suite because it is intentionally skipped. It is only proven in the separate attack record, so the CI path is not directly executing the exploit.
- The Windows junction branch remains platform-dependent: if mklink/junction creation is unavailable, the test is skipped and no boundary proof is available on that runner.
- The empty-context and trim-space behaviors are currently accepted as asymmetries; they are part of the audit but not part of the fail-closed guard set itself.
Reverse conditions:
- If any skip becomes unconditional, or if the vulnerable path executes without the documented precondition, the scan would flip to findings.
- If a production file under internal/ changes in the staged diff, or if a sentinel assertion can be satisfied by an unrelated error, the scan would also flip to findings.
- If the exploit can be reproduced without the skip gate in a normal test run, the T51-1 case is not valid as a tracked-vulnerability record and the package would fail the audit.

INDEPENDENT SCAN: PASS
```
