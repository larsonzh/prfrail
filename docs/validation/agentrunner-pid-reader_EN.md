# Slice `AGENTRUNNER-PID-READER` verification report (Directive v1.42)

> Slice: `AGENTRUNNER-PID-READER`; class `[SLICE]`, tier M-L; directive v1.42 (this slice adds one line each to §12.4 and §0.3 and bumps the header version).
> Baseline: `d:\LZProjects\prfrail`; HEAD = `7201679` (C1 closeout).
> Bilingual: this file is the CN authority; `docs/validation/agentrunner-pid-reader_EN.md` is the line-for-line mirror.

## 1 Change summary

In one sentence: this slice closes two outstanding debts — OB-110 (reader-side three-state semantics) and OB-109 (the sixth frozen fact) — and logs OB-111 under directive v1.42.

Deliverables:

- Contract A5/A6: `docs/CONTRACTS.md` and `docs/CONTRACTS_EN.md`, +2/−2 lines each (`git diff --cached --numstat` reading).
- Code: 18 paths — 10 under `internal/adapters`, 4 under `internal/chain`, 4 under `tools/b4-e2e`.
- Focused tests: the new and rewritten cases in those files, plus two named runs, `a6native` and `b3bnative`.
- Directive, three places: the v1.42 entry of §0.3, the OB-111 entry of §12.4, and the header version.
- This report pair (CN authority plus `_EN` mirror) and `docs/validation/evidence/AGENTRUNNER-PID-READER-freedom-list.md`.

**Position of v1.42**: this slice is an execution slice; v1.42 only logs one row inside the slice and bumps the version, so it is a by-product of the slice, not an independent governance slice, and it changes no criterion semantics or criterion scope.

## 2 Environment and toolchain

- Local toolchain: `go1.27.0 windows/amd64`.
- CI declaration: Go `1.22.12`.
- Repository workspace: Windows (`d:\LZProjects\prfrail`).
- This slice was not re-run on a Linux VM: its changes carry no platform branch, unlike the cross-platform re-run of C1.
- Command surface: `gofmt -l .`, `go build ./...`, `go vet ./...`, `go test -count=1 ./...`.

## 3 Deviations and fallbacks

**D1 (the only deviation in this slice)**: ②a created `tmp/verify-*.js` and `tmp/pid-reader-dump.md` at the root of `tmp/`, violating the closed location set of §5.1 for temporary artifacts. Handling: ②a self-reported and self-deleted them, then verified zero residue (top level of `tmp/` holds only `.gitkeep` and `AGENTRUNNER-PID-READER/`). It is logged and the run continues per §5.1 (log explicitly, do not dissolve, add no further action); this is the 1st occurrence in this slice and the 11th under the §5.1 aggregate count.

## 4 Review chain

| Round | Model | Verdict | Handling |
|---|---|---|---|
| ⑤ r1 | `GLM-5.3` | 1 Medium (**double BOM**) + 2 Low + 5 Info | fixed, then ⑤ r2 |
| ⑤ r2 | `GLM-5.3` (narrowed read) | PASS | converged |
| ⑥ r1 | `MAI-Code-1.1-Flash` | PASS (0 findings) | no rework |
| ⑦ r1 | `GPT-5.3-Codex` (blind review) | `ADJUDICATION: PASS` | converged |

The two Low findings of ⑤ r1 are the inaccurate test comment and the harness `readPIDs` still parsing loosely with `fmt.Sscanf`; they are logged in §8 together with the 5 Info findings.

## 5 OB-110 reader-side three-state semantics

The three states follow contract A6: a file that does **not exist** is an evidence gap (fail-closed ⇒ `uncertain`); an existing file that is **empty** (0 bytes) is the fact "there really was no child process"; an existing **non-empty** file is parsed line by line as the pid list. The three states must never be inferred from each other.

The four readers, one by one:

| Reader | Three-state handling | Verdict when missing |
|---|---|---|
| `internal/adapters/agent_runner_run.go`: `readAgentRunnerChildPIDs` with a triple return | missing ⇒ `present=false`; empty ⇒ empty list; non-empty ⇒ per-line `strconv.Atoi` | `agentRunnerVerifyTreeGoneWith` sees `present=false` and returns "not verified" immediately |
| `tools/b4-e2e/ports.go`: `readPIDs` plus `managedProcessesGone` | same three states | collects the missing file names and reports `adapters.ErrInvalidAgentRunnerEvidence` |
| `tools/b4-e2e/postflight.go`: `alivePIDs` | same three states | returns an error, the caller degrades to `uncertain` |
| `tools/b4-e2e/runmap.go`: `mirrorVerifyTreeGone` | mirrors the production semantics | `present=false` ⇒ `treeGone=false` |

**Relation to `uncertain`**: a missing report reaches `uncertain` only through "the process tree cannot be proven gone"; an empty report is not a gap. **Priority is unchanged**: a natural non-zero exit is still `failed`, and `uncertain` never overrides an established natural failure.

## 6 OB-109 the sixth frozen fact

- Byte caliber (quoting contract A6): `ProcessPidReportsHash = evidence.Digest("proofrail:agent-runner-pid-reports:1\n", "stub.pid " + Digest(parent report bytes) + "\nstub-child.pid " + Digest(child report bytes) + "\n")`; the inner `Digest` uses the empty domain over raw bytes and never decodes.
- Implementation site: `AgentRunnerPidReportsDigest` (`internal/adapters/agent_runner_evidence.go:369`); when either report is missing the collector writes only a gap-domain digest (`:357`) and never an empty digest.
- Chain side: `postflightFactsCount = 6` (`internal/chain/postflight.go:39`); the six facts are rebuilt from `evidence[2..7]` contiguously and without holes (`:40`, `:249`).
- Leg 7 three states: a recomputed value diverging from the frozen value ⇒ `failed`; a missing report ⇒ `uncertain`; an empty report ⇒ recomputed and frozen values are both the empty-string digest ⇒ `pass`.
- `PidReportsComplete` is an independent boolean (`internal/adapters/agent_runner_evidence.go:79`); when false, `FrozenFacts()` refuses with `ErrAgentRunnerEvidenceGap` (`:163`).

## 7 Tests and gates

③ readings: `gofmt` produced empty output; `go build ./...` and `go vet ./...` both exited 0; `go test -count=1 ./...` reported 15 packages `ok` and 0 `FAIL`; the named runs `a6native` and `b3bnative` were each `ok`.

④ readings: `tree=0`, `index=0`, `ci=0`; the `tools/gates` self-test reported `SELFTEST 588/588`; the encoding audit covered 24 paths with 0 violations.

The encoding audit had found a **double BOM** (the Medium of ⑤ r1); it was fixed and re-run successfully.

**CI readings (main commit `02d5797`, run `39740097139`)**: both legs finished `success` (the Ubuntu leg took 2m46s); package-level `ok` lines: windows 15 / ubuntu 17; in-leg `SELFTEST` 588/588 and `TOTAL_FAIL` 0 on both legs. The first run was green, so the §6.4 evidence-based rerun was not triggered and no DR entry was registered.

## 8 Known limitations

- **KL-C1-7 restated**: reusing `run-c1-faces.sh` requires a brand-new `C1_OUT`; keeping an old directory mixes the previous round's evidence into this round's reconciliation.
- **KL-C1-8 verdict**: the four failure conditions for offline operation were **not triggered** in this slice. The basis is the mechanical check by ③/①: this slice added no `net` dependency, no socket call, no cgo, and no new external process.
- **Low (mirror layer)**: `readPIDs` in `tools/b4-e2e/ports.go:525` still parses with `fmt.Sscanf`, which is not equivalent to the production `strconv.Atoi` (`internal/adapters/agent_runner_run.go:327`); it affects the harness only, never the production path.
- **Low**: the seed comment at `internal/adapters/agent_runner_postflight_test.go:515-520` contradicts the constant-write caliber.
- **Info (defense in depth)**: the missing branch of leg 7 is intercepted earlier by leg 5, so it is unreachable in practice.
- **Info**: the harness read has no upper bound check; a TOCTOU window exists between `Stat` and `ReadFile`.
- **Info**: the stale comments at `internal/console/runtime_agent_e2e_test.go:457-460` and `:643-646` (claiming the stub writes the child report only when it spawns a child) were not introduced by this slice.
- **Info**: five-fact-era routing evidence is rejected as a whole with zero writes on replay (fail-closed).
- **G4a gap (observed in this slice)**: G4a verifies only that a BOM exists, not how many BOMs exist, so a double BOM passes; a follow-up `gates` slice will handle it.
- **Nature of the four `G6-4` rewrites made during ⑨ (registered as-is, no verdict change)**: the first ④ run on `--scope=index` failed on four hits, rewritten as follows - `18` with a `个文件` suffix to "18 paths", `v1.42` with a `行` suffix to "the v1.42 entry", `OB-111` with a `行` suffix to "the OB-111 entry", `24` with a `文件` suffix to "24 paths". **Judged a false-positive shape**: such number-plus-unit statements describe facts of this slice change set, which do not expire once the report is frozen, whereas the `G6-4` pattern does not distinguish slice change-set counts from repository current totals. **Standing rule** stated in item ②: for the same situation prefer the built-in `G6-4` whitelist (explicit exemption, reason, context excerpt) and **do not default to rewording**; rewording is a record-only disposition only when the criterion is a confirmed false positive and the whitelist path is unavailable, and it must be registered as such. This is a **same-family data point for OB-111 (the `G6-4` variant)** for the next OB-111 evaluation.

## 9 Degrees of freedom

The degrees of freedom exercised by this slice, the boundaries fixed by rulings and the single deviation are logged one by one in `docs/validation/evidence/AGENTRUNNER-PID-READER-freedom-list.md` (F 10, M 6, D 1).

## 10 Items not executed

- No Linux VM re-run: this slice has no platform-related change.
- No new face: the harness has no evidence-dir tamper knob, and adding a face would widen the scope.
- OB-104 / OB-105 / OB-107 are not handled here.

## 11 Conclusion

Both outstanding debts are closed item by item: the three-state semantics of OB-110 now agree across the four readers, and the sixth fact of OB-109 is byte-for-byte equivalent to contract A6; all three gate scopes are green (`tree`, `index` and `ci` are each 0); the only deviation D1 is logged per §5.1.

```artifacts
docs/validation/agentrunner-pid-reader.md
docs/validation/agentrunner-pid-reader_EN.md
docs/validation/evidence/AGENTRUNNER-PID-READER-freedom-list.md
docs/CONTRACTS.md
docs/CONTRACTS_EN.md
docs/DELIVERY_DIRECTIVE.md
docs/DELIVERY_DIRECTIVE_EN.md
docs/t027/REMAINING_SLICES.md
docs/t027/REMAINING_SLICES_EN.md
internal/adapters/agent_runner_evidence.go
internal/adapters/agent_runner_evidence_test.go
internal/adapters/agent_runner_pinned_cli_test.go
internal/adapters/agent_runner_postflight.go
internal/adapters/agent_runner_postflight_test.go
internal/adapters/agent_runner_run.go
internal/adapters/agent_runner_run_test.go
internal/adapters/agent_runner_terminal_chain.go
internal/adapters/agent_runner_terminal_publisher_test.go
internal/adapters/agent_runner_terminal_record.go
internal/chain/agent_runner_terminal.go
internal/chain/agent_runner_terminal_test.go
internal/chain/postflight.go
internal/chain/postflight_test.go
tools/b4-e2e/main.go
tools/b4-e2e/ports.go
tools/b4-e2e/postflight.go
tools/b4-e2e/runmap.go
```

DOCS: DONE
